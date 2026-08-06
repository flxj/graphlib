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

func isConnected[K comparable, V any, W number](g Graph[K, V, W], vertex K, edges map[K]Edge[K, W]) (bool, error) {
	var nilK K
	var v K
	for {
		rv, err := g.RandomVertex()
		if err != nil {
			return false, err
		}
		if rv.Key != vertex {
			v = rv.Key
			break
		}
	}
	//
	order := g.Order()
	if vertex != nilK {
		order--
	}
	visited := make(map[K]bool)
	que := newFIFO[K]()
	que.push(v)

	for !que.empty() {
		p, _ := que.pop()
		if _, ok := visited[p]; !ok {
			visited[p] = true
		} else {
			continue
		}
		//
		es, err := g.IncidentEdges(p)
		if err != nil {
			return false, err
		}

		vs := make(map[K]bool)
		for _, e := range es {
			if e.Tail != vertex && e.Head != vertex {
				_, ok := edges[e.Key]
				if !ok {
					vs[e.Tail] = true
					vs[e.Head] = true
				}
			}
		}
		for k := range vs {
			if k != p {
				if _, ok := visited[k]; !ok {
					que.push(v)
				}
			}
		}
	}

	return len(visited) == order, nil
}

func IsCutvertex[K comparable, V any, W number](g Graph[K, V, W], vertex K) (bool, error) {
	if g == nil {
		return false, errNilGraph
	}
	vs, err := g.Neighbours(vertex)
	if err != nil {
		return false, err
	}
	if len(vs) == 0 || len(vs) == 1 {
		return false, nil
	}

	es := make(map[K]Edge[K, W])
	for _, v := range vs {
		ee, err := g.GetEdge(vertex, v.Key)
		if err != nil {
			return false, err
		}
		for _, e := range ee {
			es[e.Key] = e
		}
	}

	ok, err := isConnected(g, vertex, es)
	if err != nil {
		return false, err
	}
	return !ok, nil
}

/*
Pick an arbitrary vertex of the graph root  and run depth first search from it.
Let's say we are in the DFS, looking through the edges starting from vertex v .
The current edge (v, w)  is a bridge if and only if none of the vertices w and its descendants in the DFS traversal tree has a back-edge to
vertex v or any of its ancestors. Indeed, this condition means that there is no other way from v  to w  except for edge (v, w).
*/
func isBridge[K comparable, V any, W number](g Graph[K, V, W], edge Edge[K, W]) (bool, error) {
	var bri bool
	var t int
	visited := make(map[K]struct{})
	// let  inTime[v] denote entry time for node v.
	// We introduce an array lowTime[v] which will let us store the earliest entry time of the node found in the DFS search
	// that a node v can reach with a single edge from itself or its descendants.
	inTime := make(map[K]int)
	lowTime := make(map[K]int)

	var dfs func(K, K) error
	dfs = func(v K, p K) error {
		visited[v] = struct{}{}
		t++
		inTime[v] = t
		lowTime[v] = t
		ns, err := g.Neighbours(v)
		if err != nil {
			return err
		}
		var skip bool
		for _, w := range ns {
			if w.Key == p && !skip { // exist paralell edge (w,v)
				skip = true
				continue
			}
			if _, ok := visited[w.Key]; ok {
				// if w has already visited: which means (v,w) is a back edge,and w is a ancestors of v.
				lowTime[v] = min(lowTime[v], inTime[w.Key])
			} else {
				dfs(w.Key, v)
				lowTime[v] = min(lowTime[v], lowTime[w.Key])
				if lowTime[w.Key] > inTime[v] && v == edge.Head && w.Key == edge.Tail { // edge = (v,w) is a bridge.
					bri = true
					return nil
				}
			}
		}
		return nil
	}
	if err := dfs(edge.Head, edge.Head); err != nil {
		return false, err
	}
	return bri, nil
}

func IsBridge[K comparable, V any, W number](g Graph[K, V, W], edge K) (bool, error) {
	if g == nil {
		return false, errNilGraph
	}
	e, err := g.GetEdgeByKey(edge)
	if err != nil {
		return false, err
	}
	return isBridge(g, e)
}

func FindBridges[K comparable, V any, W number](g Graph[K, V, W]) ([]Edge[K, W], error) {
	if g == nil {
		return nil, errNilGraph
	}
	var bri []Edge[K, W]
	var t int
	visited := make(map[K]struct{})
	// let  inTime[v] denote entry time for node v.
	// We introduce an array lowTime[v] which will let us store the earliest entry time of the node found in the DFS search
	// that a node v can reach with a single edge from itself or its descendants.
	inTime := make(map[K]int)
	lowTime := make(map[K]int)

	var dfs func(K, K) error
	dfs = func(v K, p K) error {
		visited[v] = struct{}{}
		t++
		inTime[v] = t
		lowTime[v] = t
		es, err := g.IncidentEdges(v)
		if err != nil {
			return err
		}
		var skip bool
		for _, e := range es {
			var w K
			if e.Head == v {
				w = e.Tail
			} else {
				w = e.Head
			}
			if w == p && !skip { // exist paralell edge (w,v)
				skip = true
				continue
			}
			if _, ok := visited[w]; ok {
				// if w has already visited: which means (v,w) is a back edge,and w is a ancestors of v.
				lowTime[v] = min(lowTime[v], inTime[w])
			} else {
				dfs(w, v)
				lowTime[v] = min(lowTime[v], lowTime[w])
				if lowTime[w] > inTime[v] { // (v,w) is a bridge.
					bri = append(bri, e)
				}
			}
		}
		return nil
	}
	vtx, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}
	for _, v := range vtx {
		if _, ok := visited[v.Key]; !ok {
			if err := dfs(v.Key, v.Key); err != nil {
				return nil, err
			}
		}
	}
	return bri, nil
}

type FindBridgesOnline[K comparable, V any, W number] struct {
	g   Graph[K, V, W]
	vtx []K
	idx map[K]int
	bri map[int]int

	lcaIter int
	// connected components
	cc *dynamicUnionFind
	// 2-edge-connected components
	two_ecc *dynamicUnionFind
	// It is very easy to see, that the bridges partition the graph into 2-edge-connected components.
	// If we compress each of those 2-edge-connected components into vertices
	// and only leave the bridges as edges in the compressed graph, then we obtain an acyclic graph, i.e. a forest
	parent []int
	visit  []int
}

func NewFindBridgesOnline[K comparable, V any, W number](g Graph[K, V, W]) (*FindBridgesOnline[K, V, W], error) {
	if g == nil {
		return nil, errNilGraph
	}
	vs, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}
	es, err := g.AllEdges()
	if err != nil {
		return nil, err
	}
	fb := &FindBridgesOnline[K, V, W]{
		g:       g,
		vtx:     make([]K, g.Order()),
		idx:     make(map[K]int),
		bri:     make(map[int]int),
		cc:      newDynamicUnionFind(g.Order()),
		two_ecc: newDynamicUnionFind(g.Order()),
		parent:  make([]int, g.Order()),
		visit:   make([]int, g.Order()),
	}
	for i := range fb.parent {
		fb.parent[i] = -1
	}
	for i, v := range vs {
		fb.vtx[i] = v.Key
		fb.idx[v.Key] = i
	}
	for _, e := range es {
		if err := fb.AddEdge(e); err != nil {
			return nil, err
		}
	}
	return fb, nil
}

func (f *FindBridgesOnline[K, V, W]) AddVertex(v Vertex[K, V]) error {
	if err := f.g.AddVertex(v); err != nil {
		return err
	}
	f.vtx = append(f.vtx, v.Key)
	f.idx[v.Key] = len(f.vtx) - 1
	f.parent = append(f.parent, -1)
	f.visit = append(f.visit, 0)
	f.cc.Add(1)
	f.two_ecc.Add(1)
	return nil
}

func (f *FindBridgesOnline[K, V, W]) AddEdge(e Edge[K, W]) error {
	if err := f.g.AddEdge(e); err != nil {
		return err
	}
	a, b := f.idx[e.Head], f.idx[e.Tail]
	// Both vertices a and b are in the same 2-edge-connected component
	// - then this edge is not a bridge, and does not change anything in the forest structure, so we can just skip this edge.
	a, b = f.two_ecc.Find(a), f.two_ecc.Find(b)
	if a == b {
		return nil
	}
	if f.cc.Find(a) != f.cc.Find(b) {
		// The vertices a and b are in completely different connected components,
		// i.e. each one is part of a different tree.
		// In this case, the edge (a, b) becomes a new bridge, and these two trees are combined into one (and all the old bridges remain).
		f.addBridge(f.idx[e.Head], f.idx[e.Tail])
		f.makeRoot(a)
		f.parent[a] = b
		f.cc.SetParent(a, b)
	} else {
		// The vertices a and b are in one connected component, but in different 2-edge-connected components.
		// In this case, this edge forms a cycle along with some of the old bridges.
		// All these bridges end being bridges, and the resulting cycle must be compressed into a new 2-edge-connected component.
		f.connect(a, b)
	}
	return nil
}

func (f *FindBridgesOnline[K, V, W]) makeRoot(v int) {
	root, child := v, -1
	for v != -1 {
		p := f.two_ecc.Find(f.parent[v])
		f.parent[v] = child
		f.cc.SetParent(v, root)
		child, v = v, p
	}
}

func (f *FindBridgesOnline[K, V, W]) addBridge(a, b int) {
	c, ok := f.bri[a]
	if ok {
		if c == b {
			return
		}
		f.bri[b] = a
	} else {
		f.bri[a] = b
	}
}

func (f *FindBridgesOnline[K, V, W]) deleteBridge(a, b int) {
	c, ok := f.bri[a]
	if ok && c == b {
		delete(f.bri, a)
	}
	c, ok = f.bri[b]
	if ok && c == a {
		delete(f.bri, b)
	}
}

func (f *FindBridgesOnline[K, V, W]) connect(a, b int) {
	// The vertices a and b are in one connected component, but in different 2-edge-connected components.
	// In this case,this edge forms a cycle along with some of the old bridges.
	// All these bridges end being bridges, and the resulting cycle must be compressed into a new 2-edge-connected component.
	f.lcaIter++
	// Searching for the cycle formed by adding a new edge (a, b) .
	// Since a and b are already connected in the tree we need to find the Lowest Common Ancestor of the vertices a and b.
	// The cycle will consist of the paths from  $b$  to the LCA, from the LCA to  $a$  and the edge a to b.
	var pathA, pathB []int
	lca := -1
	for lca == -1 {
		if a != -1 {
			a = f.two_ecc.Find(a)
			pathA = append(pathA, a)
			if f.visit[a] == f.lcaIter {
				lca = a
				break
			}
			f.visit[a] = f.lcaIter
			a = f.parent[a]
		}
		if b != -1 {
			b = f.two_ecc.Find(b)
			pathB = append(pathB, b)
			if f.visit[b] == f.lcaIter {
				lca = b
				break
			}
			f.visit[b] = f.lcaIter
			b = f.parent[b]
		}
	}
	for i, v := range pathA {
		f.two_ecc.SetParent(v, lca)
		if v == lca {
			break
		}
		if i > 0 {
			// edge (v,pathA[i-1]) is a bridge
			f.deleteBridge(v, pathA[i-1])
		}
	}
	for i, v := range pathB {
		f.two_ecc.SetParent(v, lca)
		if v == lca {
			break
		}
		if i > 0 {
			f.deleteBridge(v, pathB[i-1])
		}
	}
}

func (f *FindBridgesOnline[K, V, W]) Graph() Graph[K, V, W] {
	return f.g
}

func (f *FindBridgesOnline[K, V, W]) IsBridge(e Edge[K, W]) bool {
	i, j := f.idx[e.Head], f.idx[e.Tail]
	if v, ok := f.bri[i]; ok {
		return v == j
	}
	if v, ok := f.bri[j]; ok {
		return v == i
	}
	return false
}

func (f *FindBridgesOnline[K, V, W]) Bridges() ([]Edge[K, W], error) {
	res := make([]Edge[K, W], len(f.bri))
	var i int
	for a, b := range f.bri {
		e, err := f.g.GetEdge(f.vtx[a], f.vtx[b])
		if err != nil {
			return nil, err
		}
		res[i] = e[0]
		i++
	}
	return res, nil
}

func (f *FindBridgesOnline[K, V, W]) BridgeCount() int {
	return len(f.bri)
}

func (f *FindBridgesOnline[K, V, W]) Components() int {
	return f.cc.Component()
}
