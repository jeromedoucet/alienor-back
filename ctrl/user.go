package ctrl

import (
	"github.com/jeromedoucet/alienor-back/component"
	"net/http"
	"github.com/jeromedoucet/alienor-back/model"
	"encoding/json"
	"errors"
)

// user handler
func handleUser(w http.ResponseWriter, r *http.Request) {
	usr, err := doCreateUser(r)
	if err != nil {
		writeError(w, err)
	} else {
		usr.Password = []byte{} // don't send the password !
		writeJsonResponse(w, usr, 201)
	}
}

// todo test me unit style !
// create a user using the request. Return an ctrlError if
// on issue
func doCreateUser(r *http.Request) (usr *model.User, cError *ctrlError) {
	usr = model.NewUser()
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(usr)
	if err != nil {
		cError = &ctrlError{httpCode:400, errorMsg:"Error during decoding the user creation request body"}
		return
	}
	if checkField(usr) != nil {
		cError = &ctrlError{httpCode:400, errorMsg:`Error during parsing the user creation request body
		 : there is missing fields`}
		return
	}
	err = userRepository.Insert(usr)
	if err != nil {
		cError = &ctrlError{httpCode:409, errorMsg:`Error during creating a new user
		 : user already exist`}
	}
	return
}

// todo test me unit style !
// check the user fields
func checkField(usr *model.User) error {
	if usr.Id == "" {
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
