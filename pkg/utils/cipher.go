package utils

import (
	"bytes"
	"strings"
)

const (
	CipherStringPrefix string = "cipher::"
)

var (
	CipherPrefix = []byte(CipherStringPrefix)
)

type Cipher interface {
	Encrypt([]byte) ([]byte, error)
	EncryptString(string) (string, error)
	Decrypt([]byte) ([]byte, error)
	DecryptString(string) (string, error)
}

type CipherHelper struct {
	cipher Cipher
}

func NewCipherHelper(cipher Cipher) *CipherHelper {
	return &CipherHelper{
		cipher: cipher,
	}
}

func (ch *CipherHelper) EncryptString(s string) string {
	if s == "" || ch.IsStringEncrypted(s) {
		return s
	}

	// шифруем
	res, err := ch.cipher.EncryptString(s)
	if err != nil {
		return s
	}

	// результат в base64 + prefix
	return CipherStringPrefix + res
}

func (ch *CipherHelper) DecryptString(s string) string {
	if s == "" || !ch.IsStringEncrypted(s) {
		return s
	}

	// убираем префикс и проверяем есть хоть что-нибудь
	encrypted := strings.TrimPrefix(s, CipherStringPrefix)
	if encrypted == "" {
		return s
	}

	// расшифровываем
	res, err := ch.cipher.DecryptString(encrypted)
	if err != nil {
		return s
	}

	// возвращаем результат
	return res
}

func (ch *CipherHelper) EncryptBinary(data []byte) []byte {
	if ch.IsEncrypted(data) {
		return data
	}

	prefix := []byte(CipherPrefix)
	encrypted, err := ch.cipher.Encrypt(data)
	if err != nil {
		return data
	}

	res := make([]byte, len(prefix)+len(encrypted))

	// Копируем части
	copy(res, prefix)
	copy(res[len(prefix):], encrypted)

	return res
}

func (ch *CipherHelper) DecryptBinary(data []byte) []byte {
	if !ch.IsEncrypted(data) {
		return data
	}

	prefix := []byte(CipherPrefix)
	prefixLen := len(prefix)

	// расшифровываем
	res, err := ch.cipher.Decrypt(data[prefixLen:])
	if err != nil {
		return data
	}

	return res
}

func (ch *CipherHelper) IsStringEncrypted(s string) bool {
	return strings.HasPrefix(s, CipherStringPrefix)
}

func (ch *CipherHelper) IsEncrypted(data []byte) bool {
	prefix := []byte(CipherPrefix)
	prefixLen := len(prefix)

	// Проверяем, что данных достаточно, чтобы в них физически мог быть префикс
	if len(data) < prefixLen {
		return false
	}

	// Сравниваем только начальную часть данных с префиксом
	return bytes.Equal(data[:prefixLen], prefix)
}
