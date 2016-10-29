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
	"golang.org/x/crypto/bcrypt"
)

// nominal user creation test case
func TestUserCreationSuccessful(t *testing.T) {
	// given
	utils.Before()
	usr := model.User{Id:"leroy.jenkins",
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
	assert.Equal(t, 201, res.StatusCode)
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))

	// http res check
	var userRes model.User
	json.NewDecoder(res.Body).Decode(&userRes)
	assert.Empty(t, userRes.Password)
	assert.Equal(t, usr.Email, userRes.Email)
	assert.Equal(t, usr.ForName, userRes.ForName)
	assert.Equal(t, usr.Name, userRes.Name)
	// check db
	actualUser := utils.GetUser(usr.Id)
	assert.Equal(t, usr.Email, actualUser.Email)
	assert.Equal(t, usr.ForName, actualUser.ForName)
	assert.Equal(t, usr.Name, actualUser.Name)
	assert.Nil(t, bcrypt.CompareHashAndPassword(actualUser.Password, usr.Password))
}

func TestUserCreationMalFormedJson(t *testing.T) {
	// given
	utils.Before()
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
	usr := model.User{Id:"leroy.jenkins",
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}
	utils.Populate(map[string]interface{}{"user:" + usr.Id: model.User{Id:usr.Id}})

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
	usr := model.User{Id:"leroy.jenkins",
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
	usr := model.User{Id:"leroy.jenkins",
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
	usr := model.User{Id:"leroy.jenkins",
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
	usr := model.User{Id:"leroy.jenkins",
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

