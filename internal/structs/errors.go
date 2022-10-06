package structs

import "errors"

// ErrUserAlreadyExists returns when
// a user with the same name already exists
var ErrUserAlreadyExists = errors.New("user already exists")

// ErrUserAuth returns when there are some authentication error
var ErrUserAuth = errors.New("authentication failed")
