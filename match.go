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
	"errors"
)

// Hungarian algorithm (also known as the Kuhn–Munkres algorithm) solving the maximum/minimum weighted matching problem for bipartite graphs.
// we assume that all weight of edges are non-negative.
func BipartiteWeightedMatching[K comparable, V any, W number](g *Bipartite[K, V, W], maximum bool) ([]Edge[K, W], error) {
	if g.Order() == 0 || g.Size() == 0 {
		return []Edge[K, W]{}, nil
	}
	var A, B []Vertex[K, V]
	var err error
	if A, err = g.Part(true); err != nil {
		return nil, err
	}
	if B, err = g.Part(false); err != nil {
		return nil, err
	}
	if len(A) != len(B) { // TODO: add dummy verties to make equal
		return nil, errors.New("not regular complete bipartite")
	}
	var M []int
	if maximum {
		weight := make([][]W, len(A)+1)
		for i := range weight {
			weight[i] = make([]W, len(A)+1)
		}
		for i := 1; i <= len(A); i++ {
			for j := 1; j <= len(A); j++ {
				es, err := g.GetEdge(A[i-1].Key, B[j-1].Key)
				if err != nil {
					if err != errEdgeNotExists {
						return nil, err
					} else {
						weight[i][j] = weight[0][0] // zero
					}
				} else {
					weight[i][j] = es[0].Weight
				}
			}
		}
		M = maxMatchingHungarian(len(A), weight)
	} else {
		// find minimum matching
		weight := make([][]W, len(A)+1)
		for i := range weight {
			weight[i] = make([]W, len(A)+1)
		}
		inf := maxValue(weight[0][0])
		for i := 1; i <= len(A); i++ {
			for j := 1; j <= len(A); j++ {
				es, err := g.GetEdge(A[i-1].Key, B[j-1].Key)
				if err != nil {
					if err != errEdgeNotExists {
						return nil, err
					} else {
						weight[i][j] = inf
					}
				} else {
					weight[i][j] = es[0].Weight
				}
			}
		}
		M = minMatchingHungarian(len(A), weight)
	}
	// generate results.
	res := make([]Edge[K, W], len(M)-1)
	for i := 1; i < len(M); i++ {
		u, v := i-1, M[i]-1
		es, err := g.GetEdge(A[u].Key, B[v].Key)
		if err != nil {
			if err == errEdgeNotExists {
				res[i-1] = Edge[K, W]{Head: A[u].Key, Tail: B[v].Key}
			} else {
				return nil, err
			}
		} else {
			res[i-1] = es[0]
		}
	}
	return res, nil
}

func maxMatchingHungarian0[K comparable, V any, W number](V1, V2 []Vertex[K, V], weight [][]W, g *Bipartite[K, V, W]) ([]Edge[K, W], error) {
	// init
	n := len(V1)
	Y, Z := make([]Vertex[int, any], n), make([]Vertex[int, any], n)
	for i := 0; i < n; i++ {
		Y[i].Key = i
		Z[i].Key = i + n
	}
	y, z := make([]W, n), make([]W, n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if weight[i][j] > y[i] {
				y[i] = weight[i][j]
			}
		}
	}
	ek := 0
	edge := func(u, v int) Edge[int, int] {
		ek++
		return Edge[int, int]{Key: ek, Head: u, Tail: v}
	}
	Gyz, _ := NewBipartite[int, any, int](false, "")
	for i := 0; i < n; i++ {
		_ = Gyz.AddVertexTo(Y[i], true)
		_ = Gyz.AddVertexTo(Z[i], false)
	}
	equalSubgraph := func() {
		_ = Gyz.RemoveAllEdge()
		for u := 0; u < n; u++ {
			for v := 0; v < n; v++ {
				if y[u]+z[v] == weight[u][v] {
					_ = Gyz.AddEdge(edge(u, v))
				}
			}
		}
	}
	//
	var err error
	var M []Edge[int, int]
	B := make(map[int]struct{})
	for {
		// Construct the equality subgraph Gyz associated with (y,z)
		equalSubgraph()
		// Set M and B as a maximum matching and a minimum vertex cover in Gyz, respectively.
		if M, err = BipartiteMaxMatching(Gyz); err != nil {
			return nil, err
		}
		if len(M) == n {
			break
		}
		if B, err = bipartiteMVC(Y, Z, M, Gyz); err != nil {
			return nil, err
		}
		// set R = V1 ∩ B and T = V2 ∩ B; Compute d = min{y(u)+z(v)−W(uv) : u ∈ V1−R,v ∈ V2−T}
		d := maxValue(weight[0][0])
		for u := 0; u < n; u++ {
			if _, ok := B[u]; ok {
				continue
			}
			for v := 0; v < n; v++ {
				if _, ok := B[v+n]; !ok {
					d = min(d, y[u]+z[v]-weight[u][v])
				}
			}
		}
		// Update y(u)−=d for u ∈ V1−R, and z(v) += d for v∈T
		for u := 0; u < n; u++ {
			if _, ok := B[u]; !ok {
				y[u] -= d
			}
		}
		for v := 0; v < n; v++ {
			if _, ok := B[v+n]; !ok {
				z[v] += d
			}
		}
	}
	// generate results.
	getEdge := func(u, v int) Edge[K, W] {
		if weight[u][v] == 0 {
			return Edge[K, W]{Head: V1[u].Key, Tail: V2[v].Key}
		}
		es, _ := g.GetEdge(V1[u].Key, V2[v].Key)
		return es[0]
	}
	res := make([]Edge[K, W], len(M))
	for i, e := range M {
		res[i] = getEdge(e.Head, e.Tail)
	}
	return res, nil
}

/*
G = (U,V,E),|U|=|V|=n,|E|=n^2

Let's call a potential two arbitrary arrays of numbers u[1...n] and v[1...n], such that the following condition is satisfied:

	u[i]+v[j] >= weight[i][j]

Let's call the value f of the potential the sum of its elements:

	f = Σu[i] + Σv[j]

On one hand, it is easy to see that the cost of the desired solution 's'  is not larger than the value of any potential:

	s <= f

On the other hand, it turns out that there is always a solution and a potential that turns this inequality into equality.
so we should find the minimum value of f, which is just the maximum solution.

Let's fix some potential. Let's call an edge (i,j) rigid if u[i]+v[j]=weight[i][j].
Denote with H a bipartite graph composed only of rigid edges.
The Hungarian algorithm will maintain, for the current potential, the maximum-number-of-edges matching M of the graph H.
As soon as M contains n edges, then the solution to the problem will be just M.

Step 1. At the beginning, the potential is assumed to be zero (u[i]=v[i]=0  for all i), and the matching M is assumed to be empty.

Step 2. At each step of the algorithm, we try, without changing the potential, to increase the cardinality of the current matching M by one.
To do this, the usual Kuhn Algorithm for finding the maximum matching in bipartite graphs H is used.

Step 3. If at the current step, it is not possible to increase the cardinality of the current matching,
then a recalculation of the potential is performed in such a way that, at the next steps, there will be more opportunities to increase the matching.

Denote by Z1 the set of vertices of the part U that were visited during the last traversal of Kuhn's algorithm, and through Z2 the set of visited vertices of the part V.

Calculate the value delta:

	delta = min{ u[i]+v[j]-weight[i][j] | i ∈ Z1, j ∈ U\Z2 } // matching vertices of U to unmatching vertices of V

Now recalculate the potential in this way:

	u[i] = u[i] - delta, i ∈ Z1
	v[j] = v[j] + delta, j ∈ Z2

(recalculate donot change the potential of current matched edges i.e. all edges of the matching will remain rigid
but it introducing new assemble edges by u[i]-delta to current H,so we can explore more edges in next round matching)
*/
func maxMatchingHungarian[W number](n int, weight [][]W) []int {
	// store potential.
	u, v := make([]W, n+1), make([]W, n+1)
	// p[1...n] record matching. p[0] record current row.
	p := make([]int, n+1)
	//
	// quickly recalculate the potential,maintain auxiliary minima for each of the columns:
	// minv[j] = min{ u[i]+v[j]-weighjt[i][j] | i ∈ Z1,j ∈ U\Z2}
	// delta = min{ minv[j] | j ∈ U\Z2 }
	inf := maxValue(weight[0][0])
	initMinv := func(l int) []W {
		minv := make([]W, l)
		for i := 0; i < l; i++ {
			minv[i] = inf
		}
		return minv
	}
	// The array wpath[1...n] contains information about where these minimums are reached so that we can later reconstruct the augmenting path.
	// path[j] = j'  means j' is prev node of j in a path.
	path := make([]int, n+1)
	//
	for i := 1; i <= n; i++ {
		/*
			consider the row i of the weight.
			if there is no increasing path starting in this row, recalculate the potential.
			if an augmenting path is found, extend the matching by it (thus including the last edge in the matching),
			and restart to consider the next row.
		*/
		p[0] = i
		/*
			To check for the presence of an augmenting path, there is no need to start the Kuhn traversal again after each potential recalculation.
			Instead, you can make the Kuhn traversal in an iterative form:
			after each recalculation of the potential, look at the added rigid edges and,
			if their left ends(U) were reachable, mark their right ends(V) reachable as well and continue the traversal from them.

			At each step of the loop, the potential is recalculated. Subsequently,
			a column that has become reachable is identified (which will always exist as new reachable vertices emerge after every potential recalculation).
			If the column is unsaturated, an augmenting chain is discovered. Conversely, if the column is saturated, the matching row also becomes reachable.
		*/
		var j0 int
		minv := initMinv(n + 1)
		used := make([]bool, n+1)
		for {
			// start from j0 to check augmenting path: j0->i0->j...
			used[j0] = true
			i0, j1 := p[j0], 0
			delta := inf
			for j := 1; j <= n; j++ {
				if !used[j] { // i0 ∈ Z1, j ∈ V\Z2
					d := u[i0] + v[j] - weight[i0][j]
					if minv[j] > d {
						minv[j] = d
						path[j] = j0
					}
					if minv[j] < delta {
						delta = minv[j]
						j1 = j // j1 to record the node with minimum delta.
					}
				}
			}
			// potential recalculation
			for j := 0; j <= n; j++ {
				if used[j] { // j ∈ Z2, so p[j] ∈ Z1
					u[p[j]] -= delta
					v[j] += delta
				} else { // j ∈ V\Z2
					minv[j] -= delta // because u[i0] has -delta, so for edge (i0,j) the potential should -delta,so the minv[j]-delta.
				}
			}
			j0 = j1 // move to next node j1 along the alternating path.
			if p[j0] == 0 {
				// p[j0] == 0 means j0 column not in current matching.
				break
			}
		}
		// update matching by flip alternating path,to add current row i into current matching. (note p[0] == i)
		for {
			j1 := path[j0]
			p[j0] = p[j1]
			j0 = j1
			if j0 == 0 {
				break
			}
		}
	}
	return p
}

/*
G = (U,V,E)

dual feasibility condition requires: u[i]+v[j] <= weight[i][j]
reduced cost of edge (i, j) is defined as: weight[i][j] − u[i] − v[j]  >= 0
An edge is tight (or admissible) when its reduced cost is exactly zero.
The dual objective is to maximise Σu[i] + Σv[j], which gives a lower bound on the optimal primal cost.

Step 1. At the beginning, the potential is assumed to be zero (u[i]=v[i]=0  for all i), and the matching  M  is assumed to be empty.

Step 2. Further, at each step of the algorithm, we try, without changing the potential, to increase the cardinality of the current matching M by one
(recall that the matching is searched in the graph of rigid edges H). To do this, the usual Kuhn Algorithm for finding the maximum matching in
bipartite graphs is used. Let us recall the algorithm here.
All edges of the matching M are oriented in the direction from the right part to the left one,
and all other edges of the graph H are oriented in the opposite direction.

Recall (from the terminology of searching for matchings) that a vertex is called saturated if an edge of the current matching is adjacent to it.
A vertex that is not adjacent to any edge of the current matching is called unsaturated. A path of odd length,
in which the first edge does not belong to the matching, and for all subsequent edges there is an alternating belonging to the matching
(belongs/does not belong) - is called an augmenting path.

From all unsaturated vertices in the left part, a depth-first or breadth-first traversal is started.
If, as a result of the search, it was possible to reach an unsaturated vertex of the right part,
we have found an augmenting path from the left part to the right one.
If we include odd edges of the path and remove the even ones in the matching
(i.e. include the first edge in the matching, exclude the second, include the third, etc.),
then we will increase the matching cardinality by one. (If there was no augmenting path, then the current matching M is maximal in the graph H.)

Step 3. If at the current step, it is not possible to increase the cardinality of the current matching,
then a recalculation of the potential is performed in such a way that, at the next steps, there will be more opportunities to increase the matching.

Denote by Z1 the set of vertices of the left part that were visited during the last traversal of Kuhn's algorithm,and through Z2 the set of visited vertices of the right part.

	    Let's calculate the value:
		    d = min{ weight[i][j]-u[i]-v[j] | i ∈ Z1, j ∈ U\Z2  }
		Now let's recalculate the potential in this way:
		    u[i] = u[i] + d, for all i ∈ Z1
			v[j] = v[j] - d, for all j ∈ Z2
*/
func minMatchingHungarian[W number](n int, weight [][]W) []int {
	// weight[i][j] ==> weight of edge (i,j) ==> (worker[i]'s coset for doing job[j])
	// The weight is a (n+1)*(n+1) matrix,and is ​1-based for convenience and code brevity:
	// this implementation introduces a dummy zero row and zero column,
	// which allows us to write many cycles in a general form, without additional checks.

	// Arrays u[0 ...n]  and  v[0 ... m] store potential.(such that u[i]+u[j] <= weight[i][j])
	// Initially, they are set to zero, which is consistent with a matrix of zero rows
	// (Note that it is unimportant for this implementation whether or not the matrix  weight  contains negative numbers).
	u, v := make([]W, n+1), make([]W, n+1)
	// The array p[0... n]  contains a matching M:
	// for each column  j = 1 ... n , it stores the number  p[j]  of the selected row (or 0  if nothing has been selected yet).
	// For the convenience of implementation, p[0] is assumed to be equal to the number of the current row.
	p := make([]int, n+1)
	// Maintain auxiliary array minv[1 .... n] contains, for each column j:
	//    minv[j] = min{ weight[i][j]-u[i]-v[j] | i ∈ Z1 }
	// the auxiliary minima necessary for a quick recalculation of the potential.
	// It's easy to see that the desired value d is expressed in terms of them as follows:
	//    d = min{ minv[j] | j ∈ U\Z2}
	inf := maxValue(weight[0][0])
	initMinv := func(l int) []W {
		minv := make([]W, l)
		for i := 0; i < l; i++ {
			minv[i] = inf
		}
		return minv
	}

	// The array  way[1 ... n] contains information about where these minimums are reached so that we can later reconstruct the augmenting path.
	// Note that, to reconstruct the path, it is sufficient to store only column values,
	// since the row numbers can be taken from the matching (i.e., from the array p).
	// Thus, way[j] , for each column j, contains the number of the previous column in the path (or 0 if there is none).
	way := make([]int, n+1)

	// in the outer loop, we consider matrix rows one by one.
	for i := 1; i <= n; i++ {
		p[0] = i
		var j0 int
		minv := initMinv(n + 1)
		used := make([]bool, n+1)
		for {
			used[j0] = true
			i0, delta, j1 := p[j0], inf, 0
			for j := 1; j <= n; j++ {
				if !used[j] {
					cur := weight[i0][j] - u[i0] - v[j]
					if cur < minv[j] {
						minv[j], way[j] = cur, j0
					}
					if minv[j] < delta {
						delta, j1 = minv[j], j
					}
				}
			}
			for j := 0; j <= n; j++ {
				if used[j] {
					u[p[j]] += delta
					v[j] -= delta
				} else {
					minv[j] -= delta
				}
			}
			j0 = j1
			if p[j0] == 0 {
				break
			}
		}
		for {
			j1 := way[j0]
			p[j0], j0 = p[j1], j1
			if j0 == 0 {
				break
			}
		}
	}
	res := make([]int, n+1)
	for i := 1; i <= n; i++ {
		res[p[i]] = i // res[i] = j means we select edge (row=i,col=j)
	}
	return res
}

// Obtain the minimum vertex cover of a bipartite graph.
func BipartiteMinimumVertexCover[K comparable, V any, W number](g *Bipartite[K, V, W]) (map[K]struct{}, error) {
	if g.Order() == 0 || g.Size() == 0 {
		return make(map[K]struct{}), nil
	}
	var A, B []Vertex[K, V]
	var err error
	if A, err = g.Part(true); err != nil {
		return nil, err
	}
	if B, err = g.Part(false); err != nil {
		return nil, err
	}
	M, err := mmHopcroftKarp(A, B, g)
	if err != nil {
		return nil, err
	}
	return bipartiteMVC(A, B, M, g)
}

func bipartiteMVC[K comparable, V any, W number](A, B []Vertex[K, V], M []Edge[K, W], g *Bipartite[K, V, W]) (map[K]struct{}, error) {
	VA := make(map[K]int8)
	for _, v := range A {
		VA[v.Key] = 1
	}
	// VA[v] == 2 means v in A and be covered by M.
	for _, e := range M {
		if g.InPartA(e.Head) {
			VA[e.Head] = 2
		} else {
			VA[e.Tail] = 2
		}
	}
	//
	S := make(map[K]struct{})
	visited := make(map[K]struct{})
	var dfs func(K, bool) bool
	dfs = func(v K, free bool) bool {
		if _, ok := visited[v]; ok {
			return true
		}
		visited[v] = struct{}{}
		S[v] = struct{}{} // add it to S
		if free {
			// find a unmatch edge (v,w)
			ns, err := g.Neighbours(v)
			if err != nil {
				return false
			}
			for _, w := range ns {
				if !dfs(w.Key, !free) {
					return false
				}
			}
		} else {
			// get a match edge from M
			for _, e := range M {
				if e.Head == v {
					return dfs(e.Tail, !free)
				} else if e.Tail == v {
					return dfs(e.Head, !free)
				}
			}
		}
		return true
	}
	// init a set S, its contains all unmatch vertices of A
	// then use dfs process find alternating path from A, and add all reachable vertices to S
	for _, v := range A {
		if VA[v.Key] != 2 {
			if !dfs(v.Key, true) {
				return nil, errors.New("dfs error")
			}
		}
	}
	// C = (A\S ) U (B ∩ S) is the MVC.
	C := make(map[K]struct{})
	for _, v := range A {
		if _, ok := S[v.Key]; !ok {
			C[v.Key] = struct{}{}
		}
	}
	for _, v := range B {
		if _, ok := S[v.Key]; ok {
			C[v.Key] = struct{}{}
		}
	}
	return C, nil
}

func mmBlossom[K comparable, V any, W number](_ Graph[K, V, W]) ([]Edge[K, W], error) {
	return nil, errNotImplement
}

// Calculate the maximum matching of any graph and return the set of edges.
func MaxMatching[K comparable, V any, W number](g Graph[K, V, W]) ([]Edge[K, W], error) {
	return mmBlossom(g)
}

// Calculate the perfect matching of any graph, if it exists, return the set of edges,
// otherwise return non-existent.
func PerfectMatching[K comparable, V any, W number](g Graph[K, V, W]) ([]Edge[K, W], error) {
	mm, err := MaxMatching(g)
	if err != nil {
		return nil, err
	}
	vs := make(map[K]bool)
	for _, k := range mm {
		e, err := g.GetEdgeByKey(k.Key)
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
