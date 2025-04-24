package storage

type User struct {
	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
