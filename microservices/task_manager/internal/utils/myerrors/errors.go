package myerrors

import "errors"

var (
	ErrNoRequestIdInContext = errors.New("no request_id in request context")
	ErrNoMetaInContext      = errors.New("no meta in request context")
	ErrInternal             = errors.New("internal server error, please try again later")
)

// DTO
var (
	ErrInvalidEmail     = errors.New("incorrect email was sent, correct format is username@domain.extension, e.g.: gref@sber.ru")
	ErrInvalidStatus    = errors.New("status must be in range(approved, rejected)")
	ErrInvalidFirstName = errors.New("incorrect first name was sent, it must start with a capital letter and be between 2 and 50 characters long")
	ErrInvalidLastName  = errors.New("incorrect last name was sent, it must start with a capital letter and be between 2 and 50 characters long")
	ErrPasswordTooLong  = errors.New("password is too long, it must be between 8 and 30 characters long")
	ErrPasswordTooShort = errors.New("password is too short, it must be between 8 and 30 characters long")
	ErrPasswordFormat   = errors.New("password must contain at least one digit and one capital letter")
)
