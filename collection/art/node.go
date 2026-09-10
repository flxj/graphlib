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

package art

const (
	leaf nodeKind = iota
	node4
	node16
	node48
	node256
)

const (
	//maxPrefixLen int    = 16
	capNode4   uint16 = 4
	capNode16  uint16 = 16
	capNode48  uint16 = 48
	capNode256 uint16 = 256
)

// Node serves as an interface to uniformly represent the internal
// nodes and leaf nodes of ART, while providing a series of methods
// for modifying/querying node information.
type node[T any] interface {
	// node type.
	kind() nodeKind
	//Is the inner node already full.
	isFull() bool
	// The prefix string length contained in the inner node.
	prefixLen() int
	// Find the character at the specified position in the prefix string.
	prefixChar(int) byte
	// Truncate the prefix string from the specified position and discard
	// the part between the starting position and the truncation point.
	prefixCut(int)
	// Update node prefix.
	setPrefix([]byte)
	// The number of child nodes in the current node.
	count() uint16
	// Given a character, search for the corresponding child
	// node and return its position in the child list.
	find(byte) (node[T], uint16)
	// Update the child pointer at the specified location in the node child list.
	set(uint16, node[T])
	// Insert a new record (character, child pointer) into the current node.
	add(byte, node[T])
	// Delete specified child branch.
	del(byte)
	// end is the leaf node attach current inner node.
	end() node[T]
	// update end node.
	setEnd(node[T])
	// Reduce the current node.
	shrink() node[T]
	// Find the nth child node.
	nthChild(uint16) node[T]
	// Similar to the find method but returns the ranking of child nodes.
	rank(byte) (node[T], uint16)
	// Return the key of the current node.
	// Note that for inner nodes, this method returns a
	// prefix string, while for leaf nodes, it returns the entire key.
	key() Key
	// Return the value of the current node,
	// this method is only valid for leaf nodes.
	value() T
	// Update the value of leaf nodes.
	update(T)
	// Copy the content of the current node.
	clone() ([]byte, []node[T])
}

/*
Additionally, at the front of each inner node, a header of
constant size (e.g., 16 bytes) stores the node type, the number
of children, and the compressed path (cf. Section III-E)
*/
type header[T any] struct {
	nKind    nodeKind
	childLen uint16
	tail     node[T]
	preLen   int
	prefix   []byte
}

func (h *header[T]) kind() nodeKind        { return h.nKind }
func (h *header[T]) count() uint16         { return h.childLen }
func (h *header[T]) key() Key              { return Key(h.prefix[:h.preLen:h.preLen]) }
func (h *header[T]) prefixLen() int        { return h.preLen }
func (h *header[T]) prefixChar(i int) byte { return h.prefix[i] }
func (h *header[T]) end() node[T]          { return h.tail }
func (h *header[T]) setEnd(n node[T])      { h.tail = n }
func (h *header[T]) value() (v T)          { return }
func (h *header[T]) update(v T)            {}

func (h *header[T]) setPrefix(p []byte) {
	h.prefix = p
	h.preLen = len(p)
}

func (h *header[T]) prefixCut(i int) {
	if i > 0 && i < h.preLen {
		copy(h.prefix, h.prefix[i:])
		h.preLen -= i
	} else if i >= h.preLen {
		h.prefix = nil
		h.preLen = 0
	}
}

func (h *header[T]) isFull() bool {
	switch h.nKind {
	case node4:
		return h.childLen == capNode4
	case node16:
		return h.childLen == capNode16
	case node48:
		return h.childLen == capNode48
	case node256:
		return h.childLen == capNode256
	default:
		return false
	}
}

/*
Node4:
The smallest node type can store up to 4 child
pointers and uses an array of length 4 for keys and another
array of the same length for pointers. The keys and pointers
are stored at corresponding positions and the keys are sorted
*/
type artNode4[T any] struct {
	*header[T]
	keys  [capNode4]byte
	child [capNode4]node[T]
}

func (n *artNode4[T]) find(b byte) (node[T], uint16) {
	for i := uint16(0); i < n.childLen; i++ {
		if n.keys[i] > b {
			break
		} else if n.keys[i] == b {
			return n.child[i], i
		}
	}
	return nil, 0
}

func (n *artNode4[T]) add(b byte, c node[T]) {
	if n.childLen == 0 {
		n.keys[0] = b
		n.child[0] = c
		n.childLen = 1
		return
	}
	for i := uint16(0); i < n.childLen; i++ {
		if n.keys[i] > b {
			for j := n.childLen; j > i; j-- {
				n.keys[j] = n.keys[j-1]
				n.child[j] = n.child[j-1]
			}
			n.childLen++
			n.keys[i] = b
			n.child[i] = c
			return
		} else if n.keys[i] == b {
			n.child[i] = c
			return
		}
	}
	n.keys[n.childLen] = b
	n.child[n.childLen] = c
	n.childLen++
}

func (n *artNode4[T]) set(i uint16, ptr node[T]) {
	if i < capNode4 {
		n.child[i] = ptr
	}
}

func (n *artNode4[T]) del(b byte) {
	for i := uint16(0); i < n.childLen; i++ {
		if n.keys[i] == b {
			for j := i; j < n.childLen-1; j++ {
				n.keys[j] = n.keys[j+1]
				n.child[j] = n.child[j+1]
			}
			n.child[n.childLen-1] = nil
			n.childLen--
			return
		}
	}
}

func (n *artNode4[T]) nthChild(m uint16) node[T] {
	if m <= 0 || m > n.childLen {
		return nil
	}
	return n.child[m-1]
}

func (n *artNode4[T]) rank(b byte) (node[T], uint16) {
	ch, i := n.find(b)
	return ch, i + 1
}

func (n *artNode4[T]) clone() ([]byte, []node[T]) {
	return n.keys[:n.childLen], n.child[:n.childLen]
}

// if node4 has the "1 child + no terminal" case
// where the lone child absorbs the parent's prefix and branch byte.
func (n *artNode4[T]) shrink() node[T] {
	if n.childLen == 0 {
		return n.tail
	}
	if n.childLen == 1 && n.tail == nil {
		ch := n.child[0]
		if ch.kind() == leaf {
			return ch
		}
		return mergeTo([]byte(n.key()), n.keys[0], ch)
	}
	return n
}

func newNode4[T any](perLen int) *artNode4[T] {
	n := &artNode4[T]{header: &header[T]{nKind: node4}}
	if perLen > 0 {
		n.prefix = make([]byte, perLen)
	}
	return n
}

/*
Node16:
his node type is used for storing between 5 and
16 child pointers. Like the Node4, the keys and pointers
are stored in separate arrays at corresponding positions, but
both arrays have space for 16 entries. A key can be found
efficiently with binary search.
*/
type artNode16[T any] struct {
	*header[T]
	keys  [capNode16]byte
	child [capNode16]node[T]
}

func (n *artNode16[T]) find(b byte) (node[T], uint16) {
	l, r := uint16(0), n.childLen
	for l < r {
		m := (l + r) / 2
		if n.keys[m] == b {
			return n.child[m], m
		} else if n.keys[m] < b {
			l = m + 1
		} else {
			r = m
		}
	}
	return nil, l
}

func (n *artNode16[T]) add(b byte, c node[T]) {
	if n.childLen == 0 {
		n.keys[0] = b
		n.child[0] = c
		n.childLen = 1
		return
	}
	ch, i := n.find(b)
	if i < n.childLen && n.keys[i] == b {
		if ch == nil {
			n.childLen++
		}
		n.child[i] = c
		return
	}
	for j := n.childLen; j > i; j-- {
		n.keys[j] = n.keys[j-1]
		n.child[j] = n.child[j-1]
	}
	n.keys[i] = b
	n.child[i] = c
	n.childLen++
}

func (n *artNode16[T]) set(i uint16, ptr node[T]) {
	if i < capNode16 {
		n.child[i] = ptr
	}
}

func (n *artNode16[T]) del(b byte) {
	_, i := n.find(b)
	if i < n.childLen && n.keys[i] == b {
		for j := i; j < n.childLen-1; j++ {
			n.keys[j] = n.keys[j+1]
			n.child[j] = n.child[j+1]
		}
		n.child[n.childLen-1] = nil
		n.childLen--
	}
}

func (n *artNode16[T]) nthChild(m uint16) node[T] {
	if m <= 0 || m > n.childLen {
		return nil
	}
	return n.child[m-1]
}

func (n *artNode16[T]) rank(b byte) (node[T], uint16) {
	ch, i := n.find(b)
	return ch, i + 1
}

func (n *artNode16[T]) clone() ([]byte, []node[T]) {
	return n.keys[:n.childLen], n.child[:n.childLen]
}

func (n *artNode16[T]) shrink() node[T] {
	if n.childLen == 0 {
		return n.tail
	}
	if n.childLen <= capNode4 {
		n4 := newNode4[T](0)
		n4.header = n.header
		n4.header.nKind = node4
		copy(n4.keys[:n.childLen], n.keys[:n.childLen])
		copy(n4.child[:n.childLen], n.child[:n.childLen])
		return n4
	}
	return n
}

func newNode16[T any](perLen int) *artNode16[T] {
	n := &artNode16[T]{header: &header[T]{nKind: node16}}
	if perLen > 0 {
		n.prefix = make([]byte, perLen)
	}
	return n
}

/*
Node48:
As the number of entries in a node increases,
searching the key array becomes expensive. Therefore, nodes
with more than 16 pointers do not store the keys explicitly.
Instead, a 256-element array is used, which can be indexed
with key bytes directly. If a node has between 17 and 48 child
pointers, this array stores indexes into a second array which
contains up to 48 pointers. This indirection saves space in
comparison to 256 pointers of 8 bytes, because the indexes
only require 6 bits (we use 1 byte for simplicity).
*/
type artNode48[T any] struct {
	*header[T]
	idx   [capNode256]byte
	child [capNode48 + 1]node[T] // use [1....48]
}

func (n *artNode48[T]) find(b byte) (node[T], uint16) {
	if n.idx[b] != 0 {
		return n.child[n.idx[b]], uint16(n.idx[b])
	}
	return nil, 0
}

func (n *artNode48[T]) add(b byte, c node[T]) {
	if n.idx[b] != 0 { // b already exists.
		n.child[n.idx[b]] = c
	} else {
		n.childLen++
		n.idx[b] = byte(n.childLen)
		n.child[n.childLen] = c
	}
}

func (n *artNode48[T]) set(i uint16, ptr node[T]) {
	if i < capNode48 {
		n.child[i] = ptr
	}
}

func (n *artNode48[T]) del(b byte) {
	i := n.idx[b]
	if i != 0 && n.child[i] != nil {
		n.childLen--
		n.idx[b] = 0
		n.child[i] = nil
	}
}

func (n *artNode48[T]) nthChild(m uint16) node[T] {
	if m <= 0 || m > n.childLen {
		return nil
	}
	var cnt uint16
	for _, i := range n.idx {
		if i != 0 {
			cnt++
			if cnt == m {
				return n.child[i]
			}
		}
	}
	return nil
}

func (n *artNode48[T]) rank(b byte) (node[T], uint16) {
	idx := n.idx[b]
	if idx != 0 && n.child[idx] != nil {
		var r uint16
		for i := byte(0); i < b; i++ {
			if n.idx[i] != 0 {
				r++
			}
		}
		return n.child[idx], r + 1
	}
	return nil, 0
}

func (n *artNode48[T]) clone() ([]byte, []node[T]) {
	return n.idx[:capNode256], n.child[:capNode48+1]
}

func (n *artNode48[T]) shrink() node[T] {
	if n.childLen == 0 {
		return n.tail
	}
	if n.childLen <= capNode16 {
		n16 := newNode16[T](0)
		n16.header = n.header
		n16.header.nKind = node16
		var i int
		for k, j := range n.idx {
			if j != 0 && n.child[j] != nil {
				n16.keys[i] = byte(k)
				n16.child[i] = n.child[j]
				i++
			}
		}
		return n16
	}
	return n
}

func newNode48[T any](perLen int) *artNode48[T] {
	n := &artNode48[T]{header: &header[T]{nKind: node48}}
	if perLen > 0 {
		n.prefix = make([]byte, perLen)
	}
	return n
}

/*
Node256:
The largest node type is simply an array of 256
pointers and is used for storing between 49 and 256 entries.
With this representation, the next node can be found very
efficiently using a single lookup of the key byte in that array.
No additional indirection is necessary. If most entries are not
null, this representation is also very space efficient because
only pointers need to be stored.
*/
type artNode256[T any] struct {
	*header[T]
	child [capNode256]node[T]
}

func (n *artNode256[T]) find(b byte) (node[T], uint16) {
	return n.child[b], uint16(b)
}

func (n *artNode256[T]) add(b byte, c node[T]) {
	if n.child[b] != nil {
		n.child[b] = c
	} else {
		n.child[b] = c
		n.childLen++
	}
}

func (n *artNode256[T]) set(i uint16, ptr node[T]) {
	if i < capNode256 {
		n.child[i] = ptr
	}
}

func (n *artNode256[T]) del(b byte) {
	if n.child[b] != nil {
		n.childLen--
		n.child[b] = nil
	}
}

func (n *artNode256[T]) nthChild(m uint16) node[T] {
	if m <= 0 || m > n.childLen {
		return nil
	}
	var cnt uint16
	for _, ch := range n.child {
		if ch != nil {
			cnt++
			if cnt == m {
				return ch
			}
		}
	}
	return nil
}

func (n *artNode256[T]) rank(b byte) (node[T], uint16) {
	if n.child[b] != nil {
		var r uint16
		for i := byte(0); i < b; i++ {
			if n.child[i] != nil {
				r++
			}
		}
		return n.child[b], r + 1
	}
	return nil, 0
}

func (n *artNode256[T]) clone() ([]byte, []node[T]) {
	return nil, n.child[:capNode256]
}

func (n *artNode256[T]) shrink() node[T] {
	if n.childLen == 0 {
		return n.tail
	}
	if n.childLen <= capNode48 {
		n48 := newNode48[T](0)
		n48.header = n.header
		n48.header.nKind = node48
		i := byte(1)
		for k, ch := range n.child {
			if ch != nil {
				n48.idx[k] = i
				n48.child[i] = ch
				i++
			}
		}
		return n48
	}
	return n
}

/*
Besides storing paths using the inner nodes as discussed
in the previous section, radix trees must also store the values
associated with the keys. We assume that only unique keys
are stored, because non-unique indexes can be implemented
by appending the tuple identifier to each key as a tie-breaker.

The values can be stored in different ways:
* Single-value leaves: The values are stored using an
    additional leaf node type which stores one value.
* Multi-value leaves: ...
* Combined pointer/value slots: ...

Using single-value leaves is the most general method,
because it allows keys and values of varying length within one
tree.
*/
type dumyHeader[T any] struct{}

func (h dumyHeader[T]) prefixCut(int)                 {}
func (h dumyHeader[T]) setPrefix([]byte)              {}
func (h dumyHeader[T]) add(b byte, c node[T])         {}
func (h dumyHeader[T]) set(i uint16, c node[T])       {}
func (h dumyHeader[T]) del(byte)                      {}
func (h dumyHeader[T]) setEnd(node[T])                {}
func (h dumyHeader[T]) count() uint16                 { return 0 }
func (h dumyHeader[T]) prefixLen() int                { return 0 }
func (h dumyHeader[T]) prefixChar(int) byte           { return 0 }
func (h dumyHeader[T]) isFull() bool                  { return false }
func (h dumyHeader[T]) kind() nodeKind                { return leaf }
func (h dumyHeader[T]) find(b byte) (node[T], uint16) { return nil, 0 }
func (h dumyHeader[T]) clone() ([]byte, []node[T])    { return nil, nil }
func (h dumyHeader[T]) nthChild(uint16) node[T]       { return nil }
func (h dumyHeader[T]) end() node[T]                  { return nil }
func (h dumyHeader[T]) rank(b byte) (node[T], uint16) { return nil, 0 }

type artLeaf[T any] struct {
	dumyHeader[T]
	k Key
	v T
}

func (n *artLeaf[T]) key() Key        { return n.k }
func (n *artLeaf[T]) value() T        { return n.v }
func (n *artLeaf[T]) update(v T)      { n.v = v }
func (n *artLeaf[T]) shrink() node[T] { return n }

func newLeaf[T any](key Key, val T) *artLeaf[T] {
	lk := make([]byte, len(key))
	copy(lk, key)
	return &artLeaf[T]{dumyHeader: dumyHeader[T]{}, k: lk, v: val}
}

// merge the parent node's prefix to it's only one child ch.
func mergeTo[V any](prefix []byte, k byte, ch node[V]) node[V] {
	p := make([]byte, len(prefix)+ch.prefixLen()+1)
	copy(p, prefix)
	p[len(prefix)] = k
	copy(p[len(prefix)+1:], ch.key())
	ch.setPrefix(p)
	return ch
}

// checkPrefix compares the compressed path of a
// node with the key and returns the number of equal bytes
func checkPrefix[V any](n node[V], key Key, depth int) int {
	if n == nil {
		return 0
	}
	i, k := 0, n.key()
	for i < len(k) && i+depth < len(key) {
		if k[i] != key[i+depth] {
			break
		}
		i++
	}
	return i
}

// Longest common prefix matching.
func leafMatches[V any](n node[V], key Key, depth int) bool {
	if n == nil {
		return false
	}
	lk := n.key()
	if len(lk) != len(key) {
		return false
	}
	for i := depth; i < len(key); i++ {
		if lk[i] != key[i] {
			return false
		}
	}
	return true
}

func leafHasPerfix[V any](n node[V], key Key, depth int) bool {
	if n == nil {
		return false
	}
	lk := n.key()
	if len(lk) < len(key) {
		return false
	}
	for i := depth; i < len(key); i++ {
		if lk[i] != key[i] {
			return false
		}
	}
	return true
}

// grow method replaces a node n  by a larger node type.
func graw[V any](n node[V]) node[V] {
	switch n.kind() {
	case node4:
		nn := &artNode16[V]{header: &header[V]{nKind: node16}}
		nn.prefix = n.key()
		nn.preLen = n.prefixLen()
		nn.childLen = n.count()
		nn.tail = n.end()
		ks, ch := n.clone()
		copy(nn.keys[:len(ks)], ks)
		copy(nn.child[:len(ch)], ch)
		return nn
	case node16:
		nn := &artNode48[V]{header: &header[V]{nKind: node48}}
		nn.prefix = n.key()
		nn.preLen = n.prefixLen()
		nn.childLen = n.count()
		nn.tail = n.end()
		ks, ch := n.clone()
		for i, k := range ks {
			if ch[i] != nil {
				nn.idx[k] = byte(i + 1)
				nn.child[byte(i+1)] = ch[i]
			}
		}
		return nn
	case node48:
		nn := &artNode256[V]{header: &header[V]{nKind: node256}}
		nn.prefix = n.key()
		nn.preLen = n.prefixLen()
		nn.childLen = n.count()
		nn.tail = n.end()
		idx, chs := n.clone()
		for i, j := range idx {
			if chs[j] != nil {
				nn.child[i] = chs[j]
			}
		}
		return nn
	default:
		return n
	}
}
