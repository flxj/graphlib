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

func CartesianProduct[K comparable, V any, W number](g1, g2 Graph[K, V, W]) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

// page 55
func Identify[K comparable, V any, W number](g Graph[K, V, W], v1, v2 K, newVertex K) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

// page 55
func Contract[K comparable, V any, W number](g Graph[K, V, W], v1, v2 K, newVertex K) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

// page 55
func Split[K comparable, V any, W number](g Graph[K, V, W], vertex K, v1, v2, edge K) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

// page 55
func Subdivide[K comparable, V any, W number](g Graph[K, V, W], edge K, newVertex K) (Graph[K, V, W], error) {
	return nil, errNotImplement
}
