package ctrl_test

import (
	"bytes"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"github.com/jeromedoucet/alienor-back/test"
)

// todo => check the error msg !

func TestHandleAuthSuccess(t *testing.T) {
	// given
	test.Before()
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	usr := model.User{Id: login, Type: model.USER, Password: string(hPwd)}

	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login: login, Pwd: pwd})

	// when
	res, err := test.DoReq(s.URL+"/login", "POST", bytes.NewBuffer(body))

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if res.StatusCode != 200 {
		t.Error("expect status code to be equals to 200")
	} else if len(res.Cookies()) != 1 {
		t.Error("expect to have only one cookie")
	}
	cookie := res.Cookies()[0]
	if !cookie.HttpOnly {
		t.Error("expect cookie to be http only")
	} else if cookie.Name != "ALIENOR_SESS" {
		t.Error("expect cookie name to be 'ALIENOR_SESS'")
	}
	jwtToken := cookie.Value
	_, jwtError := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		assert.Equal(t, true, ok)
		return []byte(test.Secret), nil
	})
	if jwtError != nil {
		t.Error("expect token to be successfuly unmarshall")
	}
}

func TestHandleBadPassword(t *testing.T) {
	// given
	test.Before()
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	usr := model.User{Id: login, Type: model.USER, Password: string(hPwd)}

	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login: login, Pwd: "roxx"})

	// when
	res, err := test.DoReq(s.URL+"/login", "POST", bytes.NewBuffer(body))

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if res.StatusCode != 400 {
		t.Error("expect status code to be equals to 400")
	}
}

func TestHandleUnknownUser(t *testing.T) {
	// given
	test.Before()

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login: "leroy.jenkins", Pwd: "test"})

	// when
	res, err := test.DoReq(s.URL+"/login", "POST", bytes.NewBuffer(body))

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if res.StatusCode != 404 {
		t.Error("expect status code to equals 404")
	}
}
