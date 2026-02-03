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
	return fmt.Sprintf("invalid config item [%s] with error [%v]", ch.param, ch.err)
}

func (ch *AppConfigItemError) Param() string {
	return ch.param
}
