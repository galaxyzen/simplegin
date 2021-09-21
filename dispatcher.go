package simplegin

import (
	"net/http"
	"path"
	"strings"
)

type Dispatcher struct {
	basePath    string
	router      *router
	root        bool
	middlewares HandlersChain
	notFound404 HandlersChain
}

func NewDispatcher() *Dispatcher {
	dispatcher := &Dispatcher{
		basePath: "/",
		router:   newRouter(),
		root:     true,
	}
	dispatcher.Use(Logger())

	return dispatcher
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req)
	handlers, params := d.router.query(ctx.Method, ctx.Path)
	if handlers != nil {
		ctx.Params = params
		ctx.handlers = handlers
		ctx.Next()
		return
	}

	ctx.String(http.StatusNotFound, "404 NOT FOUND: %s: %s", ctx.Method, ctx.Path)
	ctx.handlers = d.notFound404
	ctx.Next()
}

func (d *Dispatcher) Group(relativePath string, middlewares ...HandlerFunc) *Dispatcher {
	return &Dispatcher{
		basePath:    d.calculateAbsolutePath(relativePath),
		router:      d.router,
		middlewares: d.combineMiddlewares(middlewares),
	}
}

func (d *Dispatcher) Use(middlewares ...HandlerFunc) *Dispatcher {
	d.middlewares = d.combineMiddlewares(middlewares)
	if d.root {
		d.notFound404 = d.combineMiddlewares(nil) // TODO a api for adding customized 404 middlewares
	}
	return d
}

func (d *Dispatcher) GET(relativePath string, handler HandlerFunc) {
	d.register("GET", relativePath, handler)
}

func (d *Dispatcher) POST(relativePath string, handler HandlerFunc) {
	d.register("POST", relativePath, handler)
}

func (d *Dispatcher) Static(relativePath string, root string) {
	// create handlerFunc
	var fs http.FileSystem = http.Dir(root)
	handler := d.createStaticHandler(relativePath, fs)

	// create *wildcard pattern
	pattern := path.Join(relativePath, "/*filepath")
	d.GET(pattern, handler)
}

func (d *Dispatcher) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := d.calculateAbsolutePath(relativePath)
	fileHandler := http.FileServer(fs)
	fileHandler = http.StripPrefix(absolutePath, fileHandler)

	return func(ctx *Context) {
		fileHandler.ServeHTTP(ctx.Writer, ctx.Req)
	}
}

func (d *Dispatcher) register(method string, relativePath string, handler HandlerFunc) {
	d.router.register(method, d.calculateAbsolutePath(relativePath), append(d.middlewares, handler))
}

func (d *Dispatcher) calculateAbsolutePath(relativePath string) string {
	absolutePath := path.Join(d.basePath, relativePath) // without end slash
	if strings.HasSuffix(relativePath, "/") && !strings.HasSuffix(absolutePath, "/") {
		absolutePath += "/"
	}

	return absolutePath
}

func (d *Dispatcher) combineMiddlewares(middlewares HandlersChain) HandlersChain {
	finalSize := len(d.middlewares) + len(middlewares)

	mergedMiddlewares := make(HandlersChain, finalSize)
	copy(mergedMiddlewares, d.middlewares)
	copy(mergedMiddlewares[len(d.middlewares):], middlewares)
	return mergedMiddlewares
}

func (d *Dispatcher) Run(addr string) error {
	return http.ListenAndServe(addr, d)
}
