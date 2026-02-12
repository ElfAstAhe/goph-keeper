package rest

import (
	"fmt"
)

type ClientError struct {
	Op         string
	StatusCode int
	Message    string
	Err        error
}

var ErrClient *ClientError

func NewClientError(op string, statusCode int, message string, err error) *ClientError {
	code := statusCode
	if err != nil {
		clientErr, ok := err.(*ClientError)
		if ok {
			code = clientErr.StatusCode
		}
	}
	return &ClientError{
		Op:         op,
		StatusCode: code,
		Message:    message,
		Err:        err,
	}
}

func (e *ClientError) Error() string {
	msg := fmt.Sprintf("CLIENT: [%s] response error", e.Op)
	if e.StatusCode > 0 {
		msg = fmt.Sprintf("%s (HTTP code %d)", msg, e.StatusCode)
	}
	if e.Message != "" {
		msg = fmt.Sprintf("%s (%s)", msg, e.Message)
	}
	if e.Err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Err)
	}

	return msg
}

func (e *ClientError) Unwrap() error {
	return e.Err
}
