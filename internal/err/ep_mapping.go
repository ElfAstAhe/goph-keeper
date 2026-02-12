package err

import "fmt"

type EpMappingError struct {
	Op     string
	Source string
	Dest   string
	Msg    string
	Err    error
}

var ErrEpMapping *EpMappingError

func NewEpMappingError(op, source, dest, msg string, err error) *EpMappingError {
	return &EpMappingError{
		Op:     op,
		Source: source,
		Dest:   dest,
		Msg:    msg,
		Err:    err,
	}
}

func (e *EpMappingError) Error() string {
	msg := fmt.Sprintf("EP: %s map error from [%s] to [%s]", e.Op, e.Source, e.Dest)
	if e.Msg != "" {
		msg = fmt.Sprintf("%s with message [%s]", msg, e.Msg)
	}
	if e.Err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Err)
	}

	return msg
}

func (e *EpMappingError) Unwrap() error {
	return e.Err
}
