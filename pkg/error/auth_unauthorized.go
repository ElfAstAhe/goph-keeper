package error

import "fmt"

type AuthUnauthorizedError struct {
	message string
	err     error
}

var ErrAuthUnauthorized *AuthUnauthorizedError

func NewAuthUnauthorizedError(message string, err error) AuthUnauthorizedError {
	return AuthUnauthorizedError{message: message, err: err}
}

func (e AuthUnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: message [%s] error [%v]", e.message, e.err)
}

func (e AuthUnauthorizedError) Unwrap() error {
	return e.err
}
