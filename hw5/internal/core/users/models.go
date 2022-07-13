package users

type User struct {
	UserId int    `json:"userId,omitempty" param:"userId"`
	Name   string `json:"name,omitempty"`
	Age    int    `json:"age,omitempty"`
	Spouse int    `json:"spouse,omitempty"`
}
