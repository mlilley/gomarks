package services

import "errors"

var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

type ErrValidation struct {
	msg string
}
func (e *ErrValidation) Error() string {
	return e.msg
}
func NewValidationError(msg string) error {
	return &ErrValidation{msg: msg}
}

