package client

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"regexp"

	"github.com/zklevsha/gophkeeper/internal/errs"
)

func any(input string) error {
	return nil
}

func notEmpty(input string) error {
	if len(input) <= 0 {
		return errs.ErrEmptyInput
	}
	return nil
}

func isEmail(input string) error {
	_, err := mail.ParseAddress(input)
	if err != nil {
		return errs.ErrInvalidEmail
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
		return errs.ErrInvalidCardNumber
	}
	return nil
}

func isCardNumberOrEmpty(input string) error {
	if len(input) == 0 {
		return nil
	}
	return isCardNumber(input)
}

func isCardHolder(input string) error {
	matched, err := regexp.MatchString("[A-Z]+ [A-Z]+", input)
	if err != nil {
		return fmt.Errorf("matchString have returned an error: %s", err.Error())
	}
	if !matched {
		return errs.ErrInvalidCardHolder
	}
	return nil
}

func isCardHolderOrEmpty(input string) error {
	if len(input) == 0 {
		return nil
	}
	return isCardHolder(input)
}

func isCardExire(input string) error {
	matched, err := regexp.MatchString("^[0-1][1-2]/[0-9]{2}$", input)
	if err != nil {
		return fmt.Errorf("matchString have returned an error: %s", err.Error())
	}
	if !matched {
		return errs.ErrInvalidCardExpire
	}
	return nil
}

func isCardExpireOrEmpty(input string) error {
	if len(input) == 0 {
		return nil
	}
	return isCardExire(input)
}

func isCardCVC(input string) error {
	matched, err := regexp.MatchString("^[0-9]{3}$", input)
	if err != nil {
		return fmt.Errorf("matchString have returned an error: %s", err.Error())
	}
	if !matched {
		return errs.ErrInvalidCardCVV
	}
	return nil
}

func isCardCVCorEmpty(input string) error {
	if len(input) == 0 {
		return nil
	}
	return isCardCVC(input)
}
