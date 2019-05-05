package apibuilder

import (
	"errors"
	"gopkg.in/yaml.v3"
	"strings"
)

const (
	RETURN_TYPE_NOCONTENT = "nocontent"
	RETURN_TYPE_JSON = "json"
	RETURN_TYPE_FILE = "file"
)

var (
	ErrNoDocVersion = errors.New("no doc version")
	ErrNoBaseUrl = errors.New("no baseUrl")
	ErrNoApis = errors.New("no apis")
	ErrInvalidDataType = errors.New("invalid data type")
	ErrNoMemberInDataType = errors.New("no member in data type")
	ErrNoDataTypeName = errors.New("no data type name")
	ErrInvalidDataTypeName = errors.New("invalid data type name")
	ErrDuplicateDataType = errors.New("duplicate data type")
	ErrInvalidApiDef = errors.New("invalid api definition")
	ErrApiNoUrl = errors.New("api has no url")
	ErrApiNoReturn = errors.New("api no return")
	ErrApiUrlNotString = errors.New("api url is not string")
	ErrReturnTypeUnsupported = errors.New("return type unsupported")
	ErrInvalidMember = errors.New("invalid member")
	ErrInvalidMemberAttr = errors.New("invalid member attributes")
	ErrInvalidReturnDef = errors.New("invalid return definition")
	ErrInvalidParamsDef = errors.New("invalid param definition")
	ErrInvalidApiMethod = errors.New("invalid api method")
	ErrDataTypeNotExist = errors.New("data type is not exist")
)

// API字段成员
type Members map[string]*MemberAttr

func NewMember() Members {
	return make(Members)
}

// API数据复合类型
type DataType struct {
	doc     *ApiDoc
	Name    string  // 数据类型
	Members Members // 类型字段描述
}

// 添加API类型
func NewDataType(doc *ApiDoc, name string) *DataType {
	return &DataType{doc: doc, Name: name, Members: NewMember()}
}

func NewBaseDataType(name string) *DataType {
	return &DataType{Name: name}
}

// 向API类型中增加成员
func (t *DataType) AddMember(name string, attr *MemberAttr) {
	t.Members[name] = attr
}

func (t *DataType) CheckValue(typeValue interface{}) (success bool, err error) {
	return
}

type MemberEachFunc func(name string, attr *MemberAttr)

// 遍历数据类型中的成员
func (t *DataType) MembersEach(eachFunc MemberEachFunc) {
	for name, attr := range t.Members {
		eachFunc(name, attr)
	}
}

type AttrLength struct {
	Checked bool
	Value int
}

// API字段成员属性
type MemberAttr struct {
	Type        string // 字段数据基本类型
	IsArray     bool   // 是否为数组
	Required    bool   // 是否为必须字段
	Description string // 字段描述
	Length      AttrLength	// 字段长度
	MinLength   AttrLength    // 字段最小长度
	MaxLength   AttrLength    // 字段长度
}

func (ma *MemberAttr) load(attrs map[interface{}]interface{}) (err error) {
	t, hasType := attrs["type"]
	if !hasType {
		err = ErrNoDataTypeName
		return
	}

	switch t.(type) {
	case string:
		typeName := t.(string)
		ma.Type = typeName
		ma.IsArray = false
	case []string:
		ma.IsArray = true
		ma.Type = t.([]string)[0]
	}

	requiredVal, hasRequired := attrs["required"]
	if !hasRequired {
		ma.Required = false
	} else {
		ma.Required, _ = ToBool(requiredVal)
	}

	descriptionVal, hasDescription := attrs["description"]
	if hasDescription {
		ma.Description = ""
	} else {
		ma.Description, _ = ToString(descriptionVal)
	}

	lengthVal, hasLength := attrs["length"]
	if !hasLength {
		ma.Length.Value = 0
		ma.Length.Checked = false
	} else {
		ma.Length.Checked = true
		val, _ := ToInt(lengthVal)
		ma.Length.Value = int(val)
	}

	minLengthVal, hasMinLength := attrs["minLength"]
	if !hasMinLength {
		ma.MinLength.Value = 0
		ma.MinLength.Checked = false
	} else {
		ma.MinLength.Checked = true
		val, _ := ToInt(minLengthVal)
		ma.MinLength.Value = int(val)
	}

	maxLengthVal, hasMaxLength := attrs["maxLength"]
	if !hasMaxLength {
		ma.MaxLength.Value = 0
		ma.MaxLength.Checked = false
	} else {
		ma.MaxLength.Checked = true
		val, _ := ToInt(maxLengthVal)
		ma.MaxLength.Value = int(val)
	}

	return nil
}

// API返回值
type ApiReturn struct {
	ReturnType string		// 返回类型
	Data interface{} // 返回为JSON对象
}

type ApiParam struct {
	Members
	From string
}

// API入口
type ApiEntry struct {
	Url         string             // API接口路径
	Method      string             // API方法
	Params      []*ApiParam // API参数映射
	Returns     map[string]*ApiReturn  // API返回值
	Description string             // API描述
}

// API描述文档
type ApiDoc struct {
	Description    string               // API描述
	Version string               // API版本（语义化版本）：major.minor.revision
	BaseUrl string               // 基本URL
	Apis    []*ApiEntry          // API入口
	Types   map[string]*DataType // API中引用的数据类型定义
}

func NewApiDoc() *ApiDoc {
	return &ApiDoc{
		Types: make(map[string]*DataType),
		Apis: make([]*ApiEntry, 0),
	}
}

func (doc *ApiDoc) getDataType(name string) *DataType{
	dt, exists := doc.Types[name]
	if exists {
		return dt
	}
	return nil
}

func (doc *ApiDoc) addDataType(dataType *DataType) (err error) {
	_, exists := doc.Types[dataType.Name]
	if exists {
		err = ErrDuplicateDataType
		return
	}
	doc.Types[dataType.Name] = dataType

	return
}

func (doc *ApiDoc) addApiEntry(entry *ApiEntry) {
	if entry == nil {
		// ignore nil
		return
	}

	doc.Apis = append(doc.Apis, entry)
}

func (doc *ApiDoc) Parse(content []byte) (err error) {
	yamlDoc := make(map[string]interface{})
	err = yaml.Unmarshal(content, &yamlDoc)
	if err != nil {
		return
	}
	doc.parseBaseInfo(yamlDoc)

	dataTypes, hasDataTypes := yamlDoc["types"]
	if hasDataTypes {
		if err = doc.parseDataTypes(dataTypes.([]interface{})); err != nil {
			return
		}
	}

	apis, hasApis := yamlDoc["apis"]
	if !hasApis {
		err = ErrNoApis
		return
	}

	err = doc.parseApis(apis.([]interface{}))

	return
}


func (doc *ApiDoc) parseBaseInfo(d map[string]interface{}) (err error) {
	// 文档描述
	descVal, _ := d["description"]
	doc.Description, _ = ToString(descVal)

	// 文档版本
	versionVal, hasVersion := d["version"]
	if !hasVersion {
		err = ErrNoDocVersion
		return
	}
	doc.Version, _ = ToString(versionVal)

	// 接口baseUrl
	baseUrlVal, hasBaseUrl := d["baseUrl"]
	if !hasBaseUrl {
		err = ErrNoBaseUrl
		return
	}
	doc.BaseUrl, _ = ToString(baseUrlVal)

	return
}

func parseDataTypeMembers(dataType *DataType, members map[interface{}]interface{}) (err error) {
	for name, attrs := range members {
		mn := name.(string)
		memberAttr := &MemberAttr{}
		a, ok := attrs.(map[interface{}]interface{})
		if !ok {
			err = ErrInvalidMemberAttr
		}
		err = memberAttr.load(a)
		if err != nil {
			return
		}
		dataType.AddMember(mn, memberAttr)
	}
	return
}

func (doc *ApiDoc) parseDataTypes(dataTypes []interface{}) (err error) {
	for _, dt := range dataTypes {
		dataTypeMeta, ok := dt.(map[interface{}]interface{})
		if !ok {
			err = ErrInvalidDataType
			return
		}
		nameVal, hasName := dataTypeMeta["name"]
		if !hasName {
			err = ErrNoDataTypeName
			return
		}
		name, nameErr := ToString(nameVal)
		if nameErr != nil {
			err = ErrInvalidDataTypeName
			return
		}
		dataType := NewDataType(doc, name)
		members, hasMembers := dataTypeMeta["members"]
		if !hasMembers {
			err = ErrNoMemberInDataType
			return
		}
		dstMembers, ok := members.(map[interface{}]interface{})
		if !ok {
			err = ErrInvalidMember
			return
		}
		err = parseDataTypeMembers(dataType, dstMembers)
		if err != nil {
			err = ErrInvalidMember
			return
		}
		if err = doc.addDataType(dataType); err != nil {
			return
		}
	}
	return
}

func parseApiParamDetail(param *ApiParam, members map[interface{}]interface{}) (err error) {
	param.Members = make(Members)
	for memberName, memberAttr := range members {
		mn, _ := memberName.(string)
		attr := &MemberAttr{}
		ma, ok := memberAttr.(map[interface{}]interface{})
		if !ok {
			err = ErrInvalidMemberAttr
		}
		err = attr.load(ma)
		if err != nil {
			return
		}
		param.Members[mn] = attr
	}
	return
}

// 解析API参数
func parseApiParams(paramsDef map[interface{}]interface{}) (params []*ApiParam, err error) {
	for from, def := range paramsDef {
		strFrom, _ := from.(string)
		param := &ApiParam{From: strFrom}

		d, ok := def.(map[interface{}]interface{})
		if !ok {
			err = ErrInvalidParamsDef
		}

		err = parseApiParamDetail(param, d)
		if err != nil {
			return
		}

		params = append(params, param)
	}
	return
}


func parseApiReturnData(apiReturn *ApiReturn, data interface{}) (err error) {
	switch data.(type) {
	case string:
		apiReturn.Data = data.(string)
	case map[interface{}]interface{}:
		dataMembers := make(map[string]*MemberAttr)
		for memberName, memberAttr := range data.(map[interface{}]interface{}) {
			attr := &MemberAttr{}

			var a map[interface{}]interface{}
			ok := true
			switch memberAttr.(type) {
			case map[interface{}]interface{}:
				a = memberAttr.(map[interface{}]interface{})
			case map[string]interface{}:
				t := memberAttr.(map[string]interface{})
				a = make(map[interface{}]interface{})
				for key, attrVal := range t {
					a[key] = attrVal
				}
			default:
				ok = false
			}

			if !ok {
				err = ErrInvalidMemberAttr
				return
			}

			err = attr.load(a)
			if err != nil {
				return
			}

			mn := memberName.(string)
			dataMembers[mn] = attr
		}
		apiReturn.Data = dataMembers
	}

	return
}

// 解析API的返回值
func parseApiReturns(returnsDef map[interface{}]interface{}) (returns map[string]*ApiReturn, err error) {
	returns = make(map[string]*ApiReturn)
	for status, returnDef := range returnsDef {
		statusCode := status.(string)
		ret := &ApiReturn{}
		def, ok := returnDef.(map[interface{}]interface{})
		if !ok {
			err = ErrInvalidReturnDef
			return
		}

		// Ignore description
		// type: nocontent, json, file
		retTypeVal, hasRetType := def["type"]
		if !hasRetType {
			ret.ReturnType = RETURN_TYPE_JSON
		} else {
			ret.ReturnType, err = ToString(retTypeVal)
			if ret.ReturnType != RETURN_TYPE_JSON &&
				ret.ReturnType != RETURN_TYPE_FILE &&
				ret.ReturnType != RETURN_TYPE_NOCONTENT {
				err = ErrReturnTypeUnsupported
				return
			}
		}
		// data
		if ret.ReturnType == RETURN_TYPE_JSON {
			returnDataDef, hasReturnData := def["data"]
			if !hasReturnData {
				err = ErrApiNoReturn
				return
			}
			err = parseApiReturnData(ret, returnDataDef)
			if err != nil {
				return
			}
		}

		returns[statusCode] = ret
	}
	return
}

//  解析API到服务的映射
func parseApiMapper() (err error) {
	return
}

var httpMethods = map[string]bool {
	"get": true,
	"post": true,
	"put": true,
	"delete": true,
	"options": true,
	"head": true,
	"patch": true,
	"trace": true,
}

func checkMethod(method string) error {
	method = strings.ToLower(method)
	_, exist := httpMethods[method]
	if !exist {
		return ErrInvalidApiMethod
	}
	return nil
}

func parseApi(apiDef map[interface{}]interface{}) (entry *ApiEntry, err error) {
	entry = &ApiEntry{}

	urlVal, hasUrl := apiDef["url"]
	if !hasUrl {
		err = ErrApiNoUrl
		return
	}
	entry.Url, err = ToString(urlVal)
	if err != nil {
		err = ErrApiUrlNotString
		return
	}

	methodVal, hasMethod := apiDef["method"]
	if !hasMethod {
		entry.Method = "get" // default method is get
	} else {
		entry.Method, err = ToString(methodVal)
		if err != nil {
			return
		}
		if err = checkMethod(entry.Method); err != nil {
			return
		}
	}
	// Ignore description

	// 允许不存在参数的调用
	paramsVal, hasParams := apiDef["params"]
	if hasParams {
		params, ok := paramsVal.(map[interface{}]interface{})
		if !ok {
			err = ErrInvalidParamsDef
		}
		entry.Params, err = parseApiParams(params)
	}

	returnsVal, hasReturns := apiDef["returns"]
	if !hasReturns {
		err = ErrApiNoReturn
		return
	}

	returns, ok := returnsVal.(map[interface{}]interface{})
	if !ok {
		err = ErrInvalidReturnDef
	}

	entry.Returns, err = parseApiReturns(returns)

	// TODO: add api mapper
	parseApiMapper()

	return
}

func (doc *ApiDoc) parseApis(apis []interface{}) (err error) {
	if len(apis) == 0 {
		err = ErrNoApis
		return
	}

	for _, apiMeta := range apis {
		api, ok := apiMeta.(map[interface{}]interface{})
		if !ok {
			err = ErrInvalidApiDef
			return
		}

		apiEntry, e := parseApi(api)
		if e != nil {
			err = e
			return
		}

		doc.Apis = append(doc.Apis, apiEntry)
	}

	return
}

