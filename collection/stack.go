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
