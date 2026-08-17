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
	"fmt"
	"math/rand"
	"strconv"
)

type ThinGraph[K comparable] interface {
	Graph[K, any, int]
}

type ThinDigraph[K comparable] interface {
	Digraph[K, any, int]
}

type thinGraph[K comparable] struct {
	Graph[K, any, int]
}

type thinDigraph[K comparable] struct {
	Digraph[K, any, int]
}

func NewThinGraph[K comparable](digraph bool) ThinGraph[K] {
	g := newGraph[K, any, int](digraph, "")
	return &thinGraph[K]{g}
}

func NewThinDigraph[K comparable]() ThinDigraph[K] {
	g := NewDigraph[K, any, int]("")
	return &thinDigraph[K]{g}
}

type ThinTree[K comparable] struct {
	Graph[K, any, int]
	root K
	vtx  []K
	idx  map[K]int
	duf  *dynamicUnionFind
}

func NewThinTree[K comparable]() *ThinTree[K] {
	g := newGraph[K, any, int](false, "")
	t := &ThinTree[K]{Graph: g}
	t.idx = make(map[K]int)
	t.duf = newDynamicUnionFind(0)
	return t
}

func (t *ThinTree[K]) AddVertex(v Vertex[K, any]) error {
	err := t.Graph.AddVertex(v)
	if err != nil {
		if err == errVertexExists {
			return nil
		}
		return err
	}
	t.vtx = append(t.vtx, v.Key)
	t.idx[v.Key] = len(t.vtx) - 1
	t.duf.Add(1)
	return nil
}

func (t *ThinTree[K]) AddEdge(e Edge[K, int]) error {
	var ok bool
	var u, v int
	if u, ok = t.idx[e.Head]; !ok {
		return errVertexNotExists
	}
	if v, ok = t.idx[e.Tail]; !ok {
		return errVertexNotExists
	}
	if t.duf.Find(u) == t.duf.Find(v) {
		return errExistsCycle
	}
	if err := t.Graph.AddEdge(e); err != nil {
		return err
	}
	t.duf.Union(u, v)
	return nil
}

func (t *ThinTree[K]) RemoveVertex(k K) error {
	if err := t.Graph.RemoveVertex(k); err != nil {
		return err
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveVertexs(keys ...K) error {
	for _, k := range keys {
		if err := t.Graph.RemoveVertex(k); err != nil {
			return err
		}
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveEdge(endpoint1, endpoint2 K) error {
	if err := t.Graph.RemoveEdge(endpoint1, endpoint2); err != nil {
		return err
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveEdges(endpoint1, endpoint2 []K) error {
	if len(endpoint1) != len(endpoint2) {
		return errors.New("")
	}
	for i := 0; i < len(endpoint1); i++ {
		if err := t.Graph.RemoveEdge(endpoint1[i], endpoint2[i]); err != nil {
			return err
		}
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveEdgeByKey(k K) error {
	if err := t.Graph.RemoveEdgeByKey(k); err != nil {
		return err
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveEdgeByKeys(keys ...K) error {
	for _, k := range keys {
		if err := t.Graph.RemoveEdgeByKey(k); err != nil {
			return err
		}
	}
	return t.rebuild()
}

func (t *ThinTree[K]) rebuild() error {
	vs := t.AllVertexes()
	es := t.AllEdges()
	t.vtx = make([]K, len(vs))
	for i, v := range vs {
		t.vtx[i] = v.Key
		t.idx[v.Key] = i
	}
	t.duf = newDynamicUnionFind(len(t.vtx))
	for _, e := range es {
		t.duf.Union(t.idx[e.Head], t.idx[e.Tail])
	}
	return nil
}

func (t *ThinTree[K]) SetRoot(k K) error {
	if _, ok := t.idx[k]; !ok {
		return errVertexNotExists
	}
	t.root = k
	return nil
}

func (t *ThinTree[K]) GetRoot() (K, bool) {
	_, ok := t.idx[t.root]
	return t.root, ok
}

// Tarjan
func (t *ThinTree[K]) LeastCommonAncestor(k1, k2 K) (k K, b bool) {
	if k1 == k2 || t.root == k1 {
		return k1, true
	} else if t.root == k2 {
		return k2, true
	}

	if len(t.vtx) < 2 {
		return
	}
	var ok bool
	var v1, v2 int
	if v1, ok = t.idx[k1]; !ok {
		return
	}
	if v2, ok = t.idx[k2]; !ok {
		return
	}

	vst := make([]bool, len(t.vtx))
	duf := newDynamicUnionFind(len(t.vtx))

	var dfs func(u int)
	dfs = func(u int) {
		if vst[u] {
			return
		}
		vst[u] = true // visited u
		ns, err := t.Neighbours(t.vtx[u])
		if err != nil {
			return
		}
		for _, n := range ns {
			v := t.idx[n.Key]
			if !vst[v] {
				dfs(v)
				duf.SetParent(v, u)
			}
		}
		if u == v1 && vst[v2] {
			k, b = t.vtx[duf.Find(v2)], true
			return
		}
		if u == v2 && vst[v1] {
			k, b = t.vtx[duf.Find(v1)], true
			return
		}
	}
	dfs(t.idx[t.root])
	return
}

func PetersenGraph() Graph[int, any, int] {
	g := NewGraph[int, any, int](false, "petersen_graph")
	for i := 0; i < 10; i++ {
		_ = g.AddVertex(Vertex[int, any]{Key: i})
	}
	edges := []Edge[int, int]{
		{Key: 0, Head: 0, Tail: 1},
		{Key: 1, Head: 1, Tail: 2},
		{Key: 2, Head: 2, Tail: 3},
		{Key: 3, Head: 3, Tail: 4},
		{Key: 4, Head: 4, Tail: 0},
		{Key: 5, Head: 0, Tail: 5},
		{Key: 6, Head: 1, Tail: 6},
		{Key: 7, Head: 2, Tail: 7},
		{Key: 8, Head: 3, Tail: 8},
		{Key: 9, Head: 4, Tail: 9},
		{Key: 10, Head: 5, Tail: 7},
		{Key: 11, Head: 7, Tail: 9},
		{Key: 12, Head: 9, Tail: 6},
		{Key: 13, Head: 6, Tail: 8},
		{Key: 14, Head: 8, Tail: 5},
	}
	for _, e := range edges {
		_ = g.AddEdge(e)
	}
	return g
}

func CompleteGraph(n int) Graph[int, any, int] {
	if n < 0 {
		return nil
	}
	g := NewGraph[int, any, int](false, "k"+strconv.Itoa(n))
	for v := 0; v < n; v++ {
		_ = g.AddVertex(Vertex[int, any]{Key: v})
	}
	var k int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			_ = g.AddEdge(Edge[int, int]{Key: k, Head: i, Tail: j})
			k++
		}
	}
	return g
}

func CompleteBipartite(a, b int) Bipartite[int, any, int] {
	if a < 1 || b < 1 {
		return nil
	}
	g := NewBipartite[int, any, int](false, fmt.Sprintf("K%d_%d", a, b))
	for i := 0; i < a; i++ {
		_ = g.AddVertexTo(Vertex[int, any]{Key: i}, true)
	}
	for i := a; i < a+b; i++ {
		_ = g.AddVertexTo(Vertex[int, any]{Key: i}, false)
	}
	var k int
	for i := 0; i < a; i++ {
		for j := a; j < a+b; j++ {
			_ = g.AddEdge(Edge[int, int]{Key: k, Head: i, Tail: j})
			k++
		}
	}
	return g
}

func HajosGraph() Graph[int, any, int] {
	g := NewGraph[int, any, int](false, "hajos_graph")
	for i := 0; i < 7; i++ {
		_ = g.AddVertex(Vertex[int, any]{Key: i})
	}
	edges := []Edge[int, int]{
		{Head: 0, Tail: 1},
		{Head: 0, Tail: 2},
		{Head: 0, Tail: 4},
		{Head: 1, Tail: 3},
		{Head: 1, Tail: 6},
		{Head: 2, Tail: 4},
		{Head: 2, Tail: 5},
		{Head: 3, Tail: 5},
		{Head: 3, Tail: 6},
		{Head: 4, Tail: 5},
		{Head: 5, Tail: 6},
	}
	for i, e := range edges {
		e.Key = i
		_ = g.AddEdge(e)
	}
	return g
}

func Hypercube(n int) Graph[int, any, int] {
	if n < 0 || n > 31 {
		return nil
	}
	g := NewGraph[int, any, int](false, strconv.Itoa(n)+"-cube")
	for i := 0; i < (1 << n); i++ {
		_ = g.AddVertex(Vertex[int, any]{Key: i})
	}
	var k int
	for u := 0; u < (1 << n); u++ {
		// find u's neighbour v such that larger than u
		for b := 0; b < n; b++ {
			if (u>>b)&1 == 0 {
				v := u | (1 << b)
				_ = g.AddEdge(Edge[int, int]{Key: k, Head: u, Tail: v})
				k++
			}
		}
	}
	return g
}

func RandomGraph(n int, p float64) Graph[int, any, int] {
	if p < 0.0 || p > 1.0 {
		return nil
	}
	g := NewGraph[int, any, int](false, "k"+strconv.Itoa(n))
	for v := 0; v < n; v++ {
		_ = g.AddVertex(Vertex[int, any]{Key: v})
	}
	M := int(p * 10000.0)
	var k int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// probability
			if rand.Intn(10000) < M {
				_ = g.AddEdge(Edge[int, int]{Key: k, Head: i, Tail: j})
				k++
			}
		}
	}
	return nil
}

func FanoPlane() Graph[int, any, int] {
	g := NewGraph[int, any, int](false, "fano_plane")
	for v := 0; v < 7; v++ {
		_ = g.AddVertex(Vertex[int, any]{Key: v})
	}
	E := []Edge[int, int]{
		{Head: 0, Tail: 3},
		{Head: 3, Tail: 1},
		{Head: 1, Tail: 4},
		{Head: 4, Tail: 2},
		{Head: 2, Tail: 5},
		{Head: 5, Tail: 0},
		{Head: 0, Tail: 6},
		{Head: 1, Tail: 6},
		{Head: 2, Tail: 6},
		{Head: 3, Tail: 6},
		{Head: 4, Tail: 6},
		{Head: 5, Tail: 6},
		{Head: 3, Tail: 5},
		{Head: 3, Tail: 4},
		{Head: 4, Tail: 5},
	}
	for i, e := range E {
		e.Key = i
		_ = g.AddEdge(e)
	}
	return g
}

func HeawoodGraph() Graph[int, any, int] {
	g := NewGraph[int, any, int](false, "fano_plane")
	for v := 0; v < 14; v++ {
		_ = g.AddVertex(Vertex[int, any]{Key: v})
	}
	ek := 0
	for v := 0; v < 14; v++ {
		_ = g.AddEdge(Edge[int, int]{Key: ek, Head: v, Tail: (v + 1) % 14})
		ek++
	}
	E := []Edge[int, int]{
		{Head: 0, Tail: 5},
		{Head: 1, Tail: 10},
		{Head: 2, Tail: 7},
		{Head: 3, Tail: 12},
		{Head: 4, Tail: 9},
		{Head: 6, Tail: 11},
		{Head: 8, Tail: 13},
	}
	for _, e := range E {
		e.Key = ek
		_ = g.AddEdge(e)
		ek++
	}
	return g
}

func RandomTournament(n int) Digraph[int, any, int] {
	if n < 0 {
		return nil
	}
	g := NewDigraph[int, any, int]("tournament" + strconv.Itoa(n))
	for v := 0; v < n; v++ {
		_ = g.AddVertex(Vertex[int, any]{Key: v})
	}
	var ek int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rand.Intn(2) == 1 {
				_ = g.AddEdge(Edge[int, int]{Key: ek, Head: i, Tail: j})
			} else {
				_ = g.AddEdge(Edge[int, int]{Key: ek, Head: j, Tail: i})
			}
			ek++
		}
	}
	return g
}

func CompleteKpartite(parts []int) Graph[int, any, int] {
	g := NewGraph[int, any, int](false, fmt.Sprintf("kpartite_%d", len(parts)))
	var vtx [][]int
	var idx int
	for _, p := range parts {
		var v []int
		if p <= 0 {
			return nil
		}
		for i := 0; i < p; i++ {
			_ = g.AddVertex(Vertex[int, any]{Key: idx})
			v = append(v, idx)
			idx++
		}
		vtx = append(vtx, v)
	}
	idx = 0
	for i := 0; i < len(vtx); i++ {
		for _, u := range vtx[i] {
			for j := i + 1; j < len(vtx); j++ {
				for _, v := range vtx[j] {
					_ = g.AddEdge(Edge[int, int]{Key: idx, Head: u, Tail: v})
					idx++
				}
			}
		}
	}
	return g
}

func Cycle(n int) Graph[int, any, int] {
	if n < 0 {
		return nil
	}
	g := NewGraph[int, any, int](false, fmt.Sprintf("c_%d", n))
	for i := 0; i < n; i++ {
		_ = g.AddVertex(Vertex[int, any]{Key: i})
	}
	for i := 0; i < n; i++ {
		_ = g.AddEdge(Edge[int, int]{Key: i, Head: i, Tail: (i + 1) % n})
	}
	return g
}

func QueenGraph(m, n int) Graph[int, any, int] {
	if m <= 0 || n <= 0 {
		return nil
	}
	g := NewGraph[int, any, int](false, fmt.Sprintf("queen_%d_%d", m, n))
	vtx := make([][]int, m)
	for i := range vtx {
		vtx[i] = make([]int, n)
		for j := 0; j < n; j++ {
			vtx[i][j] = n*i + j
			_ = g.AddVertex(Vertex[int, any]{
				Key: vtx[i][j],
				Labels: map[string]string{
					"row":    strconv.Itoa(i),
					"column": strconv.Itoa(j),
				}})
		}
	}
	// add edge
	var ek int
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			// (i,j) - (i,k)
			for k := j - 1; k >= 0; k-- {
				_ = g.AddEdge(Edge[int, int]{Key: ek, Head: vtx[i][j], Tail: vtx[i][k]})
				ek++
			}
			// (i,j) - (k,j)
			for k := i - 1; k >= 0; k-- {
				_ = g.AddEdge(Edge[int, int]{Key: ek, Head: vtx[i][j], Tail: vtx[k][j]})
				ek++
			}
			// (i,j) - (i-k,j-k)
			for k := 1; k <= i && k <= j; k++ {
				_ = g.AddEdge(Edge[int, int]{Key: ek, Head: vtx[i][j], Tail: vtx[i-k][j-k]})
				ek++
			}
			// (i,j) - (i-k,j+k)
			for k := 1; k <= i && j+k < n; k++ {
				_ = g.AddEdge(Edge[int, int]{Key: ek, Head: vtx[i][j], Tail: vtx[i-k][j+k]})
				ek++
			}
		}
	}
	return g
}

/*
An n-Hadamard graph is a graph on 4n vertices defined in terms of a Hadamard matrix H=h[i][j] as follows.
Define 4n symbols (r_i+, r_i-, c_i+, c_i-) for i=1...n, where r stands for "row" and c stands for "columns,"
and take these as the vertices of the graph.
Then construct two edges (r_i+,c_j+),(r_i-,c_j-), for each matrix element such that h[i][j]=1
and (r_i+,c_j-),(r_i-,c_j+) for each matrix element such that h[i][j]=-1.
*/
func HadamardGraph(h [][]int) Graph[int, any, int] {
	n := len(h)
	if n == 0 {
		return nil
	}
	g := NewGraph[int, any, int](false, fmt.Sprintf("hadamard_%d", n))
	for i := 0; i < 4*n; i++ {
		_ = g.AddVertex(Vertex[int, any]{Key: i})
	}
	var ek int
	for i := 0; i < n; i++ {
		if len(h[i]) < n {
			return nil
		}
		for j := 0; j < n; j++ {
			switch h[i][j] {
			case 1:
				_ = g.AddEdge(Edge[int, int]{Key: ek, Head: i * 4, Tail: j*4 + 2})
				_ = g.AddEdge(Edge[int, int]{Key: ek + 1, Head: i*4 + 1, Tail: j*4 + 3})
			case -1:
				_ = g.AddEdge(Edge[int, int]{Key: ek, Head: i * 4, Tail: j*4 + 3})
				_ = g.AddEdge(Edge[int, int]{Key: ek + 1, Head: i*4 + 1, Tail: j*4 + 2})
			default:
				return nil
			}
			ek += 2
		}
	}
	return g
}

func WheelGraph(n int) Graph[int, any, int] {
	if n <= 0 {
		return nil
	}
	g := NewGraph[int, any, int](false, fmt.Sprintf("wheel_%d", n))
	for i := 0; i < n; i++ {
		_ = g.AddVertex(Vertex[int, any]{Key: i})
	}
	var ek int
	for i := 0; i < n-1; i++ {
		_ = g.AddEdge(Edge[int, int]{Key: ek, Head: i, Tail: (i + 1) % (n - 1)})
		_ = g.AddEdge(Edge[int, int]{Key: ek + 1, Head: n, Tail: i})
		ek += 2
	}
	return g
}
