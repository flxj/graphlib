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
	"fmt"
)

// Create weight matrix for graph.
func NewWeightMatrix[K comparable, W number](g Graph[K, W]) (*WeightMatrix[K, W], error) {
	if g == nil {
		return nil, errNilGraph
	}
	p, _ := g.Property(ProSimple)
	if !p.Value.(bool) {
		return nil, errNotSimple
	}

	vs := g.AllVertexes()
	es := g.AllEdges()

	var n W
	none := getMaxValue(n)
	wm := &WeightMatrix[K, W]{
		none:     none,
		vertexes: make([]K, len(vs)),
		data:     make([][]W, len(vs)),
	}
	idx := make(map[K]int)
	for i, v := range vs {
		idx[v.Key] = i
		wm.vertexes[i] = v.Key
		wm.data[i] = make([]W, len(vs))
		for j := 0; j < len(vs); j++ {
			wm.data[i][j] = none
			if i == j {
				wm.data[i][j] = any(0).(W)
			}
		}
	}
	//
	for _, e := range es {
		i := idx[e.Head]
		j := idx[e.Tail]
		wm.data[j][i] = e.Weight
		if !g.IsDigraph() {
			wm.data[i][j] = e.Weight
		}
	}

	return wm, nil
}

// Create adjacency matrix for graph.
func NewAdjacencytMatrix[K comparable, W number](g Graph[K, W]) (*AdjacencyMatrix[K], error) {
	if g == nil {
		return nil, errNilGraph
	}

	vs := g.AllVertexes()
	es := g.AllEdges()

	am := &AdjacencyMatrix[K]{
		vertexes: make([]K, len(vs)),
		data:     make([][]int, len(vs)),
	}
	idx := make(map[K]int)
	for i, v := range vs {
		idx[v.Key] = i
		am.vertexes[i] = v.Key
		am.data[i] = make([]int, len(vs))
	}
	//
	for _, e := range es {
		i := idx[e.Head]
		j := idx[e.Tail]
		am.data[j][i] = 1
		if !g.IsDigraph() {
			am.data[i][j] = 1
		}
	}

	return am, nil
}

func NewDegreeMatrix[K comparable, W number](g Graph[K, W]) *DegreeMatrix[K] {
	if g == nil {
		return nil
	}
	vs := g.AllVertexes()
	dm := &DegreeMatrix[K]{
		vertexes: make([]K, len(vs)),
		data:     make([][]int, len(vs)),
	}
	for i, v := range vs {
		dm.vertexes[i] = v.Key
		dm.data[i] = make([]int, len(vs))
		d, ok := g.Degree(v.Key)
		if !ok {
			return nil
		}
		dm.data[i][i] = d
	}
	return dm
}

type AdjacencyMatrix[K comparable] struct {
	vertexes []K
	data     [][]int
}

func (m *AdjacencyMatrix[K]) Matrix() [][]int {
	return m.data
}

func (m *AdjacencyMatrix[K]) Columns() []K {
	return m.vertexes
}

type DegreeMatrix[K comparable] struct {
	vertexes []K
	data     [][]int
}

func (m *DegreeMatrix[K]) Degree() [][]int {
	return m.data
}

func (m *DegreeMatrix[K]) Columns() []K {
	return m.vertexes
}

type WeightMatrix[K comparable, W number] struct {
	none     W
	vertexes []K
	data     [][]W
}

func (m *WeightMatrix[K, W]) Weight(none W) [][]W {
	w := make([][]W, len(m.data))
	for i, d := range m.data {
		ds := make([]W, len(d))
		for j, p := range d {
			if p == m.none {
				ds[j] = none
			} else {
				ds[j] = p
			}
		}
		w[i] = ds
	}
	return w
}

func (m *WeightMatrix[K, W]) Distance(infinite float64) [][]float64 {
	w := make([][]float64, len(m.data))
	for i, d := range m.data {
		ds := make([]float64, len(d))
		for j, p := range d {
			if p == m.none {
				ds[j] = infinite
			} else {
				ds[j] = any(p).(float64)
			}
		}
		w[i] = ds
	}
	return w
}

func (m *WeightMatrix[K, W]) Columns() []K {
	return m.vertexes
}

type endpoint[K comparable, W number] struct {
	vtx    K // vertex key
	edge   K // edge key
	weight W
	next   *endpoint[K, W]
}

type edge[K comparable, W number] struct {
	key    K
	head   K
	tail   K
	weight W
}

type adjList[K comparable, W number] struct {
	digraph bool
	outAdj  map[K]*endpoint[K, W] // adjacency list
	inAdj   map[K]*endpoint[K, W] // contrary adjacency list
}

func newAdjacencyLis[K comparable, W number](digraph bool) *adjList[K, W] {
	adj := &adjList[K, W]{
		digraph: digraph,
		outAdj:  make(map[K]*endpoint[K, W]),
	}
	if adj.digraph {
		adj.inAdj = make(map[K]*endpoint[K, W])
	}
	return adj
}

func newAdjacencyListFromGraph[K comparable, W number](g Graph[K, W]) *adjList[K, W] {
	var adj *adjList[K, W]
	vs := g.AllVertexes()
	es := g.AllEdges()
	adj = newAdjacencyLis[K, W](g.IsDigraph())

	for _, v := range vs {
		adj.addVertexes(v.Key)
	}
	for _, e := range es {
		adj.addEdge(e.Head, e.Tail, e.Key, e.Weight)
	}
	return adj
}

func (l *adjList[K, W]) reverse() {
	var out = l.outAdj
	l.outAdj = l.inAdj
	l.inAdj = out
}

func (l *adjList[K, W]) addVertexes(vs ...K) {
	for _, v := range vs {
		if _, ok := l.outAdj[v]; !ok {
			l.outAdj[v] = nil
		}
		if l.digraph {
			if _, ok := l.inAdj[v]; !ok {
				l.inAdj[v] = nil
			}
		}
	}
}

func (l *adjList[K, W]) delVertex(v K) bool {
	del := func(v K, adj map[K]*endpoint[K, W]) {
		delete(adj, v)
		for k, p := range adj {
			var head = p
			var prev = &endpoint[K, W]{next: head}

			for q := head; q != nil; {
				if q.vtx == v {
					if q == head {
						// remove head element
						prev = q
						q = q.next
						head = q
						prev.next = nil
					} else {
						//
						prev.next = q.next
						q = q.next
					}
				} else {
					prev = q
					q = q.next
				}
			}
			if head != p {
				adj[k] = head
			}
		}
	}
	del(v, l.outAdj)
	if l.digraph {
		del(v, l.inAdj)
	}
	return true
}

func (l *adjList[K, W]) delVertexes(vs ...K) bool {
	for _, v := range vs {
		if _, ok := l.outAdj[v]; !ok {
			return false
		}
	}
	for _, v := range vs {
		if ok := l.delVertex(v); !ok {
			return false
		}
	}
	return true
}

func (l *adjList[K, W]) addEdge(head, tail, key K, weight W) bool {
	insert := func(v1, v2, edge K, w W, adj map[K]*endpoint[K, W]) error {
		p, ok := adj[v1]
		if !ok {
			return fmt.Errorf("vertex %v not exists", v1)
		}
		var exists bool
		for q := p; q != nil; q = q.next {
			if q.vtx == v2 && q.edge == edge {
				q.weight = w
				exists = true
				break
			}
		}
		if !exists {
			q := &endpoint[K, W]{
				vtx:    v2,
				edge:   edge,
				weight: w,
			}
			if p != nil {
				q.next = p
			}
			adj[v1] = q
		}
		return nil
	}
	// insert to outAdj
	if err := insert(tail, head, key, weight, l.outAdj); err != nil {
		return false
	}
	if l.digraph {
		// insert to inAdj
		if err := insert(head, tail, key, weight, l.inAdj); err != nil {
			return false
		}
	} else {
		if err := insert(head, tail, key, weight, l.outAdj); err != nil {
			return false
		}
	}
	return true
}

func (l *adjList[K, W]) delEdge(head, tail, key K) bool {
	del := func(v1, v2, edge K, adj map[K]*endpoint[K, W]) error {
		p, ok := adj[v1]
		if !ok {
			return fmt.Errorf("vertex %v not exists", v1)
		}
		if p == nil {
			return fmt.Errorf("edge %v not exists", edge)
		}
		var prev = &endpoint[K, W]{next: p}
		var q *endpoint[K, W]
		for e := p; e != nil; e = e.next {
			if e.vtx == v2 && e.edge == edge {
				q = e
				break
			}
			prev = e
		}
		if q != nil {
			// remove head element of list.
			if q == p {
				adj[v1] = q.next
			} else {
				prev.next = q.next
				q.next = nil
			}
			return nil
		}
		return fmt.Errorf("edge %v not exists", edge)
	}
	//
	if err := del(tail, head, key, l.outAdj); err != nil {
		return false
	}
	if l.digraph {
		if err := del(head, tail, key, l.inAdj); err != nil {
			return false
		}
	} else {
		if err := del(head, tail, key, l.outAdj); err != nil {
			return false
		}
	}
	return true
}

func (l *adjList[K, W]) addEdges(es ...*edge[K, W]) bool {
	for _, e := range es {
		if _, ok := l.outAdj[e.tail]; !ok {
			//return fmt.Errorf("vertex %v not exists", e.tail)
			return false
		}
		if _, ok := l.outAdj[e.head]; !ok {
			//return fmt.Errorf("vertex %v not exists", e.head)
			return false
		}
	}
	for _, e := range es {
		if ok := l.addEdge(e.head, e.tail, e.key, e.weight); !ok {
			return false
		}
	}
	return true
}

func (l *adjList[K, W]) delEdges(es ...Edge[K, W]) bool {
	for _, e := range es {
		if ok := l.delEdge(e.Head, e.Tail, e.Key); !ok {
			return false
		}
	}
	return true
}

func (l *adjList[K, W]) degree(v K) (int, bool) {
	d, ok := l.outDegree(v)
	if !ok {
		return 0, false
	}
	if l.digraph {
		in, ok := l.inDegree(v)
		if !ok {
			return 0, false
		}
		d += in
	}
	return d, true
}

func (l *adjList[K, W]) outDegree(v K) (int, bool) {
	p, ok := l.outAdj[v]
	if !ok {
		//return 0, fmt.Errorf("vertex %v not exists", v)
		return 0, false
	}
	var d int
	for q := p; q != nil; q = q.next {
		d++
	}
	return d, true
}

func (l *adjList[K, W]) inDegree(v K) (int, bool) {
	var adj map[K]*endpoint[K, W]
	if l.digraph {
		adj = l.inAdj
	} else {
		adj = l.outAdj
	}
	p, ok := adj[v]
	if !ok {
		//return 0, fmt.Errorf("vertex %v not exists", v)
		return 0, false
	}
	var d int
	for q := p; q != nil; q = q.next {
		d++
	}
	return d, true
}

func (l *adjList[K, W]) neighbours(v K, multiple bool) (map[K]struct{}, bool) {
	ks := make(map[K]struct{})
	p, ok := l.outAdj[v]
	if !ok {
		//return nil, fmt.Errorf("vertex %v not exists", v)
		return nil, false
	}
	//
	for q := p; q != nil; q = q.next {
		ks[q.vtx] = struct{}{}
	}
	if l.digraph {
		p, ok = l.inAdj[v]
		if !ok {
			//return nil, fmt.Errorf("vertex %v not exists", v)
			return nil, false
		}
		for q := p; q != nil; q = q.next {
			ks[q.vtx] = struct{}{}
		}
	}
	return ks, true
}

func (l *adjList[K, W]) inNeighbours(v K, multiple bool) (map[K]int, bool) {
	var adj map[K]*endpoint[K, W]
	if l.digraph {
		adj = l.inAdj
	} else {
		adj = l.outAdj
	}
	ks := make(map[K]int)
	p, ok := adj[v]
	if !ok {
		//return nil, fmt.Errorf("vertex %v not exists", v)
		return nil, false
	}
	for q := p; q != nil; q = q.next {
		ks[q.vtx] = ks[q.vtx] + 1
	}
	return ks, true
}

func (l *adjList[K, W]) outNeighbours(v K, multiple bool) (map[K]int, bool) {
	ks := make(map[K]int)
	p, ok := l.outAdj[v]
	if !ok {
		//return nil, fmt.Errorf("vertex %v not exists", v)
		return nil, false
	}
	for q := p; q != nil; q = q.next {
		ks[q.vtx] = ks[q.vtx] + 1
	}
	return ks, true
}

func (l *adjList[K, W]) inEdges(v K) ([]K, bool) {
	var adj map[K]*endpoint[K, W]
	if l.digraph {
		adj = l.inAdj
	} else {
		adj = l.outAdj
	}
	var ks []K
	p, ok := adj[v]
	if !ok {
		//return nil, fmt.Errorf("vertex %v not exists", v)
		return nil, false
	}
	for q := p; q != nil; q = q.next {
		ks = append(ks, q.edge)
	}
	return ks, true
}

func (l *adjList[K, W]) outEdges(v K) ([]K, bool) {
	p, ok := l.outAdj[v]
	if !ok {
		//return nil, fmt.Errorf("vertex %v not exists", v)
		return nil, false
	}
	var ks []K
	for q := p; q != nil; q = q.next {
		ks = append(ks, q.edge)
	}
	return ks, true
}

func (l *adjList[K, W]) sources() ([]K, bool) {
	if !l.digraph {
		//return nil, errNotDigraph
		return nil, false
	}
	var vs []K
	for k, v := range l.inAdj {
		if v == nil {
			vs = append(vs, k)
		}
	}
	return vs, true
}

func (l *adjList[K, W]) sinks() ([]K, bool) {
	if !l.digraph {
		//return nil, errNotDigraph
		return nil, false
	}
	var vs []K
	for k, v := range l.outAdj {
		if v == nil {
			vs = append(vs, k)
		}
	}
	return vs, true
}

func (l *adjList[K, W]) minDegree() (int, bool) {
	minD := len(l.outAdj)
	for v := range l.outAdj {
		d, ok := l.degree(v)
		if !ok {
			return 0, false
		}
		if d < minD {
			minD = d
		}
	}
	return minD, true
}

func (l *adjList[K, W]) maxDegree() (int, bool) {
	maxD := len(l.outAdj)
	for v := range l.outAdj {
		d, ok := l.degree(v)
		if !ok {
			return 0, false
		}
		if d > maxD {
			maxD = d
		}
	}
	return maxD, true
}

func (l *adjList[K, W]) avgDegree() (float64, bool) {
	if len(l.outAdj) == 0 {
		return 0, false
	}
	var sum int
	for v := range l.outAdj {
		d, ok := l.degree(v)
		if !ok {
			return 0, false
		}
		sum += d
	}
	return float64(sum) / float64(len(l.outAdj)), true
}

func (l *adjList[K, W]) isDAG() bool {
	if len(l.outAdj) == 0 {
		return true
	}
	inDegrees := make(map[K]int)
	for k := range l.outAdj {
		dk, ok := l.inDegree(k)
		if !ok {
			return false
		}
		inDegrees[k] = dk
	}
	//
	for len(inDegrees) != 0 {
		var ks []K
		for k, d := range inDegrees {
			if d == 0 {
				ks = append(ks, k)
			}
		}
		if len(ks) == 0 {
			return false
		}
		for _, k := range ks {
			vs, ok := l.outNeighbours(k, true)
			if !ok {
				return false
			}
			for v := range vs {
				// loop
				if v == k {
					return false
				}
				inDegrees[v] = inDegrees[v] - 1
			}

			delete(inDegrees, k)
		}
	}
	return true
}

func (l *adjList[K, W]) isAcyclic() bool {
	if l.digraph {
		return l.isDAG()
	}

	if len(l.outAdj) == 0 {
		return true
	}

	var start K
	for k := range l.outAdj {
		start = k
		break
	}
	//
	visited := make(map[K]bool)
	prev := make(map[K]K)

	stack := newStack[K]()
	stack.push(start)
	for !stack.empty() {
		v, _ := stack.pop()
		if _, ok := visited[v]; !ok {
			visited[v] = true
		}
		vs, ok := l.neighbours(v, false)
		if !ok {
			return false
		}
		for k := range vs {
			// loop
			if k == v {
				return false
			}
			// exclude the parent vertex that visited just now.
			// (undigraph need this)
			if v != prev[k] && prev[v] != k {
				// if v has a prev,and k already visited,which means find a back edge.
				_, pv := prev[v]
				_, vk := visited[k]
				if vk && pv {
					return false
				} else {
					stack.push(k)
					prev[k] = v
				}
			}
		}
		// to dfs another components.
		if stack.empty() && len(visited) < len(l.outAdj) {
			for k := range l.outAdj {
				if _, ok := visited[k]; !ok {
					stack.push(k)
					break
				}
			}
		}
	}
	return true
}

func (l *adjList[K, W]) isUC() bool {
	if len(l.outAdj) == 0 {
		return false
	}
	var (
		ok     bool
		source []K
		sink   []K
	)
	if source, ok = l.sources(); !ok {
		return false
	}
	if sink, ok = l.sinks(); !ok {
		return false
	}
	return len(source) <= 1 && len(sink) <= 1
}

func (l *adjList[K, W]) isConnected(unidirectional bool) bool {
	if unidirectional && l.digraph {
		return l.isUC()
	}
	if len(l.outAdj) == 0 {
		return false
	}
	// bfs
	var start K
	for k := range l.outAdj {
		start = k
		break
	}
	visited := make(map[K]bool)
	que := newFIFO[K]()
	que.push(start)

	for !que.empty() {
		v, _ := que.pop()
		if _, ok := visited[v]; !ok {
			visited[v] = true
		}
		vs, ok := l.neighbours(v, false)
		if !ok {
			return false
		}
		for v := range vs {
			if _, ok := visited[v]; !ok {
				que.push(v)
			}
		}
	}
	return len(visited) == len(l.outAdj)
}

func (l *adjList[K, W]) isSimple() bool {
	if l.digraph {
		for k, v := range l.outAdj {
			heads := make(map[K]int)
			for p := v; p != nil; p = p.next {
				// loop
				if p.vtx == k {
					return false
				}
				//
				t := heads[p.vtx]
				if t >= 1 {
					return false
				} else {
					heads[p.vtx] = t + 1
					in := l.inAdj[k]
					for q := in; q != nil; q = q.next {
						if q.vtx == p.vtx {
							return false
						}
					}
				}
			}
		}
		return true
	}
	//
	for k, v := range l.outAdj {
		vs := make(map[K]struct{})
		for p := v; p != nil; p = p.next {
			if p.vtx == k {
				return false
			}
			if _, ok := vs[p.vtx]; ok {
				return false
			}
			vs[p.vtx] = struct{}{}
		}
	}
	return true
}

func (l *adjList[K, W]) isRegular() bool {
	d := -1
	for k := range l.outAdj {
		n, ok := l.degree(k)
		if !ok {
			return false
		}
		if d >= 0 {
			if n != d {
				return false
			}
		} else {
			d = n
		}
	}
	return true
}

func (l *adjList[K, W]) isForest() bool {
	return l.isAcyclic()
}

func (l *adjList[K, W]) hasLoop() bool {
	for k, v := range l.outAdj {
		for p := v; p != nil; p = p.next {
			if p.vtx == k {
				return true
			}
		}
	}
	return false
}

func (l *adjList[K, W]) hasNegativeWeight() bool {
	for _, v := range l.outAdj {
		for p := v; p != nil; p = p.next {
			if p.weight < 0 {
				return true
			}
		}
	}
	if l.digraph {
		for _, v := range l.inAdj {
			for p := v; p != nil; p = p.next {
				if p.weight < 0 {
					return true
				}
			}
		}
	}
	return false
}

func (l *adjList[K, W]) property(p int) (property[bool], bool) {
	var r bool
	switch p {
	case acyclic:
		r = l.isAcyclic()
	case connected:
		r = l.isConnected(false)
	case unilateralConnected:
		r = l.isConnected(true)
	case simple:
		r = l.isSimple()
	case regular:
		r = l.isRegular()
	case forest:
		r = l.isForest()
	case negativeWeight:
		r = l.hasNegativeWeight()
	case loop:
		r = l.hasLoop()
	default:
	}
	return property[bool]{
		name:  p,
		value: r,
	}, r
}

func (l *adjList[K, W]) incidentEdges(v K) ([]K, bool) {
	var ks []K
	p, ok := l.outAdj[v]
	if !ok {
		//return nil, fmt.Errorf("vertex %v not exists", v)
		return nil, false
	}
	for q := p; q != nil; q = q.next {
		ks = append(ks, q.edge)
	}
	if l.digraph {
		p, ok = l.inAdj[v]
		if !ok {
			//return nil, fmt.Errorf("vertex %v not exists", v)
			return nil, false
		}
		for q := p; q != nil; q = q.next {
			ks = append(ks, q.edge)
		}
	}
	return ks, true
}

func (l *adjList[K, W]) delAllEdge() {
	for k := range l.outAdj {
		l.outAdj[k] = nil
	}
	if l.digraph {
		for k := range l.inAdj {
			l.inAdj[k] = nil
		}
	}
}

func (l *adjList[K, W]) multiplicity() int {
	var m int
	for k, v := range l.outAdj {
		cnt := make(map[K]int)
		for p := v; p != nil; p = p.next {
			cnt[p.vtx] += 1
			if cnt[p.vtx] > m {
				m = cnt[p.vtx]
			}
		}
		if l.digraph {
			for p := l.inAdj[k]; p != nil; p = p.next {
				if _, ok := cnt[p.vtx]; ok {
					cnt[p.vtx] += 1
					if cnt[p.vtx] > m {
						m = cnt[p.vtx]
					}
				}
			}
		}
	}
	return m
}
