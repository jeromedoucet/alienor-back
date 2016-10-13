package ctrl_test

import (
	"encoding/json"
	"testing"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/bcrypt"
	"github.com/stretchr/testify/assert"
	"bytes"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/alienor-back/component"
)

func TestHandleAuthSuccess(t *testing.T) {
	// given
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost);
	usr := model.User{Identifier:login, Password:hPwd, Roles:[]model.Role{model.TRANSLATOR}}

	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	populate(c, map[string]interface{}{usr.Identifier: usr})

	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, rAddr, secret)})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login:login, Pwd:pwd})

	// when
	res, err := doReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	// check jwt token
	dec := json.NewDecoder(res.Body)
	var authRes ctrl.AuthRes
	dec.Decode(&authRes)
	_, jwtError := jwt.Parse(authRes.Token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC);
		assert.Equal(t, true, ok)
		return []byte(secret), nil
	})
	assert.Nil(t, jwtError)
}

func TestHandleRedisConError(t *testing.T) {
	// given
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost);
	usr := model.User{Identifier:login, Password:hPwd}

	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	populate(c, map[string]interface{}{usr.Identifier: usr})

	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, "192.168.99.100:1234", secret)})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login:login, Pwd:pwd})

	// when
	res, err := doReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 503, res.StatusCode)
}

func TestHandleBadPassword(t *testing.T) {
	// given
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost);
	usr := model.User{Identifier:login, Password:hPwd}

	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	populate(c, map[string]interface{}{usr.Identifier: usr})

	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, rAddr, secret)})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login:login, Pwd:"roxx"})

	// when
	res, err := doReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

func TestHandleUnknownUser(t *testing.T) {
	// given
	pwd := "lichking"
	login := "arthas.menethil"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost);
	usr := model.User{Identifier:login, Password:hPwd}

	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	populate(c, map[string]interface{}{usr.Identifier: usr})

	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, rAddr, secret)})
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login:"leroy.jenkins", Pwd:"test"})

	// when
	res, err := doReq(s.URL + "/login", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
}