package error

import "fmt"

type UtlDBError struct {
	message string
	err     error
}

var ErrUtlDB *UtlDBError

func NewUtlDBError(message string, err error) *UtlDBError {
	return &UtlDBError{message: message, err: err}
}

func (err *UtlDBError) Error() string {
	return fmt.Sprintf("error db util with message [%s] with error [%v]", err.message, err.err)
}

func (err *UtlDBError) Unwrap() error {
	return err.err
}
