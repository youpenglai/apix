package apibuilder

import (
	"errors"
	"strconv"
	"fmt"
	"bytes"
)

var (
	ErrNilCannotSetValue = errors.New("nil cannot set value")
	ErrInvalidArrayValue = errors.New("invalid array value")
	ErrValidationRequired = errors.New("required error")
	ErrValidationLength = errors.New("length error")
	ErrValidationMinLength = errors.New("min-length error")
	ErrValidationMaxLength = errors.New("max-length error")
	ErrNestArrayNotImplemented = errors.New("nest array not implemented")
)

// 参数读取/获取接口
type ParamReader interface {
	// name: 读取的参数名
	// from: 从什么位置读取
	Get(name, from string) interface{}
}

type ParamValues struct {
	schema *ApiParam
	reader ParamReader
}

// 变量构造器
type VariableConstructor func() Variable

// 变量接口
type Variable interface {
	// 序列化为JSON
	MarshalJSON() ([]byte, error)
	// 校验参数
	Validation() error
	// 设置变量值
	SetValue(interface{}) error
	SetAttr(attr *MemberAttr)
}

// 变量基础属性
type VariableBase struct {
	typeName string
	required bool
	length AttrLength
	minLength AttrLength
	maxLength AttrLength
	// TODO: 添加更多的变量值
}

func (vb *VariableBase) SetAttr(attr *MemberAttr) {
	vb.typeName = attr.Type
	vb.required = attr.Required
	vb.length = attr.Length
	vb.minLength = attr.MinLength
	vb.maxLength = attr.MaxLength
}

// 整型变量
type IntVar struct {
	VariableBase
	val int64
}

func (i *IntVar) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(i.val, 10)), nil
}

func (i *IntVar) Validation() error {
	return nil
}

func (i *IntVar) SetValue(val interface{}) error {
	var err error
	i.val, err = ToInt(val)
	return err
}

func IntConstructor() Variable {
	return &IntVar{}
}

// 浮点型变量
type FloatVar struct {
	VariableBase
	val float64
}

func (f *FloatVar) MarshalJSON() (data []byte, err error) {
	data = []byte(strconv.FormatFloat(f.val, 'e', -1, 64))
	return
}

func (f *FloatVar) Validation() (err error) {
	return
}

func (f *FloatVar) SetValue(val interface{}) (err error) {
	f.val, err = ToFloat(val)
	return
}

func FloatConstructor() Variable {
	return &FloatVar{}
}

// 字符串变量
type StringVar struct {
	VariableBase
	val string
}

func StringConstructor() Variable {
	return &StringVar{}
}

func (sv *StringVar) MarshalJSON() (data []byte, err error) {
	data = []byte(fmt.Sprintf(`"%s"`, sv.val))
	return
}

func (sv *StringVar) Validation() (err error) {
	// check required
	if sv.VariableBase.required && len(sv.val) == 0 {
		err = ErrValidationRequired
		return
	}
	// check length
	if sv.VariableBase.length.Checked && len(sv.val) < sv.VariableBase.length.Value {
		err = ErrValidationLength
		return
	}
	// check minLength
	if sv.VariableBase.minLength.Checked && len(sv.val) < sv.VariableBase.minLength.Value {
		err = ErrValidationMinLength
		return
	}
	// check maxLength
	if sv.VariableBase.maxLength.Checked && len(sv.val) < sv.VariableBase.maxLength.Value {
		err = ErrValidationMaxLength
		return
	}
	return
}

func (sv *StringVar) SetValue(val interface{}) (err error) {
	sv.val, err = ToString(val)
	return
}

// 布尔型变量
type BooleanVar struct {
	VariableBase
	val bool
}

func (bv *BooleanVar) MarshalJSON() (data []byte, err error) {
	return []byte(strconv.FormatBool(bv.val)), nil
}

func (bv *BooleanVar) Validation() (err error) {
	return
}

func (bv *BooleanVar) SetValue(val interface{}) (err error) {
	bv.val, err = ToBool(val)
	return
}

func BooleanConstructor() Variable {
	return &BooleanVar{}
}

// 数组变量，这是个容器类型，很麻烦
type ArrayVar struct {
	VariableBase
	code *ApiCode
	val []Variable
}

func ArrayConstructor() Variable {
	return &ArrayVar{val:make([]Variable, 0)}
}

func (av *ArrayVar) setCode(code *ApiCode) {
	av.code = code
}

func (av *ArrayVar) MarshalJSON() (data []byte, err error) {
	buf := bytes.NewBuffer(data)
	buf.WriteByte('[')

	last := len(av.val) - 1
	for i, v := range av.val {
		j, e := v.MarshalJSON()
		if e != nil {
			err = e
			return
		}
		buf.Write(j)
		if i != last {
			buf.WriteByte(',')
		}
	}

	buf.WriteByte(']')
	data = buf.Bytes()
	return
}

func (av *ArrayVar) Validation() (err error) {
	// 检查本身
	if av.required && len(av.val) == 0{
		err = ErrValidationRequired
		return
	}
	if av.length.Checked && len(av.val) != av.length.Value {
		err = ErrValidationLength
		return
	}
	if av.minLength.Checked && len(av.val) < av.minLength.Value {
		err = ErrValidationMinLength
		return
	}
	if av.maxLength.Checked && len(av.val) > av.maxLength.Value {
		err = ErrValidationMaxLength
		return
	}
	// 检查内部成员
	for _, item := range av.val {
		if err = item.Validation(); err != nil {
			return
		}
	}
	return
}

func (av *ArrayVar) SetValue(val interface{}) (err error) {
	//f.val, err = ToString(val)
	arr, ok := val.([]interface{})
	if !ok {
		err = ErrInvalidArrayValue
		return
	}

	var constructor *DataTypeConstructor
	constructor, err = av.code.getDataType(av.VariableBase.typeName)
	if err != nil {
		return
	}

	for _, item := range arr {
		v := constructor.constructor()
		if ov, ok := v.(*ObjectVar); ok {
			// 变量为对象
			ov.setCode(av.code)
			ov.setType(av.typeName)
		} else if av, ok := v.(*ArrayVar); ok {
			// TODO：传入数据为数组的情况处理，多维数组
			av = av
			err = ErrNestArrayNotImplemented
			return
		}

		v.SetValue(item)
	}

	return
}

// 对象变量，同样是复合类型啊
type ObjectVar struct {
	VariableBase
	code *ApiCode
	Attrs map[string]Variable
}

func (ov *ObjectVar) setType(typeName string) {
	ov.VariableBase.typeName = typeName
}

func (ov *ObjectVar) addAttr(name string, variable Variable) {
	ov.Attrs[name] = variable
}

func (ov *ObjectVar) setCode(code *ApiCode) {
	ov.code = code
}

func (ov *ObjectVar) MarshalJSON() (data []byte, err error) {
	buff := bytes.NewBuffer(data)
	buff.WriteByte('{')
	last := len(ov.Attrs) - 1
	for k, v := range ov.Attrs {
		j, e := v.MarshalJSON()
		if e != nil {
			err = e
			return
		}

		buff.WriteString("\"" + k + "\":")
		buff.Write(j)
		if last != 0 {
			buff.WriteByte(',')
			last--
		}
	}
	buff.WriteByte('}')
	data = buff.Bytes()
	return
}

func (ov *ObjectVar) Validation() (err error) {
	for _, v := range ov.Attrs {
		err = v.Validation()
		if err != nil {
			return
		}
	}
	return
}

func (ov *ObjectVar) SetValue(val interface{}) (err error) {
	members, ok := val.(map[string]interface{})
	if !ok {
		// TODO: member
	}

	typeDef, _ := ov.code.getDataType(ov.VariableBase.typeName)
	for mn, ma := range typeDef.dataType.Members {
		val, exist := members[mn]
		if !exist && ma.Required {
			err = ErrValidationRequired
			return
		}
		var v Variable
		if v, err = readData(ov.code, val, ma); err != nil {
			return
		}
		ov.Attrs[mn] = v
	}
	return
}

func ObjectConstructor() Variable {
	return &ObjectVar{Attrs:make(map[string]Variable)}
}

// 空值变量
type NilParam struct{}

func (np *NilParam) MarshalJSON() ([]byte, error){
	return []byte{'n', 'u', 'l', 'l'}, nil
}

func (np *NilParam) Validation() error {
	return nil
}

// 空值不能赋值
func (np *NilParam) SetValue(val interface{}) error {
	return ErrNilCannotSetValue
}

func (np *NilParam) SetAttr(attr *MemberAttr) {
	// do nothing
}

// 参数变量
type ParamVar Variable

// Api接口代码块，每个接口生成一个代码块
// 不是真正意义的代码块
type ApiCodeBlock struct {
	code *ApiCode
	params []*ApiParam
	paramReader ParamReader
}

// 绑定参数读取接口
func (acb *ApiCodeBlock) BindParamReader(reader ParamReader) {
	acb.paramReader = reader
}

// 读取数据
func readData(code *ApiCode, val interface{}, attr *MemberAttr) (v Variable, err error) {
	var typeConstructor *DataTypeConstructor
	if typeConstructor, err = code.getDataType(attr.Type); err != nil {
		return
	}

	if attr.IsArray {
		v = ArrayConstructor()
		arrayContainer, _ := v.(*ArrayVar)
		arrayContainer.setCode(code)
	} else {
		v = typeConstructor.constructor()
		ov, ok := v.(*ObjectVar)
		if ok {
			ov.setCode(code)
		}
	}
	v.SetValue(val)
	v.SetAttr(attr)
	return
}

// 读取参数
func (acb *ApiCodeBlock) ReadParams() (param ParamVar, err error) {
	// no params
	if len(acb.params) == 0 {
		param = &NilParam{}
		return
	}

	param = ObjectConstructor()
	paramVal, _ := param.(*ObjectVar)

	for _, inParam := range acb.params {
		for name, memberAttr := range inParam.Members {
			val := acb.paramReader.Get(name, inParam.From)
			variable , e:= readData(acb.code, val, memberAttr)
			if e != nil {
				err = e
				return
			}
			paramVal.addAttr(name, variable)
		}
	}

	return
}

// 读取返回值
func (acb *ApiCodeBlock) ReadReturn() {}

// 数据类型构造器
// constructor: 构造函数
// dataType: 数据类型
type DataTypeConstructor struct {
	constructor VariableConstructor
	dataType *DataType
}

// 获取构造器数据类型名称
func (c *DataTypeConstructor) GetName() string {
	return c.dataType.Name
}

func NewDataTypeConstructor(constructorFunc VariableConstructor, dataType *DataType) *DataTypeConstructor {
	return &DataTypeConstructor{constructor:constructorFunc, dataType:dataType}
}

// 基本数据类型
var baseDataTypeConstructor = map[string]*DataTypeConstructor{
	"integer": &DataTypeConstructor{dataType:NewBaseDataType("integer"), constructor:IntConstructor},
	"float": &DataTypeConstructor{dataType:NewBaseDataType("float"), constructor:FloatConstructor},
	"boolean": &DataTypeConstructor{dataType:NewBaseDataType("boolean"), constructor:BooleanConstructor},
	"string": &DataTypeConstructor{dataType:NewBaseDataType("string"), constructor:StringConstructor},
}

// Api代码
// 将输入的数据处理为何数据类型关联的数据
type ApiCode struct {
	dataTypeConstructor map[string]*DataTypeConstructor
	entries             map[string]*ApiCodeBlock
}

func NewApiCode() *ApiCode {
	return &ApiCode{
		dataTypeConstructor:make(map[string]*DataTypeConstructor),
		entries: make(map[string]*ApiCodeBlock),
	}
}

func (ac *ApiCode) addApiEntry(path string, block *ApiCodeBlock) {
	ac.entries[path] = block
}

func (ac *ApiCode) addDataTypeConstructor(constructor *DataTypeConstructor) {

	ac.dataTypeConstructor[constructor.GetName()] = constructor
}

func (ac *ApiCode) getDataType(name string) (constructor *DataTypeConstructor, err error) {
	var hasDataType bool
	constructor, hasDataType = ac.dataTypeConstructor[name]
	if !hasDataType {
		err = ErrDataTypeNotExist
	}
	return
}

func (ac *ApiCode) GetApiCode(path string) (codeBlock *ApiCodeBlock, err error) {
	var exists bool
	codeBlock, exists = ac.entries[path]
	if !exists {
		// TODO: err
		return
	}
	return
}

// 生成Api处理代码逻辑
// 额，非传统意义上的代码生成
func GenApiCode(doc *ApiDoc) (code *ApiCode, err error){
	code = NewApiCode()

	// Install base data types
	for _, base := range baseDataTypeConstructor {
		code.addDataTypeConstructor(base)
	}

	// Install user data types
	for _, dt := range doc.Types {
		dtc := &DataTypeConstructor{dataType: dt, constructor:ObjectConstructor}
		code.addDataTypeConstructor(dtc)
	}

	for _, api := range doc.Apis {
		block := &ApiCodeBlock{code: code}
		block.params = api.Params
		code.addApiEntry(api.Url, block)
	}
	return
}