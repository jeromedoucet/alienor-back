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
	defer utils.Clean()
	usr := model.User{Identifier:"leroy.jenkins",
		Type:model.USER,
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
	}
	utils.Populate(map[string]interface{}{"user:" + usr.Identifier: usr})
	rep.InitRepo(utils.CouchBaseAddr, "")
	actualUser := model.NewUser()
	userRepository := new(rep.UserRepository)
	// when
	cas, err := userRepository.Get(usr.Identifier, actualUser)

	// then
	assert.Nil(t, err)
	assert.NotNil(t, cas)
	assert.Equal(t, usr.Email, actualUser.Email)
	assert.Equal(t, usr.ForName, actualUser.ForName)
	assert.Equal(t, usr.Name, actualUser.Name)
	assert.Equal(t, string(actualUser.Password), string(usr.Password))
}

func TestGetUserShouldUserWithError(t *testing.T) {
	// given
	utils.Before()
	defer utils.Clean()
	rep.InitRepo(utils.CouchBaseAddr, "")
	actualUser := model.NewUser()
	userRepository := new(rep.UserRepository)

	// when
	_, err := userRepository.Get("leroy.jenkins", actualUser)

	// then
	assert.NotNil(t, err)
}

func TestInsertANonUserEntity(t *testing.T) {
	// given
	utils.Before()
	defer utils.Clean()
	userRepository := new(rep.UserRepository)

	// when
	err := userRepository.Insert("A string is not a user :)")

	// then
	assert.NotNil(t, err)
	assert.Equal(t, "Cannot Insert a non user entity !", err.Error())
}
