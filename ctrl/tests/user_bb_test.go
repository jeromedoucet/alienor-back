package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"encoding/json"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/jeromedoucet/alienor-back/utils"
)

// nominal user creation test case
func TestUserCreationSuccessful(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	utils.Clean([]string{"user:" + "leroy.jenkins"})
	usr := model.User{Identifier:"leroy.jenkins",
		Type:model.USER,
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
		}

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := utils.DoReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	// http res check
	var userRes model.User
	json.NewDecoder(res.Body).Decode(&userRes)
	assert.Equal(t, usr, userRes)

	// check db
	assert.Equal(t, usr, *utils.GetUser(usr.Identifier))
}

func TestUserCreationMalFormedJson(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()
	body := []byte("a malformed json")

	// when
	_, err := utils.DoReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
}

// already used identifier
func TestUserCreationExistingIdentifier(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	usr := model.User{Identifier:"leroy.jenkins",
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}
	utils.Populate(map[string]interface{}{"user:" + usr.Identifier: model.User{Identifier:usr.Identifier}})

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := utils.DoReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 409, res.StatusCode)
}

// when Identifier is missing
func TestUserCreationMissingIdentifier(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	usr := model.User{ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := utils.DoReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

// when ForName is missing
func TestUserCreationMissingForName(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	usr := model.User{Identifier:"leroy.jenkins",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := utils.DoReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

// when forName is missing
func TestUserCreationMissingName(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	usr := model.User{Identifier:"leroy.jenkins",
		ForName:"Leroy",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := utils.DoReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

// when email is missing
func TestUserCreationMissingEmail(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	usr := model.User{Identifier:"leroy.jenkins",
		ForName:"Leroy",
		Name:"Jenkins",
		Password:[]byte("wipe"),
	}

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := utils.DoReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

// when password is missing
func TestUserCreationMissingPassword(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	usr := model.User{Identifier:"leroy.jenkins",
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
	}

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := utils.DoReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

