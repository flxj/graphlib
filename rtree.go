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
	"math"
	"sync"
)

var (
	DefaultRTreeM = 64
)

// [4]N{minX, minY, maxX, maxY}
type Rectangle[N number] [4]N

func (t Rectangle[N]) Intersect(r Rectangle[N]) bool {
	return !(r[0] > t[2] || r[2] < t[0] || r[1] > t[3] || r[3] < t[1])
}

func (t Rectangle[N]) Cover(r Rectangle[N]) bool {
	return t[0] <= r[0] && t[1] <= r[1] && t[2] >= r[2] && t[3] >= r[3]
}

func (t Rectangle[N]) Area() (n N) {
	return (t[2] - t[0]) * (t[3] - t[1])
}

func (t Rectangle[N]) Equels(r Rectangle[N]) bool {
	return !(t[0] > r[0] || t[0] < r[0] || t[1] > r[1] || t[1] < r[1] ||
		t[2] > r[2] || t[2] < r[2] || t[3] > r[3] || t[3] < r[3])
}

// Each leaf node (unless it is the root) can host up to M entries, whereas the
// minimum allowed number of entries is m ≤ M/2. Each entry is of the form
// (mbr,oid), such that mbr is the MBR that spatially contains the object and
// oid is the object’s identifier.
type rEntry[T any, N number] struct {
	rect Rectangle[N]
	data T
	ptr  *rNode[T, N]
}

// The number of entries that each internal node can store is again between m≤M/2 and M.
// Eachentry is of the form (mbr,p), where p is a pointer to
// a child of the node and mbr is the MBR that spatially contains the MBRs contained in this child.
type rNode[T any, N number] struct {
	kind    uint8
	mbr     Rectangle[N]
	parent  *rNode[T, N]
	entries []*rEntry[T, N]
}

func (r *rNode[T, N]) isLeaf() bool {
	return (r.kind & 1) != 0
}

func (r *rNode[T, N]) isRoot() bool {
	return (r.kind & 2) != 0
}

func (r *rNode[T, N]) update(rect Rectangle[N]) {
	//r.entries[i].rect = rect
	if rect[0] < r.mbr[0] {
		r.mbr[0] = rect[0]
	}
	if rect[1] < r.mbr[1] {
		r.mbr[1] = rect[1]
	}
	if rect[2] > r.mbr[2] {
		r.mbr[2] = rect[2]
	}
	if rect[3] > r.mbr[3] {
		r.mbr[3] = rect[3]
	}
}

func (r *rNode[T, N]) reset() {
	// re calculate the mbr
	var n N
	r.mbr = Rectangle[N]{maxValue(n), maxValue(n), minValue(n), minValue(n)}
	for _, e := range r.entries {
		r.update(e.rect)
	}
}

func (r *rNode[T, N]) del(i int) {
	n := len(r.entries)
	if i < 0 || i >= n {
		return
	}
	if i < n-1 {
		copy(r.entries[i:], r.entries[i+1:])
	}
	r.entries = r.entries[:n-1]
	r.reset()
}

type rPath[T any, N number] struct {
	rn  *rNode[T, N]
	idx int
}

type DistFunc[N number] func(Rectangle[N], Rectangle[N]) N

func BoxDist[N number](r1, r2 Rectangle[N]) N {
	dx := r1[0] + r1[2] - r2[0] - r2[2]
	dy := r1[1] + r1[3] - r2[1] - r2[3]
	s := math.Sqrt(float64(dx*dx+dy*dy)) / 2.0
	return N(s)
}

/*
The original R-tree has two important disadvantages:

The execution of a point location query in an R-tree may lead to the investigation
of several paths from the root to the leaf level. This characteristic
may lead to performance deterioration, specifically when the overlap of
the MBRs is significant.

A few large rectangles may increase the degree of overlap significantly,
leading to performance degradation during range query execution, due to
empty space.
*/
type RTree[T any, N number] struct {
	dist DistFunc[N]
	M    int
	lock bool
	mu   sync.RWMutex
	cnt  int
	// The minimum allowed number of entries in the root node is 2, unless it is a leaf
	// (in this case, it may contain zero or a single entry).
	// All leaves of the R-tree are at the same level.
	root *rNode[T, N]
	pool sync.Pool
	maxN N
}

func NewRTree[T any, N number](M int, lock bool, dist DistFunc[N]) *RTree[T, N] {
	var n N
	t := &RTree[T, N]{
		dist: dist,
		M:    M,
		lock: lock,
		maxN: maxValue(n),
	}
	t.pool = sync.Pool{
		New: func() any {
			return &rEntry[T, N]{}
		},
	}
	return t
}

func (r *RTree[T, N]) Len() int {
	if r.lock {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}
	return r.cnt
}

// Finds all rectangles that are stored in an R-tree with root node node, which are intersected by a query rectangle q.
func (r *RTree[T, N]) search(node *rNode[T, N], q Rectangle[N]) (ret []Rectangle[N], dat []T) {
	if node == nil {
		return nil, nil
	}
	if node.isLeaf() {
		for _, e := range node.entries {
			if q.Intersect(e.rect) {
				ret = append(ret, e.rect)
				dat = append(dat, e.data)
			}
		}
	} else {
		for _, e := range node.entries {
			if q.Intersect(e.rect) {
				a, b := r.search(e.ptr, q)
				ret = append(ret, a...)
				dat = append(dat, b...)
			}
		}
	}
	return
}

// Query all data that overlaps with the rect parameter.
func (r *RTree[T, N]) Search(region Rectangle[N]) ([]Rectangle[N], []T) {
	if r.lock {
		r.mu.RLock()
		defer r.mu.RUnlock()
	}
	return r.search(r.root, region)
}

func (r *RTree[T, N]) searchPath(mbr Rectangle[N]) *stack[*rPath[T, N]] {
	// Traverse the tree from root to the appropriate leaf.
	// At each level,select the node, L,whose MBR will require the minimum area enlargement to cover mbr.
	stk := newStack[*rPath[T, N]]()
	for p := r.root; p != nil; {
		i, j := -1, 0
		si, sj := r.maxN, r.maxN
		for k, e := range p.entries {
			// In case of ties, select the node whose MBR has the minimum area.
			if e.rect.Cover(mbr) && e.rect.Area() < si {
				i = k
				si = e.rect.Area()
			}
			// In case of ties, select the node whose MBR has the minimum area.
			if i < 0 && e.rect.Area() < sj {
				j = k
				sj = e.rect.Area()
			}
		}
		rp := &rPath[T, N]{rn: p, idx: i}
		if i < 0 {
			rp.idx = j
		}
		stk.push(rp)
		if p.isLeaf() {
			break
		}
		p = p.entries[rp.idx].ptr
	}
	return stk
}

// Insert data, the rect parameter represents the bounding rectangle of the data.
func (r *RTree[T, N]) Insert(data T, rect Rectangle[N]) {
	p := r.searchPath(rect)
	e := r.pool.Get().(*rEntry[T, N])
	e.data, e.rect = data, rect
	if p.empty() {
		// tree is null, create a new root
		r.root = &rNode[T, N]{
			kind:    3, // root+leaf
			mbr:     rect,
			entries: []*rEntry[T, N]{e},
		}
	} else {
		L, _ := p.pop()
		L.rn.entries = append(L.rn.entries, e)
		// if the selected leaf L can accommodate obj. Insert obj into L.
		// Update all MBRs in the path from the root to L,so that all of them cover obj.mbr
		if len(L.rn.entries) <= r.M {
			L.rn.update(e.rect)
			mbr := L.rn.mbr
			for !p.empty() {
				// update mbr from leaf to root
				rp, _ := p.pop()
				rp.rn.entries[rp.idx].rect = mbr
				rp.rn.update(mbr)
				mbr = rp.rn.mbr
			}
		} else {
			// if L is already full
			// Let E be the set consisting of all L’s entries and the new entry obj.
			// Select as seeds two entries e1,e2 ∈ E, where the distance between
			// e1 and e2 is the maximum among all other pairs of entries from E
			// Form two nodes, L1 and L2, where the first contains e1 and the second e2.
			// Examine the remaining members of E one by one and assign them to L1 or L2,
			// depending on which of the MBRs of these nodes will require the minimum area
			// enlargement so as to cover this entry.
			L1, L2 := r.split(L.rn)
			for !p.empty() {
				// Update the MBRs of nodes that are in the path from root to L, so as to cover L1 and accommodate L2.
				rp, _ := p.pop()
				rp.rn.entries[rp.idx].rect = L1.mbr
				rp.rn.entries[rp.idx].ptr = L1
				rp.rn.update(L1.mbr)
				// insert L2 to parent node.
				if L2 != nil {
					e2 := r.pool.Get().(*rEntry[T, N])
					e2.rect, e2.ptr = L2.mbr, L2
					rp.rn.entries = append(rp.rn.entries, e2)
					rp.rn.update(L2.mbr)
				}
				// Perform splits at the upper levels if necessary.
				if len(rp.rn.entries) > r.M {
					L1, L2 = r.split(rp.rn)
					if p.empty() {
						// In case the root has to be split, create a new root,which increase the height of the tree by one.
						root := &rNode[T, N]{
							kind:    2, // root+nonLedf
							entries: []*rEntry[T, N]{nil, nil},
						}
						e1 := r.pool.Get().(*rEntry[T, N])
						e1.rect, e1.ptr = L1.mbr, L1
						e2 := r.pool.Get().(*rEntry[T, N])
						e2.rect, e2.ptr = L2.mbr, L2
						root.entries[0], root.entries[1] = e1, e2

						root.update(L1.mbr)
						root.update(L2.mbr)
						L1.parent, L2.parent = root, root
						r.root = root
					}
				} else {
					L1, L2 = rp.rn, nil
				}
			} //end for
		}
	}
	r.cnt++
}

func (r *RTree[T, N]) split(node *rNode[T, N]) (*rNode[T, N], *rNode[T, N]) {
	// The split is based on finding the two most distant points (seeds), and then
	// assigning the remaining points to the node whose bounding box requires the least enlargement.
	var maxDist N
	var s1, s2 int
	for i := 0; i < len(node.entries); i++ {
		for j := i + 1; j < len(node.entries); j++ {
			d := r.dist(node.entries[i].rect, node.entries[j].rect)
			if d > maxDist {
				s1, s2 = i, j
			}
		}
	}
	if s1 == s2 {
		s1, s2 = 0, 1
	}
	L1 := &rNode[T, N]{mbr: node.entries[s1].rect, entries: []*rEntry[T, N]{node.entries[s1]}}
	L2 := &rNode[T, N]{mbr: node.entries[s2].rect, entries: []*rEntry[T, N]{node.entries[s2]}}
	L1.parent, L2.parent = node.parent, node.parent
	L1.kind, L2.kind = (node.kind & 1), (node.kind & 1)
	for i, e := range node.entries {
		if i == s1 || i == s2 {
			continue
		}
		d1 := r.dist(L1.entries[0].rect, e.rect)
		d2 := r.dist(L2.entries[0].rect, e.rect)
		if d1 <= d2 {
			L1.entries = append(L1.entries, e)
		} else {
			L2.entries = append(L2.entries, e)
		}
		// TODO:
		// Assign the entry to the node whose MBR has the smaller area.
		// Assign the entry to the node that contains the smaller number of entries.
	}
	node.parent = nil
	// if during the assignment of entries, there remain λ entries to be assigned and the one node contains m−λ entries.
	// Assign all the remaining entries to this node without considering the aforementioned criteria.
	// so that the node will contain at least m entries.
	m := r.M / 2
	if len(L1.entries) < m {
		n := len(L2.entries) + len(L1.entries) - m
		L1.entries = append(L1.entries, L2.entries[n:]...)
		L2.entries = L2.entries[:n:n]
	} else if len(L2.entries) < m {
		n := len(L1.entries) + len(L2.entries) - m
		L2.entries = append(L2.entries, L1.entries[n:]...)
		L1.entries = L1.entries[:n:n]
	}
	L1.reset()
	L2.reset()
	for _, e := range L1.entries {
		if e.ptr != nil {
			e.ptr.parent = L1
		}
	}
	for _, e := range L2.entries {
		if e.ptr != nil {
			e.ptr.parent = L2
		}
	}
	return L1, L2
}

func (r *RTree[T, N]) find(node *rNode[T, N], rect Rectangle[N]) (*rNode[T, N], int) {
	if node == nil {
		return nil, 0
	}
	if node.isLeaf() {
		for i, e := range node.entries {
			if e.rect.Equels(rect) {
				return node, i
			}
		}
	} else {
		for _, e := range node.entries {
			if e.rect.Cover(rect) {
				n, i := r.find(e.ptr, rect)
				if n != nil {
					return n, i
				}
			}
		}
	}
	return nil, 0
}

// Delete the data corresponding to the region parameter.
func (r *RTree[T, N]) Delete(region Rectangle[N]) {
	L, i := r.find(r.root, region)
	if L == nil {
		return
	}
	L.del(i)
	r.condense(L)
	r.cnt--
	// if the root has only one child:Remove the root,Set as new root its only child
	if len(r.root.entries) == 1 {
		child := r.root.entries[0].ptr
		if child != nil {
			child.kind |= 2
			child.parent = nil
			// pool
			r.pool.Put(r.root.entries[0])
			r.root = child
		}
	}
}

/*
Given is the leaf L from which an entry E has been deleted. If after
the deletion of E, L has fewer than m entries, then remove entirely
leaf L and reinsert all its entries. Updates are propagated upwards and
the MBRs in the path from root to L are modified (possibly become smaller)
*/
func (r *RTree[T, N]) condense(L *rNode[T, N]) {
	// let RN be the set of nodes that are going to be removed from the tree (initially, RN is empty)
	var RN []*rNode[T, N]
	x := L
	for !x.isRoot() {
		px := x.parent
		i := -1
		for j, e := range px.entries {
			if e.ptr == x {
				i = j
				break
			}
		}
		// i be the entry of ParentX that corresponds to X
		if len(x.entries) < r.M/2 {
			//r.pool.Put(px.entries[i])
			px.del(i)
			RN = append(RN, x)
		} else {
			px.entries[i].rect = x.mbr
			px.reset()
		}
		x = px
	}
	// Reinsert all the entries of nodes that are in the set RN
	for _, tree := range RN {
		for _, node := range r.leaf(tree) {
			for _, e := range node.entries {
				r.Insert(e.data, e.rect)
				r.pool.Put(e)
			}
		}
	}
}

func (r *RTree[T, N]) leaf(node *rNode[T, N]) []*rNode[T, N] {
	if node == nil {
		return nil
	}
	if node.isLeaf() {
		return []*rNode[T, N]{node}
	} else {
		var res []*rNode[T, N]
		for _, e := range node.entries {
			res = append(res, r.leaf(e.ptr)...)
		}
		return res
	}
}

// Scan all data (note that the scanning order is random).
func (r *RTree[T, N]) Scan(fn func(Rectangle[N], T) error) error {
	stk := newStack[*rPath[T, N]]()
	p := r.root
	for !stk.empty() || p != nil {
		for p != nil {
			rp := &rPath[T, N]{rn: p, idx: 0}
			stk.push(rp)
			if p.isLeaf() {
				p = nil
				break
			}
			p = p.entries[0].ptr
		}
		if !stk.empty() {
			rp := stk.top()
			if rp.rn.isLeaf() {
				// visited rn
				L, _ := stk.pop()
				for _, e := range L.rn.entries {
					if err := fn(e.rect, e.data); err != nil {
						return err
					}
				}
			} else {
				rp.idx++
				if rp.idx >= len(rp.rn.entries) {
					_, _ = stk.pop()
				} else {
					p = rp.rn.entries[rp.idx].ptr
				}
			}
		}
	}
	return nil
}

// Query the k nearest data points to the given obj.
func (r *RTree[T, N]) NearestNeighbors(obj Rectangle[N], dist DistFunc[N], k int) ([]Rectangle[N], []T) {
	if k <= 0 {
		return nil, nil
	}
	hp := newBinaryHeap[Rectangle[N], T, N](func(a, b N) bool { return a > b })
	hp.init()

	_ = r.Scan(func(rect Rectangle[N], data T) error {
		d := dist(rect, obj)
		if hp.length() < k || hp.top().rank > d {
			if hp.length() >= k {
				_ = hp.pop()
			}
			hp.push(&element[Rectangle[N], T, N]{
				key:  rect,
				val:  data,
				rank: d,
			})
		}
		return nil
	})

	var rs []Rectangle[N]
	var ds []T
	for hp.length() > 0 {
		p := hp.pop()
		rs = append(rs, p.key)
		ds = append(ds, p.val)
	}
	return rs, ds
}

// Query all data whose distance from the given obj is less than or equal to the limit.
func (r *RTree[T, N]) NearestNeighborsByDist(obj Rectangle[N], dist DistFunc[N], limit N) ([]Rectangle[N], []T) {
	if limit < 0 {
		return nil, nil
	}
	var rs []Rectangle[N]
	var ds []T
	_ = r.Scan(func(rect Rectangle[N], data T) error {
		if dist(rect, obj) <= limit {
			rs = append(rs, rect)
			ds = append(ds, data)

		}
		return nil
	})
	return rs, ds
}

// R+-trees do not allow overlapping of MBRs at the same tree level.
// In turn, to achieve this,inserted objects have to be divided in two or more MBRs,
// which means that a specific object’s entries may be duplicated and redundantly stored in several nodes.
type RPlusTree[T any, N number] struct {
}
