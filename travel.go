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
func dfs[K comparable, V any, W number](g Graph[K, V, W], start K, visitor func(v Vertex[K, V]) error, neighbours func(K) ([]Vertex[K, V], error)) error {
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
func DFS[K comparable, V any, W number](g Graph[K, V, W], start K, visitor func(v Vertex[K, V]) error) error {
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
func DFSDigraph[K comparable, V any, W number](dg Digraph[K, V, W], start K, in bool, visitor func(v Vertex[K, V]) error) error {
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
func bfs[K comparable, V any, W number](g Graph[K, V, W], start K, visitor func(v Vertex[K, V]) error, neighbours func(K) ([]Vertex[K, V], error)) error {
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
func BFS[K comparable, V any, W number](g Graph[K, V, W], start K, visitor func(v Vertex[K, V]) error) error {
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
func BFSDigraph[K comparable, V any, W number](dg Digraph[K, V, W], start K, in bool, visitor func(v Vertex[K, V]) error) error {
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
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}

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
	return topologicalSort(g)
}

/*
In the first step of the algorithm, we perform a sequence of depth first searches (dfs), visiting the entire graph.
That is, as long as there are still unvisited vertices, we take one of them, and initiate a depth first search from that vertex.
For each vertex, we keep track of the exit_time[v] . This is the 'timestamp' at which the execution of dfs on vertex v finishes,
i.e., the moment at which all vertices reachable from v have been visited and the algorithm is back at v.

It means that any edge in the condensation graph goes from a component with a larger value of  exit_time  to a component with a smaller value.

If we sort all vertices V  in decreasing order of their exit_time[v] ,
then the first vertex u will belong to the "root" strongly connected component, which has no incoming edges in the condensation graph.

Now we want to run some type of search from this vertex u so that it will visit all vertices in its strongly connected component, but not other vertices.
*/
func sccKosaraju[K comparable, V any, W number](g Digraph[K, V, W], condensation bool) ([][]K, Digraph[K, []K, W], error) {
	if g == nil {
		return nil, nil, errEmptyGraph
	}
	vtx, err := g.AllVertexes()
	if err != nil {
		return nil, nil, err
	}

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
	cond, _ := NewDigraph[K, []K, W](g.Name() + "_condensation")
	for _, c := range components {
		_ = cond.AddVertex(Vertex[K, []K]{Key: c[0], Value: c})
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
func StronglyConnectedComponent[K comparable, V any, W number](g Digraph[K, V, W], condensation bool) ([][]K, Digraph[K, []K, W], error) {
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
func sccTarjan[K comparable, V any, W number](g Digraph[K, V, W], condensation bool) ([][]K, Digraph[K, []K, W], error) {
	if g == nil {
		return nil, nil, errEmptyGraph
	}
	vtx, err := g.AllEdges()
	if err != nil {
		return nil, nil, err
	}
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
	cond, _ := NewDigraph[K, []K, W](g.Name() + "_condensation")
	for _, c := range components {
		_ = cond.AddVertex(Vertex[K, []K]{Key: c[0], Value: c})
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
func StronglyConnectedComponentTarjan[K comparable, V any, W number](g Digraph[K, V, W], condensation bool) ([][]K, Digraph[K, []K, W], error) {
	return sccTarjan(g, condensation)
}

// Calculate the strongly connected components of a directed graph and
// return the set of vertices for each strongly connected component.
func StronglyConnectedComponentKosaraju[K comparable, V any, W number](g Digraph[K, V, W], condensation bool) ([][]K, Digraph[K, []K, W], error) {
	return sccKosaraju(g, condensation)
}
