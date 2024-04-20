package graphlib

func Contains[K comparable, N number](g Graph[K, any, N], subg Graph[K, any, N]) bool {
	return false
}

// page 46
func SpanningSubgraph[K comparable, N number](g Graph[K, any, N], edges []K) (Graph[K, any, N], error) {
	return nil, errNotImplement
}

// page 46
func SpanningSupergraph[K comparable, N number](g Graph[K, any, N], edges []*Edge[K, N]) (Graph[K, any, N], error) {
	return nil, errNotImplement
}

// page 46
func InducedSubgraph[K comparable, N number](g Graph[K, any, N], vertices []K) (Graph[K, any, N], error) {
	return nil, errNotImplement
}
