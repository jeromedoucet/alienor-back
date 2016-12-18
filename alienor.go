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
	// the CORS request are here allowed
	//if origin := req.Header.Get("Origin"); origin != "" {
	//	res.Header().Set("Access-Control-Allow-Origin", origin)
	//	res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//	res.Header().Set("Access-Control-Allow-Headers",
	//		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//}
	//if req.Method == "OPTIONS" {
	//	return
	//}
	f.mux.ServeHTTP(res, req)
}



