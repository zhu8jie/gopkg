package xutils

import (
	"fmt"
	"testing"
	"time"
)

func TestMd5(t *testing.T) {
	s := Md5("xxx999995fdsfsdandroid1704254393")
	fmt.Println(s)
}

func TestMaxMap(t *testing.T) {
	a, _ := NewMaxMap(10000)

	go func() {
		for i := 0; i < 20000; i++ {
			go func(num int) {
				a.LoadOrStore("a_"+IntToStr(num), struct{}{})
			}(i)
		}
	}()

	for {
		time.Sleep(time.Millisecond * 500)
		fmt.Println(a.Count())
		if a.Overtop() {
			a, _ = NewMaxMap(10000)
			fmt.Println("yanw_test")
			break
		}
	}
	time.Sleep(time.Second * 10)
	fmt.Println("done...")
}
