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
