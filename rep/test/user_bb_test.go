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
	expectedUsr, err := rep.GetUser(usr.Identifier)

	// then
	assert.Nil(t ,err)
	assert.Equal(t, usr, *expectedUsr)
}

func TestGetUserShouldUserWithError(t *testing.T) {
	// given
	utils.Before()
	defer utils.After()
	rep.InitRepo(utils.CouchBaseAddr, "")

	// when
	_, err := rep.GetUser("leroy.jenkins")

	// then
	assert.NotNil(t ,err)
}
