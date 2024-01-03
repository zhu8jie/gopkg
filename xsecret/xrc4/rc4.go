package xrc4

import (
	"crypto/rc4"
)

func Encrypt(in, key []byte) []byte {
	return do(in, key)
}

func Decrypt(in, key []byte) []byte {
	return do(in, key)
}

func do(in, key []byte) []byte {
	cipher2, _ := rc4.NewCipher(key)
	cipher2.XORKeyStream(in, in) // 解密后的数据直接覆盖到str中
	return in
}
