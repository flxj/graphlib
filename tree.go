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
func mstPrim[K comparable, V any, W number](g Graph[K, V, W]) ([]K, []Edge[K, W], W, error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, nil, 0.0, err
	}
	if len(vertexes) == 0 {
		return nil, nil, 0.0, errEmptyGraph
	}
	var (
		wT    W
		keys  []K
		edges = []Edge[K, W]{}
		maxW  = getMaxValue(wT)
	)
	// record selected vertexes.
	prev := make(map[K]K)
	// record current costs of uncoloured vertexes.
	cost := make(map[K]W)
	for _, v := range vertexes {
		cost[v.Key] = maxW
	}
	// randomly select a vertex as mst root.
	prev[vertexes[0].Key] = vertexes[0].Key
	cost[vertexes[0].Key] = 0.0

	// loop until no unselected vertexes.
	for len(cost) != 0 {
		// find a minimum cost u.
		var u K
		cU := maxW
		for k, cK := range cost { // TODO: use a priority queue
			if cK < cU {
				cU = cK
				u = k
			}
		}
		if cU == maxW {
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
func mstPrimWithPQ[K comparable, V any, W number](g Graph[K, V, W]) ([]K, []Edge[K, W], W, error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, nil, 0.0, err
	}
	if len(vertexes) == 0 {
		return nil, nil, 0.0, errEmptyGraph
	}
	var (
		wT    W
		keys  []K
		edges = []Edge[K, W]{}
		maxW  = getMaxValue(wT)
	)
	// record selected vertexes.
	prev := make(map[K]K)
	unvisited := make(map[K]bool)
	// record current costs of uncoloured vertexes.
	cost := newPriorityQueue[K, int, W](func(p1, p2 W) bool { return p1 < p2 })
	for _, v := range vertexes {
		cost.Push(v.Key, 0, maxW)
		unvisited[v.Key] = true
	}
	cost.Update(vertexes[0].Key, 0.0)

	for cost.Len() != 0 {
		// find a minimum cost u.
		u, _, c, _ := cost.Pop()
		if c == maxW {
			return nil, nil, 0.0, errNotConnected
		}
		// join the u to tree.
		keys = append(keys, u)
		delete(unvisited, u)

		for v := range unvisited {
			_, w, err := getMinWeightEdge(g, u, v)
			if err != nil {
				return nil, nil, 0.0, err
			}
			if w < cost.Get(v) {
				prev[v] = u
				cost.Update(v, w)
			}
		}
		// update weight sum.
		wT += c
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

func msfPrim[K comparable, V any, W number](g Graph[K, V, W]) ([][]K, [][]Edge[K, W], []W, error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, nil, nil, err
	}
	if len(vertexes) == 0 {
		return nil, nil, nil, errEmptyGraph
	}
	var (
		n     W
		wTs   []W
		trees [][]K
		edges [][]Edge[K, W]
		maxW  = getMaxValue(n)
	)
	// record selected vertexes.
	prev := make(map[K]K)
	// record current costs of uncoloured vertexes.
	cost := make(map[K]W)
	for _, v := range vertexes {
		cost[v.Key] = maxW
	}
	// randomly select a vertex as mst root.
	root := vertexes[0]
	prev[root.Key] = root.Key
	cost[root.Key] = 0.0

	var (
		wT     W
		tree   []K
		branch []Edge[K, W]
	)
	// loop until no unselected vertexes.
	for len(cost) != 0 {
		// find a minimum cost u.
		var u K
		cU := maxW
		for k, cK := range cost { // TODO: use a priority queue
			if cK < cU {
				cU = cK
				u = k
			}
		}
		if cU == maxW {
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
func mstKruskal[K comparable, V any, W number](g Graph[K, V, W]) ([]K, []Edge[K, W], W, error) {
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
		wT    W
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
				wT += e.Weight
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
func MinWeightSpanningTree[K comparable, V any, W number](g Graph[K, V, W]) ([]Edge[K, W], W, error) {
	_, es, w, err := mstPrimWithPQ(g)
	return es, w, err
}

// Generate the minimum spanning forest of a weighted graph,
// which generates a corresponding minimum spanning tree for each connected component,
// returns the set of vertices and edges for each tree, as well as the sum of tree weights.
func MinWeightSpanningForest[K comparable, V any, W number](g Graph[K, V, W]) ([][]K, [][]Edge[K, W], []W, error) {
	return msfPrim(g)
}
