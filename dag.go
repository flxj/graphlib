package graphlib

import "io"

type Digraph[K comparable, V any, W number] interface {
	Graph[K, V, W]
	InDegree(vertex K) int
	OutDegree(vertex K) int
	InNeighbours(vertex K) []*Vertex[K, V]
	OutNeighbours(vertex K) []*Vertex[K, V]
	Sources() ([]*Vertex[K, V], error)
	Sinks() ([]*Vertex[K, V], error)
}

func NewDigraph[K comparable, V any, W number]() (Digraph[K, V, W], error) {
	return nil, nil
}
func NewDigraphFromFile[K comparable, V any, W number](r io.Reader) (Digraph[K, V, W], error) {
	return nil, nil
}
