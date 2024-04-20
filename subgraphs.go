package graphlib

func Contains[K comparable, N number](g Graph[K, any, N], subg Graph[K, any, N]) bool {
	return false
}

func SpanningSubgraph[K comparable, N number](g Graph[K, any, N], edges []K) (Graph[K, any, N], error) {
	return nil, nil
}

func SpanningSupergraph[K comparable, N number](g Graph[K, any, N], edges []*Edge[K, N]) (Graph[K, any, N], error) {
	return nil, nil
}

func InducedSubgraph[K comparable, N number](g Graph[K, any, N], vertices []K) (Graph[K, any, N], error) {
	return nil, nil
}
