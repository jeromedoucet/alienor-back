package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"encoding/json"
	"bytes"
	"github.com/jeromedoucet/alienor-back/test"
	"golang.org/x/crypto/bcrypt"
)

// nominal user creation test case
func TestUserCreationSuccessful(t *testing.T) {
	// given
	test.Before()
	usr := model.User{Id: "leroy.jenkins",
		Type:         model.USER,
		ForName:      "Leroy",
		Name:         "Jenkins",
		Email:        "leroy.jenkins@wipe-guild.org",
		Password:     "wipe",
	}

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := test.DoReq(s.URL+"/user", "POST", bytes.NewBuffer(body))

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if res.StatusCode != 201 {
		t.Error("expect the status code to equals 201")
	} else if res.Header.Get("Content-Type") != "application/json" {
		t.Error("expect the content type to equals 'application/json'")
	}

	// http res check
	var userRes model.User
	json.NewDecoder(res.Body).Decode(&userRes)
	if len(userRes.Password) != 0 {
		t.Error("expect the password to be empty")
	} else if userRes.Email != usr.Email {
		t.Error("expect the email to be the same")
	} else if usr.ForName != userRes.ForName {
		t.Error("expect the forName to be the same")
	} else if usr.Name != userRes.Name {
		t.Error("expect the name to be the same")
	}
	// check db
	actualUser := test.GetUser(usr.Id)
	if usr.Email != actualUser.Email {
		t.Error("expect the email to be the same")
	} else if usr.ForName != actualUser.ForName {
		t.Error("expect the forName to be the same")
	} else if usr.Name != actualUser.Name {
		t.Error("expect the name to be the same")
	} else if bcrypt.CompareHashAndPassword([]byte(actualUser.Password), []byte(usr.Password)) != nil {
		t.Error("expect the password to be the same")
	}
}

func TestUserCreationMalFormedJson(t *testing.T) {
	// given
	test.Before()
	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body := []byte("a malformed json")

	// when
	res, err := test.DoReq(s.URL+"/user", "POST", bytes.NewBuffer(body))
	var resBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resBody)

	// then
	if err != nil {
		t.Error("expect the error to be nil")
	}
	if res.StatusCode != 400 {
		t.Error("expect the status code to equals 400")
	} else if resBody.Msg != "Error during decoding the user creation request body" {
		t.Error("expect the body msg to eqauls 'Error during decoding the user creation request body'")
	}
}

// already used identifier
func TestUserCreationExistingIdentifier(t *testing.T) {
	// given
	test.Before()
	usr := model.User{Id: "leroy.jenkins",
		ForName:      "Leroy",
		Name:         "Jenkins",
		Email:        "leroy.jenkins@wipe-guild.org",
		Password:     "wipe",
	}
	test.Populate(map[string]interface{}{"user:" + usr.Id: model.User{Id: usr.Id}})

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := test.DoReq(s.URL+"/user", "POST", bytes.NewBuffer(body))
	var resBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resBody)

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if res.StatusCode != 409 {
		t.Error("expect the status code to equals 409")
	} else if resBody.Msg != "Error during creating a new user : user already exist" {
		t.Error("expect the body msg to equals 'Error during creating a new user : user already exist'")
	}
}

// when Identifier is missing
func TestUserCreationMissingIdentifier(t *testing.T) {
	// given
	test.Before()
	usr := model.User{ForName: "Leroy",
		Name:              "Jenkins",
		Email:             "leroy.jenkins@wipe-guild.org",
		Password:          "wipe",
	}

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := test.DoReq(s.URL+"/user", "POST", bytes.NewBuffer(body))
	var resBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resBody)

	// then
	if err != nil {
		t.Error("expec the error to be nil")
	} else if res.StatusCode != 400 {
		t.Error("expect the status code to equals 400")
	} else if resBody.Msg != "invalid identifier" {
		t.Error("expect the body msg to equals 'invalid identifier'")
	}
}

// when ForName is missing
func TestUserCreationMissingForName(t *testing.T) {
	// given
	test.Before()
	usr := model.User{Id: "leroy.jenkins",
		Name:         "Jenkins",
		Email:        "leroy.jenkins@wipe-guild.org",
		Password:     "wipe",
	}

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := test.DoReq(s.URL+"/user", "POST", bytes.NewBuffer(body))
	var resBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resBody)

	// then
	if err != nil {
		t.Error("expect the error to be nil")
	} else if res.StatusCode != 400 {
		t.Error("expect the status code to equals 400")
	} else if resBody.Msg != "invalid forname" {
		t.Error("body msg to equals 'invalid forname'")
	}
}

// when forName is missing
func TestUserCreationMissingName(t *testing.T) {
	// given
	test.Before()
	usr := model.User{Id: "leroy.jenkins",
		ForName:      "Leroy",
		Email:        "leroy.jenkins@wipe-guild.org",
		Password:     "wipe",
	}

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := test.DoReq(s.URL+"/user", "POST", bytes.NewBuffer(body))
	var resBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resBody)

	// then
	if err != nil {
		t.Error("expect the error to be nil")
	} else if res.StatusCode != 400 {
		t.Error("expect the status code to equals 400")
	} else if resBody.Msg != "invalid name" {
		t.Error("expect the body msg to equals 'invalid name'")
	}
}

// when email is missing
func TestUserCreationMissingEmail(t *testing.T) {
	// given
	test.Before()
	usr := model.User{Id: "leroy.jenkins",
		ForName:      "Leroy",
		Name:         "Jenkins",
		Password:     "wipe",
	}

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := test.DoReq(s.URL+"/user", "POST", bytes.NewBuffer(body))
	var resBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resBody)

	// then
	if err != nil {
		t.Error("expect the error to be nil")
	} else if res.StatusCode != 400 {
		t.Error("expect the status code to equals 400")
	} else if resBody.Msg != "invalid email" {
		t.Error("expect the body msg to equals 'invalid email'")
	}
}

// when password is missing
func TestUserCreationMissingPassword(t *testing.T) {
	// given
	test.Before()
	usr := model.User{Id: "leroy.jenkins",
		ForName:      "Leroy",
		Name:         "Jenkins",
		Email:        "leroy.jenkins@wipe-guild.org",
	}

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := test.DoReq(s.URL+"/user", "POST", bytes.NewBuffer(body))
	var resBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resBody)

	// then
	if err != nil {
		t.Error("expect he error to be nil")
	} else if res.StatusCode != 400 {
		t.Error("expect the status code to equals 400")
	} else if resBody.Msg != "invalid password" {
		t.Error("expect the body msg to equals 'invalid password'")
	}
}
