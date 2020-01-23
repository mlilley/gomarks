package app

import (
	"errors"
)

var ErrNotFound = errors.New("not found")


type errAppError struct {
	msg string
}
func (e *errAppError) Error() string {
	return e.msg
}
func NewAppError(msg string) error {
	return &errAppError{msg: msg}
}

type Mark struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	UserID string `json:"user_id"`
}

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Active       bool   `json:"active"`
}

type Device struct {
	UserId    string   `json:"user_id"`
	DeviceId  string   `json:"device_id"`
	TokenHash [32]byte `json:"token_hash"`
}
