package main

import "net/http"

// todo test me
func serveStaticFiles(root string, router *http.ServeMux) {
	fs := http.FileServer(http.Dir(root))
	router.Handle("/", fs)
}
