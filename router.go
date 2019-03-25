package apix

type Handler func(ctx Context, next ...Handler)

type routerEntry struct {
	path string
	isParam bool

	// TODO: middlewares
	// TODO: handler
	subEntries map[string]*routerEntry
	paramEntry *routerEntry
}

func (re *routerEntry) match(path string) *routerEntry {
	if (re.isParam) {
		return re.paramEntry
	}

	entry, _ := re.subEntries[path]

	return entry
}

func (re *routerEntry) paramName() string {
	return re.path
}

type Router struct {
	path string
}

func (r *Router) Use(path string) {

}

func (r *Router) Get(path string, handler Handler) {

}

func (r *Router) Post(path string, handler Handler) {

}

func (r *Router) Put(path string, handler Handler) {

}

func (r *Router) Delete(path string, handler Handler) {

}

func (r *Router) match(path string) Handler {
	return nil
}
