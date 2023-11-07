package xsecret

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
)

type ScType int

const (
	SCTYPE_DES ScType = iota
	SCTYPE_3DES
	SCTYPE_AES
)

// 末尾填充字节
func PKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize // 要填充的值和个数
	slice1 := []byte{byte(padding)}            // 要填充的单个二进制值
	slice2 := bytes.Repeat(slice1, padding)    // 要填充的二进制数组
	return append(data, slice2...)             // 填充到数据末端
}

// 去除填充的字节
func PKCS5UnPadding(data []byte) []byte {
	unpadding := data[len(data)-1]                // 获取二进制数组最后一个数值
	result := data[:(len(data) - int(unpadding))] // 截取开始至总长度减去填充值之间的有效数据
	return result
}

// // 末尾填充字节
// func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
// 	padding := blockSize - len(ciphertext)%blockSize
// 	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
// 	return append(ciphertext, padtext...)
// }

// // 去除填充的字节
// func PKCS7UnPadding(origData []byte) []byte {
// 	length := len(origData)
// 	unpadding := int(origData[length-1])
// 	return origData[:(length - unpadding)]
// }

func GetBlock(scType ScType, key []byte) (cipher.Block, error) {
	var err error
	var block cipher.Block
	switch scType {
	case SCTYPE_DES:
		block, err = des.NewCipher(key)
	case SCTYPE_3DES:
		block, err = des.NewTripleDESCipher(key)
	case SCTYPE_AES:
		block, err = aes.NewCipher(key)
	}
	if err != nil {
		return nil, err
	}
	return block, nil
}

func Encrypt(scType ScType, data, key string) (string, error) {
	bData := []byte(data)
	bKey := []byte(key)

	block, err := GetBlock(scType, bKey)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()

	// 2、对明文进行填充（参数为原始字节切片和密码对象的区块个数）
	paddingBytes := PKCS5Padding(bData, blockSize)

	// 3、实例化加密模式（参数为密码对象和密钥）
	blockMode := cipher.NewCBCEncrypter(block, bKey[:blockSize])

	// 4、对填充字节后的明文进行加密（参数为加密字节切片和填充字节切片）
	cipherBytes := make([]byte, len(paddingBytes))
	blockMode.CryptBlocks(cipherBytes, paddingBytes)
	return base64.RawURLEncoding.EncodeToString(cipherBytes), nil
}

func Decrypt(scType ScType, data string, key string) (string, error) {
	bData, _ := base64.StdEncoding.DecodeString(data)
	bKey := []byte(key)

	block, err := GetBlock(scType, bKey)
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()

	// 2、实例化解密模式（参数为密码对象和密钥）
	blockMode := cipher.NewCBCDecrypter(block, bKey[:blockSize])

	// 3、对密文进行解密（参数为填充字节切片和加密字节切片）
	paddingBytes := make([]byte, len(bData))
	blockMode.CryptBlocks(paddingBytes, bData)

	// 4、去除填充的字节（参数为填充切片）
	originalBytes := PKCS5UnPadding(paddingBytes)
	return string(originalBytes), nil
}
