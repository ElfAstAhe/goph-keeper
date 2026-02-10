package dto

type ChangeKeysResultDto struct {
	PublicKey string `json:"public_key"`
}

func NewChangeKeysResultDto(publicKey string) *ChangeKeysResultDto {
	return &ChangeKeysResultDto{
		PublicKey: publicKey,
	}
}
