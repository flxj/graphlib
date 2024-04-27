package graphlib

type element[K comparable, V any, P comparable] struct {
	key   K
	value V
	rank  P
	index int
}

type binaryHeap[K comparable, V any, P comparable] struct {
	ready bool
	elems []*element[K, V, P]
	less  func(k1, k2 P) bool
}

func newBinaryHeap[K comparable, V any, P comparable](less func(k1, k2 P) bool) *binaryHeap[K, V, P] {
	return &binaryHeap[K, V, P]{
		less: less,
	}
}

func (h *binaryHeap[K, V, P]) valid() bool {
	return h.ready
}

func (h *binaryHeap[K, V, P]) init() {
	if !h.ready {
		for i := h.length() / 2; i > 0; i-- {
			h.shiftDown(i)
		}
		h.ready = true
	}
}

func (h *binaryHeap[K, V, P]) empty() bool {
	return len(h.elems) == 0
}

func (h *binaryHeap[K, V, P]) length() int {
	return len(h.elems)
}

func (h *binaryHeap[K, V, P]) push(e *element[K, V, P]) {
	if h.ready {
		// put the element to tail of heap.
		n := h.length()
		e.index = n
		h.elems = append(h.elems, e)
		h.shift(n)
	}
}

func (h *binaryHeap[K, V, P]) pop() *element[K, V, P] {
	if !h.ready || h.empty() {
		return nil
	}

	// exchange head and tail of elems.
	h.elems[0], h.elems[h.length()-1] = h.elems[h.length()-1], h.elems[0]
	h.elems[0].index = 0

	// remove the latest element.
	e := h.elems[h.length()-1]
	e.index = -1

	h.elems = h.elems[0 : h.length()-1]

	// top-down shift the heap from 0.
	h.shiftDown(0)

	return e
}

// shift up
func (h *binaryHeap[K, V, P]) shift(idx int) {
	h.shiftUp(idx)
	h.shiftDown(idx)
}

func (h *binaryHeap[K, V, P]) shiftUp(idx int) {
	var p int
	for i := idx; i > 0; {
		p = (i - 1) / 2
		if h.less(h.elems[p].rank, h.elems[i].rank) {
			return
		}
		h.elems[i], h.elems[p] = h.elems[p], h.elems[i]
		h.elems[i].index = i
		h.elems[p].index = p
		i = p
	}
}

func (h *binaryHeap[K, V, P]) shiftDown(idx int) {
	min := func(i, j int) int {
		if h.less(h.elems[i].rank, h.elems[j].rank) {
			return i
		}
		return j
	}
	var p1, p2 int
	for i := idx; i < h.length(); {
		p1 = 2*i + 1
		p2 = 2*i + 2
		if p2 < h.length() {
			p := min(p1, p2)
			if h.less(h.elems[i].rank, h.elems[p].rank) {
				return
			}
			h.elems[i], h.elems[p] = h.elems[p], h.elems[i]
			h.elems[i].index = i
			h.elems[p].index = p
			i = p
		} else if p1 < h.length() {
			if h.less(h.elems[i].rank, h.elems[p1].rank) {
				return
			}
			h.elems[i], h.elems[p1] = h.elems[p1], h.elems[i]
			h.elems[i].index = i
			h.elems[p1].index = p1
			i = p1
		} else {
			return
		}
	}
}

type priorityQueue[K comparable, V any, P comparable] struct {
	items map[K]*element[K, V, P]
	heap  *binaryHeap[K, V, P]
}

func newPriorityQueue[K comparable, V any, P comparable](less func(p1, p2 P) bool) *priorityQueue[K, V, P] {
	q := &priorityQueue[K, V, P]{
		items: make(map[K]*element[K, V, P]),
		heap:  newBinaryHeap[K, V, P](less),
	}
	q.heap.init()
	return q
}

// update modifies the priority and value of an Item in the queue.
func (q *priorityQueue[K, V, P]) Update(k K, priority P) {
	v, ok := q.items[k]
	if ok {
		v.rank = priority
		q.heap.shift(v.index)
	}
}

// update modifies the priority and value of an Item in the queue.
func (q *priorityQueue[K, V, P]) Push(key K, value V, priority P) {
	v, ok := q.items[key]
	if ok {
		v.value = value
		v.rank = priority
		q.heap.shift(v.index)
		return
	}
	//
	item := &element[K, V, P]{key: key, value: value, rank: priority}
	q.items[key] = item
	q.heap.push(item)
}

func (q *priorityQueue[K, V, P]) Pop() (K, V, P, bool) {
	var k K
	var v V
	var p P
	if len(q.items) != 0 {
		p := q.heap.pop()
		if p != nil {
			delete(q.items, p.key)
		}
		return p.key, p.value, p.rank, true
	}
	return k, v, p, false
}

func (q *priorityQueue[K, V, P]) Get(k K) P {
	var p P
	if v, ok := q.items[k]; ok {
		return v.rank
	}
	return p
}

func (q *priorityQueue[K, V, P]) Len() int {
	return len(q.items)
}
