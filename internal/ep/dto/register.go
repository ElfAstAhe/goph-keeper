package dto

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
