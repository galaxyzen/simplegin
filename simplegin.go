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
	middlewares []HandlerFunc
}

type Engine struct {
	RouteGroup
	router *router
}

func New() *Engine {
	engine := &Engine{
		RouteGroup: RouteGroup{basePath: "/"},
		router:     newRouter(),
	}
	engine.Engine = engine // engine.RouteGroup.Engine = engine

	return engine
}

func (group *RouteGroup) Group(relativePath string, middlewares ...HandlerFunc) *RouteGroup {
	return &RouteGroup{
		Engine:      group.Engine,
		basePath:    group.calculateAbsolutePath(relativePath),
		middlewares: append(group.middlewares, middlewares...),
	}
}

func (group *RouteGroup) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := newContext(w, req)
	group.router.handle(ctx)
}

func (group *RouteGroup) register(method string, relativePath string, handler HandlerFunc) {
	// 添加为 handlers
	group.router.register(method, group.calculateAbsolutePath(relativePath), append(group.middlewares, handler))
}

func (group *RouteGroup) GET(relativePath string, handler HandlerFunc) {
	group.register("GET", relativePath, handler)
}

func (group *RouteGroup) POST(relativePath string, handler HandlerFunc) {
	group.register("POST", relativePath, handler)
}

func (group *RouteGroup) calculateAbsolutePath(relativePath string) string {
	absolutePath := path.Join(group.basePath, relativePath) // without end slash
	if strings.HasSuffix(relativePath, "/") && !strings.HasSuffix(absolutePath, "/") {
		absolutePath += "/"
	}

	return absolutePath
}

func (group *RouteGroup) Run(addr string) error {
	return http.ListenAndServe(addr, group)
}
