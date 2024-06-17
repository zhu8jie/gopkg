package xaes

import (
	"fmt"
	"testing"

	"github.com/zhu8jie/gopkg/xsecret/xhex"
)

func TestAesEncryptCBC(t *testing.T) {
	str := "110"
	key := "d39656591bc2c499cddba38da2b9da38"
	// iv := "FD2718DD5C312460"
	// fmt.Println(len(key))
	// tmpKey := []byte(key)
	// fmt.Println(len(tmpKey), tmpKey)

	encryptStr := AesEncryptECB([]byte(str), []byte(key))
	s := xhex.HexEncode(string(encryptStr))

	fmt.Println(string(s))
}
