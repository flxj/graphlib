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

// Determine whether the start and end vertices in graph g are connected.
// If it is a directed graph, determine if there is a directed path from start to end.
func Connected[K comparable, V any, W number](g Graph[K, V, W], start, end K) (bool, error) {
	if g == nil {
		return false, errNilGraph
	}
	var connected bool
	visitor := func(v Vertex[K, V]) error {
		if v.Key == end {
			connected = true
			return errNone
		}
		return nil
	}
	err := DFS(g, start, visitor)
	if !connected {
		return false, err
	}
	return true, nil
}

func auxiliaryGraphEDP[K comparable, V any, W number](g Graph[K, V, W]) (Graph[K, V, int], error) {
	if g == nil {
		return nil, errNilGraph
	}
	aux, _ := NewGraph[K, V, int](g.IsDigraph(), "")
	vs, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}
	es, err := g.AllEdges()
	if err != nil {
		return nil, err
	}
	for _, v := range vs {
		if err := aux.AddVertex(Vertex[K, V]{Key: v.Key}); err != nil {
			return nil, err
		}
	}
	for _, e := range es {
		if err := aux.AddEdge(Edge[K, int]{Key: e.Key, Head: e.Head, Tail: e.Tail, Weight: 1}); err != nil {
			return nil, err
		}
	}
	return aux, nil
}

// The edge disjoint paths in the auxiliary digraph correspond to the node disjoint paths in the original graph.
func auxiliaryGraphVDP[K comparable, V any, W number](g Graph[K, V, W], source, target K) (Graph[int, any, int], int, int, error) {
	if g == nil {
		return nil, 0, 0, errNilGraph
	}
	aux, err := NewGraph[int, any, int](g.IsDigraph(), "")
	if err != nil {
		return nil, 0, 0, err
	}
	vs, err := g.AllVertexes()
	if err != nil {
		return nil, 0, 0, err
	}
	es, err := g.AllEdges()
	if err != nil {
		return nil, 0, 0, err
	}
	idx := make(map[K]int)
	var ek, s, t int
	for i, v := range vs {
		if err := aux.AddVertex(Vertex[int, any]{Key: -(i + 1)}); err != nil {
			return nil, 0, 0, err
		}
		if err := aux.AddVertex(Vertex[int, any]{Key: i + 1}); err != nil {
			return nil, 0, 0, err
		}
		if v.Key == source {
			s = -(i + 1)
		}
		if v.Key == target {
			t = i + 1
		}
		idx[v.Key] = i + 1
		if err := aux.AddEdge(Edge[int, int]{Key: ek, Head: i + 1, Tail: -(i + 1), Weight: 1}); err != nil {
			return nil, 0, 0, err
		}
		ek++
	}
	// add edge
	for _, e := range es {
		i, j := idx[e.Head], idx[e.Tail]
		if err := aux.AddEdge(Edge[int, int]{Key: ek, Head: -i, Tail: j, Weight: 1}); err != nil {
			return nil, 0, 0, err
		}
		ek++
	}
	return aux, s, t, nil
}

// Find maximum number of edge disjoint paths between two vertices.
func EdgeDisjointPath[K comparable, V any, W number](g Graph[K, V, W], source, target K) (int, error) {
	aux, err := auxiliaryGraphEDP(g)
	if err != nil {
		return 0, err
	}
	return MaxFlow(aux, source, target)
}

// Find maximum number of vertex disjoint paths between two vertices.
func VertexDisjointPath[K comparable, V any, W number](g Graph[K, V, W], source, target K) (int, error) {
	aux, s, t, err := auxiliaryGraphVDP(g, source, target)
	if err != nil {
		return 0, err
	}
	return MaxFlow(aux, s, t)
}

// Find maximum number of edge disjoint paths between two vertices.
func DigraphEdgeDisjointPath[K comparable, V any, W number](g Digraph[K, V, W], source, target K) (int, error) {
	return EdgeDisjointPath(g, source, target)
}

// Find maximum number of vertex disjoint paths between two vertices.
func DigraphVertexDisjointPath[K comparable, V any, W number](g Graph[K, V, W], source, target K) (int, error) {
	return VertexDisjointPath(g, source, target)
}
