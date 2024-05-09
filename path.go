/*
	Copyright (C) 2023 flxj(https://github.com/flxj)

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package graphlib

// Path represents a path on the graph,
// starting from Source and ending at Target.
// It contains edges (the key for recording edges),
// and the weighted sum of path lengths is Weight.
type Path[K comparable, W number] struct {
	Source K
	Target K
	Edges  []K
	Weight W
}

// Calculate the shortest path from the source vertex to the target vertex in the graph.
// If the source or target vertex does not exist, an error will be reported.
// g can be an undirected graph or a directed graph, and negative weights are allowed
// (but if negative loops are detected during the calculation process, an error will be returned).
// If the source and target are not connected, the shortest path does not exist,
// and the corresponding length is MaxDistance.
func ShortestPath[K comparable, V any, W number](g Graph[K, V, W], source K, target K) (Path[K, W], error) {
	p, err := g.Property(PropertyNegativeWeight)
	if err != nil {
		return Path[K, W]{}, err
	}
	var paths []Path[K, W]
	if p.Value.(bool) {
		paths, err = shortestPathBellmanFord(g, source, target, false)
	} else {
		paths, err = shortestPathDijkstraWithPQ(g, source, target, false)
	}
	if err != nil {
		return Path[K, W]{}, err
	}
	if len(paths) == 0 {
		return Path[K, W]{}, errVertexNotExists
	}
	return paths[0], nil
}

// Calculate the shortest path from the source vertex to all other vertices in the graph,
// where g can be an undirected or directed graph, with negative weights allowed
// (however, if negative loops are detected during the calculation process, an error will be returned)。
func ShortestPaths[K comparable, V any, W number](g Graph[K, V, W], source K) ([]Path[K, W], error) {
	p, err := g.Property(PropertyNegativeWeight)
	if err != nil {
		return nil, err
	}
	if p.Value.(bool) {
		return shortestPathBellmanFord(g, source, source, true)
	}
	return shortestPathDijkstraWithPQ(g, source, source, true)
}

func getMinWeightEdge[K comparable, V any, W number](g Graph[K, V, W], v1, v2 K) (*Edge[K, W], W, error) {
	es, err := g.GetEdge(v1, v2)
	if err != nil {
		if !IsNotExists(err) {
			return nil, 0, err
		}
	}
	var edge *Edge[K, W]
	var n W
	w := getMaxValue(n)
	for _, e := range es {
		if g.IsDigraph() {
			if e.Head == v1 && e.Tail == v2 {
				if e.Weight < w {
					w = e.Weight
					edge = &e
				}
			}
			continue
		}
		if e.Weight < w {
			e.Head = v1
			e.Tail = v2
			w = e.Weight
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
func shortestPathDijkstra[K comparable, V any, W number](g Graph[K, V, W], source K, target K, all bool) ([]Path[K, W], error) {
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
	var n W
	maxDist := getMaxValue(n)
	//
	dist := make([]W, len(vertexes))
	for i, v := range vertexes {
		if v.Key == source {
			dist[i] = 0
			trace[v.Key] = nil
			continue
		}
		dist[i] = maxDist
	}
	//
	visited := make(map[K]bool)
	for len(visited) < g.Order() {
		// select a vertex u from unvisited set whith min distance.
		distU := maxDist
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
				if distU < maxDist && w < maxDist { // to avoid overflow
					if dist[i] > distU+w {
						dist[i] = distU + w
						trace[v.Key] = e
					}
				}
			}
		}
	}
	//
	paths := []Path[K, W]{}
	for k, e := range trace {
		if all || (!all && k == target) {
			edges := []K{}
			for p := e; p != nil; {
				edges = append(edges, p.Key)
				p = trace[p.Head]
			}
			var w = getMaxValue(n)
			for i, d := range dist {
				if vertexes[i].Key == k {
					w = d
					break
				}
			}
			paths = append(paths, Path[K, W]{
				Source: source,
				Target: k,
				Edges:  edges,
				Weight: w,
			})
		}
	}

	return paths, nil
}

// Implement Dijkstra algorithm using priority queue.
func shortestPathDijkstraWithPQ[K comparable, V any, W number](g Graph[K, V, W], source K, target K, all bool) ([]Path[K, W], error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}
	trace := make(map[K]*Edge[K, W])
	unvisited := make(map[K]bool)
	//
	var n W
	maxDist := getMaxValue(n)
	//
	dist := newPriorityQueue[K, int, W](func(p1, p2 W) bool { return p1 < p2 })
	for _, v := range vertexes {
		p := maxDist
		if v.Key == source {
			p = 0.0
		}
		dist.Push(v.Key, 0, p)
		unvisited[v.Key] = true
	}

	for len(unvisited) > 0 {
		// select a vertex u from unvisited set whith min distance.
		u, _, distU, _ := dist.Pop()

		// coloured the vertex u
		delete(unvisited, u)
		if !all && u == target {
			break
		}
		// change all unvisited vertexes distance by u
		for v := range unvisited {
			e, w, err := getMinWeightEdge(g, u, v)
			if err != nil {
				return nil, err
			}
			if distU < maxDist && w < maxDist {
				if dist.Get(v) > distU+w {
					dist.Update(v, distU+w)
					trace[v] = e
				}
			}
		}
	}
	//
	paths := []Path[K, W]{}
	for k, e := range trace {
		if all || (!all && k == target) {
			var w W
			edges := []K{}
			for p := e; p != nil; {
				edges = append(edges, p.Key)
				w += p.Weight
				p = trace[p.Head]
			}
			if len(edges) == 0 {
				w = maxDist
			}
			paths = append(paths, Path[K, W]{
				Source: source,
				Target: k,
				Edges:  edges,
				Weight: w,
			})
		}
	}

	return paths, nil
}

func shortestPathDijkstraByMatrix[K comparable, W number](g WeightMatrix[K, W], source K, target K, all bool) ([][]*Edge[K, W], error) {
	vertexes := g.Columns()
	var n W
	maxDist := getMaxValue(n)

	w := g.Weight(maxDist)
	trace := make(map[K]*Edge[K, W])
	unvisited := make(map[int]bool)
	//
	dist := newPriorityQueue[int, int, W](func(p1, p2 W) bool { return p1 < p2 })
	for i, v := range vertexes {
		p := maxDist
		if v == source {
			p = 0.0
		}
		dist.Push(i, 0, p)
		unvisited[i] = true
	}

	for len(unvisited) > 0 {
		// select a vertex u from unvisited set whith min distance.
		u, _, distU, _ := dist.Pop()

		// coloured the vertex u
		delete(unvisited, u)
		if !all && vertexes[u] == target {
			break
		}
		// change all unvisited vertexes distance by u
		for v := range unvisited {
			if distU < maxDist && w[u][v] < maxDist {
				if dist.Get(v) > distU+w[u][v] {
					dist.Update(v, distU+w[u][v])
					trace[vertexes[v]] = &Edge[K, W]{
						Head:   vertexes[u],
						Tail:   vertexes[v],
						Weight: any(w[u][v]).(W),
					}
				}
			}
		}
	}
	//
	paths := [][]*Edge[K, W]{}
	for k, e := range trace {
		if all || (!all && k == target) {
			edges := []*Edge[K, W]{}
			for p := e; p != nil; {
				edges = append(edges, p)
				p = trace[p.Head]
			}
			paths = append(paths, edges)
		}
	}

	return paths, nil
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
func shortestPathBellmanFord[K comparable, V any, W number](g Graph[K, V, W], source K, target K, all bool) ([]Path[K, W], error) {
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
	var n W
	maxDist := getMaxValue(n)
	//
	dist := make(map[K]W)
	for _, v := range vertexes {
		if v.Key == source {
			dist[v.Key] = 0
			trace[v.Key] = nil
			continue
		}
		dist[v.Key] = maxDist
	}
	//
	for i := 0; i < g.Order(); i++ {
		for _, e := range edges {
			edge := e
			var ok bool
			var du, dv W
			if du, ok = dist[edge.Head]; !ok {
				return nil, errVertexNotExists
			}
			if dv, ok = dist[edge.Tail]; ok {
				return nil, errVertexNotExists
			}
			uv := edge.Weight
			if du < maxDist && uv < maxDist {
				if dv > du+uv {
					dist[edge.Tail] = du + uv
					trace[edge.Tail] = &edge
				}
			}
		}
	}
	for _, e := range edges {
		if dist[e.Tail] > dist[e.Head]+e.Weight {
			return nil, errHasNegativeCycle
		}
	}
	//
	paths := []Path[K, W]{}
	for k, e := range trace {
		if all || (!all && k == target) {
			edges := []K{}
			for p := e; p != nil; {
				edges = append(edges, p.Key)
				p = trace[p.Head]
			}
			paths = append(paths, Path[K, W]{
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
func shortestPathsFloyd[K comparable, V any, W number](g Graph[K, V, W]) ([]Path[K, W], error) {
	WM, err := NewWeightMatrix(g)
	if err != nil {
		return nil, err
	}
	var n W
	maxDist := getMaxValue(n)
	//
	D := WM.Weight(maxDist)
	// P is a matrix to record prev vertex of shortest path.
	// P[i][j] == v ,means the second last vertex of shortest path from i to j is v.
	// if want to find all vertexes of a path i->j, should
	// starting from P [i] [j], recursively access all intermediate vertices in the path until P [i] [v]==i.
	P := make([][]int, len(D))
	for i := range D {
		p := make([]int, len(D))
		for j := 0; j < len(D); j++ {
			if D[i][j] < maxDist {
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
				if D[i][k] < maxDist && D[k][j] < maxDist {
					if D[i][j] > D[i][k]+D[k][j] {
						D[i][j] = D[i][k] + D[k][j]
						P[i][j] = P[k][j]
					}
				}
			}
		}
	}

	vs := WM.Columns()
	var paths []Path[K, W]
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
			paths = append(paths, Path[K, W]{
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
// the corresponding shortest path value is MaxDistance.
func AllShortestPaths[K comparable, V any, W number](g Graph[K, V, W]) ([]Path[K, W], error) {
	return shortestPathsFloyd(g)
}

//
func CountCycles[K comparable, V any, W number](g Graph[K, V, W], length int) (int, error) {
	if length <= 0 {
		return 0, nil
	} else if length == 1 {
		edges, err := g.AllEdges()
		if err != nil {
			return 0, err
		}
		var count int
		for _, e := range edges {
			if e.Head == e.Tail {
				count++
			}
		}
		return count, nil
	} else if length == 2 {
		// find multi number of a edge.
		edges, err := g.AllEdges()
		if err != nil {
			return 0, err
		}
		counts := make(map[K]map[K]int)
		for _, e := range edges {
			ts, ok1 := counts[e.Head]
			if !ok1 {
				counts[e.Head] = make(map[K]int)
			}
			counts[e.Head][e.Tail] = ts[e.Tail] + 1

			hs, ok2 := counts[e.Tail]
			if !ok2 {
				counts[e.Tail] = make(map[K]int)
			}
			counts[e.Tail][e.Head] = hs[e.Head] + 1
		}
		var count int
		for _, v := range counts {
			for _, n := range v {
				if n >= 2 {
					count += n * (n - 1) / 2
				}
			}
		}

		return count / 2, nil
	}
	var count int
	visited := make(map[K]bool)

	vertexes, err := g.AllVertexes()
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(vertexes)-length+1; i++ {
		cp, err := countPaths(g, vertexes[i].Key, vertexes[i].Key, length-1, visited)
		if err != nil {
			return 0, err
		}
		count += cp
		visited[vertexes[i].Key] = true
	}
	return count / 2, nil
}

// n is the vertexes number in search path.
func countPaths[K comparable, V any, W number](g Graph[K, V, W], start, end K, n int, visited map[K]bool) (int, error) {
	visited[start] = true
	defer func() { delete(visited, start) }()

	if n == 0 {
		// check edge start->end exists.
		if _, err := g.GetEdge(start, end); err != nil {
			if !IsNotExists(err) {
				return 0, err
			}
			return 0, nil
		}
		return 1, nil
	}

	var count int
	vs, err := g.Neighbours(start)
	if err != nil {
		return 0, err
	}
	for _, v := range vs {
		if _, ok := visited[v.Key]; !ok {
			cp, err := countPaths(g, v.Key, end, n-1, visited)
			if err != nil {
				return 0, err
			}
			count += cp
		}
	}

	return count, nil
}
