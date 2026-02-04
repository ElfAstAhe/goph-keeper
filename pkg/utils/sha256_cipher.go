package utils

import (
	"crypto/sha256"
	"hash"

	errs "github.com/ElfAstAhe/goph-keeper/pkg/error"
)

type SHA256Cipher struct {
	hash hash.Hash
}

func NewSHA256Cipher() *SHA256Cipher {
	return &SHA256Cipher{
		hash: sha256.New(),
	}
}

func (sh *SHA256Cipher) Encrypt(data []byte) ([]byte, error) {
	sh.hash.Reset()
	_, err := sh.hash.Write(data)
	if err != nil {
		return nil, errs.NewUtlCipherError("error write into buffer", err)
	}

	var sum [32]byte
	sh.hash.Sum(sum[:0])

	return sum[:], nil
}

func (sh *SHA256Cipher) EncryptString(s string) (string, error) {
	res, err := sh.Encrypt([]byte(s))

	return string(res), err
}

func (sh *SHA256Cipher) Decrypt(data []byte) ([]byte, error) {
	return nil, errs.NewAppCommonError("SHA256Cipher.Decrypt not implemented", nil)
}

func (sh *SHA256Cipher) DecryptString(s string) (string, error) {
	return "", errs.NewAppCommonError("SHA256Cipher.DecryptString not implemented", nil)
}
