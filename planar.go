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
	"sort"
)

// Check if the given graph g is a planar graph.
func CheckPlanarity[K comparable, W number](g Graph[K, W]) bool {
	if g == nil {
		return false
	}
	return CheckPlanarityLR(g)
}

// Boyer-Myrvold algorithm check if the given graph g is a planar graph.(This algorithm has not yet been implemented)
func CheckPlanarityBM[K comparable, W number](g Graph[K, W]) bool { // TODO
	panic("not implement now")
}

// Hopcroft-Tarjan algorithm check if the given graph g is a planar graph.(This algorithm has not yet been implemented)
func CheckPlanarityHT[K comparable, W number](g Graph[K, W]) bool { // TODO
	panic("not implement now")
}

// Left-Right algorithm check if the given graph g is a planar graph.
func CheckPlanarityLR[K comparable, W number](g Graph[K, W]) bool {
	if g == nil {
		return false
	}
	p := &planarTestLR[K, W]{g: g}
	return p.planarity(false)
}

type interval struct {
	low  int
	high int
}

func (i interval) isEmpty() bool  { return i.low == -1 && i.high == -1 }
func (i interval) copy() interval { return interval{low: i.low, high: i.high} }

type pair struct {
	L interval // [low,high]
	R interval
}

type planarTestLR[K comparable, W number] struct {
	g       Graph[K, W]
	vtx     []Vertex[K, W]
	edges   []Edge[K, W]
	vtxIdx  map[K]int
	edgeIdx map[K]int

	roots        []int // For each connected component, the root of its spanning DFS tree is appended to a slice.
	height       []int // vertex tree-path distance from its root.
	lowpt        []int
	lowpt2       []int
	nestingDepth []int
	parentEdge   []int
	lowptEdge    []int
	ref          []int
	side         []int
	leftRef      []int
	rightRef     []int
	oriented     [][2]int      // [e][0] -> [e][1]
	orderedAdj   map[int][]int // key:vertex index,val: ordered slice of edges

	top         int
	S           stack[*pair]
	stackBottom []*pair
}

func (p *planarTestLR[K, W]) orient(e int, tail, head int) {
	p.oriented[e] = [2]int{tail, head}
	p.orderedAdj[tail] = append(p.orderedAdj[tail], e)
}

func (p *planarTestLR[K, W]) getEdge(u, v int) int {
	es, _ := p.g.GetEdge(p.vtx[u].Key, p.vtx[v].Key)
	return p.edgeIdx[es[0].Key]
}

func (p *planarTestLR[K, W]) init(embedding bool) {
	p.vtx = p.g.AllVertexes()
	p.edges = p.g.AllEdges()
	p.vtxIdx = make(map[K]int)
	p.edgeIdx = make(map[K]int)
	for i, v := range p.vtx {
		p.vtxIdx[v.Key] = i
	}
	for i, e := range p.edges {
		p.edgeIdx[e.Key] = i
	}
	m, n := p.g.Size(), p.g.Order()
	p.height = make([]int, n)
	for i := 0; i < n; i++ {
		p.height[i] = math.MaxInt
	}
	// lowpt: height of lowest return point.
	// lowpt2: height of next-to-lowest return point(tree edges only).
	p.lowpt, p.lowpt2 = make([]int, m), make([]int, m)
	p.nestingDepth = make([]int, m)
	p.parentEdge = make([]int, n)
	for i := 0; i < n; i++ {
		p.parentEdge[i] = -1
	}
	p.oriented = make([][2]int, m)
	p.orderedAdj = make(map[int][]int)

	// for testing phase
	p.lowptEdge = make([]int, m)
	p.ref, p.side = make([]int, m), make([]int, m)
	for i := 0; i < m; i++ {
		p.ref[i] = -1
		p.side[i] = 1
	}
	p.stackBottom = make([]*pair, m)
	//
	if embedding {
		p.leftRef, p.rightRef = make([]int, n), make([]int, n)
	}
}

func (p *planarTestLR[K, W]) dfsOrient(v int) {
	e := p.parentEdge[v]
	nbs, _ := p.g.Neighbours(p.vtx[v].Key)
	for _, n := range nbs {
		w := p.vtxIdx[n.Key]
		vw := p.getEdge(v, w)
		if p.oriented[vw] == [2]int{0, 0} { // there exists some non-oriented {v,w} ∈ E
			// orient v -> w
			p.orient(vw, v, w)
			p.lowpt[vw], p.lowpt2[vw] = p.height[v], p.height[v]
			if p.height[w] == math.MaxInt {
				p.parentEdge[w] = vw
				p.height[w] = p.height[v] + 1
				p.dfsOrient(w)
			} else {
				p.lowpt[vw] = p.height[w]
			}
			// determine nesting depth
			p.nestingDepth[vw] = 2 * p.lowpt[vw]
			if p.lowpt2[vw] < p.height[v] { //  chordal
				p.nestingDepth[vw]++
			}
			// update lowpoints of parent edge e
			if e != -1 {
				if p.lowpt[vw] < p.lowpt[e] {
					p.lowpt2[e] = min(p.lowpt[e], p.lowpt2[vw])
					p.lowpt[e] = p.lowpt[vw]
				} else if p.lowpt[vw] > p.lowpt[e] {
					p.lowpt2[e] = min(p.lowpt2[e], p.lowpt[vw])
				} else {
					p.lowpt2[e] = min(p.lowpt2[e], p.lowpt2[vw])
				}
			}
		}
	}
}

func (p *planarTestLR[K, W]) getOutEdges(u int) []int {
	return p.orderedAdj[u]
}

func (p *planarTestLR[K, W]) target(e int) int {
	return p.oriented[e][1]
}

func (p *planarTestLR[K, W]) source(e int) int {
	return p.oriented[e][0]
}

func (p *planarTestLR[K, W]) dfsTesting(v int) bool {
	e := p.parentEdge[v]
	for i, ei := range p.getOutEdges(v) {
		p.stackBottom[ei] = p.S.top()
		if ei == p.parentEdge[p.target(ei)] {
			if !p.dfsTesting(p.target(ei)) {
				return false
			}
		} else {
			p.lowptEdge[ei] = ei
			p.S.push(&pair{L: interval{-1, -1}, R: interval{ei, ei}})
		}
		// integrate new return edges
		if p.lowpt[ei] < p.height[v] {
			if i == 0 {
				p.lowptEdge[e] = p.lowptEdge[ei]
			} else {
				// add constraints of ei (Algorithm 4)
				if !p.addConstraints(e, ei) {
					return false
				}
			}
		}
	}
	// remove back edges returning to parent
	if e != -1 { //  v is not root
		u := p.source(e)
		// trim back edges ending at parent u (Algorithm 5)
		p.trimBackEdges(u)
		// side of e is side of a highest return edge
		if p.lowpt[e] < p.height[u] {
			hL, hR := p.S.top().L.high, p.S.top().R.high
			if hL != -1 && (hR == -1 || p.lowpt[hL] > p.lowpt[hR]) {
				p.ref[e] = hL
			} else {
				p.ref[e] = hR
			}
		}
	}
	return true
}

func (p *planarTestLR[K, W]) addConstraints(e, ei int) bool {
	P := &pair{L: interval{-1, -1}, R: interval{-1, -1}}
	// merge return edges of ei into P.R
	for {
		Q, _ := p.S.pop()
		if !Q.L.isEmpty() {
			Q.L, Q.R = Q.R, Q.L
		}
		if !Q.L.isEmpty() {
			return false
		} else {
			if p.lowpt[Q.R.low] > p.lowpt[e] {
				if P.R.isEmpty() {
					P.R = Q.R.copy()
				} else {
					p.ref[P.R.low] = Q.R.high
				}
				P.R.low = Q.R.low
			} else {
				p.ref[Q.R.low] = p.lowptEdge[e]
			}
		}
		if p.S.top() == p.stackBottom[ei] {
			break
		}
	}
	// merge conflicting return edges of e1,...,ei−1 into P.L
	for p.conflicting(p.S.top().L, ei) || p.conflicting(p.S.top().R, ei) {
		Q, _ := p.S.pop()
		if p.conflicting(Q.R, ei) {
			Q.L, Q.R = Q.R, Q.L
		}
		if p.conflicting(Q.R, ei) {
			return false
		}
		p.ref[P.R.low] = Q.R.high
		if Q.R.low != -1 {
			P.R.low = Q.R.low
		}
		if P.L.isEmpty() {
			P.L = Q.L.copy()
		} else {
			p.ref[P.L.low] = Q.L.high
		}
		P.L.low = Q.L.low
	}
	if !(P.L.isEmpty() && P.R.isEmpty()) {
		p.S.push(P)
	}
	return true
}

func (p *planarTestLR[K, W]) conflicting(it interval, e int) bool {
	return !it.isEmpty() && p.lowpt[it.high] > p.lowpt[e]
}

func (p *planarTestLR[K, W]) trimBackEdges(u int) {
	// drop entire conflict pairs
	for !p.S.empty() && p.lowest(p.S.top()) == p.height[u] {
		P, _ := p.S.pop()
		if P.L.low != -1 {
			p.side[P.L.low] = -1
		}
	}
	if !p.S.empty() {
		P, _ := p.S.pop()
		// trim left interval
		for P.L.high != -1 && p.target(P.L.high) == u {
			P.L.high = p.ref[P.L.high]
		}
		if P.L.high == -1 && P.L.low != -1 {
			p.ref[P.L.low] = P.R.low
			p.side[P.L.low] = -1
			P.L.low = -1
		}
		// trim right interval
		for P.R.high != -1 && p.target(P.R.high) == u {
			P.R.high = p.ref[P.R.high]
		}
		if P.R.high == -1 && P.R.low != -1 {
			p.ref[P.R.low] = P.L.low
			p.side[P.R.low] = -1
			P.R.low = -1
		}
		p.S.push(P)
	}
}

func (p *planarTestLR[K, W]) lowest(P *pair) int {
	if P.L.isEmpty() {
		return p.lowpt[P.R.low]
	}
	if P.R.isEmpty() {
		return p.lowpt[P.L.low]
	}
	return min(p.lowpt[P.L.low], p.lowpt[P.R.low])
}

func (p *planarTestLR[K, W]) sign(e int) int {
	if p.ref[e] != -1 {
		p.side[e] = p.side[e] * p.sign(p.ref[e])
		p.ref[e] = -1
	}
	return p.side[e]
}

func (p *planarTestLR[K, W]) planarity(embedding bool) bool {
	n, m := p.g.Order(), p.g.Size()
	if n > 2 && m > 3*n-6 {
		return false
	} else if n < 5 {
		return true
	}
	p.init(embedding)

	// 1.orientation phase
	for s := 0; s < n; s++ {
		if p.height[s] == math.MaxInt {
			p.height[s] = 0
			p.roots = append(p.roots, s)
			p.dfsOrient(s)
		}
	}
	// g is not need,release it
	p.g = nil

	// 2. testing phase
	p.sortEdges()
	for _, r := range p.roots {
		if !p.dfsTesting(r) {
			return false
		}
	}
	// 3. embedding phase
	if !embedding {
		return true
	}
	for e := 0; e < m; e++ {
		p.nestingDepth[e] = p.sign(e) * p.nestingDepth[e]
	}
	// sort adjacency lists according to non-decreasing nesting_depth
	p.sortEdges()
	for _, r := range p.roots {
		p.dfsEmbedding(r)
	}
	return true
}

// sort adjacency lists according to non-decreasing nesting_depth
func (p *planarTestLR[K, W]) sortEdges() {
	for _, e := range p.orderedAdj {
		sort.Slice(e, func(i, j int) bool {
			return p.nestingDepth[e[i]] < p.nestingDepth[e[j]]
		})
	}
}

func (p *planarTestLR[K, W]) dfsEmbedding(v int) {
	for _, ei := range p.getOutEdges(v) {
		w := p.target(ei)
		if ei == p.parentEdge[w] {
			// TODO make ei first edge in adjacency list of w

			p.leftRef[v], p.rightRef[v] = ei, ei
			p.dfsEmbedding(v)
		} else {
			if p.side[ei] == 1 {
				// TODO place ei directly after rightRef[w] in adjacency list of w
			} else {
				// TODO place ei directly before leftRef[w] in adjacency list of w
				p.leftRef[w] = ei
			}
		}
	}
}

type planarTestHT[K comparable, W number] struct {
	g       Graph[K, W]
	vtx     []Vertex[K, W]
	edges   []Edge[K, W]
	vtxIdx  map[K]int // key: v.Key, value: index of vtx
	edgeIdx map[K]int // key: e.Key, value: index of edges

	num    int
	roots  []int // root vertex index of each connected component.
	number []int // vertex number (root vertex number is 1).
	lowpt1 []int
	lowpt2 []int

	oriented   [][3]int         // [e]{tail,head,frond}
	orderedAdj map[int][][2]int // key:vertex index(tail),value: {e,head}

	s     int   // a global variable,record the start vertex's number of the current path, and is initialized to 0.
	pathN int   // path number
	path  []int // path number of vertexes

	L    []int
	R    []int
	B    [][2]int
	FREE int
}

func (p *planarTestHT[K, W]) getEdge(u, v int) int {
	es, _ := p.g.GetEdge(p.vtx[u].Key, p.vtx[v].Key)
	return p.edgeIdx[es[0].Key]
}

func (p *planarTestHT[K, W]) init() {
	p.vtx = p.g.AllVertexes()
	p.edges = p.g.AllEdges()
	p.vtxIdx = make(map[K]int)
	p.edgeIdx = make(map[K]int)
	for i, v := range p.vtx {
		p.vtxIdx[v.Key] = i
	}
	for i, e := range p.edges {
		p.edgeIdx[e.Key] = i
	}
	m, n := p.g.Size(), p.g.Order()
	p.number = make([]int, n)
	p.lowpt1 = make([]int, n)
	p.lowpt2 = make([]int, n)
	p.oriented = make([][3]int, m)
}

func (p *planarTestHT[K, W]) orient(e int, tail, head int, frond int) {
	p.oriented[e] = [3]int{tail, head, frond}
}

func (p *planarTestHT[K, W]) dfs(u, v int) { // u is father node of v
	p.number[v] = p.num
	p.num++

	p.lowpt1[v], p.lowpt2[v] = p.number[v], p.number[v]

	vk := p.vtx[v].Key
	es, err := p.g.IncidentEdges(vk)
	if err != nil {
		return
	}
	for _, e := range es {
		var w int
		if e.Head == vk {
			w = p.vtxIdx[e.Tail]
		} else {
			w = p.vtxIdx[e.Head]
		}
		vw := p.edgeIdx[e.Key]
		if p.number[w] == 0 {
			p.orient(vw, v, w, 0) // orient v -> w
			p.dfs(v, w)
			if p.lowpt1[w] < p.lowpt1[v] {
				p.lowpt2[v] = min(p.lowpt1[v], p.lowpt2[w])
				p.lowpt1[v] = p.lowpt1[w]
			} else if p.lowpt1[w] == p.lowpt1[v] {
				p.lowpt2[v] = min(p.lowpt2[v], p.lowpt2[w])
			} else {
				p.lowpt2[v] = min(p.lowpt2[v], p.lowpt1[w])
			}
		} else {
			if p.number[w] < p.number[v] && w != u {
				// mark (v,w) as a frond
				p.orient(vw, v, w, 1)
				if p.number[w] < p.lowpt1[v] {
					p.lowpt2[v] = p.lowpt1[v]
					p.lowpt1[v] = p.number[w]
				} else if p.number[w] > p.lowpt1[v] {
					p.lowpt2[v] = min(p.lowpt2[v], p.number[w])
				}
			}
		}
	}
}

func (p *planarTestHT[K, W]) sortEdges() {
	n, m := len(p.vtx), len(p.edges)
	bucket := make([][]int, 2*n+2)
	p.orderedAdj = make(map[int][][2]int)
	var v, w, f int
	for e := 0; e < m; e++ {
		v, w, f = p.oriented[e][0], p.oriented[e][1], p.oriented[e][2]
		if f == 1 {
			bucket[2*p.number[w]] = append(bucket[2*p.number[w]], e)
		} else {
			if p.lowpt2[w] >= p.number[v] {
				bucket[2*p.lowpt1[w]] = append(bucket[2*p.lowpt1[w]], e)
			} else {
				bucket[2*p.lowpt1[w]+1] = append(bucket[2*p.lowpt1[w]+1], e)
			}
		}
	}
	for i := 1; i < 2*n+2; i++ {
		if bucket[i] != nil {
			for _, e := range bucket[i] {
				v, w = p.oriented[e][0], p.oriented[e][1]
				p.orderedAdj[v] = append(p.orderedAdj[v], [2]int{e, w})
			}
		}
	}
}

/*Stacks L and R are stored as linked lists using arrays STACK and NEXT.
STACK (i) gives a stack entry, and NEXT(i) points to the next entry on the same stack.
NEXT(0) points to the first entry on L. NEXT(-1) points to the first entry on R.
FREE is the first unused location in STACK.

Variable p denotes the number of the current path. If v is a vertex PATH(v) denotes the number of the first path containing
v. If i is the number of a path, f(i) denotes the last vertex on the path numbered i.

Blocks are represented as ordered pairs on stack B. If (x, y) is on B, x denotes the last entry on L in the block,
and y denotes the last entry on R in the block. If x = 0 (y = 0), the block has no entries on L(R).
SAVE is a temporary variable used for switching;

Let a block B be a maximal set of entries on L and
R which correspond to fronds such that the placement of any one of the fronds determines the placement of all the others.
The blocks change as the content of the stacks change, but the blocks always partition the stack entries.
*/

func (p *planarTestHT[K, W]) STACK(i int) int { // L
	return -1
}

func (p *planarTestHT[K, W]) NEXT(i int) int { // R
	return -1
}

func (p *planarTestHT[K, W]) SETNEXT(i int, val int) {
}

func (p *planarTestHT[K, W]) SETSTACK(i int, val int) {
}

// If i is the number of a path, F(i) denotes the last vertex on the path numbered i.
func (p *planarTestHT[K, W]) F(i int) int {
	return -1
}

func (p *planarTestHT[K, W]) pathFinder(v int) bool {
	vn := p.number[v]
	for _, out := range p.orderedAdj[v] {
		e, w := out[0], out[1]
		if p.oriented[e][2] == 0 { // tree edge
			if p.s == 0 { // start a new path
				p.s = vn
				p.pathN++ // pn is a path number,we use it to trace a path.
			}
			// add (v,w) to current path.
			p.path[w] = p.pathN
			if !p.pathFinder(w) {
				return false
			}
			// delete stack entries and blocks corresponding to vertices no smaller than v;
			for _, e1 := range p.B {
				xn, yn := e1[0], e1[1]
				if (p.STACK(xn) >= vn || xn == 0) && (p.STACK(yn) >= vn || yn == 0) {
					// TODO: delete (x, y) from B;
				}
			}
			var xn, yn int
			// TODO: if (xn, yn) on B has STACK(xn) >= vn then replace (xn, yn) on B by (0, yn);
			// TODO: if (xn, yn) on B has STACK(yn) >= vn then replace (xn, yn) on B by (xn, 0);
			for p.NEXT(-1) != 0 && p.STACK(p.NEXT(-1)) >= vn {
				p.SETNEXT(-1, p.NEXT(p.NEXT(-1)))
			}
			for p.NEXT(0) != 0 && p.NEXT(p.NEXT(0)) >= vn {
				p.SETNEXT(0, p.NEXT(p.NEXT(0)))
			}
			if p.path[v] != p.path[w] {
				// all of segment with first edge (v, w) has been embedded. New blocks must be moved from right to left;
				L := 0
				for _, e1 := range p.B {
					xn, yn := e1[0], e1[1]
					if p.STACK(xn) > p.F(p.path[w]) || p.STACK(yn) > p.F(p.path[w]) && p.STACK(p.NEXT(-1)) != 0 {
						if p.STACK(xn) > p.F(p.path[w]) {
							if p.STACK(yn) > p.F(p.path[w]) {
								return false
							}
							L = xn
						} else {
							SAVE := p.NEXT(L)
							p.SETNEXT(L, p.NEXT(-1))
							p.SETNEXT(-1, p.NEXT(yn))
							p.SETNEXT(yn, SAVE)
							L = yn
						}
						// TODO: delete (x, y) from B;
					}
				}
				// block on B must be combined with new blocks just deleted;
				// TODO: delete (xn, yn) from B;
				if xn != 0 {
					// TODO: add (xn,yn) to B;
				} else if L != 0 || yn != 0 {
					// TODO: add (L,yn) to B;
				}
				// delete end-of-stack marker on right stack;
				p.SETNEXT(-1, p.NEXT(p.NEXT(-1)))
			}
		} else {
			// v --> w. Current path is complete.
			if p.s == 0 {
				p.pathN++
				p.s = p.number[v]
			}
			// TODO: f(p) := w;
			L, R := 0, -1
			for (p.NEXT(L) != 0 && p.STACK(p.NEXT(L)) > p.number[w]) || (p.NEXT(R) != 0 && p.STACK(p.NEXT(R)) > p.number[w]) {
				var xn, yn int //TODO: (x,y) on B
				if xn != 0 && yn != 0 {
					if p.STACK(p.NEXT(L)) > p.number[w] {
						if p.STACK(p.NEXT(R)) > p.number[w] {
							return false
						}
						SAVE := p.NEXT(R)
						p.SETNEXT(R, p.NEXT(L))
						p.SETNEXT(L, SAVE)
						//
						SAVE = p.NEXT(xn)
						p.SETNEXT(xn, p.NEXT(yn))
						p.SETNEXT(yn, SAVE)
						L, R = yn, xn
					} else {
						L, R = xn, yn
					}
				} else if xn != 0 {
					save := p.NEXT(xn)
					p.SETNEXT(xn, p.NEXT(R))
					p.SETNEXT(R, p.NEXT(L))
					p.SETNEXT(L, save)
					R = xn
				} else {
					R = yn
				}
				// TODO: delete (x, y) from B;
			}
			// add P to left stack if p is normal;
			if p.F(p.s) < p.number[w] {
				if L == 0 {
					L = p.FREE
				}
				// TODO: STACK(FREE) := f;  ---> p.SETSTACK(p.FREE) = f
				p.SETNEXT(p.FREE, p.NEXT(0))
				p.SETNEXT(0, p.FREE)
				p.FREE++
			}
			if R == -1 {
				R = 0
			}
			if L != 0 || R != 0 || vn != p.s {
				// TODO: add (L,R) to B
			}
			if vn != p.s {
				p.SETSTACK(p.FREE, 0)
				p.SETNEXT(p.FREE, p.NEXT(-1))
				p.SETNEXT(-1, p.FREE)
				p.FREE++
			}
			// TODO: add (v,w) to current path
			// TODO: output current path
			p.s = 0
		}
	}
	return true
}

/*
When c is removed, G falls into several connected pieces, called segments. Each segment
S consists either of a single frond (v, w), or of a tree arc (v, w) plus a subtree with root w plus all fronds leading from the subtree.

A segment must be embedded completely on one side of c by the Jordan Curve Theorem.
*/

/*
We use Lemma 9 to test planarity, in the following way:
first we embed the cycle c in the plane.
Then we embed the segments one at a time in the order they are explored during pathfinding.
To embed a segment S, we find a path in it, say p. We choose a side, say the left, on which to embed p.
We compare p with previously embedded fronds to determine if p can be embedded.
If not, we move segments which have fronds blocking p from the left to the right.
If p can be embedded after moving segments, we embed it.
However, if we move segments from the left to the right we may have to move other segments from the right to the left.
Thus it may be impossible to embed p. If so, we declare the graph nonplanar.

If p can be embedded, we try to embed the rest of S by in essence using the algorithm recursively.
Then we try to embed the next segment.
*/

/*
EMBED()
began

	L := R := B := the empty stack;
	find first cycle c;
	while (some segment is unexplored) do
		begin
			initiate search for path in next segment S;
			when backing down tree arc v -> w delete entries on L and R and blocks on B containing vertices no smaller than v;
			let p: s-->*f be first path found in segment S;
			while (position of top block determines position of p) do
				begin
					delete top block from B;
					if (block entries on left) then
						switch block of entries from L to R and from R to L by switching list pointers;
					if (block still has an entry on left in conflict with p) then
						go to nonplanar;
				end;
			if (p is normal) then
				add last vertex of p to L;
			add new block to B corresponding to p and blocks just removed from B;
			add end-of-stack marker to R;
			call embedding algorithm recursively;
			for each (new block (x, y) on B) do
				begin
					if (x != 0) and (y != 0) then
						go to nonplanar;
					if (y != 0) then
						move entries in block to L;
					delete (x, y) from B;
				end;
			delete end-of-stack marker on R;
			add one block to B to represent S minus path p;
			combine top two blocks on B;
		end;

end;
*/
func (p *planarTestHT[K, W]) embed(r int) bool {
	// integer array STACK(0 :: E), NEXT(-1 :: E), f(1 :: E - V + 1), PATH(1 :: V); B(1 :: E);
	m := len(p.edges)

	p.L, p.R = make([]int, m), make([]int, m)
	p.B = make([][2]int, m)

	p.SETNEXT(-1, 0)
	p.SETNEXT(0, 0)
	p.FREE = 1
	p.SETSTACK(0, 0)
	p.pathN = 0
	p.s = 0
	p.path[r] = 1
	return p.pathFinder(r)
}

func (p *planarTestHT[K, W]) planarity() bool {
	n, m := p.g.Order(), p.g.Size()
	if n > 2 && m > 3*n-6 {
		return false
	} else if n < 5 {
		return true
	}
	p.init()
	// TODO divide G into biconnected components;

	// 1.number vertices and transform G into a palm trees;
	for s := 0; s < n; s++ {
		if p.number[s] == 0 {
			p.roots = append(p.roots, s)
			p.dfs(0, s)
		}
	}
	// 2.
	p.sortEdges()
	p.path = make([]int, n)
	for _, r := range p.roots {
		if !p.embed(r) {
			return false
		}
	}
	return true
}
