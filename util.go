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

package graphlib

import (
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStr(n int, s rand.Source) string {
	//s := rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func readFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	return io.ReadAll(f)
}

func getMaxValue[W number](n W) W {
	switch any(n).(type) {
	case int:
		return any(math.MaxInt).(W)
	case int8:
		return any(math.MaxInt8).(W)
	case int16:
		return any(math.MaxInt16).(W)
	case int32:
		return any(math.MaxInt32).(W)
	case int64:
		return any(math.MaxInt64).(W)
	case uint:
		return any(math.MaxInt).(W)
	case uint8:
		return any(math.MaxUint8).(W)
	case uint16:
		return any(math.MaxUint16).(W)
	case uint32:
		return any(math.MaxUint32).(W)
	case uint64:
		return any(math.MaxInt64).(W)
	case float32:
		return any(math.MaxFloat32).(W)
	case float64:
		return any(math.MaxFloat64).(W)
	default:
		return n
	}
}

func edgeFormat[K comparable](v1, v2 K) K {
	switch any(v1).(type) {
	case string, []byte:
		return any(fmt.Sprintf("%v-%v-%s", v1, v2, randStr(7, rand.NewSource(rand.Int63())))).(K)
	case int:
		return any(rand.Int()).(K)
	case int8:
		return any(int8(rand.Int())).(K)
	case int16:
		return any(int16(rand.Int())).(K)
	case int32:
		return any(rand.Int31()).(K)
	case int64:
		return any(rand.Int63()).(K)
	case uint:
		return any(uint(rand.Uint64())).(K)
	case uint8:
		return any(uint8(rand.Uint32())).(K)
	case uint16:
		return any(uint16(rand.Uint32())).(K)
	case uint32:
		return any(rand.Uint32()).(K)
	case uint64:
		return any(rand.Uint64()).(K)
	default:
		return v1
	}
}

func maxValue[N number](n N) N {
	switch any(n).(type) {
	case int:
		return any(math.MaxInt).(N)
	case int8:
		return any(math.MaxInt8).(N)
	case int16:
		return any(math.MaxInt16).(N)
	case int32:
		return any(math.MaxInt32).(N)
	case int64:
		return any(math.MaxInt64).(N)
	case uint:
		return any(math.MaxInt).(N)
	case uint8:
		return any(math.MaxUint8).(N)
	case uint16:
		return any(math.MaxUint16).(N)
	case uint32:
		return any(math.MaxUint32).(N)
	case uint64:
		return any(math.MaxInt64).(N)
	case float32:
		return any(math.MaxFloat32).(N)
	case float64:
		return any(math.MaxFloat64).(N)
	default:
		return n
	}
}

func minValue[N number](n N) N {
	switch any(n).(type) {
	case int:
		return any(math.MinInt).(N)
	case int8:
		return any(math.MinInt8).(N)
	case int16:
		return any(math.MinInt16).(N)
	case int32:
		return any(math.MinInt32).(N)
	case int64:
		return any(math.MinInt64).(N)
	case uint:
		return any(math.MinInt).(N)
	case uint8:
		return any(uint8(0)).(N)
	case uint16:
		return any(uint16(0)).(N)
	case uint32:
		return any(uint32(0)).(N)
	case uint64:
		return any(uint64(0)).(N)
	case float32:
		return any(math.SmallestNonzeroFloat32).(N)
	case float64:
		return any(math.SmallestNonzeroFloat64).(N)
	default:
		return n
	}
}

var (
	errRunTimeout = errors.New("function run timeout")
)

func runWithTimeout(timeout time.Duration, f func() error) error {
	tr := time.NewTimer(timeout)
	defer tr.Stop()

	ch := make(chan error)
	go func() {
		defer close(ch)
		ch <- f()
	}()
	select {
	case <-tr.C:
		return errRunTimeout
	case err, ok := <-ch:
		if !ok {
			return nil
		}
		return err
	}
}

func runWithRetry(retry int, timeout time.Duration, f func() error) error {
	if retry <= 0 && timeout == time.Duration(0) {
		return f()
	} else if retry <= 0 {
		return runWithTimeout(timeout, f)
	} else if timeout == time.Duration(0) {
		var err error
		for i := 0; i <= retry; i++ {
			if err = f(); err == nil {
				return nil
			}
		}
		return fmt.Errorf("function runs exceeds the retry limit %d, %v", retry, err)
	} else {
		var err error
		for i := 0; i <= retry; i++ {
			if err = runWithTimeout(timeout, f); err == nil {
				return nil
			}
		}
		return fmt.Errorf("function runs exceeds the retry limit %d, %v", retry, err)
	}
}

func pow(a, n int) int {
	if a == 0 {
		return 0
	}
	if n == 0 {
		return 1
	}
	res, fac := 1, a
	for n > 0 {
		if n&1 == 1 {
			res *= fac
		}
		fac *= fac
		n >>= 1
	}
	return res
}

type stack[K comparable] struct {
	elems []K
	idx   int
}

func newStack[K comparable]() *stack[K] {
	return &stack[K]{}
}

func (s *stack[K]) size() int {
	return s.idx
}

func (s *stack[K]) empty() bool {
	return s.idx == 0
}

func (s *stack[K]) push(k K) {
	if s.idx < len(s.elems) {
		s.elems[s.idx] = k
	} else {
		s.elems = append(s.elems, k)
	}
	s.idx++
}

func (s *stack[K]) pop() (K, bool) {
	var k K
	if !s.empty() {
		k = s.elems[s.idx-1]
		s.idx--
		return k, true
	}
	return k, false
}

func (s *stack[K]) contains(k K) bool {
	for i := 0; i < s.idx; i++ {
		if s.elems[i] == k {
			return true
		}
	}
	return false
}

func (s *stack[K]) top() (k K) {
	if s.idx > 0 {
		return s.elems[s.idx-1]
	}
	return
}

func (s *stack[K]) clean() {
	s.idx = 0
}

type fifo[K comparable] struct {
	elems []K
	head  int
	tail  int
}

func newFIFO[K comparable]() *fifo[K] {
	return &fifo[K]{}
}

func (f *fifo[K]) size() int {
	return f.tail - f.head
}

func (f *fifo[K]) empty() bool {
	return f.head == f.tail
}

func (f *fifo[K]) push(k K) {
	if f.tail < len(f.elems) {
		f.elems[f.tail] = k
	} else {
		f.elems = append(f.elems, k)
	}
	f.tail++
}

func (f *fifo[K]) pop() (K, bool) {
	var k K
	if !f.empty() {
		k = f.elems[f.head]
		f.head++
		return k, true
	}
	return k, false
}

func (f *fifo[K]) front() (K, bool) {
	var k K
	if !f.empty() {
		return f.elems[f.head], true
	}
	return k, false
}

func (f *fifo[K]) back() (K, bool) {
	var k K
	if !f.empty() {
		return f.elems[f.tail-1], true
	}
	return k, false
}

type Stack[T any] struct {
	elems []T
	idx   int
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{}
}

func (s *Stack[T]) Len() int {
	return s.idx
}

func (s *Stack[K]) IsEmpty() bool {
	return s.idx == 0
}

func (s *Stack[T]) Push(v T) {
	if s.idx < len(s.elems) {
		s.elems[s.idx] = v
	} else {
		s.elems = append(s.elems, v)
	}
	s.idx++
}

func (s *Stack[T]) Pop() (T, bool) {
	var k T
	if s.idx > 0 {
		k = s.elems[s.idx-1]
		s.idx--
		return k, true
	}
	return k, false
}

func (s *Stack[T]) Contains(k T, comp CompareFunc[T]) bool {
	for i := 0; i < s.idx; i++ {
		if comp(s.elems[i], k) == 0 {
			return true
		}
	}
	return false
}

func (s *Stack[T]) Top() (v T) {
	if s.idx > 0 {
		return s.elems[s.idx-1]
	}
	return
}

func (s *Stack[T]) Clean() {
	s.idx = 0
}

type FIFO[T any] struct {
	elems []T
	head  int
	tail  int
}

func NewFIFO[T any]() *FIFO[T] {
	return &FIFO[T]{}
}

func (f *FIFO[T]) Len() int {
	return f.tail - f.head
}

func (f *FIFO[T]) IsEmpty() bool {
	return f.head == f.tail
}

func (f *FIFO[T]) Push(k T) {
	if f.tail < len(f.elems) {
		f.elems[f.tail] = k
	} else {
		f.elems = append(f.elems, k)
	}
	f.tail++
}

func (f *FIFO[T]) Pop() (k T, ok bool) {
	if f.head != f.tail {
		k = f.elems[f.head]
		f.head++
		return k, true
	}
	return k, false
}

func (f *FIFO[T]) Front() (k T, ok bool) {
	if f.head != f.tail {
		return f.elems[f.head], true
	}
	return k, false
}

func (f *FIFO[T]) Back() (k T, ok bool) {
	if f.head != f.tail {
		return f.elems[f.tail-1], true
	}
	return k, false
}

func (f *FIFO[T]) Clean() {
	f.head, f.tail = 0, 0
}
