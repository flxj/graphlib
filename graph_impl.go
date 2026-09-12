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
	"math/rand"
)

const (
	digraph = iota
	acyclic
	simple
	regular
	connected
	forest
	loop
	negativeWeight
	unilateralConnected
)

type property[T any] struct {
	version int
	name    int
	value   T
}

func (p property[T]) clone() property[T] {
	return property[T]{
		version: p.version,
		name:    p.name,
		value:   p.value,
	}
}

type boolPropertySet[T bool] struct {
	digraph    bool
	acyclic    property[T] // no cycle and no loop
	simple     property[T] // no loop and no multi edge
	regular    property[T] // every vertex has same order
	connect    property[T] // for digraph, which means strong connection
	forest     property[T]
	loop       property[T]
	negWeight  property[T]
	uniConnect property[T]
	orient     property[T]
}

/*
type intPropertySet struct {
	minDe property[int]
	maxDe property[int]
	multi property[int]
}
*/

// graph default implement base on adjacency list.
type graph[K comparable, W number] struct {
	ver   int // start from 1
	name  string
	prop  boolPropertySet[bool] // version start from 0
	minDe property[int]
	maxDe property[int]
	multi property[int]
	avgDe property[float64]
	vtx   map[K]*Vertex[K, W]
	edges map[K]*Edge[K, W]
	adj   *adjList[K, W]
}

func newGraph[K comparable, W number](digraph bool, name string) *graph[K, W] {
	g := &graph[K, W]{
		ver:   1,
		name:  name,
		vtx:   make(map[K]*Vertex[K, W]),
		edges: make(map[K]*Edge[K, W]),
	}
	g.prop.digraph = digraph
	g.adj = newAdjacencyLis[K, W](digraph)
	return g
}

// Create a new graph.
func NewGraph[K comparable, W number](digraph bool, name string) Graph[K, W] {
	return newGraph[K, W](digraph, name)
}

// Create a new undirected graph
func NewUnDigraph[K comparable, W number](name string) Graph[K, W] {
	return newGraph[K, W](false, name)
}

// Load graph from json or yaml file.
func NewGraphFromFile[K comparable, W number](path string) (Graph[K, W], error) {
	s, err := readFile(path)
	if err != nil {
		return nil, err
	}
	return UnmarshalGraph[K, W](s)
}

// Create a graph using vertex and edge sets.
func ConstructGraph[K comparable, W number](digraph bool, name string, vertexes []Vertex[K, W], edges []Edge[K, W]) (Graph[K, W], bool) {
	g := newGraph[K, W](digraph, name)
	for _, v := range vertexes {
		if ok := g.AddVertex(v); !ok {
			return nil, false
		}
	}
	for _, e := range edges {
		if ok := g.AddEdge(e); !ok {
			return nil, false
		}
	}
	return g, true
}

func (g *graph[K, W]) Name() string {
	return g.name
}

func (g *graph[K, W]) SetName(name string) {
	g.name = name
}

func (g *graph[K, W]) IsDigraph() bool {
	return g.adj.digraph
}

func (g *graph[K, W]) IsSimple() bool {
	if g.prop.simple.version == g.ver {
		return g.prop.simple.value
	}
	//
	p, _ := g.adj.property(simple)
	p.version = g.ver
	g.prop.simple = p

	return p.value
}

func (g *graph[K, W]) HasNegativeWeight() bool {
	if g.prop.negWeight.version == g.ver {
		return g.prop.negWeight.value
	}
	p, _ := g.adj.property(negativeWeight)
	p.version = g.ver
	g.prop.negWeight = p

	return p.value
}

func (g *graph[K, W]) IsRegular() bool {
	if g.prop.regular.version == g.ver {
		return g.prop.regular.value
	}
	p, _ := g.adj.property(regular)
	p.version = g.ver
	g.prop.regular = p

	return p.value
}

func (g *graph[K, W]) IsAcyclic() bool {
	if g.prop.acyclic.version == g.ver {
		return g.prop.acyclic.value
	}
	p, _ := g.adj.property(acyclic)
	p.version = g.ver
	g.prop.acyclic = p

	return p.value
}
func (g *graph[K, W]) IsConnected(unidirectional bool) bool {
	if unidirectional && g.IsDigraph() {
		if g.prop.uniConnect.version == g.ver {
			return g.prop.uniConnect.value
		}
		p, _ := g.adj.property(unilateralConnected)
		p.version = g.ver
		g.prop.uniConnect = p

		return p.value
	}
	if g.prop.connect.version == g.ver {
		return g.prop.connect.value
	}
	p, _ := g.adj.property(connected)
	p.version = g.ver
	g.prop.connect = p

	return p.value
}

func (g *graph[K, W]) IsCompleted() bool {
	if g.IsSimple() {
		return g.MinDegree() == g.Order()-1 // TODO for bipartite graph
	}
	return false
}

func (g *graph[K, W]) IsTree() bool {
	return g.IsConnected(false) && g.IsForest()
}

func (g *graph[K, W]) IsForest() bool {
	if g.prop.forest.version == g.ver {
		return g.prop.forest.value
	}
	p, _ := g.adj.property(forest)
	p.version = g.ver
	g.prop.forest = p

	return p.value
}

func (g *graph[K, W]) HasLoop() bool {
	if g.prop.loop.version == g.ver {
		return g.prop.loop.value
	}
	p, _ := g.adj.property(loop)
	p.version = g.ver
	g.prop.loop = p

	return p.value
}

func (g *graph[K, W]) Order() int {
	return len(g.vtx)
}

func (g *graph[K, W]) Size() int {
	return len(g.edges)
}

func (g *graph[K, W]) MinDegree() int {
	if g.minDe.version == g.ver {
		return g.minDe.value
	}
	d, ok := g.adj.minDegree()
	if !ok {
		return -1
	}
	g.minDe.version = g.ver
	g.minDe.value = d
	return d
}

func (g *graph[K, W]) MaxDegree() int {
	if g.maxDe.version == g.ver {
		return g.maxDe.value
	}
	d, ok := g.adj.maxDegree()
	if !ok {
		return -1
	}
	g.maxDe.version = g.ver
	g.maxDe.value = d
	return d
}

func (g *graph[K, W]) AvgDegree() float64 {
	if g.avgDe.version == g.ver {
		return g.avgDe.value
	}
	var avg float64
	if g.Order() != 0 {
		avg = float64(2*g.Size()) / float64(g.Order())
	}
	g.avgDe.version = g.ver
	g.avgDe.value = avg
	return avg
}

func (g *graph[K, W]) Multiplicity() int {
	if g.multi.version == g.ver {
		return g.multi.value
	}
	var d int
	if g.IsSimple() {
		if g.Size() > 0 {
			d = 1
		}
	} else {
		d = g.adj.multiplicity()
	}
	g.multi.version = g.ver
	g.multi.value = d
	return d
}

func (g *graph[K, W]) Orientation() bool {
	if g.prop.orient.version == g.ver {
		return g.prop.orient.value
	}
	if !g.IsSimple() || !g.IsDigraph() {
		g.prop.orient.value = false
	} else {
		g.prop.orient.value = true
	}
	g.prop.orient.version = g.ver
	return g.prop.orient.value
}

func (g *graph[K, W]) Property(p PropertyName) (GraphProperty[any], bool) {
	gp := GraphProperty[any]{Name: p}
	switch p {
	case ProDigraph:
		gp.Value = g.IsDigraph()
	case ProAcyclic:
		gp.Value = g.IsAcyclic()
	case ProSimple:
		gp.Value = g.IsSimple()
	case ProRegular:
		gp.Value = g.IsRegular()
	case ProConnected:
		gp.Value = g.IsConnected(false)
	case ProUnilateralConnected:
		gp.Value = g.IsConnected(true)
	case ProForest:
		gp.Value = g.IsForest()
	case ProLoop:
		gp.Value = g.HasLoop()
	case ProCompleted:
		gp.Value = g.IsCompleted()
	case ProTree:
		gp.Value = g.IsTree()
	case ProNegativeWeight:
		gp.Value = g.HasNegativeWeight()
	case ProGraphName:
		gp.Value = g.Name()
	case ProOrder:
		gp.Value = g.Order()
	case ProSize:
		gp.Value = g.Size()
	case ProMaxDegree:
		gp.Value = g.MaxDegree()
	case ProMinDegree:
		gp.Value = g.MinDegree()
	case ProAvgDegree:
		gp.Value = g.AvgDegree()
	case ProMultiplicity:
		gp.Value = g.Multiplicity()
	case ProOrientation:
		gp.Value = g.Orientation()
	default:
		return gp, false
	}
	return gp, true
}

func (g *graph[K, W]) AllVertexes() []Vertex[K, W] {
	vs := make([]Vertex[K, W], len(g.vtx))
	var i int
	for _, v := range g.vtx {
		vs[i] = *v
		i++
	}
	return vs
}

func (g *graph[K, W]) AllEdges() []Edge[K, W] {
	es := make([]Edge[K, W], len(g.edges))
	var i int
	for _, e := range g.edges {
		es[i] = *e
		i++
	}
	return es
}

func (g *graph[K, W]) AddVertex(v Vertex[K, W]) bool {
	if _, ok := g.vtx[v.Key]; ok {
		return false
	}
	g.adj.addVertexes(v.Key)
	g.vtx[v.Key] = &v
	g.ver++
	return true
}

func (g *graph[K, W]) RemoveVertex(key K) (Vertex[K, W], bool) {
	v, ok := g.vtx[key]
	if !ok {
		return Vertex[K, W]{}, false
	}
	if ok := g.adj.delVertex(key); !ok {
		return Vertex[K, W]{}, false
	}

	var edges []K
	for _, e := range g.edges {
		if e.Head == key || e.Tail == key {
			edges = append(edges, e.Key)
		}
	}
	for _, k := range edges {
		delete(g.edges, k)
	}
	delete(g.vtx, key)
	g.ver++
	return *v, true
}

func (g *graph[K, W]) AddEdge(edge Edge[K, W]) bool {
	if any(edge.Key) != nil {
		if _, ok := g.edges[edge.Key]; ok {
			return false
		}
	} else {
		for {
			k, flag := randEdgeKey(edge.Head, edge.Tail)
			if !flag {
				return false
			}
			edge.Key = k
			if _, ok := g.edges[edge.Key]; ok {
				break
			}
		}
	}
	ok := g.adj.addEdge(edge.Head, edge.Tail, edge.Key, edge.Weight)
	if !ok {
		return false
	}
	g.edges[edge.Key] = &edge
	g.ver++
	return true
}

func (g *graph[K, W]) RemoveEdgeByKey(key K) (Edge[K, W], bool) {
	e, ok := g.edges[key]
	if !ok {
		return Edge[K, W]{}, false
	}
	if ok := g.adj.delEdge(e.Head, e.Tail, e.Key); !ok {
		return Edge[K, W]{}, false
	}
	delete(g.edges, key)
	g.ver++
	return *e, true
}

func (g *graph[K, W]) RemoveEdge(v1, v2 K) ([]Edge[K, W], bool) {
	var edges []Edge[K, W]
	for _, v := range g.edges {
		ok := (v.Head == v1 && v.Tail == v2)
		if g.adj.digraph {
			ok = ok || (v.Head == v2 && v.Tail == v1)
		}
		if ok {
			edges = append(edges, *v)
		}
	}
	if len(edges) == 0 {
		return nil, false
	}
	if ok := g.adj.delEdges(edges...); !ok {
		return nil, false
	}
	for _, e := range edges {
		delete(g.edges, e.Key)
	}
	g.ver++
	return edges, true
}

func (g *graph[K, W]) RemoveAllVertex() {
	g.ver = 1
	g.vtx = make(map[K]*Vertex[K, W])
	g.edges = make(map[K]*Edge[K, W])
	g.adj = newAdjacencyLis[K, W](g.prop.digraph)
}

func (g *graph[K, W]) RemoveAllEdge() {
	g.adj.delAllEdge()
	g.edges = make(map[K]*Edge[K, W])
	g.ver++
}

func (g *graph[K, W]) Degree(key K) (int, bool) {
	if _, ok := g.vtx[key]; !ok {
		return 0, false
	}
	return g.adj.degree(key)
}

func (g *graph[K, W]) Neighbours(v K) ([]Vertex[K, W], bool) {
	vs, ok := g.adj.neighbours(v, false)
	if !ok {
		return nil, false
	}
	var res []Vertex[K, W]
	for key := range vs {
		ver, ok := g.vtx[key]
		if !ok {
			return nil, false
		}
		res = append(res, *ver)
	}
	return res, true
}

func (g *graph[K, W]) GetVertex(key K) (Vertex[K, W], bool) {
	v, ok := g.vtx[key]
	if !ok {
		return Vertex[K, W]{}, false
	}
	return *v, true
}

func (g *graph[K, W]) GetEdge(v1, v2 K) ([]Edge[K, W], bool) {
	var edges []Edge[K, W]
	for _, e := range g.edges {
		ok := e.Head == v1 && e.Tail == v2
		if !g.adj.digraph {
			ok = ok || e.Head == v2 && e.Tail == v1
		}
		if ok {
			edges = append(edges, *e)
		}
	}
	return edges, len(edges) != 0
}

func (g *graph[K, W]) GetEdgeByKey(key K) (Edge[K, W], bool) {
	e, ok := g.edges[key]
	if !ok {
		return Edge[K, W]{}, false
	}
	return *e, true
}

func (g *graph[K, W]) GetVertexesByLabel(labels Labels) []Vertex[K, W] {
	var ves []Vertex[K, W]
	if labels != nil {
		for _, u := range g.vtx {
			if u.Labels != nil {
				match := true
				for k, v := range labels {
					l, ok := u.Labels[k]
					if !ok || l != v {
						match = false
						break
					}
				}
				if match {
					ves = append(ves, *u)
				}
			}
		}
	}
	return ves
}

func (g *graph[K, W]) GetEdgesByLabel(labels Labels) []Edge[K, W] {
	var edges []Edge[K, W]
	if labels != nil {
		for _, e := range g.edges {
			if e.Labels != nil {
				match := true
				for k, v := range labels {
					l, ok := e.Labels[k]
					if !ok || l != v {
						match = false
						break
					}
				}
				if match {
					edges = append(edges, *e)
				}
			}
		}
	}
	return edges
}

func (g *graph[K, W]) SetVertexValue(key K, value any) bool {
	v, ok := g.vtx[key]
	if !ok {
		return false
	}
	v.Value = value
	return true
}

func (g *graph[K, W]) SetVertexLabel(key K, labelKey string, labelVal any) bool {
	v, ok := g.vtx[key]
	if !ok {
		return false
	}
	if v.Labels == nil {
		v.Labels = make(map[string]any)
	}
	v.Labels[labelKey] = labelVal
	return true
}

func (g *graph[K, W]) DeleteVertexLabel(key K, labelKey string) bool {
	v, ok := g.vtx[key]
	if !ok {
		return false
	}
	if v.Labels != nil {
		delete(v.Labels, labelKey)
	}
	return true
}

func (g *graph[K, W]) SetEdgeValueByKey(key K, value any) bool {
	e, ok := g.edges[key]
	if !ok {
		return false
	}
	e.Value = value
	return true
}

func (g *graph[K, W]) SetEdgeLabelByKey(key K, labelKey string, labelVal any) bool {
	e, ok := g.edges[key]
	if !ok {
		return false
	}
	if e.Labels == nil {
		e.Labels = make(map[string]any)
	}
	e.Labels[labelKey] = labelVal
	return true
}

func (g *graph[K, W]) DeleteEdgeLabelByKey(key K, labelKey string) bool {
	e, ok := g.edges[key]
	if !ok {
		return false
	}
	if e.Labels != nil {
		delete(e.Labels, labelKey)
	}
	return true
}

func (g *graph[K, W]) SetEdgeValue(endpoint1, endpoint2 K, value any) bool {
	edges, ok := g.GetEdge(endpoint1, endpoint2)
	if !ok {
		return false
	}
	for _, ed := range edges {
		e, ok := g.edges[ed.Key]
		if !ok {
			return false
		}
		e.Value = value
	}
	return true
}

func (g *graph[K, W]) SetEdgeLabel(endpoint1, endpoint2 K, labelKey string, labelVal any) bool {
	edges, ok := g.GetEdge(endpoint1, endpoint2)
	if !ok {
		return false
	}
	for _, ed := range edges {
		e, ok := g.edges[ed.Key]
		if !ok {
			return false
		}
		if e.Labels == nil {
			e.Labels = make(map[string]any)
		}
		e.Labels[labelKey] = labelVal
	}
	return true
}

func (g *graph[K, W]) DeleteEdgeLabel(endpoint1, endpoint2 K, labelKey string) bool {
	edges, ok := g.GetEdge(endpoint1, endpoint2)
	if !ok {
		return false
	}
	for _, ed := range edges {
		e, ok := g.edges[ed.Key]
		if !ok {
			return false
		}
		if e.Labels != nil {
			delete(e.Labels, labelKey)
		}
	}
	return true
}

func (g *graph[K, W]) Clone() Graph[K, W] {
	adjList := newAdjacencyLis[K, W](g.prop.digraph)
	ng := *g
	ng.vtx = make(map[K]*Vertex[K, W])
	ng.edges = make(map[K]*Edge[K, W])
	ng.adj = adjList

	for k, v := range g.vtx {
		nv := v.Clone()
		ng.vtx[k] = &nv
		ng.adj.addVertexes(k)
	}
	for k, v := range g.edges {
		nv := v.Clone()
		ng.edges[k] = &nv
		ok := ng.adj.addEdge(v.Head, v.Tail, v.Key, v.Weight)
		if !ok {
			return nil
		}
	}
	return &ng
}

func (g *graph[K, W]) RandomVertex() (Vertex[K, W], bool) {
	n := rand.Intn(len(g.vtx))
	i := 0
	for _, v := range g.vtx {
		if n == i {
			return *v, true
		}
		i++
	}
	return Vertex[K, W]{}, false
}

func (g *graph[K, W]) RandomEdge() (Edge[K, W], bool) {
	n := rand.Intn(len(g.edges))
	i := 0
	for _, e := range g.edges {
		if n == i {
			return *e, true
		}
		i++
	}
	return Edge[K, W]{}, false
}

func (g *graph[K, W]) NeighbourEdgesByKey(edge K) ([]Edge[K, W], bool) {
	e, ok := g.edges[edge]
	if !ok {
		return nil, false
	}
	var res []Edge[K, W]
	for _, ee := range g.edges {
		if ee.Key != e.Key {
			if ee.Tail == e.Head || ee.Tail == e.Tail || ee.Head == e.Tail || ee.Head == e.Head {
				res = append(res, *ee)
			}
		}
	}
	return res, true
}

func (g *graph[K, W]) NeighbourEdges(endpoint1, endpoint2 K) ([]Edge[K, W], bool) {
	es, ok := g.GetEdge(endpoint1, endpoint2)
	if !ok {
		return nil, false
	}
	if len(es) == 0 {
		return []Edge[K, W]{}, false
	}
	return g.NeighbourEdgesByKey(es[0].Key)
}

func (g *graph[K, W]) IncidentEdges(vertex K) ([]Edge[K, W], bool) {
	if _, ok := g.vtx[vertex]; !ok {
		return nil, false
	}
	var res []Edge[K, W]
	ks, ok := g.adj.incidentEdges(vertex)
	if !ok {
		return []Edge[K, W]{}, false
	}
	res = make([]Edge[K, W], len(ks))
	for i, e := range ks {
		res[i] = *g.edges[e]
	}
	return res, true
}

func (g *graph[K, W]) SetVertexWeight(key K, weight W) bool {
	v, ok := g.vtx[key]
	if !ok {
		return false
	}
	v.Weight = weight
	return true
}

func (g *graph[K, W]) SetEdgeWeight(key K, weight W) bool {
	e, ok := g.edges[key]
	if !ok {
		return false
	}
	e.Weight = weight
	return true
	// TODO: update weight on adjlist
}
