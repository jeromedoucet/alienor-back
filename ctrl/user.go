package ctrl

import (
	"github.com/jeromedoucet/alienor-back/component"
	"net/http"
	"github.com/jeromedoucet/alienor-back/model"
	"encoding/json"
	"errors"
	"github.com/jeromedoucet/alienor-back/rep"
	"fmt"
)

// user handler
func handleUser(w http.ResponseWriter, r *http.Request) {
	usr := model.NewUser()
	var err error
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(usr)

	if err != nil {
		w.WriteHeader(400)
		return
	}
	if checkField(usr) != nil {
		w.WriteHeader(400)
		return
	}
	err = rep.InsertUser(usr)
	if err != nil {
		w.WriteHeader(409)
		return
	}
	usr.Password = []byte{} // don't send the password !
	usrSaved, _ := json.Marshal(usr)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", usrSaved)
}

// check the user fields
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
	router.HandleFunc(USER_ENDPOINT, handleUser)
}
