package xsecret

import (
	"fmt"
	"testing"

	"github.com/zhu8jie/gopkg/xsecret/xaes"
	"github.com/zhu8jie/gopkg/xsecret/xhex"
)

func TestEncrypt(t *testing.T) {

	// sss, err := Encrypt(SCTYPE_AES_ECB, "1000", "123456789abcdefg")
	// if err != nil {
	// 	t.Error(err)
	// }
	// if sss != "Fp599QWzjlqGb0SoX+RFBA==" {
	// 	t.Errorf("%v not equal %v", sss, "Fp599QWzjlqGb0SoX+RFBA==")
	// }

	str := "100"
	key := "u5OvRj55tiWz0tLS"
	// iv := "FD2718DD5C312460"
	fmt.Println(len(key))
	tmpKey := []byte(key)
	fmt.Println(len(tmpKey), tmpKey)

	// encryptStr, err := Encrypt(SCTYPE_AES_ECB, str, key)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	encryptStr := (xaes.AesEncryptECB([]byte(str), []byte(key)))

	s := xhex.HexEncode(string(encryptStr))

	fmt.Println(string(s))
}
