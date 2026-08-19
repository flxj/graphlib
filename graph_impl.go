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
func ConstructGraph[K comparable, W number](digraph bool, name string, vertexes []Vertex[K, W], edges []Edge[K, W]) (Graph[K, W], error) {
	g := newGraph[K, W](digraph, name)
	for _, v := range vertexes {
		if err := g.AddVertex(v); err != nil {
			return nil, err
		}
	}
	for _, e := range edges {
		if err := g.AddEdge(e); err != nil {
			return nil, err
		}
	}
	return g, nil
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
	d, err := g.adj.minDegree()
	if err != nil {
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
	d, err := g.adj.maxDegree()
	if err != nil {
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

func (g *graph[K, W]) Property(p PropertyName) (GraphProperty[any], error) {
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
		return gp, errUnknownProperty
	}
	return gp, nil
}

func (g *graph[K, W]) AllVertexes() []Vertex[K, W] {
	vs := make([]Vertex[K, W], len(g.vtx))
	var i int
	for _, v := range g.vtx {
		/*
			vs[i] = Vertex[K,W]{
				Key:    v.Key,
				Value:  v.Value,
				Labels: v.Labels,
			}
		*/
		vs[i] = *v
		i++
	}
	return vs
}

func (g *graph[K, W]) AllEdges() []Edge[K, W] {
	es := make([]Edge[K, W], len(g.edges))
	var i int
	for _, e := range g.edges {
		/*
			es[i] = Edge[K, W]{
				Key:    e.Key,
				Head:   e.Head,
				Tail:   e.Tail,
				Value:  e.Value,
				Weight: e.Weight,
				Labels: e.Labels,
			}
		*/
		es[i] = *e
		i++
	}
	return es
}

func (g *graph[K, W]) AddVertex(v Vertex[K, W]) error {
	if _, ok := g.vtx[v.Key]; ok {
		return errVertexExists
	}
	if err := g.adj.addVertexes(v.Key); err != nil {
		return err
	}
	g.vtx[v.Key] = &v
	g.ver++
	return nil
}

func (g *graph[K, W]) RemoveVertex(key K) error {
	if _, ok := g.vtx[key]; !ok {
		return errVertexNotExists
	}
	if err := g.adj.delVertex(key); err != nil {
		return err
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
	return nil
}

func (g *graph[K, W]) AddEdge(edge Edge[K, W]) error {
	if any(edge.Key) != nil {
		if _, ok := g.edges[edge.Key]; ok {
			return errEdgeExists
		}
	} else {
		for {
			edge.Key = edgeFormat(edge.Head, edge.Tail)
			if _, ok := g.edges[edge.Key]; ok {
				break
			}
		}
	}
	if err := g.adj.addEdge(edge.Head, edge.Tail, edge.Key, edge.Weight); err != nil {
		return err
	}
	g.edges[edge.Key] = &edge
	g.ver++
	return nil
}

func (g *graph[K, W]) RemoveEdgeByKey(key K) error {
	e, ok := g.edges[key]
	if !ok {
		return errEdgeNotExists
	}
	if err := g.adj.delEdge(e.Head, e.Tail, e.Key); err != nil {
		return err
	}
	delete(g.edges, key)
	g.ver++
	return nil
}

func (g *graph[K, W]) RemoveEdge(v1, v2 K) error {
	var edges []*edge[K, W]
	for _, v := range g.edges {
		ok := (v.Head == v1 && v.Tail == v2)
		if g.adj.digraph {
			ok = ok || (v.Head == v2 && v.Tail == v1)
		}
		if ok {
			edges = append(edges, &edge[K, W]{
				key:  v.Key,
				head: v.Head,
				tail: v.Tail,
			})
		}
	}
	if err := g.adj.delEdges(edges...); err != nil {
		return err
	}
	for _, e := range edges {
		delete(g.edges, e.key)
	}
	g.ver++
	return nil
}

func (g *graph[K, W]) RemoveAllEdge() error {
	g.adj.delAllEdge()
	g.edges = make(map[K]*Edge[K, W])
	g.ver++
	return nil
}

func (g *graph[K, W]) Degree(key K) (int, error) {
	if _, ok := g.vtx[key]; !ok {
		return 0, errVertexNotExists
	}
	return g.adj.degree(key)
}

func (g *graph[K, W]) Neighbours(v K) ([]Vertex[K, W], error) {
	vs, err := g.adj.neighbours(v, false)
	if err != nil {
		return nil, err
	}
	var res []Vertex[K, W]
	for key := range vs {
		ver, ok := g.vtx[key]
		if !ok {
			return nil, fmt.Errorf("neighbour(%v) of %v not exists", key, v)
		}
		res = append(res, *ver)
	}
	return res, nil
}

func (g *graph[K, W]) GetVertex(key K) (Vertex[K, W], error) {
	v, ok := g.vtx[key]
	if !ok {
		return Vertex[K, W]{}, errVertexNotExists
	}
	return Vertex[K, W]{Key: v.Key, Value: v.Value, Labels: v.Labels}, nil
}

func (g *graph[K, W]) GetEdge(v1, v2 K) ([]Edge[K, W], error) {
	var edges []Edge[K, W]
	for _, e := range g.edges {
		ok := e.Head == v1 && e.Tail == v2
		if !g.adj.digraph {
			ok = ok || e.Head == v2 && e.Tail == v1
		}
		if ok {
			edges = append(edges, Edge[K, W]{
				Key:    e.Key,
				Head:   e.Head,
				Tail:   e.Tail,
				Value:  e.Value,
				Weight: e.Weight,
				Labels: e.Labels,
			})
		}
	}
	if len(edges) == 0 {
		return nil, errEdgeNotExists
	}
	return edges, nil
}

func (g *graph[K, W]) GetEdgeByKey(key K) (Edge[K, W], error) {
	e, ok := g.edges[key]
	if !ok {
		return Edge[K, W]{}, errEdgeNotExists
	}
	return Edge[K, W]{
		Key:    e.Key,
		Head:   e.Head,
		Tail:   e.Tail,
		Value:  e.Value,
		Weight: e.Weight,
		Labels: e.Labels,
	}, nil
}

func (g *graph[K, W]) GetVertexesByLabel(labels map[string]string) []Vertex[K, W] {
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

func (g *graph[K, W]) GetEdgesByLabel(labels map[string]string) []Edge[K, W] {
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

func (g *graph[K, W]) SetVertexValue(key K, value any) error {
	v, ok := g.vtx[key]
	if !ok {
		return errVertexNotExists
	}
	v.Value = value
	return nil
}

func (g *graph[K, W]) SetVertexLabel(key K, labelKey, labelVal string) error {
	v, ok := g.vtx[key]
	if !ok {
		return errVertexNotExists
	}
	if v.Labels == nil {
		v.Labels = make(map[string]string)
	}
	v.Labels[labelKey] = labelVal
	return nil
}

func (g *graph[K, W]) DeleteVertexLabel(key K, labelKey string) error {
	v, ok := g.vtx[key]
	if !ok {
		return errVertexNotExists
	}
	if v.Labels != nil {
		delete(v.Labels, labelKey)
	}
	return nil
}

func (g *graph[K, W]) SetEdgeValueByKey(key K, value any) error {
	e, ok := g.edges[key]
	if !ok {
		return errEdgeNotExists
	}
	e.Value = value
	return nil
}

func (g *graph[K, W]) SetEdgeLabelByKey(key K, labelKey, labelVal string) error {
	e, ok := g.edges[key]
	if !ok {
		return errEdgeNotExists
	}
	if e.Labels == nil {
		e.Labels = make(map[string]string)
	}
	e.Labels[labelKey] = labelVal
	return nil
}

func (g *graph[K, W]) DeleteEdgeLabelByKey(key K, labelKey string) error {
	e, ok := g.edges[key]
	if !ok {
		return errEdgeNotExists
	}
	if e.Labels != nil {
		delete(e.Labels, labelKey)
	}
	return nil
}

func (g *graph[K, W]) SetEdgeValue(endpoint1, endpoint2 K, value any) error {
	edges, err := g.GetEdge(endpoint1, endpoint2)
	if err != nil {
		return err
	}
	for _, ed := range edges {
		e, ok := g.edges[ed.Key]
		if !ok {
			return errEdgeNotExists
		}
		e.Value = value
	}
	return nil
}

func (g *graph[K, W]) SetEdgeLabel(endpoint1, endpoint2 K, labelKey, labelVal string) error {
	edges, err := g.GetEdge(endpoint1, endpoint2)
	if err != nil {
		return err
	}
	for _, ed := range edges {
		e, ok := g.edges[ed.Key]
		if !ok {
			return errEdgeNotExists
		}
		if e.Labels == nil {
			e.Labels = make(map[string]string)
		}
		e.Labels[labelKey] = labelVal
	}
	return nil
}

func (g *graph[K, W]) DeleteEdgeLabel(endpoint1, endpoint2 K, labelKey string) error {
	edges, err := g.GetEdge(endpoint1, endpoint2)
	if err != nil {
		return err
	}
	for _, ed := range edges {
		e, ok := g.edges[ed.Key]
		if !ok {
			return errEdgeNotExists
		}
		if e.Labels != nil {
			delete(e.Labels, labelKey)
		}
	}
	return nil
}

func (g *graph[K, W]) Clone() (Graph[K, W], error) {
	adjList := newAdjacencyLis[K, W](g.prop.digraph)
	ng := *g
	ng.vtx = make(map[K]*Vertex[K, W])
	ng.edges = make(map[K]*Edge[K, W])
	ng.adj = adjList

	for k, v := range g.vtx {
		nv := v.Clone()
		ng.vtx[k] = &nv
		if err := ng.adj.addVertexes(k); err != nil {
			return nil, err
		}
	}
	for k, v := range g.edges {
		nv := v.Clone()
		ng.edges[k] = &nv
		if err := ng.adj.addEdge(v.Head, v.Tail, v.Key, v.Weight); err != nil {
			return nil, err
		}
	}
	return &ng, nil
}

func (g *graph[K, W]) RandomVertex() (Vertex[K, W], error) {
	n := rand.Intn(len(g.vtx))
	i := 0
	for _, v := range g.vtx {
		if n == i {
			return *v, nil
		}
		i++
	}
	return Vertex[K, W]{}, errVertexNotExists
}

func (g *graph[K, W]) RandomEdge() (Edge[K, W], error) {
	n := rand.Intn(len(g.edges))
	i := 0
	for _, e := range g.edges {
		if n == i {
			return *e, nil
		}
		i++
	}
	return Edge[K, W]{}, errEdgeNotExists
}

func (g *graph[K, W]) NeighbourEdgesByKey(edge K) ([]Edge[K, W], error) {
	e, ok := g.edges[edge]
	if !ok {
		return nil, errEdgeNotExists
	}
	var res []Edge[K, W]
	for _, ee := range g.edges {
		if ee.Key != e.Key {
			if ee.Tail == e.Head || ee.Tail == e.Tail || ee.Head == e.Tail || ee.Head == e.Head {
				res = append(res, *ee)
			}
		}
	}
	return res, nil
}

func (g *graph[K, W]) NeighbourEdges(endpoint1, endpoint2 K) ([]Edge[K, W], error) {
	es, err := g.GetEdge(endpoint1, endpoint2)
	if err != nil {
		return es, nil
	}
	if len(es) == 0 {
		return []Edge[K, W]{}, nil
	}
	return g.NeighbourEdgesByKey(es[0].Key)
}

func (g *graph[K, W]) IncidentEdges(vertex K) ([]Edge[K, W], error) {
	if _, ok := g.vtx[vertex]; !ok {
		return nil, errVertexNotExists
	}
	var res []Edge[K, W]
	ks, err := g.adj.incidentEdges(vertex)
	if err != nil {
		return []Edge[K, W]{}, err
	}
	res = make([]Edge[K, W], len(ks))
	for i, e := range ks {
		res[i] = *g.edges[e]
	}
	return res, nil
}

func (g *graph[K, W]) SetVertexWeight(key K, weight W) error {
	v, ok := g.vtx[key]
	if !ok {
		return errVertexNotExists
	}
	v.Weight = weight
	return nil
}

func (g *graph[K, W]) SetEdgeWeight(key K, weight W) error {
	e, ok := g.edges[key]
	if !ok {
		return errEdgeNotExists
	}
	e.Weight = weight
	return nil
	// TODO: update weight on adjlist
}
