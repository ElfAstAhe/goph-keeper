package err

import "fmt"

type BllModelNotExistsError struct {
	model string
	msg   string
	err   error
}

var ErrBllModelNotExists *BllModelNotExistsError

func NewBllModelNotExistsError(model string, message string, err error) *BllModelNotExistsError {
	return &BllModelNotExistsError{
		model: model,
		err:   err,
	}
}

func (mne *BllModelNotExistsError) Error() string {
	return fmt.Sprintf("BLL: model [%s] not exists with message [%s] with error [%v]", mne.model, mne.msg, mne.err)
}

func (mne *BllModelNotExistsError) Unwrap() error {
	return mne.err
}
