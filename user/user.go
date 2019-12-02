package user

type User struct {
	ID string `json:"id"`
	Email string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Active bool `json:"active"`
}
