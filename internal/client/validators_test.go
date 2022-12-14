package client

import (
	"testing"

	"github.com/zklevsha/gophkeeper/internal/errs"
)

func TestIsEmpty(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "input empty", input: "",
			want: errs.ErrEmptyInput},
		{name: "input is not empty", input: "string",
			want: nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := notEmpty(tc.input)
			if res != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", res, tc.want)
			}
		})
	}
}

func TestIsEmail(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "good email", input: "test@test.ru",
			want: nil},
		{name: "bad email b@ad", input: "string",
			want: errs.ErrInvalidEmail},
		{name: "empty input", input: "",
			want: errs.ErrInvalidEmail},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := isEmail(tc.input)
			if res != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", res, tc.want)
			}
		})
	}
}

func TestIsCartNumber(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "good card", input: "4444 4444 4444 4444",
			want: nil},
		{name: "bad card", input: "123 31 dldsl3 11",
			want: errs.ErrInvalidCardNumber},
		{name: "empty input", input: "",
			want: errs.ErrInvalidCardNumber},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have := isCardNumber(tc.input)
			if have != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", have, tc.want)
			}
		})
	}
}

func TestIsCartNumberOrEmpty(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "good card", input: "4444 4444 4444 4444",
			want: nil},
		{name: "bad card", input: "123 31 dldsl3 11",
			want: errs.ErrInvalidCardNumber},
		{name: "empty input", input: "",
			want: nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have := isCardNumberOrEmpty(tc.input)
			if have != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", have, tc.want)
			}
		})
	}
}

func TestIsCardHolder(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "good holder", input: "JACK WHITE",
			want: nil},
		{name: "bad holder", input: "BOB ivan",
			want: errs.ErrInvalidCardHolder},
		{name: "empty input", input: "",
			want: errs.ErrInvalidCardHolder},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have := isCardHolder(tc.input)
			if have != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", have, tc.want)
			}
		})
	}
}

func TestIsCardHolderOrEmpty(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "good holder", input: "JACK WHITE",
			want: nil},
		{name: "bad holder", input: "BOB ivan",
			want: errs.ErrInvalidCardHolder},
		{name: "empty input", input: "",
			want: nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have := isCardHolderOrEmpty(tc.input)
			if have != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", have, tc.want)
			}
		})
	}
}

func TestIsCardExpire(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "good expire", input: "11/22",
			want: nil},
		{name: "bad expire", input: "100/1",
			want: errs.ErrInvalidCardExpire},
		{name: "empty input", input: "",
			want: errs.ErrInvalidCardExpire},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have := isCardExire(tc.input)
			if have != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", have, tc.want)
			}
		})
	}
}

func TestIsCardExpireOrEmpty(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "good expire", input: "11/22",
			want: nil},
		{name: "bad expire", input: "100/1",
			want: errs.ErrInvalidCardExpire},
		{name: "empty input", input: "",
			want: nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have := isCardExpireOrEmpty(tc.input)
			if have != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", have, tc.want)
			}
		})
	}
}

func TestIsCardCVC(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "good CVC", input: "162",
			want: nil},
		{name: "bad CVC", input: "1000",
			want: errs.ErrInvalidCardCVV},
		{name: "empty input", input: "",
			want: errs.ErrInvalidCardCVV},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have := isCardCVC(tc.input)
			if have != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", have, tc.want)
			}
		})
	}
}

func TestIsCardCVCorEmpty(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "good CVC", input: "162",
			want: nil},
		{name: "bad CVC", input: "1000",
			want: errs.ErrInvalidCardCVV},
		{name: "empty input", input: "",
			want: nil},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			have := isCardCVCorEmpty(tc.input)
			if have != tc.want {
				t.Errorf("validator mismatch: have: %v,  want: %v", have, tc.want)
			}
		})
	}
}
