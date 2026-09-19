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

package trie

import (
	"bytes"
	"sync"
)

type pnode[V any] struct {
	bit    int
	key    []byte
	val    V
	preLen int
	left   *pnode[V]
	right  *pnode[V]
}

func (n *pnode[V]) isLeaf() bool { return n.bit < 0 }

func bitAt(key []byte, i int) int {
	if i/8 >= len(key) {
		return 0
	}
	b := key[i/8]
	return int((b >> (7 - uint(i%8))) & 1)
}

func firstDiffBit(a, b []byte) int {
	maxBits := max(len(a), len(b)) * 8
	for i := 0; i < maxBits; i++ {
		if bitAt(a, i) != bitAt(b, i) {
			return i
		}
	}
	return -1
}

type PatriciaTrie[V any] struct {
	lock bool
	mu   sync.RWMutex
	root *pnode[V]
	size int
}

func NewPatriciaTrie[V any](lock bool) *PatriciaTrie[V] {
	return &PatriciaTrie[V]{lock: lock}
}

func (t *PatriciaTrie[V]) Len() int {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	return t.size
}

func (t *PatriciaTrie[V]) Search(key []byte) (V, bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	return t.search(key)
}

func (t *PatriciaTrie[V]) search(key []byte) (V, bool) {
	var zero V
	n := t.root
	if n == nil {
		return zero, false
	}
	for !n.isLeaf() {
		if bitAt(key, n.bit) == 0 {
			n = n.left
		} else {
			n = n.right
		}
	}
	if bytes.Equal(n.key, key) {
		return n.val, true
	}
	return zero, false
}

func (t *PatriciaTrie[V]) Insert(key []byte, value V) {
	if t.lock {
		t.mu.Lock()
		defer t.mu.Unlock()
	}

	k := make([]byte, len(key))
	copy(k, key)

	if t.root == nil {
		t.root = &pnode[V]{bit: -1, key: k, val: value}
		t.size++
		return
	}

	leaf := t.root
	for !leaf.isLeaf() {
		if bitAt(k, leaf.bit) == 0 {
			leaf = leaf.left
		} else {
			leaf = leaf.right
		}
	}

	if bytes.Equal(leaf.key, k) {
		leaf.val = value
		return
	}

	diffBit := firstDiffBit(k, leaf.key)

	var parent *pnode[V]
	pos := t.root
	for !pos.isLeaf() && pos.bit < diffBit {
		parent = pos
		if bitAt(k, pos.bit) == 0 {
			pos = pos.left
		} else {
			pos = pos.right
		}
	}

	newLeaf := &pnode[V]{bit: -1, key: k, val: value}
	newInternal := &pnode[V]{bit: diffBit}

	if bitAt(k, diffBit) == 0 {
		newInternal.left = newLeaf
		newInternal.right = pos
	} else {
		newInternal.left = pos
		newInternal.right = newLeaf
	}

	if parent == nil {
		t.root = newInternal
	} else if parent.left == pos {
		parent.left = newInternal
	} else {
		parent.right = newInternal
	}
	t.size++
}

func (t *PatriciaTrie[V]) Delete(key []byte) bool {
	if t.lock {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	deleted := false
	t.root = t.delete(t.root, key, &deleted)
	if deleted {
		t.size--
	}
	return deleted
}

func (t *PatriciaTrie[V]) delete(n *pnode[V], key []byte, deleted *bool) *pnode[V] {
	if n == nil {
		return nil
	}
	if n.isLeaf() {
		if bytes.Equal(n.key, key) {
			*deleted = true
			return nil
		}
		return n
	}

	var child *pnode[V]
	if bitAt(key, n.bit) == 0 {
		child = n.left
	} else {
		child = n.right
	}
	newChild := t.delete(child, key, deleted)

	if !*deleted {
		return n
	}

	if newChild == nil {
		if child == n.left {
			return n.right
		}
		return n.left
	}
	if child == n.left {
		n.left = newChild
	} else {
		n.right = newChild
	}
	return n
}

func (t *PatriciaTrie[V]) Scan(fn func(key []byte, value V) bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	if t.root == nil {
		return
	}
	type frame struct {
		n *pnode[V]
	}
	stack := []*pnode[V]{}
	cur := t.root
	for cur != nil || len(stack) > 0 {
		for cur != nil {
			stack = append(stack, cur)
			if !cur.isLeaf() {
				cur = cur.left
			} else {
				cur = nil
			}
		}
		cur = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if cur.isLeaf() {
			if !fn(cur.key, cur.val) {
				return
			}
		}
		if !cur.isLeaf() {
			cur = cur.right
		} else {
			cur = nil
		}
	}
}

func (t *PatriciaTrie[V]) ScanPrefix(prefix []byte, fn func(key []byte, value V) bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	if t.root == nil {
		return
	}
	t.scanFilter(t.root, prefix, fn)
}

func (t *PatriciaTrie[V]) scanFilter(n *pnode[V], prefix []byte, fn func([]byte, V) bool) bool {
	if n == nil {
		return true
	}
	if n.isLeaf() {
		if hasPrefix(n.key, prefix) {
			return fn(n.key, n.val)
		}
		return true
	}
	if !t.scanFilter(n.left, prefix, fn) {
		return false
	}
	return t.scanFilter(n.right, prefix, fn)
}

func hasPrefix(key, prefix []byte) bool {
	if len(key) < len(prefix) {
		return false
	}
	for i := range prefix {
		if key[i] != prefix[i] {
			return false
		}
	}
	return true
}

func (t *PatriciaTrie[V]) LongestPrefix(key []byte) ([]byte, V, bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}

	var (
		bestKey []byte
		bestVal V
		found   bool
	)

	n := t.root
	for n != nil {
		if n.isLeaf() {
			if hasPrefix(key, n.key) {
				bestKey = n.key
				bestVal = n.val
				found = true
			}
			break
		}
		if bitAt(key, n.bit) == 0 {
			n = n.left
		} else {
			n = n.right
		}
	}
	return bestKey, bestVal, found
}

/*
func hasPrefixBits(key, prefix []byte, prefixLen int) bool {
	for i := 0; i < prefixLen; i++ {
		if bitAt(key, i) != bitAt(prefix, i) {
			return false
		}
	}
	return true
}

func (t *PatriciaTrie[V]) LongestPrefix(key []byte) ([]byte, V, bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	var (
		bestKey []byte
		bestVal V
		found   bool
	)
	n := t.root
	for n != nil {
		if n.isLeaf() {
			if hasPrefixBits(key, n.key, n.preLen) {
				bestKey = n.key
				bestVal = n.val
				found = true
			}
			break
		}
		if bitAt(key, n.bit) == 0 {
			n = n.left
		} else {
			n = n.right
		}
	}
	return bestKey, bestVal, found
}
*/
