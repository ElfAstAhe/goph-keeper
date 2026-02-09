package error

import (
	"fmt"
)

type AppConfigItemError struct {
	param string
	err   error
}

var ErrAppConfigItem *AppConfigItemError

func NewAppConfigItemError(param string, err error) *AppConfigItemError {
	return &AppConfigItemError{param, err}
}

func (ch *AppConfigItemError) Error() string {
	msg := "CMN: config item error"
	if ch.param != "" {
		msg = fmt.Sprintf("%s: param [%s]", msg, ch.param)
	}
	if ch.err != nil {
		msg = fmt.Sprintf("%s: %v", msg, ch.err)
	}

	return msg
}

func (ch *AppConfigItemError) Param() string {
	return ch.param
}
