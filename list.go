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
	"errors"
	"math/rand"
	"sync"
)

type CompareFunc[K any] func(K, K) int

type SkipListConfig struct {
	Locked bool
}

type slNode[K any, V any] struct {
	key  K
	val  V
	next []*slNode[K, V] // next[i] means current node on level i
	prev *slNode[K, V]
}

type SkipList[K any, V any] struct {
	lock   bool
	comp   CompareFunc[K]
	mu     sync.RWMutex
	count  int
	height int
	head   *slNode[K, V]
	first  *slNode[K, V]
	last   *slNode[K, V]
}

func NewSkipList[K any, V any](cfg *SkipListConfig, comp CompareFunc[K]) *SkipList[K, V] {
	return &SkipList[K, V]{lock: cfg.Locked, comp: comp, height: -1}
}

func (s *SkipList[K, V]) Len() int {
	if s.lock {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	return s.count
}

func (s *SkipList[K, V]) Less(k1, k2 K) bool {
	return s.comp(k1, k2) < 0
}

// find the least element >= key, if key exists, will return true flag.
func (s *SkipList[K, V]) Search(key K) (k K, v V, false bool) {
	if s.lock {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	p := s.searchNode(key)
	if p == nil {
		return
	}
	if p.next[0] != nil {
		return p.next[0].key, p.next[0].val, s.comp(p.next[0].key, key) == 0
	}
	return
}

func (s *SkipList[K, V]) searchNode(key K) *slNode[K, V] {
	p := s.head
	l := s.height
	for p != nil && l >= 0 {
		for p.next[l] != nil && s.comp(p.next[l].key, key) < 0 {
			p = p.next[l]
		}
		l--
	}
	return p
}

// Query all elements within the closed interval [low,high].
// If there are no elements within this interval, return nil.
func (s *SkipList[K, V]) Range(low, high K) ([]K, []V) {
	if s.comp(high, low) < 0 {
		return nil, nil
	}
	if s.lock {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	p := s.searchNode(low)
	if p == nil {
		return nil, nil
	}
	var ks []K
	var vs []V
	for p = p.next[0]; p != nil; p = p.next[0] {
		if s.comp(p.key, high) > 0 {
			break
		}
		ks = append(ks, p.key)
		vs = append(vs, p.val)
	}
	return ks, vs
}

// Scan all elements in ascending order and call the f function for each element.
// If the call returns a non nil error, the scan will terminate immediately
func (s *SkipList[K, V]) Scan(f func(K, V) error) error {
	if s.lock {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if s.count == 0 {
		return errElemNotExists
	}
	for p := s.head.next[0]; p != nil; p = p.next[0] {
		if err := f(p.key, p.val); err != nil {
			return err
		}
	}
	return nil
}

// Scan all elements within the closed interval [low,high] in ascending order, and call the f function for each element.
// If the call returns a non nil error, the scan will terminate immediately
func (s *SkipList[K, V]) RangeScan(low, high K, f func(K, V) error) error {
	if s.comp(high, low) < 0 {
		return errElemNotExists
	}
	if s.lock {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	p := s.searchNode(low)
	if p == nil {
		return errElemNotExists
	}
	for p = p.next[0]; p != nil; p = p.next[0] {
		if s.comp(p.key, high) > 0 {
			break
		}
		if err := f(p.key, p.val); err != nil {
			return err
		}
	}
	return nil
}

// Search for the smallest element in the list. If the element does not exist (the list is empty), return false.
func (s *SkipList[K, V]) First() (k K, v V, false bool) {
	if s.lock {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if s.first == nil {
		return
	}
	return s.first.key, s.first.val, true
}

// Query the maximum element in the list. If the element does not exist (the list is empty), return false.
func (s *SkipList[K, V]) Last() (k K, v V, false bool) {
	if s.lock {
		s.mu.RLock()
		defer s.mu.RUnlock()
	}
	if s.last == nil {
		return
	}
	return s.last.key, s.last.val, true
}

func (s *SkipList[K, V]) randHeight() int {
	n, h := rand.Uint32(), 0
	for n&1 != 0 {
		h++
		n = n >> 1
	}
	if h > s.height {
		return s.height + 1
	}
	return h
}

// Update the elements in the list. If the update is successful, return true.
// If the key does not exist, do not take any action and return false.
func (s *SkipList[K, V]) Update(key K, val V) bool {
	p := s.searchNode(key)
	if p == nil || p.next[0] == nil {
		return false
	}
	if s.comp(p.next[0].key, key) == 0 {
		p.next[0].val = val
		return true
	}
	return false
}

func (s *SkipList[K, V]) newNode(key K, val V, h int) *slNode[K, V] {
	nd := &slNode[K, V]{key: key, val: val}
	nd.next = make([]*slNode[K, V], h+1)
	nd.prev = nd
	return nd
}

// Insert a new element into the list.
// Note that this method will not overwrite any existing keys,
// meaning that if the key already exists, there will be duplicate keys in the list after insertion.
func (s *SkipList[K, V]) Insert(key K, val V) {
	if s.lock {
		s.mu.Lock()
		defer s.mu.Unlock()
	}
	if s.head == nil {
		nd := s.newNode(key, val, 0)
		s.head = &slNode[K, V]{next: make([]*slNode[K, V], 0)}
		s.head.next = append(s.head.next, nd)
		s.height, s.count = 0, 1
		s.first, s.last = nd, nd
		return
	}
	// search path.
	stack := make([]*slNode[K, V], s.height+1)
	p := s.head
	l := s.height
	for l >= 0 {
		for p.next[l] != nil && s.comp(p.next[l].key, key) < 0 {
			p = p.next[l]
		}
		stack[l] = p
		l--
	}
	h := s.randHeight()
	nd := s.newNode(key, val, h)
	for i := 0; i < h-s.height; i++ {
		stack = append(stack, s.head)
		s.head.next = append(s.head.next, nil)
	}
	// insert in every level
	for i := 0; i <= h; i++ {
		nd.next[i] = stack[i].next[i]
		stack[i].next[i] = nd
		if i == 0 {
			// prev
			nd.prev = stack[i]
			if nd.next[i] != nil {
				nd.next[i].prev = nd
			}

			// update first/last ptr
			if s.first == nil || s.comp(key, s.first.key) <= 0 {
				s.first = nd
			}
			if s.last == nil || s.comp(key, s.last.key) > 0 {
				s.last = nd
			}
			if nd == s.first {
				nd.prev = s.last
			}
			if nd == s.last {
				s.first.prev = nd
			}
		}
	}
	s.count++
	s.height = max(s.height, h)
}

// delete a element and return value if its key exists, if not return false.
// If there are multiple elements with the same key in the list, delete the one closest to the head of the list.
func (s *SkipList[K, V]) Delete(key K) (v V, ok bool) {
	p, l := s.head, s.height
	for l >= 0 {
		for p.next[l] != nil && s.comp(p.next[l].key, key) < 0 {
			p = p.next[l]
		}
		q := p.next[l]
		if q != nil && s.comp(q.key, key) == 0 {
			v = q.val
			ok = true
			r := q.next[l]
			// level0 is special case
			if l == 0 {
				if r != nil {
					r.prev = p
					if q == s.first {
						s.first = r
						r.prev = s.last
					}
				} else { // q == s.last
					if p == s.head {
						s.first, s.last = nil, nil
					} else {
						s.last = p
						s.head.next[0].prev = p
					}
				}
			}
			q.next[l] = nil
			p.next[l] = r
			if p == s.head && p.next[l] == nil {
				s.height--
			}
		}
		l--
	}
	if ok {
		s.count--
	}
	return
}

// Delete all elements corresponding to keys,if key not exists return false.
func (s *SkipList[K, V]) DeleteAll(key K) bool {
	var flag bool
	for {
		_, ok := s.Delete(key)
		if ok {
			flag = true
		} else {
			break
		}
	}
	return flag
}

// Return a query cursor for more flexible traversal of list elements.
func (s *SkipList[K, V]) Cursor() Cursor[K, V] {
	return &slCuesor[K, V]{list: s}
}

type slCuesor[K any, V any] struct {
	mu    sync.Mutex
	opend bool
	list  *SkipList[K, V]
	cur   *slNode[K, V]
}

func (c *slCuesor[K, V]) Open() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.opend {
		return errors.New("current cursor already opened")
	}
	if c.list.lock {
		c.list.mu.RLock()
	}
	c.opend = true
	return nil
}

func (c *slCuesor[K, V]) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.opend {
		if c.list.lock {
			c.list.mu.RUnlock()
		}
	}
	c.opend = false
}

// The Seek(key) method locates the cursor at the key.
// If the key does not exist, it locates at the next key and returns it
func (c *slCuesor[K, V]) Seek(key K) (k K, v V, false bool) {
	if c.list.lock && !c.opend {
		return
	}
	p := c.list.searchNode(key)
	if p == nil {
		return
	}
	c.cur = p.next[0]
	if c.cur != nil {
		return c.cur.key, c.cur.val, c.list.comp(c.cur.key, key) == 0
	}
	return
}

// The First method locates the cursor at the minimum element of the set.
// If there is no minimum element (the set is empty), it returns false
func (c *slCuesor[K, V]) First() (k K, v V, false bool) {
	if c.list.lock && !c.opend {
		return
	}
	c.cur = c.list.first
	if c.cur == nil {
		return
	}
	return c.cur.key, c.cur.val, true
}

// The Last method locates the cursor at the maximum element of the set.
// If there is no maximum element (the set is empty), it returns false.
func (c *slCuesor[K, V]) Last() (k K, v V, false bool) {
	if c.list.lock && !c.opend {
		return
	}
	c.cur = c.list.last
	if c.cur == nil {
		return
	}
	return c.cur.key, c.cur.val, true
}

// HasNext returns whether the next element exists relative to the current cursor position.
func (c *slCuesor[K, V]) HasNext() bool {
	if c.list.lock && !c.opend {
		return false
	}
	return c.cur != nil && c.cur.next[0] != nil
}

// Next() moves the cursor backwards and returns the element.
// If the element does not exist, it returns a type zero value.
func (c *slCuesor[K, V]) Next() (k K, v V) {
	if c.list.lock && !c.opend {
		return
	}
	if c.cur != nil && c.cur.next[0] != nil {
		c.cur = c.cur.next[0]
		return c.cur.key, c.cur.val
	}
	return
}

// HasPrev() returns whether the previous element exists relative to the current cursor position.
func (c *slCuesor[K, V]) HasPrev() bool {
	if c.list.lock && !c.opend {
		return false
	}
	if c.cur == nil || c.cur.prev == c.list.last {
		return false
	}
	return true
}

// Prev() moves the cursor forward and returns the element.
// If the element does not exist, it returns a type value of zero.
func (c *slCuesor[K, V]) Prev() (k K, v V) {
	if c.list.lock && !c.opend {
		return
	}
	if c.cur != nil || c.cur.prev != c.list.last {
		c.cur = c.cur.prev
		return c.cur.key, c.cur.val
	}
	return
}
