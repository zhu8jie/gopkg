package xrc4

import (
	"crypto/rc4"
	"errors"
)

func Encrypt(in, key []byte) ([]byte, error) {
	return do(in, key)
}

func Decrypt(in, key []byte) ([]byte, error) {
	return do(in, key)
}

func do(in, key []byte) ([]byte, error) {
	if len(in) == 0 || len(key) == 0 {
		return nil, errors.New("input or key must be not null.")
	}
	cipher2, err := rc4.NewCipher(key)
	cipher2.XORKeyStream(in, in) // 解密后的数据直接覆盖到str中
	return in, err
}
