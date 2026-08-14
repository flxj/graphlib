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

type trieNode interface {
	child(byte) trieNode
	tail() bool
	set(byte, trieNode)
}

type tNode256[T any] struct {
	flag int8
	val  T
	ch   [256]*tNode256[T]
}

type Trie[T any] struct {
	cnt  int
	root *tNode256[T]
}

func (t *Trie[T]) Len() int {
	return t.cnt
}

func (t *Trie[T]) Insert(s []byte, v T) {
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
	p.val = v
	if (p.flag & 1) != 0 {
		return
	}
	p.flag |= 1 // setting tail flag
	t.cnt++
}

func (t *Trie[T]) Search(s []byte) (v T, ok bool) {
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

func (t *Trie[T]) Prefix(s []byte) ([][]byte, []T) {
	p := t.root
	for _, b := range s {
		if p == nil {
			return nil, nil
		}
		p = p.ch[b]
	}
	var keys [][]byte
	var vals []T
	ks, vs := t.all(p)
	for i := 0; i < len(ks); i++ {
		key := make([]byte, len(s)+len(ks[i]))
		copy(key, s)
		copy(key[len(s):], ks[i])
		keys = append(keys, key)
		vals = append(vals, vs[i])
	}
	return keys, vals
}

func (t *Trie[T]) all(node *tNode256[T]) ([][]byte, []T) {
	if node == nil {
		return nil, nil
	}
	var keys [][]byte
	var vals []T
	if (node.flag & 3) == 1 {
		keys = append(keys, []byte{})
		vals = append(vals, node.val)
	}

	for b, q := range node.ch {
		ks, vs := t.all(q)
		for i, s := range ks {
			key := make([]byte, len(s)+1)
			key[0] = byte(b)
			copy(key[1:], s)
			keys = append(keys, key)
			vals = append(vals, vs[i])
		}
	}
	return keys, vals
}

func (t *Trie[T]) Delete(s []byte) bool {
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

func (t *Trie[T]) scan(pre []byte, node *tNode256[T], fn func([]byte, T) error) error {
	if node == nil {
		return nil
	}
	if (node.flag & 1) != 0 {
		if err := fn(pre, node.val); err != nil {
			return err
		}
	}
	for i, p := range node.ch {
		if p != nil {
			s := make([]byte, len(pre)+1)
			copy(s, pre)
			s[len(pre)] = byte(i)
			err := t.scan(s, p, fn)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Trie[T]) Scan(fn func([]byte, T) error) error {
	per := []byte{}
	return t.scan(per, t.root, fn)
}

// Compressed Trie
type PatriciaTrie[T any] struct {
}
