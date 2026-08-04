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
	"sort"
)

type lctNode[K any, V any, W number] struct {
	key K
	val V
	wt  W
	// Remember that the nodes in auxiliary trees are keyed by their depth in T.
	// Thus nodes to the left of v are higher(more closer to T's root) than v and nodes to the right are lower.
	pathp  *lctNode[K, V, W] // if current node is a root of aux tree, the pathp point to  xxxxxx on T.
	parent *lctNode[K, V, W] // parent node on current auxiliary tree.
	left   *lctNode[K, V, W]
	right  *lctNode[K, V, W]
}

func (n *lctNode[K, V, W]) isAuxRoot() bool {
	return n.parent == nil
}

func (n *lctNode[K, V, W]) isLeftChild() bool {
	return n.parent != nil && n == n.parent.left
}

/*
Using link-cut trees we want to maintain a forest of rooted trees whose each node has an arbitrary number of unordered child nodes.
Link-Cut Trees are a forest of trees, with splay trees—self-adjusting binary search trees—representing each tree.

preferred child：

	    the preferred child  of node v is equal to its i-th child if the last access within v’s subtree was
		in the i-th subtree and it is equal to null if the last access within v’s subtree was to v itself or if
		there were no accesses to v’s subtree at all.

preferred edge:

	A preferred edge is an edge between a preferred child and its parent.

preferred path:

	    A preferred path is a maximal continuous path of preferred edges in a tree, or a single node if there is no preferred edges incident on it.
		Thus preferred paths partition the nodes of the represented tree.

forest consist of represented trees,Link-Cut-Tree represent each tree T in the forest as a tree of auxiliary trees, one auxiliary tree for each preferred path in T.
Auxiliary trees are splay trees with each node keyed by its depth in its represented tree.
Thus for each node v in its auxiliary tree all the elements in its left subtree are higher(closer to the root)
than v in v’s represented tree and all the elements in its right subtree are lower.

Auxiliary trees are joined together using path-parent pointers. T
here is one path-parent pointer per auxiliary tree and it is stored in the root of the auxiliary tree.
*/

type LinkCutTree[K any, V any, W number] struct {
	comp   CompareFunc[K]
	nodes  []*lctNode[K, V, W] // auxiliary trees
	emptyK K
	emptyV V
}

// Create an LCT, in order to quickly find elements, a comparison function for element keys needs to be provided.
func NewLinkCutTree[K any, V any, W number](comp CompareFunc[K]) *LinkCutTree[K, V, W] {
	return &LinkCutTree[K, V, W]{comp: comp}
}

func (lct *LinkCutTree[K, V, W]) makeTree(k K, v V, w W, replace bool) bool {
	i, ok := lct.index(k)
	if ok {
		if replace {
			lct.nodes[i].val = v
			lct.nodes[i].wt = w
		}
		return replace
	}
	nd := &lctNode[K, V, W]{key: k, val: v, wt: w}
	if i == len(lct.nodes) {
		lct.nodes = append(lct.nodes, nd)
	} else {
		lct.nodes = append(lct.nodes, nil)
		copy(lct.nodes[i+1:], lct.nodes[i:])
		lct.nodes[i] = nd
	}
	return true
}

// Add a new node as a singleton tree.
// This operation allows us to add elements and later manipulate them.
// If the key already exists, do not take any action and return false.
func (lct *LinkCutTree[K, V, W]) MakeTree(k K, v V, w W) bool {
	return lct.makeTree(k, v, w, false)
}

// Add new elements to LCT and update in place if the key already exists.
func (lct *LinkCutTree[K, V, W]) AddOrUpdate(k K, v V, w W) {
	_ = lct.makeTree(k, v, w, true)
}

// Accessing a node.
func (lct *LinkCutTree[K, V, W]) Access(v K) bool {
	i, ok := lct.index(v)
	if !ok {
		return false
	}
	lct.access(lct.nodes[i])
	return true
}

/*
ACCESS(v):

 1. Splay v within its auxiliary tree,
    i.e. bring it to the root.
    The left subtree will contain all the elements higher than v
    and right subtree will contain all the elements lower than v.

 2. Remove v’s preferred child.
    path-parent(right(v)) ← v
    right(v) ← null ( + symmetric setting of parent pointer)

 3. loop until we reach the root of T
    w ← path-parent(v)
    splay w
    switch w’s preferred child
    path-parent(right(w)) ← w
    right(w) ← v ( + symmetric setting of parent pointer)
    path-parent(v) ← null
    splay v
*/
func (lct *LinkCutTree[K, V, W]) access(v *lctNode[K, V, W]) {
	if v == nil {
		return
	}
	// When we access a vertex v some of the preferred paths change.
	// A preferred path from the root of T down to v is formed.
	// When this preferred path is formed every edge on the path becomes preferred
	// and all the old preferred edges in T that had an endpoint on this path are destroyed,
	// and replaced by path-parent pointers.

	// Since we access v, its preferred child becomes null.
	// Thus, if before the access v was in the middle of a preferred path,
	// after the access the lower part of this path becomes a separate path.
	// This means that we have to separate all the nodes less than v in a separate auxiliary tree.

	// The easiest way to do this is to splay on v,
	// i.e. bring v to the root and then disconnect its right subtree, making it a separate auxiliary tree.
	lct.splay(v)
	if v.right != nil {
		v.right.pathp = v
		v.right.parent = nil // make right subtree be a new splay tree.
		v.right = nil        // cut the connection to right subtree.
	}
	// After dealing with v’s descendants, we have to make a preferred path from v up to the root of T.
	// This is where path-parent pointer will be useful in guiding us up from one auxiliary tree to another.
	// After splaying, v is the root and hence has a path-parent pointer to its parent in T, call it w.
	// We need to connect v’s preferred path with the w’s preferred path, making w the real parent of v,
	// not just path-parent(we need to set w’s preferred child to v.)

	// First, we have to disconnect the lower part of w’s preferred path the same way we did for v
	// (splay on w and disconnect its right subtree).
	// Second we have to connect v’s auxiliary tree to w’s. Since all nodes in v’s auxiliary tree are lower
	// than any node in w’s, all we have to do is to make v auxiliary tree the right subtree of w.
	for v.pathp != nil {
		w := v.pathp
		lct.splay(w)
		if w.right != nil {
			w.right.parent = nil
			w.right.pathp = w //disconnect w with it’s right child
		}
		w.right = v
		v.parent = w
		v.pathp = nil
		// Finally, we have to do another splay to finish one iteration: we do a second splay of v.
		// Since v is a child of the root w, splaying simply means rotating v to the root.
		lct.splay(v)
	}
	// We continue building up the preferred path in the same way, until we reach the root of T.
	// v will have no right child in the tree of auxiliary trees.
}

// splay a node within its auxiliary tree.
func (lct *LinkCutTree[K, V, W]) splay(x *lctNode[K, V, W]) {
	rightRotate := func(x, p, gp *lctNode[K, V, W]) {
		x.parent = gp
		p.parent = x
		p.left = x.right
		if p.left != nil {
			p.left.parent = p
		}
		x.right = p
		x.pathp = p.pathp
		if gp != nil {
			if gp.left == p {
				gp.left = x
			} else {
				gp.right = x
			}
		}
	}
	leftRotate := func(x, p, gp *lctNode[K, V, W]) {
		x.parent = gp
		p.parent = x
		p.right = x.left
		if p.right != nil {
			p.right.parent = p
		}
		x.left = p
		x.pathp = p.pathp
		if gp != nil {
			if gp.left == p {
				gp.left = x
			} else {
				gp.right = x
			}
		}
	}
	for !x.isAuxRoot() {
		p := x.parent
		if p.isAuxRoot() { // x the root’s child, Zig.
			if x.isLeftChild() {
				rightRotate(x, p, nil)
			} else {
				leftRotate(x, p, nil)
			}
		} else {
			gp := p.parent
			if x.isLeftChild() {
				if p.isLeftChild() { // x a left-left child,Zig-Zig
					rightRotate(p, gp, gp.parent)
					rightRotate(x, p, p.parent)
				} else { // x is a right-left child,Zig-Zag
					rightRotate(x, p, gp)
					leftRotate(x, gp, gp.parent)
				}
			} else {
				if p.isLeftChild() { // x is a left-right child,Zig-Zag
					leftRotate(x, p, gp)
					rightRotate(x, gp, gp.parent)
				} else { // x is a right-right child,Zig-Zig
					leftRotate(p, gp, gp.parent)
					leftRotate(x, p, p.parent)
				}
			}
		}
	}
}

func (lct *LinkCutTree[K, V, W]) index(k K) (int, bool) {
	j := sort.Search(len(lct.nodes), func(i int) bool {
		return lct.comp(lct.nodes[i].key, k) >= 0
	})
	if j >= len(lct.nodes) || lct.comp(lct.nodes[j].key, k) != 0 {
		return j, false
	}
	return j, true
}

// Returns the root of the tree that vertex v is a node of.
// This operation is interesting because path to root can be very long.
// The operation can be used to determine if two nodes u and v are connected.
func (lct *LinkCutTree[K, V, W]) FindRoot(v K) (K, V, bool) {
	i, ok := lct.index(v)
	if !ok {
		return lct.emptyK, lct.emptyV, false
	}
	// First, to find the root of v’s represented tree,
	// we access v thus make it on the same auxiliary tree as the root of the represented tree.
	// Since the root of the represented tree is the highest node, its depth in the auxiliary tree is the lowest.
	// Therefore, we go left from v as much as we can. When we stop, we have found the root.
	// We splay on it and return it.
	r := lct.nodes[i]
	lct.access(r)
	// Set v to the smallest element in the auxiliary tree, i.e. to the root of the represented tree
	for ; r.left != nil; r = r.left {
	}
	lct.access(r)
	return r.key, r.val, true
}

func (lct *LinkCutTree[K, V, W]) IsRoot(v K) bool {
	r, _, ok := lct.FindRoot(v)
	return ok && lct.comp(r, v) == 0
}

// Returns an aggregate, such as max/min/sum, of the weights of the vertex，
// on the path from the root of the tree to node v. It is also possible to augment the data
// structure to return many kinds of statistics about the path.
func (lct *LinkCutTree[K, V, W]) PathAggregate(v K, f func(K, V, W)) bool {
	i, ok := lct.index(v)
	if !ok {
		return false
	}
	lct.access(lct.nodes[i])
	lct.aggregate(lct.nodes[i].left, f)
	f(lct.nodes[i].key, lct.nodes[i].val, lct.nodes[i].wt)
	return true
}

func (lct *LinkCutTree[K, V, W]) aggregate(root *lctNode[K, V, W], f func(K, V, W)) {
	if root == nil {
		return
	}
	if root.left != nil {
		lct.aggregate(root.left, f)
	}
	f(root.key, root.val, root.wt)
	if root.right != nil {
		lct.aggregate(root.right, f)
	}
}

// Returns an sum of the weights of the vertex，
// on the path from the root of the tree to node v. It is also possible to augment the data
// structure to return many kinds of statistics about the path.
func (lct *LinkCutTree[K, V, W]) PathSum(v K) (W, bool) {
	var s W
	sum := func(_ K, _ V, w W) {
		s += w
	}
	ok := lct.PathAggregate(v, sum)
	return s, ok
}

// Deletes the edge between vertex v and its parent, parent(v) where v is not the root.
func (lct *LinkCutTree[K, V, W]) Cut(v K) bool {
	// To cut (v, parent(v)) edge in the represented tree means that we have to separate nodes in v’s
	// subtree (in represented tree T) from the tree of auxiliary trees into a separate tree of auxiliary trees.
	// To do this we access v first, since it gathers all the nodes higher(more closer to T's root) than v in v’s left subtree.
	// Then all we need to do is to disconnect v’s left subtree (in auxiliary tree) from v.
	// Note that v becomes in an auxiliary tree all by itself, but path-parent pointer from v’s children (in represented tree) still
	// point to v and hence we have the tree of auxiliary trees with the elements we wanted.
	// Therefore there are two trees of aux trees after the cut.
	i, ok := lct.index(v)
	if !ok {
		return false
	}
	x := lct.nodes[i]
	lct.access(x)
	if x.left != nil {
		x.left.parent = nil
		x.left = nil
	}
	return true
}

// How many nodes are there in the current forest.
func (lct *LinkCutTree[K, V, W]) Len() int {
	return len(lct.nodes)
}

// How many trees are there in the current forest.
func (lct *LinkCutTree[K, V, W]) Component() int {
	var cnt int
	for _, n := range lct.nodes {
		if n.parent == nil && n.pathp == nil {
			cnt++
		}
	}
	return cnt
}

// Query node data.
func (lct *LinkCutTree[K, V, W]) Get(v K) (V, bool) {
	i, ok := lct.index(v)
	if !ok {
		return lct.emptyV, false
	}
	lct.access(lct.nodes[i])
	return lct.nodes[i].val, true
}

// Update node value.
func (lct *LinkCutTree[K, V, W]) UpdateValue(k K, v V) bool {
	i, ok := lct.index(k)
	if !ok {
		return false
	}
	lct.access(lct.nodes[i])
	lct.nodes[i].val = v
	return true
}

// Modify the weight of vertex v.
func (lct *LinkCutTree[K, V, W]) UpdateWeight(k K, w W) bool {
	i, ok := lct.index(k)
	if !ok {
		return false
	}
	lct.access(lct.nodes[i])
	lct.nodes[i].wt = w
	return true
}

// Specify the node as the root of the tree.
func (lct *LinkCutTree[K, V, W]) MakeRoot(v K) bool {
	i, ok := lct.index(v)
	if !ok {
		return false
	}
	x := lct.nodes[i]
	lct.access(x)

	// inorder visit x's aux tree, flip the path from old-root to x.
	var inorder []*lctNode[K, V, W]
	lct.inOrderAuxTree(x, &inorder)
	// inorder[0] as the new root,then rebulid its left subtree by inorder[1:]
	for i, p := range inorder {
		p.right = nil
		p.left = nil
		p.parent = nil
		if i < len(inorder)-1 {
			p.left = inorder[i+1]
		}
		if i != 0 {
			p.parent = inorder[i-1]
		}
	}
	inorder[0].pathp = nil
	return true
}

func (lct *LinkCutTree[K, V, W]) inOrderAuxTree(root *lctNode[K, V, W], arr *[]*lctNode[K, V, W]) {
	if root.left != nil {
		lct.inOrderAuxTree(root.left, arr)
	}
	*arr = append(*arr, root)
	if root.right != nil {
		lct.inOrderAuxTree(root.right, arr)
	}
}

// Maintain the path from u to v.
func (lct *LinkCutTree[K, V, W]) Split(u, v K) (ok bool) {
	ok = lct.MakeRoot(u)
	if !ok {
		return
	}
	return lct.Access(v)
}

func (lct *LinkCutTree[K, V, W]) findNodeRoot(v *lctNode[K, V, W]) K {
	lct.access(v)
	for ; v.left != nil; v = v.left {
	}
	lct.access(v)
	return v.key
}

// Makes vertex v a new child of vertex w, i.e. adds an edge (v, w). In order for the
// representation to remain valid this operation assumes that v is the root of its tree and that
// v and w are nodes of distinct trees.
func (lct *LinkCutTree[K, V, W]) Link(parent, child K) (ok bool) {
	// Linking two represented trees is also easy. All we need to do is to access both v and w so that they
	// are at the roots of their trees of auxiliary trees, and make latter left child of the former.
	var i, j int
	if i, ok = lct.index(parent); !ok {
		return
	}
	if j, ok = lct.index(child); !ok {
		return
	}
	p, c := lct.nodes[i], lct.nodes[j]
	if lct.comp(lct.findNodeRoot(p), lct.findNodeRoot(c)) == 0 {
		return false
	}
	// link
	lct.access(c)
	lct.access(p)
	c.left = p
	p.parent = c
	return
}

// Determine whether two nodes are in the same tree.
func (lct *LinkCutTree[K, V, W]) Connected(v, w K) (ok bool) {
	var i, j int
	if i, ok = lct.index(v); !ok {
		return
	}
	if j, ok = lct.index(w); !ok {
		return
	}
	rv := lct.findNodeRoot(lct.nodes[i])
	rw := lct.findNodeRoot(lct.nodes[j])
	return lct.comp(rv, rw) == 0
}

func (lct *LinkCutTree[K, V, W]) lastPathParent(v *lctNode[K, V, W]) *lctNode[K, V, W] {
	if v == nil {
		return nil
	}
	lct.splay(v)
	if v.right != nil {
		v.right.pathp = v
		v.right.parent = nil // make right subtree be a new splay tree.
		v.right = nil        // cut the connection to right subtree.
	}
	pp := v
	for v.pathp != nil {
		pp = v.pathp
		lct.splay(pp)
		if pp.right != nil {
			pp.right.parent = nil
			pp.right.pathp = pp //disconnect w with it’s right child
		}
		pp.right = v
		v.parent = pp
		v.pathp = nil
		// Finally, we have to do another splay to finish one iteration: we do a second splay of v.
		// Since v is a child of the root w, splaying simply means rotating v to the root.
		lct.splay(v)
	}
	return pp
}

// For finding the LCA, we access v,
// then return the last path-parent of the node (before it becomes the root of the splay tree containing the represent tree's root).
// This last path-parent node is the node separating the subtree containing v from the subtree containing w.
func (lct *LinkCutTree[K, V, W]) LeastCommonAncestor(v, w K) (a K, ok bool) {
	var i, j int
	if i, ok = lct.index(v); !ok {
		return
	}
	if j, ok = lct.index(w); !ok {
		return
	}
	nv, nw := lct.nodes[i], lct.nodes[j]
	if lct.comp(lct.findNodeRoot(nw), lct.findNodeRoot(nv)) != 0 {
		return a, false
	}
	lct.access(nv)
	pp := lct.lastPathParent(nw)
	return pp.key, true
}
