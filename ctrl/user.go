package ctrl

import (
	"github.com/jeromedoucet/alienor-back/component"
	"net/http"
	"github.com/jeromedoucet/alienor-back/model"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"errors"
)

// user handler
func handleUser(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var usr model.User
	unmarshalError := dec.Decode(&usr)

	if unmarshalError != nil {
		w.WriteHeader(400)
		return
	}
	if checkField(&usr) != nil {
		w.WriteHeader(400)
		return
	}
	c, connError := conn.Connect("tcp", rAdr)

	if connError != nil {
		w.WriteHeader(503)
		return
	}
	exist, _ := redis.Bool(c.Do("EXISTS", usr.Identifier))
	if exist {
		w.WriteHeader(409)
		return
	}
	usrToSave, _ := json.Marshal(usr)
	c.Do("SET", usr.Identifier, string(usrToSave))
	w.Write(usrToSave)
}

func checkField(usr *model.User) error {
	if usr.Identifier == "" {
		return errors.New("invalid identifier")
	}
	if usr.ForName == "" {
		return errors.New("invalid forname")
	}
	if usr.Name == "" {
		return errors.New("invalid name")
	}
	if usr.Email == "" {
		return errors.New("invalid email")
	}
	if len(usr.Password) < 1 {
		return errors.New("invalid password")
	}
	return nil
}

func initUserEndPoint(router component.Router) {
	router.HandleFunc(UserEndPoint, handleUser)
}
