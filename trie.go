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
	var params map[string]string
	var childrenParams map[string]string
	var handlers HandlersChain
	var lead *trieNode

	if len(parts) == 0 {
		return t.handlers, params
	}

	part := parts[0]
	if strings.HasPrefix(t.part, "*") {
		if handlers, childrenParams = search(t, parts[1:]); handlers != nil {
			lead = t
		}
	}

	if handlers == nil {
		for _, child := range matchChildren(t, part) {
			if handlers, childrenParams = search(child, parts[1:]); handlers != nil {
				lead = child
				break
			}
		}
	}

	if lead != nil { // equal to `match != nil`
		params = make(map[string]string)
		for k, v := range childrenParams {
			params[k] = v
		}

		if strings.HasPrefix(lead.part, ":") {
			key := lead.part[1:]
			params[key] = part
		}

		if strings.HasPrefix(lead.part, "*") {
			key := lead.part[1:]
			if v, ok := childrenParams[key]; ok {
				params[key] = part + "/" + v
			} else {
				params[key] = part
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
	if len(parts) == 0 {
		return t.pattern
	}

	part := parts[0]
	if strings.HasPrefix(t.part, "*") {
		match := checkConflict(t, parts[1:])
		if match != "" {
			return match
		}
	}

	for _, child := range matchChildren(t, part) {
		match := checkConflict(child, parts[1:])
		if match != "" {
			return match
		}
	}

	// case for *wildcard
	if strings.HasPrefix(part, "*") {
		for _, child := range t.children {
			match := checkConflict(child, parts)
			if match != "" {
				return match
			}
		}

		for _, child := range t.children {
			match := checkConflict(child, parts[1:])
			if match != "" {
				return match
			}
		}
	}

	// case for :wildcard
	if strings.HasPrefix(part, ":") {
		for _, child := range t.children {
			match := checkConflict(child, parts[1:])
			if match != "" {
				return match
			}
		}
	}

	return ""
}

func splitURLPath(pattern string) []string {
	pattern = strings.TrimPrefix(strings.TrimSuffix(pattern, "/"), "/")
	return strings.Split(pattern, "/")
}
