package apix

import (
	"strings"
	"errors"
)

var (
	ErrRouterNotFound = errors.New("router not found")
	ErrMethodNotFound = errors.New("method not found")
)

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

	if re.subEntries == nil {
		re.subEntries = make(map[string]*routerEntry)
	}

	re.subEntries[path] = sub
}

func (re *routerEntry) setParamEntry(sub *routerEntry) {
	if sub == nil {
		return
	}

	re.paramEntry = sub
}

func (re *routerEntry) setMiddlewares(handlers ...Handler) {
	if len(handlers) == 0 {
		return
	}

	re.middlewares = make([]Handler, len(handlers))
	copy(re.middlewares, handlers)
}

func (re *routerEntry) appendMiddlewares(handlers ...Handler) {
	re.middlewares = append(re.middlewares, handlers...)
}

func (re *routerEntry) bindMethod(method string, handlers ...Handler) {
	if len(handlers) == 0 {
		return
	}

	m := strings.ToUpper(method)
	switch m {
	case "GET", "PUT", "POST", "DELETE", "OPTIONS", "PATCH", "HEAD":
		if re.handlers == nil {
			re.handlers = make(map[string][]Handler)
		}
		re.handlers[m] = handlers
	// unsupported methods or invalid method
	default:
		panic("invalid http method")
	}
}

func (re *routerEntry) paramName() string {
	return re.name
}

func newRouterEntry(name string) *routerEntry {
	return &routerEntry{name:name}
}

func (re *routerEntry) mount(method, name string, handler ...Handler) {

}

type Router struct {
	routerEntry
}

func (r *Router) Group(path string, handlers ...Handler) *Router {
	re := r.buildEntries(path)
	re.setMiddlewares(handlers...)

	// copy router entry to new Router
	nr := &Router{routerEntry: *re}

	return nr
}

func isParam(pathName string) bool {
	// TODO: check param: :{1}\w
	return pathName[0] == ':' && len(pathName) > 1
}

func (r *Router) buildEntries(path string) *routerEntry {
	parts := strings.Split(strings.TrimSpace(path), "/")
	if len(parts) == 0 {
		return &r.routerEntry
	}

	p := &(r.routerEntry)
	var newEntry *routerEntry
	for _, part := range parts {
		if part == "" {
			continue
		}

		newEntry = newRouterEntry(part)
		if isParam(part) {
			p.setParamEntry(newEntry)
		} else {
			p.addSubEntry(part, newEntry)
		}
		// next level
		p = newEntry
	}

	return newEntry
}

// set global handlers
func (r *Router) Use(handlers ...Handler) {
	if len(handlers) == 0 {
		// panic?
		return
	}
	r.appendMiddlewares(handlers...)
}

func (r *Router) Get(path string, handlers ...Handler) {
	finalEntry := r.buildEntries(path)
	finalEntry.bindMethod("GET", handlers...)
}

func (r *Router) Post(path string, handlers ...Handler) {
	finalEntry := r.buildEntries(path)
	finalEntry.bindMethod("POST", handlers...)
}

func (r *Router) Put(path string, handlers ...Handler) {
	finalEntry := r.buildEntries(path)
	finalEntry.bindMethod("PUT", handlers...)
}

func (r *Router) Delete(path string, handlers ...Handler) {
	finalEntry := r.buildEntries(path)
	finalEntry.bindMethod("DELETE", handlers...)
}

func (r *Router) Options(path string, handlers ...Handler) {
	finalEntry := r.buildEntries(path)
	finalEntry.bindMethod("OPTIONS", handlers...)
}

func (r *Router) Patch(path string, handlers ...Handler) {
	finalEntry := r.buildEntries(path)
	finalEntry.bindMethod("PATCH", handlers...)
}

func (r *Router) Head(path string, handlers ...Handler) {
	finalEntry := r.buildEntries(path)
	finalEntry.bindMethod("HEAD", handlers...)
}
// TODO: add more http method handler

func (r *Router) match(path string, method string) (handlers []Handler, urlParams map[string]string, err error) {
	parts := strings.Split(path, "/")

	urlParams = make(map[string]string)

	re := &r.routerEntry
	var exist bool
	for _, part := range parts {
		if part == "" {
			return
		}

		handlers = append(handlers, re.middlewares...)
		if re, exist = re.subEntries[part]; exist {
			continue
		}

		if re.paramEntry != nil {
			re = re.paramEntry
			urlParams[re.name[1:]] = part
		} else {
			err = ErrRouterNotFound
			return
		}
	}

	if methodHandlers, exist := re.handlers[method]; !exist {
		err = ErrMethodNotFound
		return
	} else {
		handlers = append(handlers, methodHandlers...)
	}

	return
}
