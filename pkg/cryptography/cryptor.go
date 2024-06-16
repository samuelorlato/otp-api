package cryptography

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"

	"github.com/samuelorlato/otp-api.git/pkg/errors"
)

type Cryptor struct {
	key string
	iv  string
}

func NewCryptor(key string, iv string) *Cryptor {
	return &Cryptor{
		key: key,
		iv:  iv,
	}
}

func (c *Cryptor) Encrypt(text string) (*errors.HTTPError, *string) {
	var plainTextBlock []byte

	length := len(text)
	if length%16 != 0 {
		extendBlock := 16 - (length % 16)
		plainTextBlock = make([]byte, length+extendBlock)
		copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	} else {
		plainTextBlock = make([]byte, length)
	}

	copy(plainTextBlock, text)
	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		err := errors.NewGenericError(err)
		return err, nil
	}

	cipherText := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, []byte(c.iv))
	mode.CryptBlocks(cipherText, plainTextBlock)

	encryptedText := base64.StdEncoding.EncodeToString(cipherText)

	return nil, &encryptedText
}

func (c *Cryptor) PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

func (c *Cryptor) Decrypt(text string) (*errors.HTTPError, *string) {
	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		err := errors.NewGenericError(err)
		return err, nil
	}

	block, err := aes.NewCipher([]byte(c.key))
	if err != nil {
		err := errors.NewGenericError(err)
		return err, nil
	}

	if len(cipherText)%aes.BlockSize != 0 {
		err := errors.NewGenericError(fmt.Errorf("Block size can not be zero"))
		return err, nil
	}

	mode := cipher.NewCBCDecrypter(block, []byte(c.iv))
	mode.CryptBlocks(cipherText, cipherText)
	decryptedText := string(c.PKCS5UnPadding(cipherText))

	return nil, &decryptedText
}
