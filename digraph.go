package graphlib

import (
	"fmt"
	"io"
)

type Digraph[K comparable, V any, W number] interface {
	Graph[K, V, W]
	InDegree(vertex K) (int, error)
	OutDegree(vertex K) (int, error)
	InNeighbours(vertex K) ([]Vertex[K, V], error)
	OutNeighbours(vertex K) ([]Vertex[K, V], error)
	InEdges(vertex K) ([]Edge[K, W], error)
	OutEdges(vertex K) ([]Edge[K, W], error)
	Sources() ([]Vertex[K, V], error)
	Sinks() ([]Vertex[K, V], error)
	DetectCycle() ([][]K, error)
}

func NewDigraph[K comparable, V any, W number](name string) (Digraph[K, V, W], error) {
	return newGraph[K, V, W](true, name)
}

func NewDigraphFromFile[K comparable, V any, W number](r io.Reader) (Digraph[K, V, W], error) {
	return nil, errNotImplement
}

func (g *graph[K, V, W]) InDegree(vertex K) (int, error) {
	return g.adjList.inDegree(vertex)
}

func (g *graph[K, V, W]) OutDegree(vertex K) (int, error) {
	return g.adjList.outDegree(vertex)
}

func (g *graph[K, V, W]) InNeighbours(vertex K) ([]Vertex[K, V], error) {
	vs, err := g.adjList.inNeighbours(vertex)
	if err != nil {
		return nil, err
	}
	return g.getVertexes(vs)
}

func (g *graph[K, V, W]) OutNeighbours(vertex K) ([]Vertex[K, V], error) {
	vs, err := g.adjList.outNeighbours(vertex)
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
