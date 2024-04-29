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

// Does g1 include g2 as a subgraph.
func Contains[K comparable, W number](g1 Graph[K, any, W], g2 Graph[K, any, W]) (bool, error) {
	if g1.IsDigraph() != g2.IsDigraph() {
		return false, errNotSameType
	}
	var (
		err error
		vs1 []Vertex[K, any]
		es1 []Edge[K, W]
	)
	if vs1, err = g2.AllVertexes(); err != nil {
		return false, err
	}
	if es1, err = g2.AllEdges(); err != nil {
		return false, err
	}

	for _, v := range vs1 {
		_, err = g1.GetVertex(v.Key)
		if err != nil {
			if !IsNotExists(err) {
				return false, err
			}
			return false, nil
		}
	}

	for _, e := range es1 {
		es, err := g1.GetEdge(e.Head, e.Tail)
		if err != nil {
			if !IsNotExists(err) {
				return false, err
			}
			return false, nil
		}
		//
		if len(es) == 0 {
			return false, nil
		}
	}

	return true, nil
}

// Generate a spanning subgraph of g, and the new graph will not include edges in the edges list.
// The format of the edges list is [] [] K {head1, tail1} {headN, tailN}}.
func SpanningSubgraph[K comparable, W number](g Graph[K, any, W], edges [][]K) (Graph[K, any, W], error) {
	ng, err := g.Clone()
	if err != nil {
		return nil, err
	}
	for _, es := range edges {
		if len(es) >= 2 {
			if err := ng.RemoveEdge(es[0], es[1]); err != nil {
				return nil, err
			}
		}
	}
	return ng, nil
}

// Generate a spanning supergraph of g, and add edges to the edges list in the new graph.
func SpanningSupergraph[K comparable, W number](g Graph[K, any, W], edges []*Edge[K, W]) (Graph[K, any, W], error) {
	ng, err := g.Clone()
	if err != nil {
		return nil, err
	}
	for _, e := range edges {
		ee := *e
		if err := ng.AddEdge(ee); err != nil {
			return nil, err
		}
	}
	return ng, nil
}

// Generate an induced subgraph of g, where the new graph will not contain vertices in vertices.
func InducedSubgraph[K comparable, W number](g Graph[K, any, W], vertexes []K) (Graph[K, any, W], error) {
	ng, err := g.Clone()
	if err != nil {
		return nil, err
	}
	for _, v := range vertexes {
		if err := ng.RemoveVertex(v); err != nil {
			return nil, err
		}
	}
	return ng, nil
}
