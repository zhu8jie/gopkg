package xaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func zeroPad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{0}, padding)
	return append(data, padText...)
}

func unZeroPad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length%blockSize != 0 {
		return nil, fmt.Errorf("data is not block-size aligned")
	}

	for i := length - 1; i >= 0; i-- {
		if data[i] != 0 {
			return data[:i+1], nil
		}
	}

	return nil, fmt.Errorf("data not padded")
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// =================== CBC ======================
func AesEncryptCBC(origData []byte, key, iv []byte) (encrypted []byte) {

	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()               // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize) // 补全码

	if len(iv) == 0 {
		iv = key[:blockSize]
	}

	blockMode := cipher.NewCBCEncrypter(block, iv) // 加密模式
	encrypted = make([]byte, len(origData))        // 创建数组
	blockMode.CryptBlocks(encrypted, origData)     // 加密
	return encrypted
}
func AesDecryptCBC(encrypted []byte, key, iv []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key) // 分组秘钥
	blockSize := block.BlockSize() // 获取秘钥块的长度

	if len(iv) == 0 {
		iv = key[:blockSize]
	}

	blockMode := cipher.NewCBCDecrypter(block, iv) // 加密模式
	decrypted = make([]byte, len(encrypted))       // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)    // 解密
	decrypted = pkcs5UnPadding(decrypted)          // 去除补全码
	return decrypted
}

// =================== ECB ======================
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {

	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))

	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

func AESEncryptZeroPadEcb(data, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// plaintext := []byte(data)
	// Padding
	data = zeroPad(data, aes.BlockSize)
	ciphertext := make([]byte, len(data))

	for start := 0; start < len(data); start += aes.BlockSize {
		block.Encrypt(ciphertext[start:start+aes.BlockSize], data[start:start+aes.BlockSize])
	}

	return ciphertext, nil
}

func AESDecryptZeroPadEcb(data, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// ciphertext, err := base64.StdEncoding.DecodeString(data)
	// if err != nil {
	// 	return "", err
	// }

	plaintext := make([]byte, len(data))

	for start := 0; start < len(data); start += aes.BlockSize {
		block.Decrypt(plaintext[start:start+aes.BlockSize], data[start:start+aes.BlockSize])
	}

	// Unpadding
	plaintext, err = unZeroPad(plaintext, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// =================== CFB ======================
func AesEncryptCFB(origData []byte, key []byte) (encrypted []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}
func AesDecryptCFB(encrypted []byte, key []byte) (decrypted []byte) {
	block, _ := aes.NewCipher(key)
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}

// 秘钥长度验证
func validKey(key []byte) bool {
	k := len(key)
	switch k {
	default:
		return false
	case 16, 24, 32:
		return true
	}
}
