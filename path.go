package graphlib

import "math"

const (
	MaxFloatDistance = math.MaxFloat64
	MaxIntDistance   = math.MaxInt64
)

// Path represents a path on the graph,
// starting from Source and ending at Target.
// It contains edges (the key for recording edges),
// and the weighted sum of path lengths is Weight.
type Path[K comparable] struct {
	Source K
	Target K
	Edges  []K
	Weight float64
}

// Calculate the shortest path from the source vertex to the target vertex in the graph.
// If the source or target vertex does not exist, an error will be reported.
// g can be an undirected graph or a directed graph, and negative weights are allowed
// (but if negative loops are detected during the calculation process, an error will be returned).
// If the source and target are not connected, the shortest path does not exist,
// and the corresponding length is MaxFloatDistance.
func ShortestPath[K comparable, W number](g Graph[K, any, W], source K, target K) (Path[K], error) {
	p, err := g.Property(PropertyNegativeWeight)
	if err != nil {
		return Path[K]{}, err
	}
	var paths []Path[K]
	if p.Value.(bool) {
		paths, err = shortestPathBellmanFord(g, source, target, false)
	} else {
		paths, err = shortestPathDijkstra(g, source, target, false)
	}
	if err != nil {
		return Path[K]{}, err
	}
	if len(paths) == 0 {
		return Path[K]{}, errVertexNotExists
	}
	return paths[0], nil
}

// Calculate the shortest path from the source vertex to all other vertices in the graph,
// where g can be an undirected or directed graph, with negative weights allowed
// (however, if negative loops are detected during the calculation process, an error will be returned)。
func ShortestPaths[K comparable, W number](g Graph[K, any, W], source K) ([]Path[K], error) {
	p, err := g.Property(PropertyNegativeWeight)
	if err != nil {
		return nil, err
	}
	if p.Value.(bool) {
		return shortestPathBellmanFord(g, source, source, true)
	}
	return shortestPathDijkstra(g, source, source, true)
}

func getMinWeightEdge[K comparable, W number](g Graph[K, any, W], v1, v2 K) (*Edge[K, W], float64, error) {
	es, err := g.GetEdge(v1, v2)
	if err != nil {
		if !IsNotExists(err) {
			return nil, 0, err
		}
	}
	var edge *Edge[K, W]
	w := MaxFloatDistance
	for _, e := range es {
		if g.IsDigraph() {
			if e.Head == v1 && e.Tail == v2 {
				if any(e.Weight).(float64) < w {
					w = any(e.Weight).(float64)
					edge = &e
				}
			}
			continue
		}
		if any(e.Weight).(float64) < w {
			e.Head = v1
			e.Tail = v2
			w = any(e.Weight).(float64)
			edge = &e
		}
	}
	return edge, w, nil
}

// Dijkstra’s Algorithm
// Input: a positively weighted digraph (D, w) with a specified vertex r
// Output: an r-branching in D with predecessor function p, and a function L: V->R+ such that L(v)=d(r, v) for all v ∈ V
// 1: set p(v):=∅, v ∈ V, L(r):=0, and L(v):=∞, v ∈ V\{r}
// 2: while there is an uncoloured vertex u with L(u) < ∞ do
// 3:     choose such a vertex u for which L(u) is minimum
// 4:     colour u black
// 5:     for each uncoloured outneighbour v of u with L(v) > L(u) + w(u, v) do
// 6:         replace p(v) by u and L(v) by L(u) + w(u, v)
// 7:     end for
// 8: end while
// 9: return (p, L)
func shortestPathDijkstra[K comparable, W number](g Graph[K, any, W], source K, target K, all bool) ([]Path[K], error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}

	// use the map to record edges of the shortest paths.
	// if trace[v] == e,means the edges of shorest path from source to v is:
	// "edge1-....-e".
	trace := make(map[K]*Edge[K, W])
	//
	// use the slice to record vertexes of the shortest paths.
	// if prev[i] == v,means the vertexes of shorest path from source to vertexes[i] is：
	// "source-....-prev[i]-vertexes[i]".
	//prev := make([]K,len(vertexes))
	//
	dist := make([]float64, len(vertexes))
	for i, v := range vertexes {
		if v.Key == source {
			dist[i] = 0
			trace[v.Key] = nil
			continue
		}
		dist[i] = MaxFloatDistance
	}
	//
	visited := make(map[K]bool)
	for len(visited) < g.Order() {
		// select a vertex u from unvisited set whith min distance.
		distU := MaxFloatDistance
		var u K
		for i, v := range vertexes {
			if _, ok := visited[v.Key]; !ok {
				if dist[i] <= distU {
					distU = dist[i]
					u = v.Key
				}
			}
		}
		// coloured the vertex u
		visited[u] = true
		if !all && u == target {
			break
		}
		// change all unvisited vertexes distance by u
		for i, v := range vertexes {
			if _, ok := visited[v.Key]; !ok {
				e, w, err := getMinWeightEdge(g, u, v.Key)
				if err != nil {
					return nil, err
				}
				if distU < MaxFloatDistance && w < MaxFloatDistance { // to avoid overflow
					if dist[i] > distU+w {
						dist[i] = distU + w
						trace[v.Key] = e
					}
				}
			}
		}
	}
	//
	paths := []Path[K]{}
	for k, e := range trace {
		if all || (!all && k == target) {
			edges := []K{}
			for p := e; p != nil; {
				edges = append(edges, p.Key)
				p = trace[p.Head]
			}
			var w = MaxFloatDistance
			for i, d := range dist {
				if vertexes[i].Key == k {
					w = any(d).(float64)
					break
				}
			}
			paths = append(paths, Path[K]{
				Source: source,
				Target: k,
				Edges:  edges,
				Weight: w,
			})
		}
	}

	return paths, nil
}

// TODO:
func shortestPathDijkstraByMatrix[K comparable, W number](g WeightMatrix[K, W], source K, target K, all bool) ([]Path[K], error) {
	return nil, errNotImplement
}

/*
INITIALIZE-SINGLE-SOURCE(G,s)
1 for each vertex v ∈ G.V
2     v.d = ∞
3     v.p = NIL
4 s.d = 0

RELAX(u,v,w)
1 if v.d > u.d + w(u,v)
2     v.d = u.d + w(u,v)
3     v.p = u

BELLMAN-FORD(G,w,s)
1 INITIALIZE-SINGLE-SOURCE(G,s)
2 for i=1 to |G.V|-1
3     for each edge (u,v) ∈ G.E
4         RELAX(u,v,w)
5 for each edge (u,v) ∈ G.E
6     if v.d > u.d + w(u,v)
7         return FALSE
8 return TRUE
*/
func shortestPathBellmanFord[K comparable, W number](g Graph[K, any, W], source K, target K, all bool) ([]Path[K], error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}
	edges, err := g.AllEdges()
	if err != nil {
		return nil, err
	}

	// use the map to record edges of the shortest paths.
	// if trace[v] == e,means the edges of shorest path from source to v is:
	// "edge1-....-e".
	trace := make(map[K]*Edge[K, W])
	//
	dist := make(map[K]float64)
	for _, v := range vertexes {
		if v.Key == source {
			dist[v.Key] = 0
			trace[v.Key] = nil
			continue
		}
		dist[v.Key] = MaxFloatDistance
	}
	//
	for i := 0; i < g.Order(); i++ {
		for _, e := range edges {
			edge := e
			var ok bool
			var du, dv float64
			if du, ok = dist[edge.Head]; !ok {
				return nil, errVertexNotExists
			}
			if dv, ok = dist[edge.Tail]; ok {
				return nil, errVertexNotExists
			}
			uv := any(edge.Weight).(float64)
			if du < MaxFloatDistance && uv < MaxFloatDistance {
				if dv > du+uv {
					dist[edge.Tail] = du + uv
					trace[edge.Tail] = &edge
				}
			}
		}
	}
	for _, e := range edges {
		if dist[e.Tail] > dist[e.Head]+any(e.Weight).(float64) {
			return nil, errHasNegativeCycle
		}
	}
	//
	paths := []Path[K]{}
	for k, e := range trace {
		if all || (!all && k == target) {
			edges := []K{}
			for p := e; p != nil; {
				edges = append(edges, p.Key)
				p = trace[p.Head]
			}
			paths = append(paths, Path[K]{
				Source: source,
				Target: k,
				Edges:  edges,
				Weight: dist[k],
			})
		}
	}

	return paths, nil
}

/*
FLOYD-WARSHALL(W)
1 n = W.rows
2 D(0) = W
3 for k = 1 to n
4    let D(k) = d(ij)(k) be a new n*n matrix
5    for i = 1 to n
6        for j = 1 to n
7            d(ij)(k) = min{d(ij)(k-1); d(ik)(k-1)+d(kj)(k-1)}
8 return D(n)
*/
func shortestPathsFloyd[K comparable, W number](g Graph[K, any, W]) ([]Path[K], error) {
	WM, err := NewWeightMatrix(g)
	if err != nil {
		return nil, err
	}
	D := WM.Distance()
	// P is a matrix to record prev vertex of shortest path.
	// P[i][j] == v ,means the second last vertex of shortest path from i to j is v.
	// if want to find all vertexes of a path i->j, should
	// starting from P [i] [j], recursively access all intermediate vertices in the path until P [i] [v]==i.
	P := make([][]int, len(D))
	for i := range D {
		p := make([]int, len(D))
		for j := 0; j < len(D); j++ {
			if D[i][j] < MaxFloatDistance {
				p[j] = i
			}
		}
		p[i] = i
		P[i] = p
	}
	//
	for k := 0; k < len(D); k++ {
		for i := 0; i < len(D); i++ {
			for j := 0; j < len(D); j++ {
				// D[i][j] means the current shortest path from i to j.
				// if we can find a vertex k, that exists a path from i->k->j with a smaller distance.
				// then we can relax the D[i][j] and update path.
				if D[i][k] < MaxFloatDistance && D[k][j] < MaxFloatDistance {
					if D[i][j] > D[i][k]+D[k][j] {
						D[i][j] = D[i][k] + D[k][j]
						P[i][j] = P[k][j]
					}
				}
			}
		}
	}

	vs := WM.Columns()
	var paths []Path[K]
	for i := 0; i < len(D); i++ {
		var j int
		if !g.IsDigraph() {
			j = i + 1
		}
		for ; j < len(D); j++ {
			// construct the shortest path from i to j.
			var edges []K
			var t = j
			for h := P[i][j]; ; {
				// h->t
				e, _, err := getMinWeightEdge(g, vs[h], vs[t])
				if err != nil && h != t {
					return nil, err
				}
				if e != nil {
					edges = append(edges, e.Key)
				}
				if h == i {
					break
				}
				t = h
				h = P[i][h]
			}
			paths = append(paths, Path[K]{
				Source: vs[i],
				Target: vs[j],
				Weight: D[i][j],
				Edges:  edges,
			})
		}
	}

	return paths, nil
}

// Solve the shortest path between all vertex pairs in the graph
// If a pair of vertices are unreachable between them,
// the corresponding shortest path value is MaxFloatDistance.
func AllShortestPaths[K comparable, W number](g Graph[K, any, W]) ([]Path[K], error) {
	return shortestPathsFloyd(g)
}
