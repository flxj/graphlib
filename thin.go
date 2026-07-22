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

import "errors"

type ThinGraph[K comparable] struct {
	Graph[K, any, int]
}

func NewThinGraph[K comparable](digraph bool) *ThinGraph[K] {
	g, _ := newGraph[K, any, int](digraph, "")
	return &ThinGraph[K]{
		g,
	}
}

type ThinDigraph[K comparable] struct {
	Digraph[K, any, int]
}

func NewThinDigraph[K comparable]() *ThinDigraph[K] {
	g, _ := NewDigraph[K, any, int]("")
	return &ThinDigraph[K]{
		g,
	}
}

type ThinTree[K comparable] struct {
	Graph[K, any, int]
	root K
	vtx  []K
	idx  map[K]int
	duf  *dynamicUnionFind
}

func NewThinTree[K comparable]() *ThinTree[K] {
	g, _ := newGraph[K, any, int](false, "")
	t := &ThinTree[K]{Graph: g}
	t.idx = make(map[K]int)
	t.duf = newDynamicUnionFind(0)
	return t
}

func (t *ThinTree[K]) AddVertex(v Vertex[K, any]) error {
	err := t.Graph.AddVertex(v)
	if err != nil {
		if err == errVertexExists {
			return nil
		}
		return err
	}
	t.vtx = append(t.vtx, v.Key)
	t.idx[v.Key] = len(t.vtx) - 1
	t.duf.Add(1)
	return nil
}

func (t *ThinTree[K]) AddEdge(e Edge[K, int]) error {
	var ok bool
	var u, v int
	if u, ok = t.idx[e.Head]; !ok {
		return errVertexNotExists
	}
	if v, ok = t.idx[e.Tail]; !ok {
		return errVertexNotExists
	}
	if t.duf.Find(u) == t.duf.Find(v) {
		return errExistsCycle
	}
	if err := t.Graph.AddEdge(e); err != nil {
		return err
	}
	t.duf.Union(u, v)
	return nil
}

func (t *ThinTree[K]) RemoveVertex(k K) error {
	if err := t.Graph.RemoveVertex(k); err != nil {
		return err
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveVertexs(keys ...K) error {
	for _, k := range keys {
		if err := t.Graph.RemoveVertex(k); err != nil {
			return err
		}
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveEdge(endpoint1, endpoint2 K) error {
	if err := t.Graph.RemoveEdge(endpoint1, endpoint2); err != nil {
		return err
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveEdges(endpoint1, endpoint2 []K) error {
	if len(endpoint1) != len(endpoint2) {
		return errors.New("")
	}
	for i := 0; i < len(endpoint1); i++ {
		if err := t.Graph.RemoveEdge(endpoint1[i], endpoint2[i]); err != nil {
			return err
		}
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveEdgeByKey(k K) error {
	if err := t.Graph.RemoveEdgeByKey(k); err != nil {
		return err
	}
	return t.rebuild()
}

func (t *ThinTree[K]) RemoveEdgeByKeys(keys ...K) error {
	for _, k := range keys {
		if err := t.Graph.RemoveEdgeByKey(k); err != nil {
			return err
		}
	}
	return t.rebuild()
}

func (t *ThinTree[K]) rebuild() error {
	vs, err := t.AllVertexes()
	if err != nil {
		return err
	}
	es, err := t.AllEdges()
	if err != nil {
		return err
	}
	t.vtx = make([]K, len(vs))
	for i, v := range vs {
		t.vtx[i] = v.Key
		t.idx[v.Key] = i
	}
	t.duf = newDynamicUnionFind(len(t.vtx))
	for _, e := range es {
		t.duf.Union(t.idx[e.Head], t.idx[e.Tail])
	}
	return nil
}

func (t *ThinTree[K]) SetRoot(k K) error {
	if _, ok := t.idx[k]; !ok {
		return errVertexNotExists
	}
	t.root = k
	return nil
}

func (t *ThinTree[K]) GetRoot() (K, bool) {
	_, ok := t.idx[t.root]
	return t.root, ok
}

// Tarjan
func (t *ThinTree[K]) LeastCommonAncestor(k1, k2 K) (k K, b bool) {
	if k1 == k2 || t.root == k1 {
		return k1, true
	} else if t.root == k2 {
		return k2, true
	}

	if len(t.vtx) < 2 {
		return
	}
	var ok bool
	var v1, v2 int
	if v1, ok = t.idx[k1]; !ok {
		return
	}
	if v2, ok = t.idx[k2]; !ok {
		return
	}

	vst := make([]bool, len(t.vtx))
	duf := newDynamicUnionFind(len(t.vtx))

	var dfs func(u int)
	dfs = func(u int) {
		if vst[u] {
			return
		}
		vst[u] = true // visited u
		ns, err := t.Neighbours(t.vtx[u])
		if err != nil {
			return
		}
		for _, n := range ns {
			v := t.idx[n.Key]
			if !vst[v] {
				dfs(v)
				duf.SetParent(v, u)
			}
		}
		if u == v1 && vst[v2] {
			k, b = t.vtx[duf.Find(v2)], true
			return
		}
		if u == v2 && vst[v1] {
			k, b = t.vtx[duf.Find(v1)], true
			return
		}
	}
	dfs(t.idx[t.root])
	return
}
