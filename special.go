package graphlib

type Bipartite[K comparable, V any, W number] interface {
	Graph[K, V, W]
}
