package client

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"regexp"

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

func isCardNumber(input string) error {
	matched, err := regexp.MatchString("^[0-9]{4} [0-9]{4} [0-9]{4} [0-9]{4}$", input)
	if err != nil {
		return fmt.Errorf("matchString have returned an error: %s", err.Error())
	}
	if !matched {
		return structs.ErrInvalidCardNumber
	}
	return nil
}

func isCardHolder(input string) error {
	matched, err := regexp.MatchString("[A-Z]+ [A-Z]+", input)
	if err != nil {
		return fmt.Errorf("matchString have returned an error: %s", err.Error())
	}
	if !matched {
		return structs.ErrInvalidCardHolder
	}
	return nil
}

func isCardExire(input string) error {
	matched, err := regexp.MatchString("^[0-1][1-2]/[0-9]{2}$", input)
	if err != nil {
		return fmt.Errorf("matchString have returned an error: %s", err.Error())
	}
	if !matched {
		return structs.ErrInvalidCardExpire
	}
	return nil
}

func isCardCVC(input string) error {
	matched, err := regexp.MatchString("^[0-9]{3}$", input)
	if err != nil {
		return fmt.Errorf("matchString have returned an error: %s", err.Error())
	}
	if !matched {
		return structs.ErrInvalidCardCVV
	}
	return nil
}
