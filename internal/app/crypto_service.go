package app

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type cryptoService struct {
	defaultKey []byte
	defaultIv  []byte
}

type CryptoService interface {
	Encrypt(input string, password string) (string, error)
	Decrypt(input string, password string) (string, error)
}

func NewCryptoService(key string, iv string) *cryptoService {
	return &cryptoService{[]byte(key), []byte(iv)}
}

func (r *cryptoService) createIvFromPassword(password string) []byte {
	result := make([]byte, 32)
	copy(result, r.defaultKey)
	copy(result, []byte(password))
	return result
}

func (r cryptoService) Encrypt(input string, password string) (result string, err error) {
	var plainTextBlock []byte
	length := len(input)

	key := r.createIvFromPassword(password)
	var block cipher.Block
	if block, err = aes.NewCipher([]byte(key)); err != nil {
		return "", err
	}

	extendBlock := 16
	if length%16 != 0 {
		extendBlock = 16 - (length % 16)
	}

	plainTextBlock = make([]byte, length+extendBlock)
	copy(plainTextBlock[length:], bytes.Repeat([]byte{uint8(extendBlock)}, extendBlock))
	copy(plainTextBlock, input)

	ciphertext := make([]byte, len(plainTextBlock))
	mode := cipher.NewCBCEncrypter(block, r.defaultIv)
	mode.CryptBlocks(ciphertext, plainTextBlock)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return str, nil
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}

func (r cryptoService) Decrypt(encrypted string, password string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	key := r.createIvFromPassword(password)
	var block cipher.Block
	if block, err = aes.NewCipher([]byte(key)); err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, r.defaultIv)
	mode.CryptBlocks(ciphertext, ciphertext)

	cipherLength := len(ciphertext)
	if cipherLength > 0 {
		unpadding := int(ciphertext[cipherLength-1])
		if cipherLength >= unpadding {
			ciphertext = ciphertext[:(cipherLength - unpadding)]
		}
	}

	return string(ciphertext), nil
}
