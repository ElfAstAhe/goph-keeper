package dto

import "time"

type UserDataDto struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name"`
	DataKind   string    `json:"data_kind"`
	TextData   string    `json:"text_data,omitempty"`
	BinaryData []byte    `json:"binary_data,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	ModifiedAt time.Time `json:"modified_at,omitempty"`
}

func NewUserDataDto(
	id string,
	name string,
	dataKind string,
	textData string,
	binaryData []byte,
	createdAt time.Time,
	modifiedAt time.Time,
) *UserDataDto {
	return &UserDataDto{
		ID:         id,
		Name:       name,
		DataKind:   dataKind,
		TextData:   textData,
		BinaryData: binaryData,
		CreatedAt:  createdAt,
		ModifiedAt: modifiedAt,
	}
}

func NewEmptyUserDataDto() *UserDataDto {
	return &UserDataDto{
		BinaryData: make([]byte, 0),
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}
}
