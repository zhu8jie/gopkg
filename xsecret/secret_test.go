package xsecret

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {

	// sss, err := Encrypt(SCTYPE_AES_ECB, "1000", "123456789abcdefg")
	// if err != nil {
	// 	t.Error(err)
	// }
	// if sss != "Fp599QWzjlqGb0SoX+RFBA==" {
	// 	t.Errorf("%v not equal %v", sss, "Fp599QWzjlqGb0SoX+RFBA==")
	// }

	str := "110"
	key := "d39656591bc2c499cddba38da2b9da38"
	// iv := "FD2718DD5C312460"
	fmt.Println(len(key))
	tmpKey := []byte(key)
	fmt.Println(len(tmpKey), tmpKey)

	encryptStr, err := Encrypt(SCTYPE_AES_ECB, str, key)
	if err != nil {
		fmt.Println(err)
	}
	// s := xbase64.Base64Encrypt(encryptStr)

	fmt.Println(encryptStr)
}
