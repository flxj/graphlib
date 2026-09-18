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
	"time"
)

type item[K comparable, V any] struct {
	key     K
	val     V
	inCache bool
	count   int
	head    int
	times   []time.Time
}

/*
The main purpose of LRU-K is to solve the problem of "cache pollution" in LRU algorithm,
and its core idea is to extend the criterion of "recently used once" to "recently used K times".
That is to say, data that has not reached K accesses will not be cached.
This also means that the number of accesses to cached data needs to be counted,
and access records cannot be infinitely recorded, and replacement algorithms need to
be used for replacement. When data needs to be eliminated, LRU-K will eliminate the data
with the K-th access time and the longest distance from the current time.
*/
type LRUK[K comparable, V any] struct {
	mu  sync.Mutex
	k   int
	cap int

	history *list.List
	cache   *list.List

	index   map[K]*list.Element
	onEvict func(key K, val V)
}

func NewLRUK[K comparable, V any](k, capacity int, onEvict func(K, V)) *LRUK[K, V] {
	if k < 1 {
		panic("lruk: k must be >= 1")
	}
	if capacity <= 0 {
		panic("lruk: capacity must be positive")
	}
	return &LRUK[K, V]{
		k:       k,
		cap:     capacity,
		history: list.New(),
		cache:   list.New(),
		index:   make(map[K]*list.Element, capacity),
		onEvict: onEvict,
	}
}

func (c *LRUK[K, V]) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.history.Len() + c.cache.Len()
}

func (c *LRUK[K, V]) Stats() (historyLen, cacheLen int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.history.Len(), c.cache.Len()
}

/*
There are two storage spaces, one called the history space,
which stores cache with access times less than K; Another area
is called cache space, which stores data more than K times.
The historical space eliminates expired data (FIFO, LRU, etc.)
according to a certain algorithm.

When the number of data accesses reaches K,
it is added to the cache space and deleted from the historical space.
*/
func (c *LRUK[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	ele, ok := c.index[key]
	if !ok {
		return zero, false
	}
	e := ele.Value.(*item[K, V])
	c.recordAccess(e)
	if e.inCache {
		c.cache.MoveToFront(ele)
	} else if e.count >= c.k {
		c.history.Remove(ele)
		e.inCache = true
		c.cache.PushFront(ele)
		c.evictIfNeeded()
	}
	return e.val, true
}

func (c *LRUK[K, V]) Put(key K, val V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.index[key]; ok {
		e := ele.Value.(*item[K, V])
		e.val = val
		c.recordAccess(e)
		if e.inCache {
			c.cache.MoveToFront(ele)
		} else if e.count >= c.k {
			c.history.Remove(ele)
			e.inCache = true
			c.cache.PushFront(ele)
		}
		c.evictIfNeeded()
		return
	}

	e := &item[K, V]{
		key:   key,
		val:   val,
		times: make([]time.Time, c.k),
	}
	c.recordAccess(e)
	ele := c.history.PushFront(e)
	c.index[key] = ele

	c.evictIfNeeded()
}

func (c *LRUK[K, V]) recordAccess(e *item[K, V]) {
	e.times[e.head] = time.Now()
	e.head = (e.head + 1) % c.k
	if e.count < c.k {
		e.count++
	}
}

func (c *LRUK[K, V]) kthTime(e *item[K, V]) time.Time {
	return e.times[e.head]
}

/*
LRU-K records the K-th last access time of each key, and uses the
K-th last access time to determine the elimination order:

Key with access count<K: placed in the history list and eliminated by LRU.
Key with access times ≥ K: Put it in the cache list and sort it out by "K-th most recent access time".

The earlier the K-th visit, the more likely it is to be eliminated.
*/
func (c *LRUK[K, V]) evictIfNeeded() {
	for c.history.Len()+c.cache.Len() > c.cap {
		if h := c.history.Back(); h != nil {
			if c.cache.Len() > 0 && c.cache.Len() >= c.cap {
				ct := c.cache.Back()
				ce := ct.Value.(*item[K, V])
				he := h.Value.(*item[K, V])
				if c.firstTime(he).Before(c.kthTime(ce)) {
					c.removeElement(h)
					continue
				}
				c.removeElement(ct)
				continue
			}
			c.removeElement(h)
			continue
		}
		if ct := c.cache.Back(); ct != nil {
			c.removeElement(ct)
			continue
		}
		return
	}
}

func (c *LRUK[K, V]) firstTime(e *item[K, V]) time.Time {
	if e.count < c.k {
		idx := (e.head - e.count + c.k) % c.k
		return e.times[idx]
	}
	return c.kthTime(e)
}

func (c *LRUK[K, V]) removeElement(ele *list.Element) {
	e := ele.Value.(*item[K, V])
	if e.inCache {
		c.cache.Remove(ele)
	} else {
		c.history.Remove(ele)
	}
	delete(c.index, e.key)
	if c.onEvict != nil {
		c.onEvict(e.key, e.val)
	}
}
