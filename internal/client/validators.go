package client

import (
	"net/mail"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

func notEmpty(input string) error {
	if len(input) <= 0 {
		return structs.ErrEmptyInput
	}
	return nil
}

func isEmail(input string) error {
	_, err := mail.ParseAddress(input)
	if err != nil {
		return structs.ErrInvalidEmail
	}
	return nil
}
