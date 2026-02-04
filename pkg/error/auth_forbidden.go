package error

import "fmt"

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
	return fmt.Sprintf("forbidden [%s] with error [%v]", e.message, e.err)
}

func (e *AuthForbiddenError) Unwrap() error {
	return e.err
}
