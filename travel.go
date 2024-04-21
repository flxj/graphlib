package graphlib

func dfs[K comparable, W number](g Graph[K, any, W], start K, in bool, visitor func(v Vertex[K, any]) error) error {
	neighbours := g.Neighbours
	if g.IsDigraph() {
		dg, ok := g.(Digraph[K, any, W])
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
	stack := []Vertex[K, any]{startV}
	top := 1

	for top > 0 {
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
func DFS[K comparable, W number](g Graph[K, any, W], start K, visitor func(v Vertex[K, any]) error) error {
	return dfs(g, start, false, visitor)
}

// if in is true,
func DFSDigraph[K comparable, W number](g Digraph[K, any, W], start K, in bool, visitor func(v Vertex[K, any]) error) error {
	return dfs(g, start, in, visitor)
}

func bfs[K comparable, W number](g Graph[K, any, W], start K, in bool, visitor func(v Vertex[K, any]) error) error {
	neighbours := g.Neighbours
	if g.IsDigraph() {
		dg, ok := g.(Digraph[K, any, W])
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
	queue := []Vertex[K, any]{startV}
	head := 0
	tail := 1

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

//	 BFS
//	 Input: A digraph D =(V,A) and a vertex s ∈ V.
//	 Output: dist(s,v) and pred(v) for all v ∈ V.
//		1. For each v ∈ V set dist(s,v):=∞ and pred(v):=nil.
//		2. Set dist(s,s) := 0. Create a queue Q consisting of s.
//		3. While Q is not empty do the following. Delete a vertex u, the head of Q,
//		from Q and consider the out-neighbours of u in D one by one. If, for an
//		out-neighbour v of u,dist(s,v)=∞,thensetdist(s,v):=dist(s,u)+1,
//		pred(v):=u, and put v to the end of Q.
func BFS[K comparable, W number](g Graph[K, any, W], start K, visitor func(v Vertex[K, any]) error) error {
	return bfs(g, start, false, visitor)
}

func BFSDigraph[K comparable, W number](g Digraph[K, any, W], start K, in bool, visitor func(v Vertex[K, any]) error) error {
	return bfs(g, start, in, visitor)
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
func TopologicalSort[K comparable, W number](g Digraph[K, any, W]) ([]Vertex[K, W], error) {
	return nil, errNotImplement
}

/*
Let G =(V,E) be a digraph and s a root of G.

Procedure DFSM(G,s;nr,Nr,p)

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

Procedure STRONGCOMP(G,s;k)

	(1) DFSM(G,s;nr,Nr,p); k←0;
	(2) let H be the digraph with the opposite orientation of G;
	(3) repeat
	(4)     choose the vertex r in H for which Nr(r) is maximal;
	(5)     k←k+1;DFS(H,r;nr,p); Ck ←{v ∈ H :nr(v)=0};
	(6)     remove all vertices in Ck and all the edges incident with them;
	(7) until the vertex set of H is empty
*/
func StrongComponents[K comparable, W number](g Digraph[K, any, W]) ([]Digraph[K, any, W], error) {
	return nil, errNotImplement
}
