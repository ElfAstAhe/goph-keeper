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
	msg := "CMN: config error"
	if e.message != "" {
		msg = fmt.Sprintf("%s with message [%s]", msg, e.message)
	}
	if e.err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.err)
	}

	return msg
}

func (e *AppConfigError) Unwrap() error {
	return e.err
}
