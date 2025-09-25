package googleopenrtb

import (
	"fmt"
	"testing"
)

func TestGoogleEncrypt(t *testing.T) {
	// Ikey= 0hFbKtstsJ8pPviFxtS6sedi
	// eKey= 60yoA261AzXN5htjIlvxadsf
	key, err := NewKeys([]byte("0hFbKtstsJ8pPviFxtS6sedi"), []byte("60yoA261AzXN5htjIlvxadsf"))
	if err != nil {
		fmt.Println("NewKeys error:", err)
	}

	pri := NewPrice(key)
	str, err := pri.EncodePriceValue(1000, []byte(""))
	if err != nil {
		fmt.Println("EncodePriceValue error:", err)
	}
	fmt.Println(str)
}
