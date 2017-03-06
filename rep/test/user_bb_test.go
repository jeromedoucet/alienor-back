package rep_test

import (
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/rep"
	"github.com/jeromedoucet/alienor-back/test"
	"testing"
)

func TestGetUserShouldUserWithSuccess(t *testing.T) {
	// given
	test.Before()
	usr := model.User{Id: "leroy.jenkins",
		Type:     model.USER,
		ForName:  "Leroy",
		Name:     "Jenkins",
		Email:    "leroy.jenkins@wipe-guild.org",
		Password: "wipe",
	}
	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})
	rep.InitRepo(test.CouchBaseAddr, "")
	actualUser := model.NewUser()
	userRepository := new(rep.UserRepository)
	// when
	_, err := userRepository.Get(usr.Id, actualUser)

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if actualUser.Email != usr.Email {
		t.Error("expect users email to be equals")
	} else if actualUser.ForName != usr.ForName {
		t.Error("expect users ForName to be equals")
	} else if actualUser.Name != usr.Name {
		t.Error("expect users Name to be equals")
	} else if actualUser.Password != usr.Password {
		t.Error("expect users Password to be equals")
	}
}

func TestGetUserShouldUserWithError(t *testing.T) {
	// given
	test.Before()
	rep.InitRepo(test.CouchBaseAddr, "")
	actualUser := model.NewUser()
	userRepository := new(rep.UserRepository)

	// when
	_, err := userRepository.Get("leroy.jenkins", actualUser)

	// then
	if err == nil {
		t.Error("expect error not to be nil")
	}
}

func TestInsertANonUserEntity(t *testing.T) {
	// given
	test.Before()
	userRepository := new(rep.UserRepository)

	// when
	err := userRepository.Insert(&test.MockDocument{Id: "someId"})

	// then
	if err == nil {
		t.Error("expect error not to be nil")
	} else if err.Error() != "Cannot Insert a non user entity" {
		t.Error("bad error message")
	}
}

func TestInsertUserWithoutPwd(t *testing.T) {
	// given
	test.Before()
	rep.InitRepo(test.CouchBaseAddr, "")
	userRepository := new(rep.UserRepository)
	usr := &model.User{Id: "leroy.jenkins",
		Type:    model.USER,
		ForName: "Leroy",
		Name:    "Jenkins",
		Email:   "leroy.jenkins@wipe-guild.org",
	}

	// when
	err := userRepository.Insert(usr)

	// then
	if err == nil {
		t.Error("expect error not to be nil")
	} else if err.Error() != "Cannot Insert a user without password!" {
		t.Error("bad error message")
	}
}
