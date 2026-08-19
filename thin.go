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

type Forest[K comparable, W number] struct {
	Graph[K, W]
	roots map[K]struct{}
	vtx   []K
	idx   map[K]int
	duf   *dynamicUnionFind
}

func NewForest[K comparable, W number]() *Forest[K, W] {
	g := newGraph[K, W](false, "")
	t := &Forest[K, W]{Graph: g}
	t.roots = make(map[K]struct{}) // some tree is rooted,some not.
	t.idx = make(map[K]int)
	t.duf = newDynamicUnionFind(0)
	return t
}

func (f *Forest[K, W]) AllRoots() []K {
	var rs []K
	for k := range f.roots {
		rs = append(rs, k)
	}
	return rs
}

func (f *Forest[K, W]) Connected(k1, k2 K) (bool, error) {
	var ok bool
	var u, v int
	if u, ok = f.idx[k1]; !ok {
		return false, errVertexNotExists
	}
	if v, ok = f.idx[k2]; !ok {
		return false, errVertexNotExists
	}
	return f.duf.Find(u) == f.duf.Find(v), nil
}

func (f *Forest[K, W]) Root(k K) (K, error) {
	v, ok := f.idx[k]
	if !ok {
		return k, errVertexNotExists
	}
	if _, ok := f.roots[k]; ok {
		return k, nil
	}
	for r := range f.roots {
		u := f.idx[r]
		if f.duf.Find(u) == f.duf.Find(v) {
			return r, nil
		}
	}
	return k, errors.New("the vertex in a unrooted tree")
}

func (f *Forest[K, W]) SetRoot(k K) error {
	if _, ok := f.idx[k]; !ok {
		return errVertexNotExists
	}
	if _, ok := f.roots[k]; ok {
		return nil
	}
	r, err := f.Root(k)
	if err == nil {
		delete(f.roots, r)
	}
	f.roots[k] = struct{}{}
	return nil
}

func (f *Forest[K, W]) CancelRoot(k K) {
	delete(f.roots, k)
}

func (f *Forest[K, W]) IsTree() bool {
	return f.duf.Component() == 1
}

func (f *Forest[K, W]) Trees() int {
	return f.duf.Component()
}

func (f *Forest[K, W]) AddVertex(v Vertex[K, W]) error {
	err := f.Graph.AddVertex(v)
	if err != nil {
		if err == errVertexExists {
			return nil
		}
		return err
	}
	f.vtx = append(f.vtx, v.Key)
	f.idx[v.Key] = len(f.vtx) - 1
	f.duf.Add(1)
	return nil
}

func (f *Forest[K, W]) AddEdge(e Edge[K, W]) error {
	var ok bool
	var u, v int
	if u, ok = f.idx[e.Head]; !ok {
		return errVertexNotExists
	}
	if v, ok = f.idx[e.Tail]; !ok {
		return errVertexNotExists
	}
	if f.duf.Find(u) == f.duf.Find(v) {
		return errExistsCycle
	}
	if err := f.Graph.AddEdge(e); err != nil {
		return err
	}
	f.duf.Union(u, v)
	return nil
}

func (f *Forest[K, W]) RemoveVertex(k K) error {
	if err := f.Graph.RemoveVertex(k); err != nil {
		return err
	}
	delete(f.roots, k)
	return f.rebuild()
}

func (f *Forest[K, W]) RemoveVertexs(keys ...K) error {
	for _, k := range keys {
		if err := f.Graph.RemoveVertex(k); err != nil {
			return err
		}
		delete(f.roots, k)
	}
	return f.rebuild()
}

func (f *Forest[K, W]) RemoveEdge(endpoint1, endpoint2 K) error {
	if err := f.Graph.RemoveEdge(endpoint1, endpoint2); err != nil {
		return err
	}
	return f.rebuild()
}

func (f *Forest[K, W]) RemoveEdges(endpoint1, endpoint2 []K) error {
	if len(endpoint1) != len(endpoint2) {
		return errors.New("")
	}
	for i := 0; i < len(endpoint1); i++ {
		if err := f.Graph.RemoveEdge(endpoint1[i], endpoint2[i]); err != nil {
			return err
		}
	}
	return f.rebuild()
}

func (f *Forest[K, W]) RemoveEdgeByKey(k K) error {
	if err := f.Graph.RemoveEdgeByKey(k); err != nil {
		return err
	}
	return f.rebuild()
}

func (f *Forest[K, W]) RemoveEdgeByKeys(keys ...K) error {
	for _, k := range keys {
		if err := f.Graph.RemoveEdgeByKey(k); err != nil {
			return err
		}
	}
	return f.rebuild()
}

func (f *Forest[K, W]) rebuild() error {
	vs := f.AllVertexes()
	es := f.AllEdges()
	f.vtx = make([]K, len(vs))
	for i, v := range vs {
		f.vtx[i] = v.Key
		f.idx[v.Key] = i
	}
	f.duf = newDynamicUnionFind(len(f.vtx))
	for _, e := range es {
		f.duf.Union(f.idx[e.Head], f.idx[e.Tail])
	}
	return nil
}

// Tarjan
func (f *Forest[K, W]) LeastCommonAncestor(k1, k2 K) (K, bool) {
	if k1 == k2 {
		return k1, true
	}
	var ok bool
	var v1, v2 int
	if v1, ok = f.idx[k1]; !ok {
		return k1, ok
	}
	if v2, ok = f.idx[k2]; !ok {
		return k2, ok
	}
	r1, err := f.Root(k1)
	if err != nil {
		return k1, false
	}
	r2, err := f.Root(k2)
	if err != nil {
		return k2, false
	}
	if r1 != r2 {
		return k1, false
	}
	if r1 == k1 || r1 == k2 {
		return r1, true
	}
	if r2 == k1 || r2 == k2 {
		return r2, true
	}

	vst := make([]bool, len(f.vtx))
	duf := newDynamicUnionFind(len(f.vtx))

	var lca K
	var dfs func(u int)
	dfs = func(u int) {
		if vst[u] {
			return
		}
		vst[u] = true // visited u
		ns, err := f.Neighbours(f.vtx[u])
		if err != nil {
			return
		}
		for _, n := range ns {
			v := f.idx[n.Key]
			if !vst[v] {
				dfs(v)
				duf.SetParent(v, u)
			}
		}
		if u == v1 && vst[v2] {
			lca, ok = f.vtx[duf.Find(v2)], true
			return
		}
		if u == v2 && vst[v1] {
			lca, ok = f.vtx[duf.Find(v1)], true
			return
		}
	}
	dfs(f.idx[r1])
	return lca, ok
}

type ThinGraph[K comparable] interface {
	Graph[K, int]
}

type ThinDigraph[K comparable] interface {
	Digraph[K, int]
}

type thinGraph[K comparable] struct {
	Graph[K, int]
}

type thinDigraph[K comparable] struct {
	Digraph[K, int]
}

func NewThinGraph[K comparable](digraph bool) ThinGraph[K] {
	g := newGraph[K, int](digraph, "")
	return &thinGraph[K]{g}
}

func NewThinDigraph[K comparable]() ThinDigraph[K] {
	g := NewDigraph[K, int]("")
	return &thinDigraph[K]{g}
}

func CompleteTree(k int, n int) *Forest[int, int] {
	if k <= 0 || n < 0 {
		return nil
	}
	f := NewForest[int, int]()
	for v := 0; v < n; v++ {
		_ = f.AddVertex(Vertex[int, int]{Key: v})
		if v-1 >= 0 {
			_ = f.AddEdge(Edge[int, int]{
				Key:  v - 1,
				Head: v,
				Tail: (v - 1) / k,
			})
		}
	}
	_ = f.SetRoot(0)
	return f
}

func FullTree(k int, level int) *Forest[int, int] {
	if k <= 0 || level < 0 {
		return nil
	}
	if k == 1 {
		return CompleteTree(k, level)
	} else {
		return CompleteTree(k, (pow(k, level)-1)/(k-1))
	}
}

func PetersenGraph() Graph[int, int] {
	g := NewGraph[int, int](false, "petersen_graph")
	for i := 0; i < 10; i++ {
		_ = g.AddVertex(Vertex[int, int]{Key: i})
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

func CompleteGraph(n int) Graph[int, int] {
	if n < 0 {
		return nil
	}
	g := NewGraph[int, int](false, "k"+strconv.Itoa(n))
	for v := 0; v < n; v++ {
		_ = g.AddVertex(Vertex[int, int]{Key: v})
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

func CompleteBipartite(a, b int) Bipartite[int, int] {
	if a < 1 || b < 1 {
		return nil
	}
	g := NewBipartite[int, int](false, fmt.Sprintf("K%d_%d", a, b))
	for i := 0; i < a; i++ {
		_ = g.AddVertexTo(Vertex[int, int]{Key: i}, true)
	}
	for i := a; i < a+b; i++ {
		_ = g.AddVertexTo(Vertex[int, int]{Key: i}, false)
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

func HajosGraph() Graph[int, int] {
	g := NewGraph[int, int](false, "hajos_graph")
	for i := 0; i < 7; i++ {
		_ = g.AddVertex(Vertex[int, int]{Key: i})
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

func Hypercube(n int) Graph[int, int] {
	if n < 0 || n > 31 {
		return nil
	}
	g := NewGraph[int, int](false, strconv.Itoa(n)+"-cube")
	for i := 0; i < (1 << n); i++ {
		_ = g.AddVertex(Vertex[int, int]{Key: i})
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

func RandomGraph(n int, p float64) Graph[int, int] {
	if p < 0.0 || p > 1.0 {
		return nil
	}
	g := NewGraph[int, int](false, "k"+strconv.Itoa(n))
	for v := 0; v < n; v++ {
		_ = g.AddVertex(Vertex[int, int]{Key: v})
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

func FanoPlane() Graph[int, int] {
	g := NewGraph[int, int](false, "fano_plane")
	for v := 0; v < 7; v++ {
		_ = g.AddVertex(Vertex[int, int]{Key: v})
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

func HeawoodGraph() Graph[int, int] {
	g := NewGraph[int, int](false, "fano_plane")
	for v := 0; v < 14; v++ {
		_ = g.AddVertex(Vertex[int, int]{Key: v})
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

func RandomTournament(n int) Digraph[int, int] {
	if n < 0 {
		return nil
	}
	g := NewDigraph[int, int]("tournament" + strconv.Itoa(n))
	for v := 0; v < n; v++ {
		_ = g.AddVertex(Vertex[int, int]{Key: v})
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

func CompleteKpartite(parts []int) Graph[int, int] {
	g := NewGraph[int, int](false, fmt.Sprintf("kpartite_%d", len(parts)))
	var vtx [][]int
	var idx int
	for _, p := range parts {
		var v []int
		if p <= 0 {
			return nil
		}
		for i := 0; i < p; i++ {
			_ = g.AddVertex(Vertex[int, int]{Key: idx})
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

func Cycle(n int) Graph[int, int] {
	if n < 0 {
		return nil
	}
	g := NewGraph[int, int](false, fmt.Sprintf("c_%d", n))
	for i := 0; i < n; i++ {
		_ = g.AddVertex(Vertex[int, int]{Key: i})
	}
	for i := 0; i < n; i++ {
		_ = g.AddEdge(Edge[int, int]{Key: i, Head: i, Tail: (i + 1) % n})
	}
	return g
}

func QueenGraph(m, n int) Graph[int, int] {
	if m <= 0 || n <= 0 {
		return nil
	}
	g := NewGraph[int, int](false, fmt.Sprintf("queen_%d_%d", m, n))
	vtx := make([][]int, m)
	for i := range vtx {
		vtx[i] = make([]int, n)
		for j := 0; j < n; j++ {
			vtx[i][j] = n*i + j
			_ = g.AddVertex(Vertex[int, int]{
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
func HadamardGraph(h [][]int) Graph[int, int] {
	n := len(h)
	if n == 0 {
		return nil
	}
	g := NewGraph[int, int](false, fmt.Sprintf("hadamard_%d", n))
	for i := 0; i < 4*n; i++ {
		_ = g.AddVertex(Vertex[int, int]{Key: i})
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

func WheelGraph(n int) Graph[int, int] {
	if n <= 0 {
		return nil
	}
	g := NewGraph[int, int](false, fmt.Sprintf("wheel_%d", n))
	for i := 0; i < n; i++ {
		_ = g.AddVertex(Vertex[int, int]{Key: i})
	}
	var ek int
	for i := 0; i < n-1; i++ {
		_ = g.AddEdge(Edge[int, int]{Key: ek, Head: i, Tail: (i + 1) % (n - 1)})
		_ = g.AddEdge(Edge[int, int]{Key: ek + 1, Head: n, Tail: i})
		ek += 2
	}
	return g
}
