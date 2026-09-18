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
	"github.com/flxj/graphlib/collection"
	"github.com/flxj/graphlib/collection/stack"
)

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

/*
A ScapeGoat tree is a self-balancing Binary Search Tree like AVL Tree, Red-Black Tree, Splay Tree, ..etc.
Search time is O(Log n) in worst case. Time taken by deletion and insertion is amortized O(Log n)
Unlike other self-balancing BSTs, ScapeGoat tree doesn't require extra space per node.
For example, Red Black Tree nodes are required to have color. In below implementation of ScapeGoat Tree,
we only have left, right and parent pointers in Node class. Use of parent is done for simplicity of implementation and can be avoided.
*/
type ScapegoatTree[K any, V any] struct {
	comp  collection.CompareFunc[K]
	num   int
	alpha float64
	root  *sgtNode[K, V]
}

// Create a scapegoat tree.
func NewScapegoatTree[K any, V any](alpha float64, comp collection.CompareFunc[K]) *ScapegoatTree[K, V] {
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

func (s *ScapegoatTree[K, V]) Compare(k1, k2 K) int {
	return s.comp(k1, k2)
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

// Insert key value data, and if the key already exists, update its value in place.
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

// Use key to query elements. If the element does not exist, return null data value and false flag.
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

// Delete key. If the key exists, delete it and return its value. If the key does not exist, return a false flag.
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

// Query the minimum element of the key.
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

// Query the maximum element of the key.
func (s *ScapegoatTree[K, V]) Max() (k K, v V, ok bool) {
	return s.max(s.root)
}

func (s *ScapegoatTree[K, V]) flatten(r *sgtNode[K, V]) ([]K, []V) {
	if r == nil {
		return nil, nil
	}
	var keys []K
	var vals []V
	stk := stack.NewStack[*sgtNode[K, V]]()
	p := r
	for !stk.IsEmpty() || p != nil {
		for p != nil {
			stk.Push(p)
			p = p.l
		}
		p, _ = stk.Pop()
		if !p.del {
			keys = append(keys, p.key)
			vals = append(vals, p.val)
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

func (s *ScapegoatTree[K, V]) Clean() {
	s.num = 0
	s.root = nil
}

func (s *ScapegoatTree[K, V]) Cursor() collection.Cursor[K, V] {
	return &sgtCursor[K, V]{tree: s, stk: stack.NewStack[*sgtPath[K, V]]()}
}

type sgtPath[K, V any] struct {
	node    *sgtNode[K, V]
	visited bool
}

type sgtCursor[K, V any] struct {
	tree *ScapegoatTree[K, V]
	prev *sgtNode[K, V]
	stk  *stack.Stack[*sgtPath[K, V]]
}

func (c *sgtCursor[K, V]) Open() error { return nil }
func (c *sgtCursor[K, V]) Close()      { c.reset() }

func (c *sgtCursor[K, V]) reset() {
	c.prev = nil
	c.stk.Clean()
}

// The Seek(key) method locates the cursor at the key.
// If the key does not exist, it locates at the next key and returns it.
func (c *sgtCursor[K, V]) Seek(key K) (k K, v V, ok bool) {
	if c.tree == nil {
		return
	}
	c.reset()
	for p := c.tree.root; p != nil; {
		cp := c.tree.comp(p.key, key)
		if cp == 0 {
			k, v, ok = p.key, p.val, true
			c.stk.Push(&sgtPath[K, V]{node: p, visited: true})
			return
		} else if cp > 0 {
			c.stk.Push(&sgtPath[K, V]{node: p})
			p = p.l
		} else {
			c.stk.Push(&sgtPath[K, V]{node: p})
			p = p.r
		}
	}
	if tp := c.stk.Top(); tp != nil && !tp.node.del {
		tp.visited = true
		k, v = tp.node.key, tp.node.val
	}
	return
}

// The First method locates the cursor at the minimum element of the set.
// If there is no minimum element (the set is empty), it returns false
func (c *sgtCursor[K, V]) First() (k K, v V, ok bool) {
	if c.tree == nil {
		return
	}
	c.reset()
	for p := c.tree.root; p != nil; {
		c.stk.Push(&sgtPath[K, V]{node: p})
		if p.active(p.l) > 0 {
			p = p.l
		} else {
			if p.del && p.active(p.r) > 0 {
				p = p.r
			} else {
				break
			}
		}
	}
	if tp := c.stk.Top(); tp != nil && !tp.node.del {
		tp.visited = true
		k, v, ok = tp.node.key, tp.node.val, true
	}
	return
}

// The Last method locates the cursor at the maximum element of the set.
// If there is no maximum element (the set is empty), it returns false.
func (c *sgtCursor[K, V]) Last() (k K, v V, ok bool) {
	if c.tree == nil {
		return
	}
	c.reset()
	for p := c.tree.root; p != nil; {
		c.stk.Push(&sgtPath[K, V]{node: p})
		if p.active(p.r) > 0 {
			p = p.r
		} else {
			if p.del && p.active(p.l) > 0 {
				p = p.l
			} else {
				break
			}
		}
	}
	if tp := c.stk.Top(); tp != nil && !tp.node.del {
		tp.visited = true
		k, v, ok = tp.node.key, tp.node.val, true
	}
	return
}

// HasNext returns whether the next element exists relative to the current cursor position.
func (c *sgtCursor[K, V]) HasNext() bool {
	if c.tree == nil {
		return false
	}
	var prev *sgtNode[K, V]
	for !c.stk.IsEmpty() {
		p := c.stk.Top()
		if prev != nil && c.tree.comp(prev.key, p.node.key) >= 0 {
			pp, _ := c.stk.Pop()
			prev = pp.node
			continue
		}
		// try to move cursor to next element,but not visited it.
		if !p.visited && !p.node.del {
			return true
		}
		// subtree not nil
		if p.node.active(p.node.r) > 0 {
			if prev == nil || c.tree.comp(prev.key, p.node.key) < 0 {
				// p's right subtree not visited.
				for q := p.node.r; q != nil; {
					c.stk.Push(&sgtPath[K, V]{node: q})
					if q.active(q.l) > 0 {
						q = q.l
					} else {
						if q.del && q.active(q.r) > 0 {
							q = q.r
						} else {
							break
						}
					}
				}
				return true
			}
		}
		pp, _ := c.stk.Pop()
		prev = pp.node
	}
	return false
}

// Next() moves the cursor backwards and returns the element.
// If the element does not exist, it returns a type zero value.
func (c *sgtCursor[K, V]) Next() (k K, v V) {
	if c.tree == nil {
		return
	}
	var prev *sgtNode[K, V]
	for !c.stk.IsEmpty() {
		p := c.stk.Top()
		if prev != nil && c.tree.comp(prev.key, p.node.key) >= 0 {
			pp, _ := c.stk.Pop()
			prev = pp.node
			continue
		}
		// current node not visited,so return it.
		if !p.visited && !p.node.del {
			p.visited = true
			return p.node.key, p.node.val
		}
		// subtree not nil
		if p.node.active(p.node.r) > 0 {
			// p's right subtree not visited.
			if prev == nil || c.tree.comp(prev.key, p.node.key) < 0 {
				for q := p.node.r; q != nil; {
					c.stk.Push(&sgtPath[K, V]{node: q})
					if q.active(q.l) > 0 {
						q = q.l
					} else {
						if q.del && q.active(q.r) > 0 {
							q = q.r
						} else {
							break
						}
					}
				}
				if pp := c.stk.Top(); pp != nil && !pp.node.del {
					pp.visited = true
					return pp.node.key, pp.node.val
				}
			}
		}
		pp, _ := c.stk.Pop()
		prev = pp.node
	}
	return
}

// HasPrev() returns whether the previous element exists relative to the current cursor position.
func (c *sgtCursor[K, V]) HasPrev() bool {
	if c.tree == nil {
		return false
	}
	var prev *sgtNode[K, V]
	for !c.stk.IsEmpty() {
		p := c.stk.Top()
		if prev != nil && c.tree.comp(prev.key, p.node.key) <= 0 {
			pp, _ := c.stk.Pop()
			prev = pp.node
			continue
		}
		// try to move cursor to next element,but not visited it.
		if !p.visited && !p.node.del {
			return true
		}
		// subtree not nil
		if p.node.active(p.node.l) > 0 {
			if prev == nil || c.tree.comp(prev.key, p.node.key) > 0 {
				// p's left subtree not visited.
				for q := p.node.l; q != nil; {
					c.stk.Push(&sgtPath[K, V]{node: q})
					if q.active(q.r) > 0 {
						q = q.r
					} else {
						if q.del && q.active(q.l) > 0 {
							q = q.l
						} else {
							break
						}
					}
				}
				return true
			}
		}
		pp, _ := c.stk.Pop()
		prev = pp.node
	}
	return false
}

// Prev() moves the cursor forward and returns the element.
// If the element does not exist, it returns a type value of zero.
func (c *sgtCursor[K, V]) Prev() (k K, v V) {
	if c.tree == nil {
		return
	}
	var prev *sgtNode[K, V]
	for !c.stk.IsEmpty() {
		p := c.stk.Top()
		if prev != nil && c.tree.comp(prev.key, p.node.key) <= 0 {
			pp, _ := c.stk.Pop()
			prev = pp.node
			continue
		}
		// current node not visited,so return it.
		if !p.visited && !p.node.del {
			p.visited = true
			return p.node.key, p.node.val
		}
		// subtree not nil
		if p.node.active(p.node.l) > 0 {
			// p's right subtree not visited.
			if prev == nil || c.tree.comp(prev.key, p.node.key) > 0 {
				for q := p.node.l; q != nil; {
					c.stk.Push(&sgtPath[K, V]{node: q})
					if q.active(q.r) > 0 {
						q = q.r
					} else {
						if q.del && q.active(q.l) > 0 {
							q = q.l
						} else {
							break
						}
					}
				}
				if pp := c.stk.Top(); pp != nil && !pp.node.del {
					pp.visited = true
					return pp.node.key, pp.node.val
				}
			}
		}
		pp, _ := c.stk.Pop()
		prev = pp.node
	}
	return
}
