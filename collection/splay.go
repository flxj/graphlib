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

type stNode[K any, V any] struct {
	key K
	val V
	par *stNode[K, V]
	lef *stNode[K, V]
	rig *stNode[K, V]
}

func (n *stNode[K, V]) isRoot() bool {
	return n.par == nil
}

func (n *stNode[K, V]) isLeftChild() bool {
	return n.par != nil && n == n.par.lef
}

func (t *stNode[K, V]) left() bstNode[K, V] {
	return t.lef
}

func (t *stNode[K, V]) right() bstNode[K, V] {
	return t.rig
}

func (t *stNode[K, V]) getKey() K {
	return t.key
}

func (t *stNode[K, V]) getVal() V {
	return t.val
}

// A splay tree is a self-balancing binary search tree with the additional property that
// recently accessed elements are quick to access again. It performs basic operations such
// as insertion, look-up and removal in O(log(n)) amortized time. For many non-uniform sequences
// of operations, splay trees perform better than other search trees, even when the specific
// pattern of the sequence is unknown. The splay tree was invented by Daniel Sleator and Robert Tarjan.
type SplayTree[K any, V any] struct {
	comp  CompareFunc[K]
	count int
	root  *stNode[K, V]
}

// Create a splay tree. The caller needs to provide a method for comparing element keys.
func NewSplayTree[K any, V any](comp CompareFunc[K]) *SplayTree[K, V] {
	return &SplayTree[K, V]{comp: comp}
}

// splay a node when that node is not the root and we wish to transport it to the root.
func (s *SplayTree[K, V]) splay(x *stNode[K, V]) {
	if x == nil {
		return
	}
	rightRotate := func(x, p, gp *stNode[K, V]) {
		x.par = gp
		p.par = x
		p.lef = x.rig
		if p.lef != nil {
			p.lef.par = p
		}
		x.rig = p
		if gp != nil {
			if gp.lef == p {
				gp.lef = x
			} else {
				gp.rig = x
			}
		}
	}
	leftRotate := func(x, p, gp *stNode[K, V]) {
		x.par = gp
		p.par = x
		p.rig = x.lef
		if p.rig != nil {
			p.rig.par = p
		}
		x.lef = p
		if gp != nil {
			if gp.lef == p {
				gp.lef = x
			} else {
				gp.rig = x
			}
		}
	}
	for !x.isRoot() {
		p := x.par
		if p.isRoot() { // x the root’s child, Zig.
			if x.isLeftChild() {
				rightRotate(x, p, nil)
			} else {
				leftRotate(x, p, nil)
			}
		} else {
			gp := p.par
			if x.isLeftChild() {
				if p.isLeftChild() { // x a left-left child,Zig-Zig
					rightRotate(p, gp, gp.par)
					rightRotate(x, p, p.par)
				} else { // x is a right-left child,Zig-Zag
					rightRotate(x, p, gp)
					leftRotate(x, gp, gp.par)
				}
			} else {
				if p.isLeftChild() { // x is a left-right child,Zig-Zag
					leftRotate(x, p, gp)
					rightRotate(x, gp, gp.par)
				} else { // x is a right-right child,Zig-Zig
					leftRotate(p, gp, gp.par)
					leftRotate(x, p, p.par)
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
			if x.rig == nil {
				break
			}
			x = x.rig
		} else {
			if x.lef == nil {
				break
			}
			x = x.lef
		}
	}
	s.splay(x)
	return
}

// Use key to query elements. If the element does not exist, return null data value and false flag.
func (s *SplayTree[K, V]) Search(key K) (v V, ok bool) {
	root, ok := s.searchSubtree(s.root, key)
	s.root = root
	if ok {
		return root.val, true
	}
	return
}

// Insert key value data, and if the key already exists, update its value in place.
func (s *SplayTree[K, V]) Insert(key K, val V) {
	root, ok := s.searchSubtree(s.root, key)
	s.root = root
	if ok {
		s.root.val = val
		return
	}
	node := &stNode[K, V]{key: key, val: val}
	if s.root != nil {
		s.root.par = node
		if s.comp(s.root.key, key) > 0 {
			node.lef = s.root.lef
			if node.lef != nil {
				node.lef.par = node
			}
			s.root.lef = nil
			node.rig = s.root
		} else {
			node.rig = s.root.rig
			if node.rig != nil {
				node.rig.par = node
			}
			s.root.rig = nil
			node.lef = s.root
		}
	}
	s.root = node
	s.count++
}

// Delete key. If the key exists, delete it and return its value. If the key does not exist, return a false flag.
func (s *SplayTree[K, V]) Delete(key K) (V, bool) {
	v, ok := s.Search(key)
	if !ok {
		return v, false
	}
	if s.root == nil {
		return v, false
	}
	if s.root.lef == nil {
		r := s.root.rig
		s.root.rig = nil
		if r != nil {
			r.par = nil
		}
		s.root = r
	} else if s.root.rig == nil {
		l := s.root.lef
		s.root.lef = nil
		if l != nil {
			l.par = nil
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
		s.root.rig.par = nil // cut right subtree from s.root
		r, _ := s.searchSubtree(s.root.rig, key)
		s.root.lef.par = r
		r.lef = s.root.lef
		s.root.lef, s.root.rig = nil, nil
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

// Query the minimum element of the key.
func (s *SplayTree[K, V]) Min() (k K, v V, ok bool) {
	p := s.root
	for p != nil {
		if p.lef == nil {
			break
		}
		p = p.lef
	}
	if p != nil {
		k, v, ok = p.key, p.val, true
		s.splay(p)
	}
	return
}

// Query the maximum element of the key.
func (s *SplayTree[K, V]) Max() (k K, v V, ok bool) {
	p := s.root
	for p != nil {
		if p.rig == nil {
			break
		}
		p = p.rig
	}
	if p != nil {
		k, v, ok = p.key, p.val, true
		s.splay(p)
	}
	return
}

func (s *SplayTree[K, V]) Clean() {
	s.count = 0
	s.root = nil
}

func (s *SplayTree[K, V]) Cursor() Cursor[K, V] {
	return newBSTCursor(s.root, s.comp)
}
