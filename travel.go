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
func dfs[K comparable, V any, W number](g Graph[K, V, W], start K, in bool, visitor func(v Vertex[K, V]) error) error {
	neighbours := g.Neighbours
	if g.IsDigraph() {
		dg, ok := g.(Digraph[K, V, W])
		if ok {
			if in {
				neighbours = dg.InNeighbours
			} else {
				neighbours = dg.OutNeighbours
			}
		}
	}
	//
	startV, err := g.GetVertex(start)
	if err != nil {
		return err
	}
	visited := make(map[K]struct{})
	stack := []Vertex[K, V]{startV}

	for top := 1; top > 0; {
		v := stack[top-1]
		top--
		if _, ok := visited[v.Key]; !ok {
			if err := visitor(v); err != nil {
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
				if top < len(stack) {
					stack[top] = v
				} else {
					stack = append(stack, v)
				}
				top++
			}
		}
	}
	return nil
}

// Start depth first search from the specified source vertex, where g can be a directed or undirected graph.
func DFS[K comparable, W number](g Graph[K, any, W], start K, visitor func(v Vertex[K, any]) error) error {
	return dfs(g, start, false, visitor)
}

// Perform depth first search in a directed graph, and specify the search direction using the in parameter:
// if in is set to true, search from source in the order of the incident vertices of the current vertex.
func DFSDigraph[K comparable, V any, W number](dg Digraph[K, V, W], start K, in bool, visitor func(v Vertex[K, V]) error) error {
	var g Graph[K, V, W]
	g = dg
	return dfs(g, start, in, visitor)
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
func bfs[K comparable, V any, W number](g Graph[K, V, W], start K, in bool, visitor func(v Vertex[K, V]) error) error {
	neighbours := g.Neighbours
	if g.IsDigraph() {
		dg, ok := g.(Digraph[K, V, W])
		if ok {
			if in {
				neighbours = dg.InNeighbours
			} else {
				neighbours = dg.OutNeighbours
			}
		}
	}
	//
	startV, err := g.GetVertex(start)
	if err != nil {
		return err
	}
	visited := make(map[K]struct{})
	// use a fifo queue.
	queue := []Vertex[K, V]{startV}
	head := 0
	tail := 1

	// visit current vertex,and push all neighbours of it to queue.
	for head < tail {
		v := queue[head]
		head++
		if _, ok := visited[v.Key]; !ok {
			if err := visitor(v); err != nil {
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
				if tail < len(queue) {
					queue[tail] = v
				} else {
					queue = append(queue, v)
				}
				tail++
			}
		}
	}
	return nil
}

// Start breadth first search from the specified source vertex, where g can be a directed or undirected graph.
func BFS[K comparable, W number](g Graph[K, any, W], start K, visitor func(v Vertex[K, any]) error) error {
	return bfs(g, start, false, visitor)
}

// Perform breadth first search in a directed graph, and specify the search direction using the in parameter:
// if in is set to true, search from source in the order of the incident vertices of the current vertex.
func BFSDigraph[K comparable, V any, W number](dg Digraph[K, V, W], start K, in bool, visitor func(v Vertex[K, V]) error) error {
	var g Graph[K, V, W]
	g = dg
	return bfs(g, start, in, visitor)
}

// Determine whether the start and end vertices in graph g are connected.
// If it is a directed graph, determine if there is a directed path from start to end.
func Connected[K comparable, W number](g Graph[K, any, W], start, end K) (bool, error) {
	var connected bool
	visitor := func(v Vertex[K, any]) error {
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
func topologicalSort[K comparable, W number](g Digraph[K, any, W]) ([]Vertex[K, any], error) {
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

	var vs []Vertex[K, any]
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
func TopologicalSort[K comparable, W number](g Digraph[K, any, W]) ([]Vertex[K, any], error) {
	return topologicalSort(g)
}

/*
Let G =(V,E) be a digraph and s a root of G.

Procedure DFSM(G,s,nr,Nr,p)

	(1) for v ∈ V do nr(v)←0; Nr(v)←0; p(v)←0 od;
	(2) for e ∈ E do u(e)←false od;
	(3) i←1; j←0; v←s; nr(s)←1; Nr(s)←|V|;
	(4) repeat
	(5)     while there exists w ∈ Av with u(vw)=false do
	(6)         choose some w ∈ Av with u(vw)=false; u(vw)←true;
	(7)         f nr(w)=0 then p(w)←v; i←i+1;nr(w)←i;
	            v←w fi
	(8)     od;
	(9)     j←j+1;Nr(v)←j; v←p(v)
	(10)until v = s and u(sw) = true for each w ∈ As

Let G be a digraph and s a root of G. The algorithm determines the strong components of G.

Procedure STRONGCOMP(G,s)

	(1) DFSM(G,s,nr,Nr,p); k←0;
	(2) let H be the digraph with the opposite orientation of G;
	(3) repeat
	(4)     choose the vertex r in H for which Nr(r) is maximal;
	(5)     k←k+1;DFS(H,r;nr,p); Ck ←{v ∈ H :nr(v)=0};
	(6)     remove all vertices in Ck and all the edges incident with them;
	(7) until the vertex set of H is empty
*/
func sccKosaraju[K comparable, W number](g Digraph[K, any, W]) ([][]K, error) {
	return nil, errNotImplement
}

// 1.DFS search produces a DFS tree/forest
//
// 2.Strongly Connected Components form subtrees of the DFS tree.
//
// 3.If we can find the head of such subtrees, we can print/store all the nodes in that subtree (including the head) and that will be one SCC.
//
// 4.There is no back edge from one SCC to another (There can be cross edges, but cross edges will not be used while processing the graph).
//
//
// dfn[v]: This is the time when a vertex v is visited 1st time while DFS traversal.
// Assign a new number to each vertex in the graph. If a vertex v is traversed i-th in the dfs tree, its number is i, called a timestamp,
// represented by dfn[v]=i. The timestamp is unique, and the timestamp corresponding to the vertex is also unique.
//
// In the DFS tree, Tree edges take us forward, from the ancestor node to one of its descendants.
// Back edges take us backward, from a descendant node to one of its ancestors.
//
// low[v]: as the minimum timestamp that vertex v can reach,that is, the minimum timestamp that the subtrees of v and v can reach,
// and also describe it as the minimum timestamp that v can trace in the dfs stack.
// If the low of a vertex v is equal to its timestamp, then that vertex must be the "root" of its strongly connected component.
func tarjan[K comparable, W number](g Digraph[K, any, W], u K, stack *stack[K], num *int, dfn, low map[K]int, scc map[K][]K) error {
	*num++
	dfn[u] = *num
	low[u] = *num
	stack.push(u)
	//
	es, err := g.OutEdges(u)
	if err != nil {
		return err
	}
	for _, e := range es {
		v := e.Tail
		// v has not been visited.
		if dfn[v] == 0 {
			if err = tarjan(g, v, stack, num, dfn, low, scc); err != nil {
				return err
			}
			//low[u] = min(low[u],low[v])
			if low[v] < low[u] {
				low[u] = low[v]
			}
		} else if stack.contains(v) {
			// low[u] = min(low[u],dfs[v])
			if low[u] < dfn[v] {
				low[u] = dfn[v]
			}
		}
	}
	//
	if dfn[u] == low[u] {
		for {
			v, ok := stack.pop()
			if !ok {
				break
			}
			if _, ok := scc[u]; !ok {
				scc[u] = []K{v}
			} else {
				scc[u] = append(scc[u], v)
			}
			if u == v {
				break
			}
		}
	}
	return nil
}

func sccTarjan[K comparable, W number](g Digraph[K, any, W]) ([][]K, error) {
	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}
	if len(vertexes) == 0 {
		return [][]K{}, nil
	}

	stack := newStack[K]()
	//
	num := 0
	dfn := make(map[K]int)
	low := make(map[K]int)
	// record scc vertexes with root k
	scc := make(map[K][]K)
	//
	for _, v := range vertexes {
		if dfn[v.Key] == 0 {
			if err = tarjan(g, v.Key, stack, &num, dfn, low, scc); err != nil {
				return nil, err
			}
		}
	}

	var sccs [][]K
	for _, vs := range scc {
		sccs = append(sccs, vs)
	}

	return sccs, nil
}

// Calculate the strongly connected components of a directed graph and
// return the set of vertices for each strongly connected component.
func StronglyConnectedComponent[K comparable, W number](g Digraph[K, any, W]) ([][]K, error) {
	return sccTarjan(g)
}
