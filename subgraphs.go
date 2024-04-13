package graphlib

func Contains[K comparable](g Graph[K, any, any], subg Graph[K, any, any]) bool {
	return false
}

func SpanningSubgraph[K comparable](g Graph[K, any, any], edges []K) (Graph[K, any, any], error) {
	return nil, nil
}

func SpanningSupergraph[K comparable, V any, W number | any](g Graph[K, V, W], edges []*Edge[K, W]) (Graph[K, any, any], error) {
	return nil, nil
}

func InducedSubgraph[K comparable](g Graph[K, any, any], vertices []K) (Graph[K, any, any], error) {
	return nil, nil
}
