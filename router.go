package apix

type Handler func(ctx *Context)


type routerEntry struct {
	name    string
	isParam bool

	// handler
	subEntries map[string]*routerEntry
	paramEntry *routerEntry
	// handlers
	handlers map[string][]Handler
}

func (re *routerEntry) match(path string) *routerEntry {
	if re.isParam {
		return re.paramEntry
	}

	entry, _ := re.subEntries[path]

	return entry
}

func (re *routerEntry) paramName() string {
	return re.name
}

func (re *routerEntry) mount(method, name string, handler ...Handler) {

}

type Router struct {
	routerEntry
}

type RouterGroup struct {
	routerEntry
}

func (r *Router) Group() {}

func (r *Router) mountPath(path string) {}

func (r *Router) Use(handlers ...Handler) {

}

func (r *Router) Get(path string, handlers ...Handler) {

}

func (r *Router) Post(path string, handlers ...Handler) {

}

func (r *Router) Put(path string, handlers ...Handler) {

}

func (r *Router) Delete(path string, handlers ...Handler) {

}

// TODO: add more http method handler

func (r *Router) match(path string) Handler {
	return nil
}
