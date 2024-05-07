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

import "math/rand"

type Bipartite[K comparable, V any, W number] struct {
	Graph[K, V, W]
	g     *graph[K, V, W]
	partA map[K]bool
	partB map[K]bool
}

func NewBipartite[K comparable, V any, W number](digraph bool, name string) (*Bipartite[K, V, W], error) {
	g, err := newGraph[K, V, W](digraph, name)
	if err != nil {
		return nil, err
	}
	return &Bipartite[K, V, W]{
		g:     g,
		partA: make(map[K]bool),
		partB: make(map[K]bool),
	}, nil
}

func (bg *Bipartite[K, V, W]) Name() string {
	return bg.g.Name()
}

func (bg *Bipartite[K, V, W]) SetName(name string) {
	bg.g.SetName(name)
}

func (bg *Bipartite[K, V, W]) IsDigraph() bool {
	return bg.g.IsDigraph()
}

func (bg *Bipartite[K, V, W]) Order() int {
	return bg.g.Order()
}

func (bg *Bipartite[K, V, W]) Size() int {
	return bg.g.Size()
}

func (bg *Bipartite[K, V, W]) Property(p PropertyName) (GraphProperty[any], error) {
	return bg.g.Property(p)
}

func (bg *Bipartite[K, V, W]) AllVertexes() ([]Vertex[K, V], error) {
	return bg.g.AllVertexes()
}

func (bg *Bipartite[K, V, W]) AllEdges() ([]Edge[K, W], error) {
	return bg.g.AllEdges()
}

func (bg *Bipartite[K, V, W]) AddVertex(v Vertex[K, V]) error {
	if err := bg.g.AddVertex(v); err != nil {
		return err
	}
	if rand.Intn(2) == 0 {
		bg.partA[v.Key] = true
	} else {
		bg.partB[v.Key] = true
	}
	return nil
}

func (bg *Bipartite[K, V, W]) AddVertexTo(v Vertex[K, V], partA bool) error {
	if err := bg.g.AddVertex(v); err != nil {
		return err
	}
	if partA {
		bg.partA[v.Key] = true
	} else {
		bg.partB[v.Key] = true
	}
	return nil
}

func (bg *Bipartite[K, V, W]) Part(partA bool) ([]Vertex[K, V], error) {
	var vs []Vertex[K, V]
	var ks map[K]bool
	if partA {
		ks = bg.partA
	} else {
		ks = bg.partB
	}
	for k := range ks {
		v, err := bg.g.GetVertex(k)
		if err != nil {
			return nil, err
		}
		vs = append(vs, v)
	}
	return vs, nil
}

func (bg *Bipartite[K, V, W]) RemoveVertex(key K) error {
	if err := bg.g.RemoveVertex(key); err != nil {
		return err
	}
	delete(bg.partA, key)
	delete(bg.partB, key)
	return nil
}

func (bg *Bipartite[K, V, W]) AddEdge(edge Edge[K, W]) error {
	if bg.partA[edge.Head] || bg.partA[edge.Tail] {
		return errViolateBipartite
	}
	if bg.partB[edge.Head] || bg.partB[edge.Tail] {
		return errViolateBipartite
	}

	return bg.g.AddEdge(edge)
}

func (bg *Bipartite[K, V, W]) RemoveEdgeByKey(key K) error {
	return bg.g.RemoveEdgeByKey(key)
}

func (bg *Bipartite[K, V, W]) RemoveEdge(v1, v2 K) error {
	return bg.g.RemoveEdge(v1, v2)
}

func (bg *Bipartite[K, V, W]) Degree(key K) (int, error) {
	return bg.g.Degree(key)
}

func (bg *Bipartite[K, V, W]) Neighbours(v K) ([]Vertex[K, V], error) {
	return bg.g.Neighbours(v)
}

func (bg *Bipartite[K, V, W]) GetVertex(key K) (Vertex[K, V], error) {
	return bg.g.GetVertex(key)
}

func (bg *Bipartite[K, V, W]) GetEdge(v1, v2 K) ([]Edge[K, W], error) {
	return bg.g.GetEdge(v1, v2)
}

func (bg *Bipartite[K, V, W]) GetEdgeByKey(key K) (Edge[K, W], error) {
	return bg.g.GetEdgeByKey(key)
}

func (bg *Bipartite[K, V, W]) GetVertexesByLabel(labels map[string]string) ([]Vertex[K, V], error) {
	return bg.g.GetVertexesByLabel(labels)
}

func (bg *Bipartite[K, V, W]) GetEdgesByLabel(labels map[string]string) ([]Edge[K, W], error) {
	return bg.g.GetEdgesByLabel(labels)
}

func (bg *Bipartite[K, V, W]) SetVertexValue(key K, value V) error {
	return bg.g.SetVertexValue(key, value)
}

func (bg *Bipartite[K, V, W]) SetVertexLabel(key K, labelKey, labelVal string) error {
	return bg.g.SetVertexLabel(key, labelKey, labelVal)
}

func (bg *Bipartite[K, V, W]) DeleteVertexLabel(key K, labelKey string) error {
	return bg.g.DeleteVertexLabel(key, labelKey)
}

func (bg *Bipartite[K, V, W]) SetEdgeValueByKey(key K, value any) error {
	return bg.g.SetEdgeValueByKey(key, value)
}

func (bg *Bipartite[K, V, W]) SetEdgeLabelByKey(key K, labelKey, labelVal string) error {
	return bg.g.SetEdgeLabelByKey(key, labelKey, labelVal)
}

func (bg *Bipartite[K, V, W]) DeleteEdgeLabelByKey(key K, labelKey string) error {
	return bg.g.DeleteEdgeLabelByKey(key, labelKey)
}

func (bg *Bipartite[K, V, W]) SetEdgeValue(endpoint1, endpoint2 K, value any) error {
	return bg.g.SetEdgeValue(endpoint1, endpoint2, value)
}

func (bg *Bipartite[K, V, W]) SetEdgeLabel(endpoint1, endpoint2 K, labelKey, labelVal string) error {
	return bg.g.SetEdgeLabel(endpoint1, endpoint2, labelKey, labelVal)
}

func (bg *Bipartite[K, V, W]) DeleteEdgeLabel(endpoint1, endpoint2 K, labelKey string) error {
	return bg.g.DeleteEdgeLabel(endpoint1, endpoint2, labelKey)
}

func (bg *Bipartite[K, V, W]) Clone() (Graph[K, V, W], error) {
	g, err := bg.g.Clone()
	if err != nil {
		return nil, err
	}
	ng, ok := g.(*graph[K, V, W])
	if !ok {
		return nil, errCloneFailed
	}
	b := &Bipartite[K, V, W]{
		g:     ng,
		partA: make(map[K]bool),
		partB: make(map[K]bool),
	}
	for k := range bg.partA {
		b.partA[k] = true
	}
	for k := range bg.partB {
		b.partB[k] = true
	}
	return b, nil
}

func (bg *Bipartite[K, V, W]) InDegree(vertex K) (int, error) {
	return bg.g.InDegree(vertex)
}

func (bg *Bipartite[K, V, W]) OutDegree(vertex K) (int, error) {
	return bg.g.OutDegree(vertex)
}

func (bg *Bipartite[K, V, W]) InNeighbours(vertex K) ([]Vertex[K, V], error) {
	return bg.g.InNeighbours(vertex)
}

func (bg *Bipartite[K, V, W]) OutNeighbours(vertex K) ([]Vertex[K, V], error) {
	return bg.g.OutNeighbours(vertex)
}

func (bg *Bipartite[K, V, W]) InEdges(vertex K) ([]Edge[K, W], error) {
	return bg.g.InEdges(vertex)
}

func (bg *Bipartite[K, V, W]) OutEdges(vertex K) ([]Edge[K, W], error) {
	return bg.g.OutEdges(vertex)
}

func (bg *Bipartite[K, V, W]) Sources() ([]Vertex[K, V], error) {
	return bg.g.Sources()
}

func (bg *Bipartite[K, V, W]) Sinks() ([]Vertex[K, V], error) {
	return bg.g.Sinks()
}

func (bg *Bipartite[K, V, W]) DetectCycle() ([][]K, error) {
	return nil, errNotImplement
}

func (bg *Bipartite[K, V, W]) Recerse() error {
	return bg.g.Reverse()
}

/*
Following is a simple algorithm to find out whether a given graph is Bipartite or not using Breadth First Search (BFS).
1. Assign RED color to the source vertex (putting into set U).
2. Color all the neighbors with BLUE color (putting into set V).
3. Color all neighbor’s neighbor with RED color (putting into set U).
4. This way, assign color to all vertices such that it satisfies all the
   constraints of m way coloring problem where m = 2.
5. While assigning colors, if we find a neighbor which is colored with same color as current vertex,
   then the graph cannot be colored with 2 vertices (or graph is not Bipartite).
*/

// Determine whether the given graph is a bipartite graph.
func IsBipartite[K comparable, V any, W number](g Graph[K, V, W]) (bool, error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return false, err
	}
	switch len(vertexes) {
	case 0, 1:
		return false, nil
	case 2:
		return true, nil
	default:
	}

	// red:0,blue:1
	color := make(map[K]int)

	que := newFIFO[K]()
	que.push(vertexes[0].Key)
	color[vertexes[0].Key] = 0
	//
	for !que.empty() {
		u, _ := que.pop()
		cu := color[u]

		vs, err := g.Neighbours(u)
		if err != nil {
			return false, err
		}
		for _, v := range vs {
			if v.Key == u { // loop
				return false, nil
			}
			cv, ok := color[v.Key]
			if !ok {
				color[v.Key] = (cu + 1) % 2
				que.push(v.Key)
			} else {
				if cu == cv {
					return false, nil
				}
			}
		}
		if que.empty() && len(color) != len(vertexes) {
			for _, v := range vertexes {
				if _, ok := color[v.Key]; !ok {
					color[v.Key] = 0
					que.push(v.Key)
					break
				}
			}
		}
	}
	return true, nil
}
