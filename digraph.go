package graphlib

import "io"

type Digraph[K comparable, V any, W number] interface {
	Graph[K, V, W]
	InDegree(vertex K) (int, error)
	OutDegree(vertex K) (int, error)
	InNeighbours(vertex K) ([]Vertex[K, V], error)
	OutNeighbours(vertex K) ([]Vertex[K, V], error)
	Sources() ([]Vertex[K, V], error)
	Sinks() ([]Vertex[K, V], error)
	DetectCycle() ([][]K, error)
}

func NewDigraph[K comparable, V any, W number](name string) (Digraph[K, V, W], error) {
	return nil, nil
}

func NewDigraphFromFile[K comparable, V any, W number](r io.Reader) (Digraph[K, V, W], error) {
	return nil, nil
}

type digraphImpl[K comparable, V any, W number] struct {
}
