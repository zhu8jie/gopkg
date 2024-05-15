package xaes

import (
	"fmt"
	"testing"

	"github.com/zhu8jie/gopkg/xsecret/xhex"
)

func TestAesEncryptCBC(t *testing.T) {
	str := "66"
	key := "f14a77bbba9389st"
	iv := "FD2718DD5C312460"
	fmt.Println(len(key))
	tmpKey := []byte(key)
	fmt.Println(len(tmpKey), tmpKey)

	encryptStr := AesEncryptCBC([]byte(str), []byte(key), []byte(iv))

	fmt.Println(string(xhex.HexEncode(string(encryptStr))))
}
