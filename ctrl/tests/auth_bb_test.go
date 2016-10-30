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
	"fmt"
)


// todo => check the error msg !

func TestHandleAuthSuccess(t *testing.T) {
	// given
	utils.Before()
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	usr := model.User{Id: login, Type:model.USER, Password: string(hPwd)}

	utils.Populate(map[string]interface{}{"user:" + usr.Id: usr})

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
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	usr := model.User{Id: login, Type:model.USER, Password: string(hPwd)}

	utils.Populate(map[string]interface{}{"user:" + usr.Id: usr})

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

func TestToto(t *testing.T) {
	// given

	pwd, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
	hash := string(pwd)

	fmt.Println(hash)

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("test"))

	assert.Nil(t, err)

}