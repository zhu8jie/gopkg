package xutils

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	s := Md5("xxx999995fdsfsdandroid1704254393")
	fmt.Println(s)
}
