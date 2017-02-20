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
	// todo ajouter les verbes http
	if handler == nil {
		panic("handler cannot be nil")
	} else if len(paths) < 1 {
		panic("path cannot be nil")
	}
	children := r.root
	var n *node
	var ok bool
	for _, path := range paths {
		if path == "" {
			continue
		}
		/*
		 * we only consider static and dynamic identifier of the path.
		 *
		 * For static :
		 * If at a given non terminal node, the resource
		 * already exist, and if the identifier is static, we just
		 * pass to the next level.
		 *
		 * For dynamic :
		 * if the identifier of the resource is dynamic and if a
		 * dynamic identifier already exist with another name, the router will panic.
		 *
		 * Common :
		 * If the node denoted by the incoming path already has a handler, the router will panic
		 */
		if strings.HasPrefix(path, ":") {
			for m := range children {
				if strings.HasPrefix(m, ":") && path != m {
					panic("a dynamic identifier has already been registered at that level")
				}
			}
		}
		n, ok = children[path]
		if !ok {
			n = &node{}
			n.children = make(map[string]*node)
			children[path] = n
		}
		children = n.children
	}
	if n.handler != nil {
		panic("a handler is already registered for this path")
	}
	n.handler = handler
}

func (r *DynamicRouter) findEndpoint(req *http.Request) (n *node, err error) {
	// todo clean path
	// todo check url encoder
	return parseTree(r.root, splitPath(req.URL.Path))
}

func splitPath(path string) []string {
	p := strings.TrimPrefix(path, "/")
	return strings.Split(strings.TrimSuffix(p, "/"), "/")
}

func parseTree(children map[string]*node, path []string) (*node, error) {
	n, ok := children[path[0]]
	if !ok {
		return n, errors.New("unknown path")

	}
	if len(path) > 1 {
		return parseTree(n.children, path[1:])
	} else {
		return n, nil
	}
}
