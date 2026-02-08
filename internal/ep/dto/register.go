package dto

import (
	"strings"

	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

type RegisterDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Person   string `json:"person" binding:"required"`
	EMail    string `json:"e_mail" binding:"required"`
}

func NewRegisterDto(username, password, person, eMail string) *RegisterDto {
	return &RegisterDto{
		Username: username,
		Password: password,
		Person:   person,
		EMail:    eMail,
	}
}

func (r *RegisterDto) Validate() error {
	if strings.TrimSpace(r.Username) == "" {
		return errs.NewAppInvalidArgumentError("RegisterDto Username", "empty")
	}
	if strings.TrimSpace(r.Password) == "" {
		return errs.NewAppInvalidArgumentError("RegisterDto Password", "empty")
	}

	return nil
}
