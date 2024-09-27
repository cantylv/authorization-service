package myerrors

import "errors"

var (
	// HTTP RESPONSES
	ErrInternal             = errors.New("internal server error, please try again later")
	ErrInvalidData          = errors.New("you has passed invalid data in request data")
	ErrNoRequestIdInContext = errors.New("no request_id in request context")
	ErrInvalidEmail         = errors.New("you have specified invalid email in url path")
	// DATABASE
	ErrNoRowsAffected   = errors.New("no rows were affected")
	ErrUserNotExist     = errors.New("user is not exist")
	ErrUserAlreadyExist = errors.New("user with this email already exist")
	ErrPasswordMismatch = errors.New("passwords are not equal")
)
