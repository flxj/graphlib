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

//
func IsCutvertex[K comparable, V any, W number](g Graph[K, V, W], vertex K) (bool, error) {
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

/*
Pick an arbitrary vertex of the graph root  and run depth first search from it.
Let's say we are in the DFS, looking through the edges starting from vertex v .
The current edge (v, w)  is a bridge if and only if none of the vertices w and its descendants in the DFS traversal tree has a back-edge to
vertex v or any of its ancestors. Indeed, this condition means that there is no other way from v  to w  except for edge (v, w).
*/
func IsBridge[K comparable, V any, W number](g Graph[K, V, W], edge K) (bool, error) {
	e, err := g.GetEdgeByKey(edge)
	if err != nil {
		return false, err
	}
	return isBridge(g, e)
}

func FindBridges[K comparable, V any, W number](g Graph[K, V, W]) ([]Edge[K, W], error) {
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
	g Graph[K, V, W]
}
