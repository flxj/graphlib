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
type Digraph[K comparable, W number] interface {
	Graph[K, W]
	//
	// indegree of vertex v.
	InDegree(v K) (int, error)
	//
	// outdegree of vertex v.
	OutDegree(v K) (int, error)
	//
	// The set composed of head vertexes of all v's inedges.
	InNeighbours(v K) ([]Vertex[K, W], error)
	//
	// The set composed of tail vertexes of all v's outedges.
	OutNeighbours(v K) ([]Vertex[K, W], error)
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
	Sources() ([]Vertex[K, W], error)
	//
	// All vertices with degree 0.
	Sinks() ([]Vertex[K, W], error)
	//
	DetectCycle() ([][]K, error)
	//
	// Reverse all edges in a directed graph.
	Reverse() error
}

// Create a new directed graph.
func NewDigraph[K comparable, W number](name string) Digraph[K, W] {
	return newGraph[K, W](true, name)
}

func NewDigraphFromFile[K comparable, W number](path string) (Digraph[K, W], error) {
	s, err := readFile(path)
	if err != nil {
		return nil, err
	}
	return UnmarshalDigraph[K, W](s)
}

func (g *graph[K, W]) InDegree(vertex K) (int, error) {
	return g.adj.inDegree(vertex)
}

func (g *graph[K, W]) OutDegree(vertex K) (int, error) {
	return g.adj.outDegree(vertex)
}

func (g *graph[K, W]) InNeighbours(vertex K) ([]Vertex[K, W], error) {
	vs, err := g.adj.inNeighbours(vertex, false)
	if err != nil {
		return nil, err
	}
	res := make([]Vertex[K, W], len(vs))
	var i int
	for v := range vs {
		vv, ok := g.vtx[v]
		if !ok {
			return nil, fmt.Errorf("not found neighbour %v info", v)
		}
		res[i] = *vv
		i++
	}
	return res, nil
}

func (g *graph[K, W]) OutNeighbours(vertex K) ([]Vertex[K, W], error) {
	vs, err := g.adj.outNeighbours(vertex, false)
	if err != nil {
		return nil, err
	}
	res := make([]Vertex[K, W], len(vs))
	var i int
	for v := range vs {
		vv, ok := g.vtx[v]
		if !ok {
			return nil, fmt.Errorf("not found neighbour %v info", v)
		}
		res[i] = *vv
		i++
	}
	return res, nil
}

func (g *graph[K, W]) InEdges(vertex K) ([]Edge[K, W], error) {
	es, err := g.adj.inEdges(vertex)
	if err != nil {
		return nil, err
	}
	return g.getEdges(es)
}

func (g *graph[K, W]) OutEdges(vertex K) ([]Edge[K, W], error) {
	es, err := g.adj.outEdges(vertex)
	if err != nil {
		return nil, err
	}
	return g.getEdges(es)
}

func (g *graph[K, W]) Sources() ([]Vertex[K, W], error) {
	vs, err := g.adj.sources()
	if err != nil {
		return nil, err
	}
	return g.getVertexes(vs)
}

func (g *graph[K, W]) Sinks() ([]Vertex[K, W], error) {
	vs, err := g.adj.sinks()
	if err != nil {
		return nil, err
	}
	return g.getVertexes(vs)
}

func (g *graph[K, W]) DetectCycle() ([][]K, error) {
	return nil, errNotImplement
}

func (g *graph[K, W]) Reverse() error {
	if !g.IsDigraph() {
		return nil
	}
	//
	if err := g.adj.reverse(); err != nil {
		return err
	}
	for _, e := range g.edges {
		e.Head, e.Tail = e.Tail, e.Head
	}
	return nil
}

func (g *graph[K, W]) getVertexes(vs []K) ([]Vertex[K, W], error) {
	res := make([]Vertex[K, W], len(vs))
	for i, v := range vs {
		vv, ok := g.vtx[v]
		if !ok {
			return nil, fmt.Errorf("not found neighbour %v info", v)
		}
		res[i] = *vv
	}
	return res, nil
}

func (g *graph[K, W]) getEdges(es []K) ([]Edge[K, W], error) {
	res := make([]Edge[K, W], len(es))
	for i, e := range es {
		ee, ok := g.edges[e]
		if !ok {
			return nil, fmt.Errorf("not found edge %v info", e)
		}
		res[i] = *ee
	}
	return res, nil
}
