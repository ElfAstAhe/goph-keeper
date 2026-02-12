package dto

type UpdatePasswordDto struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func NewUpdatePasswordDto(oldPassword, newPassword string) *UpdatePasswordDto {
	return &UpdatePasswordDto{
		OldPassword: oldPassword,
		NewPassword: newPassword,
	}
}

func NewEmptyUpdatePasswordDto() *UpdatePasswordDto {
	return NewUpdatePasswordDto("", "")
}
