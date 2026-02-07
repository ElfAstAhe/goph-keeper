package dto

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
