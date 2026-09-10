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

import (
	"sync"

	"github.com/flxj/graphlib/collection"
)

type (
	// The Adaptive Radix Tree uses []byte as the key for elements.
	// The keys in the tree will be sorted by byte in lexicographic order.
	Key []byte
	// The data access operator is called sequentially
	// when traversing elements in the tree.
	Visitor[V any] func(Key, V) error
	// node type: node4,node16,node48,node256,leaf
	nodeKind uint8
)

/*
Adaptive cardinality tree (ART) is a data structure designed
for in memory databases, aimed at addressing the spatial efficiency
and search performance issues of traditional Trie trees.
ART provides O (k) query complexity through adaptive nodes,
high compression (path compression and lazy extension),
and special internal and leaf node structures
*/
type Tree[V any] interface {
	// Return the number of elements in the current tree.
	// Note that the current version of ART does not support duplicate keys.
	Len() int
	// Insert element, update its value in place if the key already exists.
	Insert(Key, V)
	// Query elements, return element value if key exists and true flag,
	// return false flag if key does not exist
	Search(Key) (V, bool)
	// Delete element.
	// If the key exists, delete it and return the element value and true flag.
	// If the key does not exist, return false flag.
	Delete(Key) (V, bool)
	// Scan all elements in order, and call visitor during the scanning process.
	// If error is returned when calling, the scanning process will be
	// terminated immediately and the corresponding error will be returned.
	// Note that users can determine whether to scan in ascending or
	// descending order using the desc parameter.
	Scan(fn Visitor[V], desc bool) error
	// Scan all elements with the specified prefix in sequence,
	// and call visitor during the scanning process.
	// If error is returned when calling, the scanning process will be
	// terminated immediately and the corresponding error will be returned.
	// Note that users can determine whether to scan in ascending or
	// descending order using the desc parameter.
	PrefixScan(perfix Key, fn Visitor[V], desc bool) error
	// Update an existing element, the user can specify an update function
	// that receives its old value, and updates the result to the corresponding key.
	Update(Key, func(V) V) bool
	// Return the element with the smallest lexicographic order in the tree.
	// If it does not exist, a false flag will be returned.
	Min() (Key, V, bool)
	// Return the element with the highest lexicographic order in the tree.
	// If it does not exist, a false flag will be returned.
	Max() (Key, V, bool)
	// Get the iterator of the tree.
	Cursor() collection.Cursor[Key, V]
}

// Create a new tree, and if you want a concurrency safe tree,
// you need to set the syncSafety parameter to true,
// which will use read-write locks to protect the tree structure.
func NewTree[V any](syncSafety bool) Tree[V] {
	return &art[V]{locked: syncSafety}
}

// art has implemented the Tree interface.
type art[V any] struct {
	locked bool
	mu     sync.RWMutex
	cnt    int
	root   node[V]
	oldV   V
}

func (t *art[V]) Len() int {
	if t.locked {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	return t.cnt
}

/*
The pseudo code is shown in Figure 9. The tree is traversed using the recursive call in line 29,
until the position for the new leaf is found. Usually, the leaf can simply be inserted into an
existing inner node, after growing it if necessary (lines 31-33).
If, because of lazy expansion,an existing leaf is encountered, it is replaced by a new
inner node storing the existing and the new leaf (lines 5-13).
Another special case occurs if the key of the new leaf differs from a compressed path:
A new inner node is created above the current node and the compressed paths are adjusted
accordingly (lines 17-24).
*/
func (t *art[V]) insert(cur node[V], key Key, val V, depth int) node[V] {
	if cur == nil {
		t.cnt++
		return newLeaf(key, val)
	}
	if cur.kind() == leaf {
		old := cur.key()
		i := depth
		for ; i < len(old) && i < len(key) && key[i] == old[i]; i++ {
		}
		// The leaf node element matches the key exactly,
		// means key already exists,so just update cur's value.
		if i == len(old) && i == len(key) {
			cur.update(val)
			return cur
		}
		// lazy expansion, inner nodes are only created
		// if they are required to distinguish at least two leaf nodes.

		// Create a new inner node of type node4,
		// with cur and the newly inserted element as its two leaf nodes
		// or endpoint nodes, respectively.
		newNode := newNode4[V](i - depth)
		for j := depth; j < i; j++ {
			newNode.prefix[j-depth] = key[j]
		}
		newNode.preLen = i - depth
		depth = i
		if depth >= len(key) && depth >= len(old) {
			panic("invalid create new node for same key")
		}
		if depth < len(key) {
			newNode.add(key[depth], newLeaf(key, val))
		} else {
			newNode.setEnd(newLeaf(key, val))
		}
		if depth < len(old) {
			newNode.add(old[depth], cur)
		} else {
			newNode.setEnd(cur)
		}
		t.cnt++
		return newNode
	}
	// cur is a inner node. check prefix to find the branch point.
	p := checkPrefix(cur, key, depth)
	if p != cur.prefixLen() {
		// The newly inserted key does not match the perfix of the current node,
		// so the tree needs to branch here. Therefore, a new node of type node4
		// needs to be created, and the current node needs to be made its child.
		newNode := newNode4[V](p)
		depth += p
		if depth < len(key) {
			newNode.add(key[depth], newLeaf(key, val))
		} else {
			// key already exists.
			newNode.setEnd(newLeaf(key, val))
		}
		copy(newNode.prefix, cur.key()[:p])
		newNode.preLen = p
		newNode.add(cur.prefixChar(p), cur)
		cur.prefixCut(p + 1)
		t.cnt++
		return newNode
	}
	// key contains cur.prefix.
	depth += cur.prefixLen()
	if depth < len(key) {
		// Enter the next level node to continue matching and inserting.
		if child, i := cur.find(key[depth]); child != nil {
			newChild := t.insert(child, key, val, depth+1)
			if newChild != child {
				cur.set(i, newChild)
			}
		} else {
			// should insert a new leaf,
			// before insert if cur already full,should graw first.
			if cur.isFull() {
				cur = graw(cur)
			}
			t.cnt++
			cur.add(key[depth], newLeaf(key, val))
		}
	} else {
		if term := cur.end(); term != nil {
			term.update(val)
		} else {
			t.cnt++
			cur.setEnd(newLeaf(key, val))
		}
	}
	return cur
}

func (t *art[V]) Insert(key Key, val V) {
	if t.locked {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	t.root = t.insert(t.root, key, val, 0)
}

func (t *art[V]) find(n node[V], key Key, depth int) node[V] {
	if n == nil {
		return nil
	}
	if n.kind() == leaf {
		if leafMatches(n, key, depth) {
			return n
		}
		return nil
	}
	if checkPrefix(n, key, depth) != n.prefixLen() {
		return nil
	}
	depth += n.prefixLen()
	if depth < len(key) {
		next, _ := n.find(key[depth])
		return t.find(next, key, depth+1)
	} else {
		return t.find(n.end(), key, depth)
	}
}

/*
search (node, key, depth)

	1  if node==NULL
	2    return NULL
	3  if isLeaf(node)
	4      if leafMatches(node, key, depth)
	5          return node
	6      return NULL
	7  if checkPrefix(node,key,depth)!=node.prefixLen
	8      return NULL
	9  depth = depth + node.prefixLen
	10 next = findChild(node, key[depth])
	11 return search(next, key, depth+1)

Pessimistic path compression is handled in lines 7 and 8 by aborting the
search if the compressed path does not match the key.
*/
func (t *art[V]) search(key Key) node[V] {
	depth := 0
	for cur := t.root; cur != nil; depth++ {
		if cur.kind() == leaf {
			if leafMatches(cur, key, depth) {
				return cur
			}
			break
		} else {
			p := checkPrefix(cur, key, depth)
			if p != cur.prefixLen() { // which means cur.prefix not nil.
				break
			}
			depth += p
			if depth < len(key) {
				cur, _ = cur.find(key[depth])
			} else {
				// check if cur's end node match with key.
				if leafMatches(cur.end(), key, depth) {
					return cur.end()
				}
				break
			}
		}
	}
	return nil
}

func (t *art[V]) Search(key Key) (v V, ok bool) {
	if t.locked {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	n := t.search(key)
	if n != nil {
		v, ok = n.value(), true
	}
	return
}

func (t *art[V]) Update(key Key, fn func(V) V) bool {
	if t.locked {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	n := t.search(key)
	if n != nil {
		n.update(fn(n.value()))
		return true
	}
	return false
}

/*
The implementation of deletion is symmetrical to
insertion. The leaf is removed from an inner node, which is
shrunk if necessary. If that node now has only one child, it is
replaced by its child and the compressed path is adjusted
*/
func (t *art[V]) delete(cur node[V], key Key, depth int) node[V] {
	if cur == nil {
		return nil
	}
	if cur.kind() == leaf {
		if leafMatches(cur, key, depth) {
			t.oldV = cur.value()
			t.cnt--
			return nil
		}
		return cur
	}
	if checkPrefix(cur, key, depth) != cur.prefixLen() {
		return cur
	}
	var child node[V]
	var i uint16
	depth += cur.prefixLen()
	if depth < len(key) {
		child, i = cur.find(key[depth])
		if child == nil {
			return cur
		}
	} else {
		term := cur.end()
		if !leafMatches(term, key, depth) {
			return cur
		}
		t.oldV = term.value()
		cur.setEnd(nil)
		t.cnt--
		return cur.shrink()
	}
	newChild := t.delete(child, key, depth+1)
	if newChild == child {
		return cur
	}
	if newChild == nil {
		// The subtree is empty,
		// and the corresponding branch information
		// saved by the current node needs to be deleted.
		cur.del(key[depth])
	} else {
		// The update of the subtree root node requires
		// resetting the corresponding branch information
		// saved by the current node.
		cur.set(i, newChild)
	}
	return cur.shrink()
}

func (t *art[V]) Delete(key Key) (V, bool) {
	if t.locked {
		t.mu.Lock()
		defer t.mu.Unlock()
	}
	c := t.cnt
	t.root = t.delete(t.root, key, 0)
	return t.oldV, c != t.cnt
}

func (t *art[V]) scan(c collection.Cursor[Key, V], fn Visitor[V], desc bool) error {
	if desc {
		k, v, ok := c.Last()
		if !ok {
			return nil
		}
		if err := fn(k, v); err != nil {
			return err
		}
		for c.HasPrev() {
			k, v := c.Prev()
			if err := fn(k, v); err != nil {
				return err
			}
		}
	} else {
		k, v, ok := c.First()
		if !ok {
			return nil
		}
		if err := fn(k, v); err != nil {
			return err
		}
		for c.HasNext() {
			k, v := c.Next()
			if err := fn(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *art[V]) Scan(fn Visitor[V], desc bool) error {
	if t.locked {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	c := newCursor(t, t.root)
	return t.scan(c, fn, desc)
}

func (t *art[V]) PrefixScan(prefix Key, fn Visitor[V], desc bool) error {
	if t.locked {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	// 1.find the subtree root node of prefix
	n := t.root
	for d := 0; n != nil && d < len(prefix); d++ {
		if n.kind() == leaf {
			if leafHasPerfix(n, prefix, d) {
				break
			}
			return nil
		} else {
			p := checkPrefix(n, prefix, d)
			q := min(n.prefixLen(), len(prefix)-d)
			if p != q {
				return nil
			}
			d += q
			if d >= len(prefix) {
				break
			}
			n, _ = n.find(prefix[d])
		}
	}
	// 2. create a cursor
	if n != nil {
		c := newCursor(t, n)
		return t.scan(c, fn, desc)
	}
	return nil
}

func (t *art[V]) Min() (k Key, v V, ok bool) {
	if t.locked {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	c := newCursor(t, t.root)
	return c.First()
}

func (t *art[V]) Max() (k Key, v V, ok bool) {
	if t.locked {
		t.mu.RLock()
		defer t.mu.RUnlock()
	}
	c := newCursor(t, t.root)
	return c.Last()
}

func (t *art[V]) Cursor() collection.Cursor[Key, V] {
	return newCursor(t, t.root)
}
