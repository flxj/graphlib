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

type ThinGraph[K comparable] struct {
	Graph[K, any, int]
}

func NewThinGraph[K comparable](digraph bool) *ThinGraph[K] {
	g, _ := newGraph[K, any, int](digraph, "")
	return &ThinGraph[K]{
		g,
	}
}

type ThinDigraph[K comparable] struct {
	Digraph[K, any, int]
}

func NewThinDigraph[K comparable]() *ThinDigraph[K] {
	g, _ := NewDigraph[K, any, int]("")
	return &ThinDigraph[K]{
		g,
	}
}

// TODO
type ThinTree[K comparable] struct {
	Graph[K, any, int]
}

func NewThinTree[K comparable]() *ThinTree[K] {
	g, _ := newGraph[K, any, int](false, "")
	return &ThinTree[K]{g}
}

func (t *ThinTree[K]) AddEdge(e Edge[K, int]) error {
	// TODO if construct cycle, retrurn err
	return errNotImplement
}

// Tarjan
func (t *ThinTree[K]) LeastCommonAncestors(k1, k2 K) (k K, b bool) {
	/*
		v1, err := t.Graph.GetVertex(k1)
		if err != nil {
			return
		}
		v2, err := t.Graph.GetVertex(k2)
		if err != nil {
			return
		}
	*/
	return
}
