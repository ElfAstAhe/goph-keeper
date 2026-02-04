package utils

import (
	"strings"
)

const (
	CipherPrefix string = "cipher::"
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
	if s == "" || ch.IsEncrypted(s) {
		return s
	}

	// шифруем
	res, err := ch.cipher.EncryptString(s)
	if err != nil {
		return s
	}

	// результат в base64 + prefix
	return CipherPrefix + res
}

func (ch *CipherHelper) DecryptString(s string) string {
	if s == "" || !ch.IsEncrypted(s) {
		return s
	}

	// убираем префикс и проверяем есть хоть что-нибудь
	encrypted := strings.TrimPrefix(s, CipherPrefix)
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

func (ch *CipherHelper) IsEncrypted(s string) bool {
	return strings.HasPrefix(s, CipherPrefix)
}
