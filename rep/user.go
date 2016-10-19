package rep

import (
	"fmt"
	"github.com/jeromedoucet/alienor-back/model"
	"golang.org/x/crypto/bcrypt"
	"github.com/couchbase/gocb"
)

func GetUser(identifier string) (*model.User, gocb.Cas) {
	usr := model.NewUser()
	cas, err := bucket.Get(string(model.USER) + ":" + identifier, usr)
	if err != nil {
		fmt.Println("Error returning the document:", err)
		return nil, cas
	}
	return usr, cas
}

// todo test me
func InsertUser(usr *model.User) (err error) {
	cPwd, _ := bcrypt.GenerateFromPassword(usr.Password, bcrypt.DefaultCost) //todo handle error
	usr.Password = cPwd
	_, err = bucket.Insert(string(model.USER) + ":" + usr.Identifier, usr, 0)
	return
}

// todo test me
func UpdateUser(user *model.User, cas gocb.Cas) (err error) {
	_, err = bucket.Replace(string(model.USER) + ":" + user.Identifier, user, cas, 0)
	return
}
