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

// A bipartite graph, also called a bigraph, is a set of graph vertices decomposed
// into two disjoint sets such that no two graph vertices within the same set are adjacent.
// The two vertex classes are said to form a bipartition.
type Bipartite[K comparable, W number] interface {
	Graph[K, W]
	AddVertexTo(v Vertex[K, W], partA bool) bool
	Part(partA bool) ([]Vertex[K, W], bool)
	InPartA(K) bool
	PartOrder(partA bool) int
	RemovePart(partA bool)
}

type bipartite[K comparable, W number] struct {
	Graph[K, W]
	g     *graph[K, W]
	partA map[K]bool
	partB map[K]bool
}

func NewBipartite[K comparable, W number](digraph bool, name string) Bipartite[K, W] {
	g := newGraph[K, W](digraph, name)
	return &bipartite[K, W]{
		g:     g,
		partA: make(map[K]bool),
		partB: make(map[K]bool),
	}
}

func (bg *bipartite[K, W]) Name() string {
	return bg.g.Name()
}

func (bg *bipartite[K, W]) SetName(name string) {
	bg.g.SetName(name)
}

func (bg *bipartite[K, W]) IsDigraph() bool {
	return bg.g.IsDigraph()
}

func (bg *bipartite[K, W]) Order() int {
	return bg.g.Order()
}

func (bg *bipartite[K, W]) PartOrder(partA bool) int {
	if partA {
		return len(bg.partA)
	}
	return len(bg.partB)
}

func (bg *bipartite[K, W]) Size() int {
	return bg.g.Size()
}

func (bg *bipartite[K, W]) Property(p PropertyName) (GraphProperty[any], bool) {
	return bg.g.Property(p)
}

func (bg *bipartite[K, W]) AllVertexes() []Vertex[K, W] {
	return bg.g.AllVertexes()
}

func (bg *bipartite[K, W]) AllEdges() []Edge[K, W] {
	return bg.g.AllEdges()
}

func (bg *bipartite[K, W]) AddVertex(v Vertex[K, W]) bool {
	if ok := bg.g.AddVertex(v); !ok {
		return false
	}
	if rand.Intn(2) == 0 {
		bg.partA[v.Key] = true
	} else {
		bg.partB[v.Key] = true
	}
	return true
}

func (bg *bipartite[K, W]) AddVertexTo(v Vertex[K, W], partA bool) bool {
	if ok := bg.g.AddVertex(v); !ok {
		return false
	}
	if partA {
		bg.partA[v.Key] = true
	} else {
		bg.partB[v.Key] = true
	}
	return true
}

func (bg *bipartite[K, W]) Part(partA bool) ([]Vertex[K, W], bool) {
	var vs []Vertex[K, W]
	var ks map[K]bool
	if partA {
		ks = bg.partA
	} else {
		ks = bg.partB
	}
	for k := range ks {
		v, ok := bg.g.GetVertex(k)
		if !ok {
			return nil, false
		}
		vs = append(vs, v)
	}
	return vs, true
}

func (bg *bipartite[K, W]) RemoveVertex(key K) (Vertex[K, W], bool) {
	v, ok := bg.g.RemoveVertex(key)
	if !ok {
		return Vertex[K, W]{}, false
	}
	delete(bg.partA, key)
	delete(bg.partB, key)
	return v, true
}

func (bg *bipartite[K, W]) RemovePart(partA bool) {
	if partA {
		for v := range bg.partA {
			_, _ = bg.g.RemoveVertex(v)
		}
		bg.partA = make(map[K]bool)
	} else {
		for v := range bg.partB {
			_, _ = bg.g.RemoveVertex(v)
		}
		bg.partB = make(map[K]bool)
	}
}

func (bg *bipartite[K, W]) AddEdge(edge Edge[K, W]) bool {
	if bg.partA[edge.Head] && bg.partA[edge.Tail] {
		return false
	}
	if bg.partB[edge.Head] && bg.partB[edge.Tail] {
		return false
	}
	return bg.g.AddEdge(edge)
}

func (bg *bipartite[K, W]) RemoveEdgeByKey(key K) (Edge[K, W], bool) {
	return bg.g.RemoveEdgeByKey(key)
}

func (bg *bipartite[K, W]) RemoveEdge(v1, v2 K) ([]Edge[K, W], bool) {
	return bg.g.RemoveEdge(v1, v2)
}

func (bg *bipartite[K, W]) Degree(key K) (int, bool) {
	return bg.g.Degree(key)
}

func (bg *bipartite[K, W]) Neighbours(v K) ([]Vertex[K, W], bool) {
	return bg.g.Neighbours(v)
}

func (bg *bipartite[K, W]) GetVertex(key K) (Vertex[K, W], bool) {
	return bg.g.GetVertex(key)
}

func (bg *bipartite[K, W]) GetEdge(v1, v2 K) ([]Edge[K, W], bool) {
	return bg.g.GetEdge(v1, v2)
}

func (bg *bipartite[K, W]) GetEdgeByKey(key K) (Edge[K, W], bool) {
	return bg.g.GetEdgeByKey(key)
}

func (bg *bipartite[K, W]) GetVertexesByLabel(labels Labels) []Vertex[K, W] {
	return bg.g.GetVertexesByLabel(labels)
}

func (bg *bipartite[K, W]) GetEdgesByLabel(labels Labels) []Edge[K, W] {
	return bg.g.GetEdgesByLabel(labels)
}

func (bg *bipartite[K, W]) SetVertexValue(key K, value any) bool {
	return bg.g.SetVertexValue(key, value)
}

func (bg *bipartite[K, W]) SetVertexLabel(key K, labelKey string, labelVal any) bool {
	return bg.g.SetVertexLabel(key, labelKey, labelVal)
}

func (bg *bipartite[K, W]) DeleteVertexLabel(key K, labelKey string) bool {
	return bg.g.DeleteVertexLabel(key, labelKey)
}

func (bg *bipartite[K, W]) SetEdgeValueByKey(key K, value any) bool {
	return bg.g.SetEdgeValueByKey(key, value)
}

func (bg *bipartite[K, W]) SetEdgeLabelByKey(key K, labelKey string, labelVal any) bool {
	return bg.g.SetEdgeLabelByKey(key, labelKey, labelVal)
}

func (bg *bipartite[K, W]) DeleteEdgeLabelByKey(key K, labelKey string) bool {
	return bg.g.DeleteEdgeLabelByKey(key, labelKey)
}

func (bg *bipartite[K, W]) SetEdgeValue(endpoint1, endpoint2 K, value any) bool {
	return bg.g.SetEdgeValue(endpoint1, endpoint2, value)
}

func (bg *bipartite[K, W]) SetEdgeLabel(endpoint1, endpoint2 K, labelKey string, labelVal any) bool {
	return bg.g.SetEdgeLabel(endpoint1, endpoint2, labelKey, labelVal)
}

func (bg *bipartite[K, W]) DeleteEdgeLabel(endpoint1, endpoint2 K, labelKey string) bool {
	return bg.g.DeleteEdgeLabel(endpoint1, endpoint2, labelKey)
}

func (bg *bipartite[K, W]) Clone() Graph[K, W] {
	g := bg.g.Clone()
	ng, ok := g.(*graph[K, W])
	if !ok {
		return nil
	}
	b := &bipartite[K, W]{
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
	return b
}

func (bg *bipartite[K, W]) InDegree(vertex K) (int, bool) {
	return bg.g.InDegree(vertex)
}

func (bg *bipartite[K, W]) OutDegree(vertex K) (int, bool) {
	return bg.g.OutDegree(vertex)
}

func (bg *bipartite[K, W]) InNeighbours(vertex K) ([]Vertex[K, W], bool) {
	return bg.g.InNeighbours(vertex)
}

func (bg *bipartite[K, W]) OutNeighbours(vertex K) ([]Vertex[K, W], bool) {
	return bg.g.OutNeighbours(vertex)
}

func (bg *bipartite[K, W]) InEdges(vertex K) ([]Edge[K, W], bool) {
	return bg.g.InEdges(vertex)
}

func (bg *bipartite[K, W]) OutEdges(vertex K) ([]Edge[K, W], bool) {
	return bg.g.OutEdges(vertex)
}

func (bg *bipartite[K, W]) Sources() ([]Vertex[K, W], bool) {
	return bg.g.Sources()
}

func (bg *bipartite[K, W]) Sinks() ([]Vertex[K, W], bool) {
	return bg.g.Sinks()
}

func (bg *bipartite[K, W]) DetectCycle() ([][]K, error) {
	return nil, errNotImplement
}

func (bg *bipartite[K, W]) Recerse() {
	bg.g.Reverse()
}

func (bg *bipartite[K, W]) RandomVertex() (Vertex[K, W], bool) {
	return bg.g.RandomVertex()
}

//
func (bg *bipartite[K, W]) RandomEdge() (Edge[K, W], bool) {
	return bg.g.RandomEdge()
}

//
func (bg *bipartite[K, W]) NeighbourEdgesByKey(edge K) ([]Edge[K, W], bool) {
	return bg.g.NeighbourEdgesByKey(edge)
}

//
func (bg *bipartite[K, W]) NeighbourEdges(endpoint1, endpoint2 K) ([]Edge[K, W], bool) {
	return bg.g.NeighbourEdges(endpoint1, endpoint2)
}

func (bg *bipartite[K, W]) IncidentEdges(vertex K) ([]Edge[K, W], bool) {
	return bg.g.IncidentEdges(vertex)
}

func (bg *bipartite[K, W]) InPartA(key K) bool {
	_, ok := bg.partA[key]
	return ok
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
func IsBipartite[K comparable, V any, W number](g Graph[K, W]) (bool, error) {
	if g == nil {
		return false, errNilGraph
	}
	vertexes := g.AllVertexes()
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

		vs, _ := g.Neighbours(u)
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
