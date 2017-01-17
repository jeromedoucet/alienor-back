package ctrl

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/stretchr/testify/assert"
)

func TestCheckFieldShouldCheckIdFirst(t *testing.T) {
	// given
	usr := model.NewUser()

	// when
	err := checkField(usr)

	// then
	assert.NotNil(t, err)
	assert.Equal(t, "invalid identifier", err.Error())
}

func TestCheckFieldShouldCheckForNameInSecondPlace(t *testing.T) {
	// given
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"

	// when
	err := checkField(usr)

	// then
	assert.NotNil(t, err)
	assert.Equal(t, "invalid forname", err.Error())
}

func TestCheckFieldShouldCheckNameInThirdPlace(t *testing.T) {
	// given
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"
	usr.ForName = "illidan"

	// when
	err := checkField(usr)

	// then
	assert.NotNil(t, err)
	assert.Equal(t, "invalid name", err.Error())
}

func TestCheckFieldShouldCheckEmailInFourthPlace(t *testing.T) {
	// given
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"
	usr.ForName = "illidan"
	usr.Name = "stormrage"

	// when
	err := checkField(usr)

	// then
	assert.NotNil(t, err)
	assert.Equal(t, "invalid email", err.Error())
}

func TestCheckFieldShouldCheckPasswordInFifthPlace(t *testing.T) {
	// given
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"
	usr.ForName = "illidan"
	usr.Name = "stormrage"
	usr.Email = "illidan.storage@kalimdor.org"

	// when
	err := checkField(usr)

	// then
	assert.NotNil(t, err)
	assert.Equal(t, "invalid password", err.Error())
}

func TestCheckFieldShouldPassCheck(t *testing.T) {
	// given
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"
	usr.ForName = "illidan"
	usr.Name = "stormrage"
	usr.Email = "illidan.storage@kalimdor.org"
	usr.Password = "somePasswordV1"
	// when
	err := checkField(usr)

	// then
	assert.Nil(t, err)
}