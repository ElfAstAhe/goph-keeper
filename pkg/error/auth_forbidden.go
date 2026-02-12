package error

import (
	"fmt"
)

type AuthForbiddenError struct {
	message string
	err     error
}

var ErrAuthForbidden *AuthForbiddenError

func NewAuthForbiddenError(message string, err error) *AuthForbiddenError {
	return &AuthForbiddenError{
		message: message,
		err:     err,
	}
}

func (e *AuthForbiddenError) Error() string {
	msg := "AUTH: forbidden"
	if e.message != "" {
		msg = fmt.Sprintf("%s with message [%s]", msg, e.message)
	}
	if e.err != nil {
		msg = fmt.Sprintf("%s: [%v]", msg, e.err)
	}

	return msg
}

func (e *AuthForbiddenError) Unwrap() error {
	return e.err
}
