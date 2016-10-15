package model

type DocType string

const TEAM DocType = "team"
const ROLE DocType = "role"
const USER DocType = "user"

type RoleVal string

const ADMIN RoleVal = "ADMIN";
const TRANSLATOR RoleVal = "TRANSLATOR";

// a team that is used to attach token
type Team struct {
	Identifier string `json:"identifier"`
	Type DocType `json:"type"`
	Name string `json:"name"`
}

// todo test me
// todo bench pointer vs value?
// Convenient to use because it will set the type
// of this document
func NewTeam() *Team {
	t:= new(Team)
	t.Type = TEAM
	return t
}

// a role with a team as context
type Role struct {
	Type DocType `json:"type"`
	Value RoleVal `json:"value"`
	Team *Team `json:"team"`
}

// todo test me
// todo bench pointer vs value?
// Convenient to use because it will set the type
// of this document
func NewRole() *Role {
	r := new(Role)
	r.Type = ROLE
	return r
}

type User struct {
	Identifier string `json:"identifier"`
	Type DocType `json:"type"`
	ForName string `json:"forName"`
	Name string `json:"name"`
	Email string `json:"email"`
	Roles []*Role `json:"roles"`
	Password []byte `json:"password"`
}

// todo test me
// todo bench pointer vs value?
// Convenient to use because it will set the type
// of this document
func NewUser() *User {
	u := new(User)
	u.Type = USER
	return u
}