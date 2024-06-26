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
		return any(fmt.Sprintf("%v-%v", v1, v2)).(K)
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

type item[K comparable] struct {
	key   K
	value float64
	index int // The index of the item in the heap.
}

/*
type costHeap[K comparable] []*item[K]

func (h costHeap[K]) Len() int { return len(h) }

func (h costHeap[K]) Less(i, j int) bool {
	return h[i].value > h[j].value
}

func (h costHeap[K]) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *costHeap[K]) Push(x any) {
	n := len(*h)
	v := x.(*item[K])
	v.index = n
	*h = append(*h, v)
}

func (h *costHeap[K]) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*h = old[0 : n-1]
	return item
}

type costQueue[K comparable] struct {
	items    map[K]*item[K]
	priority costHeap[K]
}

func newCostQueue[K comparable]() *costQueue[K] {
	p := &costQueue[K]{
		items: make(map[K]*item[K]),
	}
	p.Init()
	return p
}

// update modifies the priority and value of an Item in the queue.
func (pq *costQueue[K]) Update(k K, priority float64) {
	v, ok := pq.items[k]
	if ok {
		v.value = priority
		heap.Fix(&pq.priority, v.index)
	}
}

func (pq *costQueue[K]) Push(item *item[K]) {
	v, ok := pq.items[item.key]
	if ok {
		v.value = item.value
		heap.Fix(&pq.priority, v.index)
		return
	}
	//
	pq.items[item.key] = item
	pq.priority.Push(item)
	heap.Fix(&pq.priority, item.index)
}

func (pq *costQueue[K]) Pop() *item[K] {
	if len(pq.items) != 0 {
		v := pq.priority.Pop().(*item[K])
		if v != nil {
			delete(pq.items, v.key)
		}
		return v
	}
	return nil
}

func (pq *costQueue[K]) Get(k K) float64 {
	v, ok := pq.items[k]
	if ok {
		return v.value
	}
	return 0.0
}

func (pq *costQueue[K]) Len() int {
	return pq.priority.Len()
}

func (pq *costQueue[K]) Init() {
	heap.Init(&pq.priority)
}
*/

type stack[K comparable] struct {
	elems []K
	top   int
}

func newStack[K comparable]() *stack[K] {
	return &stack[K]{}
}

func (s *stack[K]) empty() bool {
	return s.top == 0
}

func (s *stack[K]) push(k K) {
	if s.top < len(s.elems) {
		s.elems[s.top] = k
	} else {
		s.elems = append(s.elems, k)
	}
	s.top++
}

func (s *stack[K]) pop() (K, bool) {
	var k K
	if !s.empty() {
		k = s.elems[s.top-1]
		s.top--
		return k, true
	}
	return k, false
}

func (s *stack[K]) contains(k K) bool {
	for i := 0; i < s.top; i++ {
		if s.elems[i] == k {
			return true
		}
	}
	return false
}

type fifo[K comparable] struct {
	elems []K
	head  int
	tail  int
}

func newFIFO[K comparable]() *fifo[K] {
	return &fifo[K]{}
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
