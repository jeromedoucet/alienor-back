package rep_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/utils"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/rep"
	"github.com/stretchr/testify/assert"
)

func TestGetUserShouldUserWithSuccess(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	defer utils.Clean([]string{"user:" + "leroy.jenkins"})
	utils.Clean([]string{"user:" + "leroy.jenkins"})
	usr := model.User{Identifier:"leroy.jenkins",
		Type:model.USER,
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}
	utils.Populate(map[string]interface{}{"user:" + usr.Identifier: usr})
	rep.InitRepo(utils.CouchBaseAddr, "")

	// when
	actualUser, cas := rep.GetUser(usr.Identifier)

	// then
	assert.NotNil(t, cas)
	assert.Equal(t, usr.Email, actualUser.Email)
	assert.Equal(t, usr.ForName, actualUser.ForName)
	assert.Equal(t, usr.Name, actualUser.Name)
	assert.Equal(t, string(actualUser.Password), string(usr.Password))
}

func TestGetUserShouldUserWithError(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	defer utils.Clean([]string{"user:" + "leroy.jenkins"})
	rep.InitRepo(utils.CouchBaseAddr, "")

	// when
	usr, _ := rep.GetUser("leroy.jenkins")

	// then
	assert.Nil(t, usr)
}
