package bloomx

import (
	"github.com/bits-and-blooms/bitset"
	"github.com/zhu8jie/gopkg/xutils"
)

const DEFAULT_SIZE = 2 << 30

var seeds = []uint{7, 11, 13, 31, 37, 61}

type BloomFilter struct {
	Set   [2]*bitset.BitSet
	Funcs [6]SimpleHash
}

func NewBloomFilter() *BloomFilter {
	bf := new(BloomFilter)
	for i := 0; i < len(bf.Funcs); i++ {
		bf.Funcs[i] = SimpleHash{DEFAULT_SIZE, seeds[i]}
	}
	for i := 0; i < len(bf.Set); i++ {
		bf.Set[i] = bitset.New(DEFAULT_SIZE)
	}
	return bf
}

func (bf BloomFilter) Add(value string) {
	idx := int(xutils.Crc(value)) % len(bf.Set)
	for _, f := range bf.Funcs {
		bf.Set[idx].Set(f.hash(value))
	}
}

func (bf BloomFilter) Contains(value string) bool {
	if value == "" {
		return false
	}
	idx := int(xutils.Crc(value)) % len(bf.Set)
	ret := true
	for _, f := range bf.Funcs {
		ret = ret && bf.Set[idx].Test(f.hash(value))
	}
	return ret
}

type SimpleHash struct {
	Cap  uint
	Seed uint
}

func (s SimpleHash) hash(value string) uint {
	var result uint = 0
	for i := 0; i < len(value); i++ {
		result = result*s.Seed + uint(value[i])
	}
	return (s.Cap - 1) & result
}
