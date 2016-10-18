package ctrl

import (
	"github.com/jeromedoucet/alienor-back/component"
	"net/http"
	"github.com/jeromedoucet/alienor-back/model"
	"encoding/json"
	"errors"
	"github.com/jeromedoucet/alienor-back/rep"
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
	_, err = rep.GetUser(usr.Identifier)
	// when error, then the key is not found
	if err == nil {
		w.WriteHeader(409)
		return
	}
	err = rep.InsertUser(usr)
	// todo test me
	if err != nil {
		w.WriteHeader(500)
		return
	}
	usr.Password = []byte{} // don't send the password !
	usrToSave, _ := json.Marshal(usr)
	w.Write(usrToSave)
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
