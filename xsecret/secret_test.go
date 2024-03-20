package xsecret

import (
	"fmt"
	"strconv"
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

	dspWinPrice := 500
	var dspWinPriceEcb string
	if dspWinPrice > 0 {
		b := xaes.AesEncryptECB([]byte(strconv.Itoa(dspWinPrice)), []byte("kJDnjt7MuK8wQ7K6"))
		dspWinPriceEcb = string(xhex.HexEncode(string(b)))
	}
	fmt.Println("---", dspWinPriceEcb)
}
