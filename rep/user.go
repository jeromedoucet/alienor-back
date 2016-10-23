package rep

import (
	"github.com/jeromedoucet/alienor-back/model"
	"golang.org/x/crypto/bcrypt"
	"github.com/couchbase/gocb"
	"errors"
)


type UserRepository struct {

}

func (UserRepository) Get(identifier string, entity interface{}) (gocb.Cas, error) {
	user, isUser := entity.(*model.User)
	if !isUser { // todo test that
		return 0, errors.New("Cannot Get a non user entity !")
	}
	return bucket.Get(string(model.USER) + ":" + identifier, user)
}

// todo test me
func (UserRepository) Insert(entity interface{}) (err error) {
	user, isUser := entity.(*model.User)
	if !isUser { // todo test that
		return errors.New("Cannot Insert a non user entity !")
	}
	cPwd, _ := bcrypt.GenerateFromPassword(user.Password, bcrypt.DefaultCost) //todo handle error
	user.Password = cPwd
	_, err = bucket.Insert(string(model.USER) + ":" + user.Identifier, entity, 0)
	return
}

// todo test me
func (UserRepository) Update(entity interface{}, cas gocb.Cas) (err error) {
	user, isUser := entity.(*model.User)
	if !isUser { // todo test that
		return errors.New("Cannot Insert a non user entity !")
	}
	// todo check if possible to update partially the document instead
	_, err = bucket.Replace(string(model.USER) + ":" + user.Identifier, user, cas, 0)
	return
}
