package dto

type LoginResultDto struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func NewLoginResultDto(token string, refreshToken string) *LoginResultDto {
	return &LoginResultDto{
		Token:        token,
		RefreshToken: refreshToken,
	}
}
