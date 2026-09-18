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
	"math"
	"sync"
)

type LFU[K comparable, V any] struct {
	mu  sync.Mutex
	cap int

	freqMap map[int]*list.List
	keyMap  map[K]*list.Element

	minFreq int
	maxFreq int

	onEvict func(K, V)
}

func NewLFU[K comparable, V any](capacity int, onEvict func(K, V)) *LFU[K, V] {
	if capacity <= 0 {
		panic("lfu: capacity must be positive")
	}
	return &LFU[K, V]{
		cap:     capacity,
		freqMap: make(map[int]*list.List),
		keyMap:  make(map[K]*list.Element, capacity),
		onEvict: onEvict,
	}
}

func (c *LFU[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero V
	ele, ok := c.keyMap[key]
	if !ok {
		return zero, false
	}
	e := ele.Value.(*entry[K, V])
	c.touch(ele, e)
	return e.val, true
}

func (c *LFU[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.keyMap[key]; ok {
		e := ele.Value.(*entry[K, V])
		e.val = value
		c.touch(ele, e)
		return
	}

	if len(c.keyMap) >= c.cap {
		c.evict()
	}

	e := &entry[K, V]{key: key, val: value, freq: 1}
	bucket := c.freqMap[1]
	if bucket == nil {
		bucket = list.New()
		c.freqMap[1] = bucket
	}
	c.keyMap[key] = bucket.PushFront(e)
	c.minFreq = 1
	if c.maxFreq < 1 {
		c.maxFreq = 1
	}
}

func (c *LFU[K, V]) touch(ele *list.Element, e *entry[K, V]) {
	oldFreq := e.freq
	oldBucket := c.freqMap[oldFreq]
	oldBucket.Remove(ele)
	if oldBucket.Len() == 0 {
		delete(c.freqMap, oldFreq)
		if oldFreq == c.minFreq {
			c.minFreq++
		}
	}

	newFreq := oldFreq + 1
	e.freq = newFreq
	bucket := c.freqMap[newFreq]
	if bucket == nil {
		bucket = list.New()
		c.freqMap[newFreq] = bucket
	}
	c.keyMap[e.key] = bucket.PushFront(e)

	if newFreq > c.maxFreq {
		c.maxFreq = newFreq
	}
}

func (c *LFU[K, V]) evict() {
	for c.minFreq <= c.maxFreq {
		bucket, ok := c.freqMap[c.minFreq]
		if !ok || bucket.Len() == 0 {
			delete(c.freqMap, c.minFreq)
			c.minFreq++
			continue
		}
		victim := bucket.Back()
		e := victim.Value.(*entry[K, V])
		bucket.Remove(victim)
		if bucket.Len() == 0 {
			delete(c.freqMap, c.minFreq)
		}
		delete(c.keyMap, e.key)
		if c.onEvict != nil {
			c.onEvict(e.key, e.val)
		}
		return
	}
}

/*
A data that used to be accessed frequently but is now cold,
due to its high frequency, will occupy the pit for a long time
and new hotspots cannot enter. This is called cache polarization.

Regularly reduce the frequency of all keys by half
(or shift to the right by one digit) to gradually invalidate
historical frequencies.
*/
func (c *LFU[K, V]) aging() {
	c.mu.Lock()
	defer c.mu.Unlock()
	minF, maxF := math.MaxInt, math.MinInt
	newFreqMap := make(map[int]*list.List)
	for f, bk := range c.freqMap {
		newFreq := f / 2
		if newFreq < 1 {
			newFreq = 1
		}
		if newFreq < minF {
			minF = newFreq
		}
		if newFreq > maxF {
			maxF = newFreq
		}
		for ele := bk.Front(); ele != nil; ele = ele.Next() {
			e := ele.Value.(*entry[K, V])
			e.freq = newFreq
			nb := newFreqMap[newFreq]
			if nb == nil {
				nb = list.New()
				newFreqMap[newFreq] = nb
			}
			c.keyMap[e.key] = nb.PushBack(e)
		}
	}
	c.freqMap = newFreqMap
	c.minFreq, c.maxFreq = minF, maxF
}
