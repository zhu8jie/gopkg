package xaes

import (
	"fmt"
	"testing"

	"github.com/zhu8jie/gopkg/xsecret/xbase64"
)

func TestAesEncryptCBC(t *testing.T) {
	str := `{"group":3}`
	key := "ZcK$BtWUj54^AR83"
	fmt.Println(len(key))
	tmpKey := []byte(key)
	fmt.Println(len(tmpKey), tmpKey)

	encryptStr := AesEncryptCBC([]byte(str), []byte(key))

	fmt.Println(xbase64.Base64Encrypt(encryptStr))
}
