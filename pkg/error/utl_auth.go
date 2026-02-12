package error

import "fmt"

type UtlAuthError struct {
	message string
	err     error
}

var ErrUtlAuth *UtlAuthError

func NewUtlAuthError(message string, err error) *UtlAuthError {
	return &UtlAuthError{
		message: message,
		err:     err,
	}
}

func (e *UtlAuthError) Error() string {
	return fmt.Sprintf("auth util error with message [%s] with error [%v]", e.message, e.err)
}

func (e *UtlAuthError) Unwrap() error {
	return e.err
}
