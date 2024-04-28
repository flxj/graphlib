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

// This interface represents a directed graph.
//
// The concept of directed graphs can be referenced:
// https://mathworld.wolfram.com/DirectedGraph.html
type Digraph[K comparable, V any, W number] interface {
	Graph[K, V, W]
	//
	// indegree of vertex v.
	InDegree(v K) (int, error)
	//
	// outdegree of vertex v.
	OutDegree(v K) (int, error)
	//
	// The set composed of head vertexes of all v's inedges.
	InNeighbours(v K) ([]Vertex[K, V], error)
	//
	// The set composed of tail vertexes of all v's outedges.
	OutNeighbours(v K) ([]Vertex[K, V], error)
	//
	// All arcs with v as the tail vertex.
	// For example [a->v, b->v,...,x->v].
	InEdges(v K) ([]Edge[K, W], error)
	//
	// All arcs with v as the head vertex.
	// For example [v->a, v->b,...,v->x].
	OutEdges(v K) ([]Edge[K, W], error)
	//
	// All vertices with an in degree of 0.
	Sources() ([]Vertex[K, V], error)
	//
	// All vertices with degree 0.
	Sinks() ([]Vertex[K, V], error)
	DetectCycle() ([][]K, error)
}

// Create a new directed graph.
func NewDigraph[K comparable, V any, W number](name string) (Digraph[K, V, W], error) {
	return newGraph[K, V, W](true, name)
}

func NewDigraphFromFile[K comparable, V any, W number](path string) (Digraph[K, V, W], error) {
	s, err := readFile(path)
	if err != nil {
		return nil, err
	}
	return UnmarshalDigraph[K, V, W](s)
}

func (g *graph[K, V, W]) InDegree(vertex K) (int, error) {
	return g.adjList.inDegree(vertex)
}

func (g *graph[K, V, W]) OutDegree(vertex K) (int, error) {
	return g.adjList.outDegree(vertex)
}

func (g *graph[K, V, W]) InNeighbours(vertex K) ([]Vertex[K, V], error) {
	vs, err := g.adjList.inNeighbours(vertex,false)
	if err != nil {
		return nil, err
	}
	return g.getVertexes(vs)
}

func (g *graph[K, V, W]) OutNeighbours(vertex K) ([]Vertex[K, V], error) {
	vs, err := g.adjList.outNeighbours(vertex,false)
	if err != nil {
		return nil, err
	}
	return g.getVertexes(vs)
}

func (g *graph[K, V, W]) InEdges(vertex K) ([]Edge[K, W], error) {
	es, err := g.adjList.inEdges(vertex)
	if err != nil {
		return nil, err
	}
	return g.getEdges(es)
}

func (g *graph[K, V, W]) OutEdges(vertex K) ([]Edge[K, W], error) {
	es, err := g.adjList.outEdges(vertex)
	if err != nil {
		return nil, err
	}
	return g.getEdges(es)
}

func (g *graph[K, V, W]) Sources() ([]Vertex[K, V], error) {
	vs, err := g.adjList.sources()
	if err != nil {
		return nil, err
	}
	return g.getVertexes(vs)
}

func (g *graph[K, V, W]) Sinks() ([]Vertex[K, V], error) {
	vs, err := g.adjList.sinks()
	if err != nil {
		return nil, err
	}
	return g.getVertexes(vs)
}

func (g *graph[K, V, W]) DetectCycle() ([][]K, error) {
	return nil, errNotImplement
}

func (g *graph[K, V, W]) getVertexes(vs []K) ([]Vertex[K, V], error) {
	res := make([]Vertex[K, V], len(vs))
	for i, v := range vs {
		vv, ok := g.vertexes[v]
		if !ok {
			return nil, fmt.Errorf("not found neighbour %v info", v)
		}
		res[i] = Vertex[K, V]{
			Key:    vv.Key,
			Value:  vv.Value,
			Labels: vv.Labels,
		}
	}
	return res, nil
}

func (g *graph[K, V, W]) getEdges(es []K) ([]Edge[K, W], error) {
	res := make([]Edge[K, W], len(es))
	for i, e := range es {
		ee, ok := g.edges[e]
		if !ok {
			return nil, fmt.Errorf("not found edge %v info", e)
		}
		res[i] = Edge[K, W]{
			Key:    ee.Key,
			Head:   ee.Head,
			Tail:   ee.Tail,
			Value:  ee.Value,
			Labels: ee.Labels,
		}
	}
	return res, nil
}
