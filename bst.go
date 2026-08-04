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
