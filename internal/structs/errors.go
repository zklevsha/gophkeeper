package structs

import "errors"

// ErrUserAlreadyExists returns when
// a user with the same name already exists
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrUserAuth returns when there are some authentication error
var ErrUserAuth = errors.New("authentication failed")

// ErrEmptyInput returns when client input is empty
var ErrEmptyInput = errors.New("input is empty")

// ErrInvalidEmail returns when we failed to parse input email
var ErrInvalidEmail = errors.New("email is invalid")

// KeyValueError returns when we cant parse key-value pairs
var ErrKeyValue = errors.New("failed to parse input as json")
