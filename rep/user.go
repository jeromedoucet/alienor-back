package rep

import (
	"fmt"
	"github.com/jeromedoucet/alienor-back/model"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(identifier string) (*model.User, error) {
	usr := model.NewUser()
	_, err := bucket.Get(string(model.USER) + ":" + identifier, usr)
	if err != nil {
		fmt.Println("Error returning the document:", err)
		return usr, err
	}
	return usr, nil
}

// todo test me
func InsertUser(usr *model.User) (err error) {
	cPwd, _ := bcrypt.GenerateFromPassword(usr.Password, bcrypt.DefaultCost) //todo handle error
	usr.Password = cPwd
	_, err = bucket.Upsert(string(model.USER) + ":" + usr.Identifier, usr, 0)
	return
}
