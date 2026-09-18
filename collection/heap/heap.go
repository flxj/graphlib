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

// Heap data element, rank is used to maintain heap structure.
type IndexHeapElem[K any, V any, P any] struct {
	Key  K
	Val  V
	Rank P
	Idx  int
}

type IndexHeap[K any, V any, P any] struct {
	elems []*IndexHeapElem[K, V, P]
	less  func(P, P) bool
}

func NewIndexHeap[K any, V any, P any](less func(P, P) bool) *IndexHeap[K, V, P] {
	return &IndexHeap[K, V, P]{less: less}
}

func (h *IndexHeap[K, V, P]) Len() int {
	return len(h.elems)
}

func (h *IndexHeap[K, V, P]) Less(i, j int) bool {
	return h.less(h.elems[i].Rank, h.elems[j].Rank)
}

func (h *IndexHeap[K, V, P]) Swap(i, j int) {
	h.elems[i], h.elems[j] = h.elems[j], h.elems[i]
	h.elems[i].Idx = i
	h.elems[j].Idx = j
}

func (h *IndexHeap[K, V, P]) Push(x any) {
	v, _ := x.(*IndexHeapElem[K, V, P])
	v.Idx = len(h.elems)
	h.elems = append(h.elems, v)
}

func (h *IndexHeap[K, V, P]) Pop() any {
	n := len(h.elems)
	v := h.elems[n-1]
	h.elems = h.elems[:n-1]
	return v
}

func (h *IndexHeap[K, V, P]) Top() *IndexHeapElem[K, V, P] {
	if len(h.elems) > 0 {
		return h.elems[0]
	}
	return nil
}

func (h *IndexHeap[K, V, P]) Index(i int) *IndexHeapElem[K, V, P] {
	if i >= 0 && i < len(h.elems) {
		return h.elems[i]
	}
	return nil
}
