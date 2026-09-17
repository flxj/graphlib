/*
	Copyright (C) 2023 flxj(https://github.com/flxj)

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package collection

import (
	"encoding/binary"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"math"
	"sync/atomic"
)

func hashPair(data []byte, seed1, seed2 uint64) (uint64, uint64) {
	h1 := fnv.New64a()
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], seed1)
	h1.Write(buf[:])
	h1.Write(data)
	v1 := h1.Sum64()

	h2 := fnv.New64()
	binary.LittleEndian.PutUint64(buf[:], seed2)
	h2.Write(buf[:])
	h2.Write(data)
	v2 := h2.Sum64()

	return v1, v2
}

func appendUint64(buf []byte, v uint64) []byte {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], v)
	return append(buf, b[:]...)
}

type BloomFilter struct {
	arr   *BitArray // bit array
	m     uint64    // length of bitarray
	k     uint64    // Number of hash functions
	seed1 uint64    // hash seed 1
	seed2 uint64    // hash seed 2
}

// Create a filter based on the expected number of elements n
// and the false positive rate p
func NewBloomFilter(n uint64, p float64) *BloomFilter {
	if n == 0 || p <= 0 || p >= 1 {
		return nil
	}
	// Optimal bits：m = -n * ln(p) / (ln2)^2
	ln2 := math.Ln2
	m := uint64(math.Ceil(-float64(n) * math.Log(p) / (ln2 * ln2)))
	// Optimal hash number：k = m/n * ln2
	k := max(uint64(math.Round(float64(m)/float64(n)*ln2)), 1)
	if m < 1 {
		m = 1
	}
	return NewBloomFilterWithParams(m, k)
}

func NewBloomFilterWithParams(m, k uint64) *BloomFilter {
	if m == 0 || k == 0 {
		return nil
	}
	return &BloomFilter{
		arr:   NewBitArray(int(m)),
		m:     m,
		k:     k,
		seed1: 0x9e3779b97f4a7c15,
		seed2: 0xc2b2ae3d27d4eb4f,
	}
}

func (f *BloomFilter) M() uint64  { return f.m }
func (f *BloomFilter) K() uint64  { return f.k }
func (f *BloomFilter) Count() int { return f.arr.Count() }

func (f *BloomFilter) Add(data []byte) {
	h1, h2 := hashPair(data, f.seed1, f.seed2)
	for i := uint64(0); i < f.k; i++ {
		// g_i = h1 + i*h2 (mod m)
		pos := (h1 + i*h2) % f.m
		f.arr.Set(int(pos))
	}
}

func (f *BloomFilter) AddString(s string) {
	f.Add([]byte(s))
}

func (f *BloomFilter) Contains(data []byte) bool {
	h1, h2 := hashPair(data, f.seed1, f.seed2)
	for i := uint64(0); i < f.k; i++ {
		pos := (h1 + i*h2) % f.m
		if !f.arr.Get(int(pos)) {
			return false
		}
	}
	return true
}

func (f *BloomFilter) ContainsString(s string) bool {
	return f.Contains([]byte(s))
}

// Estimation of misjudgment rate.
func (f *BloomFilter) FalsePositiveRate() float64 {
	set := float64(f.arr.Count())
	ratio := set / float64(f.m)
	return math.Pow(ratio, float64(f.k))
}

// Marshal format:
// [8]m | [8]k | [8]seed1 | [8]seed2 | [8]wordLen | words...
func (f *BloomFilter) Marshal() []byte {
	wordLen := uint64(len(f.arr.bits))
	buf := make([]byte, 0, 8*5+int(wordLen)*8)
	buf = appendUint64(buf, f.m)
	buf = appendUint64(buf, f.k)
	buf = appendUint64(buf, f.seed1)
	buf = appendUint64(buf, f.seed2)
	buf = appendUint64(buf, wordLen)
	for _, w := range f.arr.bits {
		buf = appendUint64(buf, w)
	}
	return buf
}

func (f *BloomFilter) Unmarshal(data []byte) error {
	if len(data) < 8*5 {
		return errors.New("bloom: data too short")
	}
	f.m = binary.LittleEndian.Uint64(data[0:8])
	f.k = binary.LittleEndian.Uint64(data[8:16])
	f.seed1 = binary.LittleEndian.Uint64(data[16:24])
	f.seed2 = binary.LittleEndian.Uint64(data[24:32])
	wordLen := binary.LittleEndian.Uint64(data[32:40])

	if uint64(len(data)) < 8*5+wordLen*8 {
		return errors.New("bloom: truncated data")
	}
	words := make([]uint64, wordLen)
	for i := uint64(0); i < wordLen; i++ {
		words[i] = binary.LittleEndian.Uint64(data[40+i*8:])
	}
	f.arr = &BitArray{
		bits: words,
		size: int(f.m),
	}
	return nil
}

func (f *BloomFilter) WriteTo(w io.Writer) (int64, error) {
	data := f.Marshal()
	n, err := w.Write(data)
	return int64(n), err
}

func (f *BloomFilter) ReadFrom(r io.Reader) (int64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return 0, err
	}
	if err := f.Unmarshal(data); err != nil {
		return 0, err
	}
	return int64(len(data)), nil
}

func (f *BloomFilter) String() string {
	return fmt.Sprintf("BloomFilter{m=%d, k=%d, bits=%d, fpRate=%.4f}",
		f.m, f.k, f.arr.Count(), f.FalsePositiveRate())
}

type syncBloomFilter struct {
	arr    *SyncBitArray // bit array
	m      uint64        // length of bitarray
	k      uint64        // Number of hash functions
	seed1  uint64        // hash seed 1
	seed2  uint64        // hash seed 2
	addCnt atomic.Int64
}

// Create a filter based on the expected number of elements n
// and the false positive rate p
func newSyncBloomFilter(n uint64, p float64) *syncBloomFilter {
	if n == 0 || p <= 0 || p >= 1 {
		return nil
	}
	// Optimal bits：m = -n * ln(p) / (ln2)^2
	ln2 := math.Ln2
	m := uint64(math.Ceil(-float64(n) * math.Log(p) / (ln2 * ln2)))
	// Optimal hash number：k = m/n * ln2
	k := uint64(math.Round(float64(m) / float64(n) * ln2))
	if k < 1 {
		k = 1
	}
	if m < 1 {
		m = 1
	}
	return &syncBloomFilter{
		arr:   NewSyncBitArray(int(m)),
		m:     m,
		k:     k,
		seed1: 0x9e3779b97f4a7c15,
		seed2: 0xc2b2ae3d27d4eb4f,
	}
}

func (f *syncBloomFilter) count() int { return f.arr.Count() }

func (f *syncBloomFilter) add(data []byte) {
	h1, h2 := hashPair(data, f.seed1, f.seed2)
	for i := uint64(0); i < f.k; i++ {
		// g_i = h1 + i*h2 (mod m)
		pos := (h1 + i*h2) % f.m
		f.arr.Set(int(pos))
	}
	f.addCnt.Add(1)
}

func (f *syncBloomFilter) contains(data []byte) bool {
	h1, h2 := hashPair(data, f.seed1, f.seed2)
	for i := uint64(0); i < f.k; i++ {
		pos := (h1 + i*h2) % f.m
		if !f.arr.Get(int(pos)) {
			return false
		}
	}
	return true
}

func (f *syncBloomFilter) added() int64 {
	return f.addCnt.Load()
}

type ShardedBloomFilter struct {
	shards []*syncBloomFilter
	mask   uint64
}

func NewShardedBloomFilter(n uint64, p float64, shards int) *ShardedBloomFilter {
	if shards < 1 {
		shards = 1
	}
	s := 1
	for s < shards {
		s <<= 1
	}
	perN := n/uint64(s) + 1
	perP := p / 2
	if perP <= 0 {
		perP = p
	}

	sf := &ShardedBloomFilter{
		shards: make([]*syncBloomFilter, s),
		mask:   uint64(s - 1),
	}
	for i := 0; i < s; i++ {
		sf.shards[i] = newSyncBloomFilter(perN, perP)
	}
	return sf
}

func (sf *ShardedBloomFilter) Shards() int { return len(sf.shards) }

func (sf *ShardedBloomFilter) shardOf(data []byte) *syncBloomFilter {
	h := fnv.New64a()
	h.Write(data)
	return sf.shards[h.Sum64()&sf.mask]
}

func (sf *ShardedBloomFilter) Add(data []byte) {
	sf.shardOf(data).add(data)
}

func (sf *ShardedBloomFilter) AddString(s string) {
	sf.Add([]byte(s))
}

func (sf *ShardedBloomFilter) Contains(data []byte) bool {
	return sf.shardOf(data).contains(data)
}

func (sf *ShardedBloomFilter) ContainsString(s string) bool {
	return sf.Contains([]byte(s))
}

func (sf *ShardedBloomFilter) Count() int {
	total := 0
	for _, sh := range sf.shards {
		total += sh.count()
	}
	return total
}

func (sf *ShardedBloomFilter) Added() int64 {
	var total int64
	for _, sh := range sf.shards {
		total += sh.added()
	}
	return total
}
