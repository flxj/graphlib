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

package heap

import (
	"container/heap"

	"github.com/flxj/graphlib/collection"
)

type PriorityQueue[T, P any] struct {
	data []T
	prio []P
	less collection.Less[P]
}

func NewPriorityQueue[T, P any](less collection.Less[P]) *PriorityQueue[T, P] {
	return &PriorityQueue[T, P]{less: less}
}

func (q *PriorityQueue[T, P]) Init(data []T, prio []P) {
	q.data = data
	q.prio = prio
	n := len(data)
	for i := n/2 - 1; i >= 0; i-- {
		q.siftDown(i)
	}
}

func (q *PriorityQueue[T, P]) Push(v T, p P) {
	q.data = append(q.data, v)
	q.prio = append(q.prio, p)
	q.siftUp(len(q.prio) - 1)
}

func (q *PriorityQueue[T, P]) Pop() (v T, p P, ok bool) {
	n := len(q.data)
	if n == 0 {
		return
	}
	v, p, ok = q.data[0], q.prio[0], true
	q.data[0], q.prio[0] = q.data[n-1], q.prio[n-1]
	var zeroT T
	q.data[n-1] = zeroT
	q.data = q.data[:n-1]
	q.prio = q.prio[:n-1]
	if len(q.prio) > 0 {
		q.siftDown(0)
	}
	return
}

func (q *PriorityQueue[T, P]) Peek() (v T, p P, ok bool) {
	if len(q.data) == 0 {
		return
	}
	return q.data[0], q.prio[0], true
}

func (q *PriorityQueue[T, P]) Len() int {
	return len(q.data)
}

func (q *PriorityQueue[T, P]) siftUp(i int) {
	for i > 0 {
		p := (i - 1) / 2
		if !q.less(q.prio[i], q.prio[p]) {
			break
		}
		q.prio[i], q.prio[p] = q.prio[p], q.prio[i]
		q.data[i], q.data[p] = q.data[p], q.data[i]
		i = p
	}
}

func (q *PriorityQueue[T, P]) siftDown(i int) {
	n := len(q.data)
	for {
		l, r := 2*i+1, 2*i+2
		s := i
		if l < n && q.less(q.prio[l], q.prio[s]) {
			s = l
		}
		if r < n && q.less(q.prio[r], q.prio[s]) {
			s = r
		}
		if s == i {
			break
		}
		q.prio[i], q.prio[s] = q.prio[s], q.prio[i]
		q.data[i], q.data[s] = q.data[s], q.data[i]
		i = s
	}
}

type PQ[K comparable, V, P any] struct {
	data map[K]*IndexHeapElem[K, V, P]
	hp   *IndexHeap[K, V, P]
}

func NewPQ[K comparable, V, P any](n int, less collection.Less[P]) *PQ[K, V, P] {
	data := make(map[K]*IndexHeapElem[K, V, P])
	hp := NewIndexHeap[K, V, P](less)
	heap.Init(hp)
	return &PQ[K, V, P]{
		data: data,
		hp:   hp,
	}
}

func (h *PQ[K, V, P]) Len() int { return len(h.data) }

func (h *PQ[K, V, P]) Contains(k K) bool {
	_, ok := h.data[k]
	return ok
}

func (h *PQ[K, V, P]) Get(k K) (v V, p P, ok bool) {
	e, in := h.data[k]
	if in {
		v, p, ok = e.Val, e.Rank, true
	}
	return
}

func (h *PQ[K, V, P]) Priority(k K) (p P, ok bool) {
	e, ok := h.data[k]
	if ok {
		p, ok = e.Rank, true
	}
	return
}

func (h *PQ[K, V, P]) Top() (k K, v V, p P, ok bool) {
	e := h.hp.Index(0)
	if e != nil {
		k, v, p, ok = e.Key, e.Val, e.Rank, true
	}
	return
}

func (h *PQ[K, V, P]) Push(k K, v V, p P) {
	e, ok := h.data[k]
	if ok {
		e.Val = v
		e.Rank = p
		heap.Fix(h.hp, e.Idx)
		return
	}
	ele := &IndexHeapElem[K, V, P]{
		Key:  k,
		Val:  v,
		Rank: p,
	}
	h.data[k] = e
	heap.Push(h.hp, ele)
}

func (h *PQ[K, V, P]) Pop() (k K, v V, p P, ok bool) {
	n := len(h.data)
	if n == 0 {
		return
	}

	e := heap.Pop(h.hp).(*IndexHeapElem[K, V, P])
	if e != nil {
		k, v, p, ok = e.Key, e.Val, e.Rank, true
		delete(h.data, e.Key)
	}
	return
}

func (h *PQ[K, V, P]) Update(k K, v V, p P) {
	e, ok := h.data[k]
	if ok {
		e.Val = v
		e.Rank = p
		heap.Fix(h.hp, e.Idx)
	}
}

func (h *PQ[K, V, P]) Remove(k K) (v V, p P, ok bool) {
	e, ok := h.data[k]
	if ok {
		heap.Remove(h.hp, e.Idx)
		v, p, ok = e.Val, e.Rank, true
		delete(h.data, k)
	}
	return
}
