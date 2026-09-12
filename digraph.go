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

// This interface represents a directed graph.
//
// The concept of directed graphs can be referenced:
// https://mathworld.wolfram.com/DirectedGraph.html
type Digraph[K comparable, W number] interface {
	Graph[K, W]
	//
	// indegree of vertex v.
	InDegree(v K) (int, bool)
	//
	// outdegree of vertex v.
	OutDegree(v K) (int, bool)
	//
	// The set composed of head vertexes of all v's inedges.
	InNeighbours(v K) ([]Vertex[K, W], bool)
	//
	// The set composed of tail vertexes of all v's outedges.
	OutNeighbours(v K) ([]Vertex[K, W], bool)
	//
	// All arcs with v as the tail vertex.
	// For example [a->v, b->v,...,x->v].
	InEdges(v K) ([]Edge[K, W], bool)
	//
	// All arcs with v as the head vertex.
	// For example [v->a, v->b,...,v->x].
	OutEdges(v K) ([]Edge[K, W], bool)
	//
	// All vertices with an in degree of 0.
	Sources() ([]Vertex[K, W], bool)
	//
	// All vertices with degree 0.
	Sinks() ([]Vertex[K, W], bool)
	//
	DetectCycle() ([][]K, bool)
	//
	// Reverse all edges in a directed graph.
	Reverse()
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

func (g *graph[K, W]) InDegree(vertex K) (int, bool) {
	return g.adj.inDegree(vertex)
}

func (g *graph[K, W]) OutDegree(vertex K) (int, bool) {
	return g.adj.outDegree(vertex)
}

func (g *graph[K, W]) InNeighbours(vertex K) ([]Vertex[K, W], bool) {
	vs, ok := g.adj.inNeighbours(vertex, false)
	if !ok {
		return nil, false
	}
	res := make([]Vertex[K, W], len(vs))
	var i int
	for v := range vs {
		vv, ok := g.vtx[v]
		if !ok {
			return nil, false
		}
		res[i] = *vv
		i++
	}
	return res, true
}

func (g *graph[K, W]) OutNeighbours(vertex K) ([]Vertex[K, W], bool) {
	vs, ok := g.adj.outNeighbours(vertex, false)
	if !ok {
		return nil, false
	}
	res := make([]Vertex[K, W], len(vs))
	var i int
	for v := range vs {
		vv, ok := g.vtx[v]
		if !ok {
			return nil, false
		}
		res[i] = *vv
		i++
	}
	return res, true
}

func (g *graph[K, W]) InEdges(vertex K) ([]Edge[K, W], bool) {
	es, ok := g.adj.inEdges(vertex)
	if !ok {
		return nil, false
	}
	return g.getEdges(es)
}

func (g *graph[K, W]) OutEdges(vertex K) ([]Edge[K, W], bool) {
	es, ok := g.adj.outEdges(vertex)
	if !ok {
		return nil, false
	}
	return g.getEdges(es)
}

func (g *graph[K, W]) Sources() ([]Vertex[K, W], bool) {
	vs, ok := g.adj.sources()
	if !ok {
		return nil, false
	}
	return g.getVertexes(vs)
}

func (g *graph[K, W]) Sinks() ([]Vertex[K, W], bool) {
	vs, ok := g.adj.sinks()
	if !ok {
		return nil, false
	}
	return g.getVertexes(vs)
}

func (g *graph[K, W]) DetectCycle() ([][]K, bool) {
	return nil, false
}

func (g *graph[K, W]) Reverse() {
	if !g.IsDigraph() {
		return
	}
	g.adj.reverse()
	for _, e := range g.edges {
		e.Head, e.Tail = e.Tail, e.Head
	}
}

func (g *graph[K, W]) getVertexes(vs []K) ([]Vertex[K, W], bool) {
	res := make([]Vertex[K, W], len(vs))
	for i, v := range vs {
		vv, ok := g.vtx[v]
		if !ok {
			return nil, false
		}
		res[i] = *vv
	}
	return res, true
}

func (g *graph[K, W]) getEdges(es []K) ([]Edge[K, W], bool) {
	res := make([]Edge[K, W], len(es))
	for i, e := range es {
		ee, ok := g.edges[e]
		if !ok {
			return nil, false
		}
		res[i] = *ee
	}
	return res, true
}
