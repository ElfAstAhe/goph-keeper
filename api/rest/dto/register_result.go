package dto

type RegisterResultDto struct {
	PublicKey string `json:"public_key" binding:"required"`
}

func NewRegisterResultDto(publicKey string) *RegisterResultDto {
	return &RegisterResultDto{
		PublicKey: publicKey,
	}
}
