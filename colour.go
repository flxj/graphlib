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

// backtracking
func vertexColouring[K comparable, W number](g Graph[K, W], n int) (map[K]int, error) {
	p, err := g.Property(ProMaxDegree)
	if err != nil {
		return nil, err
	}
	if n < p.Value.(int) {
		return nil, errNoColouring
	}

	vertexes := g.AllVertexes()
	colouring := make(map[K]int)
	/*
		if _,err = vertexColouringFrom(g,n,vertexes[0],colouring);err!=nil{
			return nil,err
		}
		return colouring,nil
	*/
	//
	safe := func(c int, vs []Vertex[K, W]) bool {
		for _, v := range vs {
			if colouring[v.Key] == c {
				return false
			}
		}
		return true
	}

	var nilK K
	prev := make(map[K]K)
	stack := newStack[K]()

	stack.push(vertexes[0].Key)
	prev[vertexes[0].Key] = nilK

	for !stack.empty() {
		if len(colouring) == len(vertexes) {
			break
		}
		//
		v, _ := stack.pop()
		if _, ok := colouring[v]; ok {
			continue
		}
		vs, err := g.Neighbours(v)
		if err != nil {
			return nil, err
		}
		//
		var col int
		for c := colouring[v]; c < n; c++ {
			if safe(c+1, vs) {
				col = c + 1
				break
			}
		}
		if col == 0 {
			// backtrack
			delete(colouring, v)
			p := prev[v]
			if p == nilK {
				return nil, errNoColouring
			}
			stack.push(p)
		} else {
			colouring[v] = col
			for _, k := range vs {
				prev[k.Key] = v
				stack.push(k.Key)
			}
		}
		// colouring another components.
		if stack.empty() && len(prev) < len(vertexes) {
			for _, v := range vertexes {
				if _, ok := prev[v.Key]; !ok {
					stack.push(v.Key)
					prev[v.Key] = nilK
					break
				}
			}
		}
	}
	return colouring, nil
}

// Graph vertex coloring, returning a feasible coloring scheme.
func TryVertexColouring[K comparable, W number](g Graph[K, W], colours int) (map[K]int, error) {
	if g == nil {
		return nil, errNilGraph
	}
	return vertexColouring(g, colours)
}

func edgeColouring[K comparable, W number](g Graph[K, W], n int) (map[K]int, error) {
	p, err := g.Property(ProMaxDegree)
	if err != nil {
		return nil, err
	}
	if n < p.Value.(int)+1 {
		return nil, errNoColouring
	}

	edges := g.AllEdges()
	colouring := make(map[K]int)

	safe := func(c int, es []Edge[K, W]) bool {
		for _, e := range es {
			if colouring[e.Key] == c {
				return false
			}
		}
		return true
	}

	var nilK K
	prev := make(map[K]K)
	stack := newStack[K]()

	stack.push(edges[0].Key)
	prev[edges[0].Key] = nilK

	for !stack.empty() {
		if len(colouring) == len(edges) {
			break
		}
		//
		e, _ := stack.pop()
		if _, ok := colouring[e]; ok {
			continue
		}
		es, err := g.NeighbourEdgesByKey(e)
		if err != nil {
			return nil, err
		}
		//
		var col int
		for c := colouring[e]; c < n; c++ {
			if safe(c+1, es) {
				col = c + 1
				break
			}
		}
		if col == 0 {
			// backtrack
			delete(colouring, e)
			p := prev[e]
			if p == nilK {
				return nil, errNoColouring
			}
			stack.push(p)
		} else {
			colouring[e] = col
			for _, k := range es {
				prev[k.Key] = e
				stack.push(k.Key)
			}
		}
		// colouring another components.
		if stack.empty() && len(prev) < len(edges) {
			for _, v := range edges {
				if _, ok := prev[v.Key]; !ok {
					stack.push(v.Key)
					prev[v.Key] = nilK
					break
				}
			}
		}
	}
	return colouring, nil
}

// Graph edge coloring, returning a feasible coloring scheme.
func TryEdgeColouring[K comparable, W number](g Graph[K, W], colours int) (map[K]int, error) {
	if g == nil {
		return nil, errNilGraph
	}
	return edgeColouring(g, colours)
}

// greedy graph vertex coloring, returning a feasible coloring scheme.
func GreedyVertexColouring[K comparable, W number](g Graph[K, W]) (map[K]int, int, error) {
	if g == nil {
		return nil, 0, errNilGraph
	}
	vtx := g.AllVertexes()
	col := make(map[K]int)
	var cnt int
	for i := 0; i < len(vtx); i++ {
		vs, err := g.Neighbours(vtx[i].Key)
		if err != nil {
			return nil, 0, err
		}
		used := make(map[int]bool)
		for _, u := range vs {
			c, ok := col[u.Key]
			if ok {
				used[c] = true
			}
		}
		for c := 0; ; c++ {
			if !used[c] {
				col[vtx[i].Key] = c
				if c > cnt {
					cnt = c
				}
				break
			}
		}
	}
	return col, cnt + 1, nil
}

// Ggreedy graph edge coloring, returning a feasible coloring scheme.
func GreedyEdgeColouring[K comparable, W number](g Graph[K, W]) (map[K]int, int, error) {
	if g == nil {
		return nil, 0, errNilGraph
	}
	edge := g.AllEdges()
	col := make(map[K]int)
	var cnt int
	for i := 0; i < len(edge); i++ {
		es, err := g.NeighbourEdges(edge[i].Head, edge[i].Tail)
		if err != nil {
			return nil, 0, err
		}
		used := make(map[int]bool)
		for _, e := range es {
			c, ok := col[e.Key]
			if ok {
				used[c] = true
			}
		}
		for c := 0; ; c++ {
			if !used[c] {
				col[edge[i].Key] = c
				if c > cnt {
					cnt = c
				}
				break
			}
		}
	}
	return col, cnt + 1, nil
}

func MaximalIndependentSet[K comparable, W number](g Graph[K, W]) ([]K, error) {
	return nil, errNotImplement
}

func MaximumIndependentSet[K comparable, W number](g Graph[K, W]) ([]K, error) {
	return nil, errNotImplement
}
