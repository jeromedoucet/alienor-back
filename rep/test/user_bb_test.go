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
	usr := model.User{Id:"leroy.jenkins",
		Type:model.USER,
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:"wipe",
	}
	utils.Populate(map[string]interface{}{"user:" + usr.Id: usr})
	rep.InitRepo(utils.CouchBaseAddr, "")
	actualUser := model.NewUser()
	userRepository := new(rep.UserRepository)
	// when
	cas, err := userRepository.Get(usr.Id, actualUser)

	// then
	assert.Nil(t, err)
	assert.NotNil(t, cas)
	assert.Equal(t, usr.Email, actualUser.Email)
	assert.Equal(t, usr.ForName, actualUser.ForName)
	assert.Equal(t, usr.Name, actualUser.Name)
	assert.Equal(t, actualUser.Password, usr.Password)
}

func TestGetUserShouldUserWithError(t *testing.T) {
	// given
	utils.Before()
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
	userRepository := new(rep.UserRepository)

	// when
	err := userRepository.Insert(&utils.MockDocument{Id:"someId"})

	// then
	assert.NotNil(t, err)
	assert.Equal(t, "Cannot Insert a non user entity !", err.Error())
}

func TestInsertUserWithoutPwd(t *testing.T) {
	// given
	utils.Before()
	rep.InitRepo(utils.CouchBaseAddr, "")
	userRepository := new(rep.UserRepository)
	usr := &model.User{Id:"leroy.jenkins",
		Type:model.USER,
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
	}

	// when
	err := userRepository.Insert(usr)

	// then
	assert.NotNil(t, err)
	assert.Equal(t, "Cannot Insert a user without password!", err.Error())
}
