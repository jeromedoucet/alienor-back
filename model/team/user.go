package team

type User struct {
	Identifier string `json:"identifier"`
	ForName string `json:"forName"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password []byte `json:"password"`
}
