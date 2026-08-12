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
	"math/rand"
)

type stNode[K any, V any] struct {
	key    K
	val    V
	parent *stNode[K, V]
	left   *stNode[K, V]
	right  *stNode[K, V]
}

func (n *stNode[K, V]) isRoot() bool {
	return n.parent == nil
}

func (n *stNode[K, V]) isLeftChild() bool {
	return n.parent != nil && n == n.parent.left
}

type SplayTree[K any, V any] struct {
	comp  CompareFunc[K]
	count int
	root  *stNode[K, V]
}

func NewSplayTree[K any, V any](comp CompareFunc[K]) *SplayTree[K, V] {
	return &SplayTree[K, V]{comp: comp}
}

// splay a node when that node is not the root and we wish to transport it to the root.
func (s *SplayTree[K, V]) splay(x *stNode[K, V]) {
	if x == nil {
		return
	}
	rightRotate := func(x, p, gp *stNode[K, V]) {
		x.parent = gp
		p.parent = x
		p.left = x.right
		if p.left != nil {
			p.left.parent = p
		}
		x.right = p
		if gp != nil {
			if gp.left == p {
				gp.left = x
			} else {
				gp.right = x
			}
		}
	}
	leftRotate := func(x, p, gp *stNode[K, V]) {
		x.parent = gp
		p.parent = x
		p.right = x.left
		if p.right != nil {
			p.right.parent = p
		}
		x.left = p
		if gp != nil {
			if gp.left == p {
				gp.left = x
			} else {
				gp.right = x
			}
		}
	}
	for !x.isRoot() {
		p := x.parent
		if p.isRoot() { // x the root’s child, Zig.
			if x.isLeftChild() {
				rightRotate(x, p, nil)
			} else {
				leftRotate(x, p, nil)
			}
		} else {
			gp := p.parent
			if x.isLeftChild() {
				if p.isLeftChild() { // x a left-left child,Zig-Zig
					rightRotate(p, gp, gp.parent)
					rightRotate(x, p, p.parent)
				} else { // x is a right-left child,Zig-Zag
					rightRotate(x, p, gp)
					leftRotate(x, gp, gp.parent)
				}
			} else {
				if p.isLeftChild() { // x is a left-right child,Zig-Zag
					leftRotate(x, p, gp)
					rightRotate(x, gp, gp.parent)
				} else { // x is a right-right child,Zig-Zig
					leftRotate(p, gp, gp.parent)
					leftRotate(x, p, p.parent)
				}
			}
		}
	}
}

func (s *SplayTree[K, V]) searchSubtree(root *stNode[K, V], key K) (x *stNode[K, V], ok bool) {
	x = root
	for x != nil {
		c := s.comp(key, x.key)
		if c == 0 {
			ok = true
			break
		} else if c > 0 {
			if x.right == nil {
				break
			}
			x = x.right
		} else {
			if x.left == nil {
				break
			}
			x = x.left
		}
	}
	s.splay(x)
	return
}

func (s *SplayTree[K, V]) Search(key K) (v V, ok bool) {
	root, ok := s.searchSubtree(s.root, key)
	s.root = root
	if ok {
		return root.val, true
	}
	return
}

func (s *SplayTree[K, V]) Insert(key K, val V) {
	root, ok := s.searchSubtree(s.root, key)
	s.root = root
	if ok {
		s.root.val = val
		return
	}
	node := &stNode[K, V]{key: key, val: val}
	if s.root != nil {
		s.root.parent = node
		if s.comp(s.root.key, key) > 0 {
			node.left = s.root.left
			if node.left != nil {
				node.left.parent = node
			}
			s.root.left = nil
			node.right = s.root
		} else {
			node.right = s.root.right
			if node.right != nil {
				node.right.parent = node
			}
			s.root.right = nil
			node.left = s.root
		}
	}
	s.root = node
	s.count++
}

func (s *SplayTree[K, V]) Delete(key K) (V, bool) {
	v, ok := s.Search(key)
	if !ok {
		return v, false
	}
	if s.root == nil {
		return v, false
	}
	if s.root.left == nil {
		r := s.root.right
		s.root.right = nil
		if r != nil {
			r.parent = nil
		}
		s.root = r
	} else if s.root.right == nil {
		l := s.root.left
		s.root.left = nil
		if l != nil {
			l.parent = nil
		}
		s.root = l
	} else {
		/*
			If neither L nor R is empty then we call splay on key but only in the subtree R
			Since key is not there (because it’s the parent of R) and because R > key the
			result will be that the new root of R, call it r, will be the inorder successor of key.
			Consequently r will have no left subtree itself (because there is nothing greater
			than key and smaller than r) but it will have a (possibly empty) right subtree.
			We simply delete s.root and shift r up to its place.
		*/
		s.root.right.parent = nil // cut right subtree from s.root
		r, _ := s.searchSubtree(s.root.right, key)
		s.root.left.parent = r
		r.left = s.root.left
		s.root.left, s.root.right = nil, nil
		s.root = r
	}
	s.count--
	return v, ok
}

func (s *SplayTree[K, V]) Compare(a, b K) int {
	return s.comp(a, b)
}

func (s *SplayTree[K, V]) Len() int {
	return s.count
}

func (s *SplayTree[K, V]) Min() (k K, v V, ok bool) {
	p := s.root
	for p != nil {
		if p.left == nil {
			break
		}
		p = p.left
	}
	if p != nil {
		k, v, ok = p.key, p.val, true
		s.splay(p)
	}
	return
}

func (s *SplayTree[K, V]) Max() (k K, v V, ok bool) {
	p := s.root
	for p != nil {
		if p.right == nil {
			break
		}
		p = p.right
	}
	if p != nil {
		k, v, ok = p.key, p.val, true
		s.splay(p)
	}
	return
}

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

// Non-rotating Treap
type Treap[K any, V any] struct {
	comp CompareFunc[K]
	root *treapNode[K, V]
}

func NewTreap[K any, V any](comp CompareFunc[K]) *Treap[K, V] {
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

// Delete Element.
func (t *Treap[K, V]) Delete(k K) (v V, ok bool) {
	t.root, v, ok = t.erase(t.root, k)
	return
}

// Clear the current Treap.
func (t *Treap[K, V]) Clean() {
	t.root = nil
}

func (t *Treap[K, V]) Min() (k K, v V, ok bool) {
	for p := t.root; p != nil; p = p.l {
		if p.l == nil {
			return p.key, p.val, true
		}
	}
	return
}

func (t *Treap[K, V]) Max() (k K, v V, ok bool) {
	for p := t.root; p != nil; p = p.r {
		if p.r == nil {
			return p.key, p.val, true
		}
	}
	return
}

func (t *Treap[K, V]) Cursor() Cursor[K, V] {
	return &treapCursor[K, V]{
		comp: t.comp,
		root: t.root,
		stk:  newStack[*treapPath[K, V]](),
	}
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

func TreapBuild[K any, V any](keys []K, vals []V, comp CompareFunc[K]) *Treap[K, V] {
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

type treapPath[K any, V any] struct {
	node    *treapNode[K, V]
	visited bool
}

type treapCursor[K any, V any] struct {
	comp CompareFunc[K]
	root *treapNode[K, V]
	prev *treapNode[K, V]
	stk  *stack[*treapPath[K, V]]
}

func (c *treapCursor[K, V]) Open() error { return nil }
func (c *treapCursor[K, V]) Close()      {}

func (c *treapCursor[K, V]) reset() {
	c.stk.clean()
	c.prev = nil
}

// The Seek(key) method locates the cursor at the key.
// If the key does not exist, it locates at the next key and returns it
func (c *treapCursor[K, V]) Seek(key K) (k K, v V, ok bool) {
	c.reset()
	for p := c.root; p != nil; {
		c.stk.push(&treapPath[K, V]{node: p})
		if c.comp(p.key, key) == 0 {
			return p.key, p.val, true
		} else if c.comp(p.key, key) > 0 {
			p = p.l
		} else {
			p = p.r
		}
	}
	if tp := c.stk.top(); tp != nil {
		tp.visited = true
		k, v = tp.node.key, tp.node.val
	}
	return
}

// The First method locates the cursor at the minimum element of the set.
// If there is no minimum element (the set is empty), it returns false
func (c *treapCursor[K, V]) First() (k K, v V, ok bool) {
	c.reset()
	for p := c.root; p != nil; p = p.l {
		c.stk.push(&treapPath[K, V]{node: p})
	}
	if tp := c.stk.top(); tp != nil {
		tp.visited = true
		k, v, ok = tp.node.key, tp.node.val, true
	}
	return
}

// The Last method locates the cursor at the maximum element of the set.
// If there is no maximum element (the set is empty), it returns false.
func (c *treapCursor[K, V]) Last() (k K, v V, ok bool) {
	c.reset()
	for p := c.root; p != nil; p = p.r {
		c.stk.push(&treapPath[K, V]{node: p})
	}
	if tp := c.stk.top(); tp != nil {
		tp.visited = true
		k, v, ok = tp.node.key, tp.node.val, true
	}
	return
}

// HasNext returns whether the next element exists relative to the current cursor position.
func (c *treapCursor[K, V]) HasNext() bool {
	var prev *treapNode[K, V]
	for !c.stk.empty() {
		p := c.stk.top()
		if prev != nil && c.comp(prev.key, p.node.key) >= 0 {
			pp, _ := c.stk.pop()
			prev = pp.node
			continue
		}
		// try to move cursor to next element,but not visited it.
		if !p.visited {
			return true
		}
		// subtree not nil
		if p.node.r != nil {
			if prev == nil || c.comp(prev.key, p.node.key) < 0 { // p's right subtree not visited.
				for q := p.node.r; q != nil; q = q.l {
					c.stk.push(&treapPath[K, V]{node: q})
				}
				return true
			}
		}
		pp, _ := c.stk.pop()
		prev = pp.node
	}
	return false
}

// Next() moves the cursor backwards and returns the element.
// If the element does not exist, it returns a type zero value.
func (c *treapCursor[K, V]) Next() (k K, v V) {
	var prev *treapNode[K, V]
	for !c.stk.empty() {
		p := c.stk.top()
		//
		if prev != nil && c.comp(prev.key, p.node.key) >= 0 {
			pp, _ := c.stk.pop()
			prev = pp.node
			continue
		}
		// current node not visited,so return it.
		if !p.visited {
			p.visited = true
			return p.node.key, p.node.val
		}
		// subtree not nil
		if p.node.r != nil {
			if prev == nil || c.comp(prev.key, p.node.key) < 0 { // p's right subtree not visited.
				for q := p.node.r; q != nil; q = q.l {
					c.stk.push(&treapPath[K, V]{node: q})
				}
				pp := c.stk.top()
				pp.visited = true
				return pp.node.key, pp.node.val
			}
		}
		pp, _ := c.stk.pop()
		prev = pp.node
	}
	return
}

// HasPrev() returns whether the previous element exists relative to the current cursor position.
func (c *treapCursor[K, V]) HasPrev() bool {
	return false
}

// Prev() moves the cursor forward and returns the element.
// If the element does not exist, it returns a type value of zero.
func (c *treapCursor[K, V]) Prev() (k K, v V) {
	return
}

var (
	DefaultScapegoatTreeAlpha = 0.75
)

type sgtNode[K any, V any] struct {
	key K
	val V
	del bool
	siz int
	act int
	l   *sgtNode[K, V]
	r   *sgtNode[K, V]
}

func (s *sgtNode[K, V]) size(n *sgtNode[K, V]) int {
	if n != nil {
		return n.siz
	}
	return 0
}

func (s *sgtNode[K, V]) active(n *sgtNode[K, V]) int {
	if n != nil {
		return n.act
	}
	return 0
}

func (s *sgtNode[K, V]) resize() {
	s.siz = s.size(s.l) + s.size(s.r) + 1
	s.act = s.active(s.l) + s.active(s.r)
	if !s.del {
		s.act++
	}
}

func (s *sgtNode[K, V]) balance(n, m int) bool {
	if s.size(s.l) > s.siz*n/m || s.size(s.r) > s.siz*n/m {
		return false
	}
	return true
}

type ScapegoatTree[K any, V any] struct {
	comp  CompareFunc[K]
	num   int
	alpha float64
	root  *sgtNode[K, V]
}

func NewScapegoatTree[K any, V any](alpha float64, comp CompareFunc[K]) *ScapegoatTree[K, V] {
	if alpha < 0.0 || alpha >= 1.0 {
		return nil
	}
	return &ScapegoatTree[K, V]{
		comp:  comp,
		alpha: alpha,
		num:   int(alpha * 1000),
	}
}

func (s *ScapegoatTree[K, V]) Len() int {
	if s.root == nil {
		return 0
	}
	return s.root.act
}

func (s *ScapegoatTree[K, V]) newNode(k K, v V) *sgtNode[K, V] {
	return &sgtNode[K, V]{
		key: k,
		val: v,
		siz: 1,
		act: 1,
	}
}

func (s *ScapegoatTree[K, V]) insert(cur *sgtNode[K, V], k K, v V) (*sgtNode[K, V], *sgtNode[K, V]) {
	if cur == nil {
		return s.newNode(k, v), nil
	}
	var sg *sgtNode[K, V]
	if s.comp(cur.key, k) == 0 {
		cur.val = v
		if cur.del {
			cur.del = false
		}
	} else if s.comp(cur.key, k) > 0 {
		cur.l, sg = s.insert(cur.l, k, v)
	} else {
		cur.r, sg = s.insert(cur.r, k, v)
	}
	cur.resize()
	if !cur.balance(s.num, 1000) {
		sg = cur
	}
	return cur, sg
}

func (s *ScapegoatTree[K, V]) rebalance(cur, sg *sgtNode[K, V], k K) *sgtNode[K, V] {
	if cur == sg {
		ks, vs := s.flatten(sg)
		root := s.build(ks, vs)
		return root
	} else if s.comp(cur.key, k) > 0 {
		cur.l = s.rebalance(cur.l, sg, k)
	} else {
		cur.r = s.rebalance(cur.r, sg, k)
	}
	cur.resize()
	return cur
}

func (s *ScapegoatTree[K, V]) Insert(k K, v V) {
	var sg *sgtNode[K, V]
	s.root, sg = s.insert(s.root, k, v)
	if sg != nil {
		if sg == s.root {
			// rebuild s.root
			ks, vs := s.flatten(s.root)
			s.root = s.build(ks, vs)
		} else {
			// rebuild sgt
			s.root = s.rebalance(s.root, sg, k)
		}
	}
}

func (s *ScapegoatTree[K, V]) Search(k K) (v V, ok bool) {
	for p := s.root; p != nil; {
		if s.comp(p.key, k) == 0 {
			if p.del {
				return
			}
			return p.val, true
		} else if s.comp(p.key, k) > 0 {
			p = p.l
		} else {
			p = p.r
		}
	}
	return
}

func (s *ScapegoatTree[K, V]) del(node *sgtNode[K, V], k K) (v V, ok bool) {
	if node == nil {
		return
	}
	if s.comp(node.key, k) == 0 {
		if !node.del {
			node.del = true
			node.act--
			return node.val, true
		}
	} else if s.comp(node.key, k) > 0 {
		v, ok = s.del(node.l, k)
	} else {
		v, ok = s.del(node.r, k)
	}
	if ok {
		node.resize()
	}
	return
}

func (s *ScapegoatTree[K, V]) Delete(k K) (v V, ok bool) {
	v, ok = s.del(s.root, k)
	if s.root != nil && s.root.act < s.root.siz*s.num/1000 {
		ks, vs := s.flatten(s.root)
		s.root = s.build(ks, vs)
	}
	return
}

func (s *ScapegoatTree[K, V]) min(node *sgtNode[K, V]) (k K, v V, ok bool) {
	if node == nil {
		return
	}
	if node.active(node.l) > 0 {
		return s.min(node.l)
	}
	if !node.del {
		return node.key, node.val, true
	}
	return s.min(node.r)
}

func (s *ScapegoatTree[K, V]) Min() (k K, v V, ok bool) {
	return s.min(s.root)
}

func (s *ScapegoatTree[K, V]) max(node *sgtNode[K, V]) (k K, v V, ok bool) {
	if node == nil {
		return
	}
	if node.active(node.r) > 0 {
		return s.max(node.r)
	}
	if !node.del {
		return node.key, node.val, true
	}
	return s.max(node.l)
}

func (s *ScapegoatTree[K, V]) Max() (k K, v V, ok bool) {
	return s.max(s.root)
}

func (s *ScapegoatTree[K, V]) flatten(r *sgtNode[K, V]) ([]K, []V) {
	if r == nil {
		return nil, nil
	}
	var keys []K
	var vals []V
	stk := newStack[*sgtNode[K, V]]()
	p := r
	for !stk.empty() || p != nil {
		for p != nil {
			stk.push(p)
			p = p.l
		}
		p, _ = stk.pop()
		if !p.del {
			keys = append(keys, p.key)
			vals = append(vals, p.val)
			//fmt.Println("add key=", p.key)
		}
		p = p.r
	}
	return keys, vals
}

func (s *ScapegoatTree[K, V]) build(keys []K, vals []V) *sgtNode[K, V] {
	switch n := len(keys); n {
	case 0:
		return nil
	case 1:
		return &sgtNode[K, V]{
			key: keys[0],
			val: vals[0],
			siz: 1,
			act: 1,
		}
	default:
		root := &sgtNode[K, V]{
			key: keys[n/2],
			val: vals[n/2],
		}
		root.l = s.build(keys[:n/2], vals[:n/2])
		root.r = s.build(keys[n/2+1:], vals[n/2+1:])
		root.resize()
		return root
	}
}

type rbNode[K any, V any] struct {
	red bool
	key K
	val V
	p   *rbNode[K, V]
	l   *rbNode[K, V]
	r   *rbNode[K, V]
}

// 1. The root of the tree is colored Black.
// 2. A Red node can have only Black children.
// 3. Every path from the root to a leaf contains the same number of Black nodes。
type RedBlackTree[K any, V any] struct {
	comp  CompareFunc[K]
	count int
	root  *rbNode[K, V]
}

func NewRedBlackTree[K any, V any](comp CompareFunc[K]) *RedBlackTree[K, V] {
	return &RedBlackTree[K, V]{comp: comp}
}

func (t *RedBlackTree[K, V]) Len() int {
	return t.count
}

func (t *RedBlackTree[K, V]) rebalance(x *rbNode[K, V]) {
	if x.p == nil { // x is root
		x.red = false
		t.root = x
		return
	}
	y := x.p
	// 1.Parent is Black:
	// In this case we can safely color the problematic node Red.
	// It is easy to see that the invariant is not violated. X cannot be the root, so Red is a valid color.
	if !y.red {
		x.red = true
		return
	} else {
		// 2. Parent is Red:
		// This case is little more involved, because we need to know the color of z, x’s uncle, or y’s sibling
		// (y is x’s parent). p cannot be the root of the tree, because it is Red, soCexists.
		// If not recall that NIL pointers are taken as Blacks. Again, there are two subcases.
		var gl bool
		var z *rbNode[K, V]
		g := y.p
		if y == g.l {
			z = g.r
			gl = true
		} else {
			z = g.l
		}
		// (a) Uncle is Red
		if z != nil && z.red {
			// In this case, we color X Red.before doing so, we have to make sure that its parent has become Black.
			y.red, z.red = false, false
			x.red = true
			t.rebalance(g)
		} else {
			// (b) Uncle is Black
			if gl && y.l == x {
				// x is g's left-left child, right rotate
				g.l = y.r
				if y.r != nil {
					y.r.p = g
				}
				y.p = g.p
				if g.p != nil {
					if g == g.p.l {
						g.p.l = y
					} else {
						g.p.r = y
					}
				}
				y.r = g
				g.p = y
				x.red, g.red = true, true
				t.rebalance(y)
			} else if !gl && y.r == x {
				// x is g's right-right child, left rotate
				g.r = y.l
				if y.l != nil {
					y.l.p = g
				}
				y.p = g.p
				if g.p != nil {
					if g == g.p.l {
						g.p.l = y
					} else {
						g.p.r = y
					}
				}
				y.l = g
				g.p = y
				x.red, g.red = true, true
				t.rebalance(y)
			} else if gl && y.r == x {
				// x is g's left-right child, left+right rotate
				y.r = x.l
				if x.l != nil {
					x.l.p = y
				}
				x.l = y
				y.p = x

				g.l = x.r
				if x.r != nil {
					x.r.p = g
				}
				x.p = g.p
				if g.p != nil {
					if g == g.p.l {
						g.p.l = x
					} else {
						g.p.r = x
					}
				}
				x.r = g
				g.p = x
				y.red, g.red = true, true
				t.rebalance(x)
			} else {
				y.l = x.r
				if x.r != nil {
					x.r.p = y
				}
				x.r = y
				y.p = x

				g.r = x.l
				if x.l != nil {
					x.l.p = g
				}
				x.p = g.p
				if g.p != nil {
					if g == g.p.l {
						g.p.l = x
					} else {
						g.p.r = x
					}
				}
				x.l = g
				g.p = x
				y.red, g.red = true, true
				t.rebalance(x)
			}
		}
	}
}

func (t *RedBlackTree[K, V]) Insert(k K, v V) {
	var p *rbNode[K, V]
	for q := t.root; q != nil; {
		if t.comp(q.key, k) == 0 {
			q.val = v
			return
		} else if t.comp(q.key, k) > 0 {
			if q.l == nil {
				// insert a new node as p's left child
				p = &rbNode[K, V]{red: true, key: k, val: v, p: q}
				q.l = p
				break
			} else {
				q = q.l
			}
		} else {
			if q.r == nil {
				p = &rbNode[K, V]{red: true, key: k, val: v, p: q}
				q.r = p
				break
			} else {
				q = q.r
			}
		}
	}
	t.count++
	if p == nil {
		t.root = &rbNode[K, V]{key: k, val: v, p: p}
		return
	}
	t.rebalance(p)
}

func (t *RedBlackTree[K, V]) Search(k K) (v V, ok bool) {
	for p := t.root; p != nil; {
		if t.comp(p.key, k) == 0 {
			return p.val, true
		} else if t.comp(p.key, k) < 0 {
			p = p.r
		} else {
			p = p.l
		}
	}
	return
}

func (t *RedBlackTree[K, V]) Min() (k K, v V, ok bool) {
	for p := t.root; p != nil; p = p.l {
		if p.l == nil {
			return p.key, p.val, true
		}
	}
	return
}

func (t *RedBlackTree[K, V]) Max() (k K, v V, ok bool) {
	for p := t.root; p != nil; p = p.r {
		if p.r == nil {
			return p.key, p.val, true
		}
	}
	return
}

func (t *RedBlackTree[K, V]) Delete(k K) (v V, ok bool) {
	var p *rbNode[K, V]
	for q := t.root; q != nil; {
		if t.comp(q.key, k) == 0 {
			v, ok = q.val, true
			for p = q.r; p != nil && p.l != nil; p = p.l {
			}
			if p != nil {
				q.key, q.val = p.key, p.val
			} else {
				for p = q.l; p != nil && p.r != nil; p = p.r {
				}
				if p != nil {
					q.key, q.val = p.key, p.val
				} else {
					p = q
				}
			}
			break
		} else if t.comp(q.key, k) < 0 {
			q = q.r
		} else {
			q = q.r
		}
	}
	if p != nil {
		t.count--
		t.del(p)
	}
	return
}

func (t *RedBlackTree[K, V]) del(x *rbNode[K, V]) {
	z := x.l
	if x.r != nil {
		z = x.r
	}
	// If X is Red, then we just remove it. The invariant is not affected,
	// since X cannot be the root, and its removal does not change the number of Black nodes in any path.
	if x.red {
		if x == x.p.l {
			x.p.l = z
		} else {
			x.p.r = z
		}
		if z != nil {
			z.p = x.p
		}
		x.p, x.l, x.r = nil, nil, nil
	} else {
		// If X is Black and its single child is Red, then we can safely remove X and color the child Black.
		if z != nil && z.red {
			if x.p == nil { // x is root
				z.p = nil
				z.red = false
				t.root = z
				return
			} else {
				if x == x.p.l {
					x.p.l = z
				} else {
					x.p.r = z
				}
				z.p = x.p
				z.red = false
				x.p, x.l, x.r = nil, nil, nil
				return
			}
		}
		// So, the only case that is really the problem is when the node to be deleted is Black and it has
		// a Black child or no children at all (remember that NIL pointers are treated as Black).
		y := x.p
		if y == nil { // x is root
			if z != nil {
				z.p = nil
				z.red = false
			}
			t.root = z
			return
		}
		// 1. Parent is Red (and Sibling is Black)
		// If the parent of X is Red, then X must have a sibling, otherwise the invariant would be violated
		// (the path through X wouldhaveat leastonemoreBlack node, that is X, than the path that
		// would terminate at the NIL pointer of X’s parent).
		// Moreover, the sibling has to be Black,given that it is a son of a Red node.
		if y.red {
			s := y.l
			sl := true
			if x == y.l {
				s = y.r
				sl = false
			}
			// (a) Close Nephew is Black
			// The idea is to pull the Black sibling s to the top to compensate for the loss of the Black X on the path (y −x −T1).
			// This is done by a single rotation.
			if !sl && (s.l == nil || !s.l.red) {
				y.l = z
				if z != nil {
					z.p = y
				}
				x.p, x.l, x.r = nil, nil, nil
				y.r = s.l
				if s.l != nil {
					s.l.p = y
				}
				s.l = y
				s.p = y.p
				if y.p.l == y { // y is red so its not root, y.p != nil
					y.p.l = s
				} else {
					y.p.r = s
				}
				y.p = s
				return
			}
			if sl && (s.r == nil || !s.r.red) {
				y.r = z
				if z != nil {
					z.p = y
				}
				x.p, x.l, x.r = nil, nil, nil
				y.l = s.r
				if s.r != nil {
					s.r.p = y
				}
				s.r = y
				s.p = y.p
				if y.p.l == y {
					y.p.l = s
				} else {
					y.p.r = s
				}
				y.p = s
				return
			}
			// (b) Close Nephew is Red
			// The problem is fully restored by a double rotation and some recoloring.
			if !sl && s.l != nil && s.l.red {
				y.l = z
				if z != nil {
					z.p = y
				}
				w := s.l
				s.l = w.r
				if w.r != nil {
					w.r.p = s
				}
				w.r = s
				s.p = w
				y.r = w.l
				if w.l != nil {
					w.l.p = y
				}
				w.l = y
				w.p = y.p
				if y.p != nil {
					if y == y.p.l {
						y.p.l = w
					} else {
						y.p.r = w
					}
				}
				y.p = w
				y.red = false
				return
			}
			if sl && s.r != nil && s.r.red {
				y.r = z
				if z != nil {
					z.p = y
				}
				w := s.r
				s.r = w.l
				if w.l != nil {
					w.l.p = s
				}
				w.l = s
				s.p = w
				y.l = w.r
				w.r = y
				w.p = y.p
				if y.p != nil {
					if y == y.p.l {
						y.p.l = w
					} else {
						y.p.r = w
					}
				}
				y.p = w
				y.red = false
				return
			}
		} else {
			// 2.Parent is Black
			// X must have a sibling, as in the previous case, otherwise the invariant would be violated (the
			// path through X wouldhaveat leastonemoreBlack node, that is X, than the path that would
			// terminate at the NIL pointer of X’s parent). However, we cannot determine the color of X’s sibling
			// given this situation. So, we distinguish among two cases:
			s := y.l
			sl := true
			if x == y.l {
				s = y.r
				sl = false
			}
			// (a) Sibling is Red
			// The problem is attacked by a single rotation that,interestingly,
			// pushes the problem down! However, notice that the resulting situation is the
			// one where the parent is Red and the sibling is Black, which we just described and can be
			// fully resolved locally.
			if s.red {
				if sl {
					y.l = s.r
					if s.r != nil {
						s.r.p = y
					}
					s.r = y
					s.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = s
						} else {
							y.p.r = s
						}
					} else {
						t.root = s
					}
					y.p = s
					y.red, s.red = true, false
					t.del(x)
				} else {
					y.r = s.l
					if s.l != nil {
						s.l.p = y
					}
					s.l = y
					s.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = s
						} else {
							y.p.r = s
						}
					} else {
						t.root = s
					}
					y.p = s
					y.red, s.red = true, false
					t.del(x)
				}
			} else {
				// (b) Sibling is Black
				// There are several cases, depending on the coloring of the nephews of X.
				// Notice that the nephews do not necessarily exist. However, the case where some
				// nephew does not exist is equivalent to that where this nephew is colored Black.

				// b.i. Far Nephew is Red
				// In this case, the problem is resolved with a single rotation
				if !sl && s.r != nil && s.r.red {
					y.l = z
					if z != nil {
						z.p = y
					}
					x.p, x.l, x.r = nil, nil, nil
					y.r = s.l
					if s.l != nil {
						s.l.p = y
					}
					s.l = y
					s.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = s
						} else {
							y.p.r = s
						}
					} else {
						// y is root
						t.root = s
					}
					s.r.red = false
					return
				}
				if sl && s.l != nil && s.l.red {
					y.r = z
					if z != nil {
						z.p = y
					}
					x.p, x.l, x.r = nil, nil, nil
					y.l = s.r
					if s.r != nil {
						s.r.p = y
					}
					s.r = y
					s.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = s
						} else {
							y.p.r = s
						}
					} else {
						t.root = s
					}
					s.l.red = false
					return
				}
				// b.ii. Far Nephew is Black
				// This case requires checking the color of the other nephew.
				// Depending on its color there are two cases:

				// b.ii.A Close Nephew is Red
				if !sl && s.l != nil && s.l.red {
					y.l = z
					if z != nil {
						z.p = y
					}
					w := s.l
					s.l = w.r
					if w.r != nil {
						w.r.p = s
					}
					w.r = s
					s.p = w
					y.r = w.l
					if w.l != nil {
						w.l.p = y
					}
					w.l = y
					w.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = w
						} else {
							y.p.r = w
						}
					} else {
						t.root = w
					}
					y.p = w
					w.red = false
					return
				}
				if sl && s.r != nil && s.r.red {
					y.r = z
					if z != nil {
						z.p = y
					}
					w := s.r
					s.r = w.l
					if w.l != nil {
						w.l.p = s
					}
					w.l = s
					s.p = w
					y.l = w.r
					if w.r != nil {
						w.r.p = y
					}
					w.r = y
					w.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = w
						} else {
							y.p.r = w
						}
					} else {
						t.root = w
					}
					y.p = w
					w.red = false
					return
				}
				// b.ii.B. Close Nephew is Black
				// This last case is tricky, because the existence of only Black nodes in the neighborhood,
				// does not provide for removal of X and local recoloring and/or restructuring. The only
				// thing we can do in case, is to push the problem up in the tree.
				if !sl && (s.l == nil || !s.l.red) {
					s.red = true
					y.l = z
					if z != nil {
						z.p = y
					}
					x.p, x.l, x.r = nil, nil, nil
					if y.p != nil {
						x.p = y.p
						if y == y.p.l {
							y.p.l = x
						} else {
							y.p.r = x
						}
						y.p = x
						x.r = y
						t.del(x)
					} else {
						t.root = y
						return
					}
				}
				if sl && (s.r == nil || !s.r.red) {
					s.red = true
					y.r = z
					if z != nil {
						z.p = y
					}
					x.p, x.l, x.r = nil, nil, nil
					if y.p != nil {
						x.p = y.p
						if y == y.p.l {
							y.p.l = x
						} else {
							y.p.r = x
						}
						y.p = x
						x.l = y
						t.del(x)
					} else {
						t.root = y
						return
					}
				}
			}
		}
	}
}
