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

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

// A differential array is a powerful data structure used to efficiently perform range update operations on an array.
// It allows multiple updates to be applied in constant time, and the final array can be reconstructed in linear time.
// This approach is particularly useful in scenarios involving frequent range updates and queries.
type DifferenceArray[N number] struct {
	diff []N
}

// Initialize a differential array
func NewDifferenceArray[N number](nums []N) *DifferenceArray[N] {
	d := &DifferenceArray[N]{
		diff: make([]N, len(nums)),
	}
	if len(nums) > 0 {
		d.diff[0] = nums[0]
		for i := 1; i < len(nums); i++ {
			d.diff[i] = nums[i] - nums[i-1]
		}
	}
	return d
}

// Append elements to the end of the array, note that this method has a time complexity of O(n)
func (d *DifferenceArray[N]) Append(n N) {
	if len(d.diff) == 0 {
		d.diff = append(d.diff, n)
	} else {
		t := d.diff[0]
		for i := 1; i < len(d.diff); i++ {
			t = t + d.diff[i]
		}
		d.diff = append(d.diff, n-t)
	}
}

func (d *DifferenceArray[N]) Len() int {
	return len(d.diff)
}

// Add n to all elements in the closed interval [left,right] (if n is negative, it means subtract |n| from all elements).
func (d *DifferenceArray[N]) Add(n N, left, right int) bool {
	if left < 0 || right >= len(d.diff) || left > right {
		return false
	}
	d.diff[left] += n
	if right+1 < len(d.diff) {
		d.diff[right+1] -= n
	}
	return true
}

// sum of closed interval [left,right].
func (d *DifferenceArray[N]) Sum(left, right int) (N, bool) {
	var a, b N
	if left < 0 || right >= len(d.diff) || left > right {
		return a, false
	}
	for i := 0; i <= right; i++ {
		if i < left {
			a += d.diff[i]
		}
		b += d.diff[i]
	}
	return b - a, true
}

// Return the current array.
func (d *DifferenceArray[N]) Array() []N {
	arr := make([]N, len(d.diff))
	if len(d.diff) > 0 {
		arr[0] = d.diff[0]
	}
	for i := 1; i < len(d.diff); i++ {
		arr[i] = arr[i-1] + d.diff[i]
	}
	return arr
}

// use subscripts to query elements
func (d *DifferenceArray[N]) Get(i int) (N, bool) {
	var n N
	if i < 0 || i >= len(d.diff) {
		return n, false
	}
	for j := 0; j <= i; j++ {
		n += d.diff[j]
	}
	return n, true
}

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
