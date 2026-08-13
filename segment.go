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

type SegmentTree[N number] struct {
	ready bool
	size  int
	tree  []N
	fn    func(N, N) N
}

func NewSegmentTree[N number]() *SegmentTree[N] {
	return &SegmentTree[N]{}
}

func (t *SegmentTree[N]) Len() int {
	return t.size
}

func (t *SegmentTree[N]) F(n1, n2 N) (n N) {
	if t.ready {
		return t.fn(n1, n2)
	}
	return
}

func (t *SegmentTree[N]) build(arr []N, i int, l, r int) {
	// use arr[l:r+1] elements to construct a subtree,and root of the tree if t.tree[i]
	if l == r {
		t.tree[i] = arr[l]
	} else {
		m := (l + r) / 2
		t.build(arr, i<<1, l, m)
		t.build(arr, (i<<1)|1, m+1, r)
		t.tree[i] = t.fn(t.tree[i<<1], t.tree[(i<<1)|1])
	}
}

// parameters information about the current vertex/segment
// (i.e. the index  i and the boundaries tl and tr) and
// also the information about the boundaries of the query,l and r
func (t *SegmentTree[N]) query(i int, tl, tr int, l, r int) (n N) {
	if l > r {
		return
	}
	if tl == l && tr == r {
		return t.tree[i]
	}
	tm := (tl + tr) / 2
	a := t.query(i<<1, tl, tm, l, min(r, tm))
	b := t.query((i<<1)|1, tm+1, tr, max(l, tm+1), r)
	return t.fn(a, b)
}

func (t *SegmentTree[N]) Build(arr []N, fn func(N, N) N) {
	t.size = len(arr)
	t.fn = fn

	if len(t.tree) >= 4*t.size {
		t.tree = t.tree[:4*t.size]
	} else {
		t.tree = make([]N, 4*t.size)
	}
	t.build(arr, 1, 0, len(arr)-1)
	t.ready = true
}

func (t *SegmentTree[N]) Query(l, r int) (n N, ok bool) {
	if l < 0 || l > r || r >= t.size || !t.ready {
		return
	}
	return t.query(1, 0, t.size-1, l, r), true
}

func (t *SegmentTree[N]) update(r int, tl, tr int, i int, fn func(N) N) {
	if tl == tr {
		t.tree[r] = fn(t.tree[r])
	} else {
		tm := (tl + tr) / 2
		if i <= tm {
			t.update(r<<1, tl, tm, i, fn)
		} else {
			t.update((r<<1)|1, tm+1, tr, i, fn)
		}
		t.tree[r] = t.fn(t.tree[r<<1], t.tree[(r<<1)|1])
	}
}

func (t *SegmentTree[N]) Set(i int, n N) {
	if t.ready {
		t.update(1, 0, t.size-1, i, func(_ N) N { return n })
	}
}

func (t *SegmentTree[N]) Update(i int, fn func(N) N) {
	if t.ready {
		t.update(1, 0, t.size-1, i, fn)
	}
}

func (t *SegmentTree[N]) Get(i int) (n N, ok bool) {
	if t.ready {
		return t.Query(i, i)
	}
	return
}
