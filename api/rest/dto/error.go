package dto

type ErrorDto struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"Invalid input data"`
	Details string `json:"details,omitempty" example:"field 'e-mail' is required"`
}

func NewErrorDto(code int, Message string, Details string) *ErrorDto {
	return &ErrorDto{
		Code:    code,
		Message: Message,
		Details: Details,
	}
}

func NewErrorDtoFromError(code int, err error) *ErrorDto {
	var msg string
	if err != nil {
		msg = err.Error()
	}

	return NewErrorDto(code, msg, "")
}
