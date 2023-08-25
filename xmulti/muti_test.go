package xmulti

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestXmulti(t *testing.T) {
	p := New(abc, 3, 4, 0, "", 0, nil)

	for i := 0; i < 30; i++ {
		o := ObjTest{
			Name: "test_" + strconv.Itoa(i),
		}
		p.Run(context.Background(), o)
	}
	time.Sleep(time.Second * 10)
}

type ObjTest struct {
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender int    `json:"gender"`
}

func abc(ctx context.Context, i []InputParam) error {
	for i := 0; i < 2; i++ {
		time.Sleep(time.Second * 1)
	}
	fmt.Println(time.Now().Unix())
	oName := make([]string, 0)
	for _, obj := range i {
		o, b := obj.(ObjTest)
		if !b {
			fmt.Println("XmultiDecode obj is error", obj)
		}
		oName = append(oName, o.Name)
	}
	fmt.Println("objTest:", oName)
	return nil
}
