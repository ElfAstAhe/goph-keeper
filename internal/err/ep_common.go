package err

import "fmt"

type EpCommonError struct {
	op      string
	message string
	err     error
}

var ErrEpCommon *EpCommonError

func NewEpCommonError(op, message string, err error) *EpCommonError {
	return &EpCommonError{
		op:      op,
		message: message,
		err:     err,
	}
}

func (e *EpCommonError) Error() string {
	msg := fmt.Sprintf("EP: [%s] common error", e.op)
	if e.message != "" {
		msg = fmt.Sprintf("%s with message [%s]", msg, e.message)
	}
	if e.err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.err)
	}

	return msg
}

func (e *EpCommonError) Unwrap() error {
	return e.err
}
