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

type BSTKind uint16

const (
	KindSplayTree BSTKind = iota
	KindTreap
	KindRedBlackTree
	KindScapegoatTree
	KindBTree
	KindSkipList
)

// This interface describes the main methods of general BST.
type BinarySearchTree[K any, V any] interface {
	// The number of elements in the current tree.
	Len() int
	// Search for the element corresponding to the specified key.
	Search(K) (V, bool)
	// Insert or update elements.
	Insert(K, V)
	// Delete specified element.
	Delete(K) (V, bool)
	//Comparing elements.
	Compare(K, K) int
	//Return the minimum element.
	Min() (K, V, bool)
	// Return the maximum element.
	Max() (K, V, bool)
	// Create an Iterator
	Cursor() collection.Cursor[K, V]
	// Clear the current BST, which means deleting all elements.
	Clean()
}

// Create a BST object based on its type.
func NewBinarySearchTree[K any, V any](comp collection.CompareFunc[K], kind BSTKind) (BinarySearchTree[K, V], bool) {
	switch kind {
	case KindBTree:
		cfg := &BTreeConfig{MinDegree: DefaultBTreeMinDegree}
		return NewBTree[K, V](cfg, comp), true
	case KindTreap:
		return NewTreap[K, V](comp), true
	case KindRedBlackTree:
		return NewRedBlackTree[K, V](comp), true
	case KindScapegoatTree:
		return NewScapegoatTree[K, V](DefaultScapegoatTreeAlpha, comp), true
	case KindSplayTree:
		return NewSplayTree[K, V](comp), true
	case KindSkipList:
		cfg := &SkipListConfig{}
		return NewSkipList[K, V](cfg, comp), true
	default:
		return nil, false
	}
}

// General BST node interface
type bstNode[K, V any] interface {
	left() bstNode[K, V]
	right() bstNode[K, V]
	getKey() K
	getVal() V
}

type bstPath[K, V any] struct {
	node    bstNode[K, V]
	visited bool
}

// Implement an iterator for general BST.
type bstCursor[K, V any] struct {
	comp collection.CompareFunc[K]
	root bstNode[K, V]
	prev bstNode[K, V]
	stk  *stack.Stack[*bstPath[K, V]]
}

func newBSTCursor[K, V any](root bstNode[K, V], comp collection.CompareFunc[K]) *bstCursor[K, V] {
	c := &bstCursor[K, V]{
		comp: comp,
		root: root,
		stk:  stack.NewStack[*bstPath[K, V]](),
	}
	return c
}

func (c *bstCursor[K, V]) Open() error { return nil }
func (c *bstCursor[K, V]) Close()      { c.reset() }

func (c *bstCursor[K, V]) reset() {
	c.prev = nil
	c.stk.Clean()
}

// The Seek(key) method locates the cursor at the key.
// If the key does not exist, it locates at the next key and returns it.
func (c *bstCursor[K, V]) Seek(key K) (k K, v V, ok bool) {
	c.reset()
	for p := c.root; p != nil; {
		cp := c.comp(p.getKey(), key)
		if cp == 0 {
			k, v, ok = p.getKey(), p.getVal(), true
			c.stk.Push(&bstPath[K, V]{node: p, visited: true})
			return
		} else if cp > 0 {
			c.stk.Push(&bstPath[K, V]{node: p})
			p = p.left()
		} else {
			c.stk.Push(&bstPath[K, V]{node: p})
			p = p.right()
		}
	}
	if tp := c.stk.Top(); tp != nil {
		tp.visited = true
		k, v = tp.node.getKey(), tp.node.getVal()
	}
	return
}

// The First method locates the cursor at the minimum element of the set.
// If there is no minimum element (the set is empty), it returns false
func (c *bstCursor[K, V]) First() (k K, v V, ok bool) {
	c.reset()
	for p := c.root; p != nil; p = p.left() {
		c.stk.Push(&bstPath[K, V]{node: p})
	}
	if tp := c.stk.Top(); tp != nil {
		tp.visited = true
		k, v, ok = tp.node.getKey(), tp.node.getVal(), true
	}
	return
}

// The Last method locates the cursor at the maximum element of the set.
// If there is no maximum element (the set is empty), it returns false.
func (c *bstCursor[K, V]) Last() (k K, v V, ok bool) {
	c.reset()
	for p := c.root; p != nil; p = p.right() {
		c.stk.Push(&bstPath[K, V]{node: p})
	}
	if tp := c.stk.Top(); tp != nil {
		tp.visited = true
		k, v, ok = tp.node.getKey(), tp.node.getVal(), true
	}
	return
}

// HasNext returns whether the next element exists relative to the current cursor position.
func (c *bstCursor[K, V]) HasNext() bool {
	var prev bstNode[K, V]
	for !c.stk.IsEmpty() {
		p := c.stk.Top()
		if prev != nil && c.comp(prev.getKey(), p.node.getKey()) >= 0 {
			pp, _ := c.stk.Pop()
			prev = pp.node
			continue
		}
		// try to move cursor to next element,but not visited it.
		if !p.visited {
			return true
		}
		// subtree not nil
		if p.node.right() != nil {
			if prev == nil || c.comp(prev.getKey(), p.node.getKey()) < 0 {
				// p's right subtree not visited.
				for q := p.node.right(); q != nil; q = q.left() {
					c.stk.Push(&bstPath[K, V]{node: q})
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
func (c *bstCursor[K, V]) Next() (k K, v V) {
	var prev bstNode[K, V]
	for !c.stk.IsEmpty() {
		p := c.stk.Top()
		if prev != nil && c.comp(prev.getKey(), p.node.getKey()) >= 0 {
			pp, _ := c.stk.Pop()
			prev = pp.node
			continue
		}
		// current node not visited,so return it.
		if !p.visited {
			p.visited = true
			return p.node.getKey(), p.node.getVal()
		}
		// subtree not nil
		if p.node.right() != nil {
			// p's right subtree not visited.
			if prev == nil || c.comp(prev.getKey(), p.node.getKey()) < 0 {
				for q := p.node.right(); q != nil; q = q.left() {
					c.stk.Push(&bstPath[K, V]{node: q})
				}
				pp := c.stk.Top()
				pp.visited = true
				return pp.node.getKey(), pp.node.getVal()
			}
		}
		pp, _ := c.stk.Pop()
		prev = pp.node
	}
	return
}

// HasPrev() returns whether the previous element exists relative to the current cursor position.
func (c *bstCursor[K, V]) HasPrev() bool {
	var prev bstNode[K, V]
	for !c.stk.IsEmpty() {
		p := c.stk.Top()
		if prev != nil && c.comp(prev.getKey(), p.node.getKey()) <= 0 {
			pp, _ := c.stk.Pop()
			prev = pp.node
			continue
		}
		// try to move cursor to next element,but not visited it.
		if !p.visited {
			return true
		}
		// subtree not nil
		if p.node.left() != nil {
			if prev == nil || c.comp(prev.getKey(), p.node.getKey()) > 0 {
				// p's left subtree not visited.
				for q := p.node.left(); q != nil; q = q.right() {
					c.stk.Push(&bstPath[K, V]{node: q})
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
func (c *bstCursor[K, V]) Prev() (k K, v V) {
	var prev bstNode[K, V]
	for !c.stk.IsEmpty() {
		p := c.stk.Top()
		if prev != nil && c.comp(prev.getKey(), p.node.getKey()) <= 0 {
			pp, _ := c.stk.Pop()
			prev = pp.node
			continue
		}
		// current node not visited,so return it.
		if !p.visited {
			p.visited = true
			return p.node.getKey(), p.node.getVal()
		}
		// subtree not nil
		if p.node.left() != nil {
			// p's right subtree not visited.
			if prev == nil || c.comp(prev.getKey(), p.node.getKey()) > 0 {
				for q := p.node.left(); q != nil; q = q.right() {
					c.stk.Push(&bstPath[K, V]{node: q})
				}
				pp := c.stk.Top()
				pp.visited = true
				return pp.node.getKey(), pp.node.getVal()
			}
		}
		pp, _ := c.stk.Pop()
		prev = pp.node
	}
	return
}
