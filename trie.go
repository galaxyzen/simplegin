package simplegin

import (
	"fmt"
	"strings"
)

type trieNode struct {
	pattern  string
	handlers HandlersChain
	part     string
	children []*trieNode
}

type trieTree = *trieNode

func newTrieTree() trieTree {
	return new(trieNode)
}

func (t *trieNode) insert(pattern string, handlers HandlersChain) {
	parts := splitURLPath(pattern)

	if conflict := checkConflict(t, parts); conflict != "" {
		s := fmt.Sprintf("conflicted patterns:  %s, %s", pattern, conflict)
		panic(s)
	}

	insert(t, parts, pattern, handlers)
}

func insert(t *trieNode, parts []string, pattern string, handlers HandlersChain) {
	if len(parts) == 0 {
		t.pattern = pattern
		t.handlers = handlers
		return
	}

	part := parts[0]
	child := matchChild(t, part)
	if child == nil {
		child = &trieNode{part: part}
		t.children = append(t.children, child)
	}

	insert(child, parts[1:], pattern, handlers)
}

func (t *trieNode) search(requestPath string) (HandlersChain, map[string]string) {
	requestPathParts := splitURLPath(requestPath)
	return search(t, requestPathParts)
}

func search(t *trieNode, parts []string) (HandlersChain, map[string]string) {
	var handlers HandlersChain
	var params = make(map[string]string)

	if len(parts) > 0 {
		part := parts[0]

		if strings.HasPrefix(t.part, "*") {
			if handlers, params = search(t, parts[1:]); handlers != nil {
				key := t.part[1:]
				if v, exist := params[key]; !exist {
					params[key] = part
				} else {
					params[key] = part + "/" + v
				}
				return handlers, params
			}
		}

		for _, child := range matchChildren(t, part) {
			if strings.HasPrefix(child.part, "*") {
				if handlers, params = search(child, parts); handlers != nil {
					key := child.part[1:]
					if _, exist := params[key]; !exist {
						params[key] = ""
					}

					return handlers, params
				}
			} else {
				if handlers, params = search(child, parts[1:]); handlers != nil {
					if strings.HasPrefix(child.part, ":") {
						key := child.part[1:]
						params[key] = part
					}

					return handlers, params
				}
			}
		}
	} else {
		if t.handlers != nil && t.pattern != "" {
			return t.handlers, params
		}

		for _, child := range t.children {
			if strings.HasPrefix(child.part, "*") {
				if handlers, params = search(child, parts); handlers != nil {
					key := child.part[1:]
					params[key] = ""
					return handlers, params
				}
			}
		}
	}

	return handlers, params
}

func matchChild(t *trieNode, part string) *trieNode {
	for _, child := range t.children {
		if child.part == part {
			return child
		}
	}

	return nil
}

func matchChildren(t *trieNode, part string) []*trieNode {
	var matchedChildren []*trieNode
	for _, child := range t.children {
		if child.part == part || strings.HasPrefix(child.part, "*") || strings.HasPrefix(child.part, ":") {
			matchedChildren = append(matchedChildren, child)
		}
	}

	return matchedChildren
}

func checkConflict(t *trieNode, parts []string) string {
	if len(parts) > 0 {
		part := parts[0]

		if strings.HasPrefix(t.part, "*") {
			if match := checkConflict(t, parts[1:]); match != "" {
				return match
			}
		}

		for _, child := range matchChildren(t, part) {
			if strings.HasPrefix(child.part, "*") {
				if match := checkConflict(child, parts); match != "" {
					return match
				}
			} else {
				if match := checkConflict(child, parts[1:]); match != "" {
					return match
				}
			}
		}

		// case for *wildcard
		if strings.HasPrefix(part, "*") {
			for _, child := range t.children {
				if match := checkConflict(child, parts); match != "" {
					return match
				}
			}

			if match := checkConflict(t, parts[1:]); match != "" {
				return match
			}
		}

		// case for :wildcard
		if strings.HasPrefix(part, ":") {
			for _, child := range t.children {
				if match := checkConflict(child, parts[1:]); match != "" {
					return match
				}
			}
		}
	} else {
		if t.pattern != "" {
			return t.pattern
		}

		for _, child := range t.children {
			if strings.HasPrefix(child.part, "*") {
				if match := checkConflict(child, parts); match != "" {
					return match
				}
			}
		}
	}

	return ""
}

func splitURLPath(pattern string) []string {
	pattern = strings.TrimPrefix(strings.TrimSuffix(pattern, "/"), "/")
	return strings.Split(pattern, "/")
}
