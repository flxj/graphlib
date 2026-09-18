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

package tree

import (
	"math/rand"

	"github.com/flxj/graphlib/collection"
)

type treapNode[K any, V any] struct {
	key K
	val V
	pri int32
	siz int
	l   *treapNode[K, V]
	r   *treapNode[K, V]
}

func (t *treapNode[K, V]) size(node *treapNode[K, V]) int {
	if node == nil {
		return 0
	}
	return node.siz
}

func (t *treapNode[K, V]) resize() {
	t.siz = t.size(t.l) + t.size(t.r) + 1
}

func (t *treapNode[K, V]) left() bstNode[K, V] {
	return t.l
}

func (t *treapNode[K, V]) right() bstNode[K, V] {
	return t.r
}

func (t *treapNode[K, V]) getKey() K {
	return t.key
}

func (t *treapNode[K, V]) getVal() V {
	return t.val
}

// Treap is a Balanced Binary Search Tree, but not guaranteed to have height as O(Log n).
// The idea is to use Randomization and Binary Heap property to maintain balance with high probability.
// The expected time complexity of search, insert and delete is O(Log n).
type Treap[K any, V any] struct {
	comp collection.CompareFunc[K]
	root *treapNode[K, V]
}

// create a Non-rotating Treap.
func NewTreap[K any, V any](comp collection.CompareFunc[K]) *Treap[K, V] {
	return &Treap[K, V]{comp: comp}
}

func (t *Treap[K, V]) Len() int {
	if t.root != nil {
		return t.root.siz
	}
	return 0
}

func (t *Treap[K, V]) Compare(a, b K) int {
	return t.comp(a, b)
}

func (t *Treap[K, V]) newNode(k K, v V) *treapNode[K, V] {
	return &treapNode[K, V]{
		key: k,
		val: v,
		pri: rand.Int31(),
		siz: 1,
	}
}

func (t *Treap[K, V]) insert(root, item *treapNode[K, V]) *treapNode[K, V] {
	if root == nil {
		return item
	}
	if t.comp(root.key, item.key) == 0 {
		root.val = item.val // just update.
		return root
	} else if root.pri < item.pri {
		l, r := t.split(root, item.key)
		p := r
		for ; p != nil; p = p.l {
			if p.l == nil {
				break
			}
		}
		if p == nil || t.comp(p.key, item.key) != 0 {
			item.l, item.r = l, r
			item.resize()
			return item
		}
		if t.comp(p.key, item.key) == 0 {
			p.val = item.val
		}
		p = t.merge(l, r)
		return p
	} else {
		if t.comp(root.key, item.key) < 0 {
			root.r = t.insert(root.r, item)
		} else {
			root.l = t.insert(root.l, item)
		}
		root.resize()
		return root
	}
}

// Insert a pair of key value data, update in place if the key already exists.
func (t *Treap[K, V]) Insert(k K, v V) {
	t.root = t.insert(t.root, t.newNode(k, v))
}

// Query the value data corresponding to the key. If the key does not exist,
// return a null value and a false flag.
func (t *Treap[K, V]) Search(k K) (v V, ok bool) {
	for p := t.root; p != nil; {
		if t.comp(p.key, k) == 0 {
			return p.val, true
		} else if t.comp(p.key, k) <= 0 {
			p = p.r
		} else {
			p = p.l
		}
	}
	return
}

// Return the current Treap root element.
func (t *Treap[K, V]) Root() (k K, v V, ok bool) {
	if t.root != nil {
		k, v, ok = t.root.key, t.root.val, true
	}
	return
}

// Query the ranking of the key.
func (t *Treap[K, V]) Rank(k K) (int, bool) {
	var n int
	for p := t.root; p != nil; {
		if t.comp(p.key, k) == 0 {
			n = n + p.size(p.l) + 1
			return n, true
		} else if t.comp(p.key, k) > 0 {
			p = p.l
		} else {
			n = n + p.size(p.l) + 1
			p = p.r
		}
	}
	return 0, false
}

// Query data based on ranking.
func (t *Treap[K, V]) Nth(n int) (k K, v V, ok bool) {
	if n <= 0 || n > t.Len() {
		return
	}
	for p := t.root; p != nil; {
		if p.l != nil {
			if p.l.siz >= n {
				p = p.l
				continue
			} else {
				n -= p.l.siz
			}
		}
		if n == 1 {
			return p.key, p.val, true
		}
		n--
		p = p.r
	}
	return
}

func (t *Treap[K, V]) erase(cur *treapNode[K, V], k K) (p *treapNode[K, V], v V, ok bool) {
	if cur == nil {
		return
	}
	if t.comp(cur.key, k) == 0 {
		p, v, ok = t.merge(cur.l, cur.r), cur.val, true
		return
	} else if t.comp(cur.key, k) > 0 {
		cur.l, v, ok = t.erase(cur.l, k)
	} else {
		cur.r, v, ok = t.erase(cur.r, k)
	}
	if ok {
		cur.resize()
	}
	return cur, v, ok
}

// Delete key. If the key exists, delete it and return its value. If the key does not exist, return a false flag.
func (t *Treap[K, V]) Delete(k K) (v V, ok bool) {
	t.root, v, ok = t.erase(t.root, k)
	return
}

// Clear the current Treap.
func (t *Treap[K, V]) Clean() {
	t.root = nil
}

// Query the minimum element of the key.
func (t *Treap[K, V]) Min() (k K, v V, ok bool) {
	for p := t.root; p != nil; p = p.l {
		if p.l == nil {
			return p.key, p.val, true
		}
	}
	return
}

// Query the maximum element of the key.
func (t *Treap[K, V]) Max() (k K, v V, ok bool) {
	for p := t.root; p != nil; p = p.r {
		if p.r == nil {
			return p.key, p.val, true
		}
	}
	return
}

func (t *Treap[K, V]) Cursor() collection.Cursor[K, V] {
	return newBSTCursor[K, V](t.root, t.comp)
}

// separates tree root in 2 subtrees l and r.
// so that l contains all elements with key < k, and r contains all elements with key>=k.
func (t *Treap[K, V]) split(root *treapNode[K, V], k K) (l, r *treapNode[K, V]) {
	if root == nil {
		return nil, nil
	} else if t.comp(root.key, k) < 0 {
		l, r := t.split(root.r, k)
		root.r = l
		root.resize()
		return root, r
	} else if t.comp(root.key, k) == 0 {
		l := root.l
		root.l = nil
		root.resize()
		return l, root
	} else {
		l, r := t.split(root.l, k)
		root.l = r
		root.resize()
		return l, root
	}
}

func (t *Treap[K, V]) splitByRank(root *treapNode[K, V], rk int) (l, m, r *treapNode[K, V]) {
	return nil, nil, nil
}

// combines two subtrees t1 t2,and returns the new tree.
// It works under the assumption that are ordered (all keys in t1 are smaller than keys in t2)
func (t *Treap[K, V]) merge(t1, t2 *treapNode[K, V]) *treapNode[K, V] {
	if t1 == nil || t2 == nil {
		if t1 == nil {
			return t2
		} else {
			return t1
		}
	}
	if t1.pri > t2.pri {
		t1.r = t.merge(t1.r, t2)
		t1.resize()
		return t1
	} else {
		t2.l = t.merge(t1, t2.l)
		t2.resize()
		return t2
	}
}

func TreapBuild[K any, V any](keys []K, vals []V, comp collection.CompareFunc[K]) *Treap[K, V] {
	// check sorted
	return nil
}

func TreapSplit[K any, V any](t *Treap[K, V], k K) (*Treap[K, V], *Treap[K, V]) {
	if t == nil {
		return nil, nil
	}
	l, r := t.split(t.root, k)
	tl := &Treap[K, V]{root: l}
	tr := &Treap[K, V]{root: r}
	return tl, tr
}

func TreapMerge[K any, V any](t1, t2 *Treap[K, V], k K) *Treap[K, V] {
	if t1 == nil {
		return t2
	}
	if t2 == nil {
		return t1
	}
	r := t1.merge(t1.root, t2.root)
	return &Treap[K, V]{root: r}
}
