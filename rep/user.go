package rep

import (
	"fmt"
	"github.com/jeromedoucet/alienor-back/model"
)

// todo test me
func GetUser(identifier string) (*model.User, error) {
	usr := model.NewUser()
	_, err := bucket.Get(identifier, usr)
	if err != nil {
		fmt.Println("Error returning the document:", err)
		return usr, err
	}
	return usr, nil
}

// todo test me
func InsertUser(usr *model.User) (err error) {
	_, err = bucket.Upsert(usr.Identifier, usr, 0)
	return
}
