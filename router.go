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

type Router struct {
	name    string
	// handler
	subEntries map[string]*Router
	paramEntry *Router
	// middlewares
	middlewares    []Handler
	handlers map[string][]Handler
}


func (re *Router) addSubEntry(path string, sub *Router) {
	if sub == nil {
		return
	}

	if re.subEntries == nil {
		re.subEntries = make(map[string]*Router)
	}

	re.subEntries[path] = sub
}

func (re *Router) setParamEntry(sub *Router) {
	if sub == nil {
		return
	}

	re.paramEntry = sub
}

func (re *Router) setMiddlewares(handlers ...Handler) {
	if len(handlers) == 0 {
		return
	}

	re.middlewares = make([]Handler, len(handlers))
	copy(re.middlewares, handlers)
}

func (re *Router) appendMiddlewares(handlers ...Handler) {
	re.middlewares = append(re.middlewares, handlers...)
}

func (re *Router) bindMethod(method string, handlers ...Handler) {
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

func (re *Router) paramName() string {
	return re.name[1:]
}

func newRouterEntry(name string) *Router {
	return &Router{name:name}
}

func (r *Router) Group(path string, handlers ...Handler) *Router {
	re := r.buildEntries(path)
	re.setMiddlewares(handlers...)

	return re
}

func isParam(pathName string) bool {
	// TODO: check param: :{1}\w
	return pathName[0] == ':' && len(pathName) > 1
}

func (r *Router) buildEntries(path string) *Router {
	parts := strings.Split(strings.TrimSpace(path), "/")
	if len(parts) == 0 {
		return r
	}

	p := r
	var newEntry *Router
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

	if newEntry == nil {
		newEntry = r
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

func (r *Router) match(path string, method string) (handlers []Handler, urlParams Params, err error) {
	parts := strings.Split(path, "/")

	urlParams = NewParams()
	handlers = make([]Handler,0)

	re := r
	var exist bool
	var sub *Router
	for _, part := range parts {
		if part == "" {
			continue
		}

		// add handlers
		handlers = append(handlers, re.middlewares...)
		if sub, exist = re.subEntries[part]; !exist {
			// 如果路径在子项目里面不存在，尝试读取参数入口
			// 并记录下当前参数
			if re.paramEntry != nil {
				re = re.paramEntry
				urlParams.AddValue(re.paramName(), part)
			} else {
				err = ErrRouterNotFound
				return
			}
		} else {
			re = sub
		}
	}

	if methodHandlers, exist := re.handlers[strings.ToUpper(method)]; !exist {
		err = ErrMethodNotFound
	} else {
		if re.middlewares != nil {
			handlers = append(handlers, re.middlewares...)
		}
		handlers = append(handlers, methodHandlers...)
	}

	return
}
