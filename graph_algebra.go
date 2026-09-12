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
	"fmt"
)

// Calculate the intersection of two graphs.
func Union[K comparable, W number](g1, g2 Graph[K, W]) (Graph[K, W], error) {
	if g1.IsDigraph() != g2.IsDigraph() {
		return nil, errNotSameType
	}
	vs1 := g1.AllVertexes()
	es1 := g1.AllEdges()
	vs2 := g2.AllVertexes()
	es2 := g2.AllEdges()

	uv := make(map[K]*Vertex[K, W])
	ue := make(map[K]*Edge[K, W])
	for _, v := range vs1 {
		vv := v
		uv[vv.Key] = &vv
	}
	for _, v := range vs2 {
		vv := v
		uv[vv.Key] = &vv
	}
	for _, e := range es1 {
		ee := e
		ue[ee.Key] = &ee
	}
	for _, e := range es2 {
		ee := e
		ue[ee.Key] = &ee
	}

	ug := NewGraph[K, W](g1.IsDigraph(), fmt.Sprintf("%s-union-%s", g1.Name(), g2.Name()))

	for _, v := range uv {
		_ = ug.AddVertex(*v)
	}
	for _, e := range ue {
		_ = ug.AddEdge(*e)
	}

	return ug, nil
}

// Calculate the union of two graphs.
func Intersection[K comparable, W number](g1, g2 Graph[K, W]) (Graph[K, W], error) {
	if g1.IsDigraph() != g2.IsDigraph() {
		return nil, errNotSameType
	}

	vs1 := g1.AllVertexes()
	es1 := g1.AllEdges()
	vs2 := g2.AllVertexes()
	es2 := g2.AllEdges()

	iv := make(map[K]bool)
	ie := make(map[K]bool)

	uv := make(map[K]*Vertex[K, W])
	ue := make(map[K]*Edge[K, W])
	for _, v := range vs1 {
		iv[v.Key] = true
	}
	for _, v := range vs2 {
		if iv[v.Key] {
			vv := v
			uv[vv.Key] = &vv
		}
	}
	for _, e := range es1 {
		ie[e.Key] = true
	}
	for _, e := range es2 {
		if ie[e.Key] {
			ee := e
			ue[ee.Key] = &ee
		}
	}
	ug := NewGraph[K, W](g1.IsDigraph(), fmt.Sprintf("%s-intersection-%s", g1.Name(), g2.Name()))

	for _, v := range uv {
		_ = ug.AddVertex(*v)
	}
	for _, e := range ue {
		_ = ug.AddEdge(*e)
	}

	return ug, nil
}

func CartesianProduct[K comparable, W number](g1, g2 Graph[K, W]) (Graph[string, W], error) {
	if g1.IsDigraph() != g2.IsDigraph() {
		return nil, errors.New("not support operation")
	}

	g1vs := g1.AllVertexes()
	g2vs := g2.AllVertexes()
	g1es := g1.AllEdges()
	g2es := g2.AllEdges()

	g := NewGraph[string, W](g1.IsDigraph(), g1.Name()+"X"+g2.Name())
	for _, v1 := range g1vs {
		for _, v2 := range g2vs {
			v := Vertex[string, W]{
				Key: fmt.Sprintf("(%v,%v)", v1.Key, v2.Key),
				Labels: map[string]any{
					g1.Name(): fmt.Sprintf("%v", v1.Key),
					g2.Name(): fmt.Sprintf("%v", v2.Key),
				},
			}
			_ = g.AddVertex(v)
		}
	}
	for _, e := range g1es {
		//(v1,v) -- (v2,v)
		for _, v := range g2vs {
			head := fmt.Sprintf("(%v,%v)", e.Head, v.Key)
			tail := fmt.Sprintf("(%v,%v)", e.Tail, v.Key)
			e := Edge[string, W]{
				Key:  head + "-" + tail,
				Head: head,
				Tail: tail,
			}
			_ = g.AddEdge(e)
		}
	}
	for _, e := range g2es {
		// (v,v1) -- (v,v2)
		for _, v := range g1vs {
			head := fmt.Sprintf("(%v,%v)", v.Key, e.Head)
			tail := fmt.Sprintf("(%v,%v)", v.Key, e.Tail)
			e := Edge[string, W]{
				Key:  head + "-" + tail,
				Head: head,
				Tail: tail,
			}
			_ = g.AddEdge(e)
		}
	}

	return g, nil
}

func Identify[K comparable, W number](g Graph[K, W], v1, v2 K, newVertex Vertex[K, W], createGraph bool) (Graph[K, W], error) {
	return Contract(g, v1, v2, newVertex, createGraph)
}

func Contract[K comparable, W number](g Graph[K, W], v1, v2 K, newVertex Vertex[K, W], createGraph bool) (Graph[K, W], error) {
	if _, ok := g.GetVertex(v1); !ok {
		return nil, errVertexNotExists
	}
	if _, ok := g.GetVertex(v2); !ok {
		return nil, errVertexNotExists
	}

	g2 := g
	if createGraph {
		g2 = g.Clone()
	}
	// add new vertex
	_ = g2.AddVertex(newVertex)

	newEdges := make(map[K]Edge[K, W])
	// if find A={v1-x,x-v1 | x!=v2}, then add new edges x-new new-x
	es1, ok := g.IncidentEdges(v1)
	if !ok {
		return nil, errVertexNotExists
	}
	for _, e := range es1 {
		if v1 == e.Head {
			if e.Tail != v2 {
				ne := Edge[K, W]{
					Key:    e.Key,
					Head:   newVertex.Key,
					Tail:   e.Tail,
					Weight: e.Weight,
				}
				newEdges[ne.Key] = ne
			}
		} else {
			if e.Head != v2 {
				ne := Edge[K, W]{
					Key:    e.Key,
					Tail:   newVertex.Key,
					Head:   e.Head,
					Weight: e.Weight,
				}
				newEdges[ne.Key] = ne
			}
		}
	}
	// if find B = {v2-x,x-v2 | x!=v1}, then add new edges x-new new-x
	es2, ok := g.IncidentEdges(v2)
	if !ok {
		return nil, errVertexNotExists
	}
	for _, e := range es2 {
		if v1 == e.Head {
			if e.Tail != v1 {
				ne := Edge[K, W]{
					Key:    e.Key,
					Head:   newVertex.Key,
					Tail:   e.Tail,
					Weight: e.Weight,
				}
				newEdges[ne.Key] = ne
			}
		} else {
			if e.Head != v1 {
				ne := Edge[K, W]{
					Key:    e.Key,
					Tail:   newVertex.Key,
					Head:   e.Head,
					Weight: e.Weight,
				}
				newEdges[ne.Key] = ne
			}
		}
	}

	// delete A,B
	for _, e := range es1 {
		_, _ = g2.RemoveEdgeByKey(e.Key)
	}
	for _, e := range es2 {
		_, _ = g2.RemoveEdgeByKey(e.Key)
	}

	// delete edge v1-v2
	_, _ = g2.RemoveEdge(v1, v2)

	for _, e := range newEdges {
		_ = g2.AddEdge(e)
	}
	return g2, nil
}

func Split[K comparable, W number](g Graph[K, W], vertex K, edge Edge[K, W], newEdgeKey func(Edge[K, W]) K, createGraph bool) (Graph[K, W], error) {
	if _, ok := g.GetVertex(edge.Head); ok {
		return nil, fmt.Errorf("vertex %v already exists", edge.Head)
	}
	if _, ok := g.GetVertex(edge.Tail); ok {
		return nil, fmt.Errorf("vertex %v already exists", edge.Tail)
	}

	g2 := g
	if createGraph {
		g2 = g.Clone()
	}

	newEdges := make(map[K]Edge[K, W])
	es, ok := g.IncidentEdges(vertex)
	if !ok {
		return nil, errVertexNotExists
	}
	for _, e := range es {
		if e.Head == vertex {
			ne := Edge[K, W]{
				Head:   edge.Head,
				Tail:   e.Tail,
				Weight: e.Weight,
				Labels: e.Labels,
			}
			ne.Key = newEdgeKey(ne)
			newEdges[ne.Key] = ne

			ne = Edge[K, W]{
				Head:   edge.Tail,
				Tail:   e.Tail,
				Weight: e.Weight,
				Labels: e.Labels,
			}
			ne.Key = newEdgeKey(ne)
			newEdges[ne.Key] = ne
		} else {
			ne := Edge[K, W]{
				Head:   e.Head,
				Tail:   edge.Head,
				Weight: e.Weight,
				Labels: e.Labels,
			}
			ne.Key = newEdgeKey(ne)
			newEdges[ne.Key] = ne

			ne = Edge[K, W]{
				Head:   e.Head,
				Tail:   edge.Tail,
				Weight: e.Weight,
				Labels: e.Labels,
			}
			ne.Key = newEdgeKey(ne)
			newEdges[ne.Key] = ne
		}
	}
	//
	_, _ = g2.RemoveVertex(vertex)

	for _, e := range newEdges {
		_ = g2.AddEdge(e)
	}
	_ = g2.AddEdge(edge)

	return g2, nil
}

func Subdivide[K comparable, W number](g Graph[K, W], edge K, vertex Vertex[K, W], newEdgeKey func(Edge[K, W]) K, createGraph bool) (Graph[K, W], error) {
	_, ok := g.GetVertex(vertex.Key)
	if ok {
		return nil, fmt.Errorf("vertex %v already exists", vertex.Key)
	}

	g2 := g
	if createGraph {
		g2 = g.Clone()
	}

	e, ok := g.GetEdgeByKey(edge)
	if !ok {
		return nil, errEdgeNotExists
	}
	ne := Edge[K, W]{
		Head: e.Head,
		Tail: vertex.Key,
	}
	ne.Key = newEdgeKey(ne)

	_ = g2.AddEdge(ne)

	ne = Edge[K, W]{
		Head: vertex.Key,
		Tail: e.Tail,
	}
	ne.Key = newEdgeKey(ne)
	_ = g2.AddEdge(ne)

	_, _ = g2.RemoveEdgeByKey(edge)

	return g2, nil
}

func Complement[K comparable, W number](g Graph[K, W]) (Graph[K, W], error) {
	if g == nil {
		return nil, errNilGraph
	}
	ng := newGraph[K, W](g.IsDigraph(), g.Name()+"_complement")
	vtx := g.AllVertexes()
	for _, v := range vtx {
		_ = ng.AddVertex(v)
	}
	ek := make(map[K]struct{})
	genKey := func(v1, v2 K) K {
		for {
			nk, _ := randEdgeKey(v1, v2)
			if _, ok := ek[nk]; !ok {
				ek[nk] = struct{}{}
				return nk
			}
		}
	}

	for i, v := range vtx {
		for j := i + 1; j < len(vtx); j++ {
			// u := vtx[j], try to add edge (v,u) to new graph.
			es, _ := g.GetEdge(v.Key, vtx[j].Key)
			if g.IsDigraph() {
				switch len(es) {
				case 0:
					ok := ng.AddEdge(Edge[K, W]{
						Key:  genKey(v.Key, vtx[j].Key),
						Head: v.Key,
						Tail: vtx[j].Key,
					})
					if !ok {
						return nil, errEdgeExists
					}
					ok = ng.AddEdge(Edge[K, W]{
						Key:  genKey(vtx[j].Key, v.Key),
						Head: vtx[j].Key,
						Tail: v.Key,
					})
					if !ok {
						return nil, errEdgeExists
					}
				case 1:
					ok := ng.AddEdge(Edge[K, W]{
						Key:  genKey(v.Key, vtx[j].Key),
						Head: es[0].Tail,
						Tail: es[0].Head,
					})
					if !ok {
						return nil, errEdgeExists
					}
				default:
				}
			} else {
				if len(es) == 0 {
					ok := ng.AddEdge(Edge[K, W]{
						Key:  genKey(v.Key, vtx[j].Key),
						Head: v.Key,
						Tail: vtx[j].Key,
					})
					if !ok {
						return nil, errEdgeExists
					}
				}
			}
		}
	}
	return ng, nil
}
