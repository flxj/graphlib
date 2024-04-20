package graphlib

// page 150: Dijkstra
// page 154: Bellman-Ford
func DigraphShortestPath[K comparable, W number](g Digraph[K, any, W], source K, target K) ([]K, W, error) {
	return nil, 0, errNotImplement
}

func DigraphShortestPaths[K comparable, W number](g Digraph[K, any, W], source K) ([][]K, []W, error) {
	return nil, nil, errNotImplement
}

// Floyd
func DigraphAllShortestPaths[K comparable, W number](g Digraph[K, any, W], source K) ([][]K, []W, error) {
	return nil, nil, errNotImplement
}

func ShortestPath[K comparable, W number](g Graph[K, any, W], source K, target K) ([]K, W, error) {
	return nil, 0, errNotImplement
}

func ShortestPaths[K comparable, W number](g Graph[K, any, W], source K) ([][]K, []W, error) {
	return nil, nil, errNotImplement
}

func AllShortestPaths[K comparable, W number](g Graph[K, any, W], source K) ([][]K, []W, error) {
	return nil, nil, errNotImplement
}
