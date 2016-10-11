package model

type Role string

const ADMIN Role = "ADMIN";
const TRANSLATOR Role = "TRANSLATOR";

// fixme : a role is meaningful only in a scope context
// fixme : add a technical id ? how to generate it ?
type User struct {
	Identifier string `json:"identifier"`
	ForName string `json:"forName"`
	Name string `json:"name"`
	Email string `json:"email"`
	Scope []string `json:"scope"`
	Roles []Role `json:"roles"`
	Password []byte `json:"password"`
}
