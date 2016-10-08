package main

import (
	"net/http"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"log"
)

var redisAddr string = "192.168.99.100:6379"

func main() {
	m := http.NewServeMux();
	ctrl.InitAuth(m, redisAddr, "some secret") // todo generate secret
	err := http.ListenAndServe(":8080", m)
	if err != nil {
		log.Fatal(err.Error())
	}
}



