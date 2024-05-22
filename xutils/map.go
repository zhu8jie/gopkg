package xutils

import (
	"errors"
	"sync"
)

type MaxMap struct {
	m       map[string]MapVal
	l       sync.Mutex
	maxCnt  int
	overtop bool
}

type MapVal struct {
	v   interface{}
	err error
}

var LockMapNotExist error = errors.New("key is not exist.")

func NewMaxMap(maxCnt int) (*MaxMap, error) {
	if maxCnt == 0 {
		return nil, errors.New("maxCnt must > 0.")
	}
	ret := &MaxMap{
		m:       make(map[string]MapVal, maxCnt+1000),
		maxCnt:  maxCnt,
		l:       sync.Mutex{},
		overtop: false,
	}
	return ret, nil
}

func (v *MaxMap) Count() int {
	v.l.Lock()
	defer v.l.Unlock()
	return len(v.m)
}

func (v *MaxMap) Overtop() bool {
	return v.overtop
}

func (v *MaxMap) Store(key string, value interface{}) {
	v.l.Lock()
	defer v.l.Unlock()
	v.m[key] = MapVal{
		v: value,
	}
	// fmt.Println("store len", len(v.m))
	if len(v.m) >= v.maxCnt {
		v.overtop = true
	}
}

func (v *MaxMap) Load(key string) MapVal {
	v.l.Lock()
	defer v.l.Unlock()
	ret, exist := v.m[key]
	if !exist {
		return MapVal{
			err: LockMapNotExist,
		}
	}
	return ret
}

func (v *MaxMap) LoadOrStore(key string, value interface{}) (ret MapVal, loaded bool) {
	l := v.Load(key)
	if l.err != nil {
		v.Store(key, value)
		return MapVal{
			v: value,
		}, false
	} else {
		return l, true
	}
}

func (v MapVal) Error() error {
	return v.err
}

func (v MapVal) String() (string, bool) {
	ret, ok := v.v.(string)
	return ret, ok
}

func (v MapVal) Int() (int, bool) {
	ret, ok := v.v.(int)
	return ret, ok
}

func (v MapVal) Int64() (int64, bool) {
	ret, ok := v.v.(int64)
	return ret, ok
}

func (v MapVal) Float64() (float64, bool) {
	ret, ok := v.v.(float64)
	return ret, ok
}

func (v MapVal) Val() interface{} {
	return v.v
}
