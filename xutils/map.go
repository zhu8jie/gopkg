package xutils

import (
	"errors"
	"sync"
)

type MaxMap struct {
	m      map[string]MapVal
	l      sync.RWMutex
	maxCnt int
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
		m:      make(map[string]MapVal, maxCnt),
		maxCnt: maxCnt,
		l:      sync.RWMutex{},
	}
	return ret, nil
}

func (v *MaxMap) Count() int {
	return len(v.m)
}

func (v *MaxMap) Store(key string, value interface{}) {
	v.l.Lock()
	v.m[key] = MapVal{
		v: value,
	}
	v.l.Unlock()
}

func (v *MaxMap) Load(key string) MapVal {
	v.l.RLock()
	defer v.l.RUnlock()
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
