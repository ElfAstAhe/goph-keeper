package utils

import (
	"encoding/base64"
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
	if s == "" {
		return s
	}

	// шифруем
	res, err := ch.cipher.Encrypt([]byte(s))
	if err != nil {
		return s
	}

	// результат в base64 + prefix
	return CipherPrefix + base64.StdEncoding.EncodeToString(res)
}

func (ch *CipherHelper) DecryptString(s string) string {
	if s == "" || !ch.isEncrypted(s) {
		return s
	}

	// убираем префикс и проверяем есть хоть что-нибудь
	encrypted := strings.TrimPrefix(CipherPrefix, s)
	if encrypted == "" {
		return s
	}

	// из base64 в набор байт
	bytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return s
	}

	// расшифровываем
	res, err := ch.cipher.Decrypt(bytes)
	if err != nil {
		return s
	}

	// возвращаем результат
	return string(res)
}

func (ch *CipherHelper) isEncrypted(s string) bool {
	return strings.HasPrefix(s, CipherPrefix)
}
