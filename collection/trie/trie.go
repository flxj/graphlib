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

import "sync"

type tNode256[T any] struct {
	flag int8
	key  []byte
	val  T
	ch   [256]*tNode256[T]
}

type Trie[T any] struct {
	lock bool
	mu   sync.RWMutex
	cnt  int
	root *tNode256[T]
}

func NewTrie[T any](lock bool) *Trie[T] {
	return &Trie[T]{lock: lock}
}

func (t *Trie[T]) Len() int {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	return t.cnt
}

func (t *Trie[T]) Insert(s []byte, v T) {
	if t.lock {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	if t.root == nil {
		t.root = &tNode256[T]{}
	}
	p := t.root
	for _, b := range s {
		if p.ch[b] == nil {
			p.ch[b] = &tNode256[T]{}
		}
		p = p.ch[b]
	}
	p.key = s
	p.val = v
	if (p.flag & 1) != 0 {
		return
	}
	p.flag |= 1 // setting tail flag
	t.cnt++
}

func (t *Trie[T]) Search(s []byte) (v T, ok bool) {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	p := t.root
	for _, b := range s {
		if p == nil {
			return
		}
		p = p.ch[b]
	}
	if p == nil {
		return
	}
	return p.val, (p.flag & 1) != 0
}

func (t *Trie[T]) Update(s []byte, fn func(T) T) {
	p := t.root
	for _, b := range s {
		if p == nil {
			return
		}
		p = p.ch[b]
	}
	if p != nil && p.flag == 1 {
		p.val = fn(p.val)
	}
}

func (t *Trie[T]) ScanPrefix(s []byte, fn func([]byte, T) bool) bool {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	p := t.root
	for _, b := range s {
		if p == nil {
			return false
		}
		p = p.ch[b]
	}
	return t.scan(p, fn)
}

func (t *Trie[T]) scan(node *tNode256[T], fn func([]byte, T) bool) bool {
	if node == nil {
		return true
	}
	if (node.flag & 3) == 1 {
		if !fn(node.key, node.val) {
			return false
		}
	}
	for _, q := range node.ch {
		if q == nil {
			continue
		}
		if !t.scan(q, fn) {
			return false
		}
	}
	return true
}

func (t *Trie[T]) Delete(s []byte) bool {
	if t.lock {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	p := t.root
	for _, b := range s {
		if p == nil {
			return false
		}
		p = p.ch[b]
	}
	if p == nil || (p.flag&1) == 0 {
		return false
	}
	if (p.flag & 2) == 0 {
		t.cnt--
	}
	p.flag |= 2 // setting delete tag
	return true
}

func (t *Trie[T]) count(node *tNode256[T]) int {
	if node == nil {
		return 0
	}
	var c int
	if node.flag == 1 {
		c++
	}
	for _, p := range node.ch {
		c += t.count(p)
	}
	return c
}

func (t *Trie[T]) DeleteByPrefix(pre []byte) bool {
	if t.lock {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	var idx int
	var pp *tNode256[T]
	p := t.root
	for i, b := range pre {
		if p == nil {
			return false
		}
		pp, idx = p, i
		p = p.ch[b]
	}
	if p == nil || pp == nil {
		return false
	}
	t.cnt -= t.count(p)
	pp.ch[idx] = nil
	return true
}

func (t *Trie[T]) Scan(fn func([]byte, T) bool) bool {
	if t.lock {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	return t.scan(t.root, fn)
}
