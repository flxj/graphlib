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

import "fmt"

// A line graph L(G) (also called an adjoint/edge-to-vertex dual graph) of a simple graph G
// is obtained by associating a vertex with each edge of the graph and connecting two vertices
// with an edge iff the corresponding edges of G have a vertex in common.
func LineGraph[K comparable, W number](g Graph[K, W]) (Graph[int, int], error) {
	if g == nil {
		return nil, errNilGraph
	}
	edges := g.AllEdges()
	idx := make(map[K]int)
	lg := NewGraph[int, int](false, g.Name()+"_line")
	for i, e := range edges {
		idx[e.Key] = i
		err := lg.AddVertex(Vertex[int, int]{
			Key:   i,
			Value: e.Labels,
			Labels: map[string]string{
				"edge":   fmt.Sprintf("%v", e.Key),
				"head":   fmt.Sprintf("%v", e.Head),
				"tail":   fmt.Sprintf("%v", e.Tail),
				"weight": fmt.Sprintf("%v", e.Weight),
			}})
		if err != nil {
			return nil, err
		}
	}
	var ek int
	addEdges := func(i int, e, v K) error {
		es, err := g.IncidentEdges(v)
		if err != nil {
			return err
		}
		for _, e1 := range es {
			j := idx[e1.Key]
			if j > i {
				if err := lg.AddEdge(Edge[int, int]{
					Key:  ek,
					Head: i,
					Tail: j,
					Labels: map[string]string{
						"edge1": fmt.Sprintf("%v", e),
						"edge2": fmt.Sprintf("%v", e1.Key),
					},
				}); err != nil {
					return err
				}
				ek++
			}
		}
		return nil
	}
	for i, e := range edges {
		if err := addEdges(i, e.Key, e.Head); err != nil {
			return nil, err
		}
		if err := addEdges(i, e.Key, e.Tail); err != nil {
			return nil, err
		}
	}
	return lg, nil
}

func LineDigraph[K comparable, W number](g Digraph[K, W]) (Digraph[int, int], error) {
	if g == nil {
		return nil, errNilGraph
	}
	edges := g.AllEdges()
	idx := make(map[K]int)
	lg := NewDigraph[int, int](g.Name() + "_line")
	for i, e := range edges {
		idx[e.Key] = i
		err := lg.AddVertex(Vertex[int, int]{
			Key:   i,
			Value: e.Labels,
			Labels: map[string]string{
				"edge":   fmt.Sprintf("%v", e.Key),
				"head":   fmt.Sprintf("%v", e.Head),
				"tail":   fmt.Sprintf("%v", e.Tail),
				"weight": fmt.Sprintf("%v", e.Weight),
			}})
		if err != nil {
			return nil, err
		}
	}
	var ek int
	addEdges := func(i int, e, head K) error {
		es, err := g.OutEdges(head)
		if err != nil {
			return err
		}
		for _, e1 := range es {
			j := idx[e1.Key]
			err := lg.AddEdge(Edge[int, int]{
				Key:  ek,
				Tail: i,
				Head: j,
				Labels: map[string]string{
					"edge1": fmt.Sprintf("%v", e),
					"edge2": fmt.Sprintf("%v", e1.Key),
				},
			})
			if err != nil {
				return err
			}
			ek++
		}
		return nil
	}
	for i, e := range edges {
		if err := addEdges(i, e.Key, e.Head); err != nil {
			return nil, err
		}
	}
	return lg, nil
}
