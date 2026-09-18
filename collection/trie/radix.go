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

type rNode[V any] struct {
	isLeaf bool
	label  []byte
	value  V
	child  map[byte]*rNode[V]
}

func newNode[V any](label string) *rNode[V] {
	return &rNode[V]{
		label: []byte(label),
		child: make(map[byte]*rNode[V]),
	}
}

func commonPrefixLen(a, b []byte) int {
	n := min(len(a), len(b))
	i := 0
	for i < n && a[i] == b[i] {
		i++
	}
	return i
}

type RadixTree[V any] struct {
	lock bool
	mu   sync.RWMutex
	root *rNode[V]
	size int
}

func NewRadixTree[V any](lock bool) *RadixTree[V] {
	return &RadixTree[V]{lock: lock, root: newNode[V]("")}
}

func (t *RadixTree[V]) Len() int {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	return t.size
}

func (t *RadixTree[V]) Insert(key []byte, value V) {
	if t.lock {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	t.insert(key, value)
}

func (t *RadixTree[V]) InsertString(key string, value V) {
	if t.lock {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	t.insert([]byte(key), value)
}

func (t *RadixTree[V]) insert(key []byte, value V) {
	if len(key) == 0 {
		if !t.root.isLeaf {
			t.size++
		}
		t.root.isLeaf = true
		t.root.value = value
		return
	}

	node := t.root
	for {
		child, ok := node.child[key[0]]
		if !ok {
			node.child[key[0]] = &rNode[V]{
				label:  key,
				child:  make(map[byte]*rNode[V]),
				isLeaf: true,
				value:  value,
			}
			t.size++
			return
		}

		common := commonPrefixLen(child.label, key)

		if common < len(child.label) {
			oldSuffix := child.label[common:]
			oldNode := &rNode[V]{
				label:  oldSuffix,
				child:  child.child,
				isLeaf: child.isLeaf,
				value:  child.value,
			}
			child.label = child.label[:common]
			child.child = map[byte]*rNode[V]{oldSuffix[0]: oldNode}
			child.isLeaf = false
			var zero V
			child.value = zero

			rest := key[common:]
			if len(rest) == 0 {
				child.isLeaf = true
				child.value = value
				t.size++
			} else {
				child.child[rest[0]] = &rNode[V]{
					label:  rest,
					child:  make(map[byte]*rNode[V]),
					isLeaf: true,
					value:  value,
				}
				t.size++
			}
			return
		}

		rest := key[common:]
		if len(rest) == 0 {
			if !child.isLeaf {
				t.size++
			}
			child.isLeaf = true
			child.value = value
			return
		}
		node = child
		key = rest
	}
}

func (t *RadixTree[V]) Search(key []byte) (V, bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	return t.search(key)
}

func (t *RadixTree[V]) SearchString(key string) (V, bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	return t.search([]byte(key))
}

func (t *RadixTree[V]) search(key []byte) (V, bool) {
	var zero V
	node := t.root
	for len(key) != 0 {
		child, ok := node.child[key[0]]
		if !ok {
			return zero, false
		}
		if len(key) < len(child.label) || !bytes.Equal(key[:len(child.label)], child.label) {
			return zero, false
		}
		key = key[len(child.label):]
		node = child
	}
	if node.isLeaf {
		return node.value, true
	}
	return zero, false
}

func (t *RadixTree[V]) ContainsString(key string) bool {
	_, ok := t.Search([]byte(key))
	return ok
}

func (t *RadixTree[V]) Contains(key []byte) bool {
	_, ok := t.Search(key)
	return ok
}

func (t *RadixTree[V]) Delete(key []byte) bool {
	if t.lock {
		t.mu.Lock()
		defer t.mu.Unlock()
	}

	deleted := false
	t.delete(t.root, key, &deleted)
	if deleted {
		t.size--
	}
	return deleted
}

func (t *RadixTree[V]) DeleteString(key string) bool {
	return t.Delete([]byte(key))
}

func (t *RadixTree[V]) delete(node *rNode[V], key []byte, deleted *bool) {
	if len(key) == 0 {
		if !node.isLeaf {
			return
		}
		node.isLeaf = false
		var zero V
		node.value = zero
		*deleted = true
		return
	}

	child, ok := node.child[key[0]]
	if !ok || len(key) < len(child.label) || !bytes.Equal(key[:len(child.label)], child.label) {
		return
	}

	rest := key[len(child.label):]
	t.delete(child, rest, deleted)

	if !*deleted {
		return
	}

	if !child.isLeaf && len(child.child) == 0 {
		delete(node.child, key[0])
		return
	}

	if !child.isLeaf && len(child.child) == 1 {
		for _, only := range child.child {
			child.label = append(child.label, only.label...)
			child.isLeaf = only.isLeaf
			child.value = only.value
			child.child = only.child
			break
		}
	}
}

func (t *RadixTree[V]) Scan(fn func(key []byte, value V) bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	t.scan(t.root, []byte{}, fn)
}

func (t *RadixTree[V]) scan(node *rNode[V], prefix []byte, fn func([]byte, V) bool) bool {
	type frame struct {
		node *rNode[V]
		acc  []byte
	}
	stack := []frame{{node, prefix}}
	for len(stack) > 0 {
		f := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		n := f.node
		full := append(f.acc, n.label...)

		if n.isLeaf {
			if !fn(full, n.value) {
				return false
			}
		}
		for _, c := range n.child {
			stack = append(stack, frame{c, full})
		}
	}
	return true
}

func (t *RadixTree[V]) ScanPrefix(prefix []byte, fn func(key []byte, value V) bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}

	node := t.root
	key := prefix
	for len(key) != 0 {
		child, ok := node.child[key[0]]
		if !ok {
			return
		}
		if len(key) >= len(child.label) && bytes.Equal(key[:len(child.label)], child.label) {
			key = key[len(child.label):]
			node = child
			continue
		}
		if len(key) < len(child.label) && bytes.Equal(child.label[:len(key)], key) {
			key = []byte{}
			node = child
			break
		}
		return
	}
	baseLen := len(prefix) - len(key)
	base := prefix[:baseLen]
	t.scan(node, base, fn)
}

func (t *RadixTree[V]) LongestPrefix(key []byte) ([]byte, V, bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}

	var (
		bestKey []byte
		bestVal V
		found   bool
	)

	node := t.root
	consumed := 0
	for {
		if node.isLeaf {
			bestKey = key[:consumed]
			bestVal = node.value
			found = true
		}
		if consumed == len(key) {
			break
		}
		child, ok := node.child[key[consumed]]
		if !ok {
			break
		}
		if len(key)-consumed < len(child.label) ||
			!bytes.Equal(key[consumed:consumed+len(child.label)], child.label) {
			break
		}
		consumed += len(child.label)
		node = child
	}
	return bestKey, bestVal, found
}
