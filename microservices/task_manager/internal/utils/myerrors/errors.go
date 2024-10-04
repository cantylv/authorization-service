package myerrors

import "errors"

var (
	ErrNoRequestIdInContext = errors.New("no request_id in request context")
)
