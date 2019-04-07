package apibuilder

type ApiField map[string]ApiFieldAttr

type ApiFieldAttr struct {

}

type ApiParam struct {

}

type ApiReturn struct {
	File string
	Data map[string]*ApiField
}

type ApiEntry struct {
	Url string
	Method string
	Params map[string]ApiParam
	Returns map[int]ApiReturn
}

type ApiDesc struct {
	Name string
	Version string
	BaseUrl string
	Apis []string
}