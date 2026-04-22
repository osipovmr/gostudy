package entity

type User struct {
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
