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
	"github.com/jeromedoucet/alienor-back/utils"
)

func TestHandleAuthSuccess(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	utils.Clean([]string{"leroy.jenkins"})
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	usr := model.User{Identifier: login, Type:model.USER, Password: hPwd, }

	utils.Populate(map[string]interface{}{usr.Identifier: usr})

	s := utils.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login: login, Pwd: pwd})

	// when
	res, err := utils.DoReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	// check jwt token
	dec := json.NewDecoder(res.Body)
	var authRes ctrl.AuthRes
	dec.Decode(&authRes)
	_, jwtError := jwt.Parse(authRes.Token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		assert.Equal(t, true, ok)
		return []byte(utils.Secret), nil
	})
	assert.Nil(t, jwtError)
}

func TestHandleBadPassword(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	usr := model.User{Identifier: login, Type:model.USER, Password: hPwd}

	utils.Populate(map[string]interface{}{usr.Identifier: usr})

	s := utils.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login: login, Pwd: "roxx"})

	// when
	res, err := utils.DoReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

func TestHandleUnknownUser(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	utils.Clean([]string{"leroy.jenkins"})

	s := utils.StartHttp(func(r component.Router) {
		ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)
	})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login: "leroy.jenkins", Pwd: "test"})

	// when
	res, err := utils.DoReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
}
