package model

type Role string

const ADMIN Role = "ADMIN";
const TRANSLATOR Role = "TRANSLATOR";

type User struct {
	Identifier string `json:"identifier"`
	ForName string `json:"forName"`
	Name string `json:"name"`
	Email string `json:"email"`
	Scope []string `json:"scope"`
	Roles []Role `json:"roles"`
	Password []byte `json:"password"`
}
