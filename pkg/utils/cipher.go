package utils

import "strings"

const (
	CipherPrefix string = "cipher:"
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

func (ch *CipherHelper) EncryptString(s string) (string, error) {
	if s == "" {
		return s, nil
	}

	return ch.cipher.EncryptString(s)
}

func (ch *CipherHelper) MustEncryptString(s string) string {
	res, err := ch.EncryptString(s)
	if err != nil {
		panic(err)
	}

	return CipherPrefix + res
}

func (ch *CipherHelper) DecryptString(s string) (string, error) {
	if s == "" || !ch.isEncrypted(s) {
		return s, nil
	}

	return ch.cipher.DecryptString(s)
}

func (ch *CipherHelper) MustDecryptString(s string) string {
	res, err := ch.DecryptString(s)
	if err != nil {
		panic(err)
	}

	return res
}

func (ch *CipherHelper) isEncrypted(s string) bool {
	return strings.HasPrefix(s, CipherPrefix)
}
