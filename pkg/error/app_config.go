package error

import "fmt"

type AppConfigError struct {
	message string
	err     error
}

var ErrAppConfig *AppConfigError

func NewAppConfigError(message string, err error) *AppConfigError {
	return &AppConfigError{message, err}
}

func (e *AppConfigError) Error() string {
	return fmt.Sprintf("config error with message [%s] with error [%v]", e.message, e.err)
}

func (e *AppConfigError) Unwrap() error {
	return e.err
}
