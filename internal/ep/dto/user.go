package dto

type UserDto struct {
	ID           string `json:"id,omitempty"`
	Username     string `json:"username,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
	PublicKey    string `json:"public_key,omitempty"`
	Active       bool   `json:"active,omitempty"`
	Person       string `json:"person,omitempty"`
	EMail        string `json:"email,omitempty"`
}

func NewUserDto(
	id string,
	username string,
	passwordHash string,
	publicKey string,
	active bool,
	person string,
	email string,
) *UserDto {
	return &UserDto{
		ID:           id,
		Username:     username,
		PasswordHash: passwordHash,
		PublicKey:    publicKey,
		Active:       active,
		Person:       person,
		EMail:        email,
	}
}
