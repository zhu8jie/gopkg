package xtask

type Task interface {
	DoSomeThingFirst() error
	DoSomeThingSecond() error
	DoSomeThingThird() error
}
