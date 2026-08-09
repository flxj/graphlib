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

import "container/heap"

// In our formal description of DFS, each vertex x of D gets two time-stamps:
// tvisit(x) once x is visited and texpl(x) once x is declared explored.
//
// DFS
// Input: A digraph D =(V,A).
// Output: pred(v), tvisit(v)andtexpl(v) for every v ∈ V.
//  1. For each v ∈ V set pred(v):=nil, tvisit(v) := 0 and texpl(v):=0.
//  2. Set time := 0.
//  3. For each vertex v ∈ V do: if tvisit(v) = 0 then perform DFS-PROC(v).
//
// DFS-PROC(v)
//  1. Set time := time+1,tvisit(v):=time.
//  2. For each u ∈ N+(v) do: if tvisit(u) = 0 then pred(u) :=v and perform
//     DFS-PROC(u).
//  3. Set time := time+1, texpl(v):=time.
func dfs[K comparable, V any, W number](g Graph[K, V, W], start K, visitor func(Vertex[K, V]) error, neighbours func(K) ([]Vertex[K, V], error)) error {
	startV, err := g.GetVertex(start)
	if err != nil {
		return err
	}
	visited := make(map[K]struct{})
	stack := newStack[*Vertex[K, V]]()
	stack.push(&startV)

	for !stack.empty() {
		v, _ := stack.pop()
		if _, ok := visited[v.Key]; !ok {
			if err := visitor(*v); err != nil {
				return err
			}
			visited[v.Key] = struct{}{}
		}
		vs, err := neighbours(v.Key)
		if err != nil {
			return err
		}
		for _, v := range vs {
			if _, ok := visited[v.Key]; !ok {
				var vv = v
				stack.push(&vv)
			}
		}
	}
	return nil
}

// Start depth first search from the specified source vertex, where g can be a directed or undirected graph.
func DFS[K comparable, V any, W number](g Graph[K, V, W], start K, visitor func(Vertex[K, V]) error) error {
	if g == nil {
		return errNilGraph
	}
	neighbours := g.Neighbours
	if g.IsDigraph() {
		dg, ok := g.(Digraph[K, V, W])
		if ok {
			neighbours = dg.OutNeighbours
		} else {
			neighbours = func(v K) ([]Vertex[K, V], error) {
				es, err := g.IncidentEdges(v)
				if err != nil {
					return nil, err
				}
				var res []Vertex[K, V]
				for _, e := range es {
					if e.Tail == v {
						w, err := g.GetVertex(e.Head)
						if err != nil {
							return nil, err
						}
						res = append(res, w)
					}
				}
				return res, nil
			}
		}
	}
	return dfs(g, start, visitor, neighbours)
}

// Perform depth first search in a directed graph, and specify the search direction using the in parameter:
// if in is set to true, search from source in the order of the incident vertices of the current vertex.
func DFSDigraph[K comparable, V any, W number](dg Digraph[K, V, W], start K, in bool, visitor func(Vertex[K, V]) error) error {
	if dg == nil {
		return errNilGraph
	}
	var neighbours func(K) ([]Vertex[K, V], error)
	if in {
		neighbours = dg.InNeighbours
	} else {
		neighbours = dg.OutNeighbours
	}
	return dfs(dg, start, visitor, neighbours)
}

//	 BFS
//	 Input: A digraph D =(V,A) and a vertex s ∈ V.
//	 Output: dist(s,v) and pred(v) for all v ∈ V.
//		1. For each v ∈ V set dist(s,v):=∞ and pred(v):=nil.
//		2. Set dist(s,s) := 0. Create a queue Q consisting of s.
//		3. While Q is not empty do the following. Delete a vertex u, the head of Q,
//		from Q and consider the out-neighbours of u in D one by one. If, for an
//		out-neighbour v of u,dist(s,v)=∞,thensetdist(s,v):=dist(s,u)+1,
//		pred(v):=u, and put v to the end of Q.
func bfs[K comparable, V any, W number](g Graph[K, V, W], start K, visitor func(Vertex[K, V]) error, neighbours func(K) ([]Vertex[K, V], error)) error {
	startV, err := g.GetVertex(start)
	if err != nil {
		return err
	}
	visited := make(map[K]struct{})
	// use a fifo queue.
	queue := newFIFO[*Vertex[K, V]]()
	queue.push(&startV)

	// visit current vertex,and push all neighbours of it to queue.
	for !queue.empty() {
		v, _ := queue.pop()
		if _, ok := visited[v.Key]; !ok {
			if err := visitor(*v); err != nil {
				return err
			}
			visited[v.Key] = struct{}{}
		}
		vs, err := neighbours(v.Key)
		if err != nil {
			return err
		}
		for _, v := range vs {
			if _, ok := visited[v.Key]; !ok {
				var vv = v
				queue.push(&vv)
			}
		}
	}
	return nil
}

// Start breadth first search from the specified source vertex, where g can be a directed or undirected graph.
func BFS[K comparable, V any, W number](g Graph[K, V, W], start K, visitor func(Vertex[K, V]) error) error {
	if g == nil {
		return errNilGraph
	}
	neighbours := g.Neighbours
	if g.IsDigraph() {
		dg, ok := g.(Digraph[K, V, W])
		if ok {
			neighbours = dg.OutNeighbours
		} else {
			neighbours = func(v K) ([]Vertex[K, V], error) {
				es, err := g.IncidentEdges(v)
				if err != nil {
					return nil, err
				}
				var res []Vertex[K, V]
				for _, e := range es {
					if e.Tail == v {
						w, err := g.GetVertex(e.Head)
						if err != nil {
							return nil, err
						}
						res = append(res, w)
					}
				}
				return res, nil
			}
		}
	}
	return bfs(g, start, visitor, neighbours)
}

// Perform breadth first search in a directed graph, and specify the search direction using the in parameter:
// if in is set to true, search from source in the order of the incident vertices of the current vertex.
func BFSDigraph[K comparable, V any, W number](dg Digraph[K, V, W], start K, in bool, visitor func(Vertex[K, V]) error) error {
	if dg == nil {
		return errNilGraph
	}
	var neighbours func(K) ([]Vertex[K, V], error)
	if in {
		neighbours = dg.InNeighbours
	} else {
		neighbours = dg.OutNeighbours
	}
	return bfs(dg, start, visitor, neighbours)
}

/*
Let G be a directed graph with vertex set {1,...,n}.The
algorithm checks whether G is acyclic; in this case, it also determines a topological sorting.

Data structures needed：

	a) adjacency lists A1,...,An;
	b) a function ind,where ind(v)=din(v);
	c) a function topnr, where topnr(v) gives the index of vertex v in the topological sorting;
	d) a list L of the vertices v having ind(v)=0;
	e) a Boolean variable acyclic and an integer variable N (for counting).

Procedure TOPSORT (G; topnr,acyclic)：

	(1)  N ←1, L←∅;
	(2)  for i=1 to n do ind(i) ← 0 od;
	(3)  for i=1 to n do
	(4)      for j ∈ Ai do ind(j) ← ind(j)+1 od
	(5)  od;
	(6)  for i =1to n do if ind(i)=0 then append i to L fi od;
	(7)  while L= ∅ do
	(8)      delete the first vertex v from L;
	(9)      topnr(v) ← N; N ← N +1;
	(10)     for w ∈ Av do
	(11)         ind(w) ← ind(w)−1;
	(12)         if ind(w)=0 then append w to L fi
	(13)     od
	(14) od;
	(15) if N = n+1 then acyclic ← true else acyclic ← false fi
*/
func topologicalSort[K comparable, V any, W number](g Digraph[K, V, W]) ([]Vertex[K, V], error) {
	vertexes := g.AllVertexes()

	inDegree := make(map[K]int)
	for _, v := range vertexes {
		d, err := g.InDegree(v.Key)
		if err != nil {
			return nil, err
		}
		inDegree[v.Key] = d
	}

	var vs []Vertex[K, V]
	for len(inDegree) > 0 {
		var d0 []K
		for k, d := range inDegree {
			if d == 0 {
				d0 = append(d0, k)
			}
		}
		if len(d0) == 0 {
			return nil, errNotDAG
		}
		for _, k := range d0 {
			for _, v := range vertexes {
				if v.Key == k {
					vs = append(vs, v)
					break
				}
			}
			ns, err := g.OutNeighbours(k)
			if err != nil {
				return nil, err
			}
			for _, v := range ns {
				inDegree[v.Key] = inDegree[v.Key] - 1
			}
			delete(inDegree, k)
		}
	}

	return vs, nil
}

// Perform topological sorting on a directed graph and return a sequence of vertices.
// If there is a cycle in the graph, return an error.
func TopologicalSort[K comparable, V any, W number](g Digraph[K, V, W]) ([]Vertex[K, V], error) {
	if g == nil {
		return nil, errNilGraph
	}
	return topologicalSort(g)
}

// LexBFS is a variant of breadth‑first search that imposes a lexicographic
// ordering on the vertices. It produces a vertex sequence that is consistent with
// the layers of a normal BFS but refines the choiceof vertices inside each layer by
// a deterministic rule.The algorithm is often used as a sub‑routine for recognizing
// chordal graphs and for computing perfect elimination orderings.
func LexBFS[K comparable, V any, W number](g Graph[K, V, W], start K, f func(Vertex[K, V]) error) error {
	// At each step a vertex is chosen according to a priority that is defined by a label stored on every vertex.
	// The label is a finite sequence of integers that records, for each previously visited vertex,
	// whether it was a neighbour of the current vertex or not.
	// When a vertex is selected, its label is updated by appending the index of the current step
	// to all of its neighbours’ labels.The next vertex to be visited is the one with the lexicographically
	// largest label among the unvisited vertices.
	if g == nil {
		return errNilGraph
	}
	// A label of a vertex  is written as
	//    L[v] = <i1,i2,...,ik>
	// where the entries are the numbers of the steps at which neighbours of  were chosen.
	// The sequence is kept in ascending order so that comparison of labels can be performed lexicographically.
	elems := make(map[K]*HeapElem[K, any, []int])
	orders := NewHeap[K, any](func(v1, v2 []int) bool {
		n := min(len(v1), len(v2))
		for i := 0; i < n; i++ {
			if v1[i] == v2[i] {
				continue
			}
			return v1[i] > v2[i]
		}
		return len(v1) > len(v2)
	})
	heap.Init(orders)
	for _, v := range g.AllVertexes() {
		if v.Key == start {
			continue
		}
		he := &HeapElem[K, any, []int]{
			Key: v.Key,
		}
		elems[v.Key] = he
		heap.Push(orders, he)
	}
	for i := g.Order(); len(elems) != 0; i-- {
		var k K
		if i != 0 {
			el := heap.Pop(orders).(*HeapElem[K, any, []int])
			k = el.Key
		} else {
			k = start
		}
		// visited u (the number of u is |G|-i+1)
		u, err := g.GetVertex(k)
		if err != nil {
			return err
		}
		if err := f(u); err != nil {
			return err
		}
		delete(elems, k)
		// After a vertex is selected, its adjacency list is traversed.
		// For every unvisited neighbour  of , the current step index  is appended to :
		// The step index is incremented after each selection.
		// All other vertices’ labels remain unchanged.
		ns, err := g.Neighbours(u.Key)
		if err != nil {
			return err
		}
		for _, v := range ns {
			ve, ok := elems[v.Key]
			if ok {
				ve.Rank = append(ve.Rank, i)
				heap.Fix(orders, ve.Idx)
			}
		}
	}
	return nil
}

func LexDFS[K comparable, V any, W number](g Graph[K, V, W], start K, f func(Vertex[K, V]) error) error {
	if g == nil {
		return errNilGraph
	}
	elems := make(map[K]*HeapElem[K, any, []int])
	orders := NewHeap[K, any](func(v1, v2 []int) bool {
		n1, n2, n := len(v1), len(v2), min(len(v1), len(v2))
		for i := 0; i < n; i++ {
			if v1[n1-i-1] == v2[n2-i-1] {
				continue
			}
			return v1[n1-i-1] > v2[n2-i-1]
		}
		return n1 > n2
	})
	heap.Init(orders)
	for _, v := range g.AllVertexes() {
		if v.Key == start {
			continue
		}
		he := &HeapElem[K, any, []int]{
			Key: v.Key,
		}
		elems[v.Key] = he
		heap.Push(orders, he)
	}
	for i := 0; len(elems) != 0; i++ {
		var k K
		if i != 0 {
			el := heap.Pop(orders).(*HeapElem[K, any, []int])
			k = el.Key
		} else {
			k = start
		}
		// visited u (the number of u is i)
		u, err := g.GetVertex(k)
		if err != nil {
			return err
		}
		if err := f(u); err != nil {
			return err
		}
		delete(elems, k)
		// After a vertex is selected, its adjacency list is traversed.
		// For every unvisited neighbour  of , the current step index  is appended to :
		// The step index is incremented after each selection.
		// All other vertices’ labels remain unchanged.
		ns, err := g.Neighbours(u.Key)
		if err != nil {
			return err
		}
		for _, v := range ns {
			ve, ok := elems[v.Key]
			if ok {
				ve.Rank = append(ve.Rank, i)
				heap.Fix(orders, ve.Idx)
			}
		}
	}
	return nil
}
