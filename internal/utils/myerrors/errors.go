package myerrors

import "errors"

var (
	// HTTP RESPONSES
	ErrInternal          = errors.New("internal server error, please try again later")
	ErrAlreadyAuthorized = errors.New("you're already authorized")
	ErrNoRefreshToken    = errors.New("you can't sign out without refresh_token")
	ErrInvalidData       = errors.New("you has passed invalid data in request data")
	// NETWORK
	ErrInvalidUserAgent     = errors.New("invalid request, you must specify User-Agent HTTP Header")
	ErrInvalidRemoteIp      = errors.New("invalid request, we can't receive your ip-address, try to connect to another wi-fi module")
	ErrInvalidJwtToken      = errors.New("invalid jwt-token")
	ErrJwtAlreadyExpired    = errors.New("jwt-token is expired")
	ErrNoRequestIdInContext = errors.New("no request_id in request context")
	// DATABASE
	ErrNoRowsAffected   = errors.New("no rows were affected")
	ErrUserNotExist     = errors.New("user is not exist")
	ErrUserAlreadyExist = errors.New("user with this email already exist")
	ErrPasswordMismatch = errors.New("passwords are not equal")
)
