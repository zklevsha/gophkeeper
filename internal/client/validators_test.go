package client

import (
	"testing"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

func TestIsEmpty(t *testing.T) {

	tt := []struct {
		name  string
		input string
		want  error
	}{
		{name: "input empty", input: "",
			want: structs.ErrEmptyInput},
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
			want: structs.ErrInvalidEmail},
		{name: "empty input", input: "",
			want: structs.ErrInvalidEmail},
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
