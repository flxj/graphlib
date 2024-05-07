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

import (
	"math"
)

// Ford-Fulkerson Algorithm and Edmonds-Karp Algorithm:
//
// 1.Start with initial flow as 0.
// 2.While there exists an augmenting path from the source to the sink:
//
//	2.1) Find an augmenting path using any path-finding algorithm,
//	    such as breadth-first search(EK) or depth-first search(FF).
//	2.2) Determine the amount of flow that can be sent along the augmenting path,
//	    which is the minimum residual capacity along the edges of the path.
//	2.3) Increase the flow along the augmenting path by the determined amount.
//
// 3.Return the maximum flow.
func mfEdmondsKarp[K comparable, V any, W number](g Graph[K, V, W], source, sink K) (W, error) {
	var (
		flow     W
		from, to int
	)
	wm, err := NewWeightMatrix(g)
	if err != nil {
		return flow, err
	}
	vertexes := wm.Columns()
	index := make(map[K]int)
	for i, v := range vertexes {
		index[v] = i
	}
	from = index[source]
	to = index[sink]

	// Create a residual graph.
	// Residual graph where rg[i][j] indicates residual capacity of edge
	// from vertexes[i] to vertexes[j] (if there is an edge.
	// If rg[i][j] is 0, then there is not)
	rg := wm.Weight(flow)
	//
	prev := make(map[int]int)
	// find a path in current residual graph.
	augmentingPath := func(s, t int) (bool, error) {
		var find bool
		prev[s] = -1
		visited := make(map[int]bool)
		visited[s] = true

		err := BFS(g, vertexes[s], func(v Vertex[K, V]) error {
			u := index[v.Key]
			for p := 0; p < len(vertexes); p++ {
				if !visited[p] && rg[u][p] > 0 {
					prev[p] = u
					if p == to {
						find = true
						return errNone
					}
					visited[p] = true
				}
			}
			return nil
		})

		if err != errNone {
			return false, err
		}
		return find, nil
	}

	for {
		ok, err := augmentingPath(from, to)
		if err != nil {
			return flow, err
		}
		if !ok {
			break
		}
		f := getMaxValue(flow)
		for v := to; v != from && v >= 0; v = prev[v] {
			p := prev[v]
			if rg[p][v] < f {
				f = rg[p][v]
			}
		}
		// update residual capacities
		for v := to; v != from && v >= 0; v = prev[v] {
			p := prev[v]
			rg[p][v] -= f
			rg[v][p] += f
		}
		flow += f
	}

	return flow, nil
}

func send[K comparable, W number](vertexes []K, flows, capacity [][]W, level map[int]int, s, t int, f W, start []int) W {
	if s == t {
		return f
	}
	for p := start[s]; p < len(vertexes) && capacity[s][p] > 0; p++ {
		//
		if level[p] == level[s]+1 && flows[s][p] < capacity[s][p] {
			sendFlow := min(f, capacity[s][p]-flows[s][p])
			flow := send(vertexes, flows, capacity, level, p, t, sendFlow, start)
			if flow > 0 {
				flows[s][p] += flow
				flows[p][s] -= flow
				return flow
			}
		}
		start[s] = p
	}
	return 0
}

// Dinic’s algorithm :
//
// 1. Initialize residual graph G as given graph.
//
// 2. Do BFS of G to construct a level graph (or assign levels to vertices) and also check if more flow is possible.
//
// If more flow is not possible, then return
// Send multiple flows in G using level graph until blocking flow is reached.
// Here using level graph means, in every flow, levels of path nodes should be 0, 1, 2…(in order) from s to t.
//
// A flow is Blocking Flow if no more flow can be sent using level graph,
// i.e., no more s-t path exists such that path vertices have current levels 0, 1, 2… in order.
//
// In Dinic’s algorithm, we use BFS to check if more flow is possible and to construct level graph.
// In level graph, we assign levels to all nodes, level of a node is shortest distance
// (in terms of number of edges) of the node from source.
// Once level graph is constructed, we send multiple flows using this level graph.
func mfDinic[K comparable, V any, W number](g Graph[K, V, W], source, sink K) (W, error) {
	var (
		flow     W
		from, to int
	)
	wm, err := NewWeightMatrix(g)
	if err != nil {
		return flow, err
	}
	vertexes := wm.Columns()
	index := make(map[K]int)
	for i, v := range vertexes {
		index[v] = i
	}
	from = index[source]
	to = index[sink]

	//
	capacity := wm.Weight(flow)
	//
	flows := make([][]W, len(capacity))
	for i := 0; i < len(flows); i++ {
		flows[i] = make([]W, len(capacity))
	}
	//
	level := make(map[int]int)
	buildLevel := func(s, t int) (bool, error) {
		for i := 0; i < len(vertexes); i++ {
			level[i] = -1
		}
		level[s] = 0

		err := BFS(g, vertexes[s], func(v Vertex[K, V]) error {
			u := index[v.Key]
			for p := 0; p < len(vertexes); p++ {
				if level[p] < 0 && flows[u][p] < capacity[u][p] {
					level[p] = level[u] + 1
				}
			}
			return nil
		})
		return level[t] >= 0, err
	}

	for {
		ok, err := buildLevel(from, to)
		if err != nil {
			return flow, err
		}
		if !ok {
			break
		}
		for {
			start := make([]int, len(vertexes)+1)
			f := send(vertexes, flows, capacity, level, from, to, getMaxValue(flow), start)
			if f == 0 {
				break
			}
			flow += f
		}
	}

	return flow, nil
}

// Highest Label Preflow Push
func mfHLPP[K comparable, V any, W number](g Graph[K, V, W], source, sink K) (W, error) {
	return 0, errNotImplement
}

// Calculate the maximum flow from the source vertex to the sink vertex.
func MaxFlow[K comparable, V any, W number](g Graph[K, V, W], source, sink K) (W, error) {
	var (
		flow W
		err  error
	)
	if _, err = g.GetVertex(source); err != nil {
		return flow, err
	}
	if _, err = g.GetVertex(sink); err != nil {
		return flow, err
	}
	return mfDinic(g, source, sink)
}

func mmBlossom[K comparable, V any, W number](g Graph[K, V, W]) ([]K, error) {
	return nil, errNotImplement
}

// Calculate the maximum matching of any graph and return the set of edges.
func MaxMatching[K comparable, V any, W number](g Graph[K, V, W]) ([]K, error) {
	return mmBlossom(g)
}

func MaxWeightMatching[K comparable, V any, W number](g Graph[K, V, W]) ([]K, error) {
	return nil, errNotImplement
}

// Calculate the perfect matching of any graph, if it exists, return the set of edges,
// otherwise return non-existent.
func PerfectMatching[K comparable, V any, W number](g Graph[K, V, W]) ([]K, error) {
	mm, err := MaxMatching(g)
	if err != nil {
		return nil, err
	}
	vs := make(map[K]bool)
	for _, k := range mm {
		e, err := g.GetEdgeByKey(k)
		if err != nil {
			return nil, err
		}
		vs[e.Head] = true
		vs[e.Tail] = true
	}
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}
	for _, v := range vertexes {
		if _, ok := vs[v.Key]; !ok {
			return nil, errMatchNotExists
		}
	}
	return mm, nil
}

var inf = math.MaxInt

// update current matching by dfs:
//
//	from vertex k(in partA and not in current matching)
//	find the successor vertex of k in augmenting path.
func updateMatching[K comparable, V any, W number](g Bipartite[K, V, W], pairU, pairV map[K]K, dist map[K]int, u, dummyK K) (bool, error) {
	if u != dummyK {
		// get u's neighbours(in partB).
		// and check the already matched neighbour
		vs, err := g.Neighbours(u)
		if err != nil {
			return false, err
		}
		// edge u-v not in current matching.
		for _, v := range vs {
			// v already matched, we check its pair vertex(in partA) if is a successor of u.
			if dist[pairV[v.Key]] == dist[u]+1 {
				// pairV[v.Key] is a successor of u.
				// we should continue update matching from it.
				ok, err := updateMatching(g, pairU, pairV, dist, pairV[v.Key], dummyK)
				if err != nil {
					return false, err
				}
				if ok {
					// add the edge u-v in current matching.
					pairV[v.Key] = u
					pairU[u] = v.Key
					return true, nil
				}
			}
		}
		// not exists a augmenting path from u.
		dist[u] = inf
		return false, nil
	}
	return true, nil
}

// Hopcroft Karp Algorithm:
//
// 1.Initialize Maximal Matching M as empty.
//
// 2.While there exists an Augmenting Path p
//      Remove matching edges of p from M and add not-matching edges of p to M
//      (This increases size of M by 1 as p starts and ends with a free vertex)
//
// 3.Return M.
//
// The idea is to use BFS (Breadth First Search) to find augmenting paths.
// Since BFS traverses level by level, it is used to divide the graph in layers of matching and not matching edges.
// A dummy vertex dummyK is added that is connected to all vertices on the left side and all vertices on the right side.
// The following maps are used to find augmenting paths. Distance to dummyK is initialized as INF (infinite).
// If we start from a dummy vertex and come back to it using alternating paths of distinct vertices, then there is an augmenting path.
//
// pairU: An map of size m where m is the number of vertices on the left side of the Bipartite Graph.
// pairU[u] stores pair of u on the right side if u is matched and dummyK otherwise.
//
// pairV: An amp of size n where n is several vertices on the right side of the Bipartite Graph.
// pairV[v] stores a pair of v on the left side if v is matched and dummyK otherwise.
//
// dist: An map of size m where m is several vertices on the left side of the Bipartite Graph.
// dist[u] is initialized as 0 if u is not matching and INF (infinite) otherwise. dist[dummyK] is also initialized as INF

// Once an augmenting path is found, DFS (Depth First Search) is used to add augmenting paths to current matching.
// DFS simply follows the distance map setup by BFS. It fills values in pairU[u] and pairV[v] if v is next to u in BFS.
func mmHopcroftKarp[K comparable, V any, W number](g Bipartite[K, V, W]) ([]K, error) {
	var (
		dummyK K
		err    error
		partA  []Vertex[K, V]
		partB  []Vertex[K, V]
	)
	if partA, err = g.Part(true); err != nil {
		return nil, err
	}
	if partB, err = g.Part(false); err != nil {
		return nil, err
	}
	//
	pairU := make(map[K]K)
	pairV := make(map[K]K)
	dist := make(map[K]int)
	for _, u := range partA {
		pairU[u.Key] = dummyK
	}
	for _, v := range partB {
		pairV[v.Key] = dummyK
	}

	// find augmenting path by bfs:
	//   search alternating path:
	//   --> start from dummyK
	//   --> vertexes of partA that matched with dummyK
	//   --> vertexes of partB that unmatched with prev vertexes
	//   --> vertexes of partA that matched with prev vertexes
	//   --> ...
	augmentingPath := func() (bool, error) {
		// queue store search vertexes in partA.
		que := newFIFO[K]()
		for _, u := range partA {
			// if u matched with dummyK,
			// which means it not in current matching.
			if pairU[u.Key] == dummyK {
				dist[u.Key] = 0
				que.push(u.Key)
			} else {
				// u has already in current matching.
				dist[u.Key] = inf
			}
		}
		// put dummyK in current matching
		dist[dummyK] = inf
		//
		for !que.empty() {
			u, _ := que.pop()
			// which means u not in current matching
			if dist[u] < dist[dummyK] {
				// so we should try to visit its neighbours(in partB) v,
				// obvious edge u-v is not in current matching.
				vs, err := g.Neighbours(u)
				if err != nil {
					return false, err
				}
				// travel vs, and find some v that already matched with some vertex(in partA),
				// which means we find some edges that in current matching.
				for _, v := range vs {
					// v has already matched.
					if dist[pairV[v.Key]] == inf {
						// set vertex dist, we will use the dist to ????
						dist[pairV[v.Key]] = dist[u] + 1
						// put the matched vertex in queue.
						// we will continue to find alternating path from it.
						que.push(pairV[v.Key])
					}
				}
			}
		}
		// after searching, if dist[dummyK] != inf, means we visited dummyK(some pairV[v.Key]) in bfs loop,
		// which means we find a alternating path start at dummyK and end at dummyK.
		// if dist[dummyK] == inf,means search path break off at some vertexes,cannot reach dummyK again.
		return dist[dummyK] != inf, nil
	}
	//
	for {
		ok, err := augmentingPath()
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
		for _, u := range partA {
			// u not in current matching.
			if pairU[u.Key] == dummyK {
				_, err = updateMatching(g, pairU, pairV, dist, u.Key, dummyK)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	//
	var edges []K
	for u, v := range pairU {
		if v != dummyK {
			es, err := g.GetEdge(u, v)
			if err != nil {
				return nil, err
			}
			if len(es) > 0 {
				edges = append(edges, es[0].Key)
			}
		}
	}
	return edges, nil
}

// Calculate the maximum matching of bipartite graph.
func BipartiteMaxMatching[K comparable, V any, W number](g Bipartite[K, V, W]) ([]K, error) {
	return mmHopcroftKarp(g)
}

func BipartitePerfectMatching[K comparable, V any, W number](g Bipartite[K, V, W]) ([]K, error) {
	mm, err := BipartiteMaxMatching(g)
	if err != nil {
		return nil, err
	}
	vs := make(map[K]bool)
	for _, k := range mm {
		e, err := g.GetEdgeByKey(k)
		if err != nil {
			return nil, err
		}
		vs[e.Head] = true
		vs[e.Tail] = true
	}
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}
	for _, v := range vertexes {
		if _, ok := vs[v.Key]; !ok {
			return nil, errMatchNotExists
		}
	}
	return mm, nil
}
