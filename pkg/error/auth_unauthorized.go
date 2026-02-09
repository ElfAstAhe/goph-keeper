package error

import (
	"fmt"
)

type AuthUnauthorizedError struct {
	message string
	err     error
}

var ErrAuthUnauthorized *AuthUnauthorizedError

func NewAuthUnauthorizedError(message string, err error) *AuthUnauthorizedError {
	return &AuthUnauthorizedError{
		message: message,
		err:     err,
	}
}

func (e *AuthUnauthorizedError) Error() string {
	msg := "AUTH: unauthorized"
	if e.message != "" {
		msg = fmt.Sprintf("%s with message %s", msg, e.message)
	}
	if e.err != nil {
		msg = fmt.Sprintf("%s: %s", msg, e.err)
	}

	return msg
}

func (e *AuthUnauthorizedError) Unwrap() error {
	return e.err
}
