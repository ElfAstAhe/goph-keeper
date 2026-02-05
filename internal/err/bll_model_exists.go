package err

import "fmt"

type BllModelExistsError struct {
	id    string
	model string
}

var ErrBllModelExists *BllModelExistsError

func NewBllModelExistsError(model string, id string) *BllModelExistsError {
	return &BllModelExistsError{model: model}
}

func (bme *BllModelExistsError) Error() string {
	if bme.id == "" {
		return fmt.Sprintf("BLL: model [%s] already exists", bme.model)
	}

	return fmt.Sprintf("BLL: model [%s] with id [%s] already exists", bme.model, bme.id)
}
