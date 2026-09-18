package heap

import (
	"github.com/flxj/graphlib/collection"
)

type Node[T, P any] struct {
	Val  T
	Rank P

	child   *Node[T, P] // first child
	sibling *Node[T, P] // next sibling
	prev    *Node[T, P] // prev，or parent（if itself is first child）
	parent  *Node[T, P] // parent node

	owner *PairingHeap[T, P]
}

type PairingHeap[T, P any] struct {
	root *Node[T, P]
	size int
	less collection.Less[P]
}

func NewPairingHeap[T, P any](less collection.Less[P]) *PairingHeap[T, P] {
	return &PairingHeap[T, P]{less: less}
}

func (h *PairingHeap[T, P]) Len() int { return h.size }

func (h *PairingHeap[T, P]) Top() (*Node[T, P], bool) {
	return h.root, h.root != nil
}

func (h *PairingHeap[T, P]) Contains(n *Node[T, P]) bool {
	return n != nil && n.owner == h
}

func (h *PairingHeap[T, P]) Push(n *Node[T, P]) {
	if n == nil || n.owner == h {
		return
	}
	n.owner = h
	h.root = merge(h.root, n, h.less)
	h.size++
}

func (h *PairingHeap[T, P]) Pop() (n *Node[T, P], ok bool) {
	if h.root != nil {
		n, ok = h.root, true
		h.root = mergeChildren(h.root, h.less)
		h.size--
		n.owner = nil
	}
	return
}

func (h *PairingHeap[T, P]) Merge(other *PairingHeap[T, P]) {
	if other == nil || other.root == nil {
		return
	}
	h.root = merge(h.root, other.root, h.less)
	h.size += other.size
	other.root = nil
	other.size = 0
}

func (h *PairingHeap[T, P]) Decrease(node *Node[T, P], newRank P) bool {
	if node == nil || node.owner != h {
		return false
	}
	if !h.less(newRank, node.Rank) {
		return false
	}
	node.Rank = newRank
	if node.parent == nil {
		return true
	}
	cut(node)
	h.root = merge(h.root, node, h.less)
	return true
}

func (h *PairingHeap[T, P]) Remove(node *Node[T, P]) bool {
	if node == nil || node.owner != h {
		return false
	}
	if node.parent == nil {
		h.root = mergeChildren(node, h.less)
		h.size--
		node.owner = nil
		return true
	}
	cut(node)
	sub := mergeChildren(node, h.less)
	h.root = merge(h.root, sub, h.less)
	h.size--
	node.owner = nil
	return true
}

func merge[T, P any](a, b *Node[T, P], less collection.Less[P]) *Node[T, P] {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if less(b.Rank, a.Rank) {
		a, b = b, a
	}
	b.parent = a
	b.prev = nil
	b.sibling = a.child
	if a.child != nil {
		a.child.prev = b
	}
	a.child = b
	return a
}

func mergeChildren[T, P any](node *Node[T, P], less collection.Less[P]) *Node[T, P] {
	first := node.child
	node.child = nil

	if first == nil {
		return nil
	}

	var pairs []*Node[T, P]
	cur := first
	for cur != nil {
		a := cur
		b := cur.sibling
		var next *Node[T, P]
		if b != nil {
			next = b.sibling
			b.sibling = nil
		}
		a.sibling = nil
		a.prev = nil
		if b != nil {
			b.prev = nil
		}

		if b != nil {
			pairs = append(pairs, merge(a, b, less))
		} else {
			pairs = append(pairs, a)
		}
		cur = next
	}

	if len(pairs) == 0 {
		return nil
	}
	result := pairs[len(pairs)-1]
	for i := len(pairs) - 2; i >= 0; i-- {
		result = merge(pairs[i], result, less)
	}
	return result
}

func cut[T, P any](n *Node[T, P]) {
	if n.parent == nil {
		return
	}
	if n.prev != nil {
		n.prev.sibling = n.sibling
	} else {
		n.parent.child = n.sibling
	}
	if n.sibling != nil {
		n.sibling.prev = n.prev
	}
	n.parent = nil
	n.sibling = nil
	n.prev = nil
}
