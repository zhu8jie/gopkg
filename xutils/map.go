package xutils

import (
	"errors"
	"sync"
)

type LockMap struct {
	m map[string]LockMapVal
	l sync.RWMutex
}

type LockMapVal struct {
	v   interface{}
	err error
}

var LockMapNotExist error = errors.New("key is not exist.")

func (v LockMap) Count() int {
	return len(v.m)
}

func (v LockMap) Store(key string, value interface{}) {
	v.l.Lock()
	v.m[key] = LockMapVal{
		v: value,
	}
	v.l.Unlock()
}

func (v LockMap) Load(key string) LockMapVal {
	v.l.RLock()
	defer v.l.RUnlock()
	ret, exist := v.m[key]
	if !exist {
		return LockMapVal{
			err: LockMapNotExist,
		}
	}
	return ret
}

func (v LockMap) StoreOrLoad(key string, value interface{}) (ret LockMapVal, loaded bool) {
	l := v.Load(key)
	if l.err != nil {
		v.Store(key, value)
		return LockMapVal{
			v: value,
		}, true
	} else {
		return l, false
	}
}

func (v LockMapVal) Error() error {
	return v.err
}

func (v LockMapVal) String() (string, bool) {
	ret, ok := v.v.(string)
	return ret, ok
}

func (v LockMapVal) Int() (int, bool) {
	ret, ok := v.v.(int)
	return ret, ok
}

func (v LockMapVal) Int64() (int64, bool) {
	ret, ok := v.v.(int64)
	return ret, ok
}

func (v LockMapVal) Float64() (float64, bool) {
	ret, ok := v.v.(float64)
	return ret, ok
}

func (v LockMapVal) Val() interface{} {
	return v.v
}
