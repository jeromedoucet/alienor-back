package main

import (
	"net/http"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"log"
	"flag"
)

var dataStoreAddr string = "127.0.0.1"

// todo make some properties external
func main() {
	httpAddr := flag.String("http", ":8080", "the address on which the http server will listen")
	webRoot := flag.String("root", "./dist", "the repository from which the web client resources are served")
	flag.Parse()
	m := http.NewServeMux();
	ctrl.InitEndPoints(m, dataStoreAddr, "", "some secret") // todo generate secret
	serveStaticFiles(*webRoot, m)
	f := new(httpFilter)
	f.mux = m
	err := http.ListenAndServe(*httpAddr, f)
	if err != nil {
		log.Fatal(err.Error())
	}
}

type httpFilter struct {
	mux *http.ServeMux
}

// todo test me !
func (f *httpFilter) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	f.mux.ServeHTTP(res, req)
}



