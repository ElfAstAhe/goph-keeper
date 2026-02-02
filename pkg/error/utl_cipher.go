package error

import (
	"fmt"
)

type UtlCipherError struct {
	message string
	err     error
}

var ErrUtlCipher *UtlCipherError

func NewUtlCipherError(msg string, err error) *UtlCipherError {
	return &UtlCipherError{
		message: msg,
		err:     err,
	}
}

func (e *UtlCipherError) Error() string {
	return fmt.Sprintf("cipher error with message [%s], with error [%s]", e.message, e.err)
}

func (e *UtlCipherError) Unwrap() error {
	return e.err
}
