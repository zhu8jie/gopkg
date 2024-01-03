package xaescbc

// import (
// 	"bytes"
// 	"crypto/aes"
// 	"crypto/cipher"
// 	"crypto/rand"
// 	"encoding/base64"
// 	"errors"
// 	"io"
// )

// func AesEncrypt(rawData, key string) (string, error) {
// 	data, err := AesCBCEncrypt([]byte(rawData), []byte(key))
// 	return base64.RawURLEncoding.EncodeToString([]byte(data)), err
// }

// //aes 加密，填充秘钥 key 的 16 位，24,32 分别对应 AES-128, AES-192, or AES-256.
// func AesCBCEncrypt(rawData, key []byte) (string, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}
// 	blockSize := block.BlockSize()
// 	rawData = PKCS7Padding(rawData, blockSize)
// 	cipherText := make([]byte, blockSize+len(rawData))
// 	iv := cipherText[:blockSize]
// 	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
// 		return "", err
// 	}
// 	mode := cipher.NewCBCEncrypter(block, iv)
// 	mode.CryptBlocks(cipherText[blockSize:], rawData)
// 	return string(cipherText), nil
// }

// func AesDecrypt(rawData string, key string) (string, error) {
// 	b, err := base64.RawURLEncoding.DecodeString(rawData)
// 	if err != nil {
// 		return "", err
// 	}
// 	return AesCBCDncrypt(b, []byte(key))
// }

// func AesCBCDncrypt(encryptData, key []byte) (string, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}
// 	blockSize := block.BlockSize()
// 	if len(encryptData) < blockSize {
// 		return "", errors.New("encrypt data length < block size")
// 	}
// 	iv := encryptData[:blockSize]
// 	encryptData = encryptData[blockSize:]
// 	// CBC mode always works in whole blocks.
// 	if len(encryptData)%blockSize != 0 {
// 		return "", errors.New("CBC mode always works in whole blocks.")
// 	}
// 	mode := cipher.NewCBCDecrypter(block, iv)
// 	// CryptBlocks can work in-place if the two arguments are the same.
// 	mode.CryptBlocks(encryptData, encryptData)
// 	encryptData = PKCS7UnPadding(encryptData)
// 	return string(encryptData), nil
// }

// func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
// 	padding := blockSize - len(ciphertext)%blockSize
// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(ciphertext, padtext...)
// }

// func PKCS7UnPadding(origData []byte) []byte {
// 	length := len(origData)
// 	unpadding := int(origData[length-1])
// 	return origData[:(length - unpadding)]
// }
