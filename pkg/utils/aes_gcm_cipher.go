package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

type AesGcmCipher struct {
	key []byte
}

func NewAesGcmCipher(key []byte) *AesGcmCipher {
	return &AesGcmCipher{
		key: key,
	}
}

func (a *AesGcmCipher) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, errs.NewUtlCipherError("error create cipher", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errs.NewUtlCipherError("error create gcm", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errs.NewUtlCipherError("error fill nonce", err)
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (a *AesGcmCipher) EncryptString(s string) (string, error) {
	res, err := a.Encrypt([]byte(s))

	return string(res), err
}

func (a *AesGcmCipher) Decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, errs.NewUtlCipherError("error create cipher", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errs.NewUtlCipherError("error create gcm", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errs.NewUtlCipherError("error data validation", errs.NewAppInvalidArgumentError("data", data))
	}

	nonce, data := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, errs.NewUtlCipherError("error decrypt data", err)
	}

	return plain, nil
}

func (a *AesGcmCipher) DecryptString(s string) (string, error) {
	res, err := a.Decrypt([]byte(s))

	return string(res), err
}
