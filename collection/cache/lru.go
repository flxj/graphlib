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

package cache

import (
	"container/list"
	"sync"
)

type entry[K comparable, V any] struct {
	key  K
	val  V
	freq int
}

type LRU[K comparable, V any] struct {
	mu      sync.Mutex
	cap     int
	ll      *list.List
	elems   map[K]*list.Element
	onEvict func(key K, value V)
}

func NewLRU[K comparable, V any](capacity int, onEvict func(K, V)) *LRU[K, V] {
	if capacity <= 0 {
		panic("lru: capacity must be positive")
	}
	return &LRU[K, V]{
		cap:     capacity,
		ll:      list.New(),
		elems:   make(map[K]*list.Element, capacity),
		onEvict: onEvict,
	}
}

func (c *LRU[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.elems[key]; ok {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry[K, V]).val, true
	}
	var zero V
	return zero, false
}

func (c *LRU[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.elems[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value.(*entry[K, V]).val = value
		return
	}

	ele := c.ll.PushFront(&entry[K, V]{key: key, val: value})
	c.elems[key] = ele

	if c.ll.Len() > c.cap {
		c.removeOldest()
	}
}

func (c *LRU[K, V]) removeOldest() {
	ele := c.ll.Back()
	if ele == nil {
		return
	}
	c.ll.Remove(ele)
	kv := ele.Value.(*entry[K, V])
	delete(c.elems, kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.val)
	}
}

func (c *LRU[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}
