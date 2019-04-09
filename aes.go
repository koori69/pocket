// Package pocket @author K·J Create at 2019-04-09 15:09
package pocket

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

// Encrypt aes Encrypt
func Encrypt(plantText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	plantText = PKCS7Padding(plantText, block.BlockSize())

	blockModel := cipher.NewCBCEncrypter(block, key[:block.BlockSize()])

	cipherText := make([]byte, len(plantText))

	blockModel.CryptBlocks(cipherText, plantText)
	return cipherText, nil
}

// PKCS7Padding aes pkcs
func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

// Decrypt aes Decrypt
func Decrypt(cipherText, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key) //选择加密算法
	if err != nil {
		return nil, err
	}
	blockModel := cipher.NewCBCDecrypter(block, key[:block.BlockSize()])
	plantText := make([]byte, len(cipherText))
	blockModel.CryptBlocks(plantText, cipherText)
	plantText = PKCS7UnPadding(plantText)
	return plantText, nil
}

// PKCS7UnPadding aes pkcs
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	if length < 1 {
		return nil
	}
	unPadding := int(plantText[length-1])
	if length-unPadding > len(plantText) || length-unPadding < 0 {
		return nil
	}
	return plantText[:(length - unPadding)]
}
