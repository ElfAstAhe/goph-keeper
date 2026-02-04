package model

// DataKinds - виды хранимой информации
const (
	// DataKindCredential - пара login/password
	DataKindCredential string = "credential"

	// DataKindPlainText - произвольные текстовые данные
	DataKindPlainText string = "plaintext"

	// DataKindBinary - произвольные бинарные данные
	DataKindBinary string = "binary"

	// DataKindBankCard - банковская карта
	DataKindBankCard string = "bankcard"
)

var DataKinds = map[string]struct{}{
	DataKindCredential: {},
	DataKindPlainText:  {},
	DataKindBinary:     {},
	DataKindBankCard:   {},
}

func IsDataKind(k string) bool {
	_, ok := DataKinds[k]

	return ok
}

type Identity interface {
	TableName() string

	GetID() string

	IsExists() bool

	ValidateCreate() error
	ValidateChange() error

	BeforeCreate() error
	BeforeChange() error
}

type IdentitySoftRemovable interface {
	IsDeleted() bool
}
