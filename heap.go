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
}

func newBinaryHeap[K comparable, V any, P comparable](elems []*element[K, V, P]) *binaryHeap[K, V, P] {
	return &binaryHeap[K, V, P]{
		elems: elems,
	}
}

func (h *binaryHeap[K, V, P]) valid() bool {
	return h.ready
}

func (h *binaryHeap[K, V, P]) init() {
	if !h.ready {
		for i := h.length() / 2; i >= 0; i-- {
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
	if idx >= h.length() {
		return
	}
	//
}

func (h *binaryHeap[K, V, P]) shiftDown(idx int) {
	if idx >= h.length() {
		return
	}
	//
}
