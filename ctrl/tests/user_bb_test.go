package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"encoding/json"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/garyburd/redigo/redis"
)

// nominal user creation test case
func TestUserCreationSuccessful(t *testing.T) {
	// given
	usr := model.User{Identifier:"leroy.jenkins",
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
		}

	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, rAddr, secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := doReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	// http res check
	var userRes model.User
	json.NewDecoder(res.Body).Decode(&userRes)
	assert.Equal(t, usr, userRes)

	// check db
	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	var userDb model.User
	bUser, _ := c.Do("GET", usr.Identifier)
	json.Unmarshal(bUser.([]byte), &userDb)
	assert.Equal(t, usr, userDb)
}

func TestUserCreationMalFormedJson(t *testing.T) {
	// given
	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, rAddr, secret)})
	defer s.Close()
	body := []byte("a malformed json")

	// when
	_, err := doReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
}

func TestUserCreationRedisUnavailable(t *testing.T) {
	// given
	usr := model.User{Identifier:"leroy.jenkins",
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}

	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, "192.168.99.100:12345", secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := doReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 503, res.StatusCode)
}

// already used identifier
func TestUserCreationExistingIdentifier(t *testing.T) {
	// given
	usr := model.User{Identifier:"leroy.jenkins",
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}
	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	populate(c, map[string]interface{}{usr.Identifier: model.User{Identifier:usr.Identifier}})

	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, rAddr, secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := doReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 409, res.StatusCode)
}

// when some mandatory fields are missing missing
func TestUserCreationMissingIdentifier(t *testing.T) {
	// given
	usr := model.User{ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}

	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, rAddr, secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := doReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}


