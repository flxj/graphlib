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

type frame[K comparable, V any] struct {
	key     K
	val     V
	pin     int
	first   time.Time
	inYoung bool
}

type BufferPool[K comparable, V any] struct {
	mu     sync.Mutex
	cap    int
	size   int // young list length limit.
	window time.Duration

	young *list.List
	old   *list.List
	cache map[K]*list.Element

	onEvict func(K, V)
}

func NewBufferPool[K comparable, V any](
	capacity int,
	youngRatio float64,
	timeWindow time.Duration,
	onEvict func(K, V),
) *BufferPool[K, V] {
	if capacity <= 0 {
		panic("pool: capacity must be positive")
	}
	if youngRatio <= 0 || youngRatio >= 1 {
		youngRatio = 0.625
	}
	return &BufferPool[K, V]{
		cap:     capacity,
		size:    int(float64(capacity) * youngRatio),
		window:  timeWindow,
		young:   list.New(),
		old:     list.New(),
		cache:   make(map[K]*list.Element, capacity),
		onEvict: onEvict,
	}
}

func (p *BufferPool[K, V]) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.young.Len() + p.old.Len()
}

func (p *BufferPool[K, V]) Get(key K) (V, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	var zero V
	ele, ok := p.cache[key]
	if !ok {
		return zero, false
	}
	f := ele.Value.(*frame[K, V])
	f.pin++

	if f.inYoung {
		p.young.MoveToFront(ele)
	} else if time.Since(f.first) >= p.window {
		p.old.Remove(ele)
		f.inYoung = true
		p.young.PushFront(ele)
		p.trimYoung()
	}
	return f.val, true
}

func (p *BufferPool[K, V]) trimYoung() {
	for p.young.Len() > p.size {
		ele := p.young.Back()
		p.young.Remove(ele)
		f := ele.Value.(*frame[K, V])
		f.inYoung = false
		f.first = time.Now()
		p.old.PushFront(ele)
	}
}

func (p *BufferPool[K, V]) Unpin(key K) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if ele, ok := p.cache[key]; ok {
		frame := ele.Value.(*frame[K, V])
		if frame.pin > 0 {
			frame.pin--
		}
	}
}

func (p *BufferPool[K, V]) Put(key K, value V) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ele, ok := p.cache[key]; ok {
		f := ele.Value.(*frame[K, V])
		f.val = value
		f.pin++
		return
	}
	p.evictIfNeeded()
	f := &frame[K, V]{key: key, val: value, first: time.Now(), pin: 1}
	p.cache[key] = p.old.PushFront(f)
}

func (p *BufferPool[K, V]) evictIfNeeded() {
	for p.young.Len()+p.old.Len() > p.cap {
		ele := p.findEvictable()
		if ele == nil {
			return
		}
		f := ele.Value.(*frame[K, V])
		if f.inYoung {
			p.young.Remove(ele)
		} else {
			p.old.Remove(ele)
		}
		delete(p.cache, f.key)
		if p.onEvict != nil {
			p.onEvict(f.key, f.val)
		}
	}
}

func (p *BufferPool[K, V]) findEvictable() *list.Element {
	for e := p.old.Back(); e != nil; e = e.Prev() {
		f := e.Value.(*frame[K, V])
		if f.pin == 0 {
			return e
		}
	}
	return nil
}
