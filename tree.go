package graphlib

import "sort"

// each uncoloured vertex v is assigned a provisional cost c(v).
// This is the least weight of an edge linking v to some black vertex u.
//
// Jarnı́k–Prim Algorithm
// Input: a weighted connected graph(G, w)
// Output: an optimal tree T of G with predecessor function p, and its weight w(T)
// 1: set p(v) := ∅ and c(v) := ∞, v ∈ V , and w(T) := 0
// 2: choose any vertex r (as root)
// 3: replace c(r) by 0
// 4: while there is an uncoloured vertex do
// 5:     choose such a vertex u of minimum cost c(u)
// 6:     colour u black
// 7:     for each uncoloured vertex v such that w(uv) < c(v) do
// 8:         replace p(v) by u and c(v) by w(uv)
// 9:     end for
// 10:    replace w(T) by w(T) + c(u)
// 11: end while
// 12: return (p, w(T))
func mstPrim[K comparable, W number](g Graph[K, any, W]) ([]K, []Edge[K, W], float64, error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, nil, 0.0, err
	}
	if len(vertexes) == 0 {
		return nil, nil, 0.0, errEmptyGraph
	}
	var (
		wT    float64
		keys  []K
		edges = []Edge[K, W]{}
	)
	// record selected vertexes.
	prev := make(map[K]K)
	// record current costs of uncoloured vertexes.
	cost := make(map[K]float64)
	for _, v := range vertexes {
		cost[v.Key] = MaxFloatDistance
	}
	// randomly select a vertex as mst root.
	prev[vertexes[0].Key] = vertexes[0].Key
	cost[vertexes[0].Key] = 0.0

	// loop until no unselected vertexes.
	for len(cost) != 0 {
		// find a minimum cost u.
		var u K
		cU := MaxFloatDistance
		for k, cK := range cost { // TODO: use a priority queue
			if cK < cU {
				cU = cK
				u = k
			}
		}
		if cU == MaxFloatDistance {
			return nil, nil, 0.0, errNotConnected
		}
		// join the u to tree.
		delete(cost, u)
		keys = append(keys, u)
		// update cost for all unselected vertexes,because joined a new vertex u to the mst,
		// maybe the costs of some unselected vertexes to current mst can be reduce by u.
		for v, c := range cost {
			_, w, err := getMinWeightEdge(g, u, v)
			if err != nil {
				return nil, nil, 0.0, err
			}
			if w < c {
				prev[v] = u
				cost[v] = w
			}
		}
		// update weight sum.
		wT += cU
		v := prev[u]
		if v != u {
			e, _, err := getMinWeightEdge(g, u, v)
			if err != nil {
				return nil, nil, 0.0, err
			}
			if e != nil {
				edges = append(edges, *e)
			}
		}
	}

	return keys, edges, wT, nil
}

// use priority queue.
func mstPrimWithPQ[K comparable, W number](g Graph[K, any, W]) ([]K, []Edge[K, W], float64, error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, nil, 0.0, err
	}
	if len(vertexes) == 0 {
		return nil, nil, 0.0, errEmptyGraph
	}
	var (
		wT    float64
		keys  []K
		edges = []Edge[K, W]{}
	)
	// record selected vertexes.
	prev := make(map[K]K)
	unvisited := make(map[K]bool)
	// record current costs of uncoloured vertexes.
	cost := newCostQueue[K]()
	for _, v := range vertexes {
		c := &item[K]{
			key:   v.Key,
			value: MaxFloatDistance,
		}
		cost.Push(c)
		unvisited[v.Key] = true
	}
	cost.Update(vertexes[0].Key, 0.0)

	for cost.Len() != 0 {
		// find a minimum cost u.
		u := cost.Pop()
		if u.value == MaxFloatDistance {
			return nil, nil, 0.0, errNotConnected
		}
		// join the u to tree.
		keys = append(keys, u.key)
		delete(unvisited, u.key)

		for v := range unvisited {
			_, w, err := getMinWeightEdge(g, u.key, v)
			if err != nil {
				return nil, nil, 0.0, err
			}
			if w < cost.Get(v) {
				prev[v] = u.key
				cost.Update(v, w)
			}
		}
		// update weight sum.
		wT += u.value
		v := prev[u.key]
		if v != u.key {
			e, _, err := getMinWeightEdge(g, u.key, v)
			if err != nil {
				return nil, nil, 0.0, err
			}
			if e != nil {
				edges = append(edges, *e)
			}
		}
	}

	return keys, edges, wT, nil
}

func msfPrim[K comparable, W number](g Graph[K, any, W]) ([][]K, [][]Edge[K, W], []float64, error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, nil, nil, err
	}
	if len(vertexes) == 0 {
		return nil, nil, nil, errEmptyGraph
	}
	var (
		wTs   []float64
		trees [][]K
		edges [][]Edge[K, W]
	)
	// record selected vertexes.
	prev := make(map[K]K)
	// record current costs of uncoloured vertexes.
	cost := make(map[K]float64)
	for _, v := range vertexes {
		cost[v.Key] = MaxFloatDistance
	}
	// randomly select a vertex as mst root.
	root := vertexes[0]
	prev[root.Key] = root.Key
	cost[root.Key] = 0.0

	var (
		wT     float64
		tree   []K
		branch []Edge[K, W]
	)
	// loop until no unselected vertexes.
	for len(cost) != 0 {
		// find a minimum cost u.
		var u K
		cU := MaxFloatDistance
		for k, cK := range cost { // TODO: use a priority queue
			if cK < cU {
				cU = cK
				u = k
			}
		}
		if cU == MaxFloatDistance {
			if len(tree) != 0 {
				trees = append(trees, tree)
				edges = append(edges, branch)
				wTs = append(wTs, wT)
			}
			tree = []K{}
			branch = []Edge[K, W]{}
			wT = 0.0

			// find a new root
			for k := range cost {
				u = k
				break
			}
			prev[u] = u
			cU = 0.0
		}
		// join the u to tree.
		delete(cost, u)
		tree = append(tree, u)
		// update cost for all unselected vertexes,because joined a new vertex u to the mst,
		// maybe the costs of some unselected vertexes to current mst can be reduce by u.
		for v, c := range cost {
			_, w, err := getMinWeightEdge(g, u, v)
			if err != nil {
				return nil, nil, nil, err
			}
			if w < c {
				prev[v] = u
				cost[v] = w
			}
		}
		wT += cU // update weight sum.
		v := prev[u]
		if v != u {
			e, _, err := getMinWeightEdge(g, v, u)
			if err != nil {
				return nil, nil, nil, err
			}
			if e != nil {
				branch = append(branch, *e)
			}
		}
	}

	return trees, edges, wTs, nil
}

// It uses a disjoint-set data structure to maintain several
// disjoint sets of elements. Each set contains the vertices in one tree
// of the current forest.
//
// MST-KRUSKAL(G,w)
// 1 A = null;
// 2 for each vertex v ∈ G.V
// 3     MAKE SET(v)
// 4 sort the edges of G.E into nondecreasing order by weight w
// 5 for each edge (u,v) ∈ G.E, taken in nondecreasing order by weight
// 6     if FIND-SET(u) != FIND-SET(v)
// 7         A = A Union {(u,v)}
// 8         UNION(u,v)
// 9 return A
func mstKruskal[K comparable, W number](g Graph[K, any, W]) ([]K, []Edge[K, W], float64, error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, nil, 0.0, err
	}
	if len(vertexes) == 0 {
		return nil, nil, 0.0, errEmptyGraph
	}
	// vertex disjoin set.
	prev := make(map[K]K)
	for _, v := range vertexes {
		prev[v.Key] = v.Key
	}

	// if v1 and v2 has common ancestor,they in same subtree,
	// which means join edge v1-v2 to it will produce a cycle.
	inSameSet := func(v1, v2 K) bool {
		vs := make(map[K]struct{})
		for p := prev[v1]; p != v1; p = prev[p] {
			vs[p] = struct{}{}
		}
		for p := prev[v2]; p != v2; p = prev[p] {
			if _, ok := vs[p]; ok {
				return true
			}
		}
		return false
	}

	// Sort all the edges in non-decreasing order of their weight.
	allEdges, err := g.AllEdges()
	if err != nil {
		return nil, nil, 0.0, err
	}
	sort.Slice(allEdges, func(i, j int) bool {
		return allEdges[i].Weight < allEdges[j].Weight
	})

	var (
		idx   int
		wT    float64
		edges = []Edge[K, W]{}
	)
	// Repeat until there are (V-1) edges in the spanning tree.
	for len(edges) != len(vertexes)-1 {
		if idx >= len(allEdges) {
			return nil, nil, 0.0, errNotConnected
		}
		// Pick the smallest edge.
		// Check if it forms a cycle with the spanning tree formed so far.
		// If cycle is not formed, include this edge. Else, discard it.
		e := allEdges[idx]
		if e.Tail != e.Head {
			// check if we can
			if !inSameSet(e.Tail, e.Head) {
				// join the edge to mst.
				prev[e.Tail] = e.Head
				edges = append(edges, e)
				// update weight sum.
				wT += any(e.Weight).(float64)
			}
		}
		idx++
	}

	keys := make([]K, len(vertexes))
	for i, v := range vertexes {
		keys[i] = v.Key
	}

	return keys, edges, wT, nil
}

// Generate the minimum spanning tree of a weighted connected graph,
// return the set of edges of the tree, and the sum of tree weights.
// If the graph is non connected, an error will be returned.
func MinWeightSpanningTree[K comparable, W number](g Graph[K, any, W]) ([]Edge[K, W], float64, error) {
	_, es, w, err := mstPrim(g)
	return es, w, err
}

// Generate the minimum spanning forest of a weighted graph,
// which generates a corresponding minimum spanning tree for each connected component,
// returns the set of vertices and edges for each tree, as well as the sum of tree weights.
func MinWeightSpanningForest[K comparable, W number](g Graph[K, any, W]) ([][]K, [][]Edge[K, W], []float64, error) {
	return msfPrim(g)
}
