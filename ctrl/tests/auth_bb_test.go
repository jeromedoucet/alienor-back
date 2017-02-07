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
	usr := model.User{Id: login, Type:model.USER, Password: string(hPwd)}

	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login: login, Pwd: pwd})

	// when
	res, err := test.DoReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	// check jwt token
	assert.Len(t, res.Cookies(), 1)
	cookie := res.Cookies()[0]
	assert.True(t, cookie.HttpOnly)
	assert.Equal(t, "ALIENOR_SESS", cookie.Name)
	jwtToken := cookie.Value
	_, jwtError := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		assert.Equal(t, true, ok)
		return []byte(test.Secret), nil
	})
	assert.Nil(t, jwtError)
}

func TestHandleBadPassword(t *testing.T) {
	// given
	test.Before()
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	usr := model.User{Id: login, Type:model.USER, Password: string(hPwd)}

	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login: login, Pwd: "roxx"})

	// when
	res, err := test.DoReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
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
	res, err := test.DoReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
}