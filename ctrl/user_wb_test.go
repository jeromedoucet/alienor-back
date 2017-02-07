package ctrl

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/model"
)

func TestCheckFieldShouldCheckIdFirst(t *testing.T) {
	// given
	t.Parallel()
	usr := model.NewUser()

	// when
	err := checkField(usr)

	// then
	if err == nil {
		t.Error("expected err not to be nil")
	} else if err.Error() != "invalid identifier" {
		t.Error("Bad error message")
	}
}

func TestCheckFieldShouldCheckForNameInSecondPlace(t *testing.T) {
	// given
	t.Parallel()
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"

	// when
	err := checkField(usr)

	// then
	if err == nil {
		t.Error("expected error not to be nil")
	} else if err.Error() != "invalid forname" {
		t.Error("Bad error message")
	}
}

func TestCheckFieldShouldCheckNameInThirdPlace(t *testing.T) {
	// given
	t.Parallel()
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"
	usr.ForName = "illidan"

	// when
	err := checkField(usr)

	// then
	if err == nil {
		t.Error("expected error not to be nil")
	} else if err.Error() != "invalid name" {
		t.Error("Bad error message")
	}
}

func TestCheckFieldShouldCheckEmailInFourthPlace(t *testing.T) {
	// given
	t.Parallel()
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"
	usr.ForName = "illidan"
	usr.Name = "stormrage"

	// when
	err := checkField(usr)

	// then
	if err == nil {
		t.Error("expected error not to be nil")
	} else if err.Error() != "invalid email" {
		t.Error("Bad Error message")
	}
}

func TestCheckFieldShouldCheckPasswordInFifthPlace(t *testing.T) {
	// given
	t.Parallel()
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"
	usr.ForName = "illidan"
	usr.Name = "stormrage"
	usr.Email = "illidan.storage@kalimdor.org"

	// when
	err := checkField(usr)

	// then
	if err == nil {
		t.Error("expected error not to be nil")
	} else if err.Error() != "invalid password" {
		t.Error("Bad Error message")
	}
}

func TestCheckFieldShouldPassCheck(t *testing.T) {
	// given
	t.Parallel()
	usr := model.NewUser()
	usr.Id = "illidan.stormrage"
	usr.ForName = "illidan"
	usr.Name = "stormrage"
	usr.Email = "illidan.storage@kalimdor.org"
	usr.Password = "somePasswordV1"
	// when
	err := checkField(usr)

	// then
	if err != nil {
		t.Error("expected error to be nil")
	}
}
