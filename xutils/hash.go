package xutils

import (
	"hash/fnv"
)

func HashString(s string, mod uint32) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32() % mod
}
