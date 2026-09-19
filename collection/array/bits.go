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

package array

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/bits"
	"sync"
	"sync/atomic"
)

// A Bitarray Implementation.
type BitArray struct {
	size int
	bits []uint64
}

// Note that if size<0, a null pointer will be returned.
func NewBitArray(size int) *BitArray {
	if size < 0 {
		return nil
	}
	return &BitArray{
		size: size,
		bits: make([]uint64, (size+63)/64),
	}
}

func (b *BitArray) Size() int {
	return b.size
}

func (b *BitArray) Array() []uint64 {
	res := make([]uint64, len(b.bits))
	copy(res, b.bits)
	return res
}

func (b *BitArray) Set(i int) {
	b.checkIndex(i)
	b.bits[i/64] |= 1 << uint(i%64)
}

func (b *BitArray) Clear(i int) {
	b.checkIndex(i)
	b.bits[i/64] &^= 1 << uint(i%64)
}

func (b *BitArray) SetTo(i int, v bool) {
	if v {
		b.Set(i)
	} else {
		b.Clear(i)
	}
}

func (b *BitArray) Get(i int) bool {
	b.checkIndex(i)
	return b.bits[i/64]&(1<<uint(i%64)) != 0
}

func (b *BitArray) Flip(i int) {
	b.checkIndex(i)
	b.bits[i/64] ^= 1 << uint(i%64)
}

func (b *BitArray) Count() int {
	n := 0
	for _, w := range b.bits {
		n += bits.OnesCount64(w)
	}
	return n
}

func (b *BitArray) ClearAll() {
	for i := range b.bits {
		b.bits[i] = 0
	}
}

func (b *BitArray) SetAll() {
	for i := range b.bits {
		b.bits[i] = ^uint64(0)
	}
	b.trim()
}

func (b *BitArray) trim() {
	if b.size == 0 {
		return
	}
	rem := b.size % 64
	if rem != 0 {
		b.bits[len(b.bits)-1] &= (1 << uint(rem)) - 1
	}
}

func (b *BitArray) checkIndex(i int) {
	if i < 0 || i >= b.size {
		panic(fmt.Sprintf("bitarray: index %d out of range [0,%d)", i, b.size))
	}
}

func (b *BitArray) And(arr *BitArray) *BitArray {
	b.checkSameSize(arr)
	res := NewBitArray(b.size)
	for i := range b.bits {
		res.bits[i] = b.bits[i] & arr.bits[i]
	}
	return res
}

func (b *BitArray) Or(arr *BitArray) *BitArray {
	b.checkSameSize(arr)
	res := NewBitArray(b.size)
	for i := range b.bits {
		res.bits[i] = b.bits[i] | arr.bits[i]
	}
	return res
}

func (b *BitArray) Xor(arr *BitArray) *BitArray {
	b.checkSameSize(arr)
	res := NewBitArray(b.size)
	for i := range b.bits {
		res.bits[i] = b.bits[i] ^ arr.bits[i]
	}
	return res
}

func (b *BitArray) Not() *BitArray {
	res := NewBitArray(b.size)
	for i := range b.bits {
		res.bits[i] = ^b.bits[i]
	}
	res.trim()
	return res
}

func (b *BitArray) checkSameSize(arr *BitArray) {
	if b.size != arr.size {
		panic("bitarray: size mismatch")
	}
}

// ForEach traverses all bits set to 1, and the callback
// returns false to terminate prematurely.
func (b *BitArray) ForEach(fn func(i int) bool) {
	for i, w := range b.bits {
		for w != 0 {
			tz := bits.TrailingZeros64(w)
			idx := i*64 + tz
			if idx >= b.size {
				return
			}
			if !fn(idx) {
				return
			}
			w &= w - 1
		}
	}
}

func (b *BitArray) ToSlice() []int {
	var res []int
	b.ForEach(func(i int) bool {
		res = append(res, i)
		return true
	})
	return res
}

func (b *BitArray) String() string {
	buf := make([]byte, b.size)
	for i := 0; i < b.size; i++ {
		if b.Get(i) {
			buf[i] = '1'
		} else {
			buf[i] = '0'
		}
	}
	return string(buf)
}

type SyncBitArray struct {
	mu  sync.RWMutex
	arr *BitArray
}

func NewSyncBitArray(size int) *SyncBitArray {
	if size < 0 {
		panic("bitarray: negative size")
	}
	return &SyncBitArray{
		arr: NewBitArray(size),
	}
}

func (b *SyncBitArray) Size() int {
	return b.arr.Size()
}

func (b *SyncBitArray) Get(i int) bool {
	b.arr.checkIndex(i)
	w := atomic.LoadUint64(&b.arr.bits[i/64])
	return w&(uint64(1)<<uint(i%64)) != 0
}

func (b *SyncBitArray) Set(i int) {
	b.arr.checkIndex(i)
	w := &b.arr.bits[i/64]
	mask := uint64(1) << uint(i%64)
	for {
		old := atomic.LoadUint64(w)
		if old&mask != 0 {
			return
		}
		if atomic.CompareAndSwapUint64(w, old, old|mask) {
			return
		}
	}
}

func (b *SyncBitArray) Clear(i int) {
	b.arr.checkIndex(i)
	w := &b.arr.bits[i/64]
	mask := uint64(1) << uint(i%64)
	for {
		old := atomic.LoadUint64(w)
		if old&mask == 0 {
			return
		}
		if atomic.CompareAndSwapUint64(w, old, old&^mask) {
			return
		}
	}
}

func (b *SyncBitArray) Flip(i int) {
	b.arr.checkIndex(i)
	w := &b.arr.bits[i/64]
	mask := uint64(1) << uint(i%64)
	for {
		old := atomic.LoadUint64(w)
		if atomic.CompareAndSwapUint64(w, old, old^mask) {
			return
		}
	}
}

func (b *SyncBitArray) SetTo(i int, v bool) {
	if v {
		b.Set(i)
	} else {
		b.Clear(i)
	}
}

func (b *SyncBitArray) Count() int {
	n := 0
	for i := range b.arr.bits {
		n += bits.OnesCount64(atomic.LoadUint64(&b.arr.bits[i]))
	}
	return n
}

func (b *SyncBitArray) All() bool {
	for i := range b.arr.bits {
		w := atomic.LoadUint64(&b.arr.bits[i])
		if w != ^uint64(0) {
			return false
		}
	}
	return b.checkTailBits(true)
}

func (b *SyncBitArray) Any() bool {
	for i := range b.arr.bits {
		if atomic.LoadUint64(&b.arr.bits[i]) != 0 {
			return true
		}
	}
	return false
}

func (b *SyncBitArray) checkTailBits(want bool) bool {
	if b.arr.size == 0 {
		return true
	}
	rem := b.arr.size % 64
	if rem == 0 {
		return true
	}
	last := atomic.LoadUint64(&b.arr.bits[len(b.arr.bits)-1])
	highMask := ^((uint64(1) << uint(rem)) - 1)
	if want {
		_ = last
		return true
	}
	return last&highMask == 0
}

func (b *SyncBitArray) And(other *SyncBitArray) *SyncBitArray {
	b.arr.checkSameSize(other.arr)
	res := NewSyncBitArray(b.arr.size)
	for i := range b.arr.bits {
		x := atomic.LoadUint64(&b.arr.bits[i])
		y := atomic.LoadUint64(&other.arr.bits[i])
		atomic.StoreUint64(&res.arr.bits[i], x&y)
	}
	return res
}

func (b *SyncBitArray) Or(other *SyncBitArray) *SyncBitArray {
	b.arr.checkSameSize(other.arr)
	res := NewSyncBitArray(b.arr.size)
	for i := range b.arr.bits {
		x := atomic.LoadUint64(&b.arr.bits[i])
		y := atomic.LoadUint64(&other.arr.bits[i])
		atomic.StoreUint64(&res.arr.bits[i], x|y)
	}
	return res
}

func (b *SyncBitArray) Xor(other *SyncBitArray) *SyncBitArray {
	b.arr.checkSameSize(other.arr)
	res := NewSyncBitArray(b.arr.size)
	for i := range b.arr.bits {
		x := atomic.LoadUint64(&b.arr.bits[i])
		y := atomic.LoadUint64(&other.arr.bits[i])
		atomic.StoreUint64(&res.arr.bits[i], x^y)
	}
	return res
}

func (b *SyncBitArray) Not() *SyncBitArray {
	res := NewSyncBitArray(b.arr.size)
	for i := range b.arr.bits {
		w := atomic.LoadUint64(&b.arr.bits[i])
		atomic.StoreUint64(&res.arr.bits[i], ^w)
	}
	res.arr.trim()
	return res
}

func (b *SyncBitArray) ClearAll() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := range b.arr.bits {
		atomic.StoreUint64(&b.arr.bits[i], 0)
	}
}

func (b *SyncBitArray) SetAll() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for i := range b.arr.bits {
		atomic.StoreUint64(&b.arr.bits[i], ^uint64(0))
	}
	b.arr.trim()
}

func (b *SyncBitArray) Resize(newSize int) {
	if newSize < 0 {
		panic("bitarray: negative size")
	}
	b.mu.Lock()
	defer b.mu.Unlock()

	newWords := (newSize + 63) / 64
	oldWords := len(b.arr.bits)

	newBits := make([]uint64, newWords)
	copy(newBits, b.arr.bits)

	b.arr.bits = newBits
	b.arr.size = newSize

	if newSize < oldWords*64 {
		b.arr.trim()
	}
}

func (b *SyncBitArray) ForEach(fn func(i int) bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for wi := range b.arr.bits {
		w := atomic.LoadUint64(&b.arr.bits[wi])
		for w != 0 {
			tz := bits.TrailingZeros64(w)
			idx := wi*64 + tz
			if idx >= b.arr.size {
				return
			}
			if !fn(idx) {
				return
			}
			w &= w - 1
		}
	}
}

func (b *SyncBitArray) ToSlice() []int {
	var res []int
	b.ForEach(func(i int) bool {
		res = append(res, i)
		return true
	})
	return res
}

func (b *SyncBitArray) Marshal() []byte {
	b.mu.RLock()
	defer b.mu.RUnlock()

	size := b.arr.size
	words := (size + 63) / 64

	buf := make([]byte, 8+words*8)
	binary.LittleEndian.PutUint64(buf[0:8], uint64(size))

	for i := 0; i < words; i++ {
		var w uint64
		if i < len(b.arr.bits) {
			w = atomic.LoadUint64(&b.arr.bits[i])
		}
		binary.LittleEndian.PutUint64(buf[8+i*8:], w)
	}
	return buf
}

func (b *SyncBitArray) Unmarshal(data []byte) error {
	if len(data) < 8 {
		return errors.New("bitarray: data too short")
	}
	size := int(binary.LittleEndian.Uint64(data[0:8]))
	if size < 0 {
		return errors.New("bitarray: invalid size")
	}
	words := (size + 63) / 64
	if len(data) < 8+words*8 {
		return errors.New("bitarray: truncated data")
	}
	newBits := make([]uint64, words)
	for i := 0; i < words; i++ {
		newBits[i] = binary.LittleEndian.Uint64(data[8+i*8:])
	}
	b.mu.Lock()
	b.arr.bits = newBits
	b.arr.size = size
	b.mu.Unlock()
	return nil
}
