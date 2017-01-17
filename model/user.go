package model

import "github.com/satori/go.uuid"

type RoleVal string

const ADMIN RoleVal = "ADMIN";
const TRANSLATOR RoleVal = "TRANSLATOR";

// a team that is used to attach item
type Team struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

// create a new team and generate team ID
func NewTeam() *Team {
	t:= new(Team)
	t.Id = uuid.NewV4().String()
	return t
}

// a role with a team as context
type Role struct {
	Value RoleVal `json:"value"`
	Team *Team `json:"team"`
}

func NewRole() *Role {
	r := new(Role)
	return r
}

type User struct {
	Id       string `json:"identifier"`
	Type     DocType `json:"type"`
	ForName  string `json:"forName"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Roles    []*Role `json:"roles"`
	Password string `json:"password"`
}

func (u *User) Identifier() string {
	return u.Id
}

func NewUser() *User {
	u := new(User)
	u.Type = USER
	return u
}