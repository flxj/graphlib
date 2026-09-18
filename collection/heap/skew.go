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

import "github.com/flxj/graphlib/collection"

type node[T, P any] struct {
	val         T
	rank        P
	left, right *node[T, P]
}

type SkewHeap[T, P any] struct {
	root *node[T, P]
	size int
	less collection.Less[P]
}

func NewSkewHeap[T, P any](less collection.Less[P]) *SkewHeap[T, P] {
	return &SkewHeap[T, P]{less: less}
}

func (h *SkewHeap[T, P]) Len() int { return h.size }

func (h *SkewHeap[T, P]) Merge(other *SkewHeap[T, P]) *SkewHeap[T, P] {
	return &SkewHeap[T, P]{
		root: mergeIter(h.root, other.root, h.less),
		size: h.size + other.size,
		less: h.less,
	}
}

func (h *SkewHeap[T, P]) Push(val T, rank P) {
	n := &node[T, P]{val: val, rank: rank}
	h.root = mergeIter(h.root, n, h.less)
	h.size++
}

func (h *SkewHeap[T, P]) Pop() (v T, p P, ok bool) {
	if h.root != nil {
		v, p, ok = h.root.val, h.root.rank, true
		h.root = mergeIter(h.root.left, h.root.right, h.less)
		h.size--
	}
	return
}

func (h *SkewHeap[T, P]) Top() (v T, p P, ok bool) {
	if h.root != nil {
		v, p, ok = h.root.val, h.root.rank, true
	}
	return
}

func mergeIter[T, P any](a, b *node[T, P], less collection.Less[P]) *node[T, P] {
	var stack []*node[T, P]

	for a != nil && b != nil {
		if less(b.rank, a.rank) {
			a, b = b, a
		}
		stack = append(stack, a)
		a = a.right
	}
	var tail *node[T, P]
	if a != nil {
		tail = a
	} else {
		tail = b
	}

	cur := tail
	for i := len(stack) - 1; i >= 0; i-- {
		n := stack[i]
		n.right = cur
		n.left, n.right = n.right, n.left
		cur = n
	}
	return cur
}
