package route

import (
	"net/http"
	"strings"
	"errors"
)

type node struct {
	handler  func(http.ResponseWriter, *http.Request)
	children map[string]*node
}

type DynamicRouter struct {
	root map[string]*node
}

func NewDynamicRouter() *DynamicRouter {
	r := new(DynamicRouter)
	r.root = make(map[string]*node)
	return r
}

func (r *DynamicRouter) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.registerHandler(splitPath(pattern), handler)
}

func (r *DynamicRouter) ServeHTTP(res http.ResponseWriter, req *http.Request) {

}

func (r *DynamicRouter) registerHandler(paths []string, handler func(http.ResponseWriter, *http.Request)) {
	// todo verifier si il existe deja (nil sur handler ?) => panic
	// todo handle paths empty
	// todo ajouter les verbes http
	children := r.root
	var n *node
	var ok bool
	for _, path := range paths {
		/*
		 * we only consider static and dynamic element of the path
		 * static are directly registered into the tree
		 * dynamic element start with ':' and are registered as it into the tree
		 * each node can only have one dynamic element
		 */
		if strings.HasPrefix(path, ":") {
			path = ":"
		}
		n, ok = children[path]
		if !ok {
			n = &node{}
			n.children = make(map[string]*node)
			children[path] = n
		}
		children = n.children
	}
	n.handler = handler
}

func (r *DynamicRouter) findEndpoint(req *http.Request) (n *node, err error) {
	// todo clean path
	// todo check url encode
	return parseTree(r.root, splitPath(req.URL.Path))
}

func splitPath(path string) []string {
	p := strings.TrimPrefix(path, "/")
	return strings.Split(strings.TrimSuffix(p, "/"), "/")
}

func parseTree(children map[string]*node, path []string) (*node, error) {
	n, ok := children[path[0]]
	if !ok {
		// if the path doesn't match a static value, try with dynamic
		n, ok = children[":"]
		if !ok {
			return n, errors.New("unknown path")
		}
	}
	if len(path) > 1 {
		return parseTree(n.children, path[1:])
	} else {
		return n, nil
	}
}
