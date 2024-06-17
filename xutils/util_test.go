package xutils

import (
	"fmt"
	"testing"
)

// func TestMd5(t *testing.T) {
// 	s := Md5("xxx999995fdsfsdandroid1704254393")
// 	fmt.Println(s)
// }

// func TestMaxMap(t *testing.T) {
// 	a, _ := NewMaxMap(10000)

// 	go func() {
// 		for i := 0; i < 20000; i++ {
// 			go func(num int) {
// 				a.LoadOrStore("a_"+IntToStr(num), struct{}{})
// 			}(i)
// 		}
// 	}()

// 	for {
// 		time.Sleep(time.Millisecond * 500)
// 		fmt.Println(a.Count())
// 		if a.Overtop() {
// 			a, _ = NewMaxMap(10000)
// 			fmt.Println("yanw_test")
// 			break
// 		}
// 	}
// 	time.Sleep(time.Second * 10)
// 	fmt.Println("done...")
// }

func TestHash(t *testing.T) {
	s := "c828b1a1686bb999dc5789747f1737627b3dfbc0bd0892e7f04e300289bf34d7"
	m := HashString(s, 10)
	fmt.Println(m)
	fmt.Println(30 % 10)
}
