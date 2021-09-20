package simplegin

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

func (r *router) query(method string, requestPath string) (HandlersChain, map[string]string) {
	var handlers HandlersChain
	var params map[string]string
	if root := r.trees[method]; root != nil {
		handlers, params = root.search(requestPath)
	}
	return handlers, params
}
