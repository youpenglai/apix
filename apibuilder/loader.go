package apibuilder

// API字段成员
type ApiMembers map[string]ApiMemberAttr

func NewApiMember() ApiMembers {
	return make(ApiMembers)
}

// API数据复合类型
type ApiType struct {
	Name    string     // 数据类型
	Members ApiMembers // 类型字段描述
}

// 添加API类型
func NewApiType(name string) *ApiType {
	return &ApiType{Name:name, Members: NewApiMember()}
}

// 向API类型中增加成员
func (t *ApiType) AddMember(name string, attr ApiMemberAttr) {
	t.Members[name] = attr
}

type MemberEachFunc func(name string, attr ApiMemberAttr)

// 遍历数据类型中的成员
func (t *ApiType) MembersEach(eachFunc MemberEachFunc) {
	for name, attr := range t.Members {
		eachFunc(name, attr)
	}
}

// API字段成员属性
type ApiMemberAttr struct {
	Type        string // 字段数据基本类型
	IsArray     bool   // 是否为数组
	Required    bool   // 是否为必须字段
	Description string // 字段描述
	Length      int    // 字段长度
	MinLength   int    // 字段最小长度
	MaxLength   int    // 字段长度
}

// API返回值
type ApiReturn struct {
	File string                // 返回为文件类型，输出的文件名
	Data map[string]ApiMembers // 返回为JSON对象
}

// API入口
type ApiEntry struct {
	Url         string                // API接口路径
	Method      string                // API方法
	Params      map[string]ApiMembers // API参数映射
	Returns     map[int]ApiReturn     // API返回值
	Description string                // API描述
}

// API描述文档
type ApiDesc struct {
	Name    string      // API名称
	Version string      // API版本（语义化版本）：major.minor.revision
	BaseUrl string      // 基本URL
	Apis    []*ApiEntry // API入口
	Types   map[string]*ApiType  // API中引用的数据类型定义
}
