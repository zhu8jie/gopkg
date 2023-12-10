package xbase64

import "encoding/base64"

func Base64Encrypt(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func Base64Decrypt(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
