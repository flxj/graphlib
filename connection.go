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
func Connected[K comparable, W number](g Graph[K, W], start, end K) (bool, error) {
	if g == nil {
		return false, errNilGraph
	}
	var connected bool
	visitor := func(v Vertex[K, W]) error {
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

func auxiliaryGraphEDP[K comparable, W number](g Graph[K, W]) (Graph[K, int], error) {
	if g == nil {
		return nil, errNilGraph
	}
	aux := NewGraph[K, int](g.IsDigraph(), "")
	vs := g.AllVertexes()
	es := g.AllEdges()
	for _, v := range vs {
		if err := aux.AddVertex(Vertex[K, int]{Key: v.Key}); err != nil {
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
func auxiliaryGraphVDP[K comparable, W number](g Graph[K, W], source, target K) (Graph[int, int], int, int, error) {
	if g == nil {
		return nil, 0, 0, errNilGraph
	}
	aux := NewGraph[int, int](g.IsDigraph(), "")
	vs := g.AllVertexes()
	es := g.AllEdges()
	idx := make(map[K]int)
	var ek, s, t int
	for i, v := range vs {
		if err := aux.AddVertex(Vertex[int, int]{Key: -(i + 1)}); err != nil {
			return nil, 0, 0, err
		}
		if err := aux.AddVertex(Vertex[int, int]{Key: i + 1}); err != nil {
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
func EdgeDisjointPath[K comparable, W number](g Graph[K, W], source, target K) (int, error) {
	aux, err := auxiliaryGraphEDP(g)
	if err != nil {
		return 0, err
	}
	return MaxFlow(aux, source, target)
}

// Find maximum number of vertex disjoint paths between two vertices.
func VertexDisjointPath[K comparable, W number](g Graph[K, W], source, target K) (int, error) {
	aux, s, t, err := auxiliaryGraphVDP(g, source, target)
	if err != nil {
		return 0, err
	}
	return MaxFlow(aux, s, t)
}

// Find maximum number of edge disjoint paths between two vertices.
func DigraphEdgeDisjointPath[K comparable, W number](g Digraph[K, W], source, target K) (int, error) {
	return EdgeDisjointPath(g, source, target)
}

// Find maximum number of vertex disjoint paths between two vertices.
func DigraphVertexDisjointPath[K comparable, W number](g Graph[K, W], source, target K) (int, error) {
	return VertexDisjointPath(g, source, target)
}

// Query the incut or outcut of vertex set X on directed graph g (the incut is composed of all directed arcs whose heads belong to X).
func DigraphCut[K comparable, W number](g Digraph[K, W], X []K, incut bool) ([]Edge[K, W], error) {
	var res []Edge[K, W]
	xm := make(map[K]struct{})
	for _, v := range X {
		xm[v] = struct{}{}
	}
	var es []Edge[K, W]
	var err error
	for v := range xm {
		if incut {
			es, err = g.InEdges(v)
		} else {
			es, err = g.OutEdges(v)
		}
		if err != nil {
			return nil, err
		}
		for _, e := range es {
			_, ok := xm[e.Head]
			_, ok1 := xm[e.Tail]
			if ok && ok1 {
				continue
			}
			res = append(res, e)
		}
	}
	return res, nil
}

// Query the directed arc between two sets of vertices, from and to, on a directed graph g.
func DigraphArcs[K comparable, W number](g Digraph[K, W], from, to []K) ([]Edge[K, W], error) {
	var res []Edge[K, W]
	xm := make(map[K]struct{})
	ym := make(map[K]struct{})
	for _, v := range from {
		xm[v] = struct{}{}
	}
	for _, v := range to {
		ym[v] = struct{}{}
	}
	for v := range xm {
		es, err := g.OutEdges(v)
		if err != nil {
			return nil, err
		}
		for _, e := range es {
			if _, ok := ym[e.Head]; ok {
				res = append(res, e)
			}
		}
	}
	return res, nil
}

// Query the boundary of vertex set X on the graph, i.e. the edge between X and G-X.
func GraphCoboundary[K comparable, W number](g Graph[K, W], X []K) ([]Edge[K, W], error) {
	var res []Edge[K, W]
	xm := make(map[K]struct{})
	for _, v := range X {
		xm[v] = struct{}{}
	}
	for v := range xm {
		es, err := g.IncidentEdges(v)
		if err != nil {
			return nil, err
		}
		for _, e := range es {
			_, ok := xm[e.Head]
			_, ok1 := xm[e.Tail]
			if ok && ok1 {
				continue
			}
			res = append(res, e)
		}
	}
	return res, nil
}

// Query the edges between two sets of vertices X and Y on the graph.
func GraphEdges[K comparable, W number](g Graph[K, W], X, Y []K) ([]Edge[K, W], error) {
	xm := make(map[K]struct{})
	ym := make(map[K]struct{})
	for _, v := range X {
		xm[v] = struct{}{}
	}
	for _, v := range Y {
		ym[v] = struct{}{}
	}
	em := make(map[K]Edge[K, W])
	for v := range xm {
		es, err := g.IncidentEdges(v)
		if err != nil {
			return nil, err
		}
		for _, e := range es {
			_, ok := ym[e.Head]
			_, ok2 := ym[e.Tail]
			if ok || ok2 {
				em[e.Key] = e
			}
		}
	}
	var res []Edge[K, W]
	for _, e := range em {
		res = append(res, e)
	}
	return res, nil
}

/*
In the first step of the algorithm, we perform a sequence of depth first searches (dfs), visiting the entire graph.
That is, as long as there are still unvisited vertices, we take one of them, and initiate a depth first search from that vertex.
For each vertex, we keep track of the exit_time[v] .
This is the 'timestamp' at which the execution of dfs on vertex v finishes,
i.e., the moment at which all vertices reachable from v have been visited and the algorithm is back at v.

It means that any edge in the condensation graph goes from a component with a larger value of exit_time
to a component with a smaller value.

If we sort all vertices V in decreasing order of their exit_time[v],
then the first vertex u will belong to the "root" strongly connected component,
which has no incoming edges in the condensation graph.

Now we want to run some type of search from this vertex u so that it will
visit all vertices in its strongly connected component, but not other vertices.
*/
func sccKosaraju[K comparable, W number](g Digraph[K, W], condensation bool) ([][]K, Digraph[K, W], error) {
	if g == nil {
		return nil, nil, errNilGraph
	}
	vtx := g.AllVertexes()

	var dfs func(K, *[]K) error
	var exitTime []K // record the exit time of a vertex when dfs.
	// first dfs
	revAdj := make(map[K][]K)
	visited := make(map[K]struct{})
	dfs = func(v K, arr *[]K) error {
		visited[v] = struct{}{}
		out, err := g.OutNeighbours(v)
		if err != nil {
			return err
		}
		// update reverse adjlist
		for _, w := range out {
			revAdj[w.Key] = append(revAdj[w.Key], v)
		}
		//
		for _, w := range out {
			if _, ok := visited[w.Key]; !ok {
				dfs(w.Key, arr)
			}
		}
		*arr = append(*arr, v)
		return nil
	}
	for _, v := range vtx {
		if _, ok := visited[v.Key]; !ok {
			if err := dfs(v.Key, &exitTime); err != nil {
				return nil, nil, err
			}
		}
	}
	// second dfs
	var components [][]K
	root := make(map[K]K)
	visited = make(map[K]struct{})
	dfs = func(v K, arr *[]K) error {
		visited[v] = struct{}{}
		for _, w := range revAdj[v] {
			if _, ok := visited[w]; !ok {
				_ = dfs(w, arr)
			}
		}
		*arr = append(*arr, v)
		return nil
	}
	for i := len(exitTime) - 1; i >= 0; i-- {
		v := exitTime[i]
		if _, ok := visited[v]; !ok {
			// new component
			c := []K{}
			_ = dfs(v, &c)
			components = append(components, c)
			if condensation {
				for _, w := range c {
					root[w] = c[0]
				}
			}
		}
	}
	if !condensation {
		return components, nil, nil
	}
	// build condensation graph
	cond := NewDigraph[K, W](g.Name() + "_condensation")
	for _, c := range components {
		_ = cond.AddVertex(Vertex[K, W]{Key: c[0], Value: c})
	}
	for _, v := range vtx {
		out, err := g.OutNeighbours(v.Key)
		if err != nil {
			return nil, nil, err
		}
		for _, w := range out {
			if root[w.Key] != root[v.Key] {
				_ = cond.AddEdge(Edge[K, W]{Head: root[w.Key], Tail: root[v.Key]})
			}
		}
	}
	return components, cond, nil
}

// Calculate the strongly connected components of a directed graph and
// return the set of vertices for each strongly connected component.
func StronglyConnectedComponent[K comparable, W number](g Digraph[K, W], condensation bool) ([][]K, Digraph[K, W], error) {
	return sccKosaraju(g, condensation)
}

/*
Let's consider the tree induced by the sequence of DFS calls, which we will call DFS tree.
Once we first call a DFS on a vertex from an SCC, all the vertices of its SCC will be visited before this call ends,
since they are all reachable from each other.
In the DFS tree, this first vertex will be a common ancestor to all other vertices of the SCC;
we define this vertex to be the root of the SCC.

Once we finish traversing the neighbours list of a vertex, we somehow are able to determine whether it is a root or not.
In case the vertex is a root, we will then immediately find and claim all the vertices of its SCC.
When all calls finish, all roots will have been detected and all vertices will have been claimed as part of some SCC.

we define the entry time in[v]  for each vertex v  which corresponds to the 'timestamp' at which the DFS was called on v.
By definition, the root is the first vertex of an SCC to be visited by the DFS so it will have the minimal value of in[v]  of its SCC.
*/
func sccTarjan[K comparable, W number](g Digraph[K, W], condensation bool) ([][]K, Digraph[K, W], error) {
	if g == nil {
		return nil, nil, errNilGraph
	}
	vtx := g.AllEdges()
	// 1.DFS search produces a DFS tree/forest
	// 2.Strongly Connected Components form subtrees of the DFS tree.
	// 3.If we can find the root of such subtrees, we can print/store all the nodes in that subtree (including the root) and that will be one SCC.
	// 4.There is no back edge from one SCC to another (There can be cross edges, but cross edges will not be used while processing the graph).
	//
	// inTime[v]: This is the time when a vertex v is visited 1st time while DFS traversal.
	// Assign a new number to each vertex in the graph. If a vertex v is traversed i-th in the dfs tree, its number is i, called a timestamp,
	// represented by inTime[v]=i. The timestamp is unique, and the timestamp corresponding to the vertex is also unique.
	//
	// In the DFS tree, Tree edges take us forward, from the ancestor node to one of its descendants.
	// Back edges take us backward, from a descendant node to one of its ancestors.
	//
	// lowTime[v]: as the minimum timestamp that vertex v can reach,that is, the minimum timestamp that the subtrees of v and v can reach,
	// and also describe it as the minimum timestamp that v can trace in the dfs stack.
	// If the low of a vertex v is equal to its timestamp, then that vertex must be the "root" of its strongly connected component.
	var t int
	lowTime := make([]int, len(vtx))
	inTime := make([]int, len(vtx))
	root := make([]int, len(vtx))
	idx := make(map[K]int)
	for i, v := range vtx {
		idx[v.Key] = i
	}
	stk := newStack[int]()
	var dfs func(K, int, *[][]K) error
	dfs = func(v K, i int, components *[][]K) error {
		lowTime[i] = t
		inTime[i] = t
		t++
		stk.push(i)
		out, err := g.OutNeighbours(v)
		if err != nil {
			return err
		}
		for _, w := range out {
			j := idx[w.Key]
			if inTime[j] == -1 { // (v,w) is tree edge
				dfs(w.Key, j, components)
				lowTime[i] = min(lowTime[i], lowTime[j])
			} else if root[j] == -1 { // // back-edge, cross-edge or forward-edge to an unclaimed vertex
				lowTime[i] = min(lowTime[i], inTime[j])
			}
		}
		if inTime[i] == lowTime[i] { // v is a root
			// create a new component.
			comp := []K{v}
			for {
				j, ok := stk.pop()
				if !ok {
					break
				}
				root[j] = i
				if i == j {
					break
				}
				comp = append(comp, vtx[j].Key)
			}
			*components = append(*components, comp)
		}
		return nil
	}
	for i := 0; i < len(vtx); i++ {
		root[i] = -1
		inTime[i] = -1
		lowTime[i] = -1
	}
	var components [][]K
	for i, v := range vtx {
		if inTime[i] == -1 {
			if err := dfs(v.Key, i, &components); err != nil {
				return nil, nil, err
			}
		}
	}
	if !condensation {
		return components, nil, nil
	}
	// build condensation graph
	cond := NewDigraph[K, W](g.Name() + "_condensation")
	for _, c := range components {
		_ = cond.AddVertex(Vertex[K, W]{Key: c[0], Value: c})
	}
	for i, v := range vtx {
		out, err := g.OutNeighbours(v.Key)
		if err != nil {
			return nil, nil, err
		}
		for _, w := range out {
			j := idx[w.Key]
			if root[i] != root[j] {
				_ = cond.AddEdge(Edge[K, W]{Head: vtx[root[j]].Key, Tail: vtx[root[i]].Key})
			}
		}
	}
	return components, cond, nil
}

// Calculate the strongly connected components of a directed graph and
// return the set of vertices for each strongly connected component.
func StronglyConnectedComponentTarjan[K comparable, W number](g Digraph[K, W], condensation bool) ([][]K, Digraph[K, W], error) {
	return sccTarjan(g, condensation)
}

// Calculate the strongly connected components of a directed graph and
// return the set of vertices for each strongly connected component.
func StronglyConnectedComponentKosaraju[K comparable, W number](g Digraph[K, W], condensation bool) ([][]K, Digraph[K, W], error) {
	return sccKosaraju(g, condensation)
}
