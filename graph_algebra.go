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
func Union[K comparable, V any, W number](g1, g2 Graph[K, V, W]) (Graph[K, V, W], error) {
	if g1.IsDigraph() != g2.IsDigraph() {
		return nil, errNotSameType
	}
	var (
		err      error
		vs1, vs2 []Vertex[K, V]
		es1, es2 []Edge[K, W]
	)
	if vs1, err = g1.AllVertexes(); err != nil {
		return nil, err
	}
	if es1, err = g1.AllEdges(); err != nil {
		return nil, err
	}
	if vs2, err = g2.AllVertexes(); err != nil {
		return nil, err
	}
	if es2, err = g2.AllEdges(); err != nil {
		return nil, err
	}

	uv := make(map[K]*Vertex[K, V])
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

	ug, err := NewGraph[K, V, W](g1.IsDigraph(), fmt.Sprintf("%s-union-%s", g1.Name(), g2.Name()))
	if err != nil {
		return nil, err
	}

	for _, v := range uv {
		if err = ug.AddVertex(*v); err != nil {
			return nil, err
		}
	}
	for _, e := range ue {
		if err = ug.AddEdge(*e); err != nil {
			return nil, err
		}
	}

	return ug, nil
}

// Calculate the union of two graphs.
func Intersection[K comparable, V any, W number](g1, g2 Graph[K, V, W]) (Graph[K, V, W], error) {
	if g1.IsDigraph() != g2.IsDigraph() {
		return nil, errNotSameType
	}
	var (
		err      error
		vs1, vs2 []Vertex[K, V]
		es1, es2 []Edge[K, W]
	)
	if vs1, err = g1.AllVertexes(); err != nil {
		return nil, err
	}
	if es1, err = g1.AllEdges(); err != nil {
		return nil, err
	}
	if vs2, err = g2.AllVertexes(); err != nil {
		return nil, err
	}
	if es2, err = g2.AllEdges(); err != nil {
		return nil, err
	}

	iv := make(map[K]bool)
	ie := make(map[K]bool)

	uv := make(map[K]*Vertex[K, V])
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
	ug, err := NewGraph[K, V, W](g1.IsDigraph(), fmt.Sprintf("%s-intersection-%s", g1.Name(), g2.Name()))
	if err != nil {
		return nil, err
	}

	for _, v := range uv {
		if err = ug.AddVertex(*v); err != nil {
			return nil, err
		}
	}
	for _, e := range ue {
		if err = ug.AddEdge(*e); err != nil {
			return nil, err
		}
	}

	return ug, nil
}

func CartesianProduct[K comparable, V any, W number](g1, g2 Graph[K, V, W]) (Graph[string, V, W], error) {
	if g1.IsDigraph() != g2.IsDigraph() {
		return nil, errors.New("not support operation")
	}

	var (
		err  error
		g1vs []Vertex[K, V]
		g2vs []Vertex[K, V]
		g1es []Edge[K, W]
		g2es []Edge[K, W]
	)
	if g1vs, err = g1.AllVertexes(); err != nil {
		return nil, err
	}
	if g2vs, err = g2.AllVertexes(); err != nil {
		return nil, err
	}
	if g1es, err = g1.AllEdges(); err != nil {
		return nil, err
	}
	if g2es, err = g2.AllEdges(); err != nil {
		return nil, err
	}

	g, _ := NewGraph[string, V, W](g1.IsDigraph(), g1.Name()+"X"+g2.Name())

	for _, v1 := range g1vs {
		for _, v2 := range g2vs {
			v := Vertex[string, V]{
				Key: fmt.Sprintf("(%v,%v)", v1.Key, v2.Key),
				Labels: map[string]string{
					g1.Name(): fmt.Sprintf("%v", v1.Key),
					g2.Name(): fmt.Sprintf("%v", v2.Key),
				},
			}
			if err = g.AddVertex(v); err != nil {
				return nil, err
			}
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
			if err = g.AddEdge(e); err != nil {
				return nil, err
			}
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
			if err = g.AddEdge(e); err != nil {
				return nil, err
			}
		}
	}

	return g, nil
}

func Identify[K comparable, V any, W number](g Graph[K, V, W], v1, v2 K, newVertex Vertex[K, V], createGraph bool) (Graph[K, V, W], error) {
	return Contract(g, v1, v2, newVertex, createGraph)
}

func Contract[K comparable, V any, W number](g Graph[K, V, W], v1, v2 K, newVertex Vertex[K, V], createGraph bool) (Graph[K, V, W], error) {
	var err error
	if _, err = g.GetVertex(v1); err != nil {
		return nil, err
	}
	if _, err = g.GetVertex(v2); err != nil {
		return nil, err
	}

	g2 := g
	if createGraph {
		if g2, err = g.Clone(); err != nil {
			return nil, err
		}
	}
	// add new vertex
	if err = g2.AddVertex(newVertex); err != nil {
		return nil, err
	}

	newEdges := make(map[K]Edge[K, W])
	// if find A={v1-x,x-v1 | x!=v2}, then add new edges x-new new-x
	es1, err := g.IncidentEdges(v1)
	if err != nil {
		return nil, err
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
	es2, err := g.IncidentEdges(v2)
	if err != nil {
		return nil, err
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
		if err = g2.RemoveEdgeByKey(e.Key); err != nil {
			if !IsNotExists(err) {
				return nil, err
			}
		}
	}
	for _, e := range es2 {
		if err = g2.RemoveEdgeByKey(e.Key); err != nil {
			if !IsNotExists(err) {
				return nil, err
			}
		}
	}

	// delete edge v1-v2
	if err = g2.RemoveEdge(v1, v2); err != nil {
		if !IsNotExists(err) {
			return nil, err
		}
	}

	for _, e := range newEdges {
		if err = g2.AddEdge(e); err != nil {
			return nil, err
		}
	}
	return g2, nil
}

func Split[K comparable, V any, W number](g Graph[K, V, W], vertex K, edge Edge[K, W], newEdgeKey func(Edge[K, W]) K, createGraph bool) (Graph[K, V, W], error) {
	var err error
	if _, err = g.GetVertex(edge.Head); err == nil {
		return nil, fmt.Errorf("vertex %v already exists", edge.Head)
	}
	if _, err = g.GetVertex(edge.Tail); err == nil {
		return nil, fmt.Errorf("vertex %v already exists", edge.Tail)
	}

	g2 := g
	if createGraph {
		if g2, err = g.Clone(); err != nil {
			return nil, err
		}
	}

	newEdges := make(map[K]Edge[K, W])
	es, err := g.IncidentEdges(vertex)
	if err != nil {
		return nil, err
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
	if err = g2.RemoveVertex(vertex); err != nil {
		return nil, err
	}

	for _, e := range newEdges {
		if err = g2.AddEdge(e); err != nil {
			return nil, err
		}
	}
	if err = g2.AddEdge(edge); err != nil {
		return nil, err
	}

	return g2, nil
}

func Subdivide[K comparable, V any, W number](g Graph[K, V, W], edge K, vertex Vertex[K, V], newEdgeKey func(Edge[K, W]) K, createGraph bool) (Graph[K, V, W], error) {
	_, err := g.GetVertex(vertex.Key)
	if err == nil {
		return nil, fmt.Errorf("vertex %v already exists", vertex.Key)
	}

	g2 := g
	if createGraph {
		if g2, err = g.Clone(); err != nil {
			return nil, err
		}
	}

	e, err := g.GetEdgeByKey(edge)
	if err != nil {
		return nil, err
	}
	ne := Edge[K, W]{
		Head: e.Head,
		Tail: vertex.Key,
	}
	ne.Key = newEdgeKey(ne)

	if err = g2.AddEdge(ne); err != nil {
		return nil, err
	}

	ne = Edge[K, W]{
		Head: vertex.Key,
		Tail: e.Tail,
	}
	ne.Key = newEdgeKey(ne)
	if err = g2.AddEdge(ne); err != nil {
		return nil, err
	}

	if err = g2.RemoveEdgeByKey(edge); err != nil {
		return nil, err
	}

	return g2, nil
}
