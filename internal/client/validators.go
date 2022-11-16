package client

import (
	"encoding/json"
	"net/mail"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

func any(input string) error {
	return nil
}

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

func isTags(input string) error {
	if input == "" {
		return nil
	}
	var tags map[string]string
	return json.Unmarshal([]byte(input), &tags)

}
