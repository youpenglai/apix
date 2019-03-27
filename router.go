package apix

import "strings"

type Handler func(ctx *Context)


type routerEntry struct {
	name    string
	// handler
	subEntries map[string]*routerEntry
	paramEntry *routerEntry
	// middlewares
	middlewares    []Handler
	handlers map[string][]Handler
}

//func (re *routerEntry) match(path string) *routerEntry {
//	if re.isParam {
//		return re.paramEntry
//	}
//
//	entry, _ := re.subEntries[path]
//
//	return entry
//}

func (re *routerEntry) addSubEntry(path string, sub *routerEntry) {
	if sub == nil {
		return
	}

	re.subEntries[path] = sub
}

func (re *routerEntry) setParamEntry(sub *routerEntry) {
	if sub == nil {
		return
	}

	re.paramEntry = sub
}

func (re *routerEntry) bindMethod(method string, handlers ...Handler) {
	if len(handlers) == 0 {
		return
	}

	m := strings.ToUpper(method)
	switch m {
	case "GET", "PUT", "POST", "DELETE", "OPTIONS", "PATCH", "HEAD":
		re.handlers[m] = handlers
	// unsupported methods or invalid method
	default:
		panic("invalid http method")
	}
}

func (re *routerEntry) paramName() string {
	return re.name
}

func newRouterEntry() *routerEntry {
	return &routerEntry{}
}

func (re *routerEntry) mount(method, name string, handler ...Handler) {

}

type Router struct {
	routerEntry
}

type RouterGroup struct {
	Router
	entry *routerEntry
}

func NewRouterGroup(groupPath string) *RouterGroup {
	g := &RouterGroup{}

	g.entry = g.buildEntries(groupPath)

	return g
}

func (rg *RouterGroup) AddRouter(router *Router) {

}

func (r *Router) Group() {}

func isParam(pathName string) bool {
	// TODO: check param: :{1}\w
	return pathName[0] == ':' && len(pathName) > 1
}

func (r *Router) buildEntries(path string) *routerEntry {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return &r.routerEntry
	}

	p := &r.routerEntry
	var newEntry *routerEntry
	for _, part := range parts {
		newEntry = newRouterEntry()
		newEntry.name = part
		if isParam(part) {
			p.paramEntry = newEntry
		} else {
			p.subEntries[part] = newEntry
		}
		// next level
		p = newEntry
	}

	return newEntry
}

func (r *Router) Use(handlers ...Handler) {
	if len(handlers) == 0 {
		// panic?
		return
	}
	r.middlewares = append(r.middlewares, handlers...)
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
