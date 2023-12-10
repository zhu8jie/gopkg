package xsecret

import (
	"testing"
)

func TestEncrypt(t *testing.T) {

	sss, err := Encrypt(SCTYPE_AES_ECB, "1000", "123456789abcdefg")
	if err != nil {
		t.Error(err)
	}
	if sss != "Fp599QWzjlqGb0SoX+RFBA==" {
		t.Errorf("%v not equal %v", sss, "Fp599QWzjlqGb0SoX+RFBA==")
	}
}
