package customer

type Customer struct {
	Id       uint   `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
