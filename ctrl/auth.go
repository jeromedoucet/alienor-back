package ctrl

import (
	"net/http"
	"github.com/jeromedoucet/alienor-back/component/endpoint"
	"github.com/jeromedoucet/alienor-back/component/db"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/jeromedoucet/alienor-back/model/team"
	"golang.org/x/crypto/bcrypt"
)

// todo create a package view for req and res struct ?
type AuthReq struct {
	Login string `json:"login"`
	Pwd   string `json:"pwd"`
}

type authenticator struct {
	rAdr string
	conn db.Connector
}

func (a *authenticator) handle(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var req AuthReq
	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	c, cErr := a.conn.Connect("tcp", a.rAdr)
	if cErr != nil {
		w.WriteHeader(503)
		return
	}
	defer c.Close()
	exist, eErr := redis.Bool(c.Do("EXISTS", req.Login))
	if eErr != nil {
		w.WriteHeader(503)
		return
	}
	if exist {
		var user team.User
		bUser, _ := c.Do("GET", req.Login)
		json.Unmarshal(bUser.([]byte), &user)
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Pwd)) == nil {
			// todo use jwt
			w.WriteHeader(200)
		} else {
			w.WriteHeader(400)
		}

	} else {
		w.WriteHeader(404)
	}

}

func InitAuth(router endpoint.Router, rAdr string) {
	auth := authenticator{rAdr:rAdr, conn:db.NewConnector()}
	router.HandleFunc("/login", auth.handle)
	// todo return router ?
}