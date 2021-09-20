package simplegin

import (
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type RouteGroup struct {
	*Engine
	basePath    string
	middlewares HandlersChain
	root        bool
}

type Engine struct {
	RouteGroup
	router      *router
	notFound404 HandlersChain
}

func New() *Engine {
	engine := &Engine{
		RouteGroup: RouteGroup{
			basePath: "/",
			root:     true,
		},
		router: newRouter(),
	}
	engine.Engine = engine // engine.RouteGroup.Engine = engine
	engine.Use(Logger())

	return engine
}

func (group *RouteGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req)
	handlers, params := group.router.query(ctx.Method, ctx.Path)
	if handlers != nil {
		ctx.Params = params
		ctx.handlers = handlers
		ctx.Next()
		return
	}

	ctx.String(http.StatusNotFound, "404 NOT FOUND: %s: %s", ctx.Method, ctx.Path)
	ctx.handlers = group.notFound404
	ctx.Next()
}

func (group *RouteGroup) Group(relativePath string, middlewares ...HandlerFunc) *RouteGroup {
	return &RouteGroup{
		Engine:      group.Engine,
		basePath:    group.calculateAbsolutePath(relativePath),
		middlewares: group.combineMiddlewares(middlewares),
	}
}

func (group *RouteGroup) Use(middlewares ...HandlerFunc) *RouteGroup {
	group.middlewares = group.combineMiddlewares(middlewares)
	if group.root {
		group.notFound404 = group.combineMiddlewares(nil) // TODO a api for adding customized 404 middlewares
	}
	return group
}

func (group *RouteGroup) GET(relativePath string, handler HandlerFunc) {
	group.register("GET", relativePath, handler)
}

func (group *RouteGroup) POST(relativePath string, handler HandlerFunc) {
	group.register("POST", relativePath, handler)
}

func (group *RouteGroup) register(method string, relativePath string, handler HandlerFunc) {
	group.router.register(method, group.calculateAbsolutePath(relativePath), append(group.middlewares, handler))
}

func (group *RouteGroup) calculateAbsolutePath(relativePath string) string {
	absolutePath := path.Join(group.basePath, relativePath) // without end slash
	if strings.HasSuffix(relativePath, "/") && !strings.HasSuffix(absolutePath, "/") {
		absolutePath += "/"
	}

	return absolutePath
}

func (group *RouteGroup) combineMiddlewares(middlewares HandlersChain) HandlersChain {
	finalSize := len(group.middlewares) + len(middlewares)

	mergedMiddlewares := make(HandlersChain, finalSize)
	copy(mergedMiddlewares, group.middlewares)
	copy(mergedMiddlewares[len(group.middlewares):], middlewares)
	return mergedMiddlewares
}

func (group *RouteGroup) Run(addr string) error {
	return http.ListenAndServe(addr, group)
}
