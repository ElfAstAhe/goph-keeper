package err

import "fmt"

type BllModelSoftDeletedError struct {
	model string
}

var ErrBllModelSoftDeleted *BllModelSoftDeletedError

func NewBllModelSoftDeletedError(model string) *BllModelSoftDeletedError {
	return &BllModelSoftDeletedError{model: model}
}

func (e *BllModelSoftDeletedError) Error() string {
	return fmt.Sprintf("BLL: model [%s] removed", e.model)
}
