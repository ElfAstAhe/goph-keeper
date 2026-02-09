package dto

import (
	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewLoginDto(username, password string) *LoginDto {
	return &LoginDto{
		Username: username,
		Password: password,
	}
}

func (ld *LoginDto) Validate() error {
	if ld.Username == "" {
		return errs.NewAppInvalidArgumentError("Username", "cannot be empty")
	}
	if ld.Password == "" {
		return errs.NewAppInvalidArgumentError("Password", "cannot be empty")
	}

	return nil
}
