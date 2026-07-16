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
	"sync"
)

var (
	errElemNotExists = errors.New("elements NotExists")
)

// Each non leaf node contains minDegree-1 to 2 * minDegree-1 key value pairs,
// arranged in ascending order.
type node[K, V any] struct {
	key    []K
	val    []V
	parent *node[K, V]
	child  []*node[K, V]
}

func (n *node[K, V]) isRoot() bool {
	return n.parent == nil
}

func (n *node[K, V]) isLeaf() bool {
	return n.child == nil
}

func (n *node[K, V]) len() int {
	return len(n.key)
}

func (n *node[K, V]) max() (k K, v V) {
	if n.len() > 0 {
		k, v = n.key[n.len()-1], n.val[n.len()-1]
	}
	return
}

func (n *node[K, V]) updateAt(idx int, k K, v V) {
	if idx >= 0 && idx < len(n.key) {
		n.key[idx] = k
		n.val[idx] = v
	}
}

func (n *node[K, V]) find(key K, comp func(K, K) int) (int, bool) {
	l, r := 0, len(n.key)
	for l < r {
		m := (l + r) / 2
		if comp(n.key[m], key) == 0 {
			return m, true
		} else if comp(n.key[m], key) < 0 { // key[m] < key
			l = m + 1
		} else {
			r = m // key[m] > key
		}
	}
	return l, false
}

func (n *node[K, V]) insert(k K, v V, p *node[K, V], comp func(K, K) int) int { // TODO
	i, ok := n.find(k, comp)
	if !ok {
		n.key, n.val, n.child = append(n.key, k), append(n.val, v), append(n.child, nil)
		copy(n.key[i+1:], n.key[i:])
		copy(n.val[i+1:], n.val[i:])
		copy(n.child[i+1:], n.child[i:])
	}
	n.key[i], n.val[i], n.child[i] = k, v, p
	return i
}

func (n *node[K, V]) insertAt(i int, k K, v V) {
	n.key, n.val = append(n.key, k), append(n.val, v)
	copy(n.key[i+1:], n.key[i:])
	copy(n.val[i+1:], n.val[i:])
	n.key[i], n.val[i] = k, v
}

func (n *node[K, V]) deleteAt(i int) {
	if i >= 0 && i < len(n.key) {
		copy(n.key[i:], n.key[i+1:])
		copy(n.val[i:], n.val[i+1:])
		n.key, n.val = n.key[:len(n.key)-1], n.val[:len(n.val)-1]
		if !n.isLeaf() {
			copy(n.child[i:], n.child[i+1:])
			n.child = n.child[:len(n.child)-1]
		}
	}
}

// Cursors are used to access ordered collections.
type Cursor[K, V any] interface {
	// If the collection object is in concurrent security mode,
	// the Open method needs to be called to attempt locking before using the cursor.
	// After use, the Close method must be called to release the lock.
	Open() error
	Close()
	// The Seek(key) method locates the cursor at the key.
	// If the key does not exist, it locates at the next key and returns it
	Seek(K) (K, V, bool)
	// The First method locates the cursor at the minimum element of the set.
	// If there is no minimum element (the set is empty), it returns false
	First() (K, V, bool)
	// The Last method locates the cursor at the maximum element of the set.
	// If there is no maximum element (the set is empty), it returns false.
	Last() (K, V, bool)
	// HasNext returns whether the next element exists relative to the current cursor position.
	HasNext() bool
	// Next() moves the cursor backwards and returns the element.
	// If the element does not exist, it returns a type zero value.
	Next() (K, V)
	// HasPrev() returns whether the previous element exists relative to the current cursor position.
	HasPrev() bool
	// Prev() moves the cursor forward and returns the element.
	// If the element does not exist, it returns a type value of zero.
	Prev() (K, V)
}

// BTree is an ordered collection of key value pairs in memory,
// structurally a multi-path balanced tree.
// Support operations such as adding, deleting, modifying, and querying.
type BTree[K any, V any] struct {
	lock   bool
	degree int
	comp   func(K, K) int

	mu    sync.RWMutex
	count int
	high  int
	root  *node[K, V]
	cur   *btreeCursor[K, V]
}

type BTreeConfig struct {
	// If the Lock field is set to true, the created BTree is concurrency safe.
	Lock bool
	// The MinDegree field is used to set the minimum outdegree (minimum number of subtrees) of internal nodes in the BTree,
	// where each node will contain MinDegree-1 to 2 * MinDegree-1 elements
	MinDegree int
}

const (
	DefaultBTreeLock      = false
	DefaultBTreeMinDegree = 32
)

/*
type BTreeOptions[K, V any] func(*BTree[K, V])
func CompareOpts[K, V any](comp func(K, K) int) BTreeOptions[K, V] {
	return func(bt *BTree[K, V]) {
		bt.comp = comp
	}
}
func CapOpts[K, V any](min, max int) BTreeOptions[K, V] {
	return func(bt *BTree[K, V]) {
		bt.minCap, bt.maxCap = min, max
	}
}
func LockOpts[K, V any](noLock bool) BTreeOptions[K, V] {
	return func(bt *BTree[K, V]) {
		bt.noLock = noLock
	}
}

func NewBTree[K, V any](ops ...BTreeOptions[K, V]) *BTree[K, V] {
	return &BTree[K, V]{}
}
*/

// NewBTree creates a btree.
// The comp parameter is used to specify the key comparison function.
// If the return value of comp (k1, k2) is less than 0, it means k1<k2;
// A return value of 0 indicates that k1=k2; A return value greater than 0 indicates that k1>k2。
func NewBTree[K, V any](cfg *BTreeConfig, comp func(K, K) int) *BTree[K, V] {
	bt := &BTree[K, V]{
		degree: cfg.MinDegree,
		lock:   cfg.Lock,
		comp:   comp,
	}
	if bt.degree < 2 {
		bt.degree = DefaultBTreeMinDegree
	}
	bt.cur = &btreeCursor[K, V]{tree: bt}
	return bt
}

// The number of elements in the current BTree.
func (bt *BTree[K, V]) Len() int {
	if bt.lock {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
	}
	return bt.count
}

func (bt *BTree[K, V]) Less(k1, k2 K) bool {
	return bt.comp(k1, k2) < 0
}

// The current height of BTree.
func (bt *BTree[K, V]) High() int {
	if bt.lock {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
	}
	return bt.high
}

// Return the configuration information of BTree.
func (bt *BTree[K, V]) Options() BTreeConfig {
	return BTreeConfig{
		Lock:      bt.lock,
		MinDegree: bt.degree,
	}
}

// Query the specified element, if it does not exist, return Not Exists error.
func (bt *BTree[K, V]) Search(key K) (V, error) {
	if bt.lock {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
	}
	_, v, ok := bt.cur.Seek(key)
	if !ok {
		return v, errElemNotExists
	}
	return v, nil
}

// Query the specified element, if it does not exist, return next.
func (bt *BTree[K, V]) Seek(key K) (K, V, bool) {
	if bt.lock {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
	}
	return bt.cur.Seek(key)
}

// Query the minimum element, if it does not exist, return NotExists error.
func (bt *BTree[K, V]) First() (K, V, error) {
	if bt.lock {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
	}
	k, v, ok := bt.cur.First()
	if ok {
		return k, v, nil
	}
	return k, v, errElemNotExists
}

// Query the maximum element, if it does not exist, return NotExists error.
func (bt *BTree[K, V]) Last() (K, V, error) {
	if bt.lock {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
	}
	k, v, ok := bt.cur.Last()
	if ok {
		return k, v, nil
	}
	return k, v, errElemNotExists
}

func (bt *BTree[K, V]) Index(n int) (k K, v V, err error) {
	if n < 0 || n >= bt.Len() {
		err = errors.New("index out of range")
		return
	}
	if n > bt.Len()/2 {
		k, v, err = bt.Last()
		if err != nil {
			return
		}
		for i := 1; i < bt.Len()-n; i++ {
			k, v = bt.cur.Prev()
		}
	} else {
		k, v, err = bt.First()
		if err != nil {
			return
		}
		for i := 1; i < n; i++ {
			k, v = bt.cur.Next()
		}
	}
	return
}

// Query the elements within the range of [left, right].
// If there are no elements within this range, empty slices will be returned.
func (bt *BTree[K, V]) Range(left K, right K) ([]K, []V) {
	if bt.lock {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
	}
	var resK []K
	var resV []V
	k, v, ok := bt.cur.Seek(left)
	if ok {
		resK = append(resK, k)
		resV = append(resV, v)
	}
	for bt.cur.HasNext() {
		k, v = bt.cur.Next()
		if bt.comp(k, right) >= 0 {
			break
		}
		resK = append(resK, k)
		resV = append(resV, v)
	}
	return resK, resV
}

// Scan all elements in ascending order and use parameter functions op to process key value pairs sequentially.
// Once the processing function returns an error, the scan terminates.
func (bt *BTree[K, V]) Scan(desc bool, op func(K, V) error) error {
	if bt.lock {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
	}
	if desc {
		k, v, ok := bt.cur.Last()
		if !ok {
			return errElemNotExists
		}
		err := op(k, v)
		if err != nil {
			return err
		}
		for bt.cur.HasPrev() {
			k, v = bt.cur.Prev()
			if err = op(k, v); err != nil {
				return err
			}
		}
		return nil
	}
	k, v, ok := bt.cur.First()
	if !ok {
		return errElemNotExists
	}
	err := op(k, v)
	if err != nil {
		return err
	}
	for bt.cur.HasNext() {
		k, v = bt.cur.Next()
		if err = op(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Insert elements into the current BTree,
// and if the corresponding key already exists,
// overwrite the original value with the parameter value
func (bt *BTree[K, V]) Insert(key K, value V) {
	if bt.lock {
		bt.mu.Lock()
		defer bt.mu.Unlock()
	}
	_, _, ok := bt.cur.Seek(key)
	nd, i := bt.cur.currentNode()
	if ok {
		// just update
		nd.val[i] = value
		return
	}
	if nd == nil { // new root
		bt.root = &node[K, V]{}
		nd = bt.root
		bt.high++
		i = 0
	}
	nd.insertAt(max(0, i), key, value) // insert always into leaf.
	bt.count++
	bt.split(nd)
}

// Delete element, return true if successful, return NotExists error if key does not exist.
func (bt *BTree[K, V]) Delete(key K) (bool, error) {
	if bt.lock {
		bt.mu.Lock()
		defer bt.mu.Unlock()
	}
	if _, _, ok := bt.cur.Seek(key); !ok {
		return false, errElemNotExists
	}
	nd, i := bt.cur.currentNode()
	if nd.isLeaf() {
		nd.deleteAt(i)
	} else {
		// find next or prev element,replace current key by it,then delete on leaf.
		if bt.cur.HasNext() {
			_, _ = bt.cur.Next()
		} else {
			_, _ = bt.cur.Prev()
		}
		leaf, j := bt.cur.currentNode()
		if leaf == nil || !leaf.isLeaf() {
			return false, errors.New("delete failed")
		}
		nd.updateAt(i, leaf.key[j], leaf.val[j])
		leaf.deleteAt(j)
		nd = leaf
	}
	bt.count--
	bt.rebalance(nd)
	return true, nil
}

// Delete all elements, BTree becomes an empty tree.
func (bt *BTree[K, V]) Clean() {
	if bt.lock {
		bt.mu.Lock()
		defer bt.mu.Unlock()
	}
	bt.cur = nil
	bt.root = nil
	bt.count, bt.high = 0, 0
	bt.cur = &btreeCursor[K, V]{tree: bt}
}

// Create a cursor for the current BTree.
func (bt *BTree[K, V]) Cursor() Cursor[K, V] {
	return &btreeCursor[K, V]{tree: bt}
}

// Deeply copy a current btree.
func (bt *BTree[K, V]) Clone() *BTree[K, V] {
	if bt.lock {
		bt.mu.RLock()
		defer bt.mu.RUnlock()
	}
	t := &BTree[K, V]{
		lock:   bt.lock,
		degree: bt.degree,
		comp:   bt.comp,
		high:   bt.high,
		count:  bt.count,
	}
	t.root = bt.clone(bt.root)
	return t
}

func (bt *BTree[K, V]) clone(nd *node[K, V]) *node[K, V] {
	if nd == nil {
		return nil
	}
	n := bt.newNode(nd.isLeaf())
	n.key, n.val = make([]K, nd.len()), make([]V, nd.len())
	copy(n.key, nd.key)
	copy(n.val, nd.val)
	if !n.isLeaf() {
		n.child = make([]*node[K, V], nd.len()+1)
		copy(n.child, nd.child)
	}
	for i, p := range n.child {
		n.child[i] = bt.clone(p)
		if n.child[i] != nil {
			n.child[i].parent = n
		}
	}
	return n
}

func (bt *BTree[K, V]) newNode(leaf bool) *node[K, V] {
	n := &node[K, V]{}
	if !leaf {
		n.child = make([]*node[K, V], 0)
	}
	return n
}

func (bt *BTree[K, V]) split(nd *node[K, V]) {
	if nd == nil || nd.len() < 2*bt.degree-1 {
		return
	}
	// split current node two nodes.
	i := bt.degree
	k, v := nd.key[i], nd.val[i]
	// new node
	newNd := bt.newNode(nd.isLeaf())
	newNd.parent = nd.parent
	newNd.key, newNd.val = nd.key[i+1:], nd.val[i+1:]
	if !nd.isLeaf() {
		newNd.child = nd.child[i+1:]
	}
	// update parent ptr for new node's child
	for _, ch := range newNd.child {
		if ch != nil {
			ch.parent = newNd
		}
	}
	// old node
	nd.key, nd.val = nd.key[:i:i], nd.val[:i:i]
	if !nd.isLeaf() {
		nd.child = nd.child[: i+1 : i+1]
	}
	// if current node is root, creat a new root.
	if nd.isRoot() {
		r := bt.newNode(false)
		r.key = append(r.key, k)
		r.val = append(r.val, v)
		r.child = append(r.child, nd, newNd)
		nd.parent, newNd.parent = r, r
		bt.high++
		bt.root = r
	} else {
		// insert (k,v,nd) to nd.parent
		j := nd.parent.insert(k, v, nd, bt.comp)
		nd.parent.child[j+1] = newNd
	}
	bt.split(nd.parent)
}

func (bt *BTree[K, V]) rebalance(nd *node[K, V]) {
	if nd == nil || nd.len() >= bt.degree-1 {
		return
	}
	if nd.isRoot() {
		if nd.len() == 0 {
			// drop current root, and lift its child as new root
			for _, p := range nd.child {
				if p != nil {
					bt.root = p
					bt.root.parent = nil
					bt.high--
					return
				}
			}
			bt.root = nil
			bt.high, bt.count = 0, 0
		}
		return
	}
	// find left or right neighbour node.
	i := -1
	for j, p := range nd.parent.child {
		if p == nd {
			i = j
			break
		}
	}
	if i < 0 {
		panic("btree construction failure")
	}
	var l, r *node[K, V]
	if i > 0 {
		l = nd.parent.child[i-1]
	}
	if i < nd.parent.len() {
		r = nd.parent.child[i+1]
	}
	// try to merge nodes.
	if l != nil && l.len()+nd.len() < 2*bt.degree-1 {
		// merge current node to left sibling:
		// get k-v=(key[i-1],val[i-1]) from l.parent(=nd.parent),
		// put k-v to tail of l, put nd on the tail of l, update l.parent.child[i] = l
		// delete k-v from l.parent.
		l.key, l.val = append(l.key, l.parent.key[i-1]), append(l.val, l.parent.val[i-1])
		for _, p := range nd.child {
			if p != nil {
				p.parent = l
			}
		}
		l.key, l.val = append(l.key, nd.key...), append(l.val, nd.val...)
		if !l.isLeaf() {
			l.child = append(l.child, nd.child...)
		}
		l.parent.child[i] = l
		l.parent.deleteAt(i - 1)
		bt.rebalance(l.parent)
	} else if r != nil && r.len()+nd.len() < 2*bt.degree-1 {
		// merge right sibling to current node.
		nd.key, nd.val = append(nd.key, nd.parent.key[i]), append(nd.val, nd.parent.val[i])
		for _, p := range r.child {
			if p != nil {
				p.parent = nd
			}
		}
		nd.key, nd.val = append(nd.key, r.key...), append(nd.val, r.val...)
		if !nd.isLeaf() {
			nd.child = append(nd.child, r.child...)
		}
		nd.parent.child[i+1] = nd
		nd.parent.deleteAt(i)
		bt.rebalance(nd.parent)
	} else { // >= 2*bt.drgree-1
		if l != nil {
			// steal some elements from left sibling:
			// get k-v=(key[i-1],val[i-1]) from l.parent, put it to tail of l,
			// cut (degree-1-nd.len()) elements from tail of l, put these to head of nd,
			// put the max key of l to l.parent.key[i-1].
			k, v := l.parent.key[i-1], l.parent.val[i-1]
			l.key, l.val = append(l.key, k), append(l.val, v)
			j := l.len() + nd.len() + 1 - bt.degree // l.len()-(bt.degree-1-nd.len())

			ks, vs := l.key[j:], l.val[j:]
			ks, vs = append(ks, nd.key...), append(vs, nd.val...)
			nd.key, nd.val = ks, vs
			if !l.isLeaf() {
				ch := l.child[j:]
				ch = append(ch, nd.child...)
				nd.child = ch
			}

			l.parent.key[i-1], l.parent.val[i-1] = l.key[j-1], l.val[j-1]
			l.key, l.val = l.key[:j-1:j-1], l.val[:j-1:j-1]
			if !l.isLeaf() {
				l.child = l.child[:j:j]
			}

			for _, p := range nd.child {
				if p != nil {
					p.parent = nd
				}
			}
		} else if r != nil {
			// steal some elements from right sibling:
			// get k-v=(key[i],val[i]) from nd.parent, put it to tail of nd,
			// cut (degree-nd.len()) elements from head of r,put it to tail of nd,
			// put the max key of nd to r.parent.key[i].
			k, v := nd.parent.key[i], nd.parent.val[i]
			nd.key, nd.val = append(nd.key, k), append(nd.val, v)
			j := bt.degree - nd.len()
			nd.key, nd.val = append(nd.key, r.key[:j]...), append(nd.val, r.val[:j]...)
			if !nd.isLeaf() {
				nd.child = append(nd.child, r.child[:j]...)
			}

			r.key, r.val = r.key[j:], r.val[j:]
			if !r.isLeaf() {
				r.child = r.child[j:]
			}

			nd.parent.key[i], nd.parent.val[i] = nd.max()
			nd.key, nd.val = nd.key[:nd.len()-1], nd.val[:nd.len()-1]

			for _, p := range nd.child {
				if p != nil {
					p.parent = nd
				}
			}
		} else {
			p := bt.mergeWithParent(i, nd)
			if p != nil && p.len() < bt.degree-1 {
				bt.rebalance(p)
			} else {
				bt.split(p)
			}
		}
	}
}

func (bt *BTree[K, V]) mergeWithParent(i int, nd *node[K, V]) *node[K, V] {
	p := nd.parent
	if p == nil {
		return nil
	}
	if i >= p.len() {
		p.key, p.val = append(p.key, nd.key...), append(p.val, nd.val...)
		if nd.isLeaf() {
			p.child = nil
		} else {
			p.child = p.child[:len(p.child)-1]
			p.child = append(p.child, nd.child...)
		}
	} else {
		p.child[i] = nil
		nd.key, nd.val = append(nd.key, p.key[i]), append(nd.val, p.val[i])
		ks, vs := p.key[i+1:], p.val[i+1:]
		p.key, p.val = p.key[:i:i], p.val[:i:i]
		p.key, p.val = append(p.key, nd.key...), append(p.val, nd.val...)
		p.key, p.val = append(p.key, ks...), append(p.val, vs...)
		if nd.isLeaf() {
			p.child = nil
		} else {
			cs := p.child[i+1:]
			p.child = p.child[:i:i]
			p.child = append(p.child, nd.child...)
			p.child = append(p.child, cs...)
		}
	}
	for _, q := range p.child {
		if q != nil {
			q.parent = p
		}
	}
	return p
}

type nodeIdx[K, V any] struct {
	node    *node[K, V]
	idx     int  // idx record current search location on the node.
	visited bool // record current location has been visited.
}

type btreeCursor[K, V any] struct {
	mu    sync.Mutex
	opend bool
	tree  *BTree[K, V]
	path  []*nodeIdx[K, V] // a stack to record search path.
	top   int
}

func (c *btreeCursor[K, V]) Open() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.opend {
		return errors.New("current cursor already opened")
	}
	if c.tree.lock {
		c.tree.mu.RLock()
	}
	c.opend = true
	return nil
}

func (c *btreeCursor[K, V]) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.opend {
		if c.tree.lock {
			c.tree.mu.RUnlock()
		}
	}
	c.opend = false
}

// stack operations.
func (c *btreeCursor[K, V]) pClean() {
	c.path = []*nodeIdx[K, V]{}
	c.top = 0
}
func (c *btreeCursor[K, V]) pAppend(idx *nodeIdx[K, V]) {
	if c.top < len(c.path) {
		c.path[c.top] = idx
	} else {
		c.path = append(c.path, idx)
	}
	c.top++
}
func (c *btreeCursor[K, V]) pRemove() *nodeIdx[K, V] {
	if c.top > 0 {
		v := c.path[c.top-1]
		c.top--
		return v
	}
	return nil
}
func (c *btreeCursor[K, V]) pTail() *nodeIdx[K, V] {
	if c.top > 0 {
		return c.path[c.top-1]
	}
	return nil
}

func (c *btreeCursor[K, V]) currentNode() (*node[K, V], int) {
	p := c.pTail()
	if p != nil {
		return p.node, p.idx
	}
	return nil, 0
}

func (c *btreeCursor[K, V]) Seek(key K) (k K, v V, false bool) {
	if c.tree.lock && !c.opend {
		return
	}
	c.pClean()
	for p := c.tree.root; p != nil; {
		i, ok := p.find(key, c.tree.comp)
		ni := &nodeIdx[K, V]{node: p, idx: i}
		c.pAppend(ni)
		if ok {
			ni.visited = true
			return key, p.val[i], true
		}
		if p.isLeaf() {
			if i < p.len() {
				ni.visited = true
				return p.key[i], p.val[i], false
			}
			return
		}
		p = p.child[i]
	}
	return
}

func (c *btreeCursor[K, V]) First() (k K, v V, false bool) {
	if c.tree.lock && !c.opend {
		return
	}
	c.pClean()
	for p := c.tree.root; p != nil; p = p.child[0] {
		c.pAppend(&nodeIdx[K, V]{node: p})
		if p.isLeaf() {
			break
		}
	}
	if p := c.pTail(); p != nil && p.node.len() > 0 {
		p.visited = true
		return p.node.key[0], p.node.val[0], true
	}
	return
}

func (c *btreeCursor[K, V]) HasNext() bool {
	if c.tree.lock && !c.opend {
		return false
	}
	for p := c.pTail(); p != nil; p = c.pTail() {
		if p.node.isLeaf() {
			if p.idx+1 < p.node.len() { // for leaf node,p.idx has been visited.
				return true
			}
		} else {
			// for non-leaf node,p.idx element not been visited.
			if p.idx < p.node.len() {
				return true
			}
		}
		_ = c.pRemove()
	}
	return false
}

func (c *btreeCursor[K, V]) Next() (k K, v V) {
	if c.tree.lock && !c.opend {
		return
	}
	for p := c.pTail(); p != nil; p = c.pTail() {
		if p.idx < p.node.len() {
			if p.node.isLeaf() {
				// try to visited next element
				p.idx++
				if p.idx < p.node.len() {
					return p.node.key[p.idx], p.node.val[p.idx]
				} else {
					_ = c.pRemove()
				}
			} else {
				// p.idx has been visited,so should move cursor to next element, next subtree root,
				// and push its left subtree in stack.
				if p.visited {
					p.idx++
					p.visited = false
					for q := p.node.child[p.idx]; q != nil; q = q.child[0] {
						ni := &nodeIdx[K, V]{node: q}
						c.pAppend(ni)
						if q.isLeaf() {
							ni.idx--
							break
						}
					}
				} else {
					p.visited = true
					return p.node.key[p.idx], p.node.val[p.idx]
				}
			}
		} else {
			_ = c.pRemove()
		}
	}
	return
}

func (c *btreeCursor[K, V]) Last() (k K, v V, false bool) {
	if c.tree.lock && !c.opend {
		return
	}
	c.pClean()
	for p := c.tree.root; p != nil; p = p.child[p.len()] {
		c.pAppend(&nodeIdx[K, V]{node: p, idx: p.len()})
		if p.isLeaf() {
			break
		}
	}
	if p := c.pTail(); p != nil && p.node.len() > 0 {
		p.idx--
		p.visited = true
		return p.node.key[p.idx], p.node.val[p.idx], true
	}
	return
}

func (c *btreeCursor[K, V]) HasPrev() bool {
	if c.tree.lock && !c.opend {
		return false
	}
	for p := c.pTail(); p != nil; p = c.pTail() {
		if p.node.isLeaf() {
			if p.idx >= 1 { // for leaf node,p.idx has benn visited.
				return true
			}
		} else {
			// if p.visited we just check it has child.
			// if p.idx not been visited,we should check it has left element.
			if (p.idx >= 0 && p.visited) || p.idx >= 1 {
				return true
			}
		}
		_ = c.pRemove()
	}
	return false
}

func (c *btreeCursor[K, V]) Prev() (k K, v V) {
	if c.tree.lock && !c.opend {
		return
	}
	for p := c.pTail(); p != nil; p = c.pTail() {
		if p.idx >= 0 {
			if p.node.isLeaf() {
				p.idx--
				if p.idx >= 0 {
					return p.node.key[p.idx], p.node.val[p.idx]
				} else {
					_ = c.pRemove()
				}
			} else {
				if p.visited {
					// if p.idx has been visited, then try to visited its subtree
					p.visited = false // set false,then next meeting it,we just ignore.
					for q := p.node.child[p.idx]; q != nil; q = q.child[q.len()] {
						c.pAppend(&nodeIdx[K, V]{node: q, idx: q.len()})
						if q.isLeaf() {
							break
						}
					}
				} else {
					// we need not visited it,just move cursor to prev element.
					p.idx--
					if p.idx >= 0 {
						p.visited = true
						return p.node.key[p.idx], p.node.val[p.idx]
					} else {
						_ = c.pRemove()
					}
				}
			}
		} else {
			_ = c.pRemove()
		}
	}
	return
}
