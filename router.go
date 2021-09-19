package simplegin

import (
	"net/http"
)

type router struct {
	trees map[string]trieTree
}

func newRouter() *router {
	return &router{
		trees: map[string]trieTree{},
	}
}

func (r *router) register(method string, pattern string, handlers HandlersChain) {
	if tree, ok := r.trees[method]; ok {
		tree.insert(pattern, handlers)
	} else {
		tree = newTrieTree()
		tree.insert(pattern, handlers)
		r.trees[method] = tree
	}
}

func (r *router) handle(ctx *Context) {
	if root := r.trees[ctx.Method]; root != nil { // GET POST PATCH PUT DELETE Trie
		if handlers, params := root.search(ctx.Path); handlers != nil {
			ctx.Params = params
			ctx.handlers = append(ctx.handlers, handlers...)
			ctx.Next()
			return
		}
	}

	ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n for %s.", ctx.Path, ctx.Method)
}
