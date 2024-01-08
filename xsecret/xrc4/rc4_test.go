package xrc4

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	key := ""
	in := "300"

	enStr, err := Encrypt([]byte(in), []byte(key))
	fmt.Println(enStr, err)

	deStr, err := Decrypt(enStr, []byte(key))
	fmt.Println(string(deStr), err)

}
