package apibuilder

import (
	"errors"
	"gopkg.in/yaml.v3"
)

var (
	ErrNoDocVersion = errors.New("no doc version")
	ErrNoBaseUrl = errors.New("no baseUrl")
	ErrNoApis = errors.New("no apis")
	ErrInvalidDataType = errors.New("invalid data type")
	ErrNoDataTypeName = errors.New("no data type name")
	ErrInvalidDataTypeName = errors.New("invalid data type name")
	ErrDuplicateDataType = errors.New("duplicate data type")
)

// API字段成员
type Members map[string]*MemberAttr

func NewApiMember() Members {
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
	return &DataType{doc: doc, Name: name, Members: NewApiMember()}
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

type MemberEachFunc func(name string, attr MemberAttr)

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
	File string             // 返回为文件类型，输出的文件名
	Data map[string]Members // 返回为JSON对象
}

// API入口
type ApiEntry struct {
	Url         string             // API接口路径
	Method      string             // API方法
	Params      map[string]Members // API参数映射
	Returns     map[int]ApiReturn  // API返回值
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

func (doc *ApiDoc) parseApis(apis []interface{}) (err error) {
	if len(apis) == 0 {
		err = ErrNoApis
		return
	}

	return
}

