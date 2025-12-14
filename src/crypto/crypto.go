package crypto

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

const (
	NoEncrypt = 0
	AesCbc    = 1
)

// AES CBC加密方式时，key长度必须为16｜24｜32字节

func Encrypt(srcStr []byte, crypto int, key []byte) []byte {
	switch crypto {
	case AesCbc:
		return cbcEncrypt(srcStr, key, iV)
	default:
		return srcStr
	}
}

func Decrypt(encryptStr []byte, crypto int, key []byte) []byte {
	switch crypto {
	case AesCbc:
		return cbcDecrypt(encryptStr, key, iV)
	default:
		return encryptStr
	}
}

func ReadObj(filePath string, obj any, encrypt int, key []byte) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "read file")
	}
	data = Decrypt(data, encrypt, key)
	err = json.Unmarshal(data, obj)
	if err != nil {
		return errors.Wrap(err, "json unmarshal")
	}
	return nil
}

func WriteObj(filePath string, obj any, encrypt int, key []byte) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return errors.Wrap(err, "json marshal")
	}
	data = Encrypt(data, encrypt, key)
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return errors.Wrap(err, "write file")
	}
	return nil
}

type Crypto struct {
	crypto int
	key    []byte
}

func NewCrypto(crypto int, key []byte) *Crypto {
	return &Crypto{
		crypto: crypto,
		key:    key,
	}
}

func (c *Crypto) Encrypt(srcStr []byte) []byte {
	return Encrypt(srcStr, c.crypto, c.key)
}

func (c *Crypto) Decrypt(encryptStr []byte) []byte {
	return Decrypt(encryptStr, c.crypto, c.key)
}

func (c *Crypto) ReadObj(filePath string, obj any) error {
	return ReadObj(filePath, obj, c.crypto, c.key)
}

func (c *Crypto) WriteObj(filePath string, obj any) error {
	return WriteObj(filePath, obj, c.crypto, c.key)
}
