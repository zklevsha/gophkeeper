// errs contains all custo, errors used in this project
package errs

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

// ErrInvalidCardNumber returns when we failed to parse input card number
var ErrInvalidCardNumber = errors.New("expected format: XXXX XXXX XXXX XXXX")

// ErrInvalidCardHolder returns when we failed to parse input card holder
var ErrInvalidCardHolder = errors.New("expected format: JOHN DOE")

// ErrInvalidCardExpire return when we failed to parse input card expiration date
var ErrInvalidCardExpire = errors.New("expected format: MM/YY")

// ErrInvalidCardCVV return when we failed to parse input card CVV/CVC
var ErrInvalidCardCVV = errors.New("expected format: XXX")

// KeyValueError returns when we cant parse key-value pairs
var ErrKeyValue = errors.New("failed to parse input as json")

// ErrPdataNotFound returns when we cant find Pdata for user
var ErrPdataNotFound = errors.New("pdata not found")

// ErrPdataAlreatyEsists reutrns when Pdata already exists
var ErrPdataAlreatyEsists = errors.New("pdata already exists")

// ErrNoToken returns when JWT is not found in the request
var ErrNoToken = errors.New("no JWT")

// ErrInvalidToken returns when JWT token is invalid
var ErrInvalidToken = errors.New("JWT is invalid")
