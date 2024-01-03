package xrc4

import (
	"fmt"
	"testing"
)

func TestDo(t *testing.T) {
	key := "abcdefg123456789"
	in := "300"

	enStr := Encrypt([]byte(in), []byte(key))
	fmt.Println(enStr)

	deStr := Decrypt(enStr, []byte(key))
	fmt.Println(string(deStr))

}
