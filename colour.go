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

func vertexColouringFrom[K comparable, V any, W number](g Graph[K, V, W], n int, v Vertex[K, V], colours map[K]int) (bool, error) {
	safe := func(c int, vs []Vertex[K, V]) bool {
		for _, v := range vs {
			if colours[v.Key] == c {
				return false
			}
		}
		return true
	}

	if _, ok := colours[v.Key]; ok {
		return true, nil
	}
	vs, err := g.Neighbours(v.Key)
	if err != nil {
		return false, err
	}
	//
	for c := 1; c <= n; c++ {
		if safe(c, vs) {
			colours[v.Key] = c
			for _, k := range vs {
				ok, err := vertexColouringFrom(g, n, k, colours)
				if err != nil {
					return false, err
				}
				if ok {
					return true, nil
				}
			}
			delete(colours, v.Key)
		}
	}
	return false, nil
}

// backtracking
func vertexColouring[K comparable, V any, W number](g Graph[K, V, W], n int) (map[K]int, error) {
	p, err := g.Property(PropertyMaxDegree)
	if err != nil {
		return nil, err
	}
	if n < p.Value.(int) {
		return nil, errNoColouring
	}

	vertexes, err := g.AllVertexes()
	if err != nil {
		return nil, err
	}
	//
	colouring := make(map[K]int)
	/*
		if _,err = vertexColouringFrom(g,n,vertexes[0],colouring);err!=nil{
			return nil,err
		}
		return colouring,nil
	*/
	//
	safe := func(c int, vs []Vertex[K, V]) bool {
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
func VertexColouring[K comparable, V any, W number](g Graph[K, V, W], colours int) (map[K]int, error) {
	return vertexColouring(g, colours)
}

func edgeColouring[K comparable, V any, W number](g Graph[K, V, W], n int) (map[K]int, error) {
	p, err := g.Property(PropertyMaxDegree)
	if err != nil {
		return nil, err
	}
	if n < p.Value.(int)+1 {
		return nil, errNoColouring
	}

	edges, err := g.AllEdges()
	if err != nil {
		return nil, err
	}
	//
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
func EdgeColouring[K comparable, V any, W number](g Graph[K, V, W], colours int) (map[K]int, error) {
	return edgeColouring(g, colours)
}
