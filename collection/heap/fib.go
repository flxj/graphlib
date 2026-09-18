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

type FibHeapElem[T any, R any] struct {
	Val  T
	Rank R
}

type fibNode[T, R any] struct {
	elem *FibHeapElem[T, R]
	// The boolean-valued attribute x.mark
	// indicates whether node x has lost a
	// child since the last time x was made
	// the child of another node.
	mark bool
	// We store the number of children in
	// the child list of node x in x.degree.
	degree int
	ptr    [4]*fibNode[T, R] // {next,prev,parent,child}
}

func fibInsert[T, R any](prev *fibNode[T, R], node *fibNode[T, R]) {
	node.ptr[0], node.ptr[2] = prev.ptr[0], prev.ptr[2]
	prev.ptr[0] = node
	node.ptr[1] = prev
	node.ptr[0].ptr[1] = node
}

func fibRemove[T, R any](node *fibNode[T, R]) {
	node.ptr[0].ptr[1] = node.ptr[1]
	node.ptr[1].ptr[0] = node.ptr[0]
	node.ptr[0], node.ptr[1] = node, node
	node.ptr[2] = nil
}

func fibMerge[T, R any](a, b *fibNode[T, R]) {
	bPrev := b.ptr[1]
	aNext := a.ptr[0]

	bPrev.ptr[0] = aNext
	a.ptr[0] = b
	aNext.ptr[1] = bPrev
	b.ptr[1] = a
	//
	par := a.ptr[2]
	for p := b; p.ptr[2] != par; p = p.ptr[0] {
		p.ptr[2] = par
	}
}

func fibLink[T, R any](y, x *fibNode[T, R]) {
	// remove y from the root list of H
	fibRemove(y)
	// make y a child of x,incrementing x.degree
	if x.ptr[3] == nil {
		y.ptr[2] = x
		x.ptr[3] = y
	} else {
		fibInsert(x.ptr[3], y)
	}
	x.degree++
	y.mark = false
}

type FibonacciHeap[T, R any] struct {
	count int
	less  func(R, R) bool
	// The pointer H.min thus points to the
	// node in the root list whose key(or priority) is minimum.
	minPtr *fibNode[T, R]
	nodes  map[*FibHeapElem[T, R]]*fibNode[T, R]
}

func NewFibonacciHeap[T, R any](less func(R, R) bool) *FibonacciHeap[T, R] {
	return &FibonacciHeap[T, R]{
		less:  less,
		nodes: make(map[*FibHeapElem[T, R]]*fibNode[T, R]),
	}
}

func (h *FibonacciHeap[T, R]) Len() int {
	return h.count
}

func (h *FibonacciHeap[T, R]) Less(a, b R) bool {
	return h.less(a, b)
}

func (h *FibonacciHeap[T, R]) Top() *FibHeapElem[T, R] {
	if h.minPtr != nil {
		return h.minPtr.elem
	}
	return nil
}

func (h *FibonacciHeap[T, R]) Push(e *FibHeapElem[T, R]) {
	if old, ok := h.nodes[e]; ok {
		old.elem.Val = e.Val
		return
	}
	n := &fibNode[T, R]{elem: e}
	n.ptr[0], n.ptr[1] = n, n
	h.nodes[e] = n
	if h.minPtr == nil {
		h.minPtr = n
	} else {
		// insert n into root list.
		fibInsert(h.minPtr, n)
		if h.less(n.elem.Rank, h.minPtr.elem.Rank) {
			h.minPtr = n
		}
	}
	h.count++
}

func (h *FibonacciHeap[T, R]) Pop() *FibHeapElem[T, R] {
	z := h.minPtr
	if z != nil {
		child := z.ptr[3]
		for p := child; p != nil; {
			// delete p from current list,and add it to root list.
			q := p.ptr[0]
			p.ptr[2] = nil
			p.ptr[0], p.ptr[1] = p, p
			fibInsert(z, p)
			if q == child {
				break
			}
			p = q
		}
		// delete z from root list.
		if z == z.ptr[0] {
			h.minPtr = nil
		} else {
			h.minPtr = z.ptr[0]
			fibRemove(z)
			h.consolidate()
		}
		h.count--
		delete(h.nodes, z.elem)
		return z.elem
	}
	return nil
}

// Consolidating the root list consists of repeatedly executing
// the following stepsuntil every root in the root list has a distinct degree value.
func (h *FibonacciHeap[T, R]) consolidate() {
	vis := make(map[*fibNode[T, R]]struct{})
	deg := make(map[int]*fibNode[T, R])
	for x := h.minPtr; x != nil; {
		q := x.ptr[0]
		d := x.degree
		for deg[d] != nil {
			//another node with the same degree as x.
			y := deg[d]
			if !h.less(x.elem.Rank, y.elem.Rank) {
				x, y = y, x
			}
			// link
			fibLink(y, x)
			delete(deg, d)
			d++
		}
		deg[d] = x
		vis[x] = struct{}{}
		// if q has already visited,stop.
		if _, ok := vis[q]; ok {
			break
		}
		x = q
	}
	// rebuild the root list.
	h.minPtr = nil
	for _, p := range deg {
		if h.minPtr == nil {
			p.ptr[2] = nil
			p.ptr[0], p.ptr[1] = p, p
			h.minPtr = p
		} else {
			fibInsert(h.minPtr, p)
			if h.less(p.elem.Rank, h.minPtr.elem.Rank) {
				h.minPtr = p
			}
		}
	}
}

func (h *FibonacciHeap[T, R]) Decrease(e *FibHeapElem[T, R], r R) bool {
	x, ok := h.nodes[e]
	if !ok {
		return false
	}
	if !h.less(r, x.elem.Rank) {
		return false
	}
	x.elem.Rank = r
	y := x.ptr[2]
	if !(y == nil || h.less(y.elem.Rank, x.elem.Rank)) {
		h.cut(x, y)
		h.cascadingCut(y)
	}
	if h.less(x.elem.Rank, h.minPtr.elem.Rank) {
		h.minPtr = x
	}
	return true
}

func (h *FibonacciHeap[T, R]) cut(x, y *fibNode[T, R]) {
	//  remove x from the child list of y, decrementing y.degree
	if x.ptr[0] == x {
		// y has only one child
		y.ptr[3] = nil
	} else {
		y.ptr[3] = x.ptr[0]
		fibRemove(x)
	}
	y.degree--
	//  add x to the root list of H
	fibInsert(h.minPtr, x)
	x.mark = false
}

func (h *FibonacciHeap[T, R]) cascadingCut(y *fibNode[T, R]) {
	z := y.ptr[2]
	if z != nil {
		if !y.mark {
			y.mark = true
		} else {
			h.cut(y, z)
			h.cascadingCut(z)
		}
	}
}

func (h *FibonacciHeap[T, R]) Delete(e *FibHeapElem[T, R]) {
	x, ok := h.nodes[e]
	if !ok {
		return
	}
	if x == h.minPtr {
		_ = h.Pop()
	} else {
		y := x.ptr[2]
		if y != nil {
			h.cut(x, y)
			h.cascadingCut(y)
		}
		// add x’s child list to the root list of H
		child := x.ptr[3]
		for p := child; p != nil; {
			q := p.ptr[0]
			fibInsert(h.minPtr, p)
			if q == child {
				break
			}
			p = q
		}
		// remove x from the root list of H
		fibRemove(x)
		h.count--
	}
}

func FibonacciHeapUnion[T, R any](h1, h2 *FibonacciHeap[T, R]) *FibonacciHeap[T, R] {
	if h1 == nil && h2 == nil {
		return nil
	} else if h1 == nil {
		return &FibonacciHeap[T, R]{minPtr: h2.minPtr, count: h2.count}
	} else if h2 == nil {
		return &FibonacciHeap[T, R]{minPtr: h1.minPtr, count: h1.count}
	}
	h := &FibonacciHeap[T, R]{}
	if h1.Len() == 0 {
		h.minPtr = h2.minPtr
		h.less = h2.less
	} else if h2.Len() == 0 {
		h.minPtr = h1.minPtr
		h.less = h1.less
	} else {
		h.minPtr = h1.minPtr
		h.less = h1.less
		fibMerge(h.minPtr, h2.minPtr)
		if !h.less(h1.minPtr.elem.Rank, h2.minPtr.elem.Rank) {
			h.minPtr = h2.minPtr
		}
	}
	h.count = h1.Len() + h2.Len()
	return h
}
