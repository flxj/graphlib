package graphlib

// TODO implement travel algorithms for graph, for example dfs,bfs...
func DFS[K comparable, W number](g Graph[K, any, W], start K, visitor func(v Vertex[K, W]) error) error {
	return errNotImplement
}

func BFS[K comparable, W number](g Graph[K, any, W], start K, visitor func(v Vertex[K, W]) error) error {
	return errNotImplement
}

func TopologicalSort[K comparable, W number](g Digraph[K, any, W]) ([]Vertex[K, W], error) {
	return nil, errNotImplement
}

func StrongComponents[K comparable, W number](g Digraph[K, any, W]) ([]Digraph[K, any, W], error) {
	return nil, errNotImplement
}
