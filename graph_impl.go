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

type basicPropertySet[T any] struct {
	digraph             bool
	acyclic             property[T] // no cycle and no loop
	simple              property[T] // no loop and no multi edge
	regular             property[T] // every vertex has same order
	connected           property[T] // for digraph, which means strong connection
	forest              property[T]
	loop                property[T]
	negativeWeight      property[T]
	unilateralConnected property[T]
}

// graph default implement base on adjacency list.
type graph[K comparable, V any, W number] struct {
	version    int // start from 1
	name       string
	properties basicPropertySet[bool] // version start from 0
	minDe      property[int]
	maxDe      property[int]
	avgDe      property[float64]
	vertexes   map[K]*Vertex[K, V]
	edges      map[K]*Edge[K, W]
	adjList    *adjacencyList[K, W]
}

func newGraph[K comparable, V any, W number](digraph bool, name string) (*graph[K, V, W], error) {
	g := &graph[K, V, W]{
		version:  1,
		name:     name,
		vertexes: make(map[K]*Vertex[K, V]),
		edges:    make(map[K]*Edge[K, W]),
	}
	g.properties.digraph = digraph

	adj, err := newAdjacencyLis[K, W](digraph)
	if err != nil {
		return nil, err
	}
	g.adjList = adj

	return g, nil
}

// Create a new graph.
func NewGraph[K comparable, V any, W number](digraph bool, name string) (Graph[K, V, W], error) {
	return newGraph[K, V, W](digraph, name)
}

// Create a new undirected graph
func NewUnDigraph[K comparable, V any, W number](name string) (Graph[K, V, W], error) {
	return newGraph[K, V, W](false, name)
}

// Load graph from json or yaml file.
func NewGraphFromFile[K comparable, V any, W number](path string) (Graph[K, V, W], error) {
	s, err := readFile(path)
	if err != nil {
		return nil, err
	}
	return UnmarshalGraph[K, V, W](s)
}

// Create a graph using vertex and edge sets.
func ConstructGraph[K comparable, V any, W number](digraph bool, name string, vertexes []Vertex[K, V], edges []Edge[K, W]) (Graph[K, V, W], error) {
	g, err := newGraph[K, V, W](digraph, name)
	if err != nil {
		return nil, err
	}
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

func (g *graph[K, V, W]) Name() string {
	return g.name
}

func (g *graph[K, V, W]) SetName(name string) {
	g.name = name
}

func (g *graph[K, V, W]) IsDigraph() bool {
	return g.adjList.digraph
}

func (g *graph[K, V, W]) IsSimple() bool {
	if g.properties.simple.version == g.version {
		return g.properties.simple.value
	}
	//
	p, _ := g.adjList.property(simple)
	p.version = g.version
	g.properties.simple = p

	return p.value
}

func (g *graph[K, V, W]) HasNegativeWeight() bool {
	if g.properties.negativeWeight.version == g.version {
		return g.properties.negativeWeight.value
	}
	p, _ := g.adjList.property(negativeWeight)
	p.version = g.version
	g.properties.negativeWeight = p

	return p.value
}

func (g *graph[K, V, W]) IsRegular() bool {
	if g.properties.regular.version == g.version {
		return g.properties.regular.value
	}
	p, _ := g.adjList.property(regular)
	p.version = g.version
	g.properties.regular = p

	return p.value
}

func (g *graph[K, V, W]) IsAcyclic() bool {
	if g.properties.acyclic.version == g.version {
		return g.properties.acyclic.value
	}
	p, _ := g.adjList.property(acyclic)
	p.version = g.version
	g.properties.acyclic = p

	return p.value
}
func (g *graph[K, V, W]) IsConnected(unidirectional bool) bool {
	if unidirectional && g.IsDigraph() {
		if g.properties.unilateralConnected.version == g.version {
			return g.properties.unilateralConnected.value
		}
		p, _ := g.adjList.property(unilateralConnected)
		p.version = g.version
		g.properties.unilateralConnected = p

		return p.value
	}
	if g.properties.connected.version == g.version {
		return g.properties.connected.value
	}
	p, _ := g.adjList.property(connected)
	p.version = g.version
	g.properties.connected = p

	return p.value
}

func (g *graph[K, V, W]) IsCompleted() bool {
	if g.IsSimple() {
		return g.MinDegree() == g.Order()-1 // TODO for bipartite graph
	}
	return false
}

func (g *graph[K, V, W]) IsTree() bool {
	return g.IsConnected(false) && g.IsForest()
}

func (g *graph[K, V, W]) IsForest() bool {
	if g.properties.forest.version == g.version {
		return g.properties.forest.value
	}
	p, _ := g.adjList.property(forest)
	p.version = g.version
	g.properties.forest = p

	return p.value
}

func (g *graph[K, V, W]) HasLoop() bool {
	if g.properties.loop.version == g.version {
		return g.properties.loop.value
	}
	p, _ := g.adjList.property(loop)
	p.version = g.version
	g.properties.loop = p

	return p.value
}

func (g *graph[K, V, W]) Order() int {
	return len(g.vertexes)
}

func (g *graph[K, V, W]) Size() int {
	return len(g.edges)
}

func (g *graph[K, V, W]) MinDegree() int {
	if g.minDe.version == g.version {
		return g.minDe.value
	}
	d, err := g.adjList.minDegree()
	if err != nil {
		return -1
	}
	g.minDe.version = g.version
	g.minDe.value = d
	return d
}

func (g *graph[K, V, W]) MaxDegree() int {
	if g.maxDe.version == g.version {
		return g.maxDe.value
	}
	d, err := g.adjList.maxDegree()
	if err != nil {
		return -1
	}
	g.maxDe.version = g.version
	g.maxDe.value = d
	return d
}

func (g *graph[K, V, W]) AvgDegree() float64 {
	if g.avgDe.version == g.version {
		return g.avgDe.value
	}
	var avg float64
	if g.Order() != 0 {
		avg = float64(2*g.Size()) / float64(g.Order())
	}
	g.avgDe.version = g.version
	g.avgDe.value = avg
	return avg
}

func (g *graph[K, V, W]) Property(p PropertyName) (GraphProperty[any], error) {
	gp := GraphProperty[any]{Name: p}
	switch p {
	case PropertyDigraph:
		gp.Value = g.IsDigraph()
	case PropertyAcyclic:
		gp.Value = g.IsAcyclic()
	case PropertySimple:
		gp.Value = g.IsSimple()
	case PropertyRegular:
		gp.Value = g.IsRegular()
	case PropertyConnected:
		gp.Value = g.IsConnected(false)
	case PropertyUnilateralConnected:
		gp.Value = g.IsConnected(true)
	case PropertyForest:
		gp.Value = g.IsForest()
	case PropertyLoop:
		gp.Value = g.HasLoop()
	case PropertyCompleted:
		gp.Value = g.IsCompleted()
	case PropertyTree:
		gp.Value = g.IsTree()
	case PropertyNegativeWeight:
		gp.Value = g.HasNegativeWeight()
	case PropertyGraphName:
		gp.Value = g.Name()
	case PropertyOrder:
		gp.Value = g.Order()
	case PropertySize:
		gp.Value = g.Size()
	case PropertyMaxDegree:
		gp.Value = g.MaxDegree()
	case PropertyMinDegree:
		gp.Value = g.MinDegree()
	case PropertyAvgDegree:
		gp.Value = g.AvgDegree()
	default:
		return gp, errUnknownProperty
	}
	return gp, nil
}

func (g *graph[K, V, W]) AllVertexes() ([]Vertex[K, V], error) {
	vs := make([]Vertex[K, V], len(g.vertexes))
	var i int
	for _, v := range g.vertexes {
		vs[i] = Vertex[K, V]{
			Key:    v.Key,
			Value:  v.Value,
			Labels: v.Labels,
		}
		i++
	}
	return vs, nil
}

func (g *graph[K, V, W]) AllEdges() ([]Edge[K, W], error) {
	es := make([]Edge[K, W], len(g.edges))
	var i int
	for _, e := range g.edges {
		es[i] = Edge[K, W]{
			Key:    e.Key,
			Head:   e.Head,
			Tail:   e.Tail,
			Value:  e.Value,
			Weight: e.Weight,
			Labels: e.Labels,
		}
		i++
	}
	return es, nil
}

func (g *graph[K, V, W]) AddVertex(v Vertex[K, V]) error {
	if _, ok := g.vertexes[v.Key]; ok {
		return errVertexExists
	}
	if err := g.adjList.addVertexes(v.Key); err != nil {
		return err
	}
	g.vertexes[v.Key] = &v
	g.version++
	return nil
}

func (g *graph[K, V, W]) RemoveVertex(key K) error {
	if _, ok := g.vertexes[key]; !ok {
		return errVertexNotExists
	}
	if err := g.adjList.delVertex(key); err != nil {
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
	delete(g.vertexes, key)
	g.version++
	return nil
}

func (g *graph[K, V, W]) AddEdge(edge Edge[K, W]) error {
	if any(edge.Key) == nil {
		edge.Key = edgeFormat(edge.Head, edge.Tail)
	}
	if _, ok := g.edges[edge.Key]; ok {
		return errEdgeExists
	}
	if err := g.adjList.addEdge(edge.Head, edge.Tail, edge.Key, edge.Weight); err != nil {
		return err
	}
	g.edges[edge.Key] = &edge
	g.version++
	return nil
}

func (g *graph[K, V, W]) RemoveEdgeByKey(key K) error {
	e, ok := g.edges[key]
	if !ok {
		return errEdgeNotExists
	}
	if err := g.adjList.delEdge(e.Head, e.Tail, e.Key); err != nil {
		return err
	}
	delete(g.edges, key)
	g.version++
	return nil
}

func (g *graph[K, V, W]) RemoveEdge(v1, v2 K) error {
	var edges []*edge[K, W]
	for _, v := range g.edges {
		ok := (v.Head == v1 && v.Tail == v2)
		if g.adjList.digraph {
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
	if err := g.adjList.delEdges(edges...); err != nil {
		return err
	}
	for _, e := range edges {
		delete(g.edges, e.key)
	}
	g.version++
	return nil
}

func (g *graph[K, V, W]) Degree(key K) (int, error) {
	if _, ok := g.vertexes[key]; !ok {
		return 0, errVertexNotExists
	}
	return g.adjList.degree(key)
}

func (g *graph[K, V, W]) Neighbours(v K) ([]Vertex[K, V], error) {
	vs, err := g.adjList.neighbours(v, false)
	if err != nil {
		return nil, err
	}
	var res []Vertex[K, V]
	for _, key := range vs {
		ver, ok := g.vertexes[key]
		if !ok {
			return nil, fmt.Errorf("neighbour(%v) of %v not exists", key, v)
		}
		res = append(res, Vertex[K, V]{
			Key:    key,
			Value:  ver.Value,
			Labels: ver.Labels,
		})
	}
	return res, nil
}

func (g *graph[K, V, W]) GetVertex(key K) (Vertex[K, V], error) {
	v, ok := g.vertexes[key]
	if !ok {
		return Vertex[K, V]{}, errVertexNotExists
	}
	return Vertex[K, V]{Key: v.Key, Value: v.Value, Labels: v.Labels}, nil
}

func (g *graph[K, V, W]) GetEdge(v1, v2 K) ([]Edge[K, W], error) {
	var edges []Edge[K, W]
	for _, e := range g.edges {
		ok := e.Head == v1 && e.Tail == v2
		if !g.adjList.digraph {
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

func (g *graph[K, V, W]) GetEdgeByKey(key K) (Edge[K, W], error) {
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

func (g *graph[K, V, W]) GetVertexesByLabel(labels map[string]string) ([]Vertex[K, V], error) {
	var ves []Vertex[K, V]
	if labels != nil {
		for _, vertex := range g.vertexes {
			if vertex.Labels != nil {
				match := true
				for k, v := range labels {
					l, ok := vertex.Labels[k]
					if !ok || l != v {
						match = false
						break
					}
				}
				if match {
					ves = append(ves, Vertex[K, V]{
						Key:    vertex.Key,
						Value:  vertex.Value,
						Labels: vertex.Labels,
					})
				}
			}
		}
	}
	return ves, nil
}

func (g *graph[K, V, W]) GetEdgesByLabel(labels map[string]string) ([]Edge[K, W], error) {
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
		}
	}
	return edges, nil
}

func (g *graph[K, V, W]) SetVertexValue(key K, value V) error {
	v, ok := g.vertexes[key]
	if !ok {
		return errVertexNotExists
	}
	v.Value = value
	return nil
}

func (g *graph[K, V, W]) SetVertexLabel(key K, labelKey, labelVal string) error {
	v, ok := g.vertexes[key]
	if !ok {
		return errVertexNotExists
	}
	if v.Labels == nil {
		v.Labels = make(map[string]string)
	}
	v.Labels[labelKey] = labelVal
	return nil
}

func (g *graph[K, V, W]) DeleteVertexLabel(key K, labelKey string) error {
	v, ok := g.vertexes[key]
	if !ok {
		return errVertexNotExists
	}
	if v.Labels != nil {
		delete(v.Labels, labelKey)
	}
	return nil
}

func (g *graph[K, V, W]) SetEdgeValueByKey(key K, value any) error {
	e, ok := g.edges[key]
	if !ok {
		return errEdgeNotExists
	}
	e.Value = value
	return nil
}

func (g *graph[K, V, W]) SetEdgeLabelByKey(key K, labelKey, labelVal string) error {
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

func (g *graph[K, V, W]) DeleteEdgeLabelByKey(key K, labelKey string) error {
	e, ok := g.edges[key]
	if !ok {
		return errEdgeNotExists
	}
	if e.Labels != nil {
		delete(e.Labels, labelKey)
	}
	return nil
}

func (g *graph[K, V, W]) SetEdgeValue(endpoint1, endpoint2 K, value any) error {
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

func (g *graph[K, V, W]) SetEdgeLabel(endpoint1, endpoint2 K, labelKey, labelVal string) error {
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

func (g *graph[K, V, W]) DeleteEdgeLabel(endpoint1, endpoint2 K, labelKey string) error {
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

func (g *graph[K, V, W]) Clone() (Graph[K, V, W], error) {
	adjList, err := newAdjacencyLis[K, W](g.properties.digraph)
	if err != nil {
		return nil, err
	}
	ng := *g
	ng.vertexes = make(map[K]*Vertex[K, V])
	ng.edges = make(map[K]*Edge[K, W])
	ng.adjList = adjList

	for k, v := range g.vertexes {
		ng.vertexes[k] = v.Clone()
		if err = ng.adjList.addVertexes(k); err != nil {
			return nil, err
		}
	}
	for k, v := range g.edges {
		ng.edges[k] = v.Clone()
		if err = ng.adjList.addEdge(v.Head, v.Tail, v.Key, v.Weight); err != nil {
			return nil, err
		}
	}
	return &ng, nil
}

func (g *graph[K, V, W]) RandomVertex() (Vertex[K, V], error) {
	n := rand.Intn(len(g.vertexes))
	i := 0
	for _, v := range g.vertexes {
		if n == i {
			return *v, nil
		}
		i++
	}
	return Vertex[K, V]{}, errVertexNotExists
}

func (g *graph[K, V, W]) RandomEdge() (Edge[K, W], error) {
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

func (g *graph[K, V, W]) NeighbourEdgesByKey(edge K) ([]Edge[K, W], error) {
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

func (g *graph[K, V, W]) NeighbourEdges(endpoint1, endpoint2 K) ([]Edge[K, W], error) {
	es, err := g.GetEdge(endpoint1, endpoint2)
	if err != nil {
		return es, nil
	}
	if len(es) == 0 {
		return []Edge[K, W]{}, nil
	}
	return g.NeighbourEdgesByKey(es[0].Key)
}

func (g *graph[K, V, W]) IncidentEdges(vertex K) ([]Edge[K, W], error) {
	if _, ok := g.vertexes[vertex]; !ok {
		return nil, errVertexNotExists
	}

	var res []Edge[K, W]
	for _, e := range g.edges {
		if e.Tail == vertex || e.Head == vertex {
			res = append(res, *e)
		}
	}

	return res, nil
}
