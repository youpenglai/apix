package apibuilder

import (
	"errors"
	"gopkg.in/yaml.v3"
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
	ErrNoDataTypeName = errors.New("no data type name")
	ErrInvalidDataTypeName = errors.New("invalid data type name")
	ErrDuplicateDataType = errors.New("duplicate data type")
	ErrInvalidApiDef = errors.New("invalid api definition")
	ErrApiNoUrl = errors.New("api has no url")
	ErrApiUrlNotString = errors.New("api url is not string")
	ErrReturnTypeUnsupported = errors.New("return type unsupported")
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

func (ma *MemberAttr) load(attrs map[string]interface{}) (err error) {
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
		ma.Required, _ = toBool(requiredVal)
	}

	descriptionVal, hasDescription := attrs["description"]
	if hasDescription {
		ma.Description = ""
	} else {
		ma.Description, _ = toString(descriptionVal)
	}

	lengthVal, hasLength := attrs["length"]
	if !hasLength {
		ma.Length.Value = 0
		ma.Length.Checked = false
	} else {
		ma.Length.Checked = true
		val, _ := toInt(lengthVal)
		ma.Length.Value = int(val)
	}

	minLengthVal, hasMinLength := attrs["minLength"]
	if !hasMinLength {
		ma.MinLength.Value = 0
		ma.MinLength.Checked = false
	} else {
		ma.MinLength.Checked = true
		val, _ := toInt(minLengthVal)
		ma.MinLength.Value = int(val)
	}

	maxLengthVal, hasMaxLength := attrs["maxLength"]
	if !hasMaxLength {
		ma.MaxLength.Value = 0
		ma.MaxLength.Checked = false
	} else {
		ma.MaxLength.Checked = true
		val, _ := toInt(maxLengthVal)
		ma.MaxLength.Value = int(val)
	}

	return nil
}

func (ma *MemberAttr) CheckType() bool {
	return false
}

func (ma *MemberAttr) CheckLength() bool {
	return false
}

func (ma *MemberAttr) CheckMinLength() bool {
	return false
}

func (ma *MemberAttr) CheckMaxLength() bool {
	return false
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
	Returns     map[string]ApiReturn  // API返回值
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
		doc.parseDataTypes(dataTypes.([]interface{}))
	}

	apis, hasApis := yamlDoc["apis"]
	if !hasApis {
		err = ErrNoApis
		return
	}

	doc.parseDataTypes(dataTypes.([]interface{}))
	doc.parseApis(apis.([]interface{}))

	return
}


func (doc *ApiDoc) parseBaseInfo(d map[string]interface{}) (err error) {
	// 文档描述
	descVal, _ := d["description"]
	doc.Description, _ = toString(descVal)

	// 文档版本
	versionVal, hasVersion := d["version"]
	if !hasVersion {
		err = ErrNoDocVersion
		return
	}
	doc.Version, _ = toString(versionVal)

	// 接口baseUrl
	baseUrlVal, hasBaseUrl := d["baseUrl"]
	if !hasBaseUrl {
		err = ErrNoBaseUrl
		return
	}
	doc.BaseUrl, _ = toString(baseUrlVal)

	return
}

func parseDataTypeMembers(dataType *DataType, members map[string]interface{}) (err error) {
	for name, attrs := range members {
		memberAttr := &MemberAttr{}
		a, ok := attrs.(map[string]interface{})
		if !ok {
			// TODO: 不符合的数据
		}
		err = memberAttr.load(a)
		if err != nil {
			return
		}
		dataType.AddMember(name, memberAttr)
	}
	return
}

func (doc *ApiDoc) parseDataTypes(dataTypes []interface{}) (err error) {
	for _, dt := range dataTypes {
		dataTypeMeta, ok := dt.(map[string]interface{})
		if !ok {
			err = ErrInvalidDataType
			return
		}
		nameVal, hasName := dataTypeMeta["name"]
		if !hasName {
			err = ErrNoDataTypeName
			return
		}
		name, nameErr := toString(nameVal)
		if nameErr != nil {
			err = ErrInvalidDataTypeName
			return
		}
		dataType := NewDataType(doc, name)
		members, hasMembers := dataTypeMeta["members"]
		if !hasMembers {
			// TODO: no members
		}
		dstMembers, ok := members.(map[string]interface{})
		if !ok {
			// TODO: invalid member
		}
		err = parseDataTypeMembers(dataType, dstMembers)
		if err != nil {
			// TODO: invalid members
		}
		if err = doc.addDataType(dataType); err != nil {
			return
		}
	}
	return
}

func parseApiParamDetail(param *ApiParam, members map[string]interface{}) (err error) {
	param.Members = make(Members)
	for memberName, memberAttr := range members {
		attr := &MemberAttr{}
		ma, ok := memberAttr.(map[string]interface{})
		if !ok {
			// TODO: invalid member attr
		}
		err = attr.load(ma)
		if err != nil {
			return
		}
		param.Members[memberName] = attr
	}
	return
}

// 解析API参数
func parseApiParams(paramsDef map[string]interface{}) (params []*ApiParam, err error) {
	for from, def := range paramsDef {
		param := &ApiParam{From: from}

		d, ok := def.(map[string]interface{})
		if !ok {
			// TODO: invalid param
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
	case map[string]interface{}:
		dataMembers := make(map[string]*MemberAttr)
		for memberName, memberAttr := range data.(map[string]interface{}) {
			attr := &MemberAttr{}

			a, ok := memberAttr.(map[string]interface{})
			if !ok {
				// TODO:
			}

			err = attr.load(a)
			if err != nil {
				return
			}

			dataMembers[memberName] = attr
		}
		apiReturn.Data = dataMembers
	}

	return
}

// 解析API的返回值
func parseApiReturns(returnsDef map[string]interface{}) (returns map[string]*ApiReturn, err error) {
	returns = make(map[string]*ApiReturn)
	for statusCode, returnDef := range returnsDef {
		ret := &ApiReturn{}
		def, ok := returnDef.(map[string]interface{})
		if !ok {
			// TODO
		}

		// Ignore description
		// type: nocontent, json, file
		retTypeVal, hasRetType := def["type"]
		if !hasRetType {
			ret.ReturnType = RETURN_TYPE_JSON
		} else {
			ret.ReturnType, err = toString(retTypeVal)
			if ret.ReturnType != RETURN_TYPE_JSON &&
				ret.ReturnType != RETURN_TYPE_FILE &&
				ret.ReturnType != RETURN_TYPE_NOCONTENT {
				err = ErrReturnTypeUnsupported
				return
			}
		}
		// data
		if ret.ReturnType == RETURN_TYPE_JSON {

		}

		returns[statusCode] = ret
	}
	return
}

//  解析API到服务的映射
func parseApiMapper() (err error) {
	return
}

func parseApi(apiDef map[string]interface{}) (entry *ApiEntry, err error) {
	entry = &ApiEntry{}

	urlVal, hasUrl := apiDef["url"]
	if !hasUrl {
		err = ErrApiNoUrl
		return
	}
	entry.Url, err = toString(urlVal)
	if err != nil {
		err = ErrApiUrlNotString
		return
	}

	methodVal, hasMethod := apiDef["method"]
	if !hasMethod {
		entry.Method = "get" // default method is get
	} else {
		entry.Method, err = toString(methodVal)
		if err != nil {
			return
		}
		// TODO: check method valid
	}
	// Ignore description

	// 允许不存在参数的调用
	paramsVal, hasParams := apiDef["params"]
	if hasParams {
		params, ok := paramsVal.(map[string]interface{})
		if !ok {
			// TODO:
		}
		entry.Params, err = parseApiParams(params)
	}

	returnsVal, hasReturns := apiDef["returns"]
	if !hasReturns {
		// TODO: no returns
	}

	returns, ok := returnsVal.(map[string]interface{})
	if !ok {
		// TODO:
	}

	parseApiReturns(returns)

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
		api, ok := apiMeta.(map[string]interface{})
		if !ok {
			err = ErrInvalidApiDef
			return
		}


	}

	return
}

