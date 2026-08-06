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
	var w W
	if g == nil {
		return w, errNilGraph
	}
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
