package rep

import (
	"github.com/jeromedoucet/alienor-back/model"
	"golang.org/x/crypto/bcrypt"
	"github.com/couchbase/gocb"
	"errors"
)


type UserRepository struct {

}

func (UserRepository) Get(identifier string, document model.Document) (gocb.Cas, error) {
	user, isUser := document.(*model.User)
	if !isUser { // todo test that
		return 0, errors.New("Cannot Get a non user entity !")
	}
	return bucket.Get(string(model.USER) + ":" + identifier, user)
}

func (UserRepository) Insert(document model.Document) (err error) {
	user, isUser := document.(*model.User)
	if !isUser {
		err = errors.New("Cannot Insert a non user entity !")
		return
	}
	if len(user.Password) == 0 {
		err = errors.New("Cannot Insert a user without password!")
		return
	}
	cPwd, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost) //todo handle error
	user.Password = string(cPwd)
	_, err = bucket.Insert(string(model.USER) + ":" + user.Id, user, 0)
	return
}

// todo test me
func (UserRepository) Update(document model.Document, cas gocb.Cas) (err error) {
	user, isUser := document.(*model.User)
	if !isUser { // todo test that
		err = errors.New("Cannot Insert a non user entity !")
		return
	}
	// todo check if possible to update partially the document instead
	_, err = bucket.Replace(string(model.USER) + ":" + user.Id, user, cas, 0)
	return
}
