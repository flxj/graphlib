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

func CheckPlanarity[K comparable, V any, W number](g Graph[K, V, W]) bool {
	if g == nil {
		return false
	}
	return CheckPlanarityLR(g)
}

// Boyer-Myrvold
func CheckPlanarityBM[K comparable, V any, W number](g Graph[K, V, W]) bool { // TODO
	panic("not implement now")
}

// Hopcroft-Tarjan
func CheckPlanarityHT[K comparable, V any, W number](g Graph[K, V, W]) bool { //TODO
	panic("not implement now")
}

// Left-Right
func CheckPlanarityLR[K comparable, V any, W number](g Graph[K, V, W]) bool {
	p := newPlanarTestLR(g)
	return p.lrPlanarity(false)
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

type planarTestLR[K comparable, V any, W number] struct {
	g       Graph[K, V, W]
	vtx     []Vertex[K, V]
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

func newPlanarTestLR[K comparable, V any, W number](g Graph[K, V, W]) *planarTestLR[K, V, W] {
	return &planarTestLR[K, V, W]{
		g: g,
	}
}

func (p *planarTestLR[K, V, W]) orient(e int, tail, head int) {
	p.oriented[e] = [2]int{tail, head}
	p.orderedAdj[tail] = append(p.orderedAdj[tail], e)
}

func (p *planarTestLR[K, V, W]) getEdge(u, v int) int {
	es, _ := p.g.GetEdge(p.vtx[u].Key, p.vtx[v].Key)
	return p.edgeIdx[es[0].Key]
}

func (p *planarTestLR[K, V, W]) init(embedding bool) {
	p.vtx, _ = p.g.AllVertexes()
	p.edges, _ = p.g.AllEdges()
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

func (p *planarTestLR[K, V, W]) dfsOrient(v int) {
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

func (p *planarTestLR[K, V, W]) getOutEdges(u int) []int {
	return p.orderedAdj[u]
}

func (p *planarTestLR[K, V, W]) target(e int) int {
	return p.oriented[e][1]
}

func (p *planarTestLR[K, V, W]) source(e int) int {
	return p.oriented[e][0]
}

func (p *planarTestLR[K, V, W]) dfsTesting(v int) bool {
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

func (p *planarTestLR[K, V, W]) addConstraints(e, ei int) bool {
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

func (p *planarTestLR[K, V, W]) conflicting(it interval, e int) bool {
	return !it.isEmpty() && p.lowpt[it.high] > p.lowpt[e]
}

func (p *planarTestLR[K, V, W]) trimBackEdges(u int) {
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

func (p *planarTestLR[K, V, W]) lowest(P *pair) int {
	if P.L.isEmpty() {
		return p.lowpt[P.R.low]
	}
	if P.R.isEmpty() {
		return p.lowpt[P.L.low]
	}
	return min(p.lowpt[P.L.low], p.lowpt[P.R.low])
}

func (p *planarTestLR[K, V, W]) sign(e int) int {
	if p.ref[e] != -1 {
		p.side[e] = p.side[e] * p.sign(p.ref[e])
		p.ref[e] = -1
	}
	return p.side[e]
}

func (p *planarTestLR[K, V, W]) lrPlanarity(embedding bool) bool {
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
func (p *planarTestLR[K, V, W]) sortEdges() {
	for _, e := range p.orderedAdj {
		sort.Slice(e, func(i, j int) bool {
			return p.nestingDepth[e[i]] < p.nestingDepth[e[j]]
		})
	}
}

func (p *planarTestLR[K, V, W]) dfsEmbedding(v int) {
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
