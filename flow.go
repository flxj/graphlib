package graphlib

// page 163
func MaxNetworkFlow[K comparable, W number](g Digraph[K, any, W]) (W, error) {
	return 0, errNotImplement
}

func MaxMatching[K comparable, W number](g Bipartite[K, any, W]) ([]Edge[K, W], error) {
	return nil, errNotImplement
}

func PerfectMatching[K comparable, W number](g Bipartite[K, any, W]) ([]Edge[K, W], error) {
	return nil, errNotImplement
}
