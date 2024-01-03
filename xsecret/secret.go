package xsecret

import (
	"errors"

	"github.com/zhu8jie/gopkg/xsecret/xaes"
	"github.com/zhu8jie/gopkg/xsecret/xbase64"
	"github.com/zhu8jie/gopkg/xsecret/xdes"
	"github.com/zhu8jie/gopkg/xsecret/xrc4"
)

type ScType int

const (
	SCTYPE_DES ScType = iota
	SCTYPE_3DES
	SCTYPE_AES_CBC
	SCTYPE_AES_ECB
	SCTYPE_AES_CFB
	SCTYPE_RC4
)

func Encrypt(t ScType, in, key string) (string, error) {
	if in == "" {
		return "", nil
	}
	switch t {
	case SCTYPE_DES:
		b, err := xdes.EncryptDES(in, key)
		if err != nil {
			return "", err
		}
		return xbase64.Base64Encrypt(b), nil
	case SCTYPE_3DES:
		b, err := xdes.Encrypt3DES(in, key)
		if err != nil {
			return "", err
		}
		return xbase64.Base64Encrypt(b), nil
	case SCTYPE_AES_CBC:
		b := xaes.AesEncryptCBC([]byte(in), []byte(key))
		return xbase64.Base64Encrypt(b), nil
	case SCTYPE_AES_ECB:
		b := xaes.AesEncryptECB([]byte(in), []byte(key))
		return xbase64.Base64Encrypt(b), nil
	case SCTYPE_AES_CFB:
		b := xaes.AesEncryptCFB([]byte(in), []byte(key))
		return xbase64.Base64Encrypt(b), nil
	case SCTYPE_RC4:
		b := xrc4.Encrypt([]byte(in), []byte(key))
		return xbase64.Base64Encrypt(b), nil
	default:
		return "", errors.New("sctype not exist.")
	}
}

func Decrypt(t ScType, in, key string) (string, error) {
	data, err := xbase64.Base64Decrypt(in)
	if err != nil {
		return "", err
	}
	switch t {
	case SCTYPE_DES:
		return xdes.DecryptDES(data, key)
	case SCTYPE_3DES:
		return xdes.Decrypt3DES(data, key)
	case SCTYPE_AES_CBC:
		b := xaes.AesDecryptCBC(data, []byte(key))
		return string(b), nil
	case SCTYPE_AES_ECB:
		b := xaes.AesDecryptECB(data, []byte(key))
		return string(b), nil
	case SCTYPE_AES_CFB:
		b := xaes.AesDecryptCFB(data, []byte(key))
		return string(b), nil
	case SCTYPE_RC4:
		b := xrc4.Decrypt(data, []byte(key))
		return string(b), nil
	default:
		return "", errors.New("sctype not exist.")
	}
}
