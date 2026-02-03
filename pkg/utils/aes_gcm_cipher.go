package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"

	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

// AesGcmCipher -AES/GCM cipher util
type AesGcmCipher struct {
	block cipher.Block
	gcm   cipher.AEAD
}

// NewAesGcmCipher - aes/gcm util constructor
func NewAesGcmCipher(key []byte) (*AesGcmCipher, error) {
	instance := new(AesGcmCipher)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errs.NewUtlCipherError("error create cipher", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, errs.NewUtlCipherError("error create gcm", err)
	}

	instance.block = block
	instance.gcm = gcm

	return instance, nil
}

// MustNewAesGcmCipher - aes/gcm util constructor, returns instance, but generate panic on error
func MustNewAesGcmCipher(key []byte) *AesGcmCipher {
	instance, err := NewAesGcmCipher(key)
	if err != nil {
		panic(err)
	}

	return instance
}

func (a *AesGcmCipher) Encrypt(data []byte) ([]byte, error) {
	nonce := make([]byte, a.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, errs.NewUtlCipherError("error fill nonce", err)
	}

	return a.gcm.Seal(nonce, nonce, data, nil), nil
}

func (a *AesGcmCipher) EncryptString(s string) (string, error) {
	res, err := a.Encrypt([]byte(s))

	return string(res), err
}

func (a *AesGcmCipher) Decrypt(data []byte) ([]byte, error) {
	nonceSize := a.gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errs.NewUtlCipherError("error data validation", errs.NewAppInvalidArgumentError("data", data))
	}

	nonce, data := data[:nonceSize], data[nonceSize:]
	plain, err := a.gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, errs.NewUtlCipherError("error decrypt data", err)
	}

	return plain, nil
}

func (a *AesGcmCipher) DecryptString(s string) (string, error) {
	res, err := a.Decrypt([]byte(s))

	return string(res), err
}
