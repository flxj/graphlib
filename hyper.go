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

import "math/rand"

// This interface represents a hypergraph,
// which differs from common graphs in that its edges can contain any number of vertices.
type HyperGraph[K comparable, W number] interface {
	Name() string
	SetName(name string)
	Order() int
	Size() int
	AllVertexes() []Vertex[K, W]
	AllEdges() []HyperEdge[K, W]
	AddVertex(vertex Vertex[K, W]) bool
	RemoveVertex(key K) (Vertex[K, W], bool)
	AddEdge(edge HyperEdge[K, W]) bool
	RemoveEdgeByKey(key K) (HyperEdge[K, W], bool)
	// Delete all edges containing vtx as the vertex subset.
	// If the exact parameter is true, delete edges with a vertex set equal to vtx.
	RemoveEdge(vtx []K, exact bool) ([]HyperEdge[K, W], bool)
	RemoveAllEdge()
	Degree(vertex K) (int, bool)
	Neighbours(vertex K) ([]Vertex[K, W], bool)
	GetVertex(key K) (Vertex[K, W], bool)
	// Query all edges containing vtx as the vertex subset.
	// If the exact parameter is true, then query the edges of the vertex subset equal to vtx.
	GetEdge(vtx []K, exact bool) ([]HyperEdge[K, W], bool)
	GetEdgeByKey(key K) (HyperEdge[K, W], bool)
	GetVertexesByLabel(labels Labels) []Vertex[K, W]
	GetEdgesByLabel(labels Labels) []HyperEdge[K, W]
	SetVertexValue(key K, value any) bool
	SetVertexLabel(key K, labelKey string, labelVal any) bool
	DeleteVertexLabel(key K, labelKey string) bool
	SetVertexWeight(key K, weight W) bool
	SetEdgeWeight(key K, weight W) bool
	SetEdgeValueByKey(key K, value any) bool
	SetEdgeLabelByKey(key K, labelKey string, labelVal any) bool
	DeleteEdgeLabelByKey(key K, labelKey string) bool
	// Set the value of all edges that contain vtx as the vertex subset.
	// If the exact parameter is true, only set the vertex set to be equal to vtx edges.
	SetEdgeValue(vtx []K, value any, exact bool) bool
	// Set the labels for all edges that contain vtx as the vertex subset.
	// If the exact parameter is true, only set the edges whose vertex set is equal to vtx.
	SetEdgeLabel(vtx []K, labelKey string, labelVal any, exact bool) bool
	DeleteEdgeLabel(vtx []K, labelKey string, exact bool) bool
	Clone() HyperGraph[K, W]
	RandomVertex() (Vertex[K, W], bool)
	RandomEdge() (HyperEdge[K, W], bool)
	NeighbourEdgesByKey(edge K) ([]HyperEdge[K, W], bool)
	NeighbourEdges(vtx []K) ([]HyperEdge[K, W], bool)
	IncidentEdges(vertex K) ([]HyperEdge[K, W], bool)
}

// Hypergraph edge, Key is used to uniquely identify the edge,
// Weight represents the weight of the edge,
// and Vtx identifies the vertices contained in the edge.
type HyperEdge[K comparable, W number] struct {
	Key    K              `json:"key" yaml:"key"`
	Weight W              `json:"weight" yaml:"weight"`
	Value  any            `json:"value" yaml:"value"`
	Vtx    map[K]struct{} `json:"vtx" yaml:"vtx"`
	Labels map[string]any `json:"labels" yaml:"labels"`
}

func (e HyperEdge[K, W]) Size() int {
	return len(e.Vtx)
}

func (e HyperEdge[K, W]) Clone() HyperEdge[K, W] {
	h := HyperEdge[K, W]{
		Key:    e.Key,
		Weight: e.Weight,
		Value:  e.Value,
		Vtx:    make(map[K]struct{}),
	}
	for k := range e.Vtx {
		h.Vtx[k] = struct{}{}
	}
	if e.Labels != nil {
		h.Labels = make(map[string]any)
		for k, v := range e.Labels {
			h.Labels[k] = v
		}
	}
	return h
}

func NewHyperGraph[K comparable, W number](name string) HyperGraph[K, W] {
	bi := NewBipartite[int, int](false, name)
	return &hypergraph[K, W]{
		bi:   bi,
		vtx:  make(map[int]Vertex[K, W]),
		edge: make(map[int]HyperEdge[K, W]),
		vIdx: make(map[K]int),
		eIdx: make(map[K]int),
	}
}

type hypergraph[K comparable, W number] struct {
	bi   Bipartite[int, int]
	seq  int
	key  int
	vtx  map[int]Vertex[K, W]
	edge map[int]HyperEdge[K, W]
	vIdx map[K]int
	eIdx map[K]int
}

func (h *hypergraph[K, W]) Name() string {
	return h.bi.Name()
}

func (h *hypergraph[K, W]) SetName(name string) {
	h.bi.SetName(name)
}

func (h *hypergraph[K, W]) Order() int {
	return h.bi.PartOrder(true)
}

func (h *hypergraph[K, W]) Size() int {
	return h.bi.PartOrder(false)
}

func (h *hypergraph[K, W]) AllVertexes() []Vertex[K, W] {
	vtx := make([]Vertex[K, W], len(h.vtx))
	var i int
	for _, v := range h.vtx {
		vtx[i] = v
		i++
	}
	return vtx
}

func (h *hypergraph[K, W]) AllEdges() []HyperEdge[K, W] {
	es := make([]HyperEdge[K, W], len(h.edge))
	var i int
	for _, e := range h.edge {
		es[i] = e
		i++
	}
	return es
}

func (h *hypergraph[K, W]) AddVertex(vertex Vertex[K, W]) bool {
	if _, ok := h.vIdx[vertex.Key]; ok {
		return false
	}
	_ = h.bi.AddVertexTo(Vertex[int, int]{Key: h.key}, true)
	h.vtx[h.key] = vertex
	h.vIdx[vertex.Key] = h.key
	h.key++
	return true
}

func (h *hypergraph[K, W]) incidentEdges(v int) []Vertex[int, int] {
	ns, _ := h.bi.Neighbours(v)
	return ns
}

func (h *hypergraph[K, W]) RemoveVertex(key K) (Vertex[K, W], bool) {
	v, ok := h.vIdx[key]
	if !ok {
		return Vertex[K, W]{}, false
	}
	for _, e := range h.incidentEdges(v) {
		_, _ = h.bi.RemoveVertex(e.Key)
		ek := h.edge[e.Key]
		delete(h.edge, e.Key)
		delete(h.eIdx, ek.Key)
	}
	_, _ = h.bi.RemoveVertex(v)
	vt := h.vtx[v]
	delete(h.vIdx, key)
	delete(h.vtx, v)
	return vt, true
}

func (h *hypergraph[K, W]) AddEdge(e HyperEdge[K, W]) bool {
	if _, ok := h.eIdx[e.Key]; ok {
		return false
	}
	if len(e.Vtx) == 0 {
		return false
	}
	for v := range e.Vtx {
		if _, ok := h.vIdx[v]; !ok {
			return false
		}
	}
	// add edge
	_ = h.bi.AddVertexTo(Vertex[int, int]{Key: h.key}, false)
	h.eIdx[e.Key] = h.key
	h.edge[h.key] = e

	for v := range e.Vtx {
		vi := h.vIdx[v]
		_ = h.bi.AddEdge(Edge[int, int]{Key: h.seq, Head: vi, Tail: h.key})
		h.seq++
	}
	h.key++
	return true
}

func (h *hypergraph[K, W]) RemoveEdgeByKey(key K) (HyperEdge[K, W], bool) {
	e, ok := h.eIdx[key]
	if !ok {
		return HyperEdge[K, W]{}, false
	}
	_, _ = h.bi.RemoveVertex(e)
	he := h.edge[e]
	delete(h.eIdx, key)
	delete(h.edge, e)
	return he, true
}

// if exact=true, then try to delete edges that
// else delete all edges that contains vtx as vertex subset.
func (h *hypergraph[K, W]) RemoveEdge(vtx []K, exact bool) ([]HyperEdge[K, W], bool) {
	if len(vtx) == 0 {
		return nil, false
	}
	v0, ok := h.vIdx[vtx[0]]
	if !ok {
		return nil, false
	}
	var res []HyperEdge[K, W]
	for _, ev := range h.incidentEdges(v0) {
		e := h.edge[ev.Key]
		flag := true
		if len(e.Vtx) >= len(vtx) {
			for _, k := range vtx {
				if _, ok := e.Vtx[k]; !ok {
					flag = false
					break
				}
			}
		}
		if flag {
			if exact && len(e.Vtx) != len(vtx) {
				continue
			}
			_, _ = h.bi.RemoveVertex(ev.Key)
			res = append(res, e)
			delete(h.eIdx, e.Key)
			delete(h.edge, ev.Key)
		}
	}
	return res, true
}

func (h *hypergraph[K, W]) RemoveAllEdge() {
	h.bi.RemovePart(false)
	h.edge = make(map[int]HyperEdge[K, W])
	h.eIdx = make(map[K]int)
}

func (h *hypergraph[K, W]) Degree(vertex K) (int, bool) {
	v, ok := h.vIdx[vertex]
	if !ok {
		return 0, false
	}
	return h.bi.Degree(v)
}

func (h *hypergraph[K, W]) Neighbours(vertex K) ([]Vertex[K, W], bool) {
	v, ok := h.vIdx[vertex]
	if !ok {
		return nil, false
	}
	idx := make(map[K]struct{})
	for _, ev := range h.incidentEdges(v) {
		e := h.edge[ev.Key]
		for u := range e.Vtx {
			idx[u] = struct{}{}
		}
	}
	var vs []Vertex[K, W]
	for k := range idx {
		i := h.vIdx[k]
		vs = append(vs, h.vtx[i])
	}
	return vs, true
}

func (h *hypergraph[K, W]) GetVertex(key K) (Vertex[K, W], bool) {
	i, ok := h.vIdx[key]
	if !ok {
		return Vertex[K, W]{}, false
	}
	return h.vtx[i], true
}

func (h *hypergraph[K, W]) GetEdge(vtx []K, exact bool) ([]HyperEdge[K, W], bool) {
	if len(vtx) == 0 {
		return nil, false
	}
	v0, ok := h.vIdx[vtx[0]]
	if !ok {
		return nil, false
	}
	var res []HyperEdge[K, W]
	for _, ev := range h.incidentEdges(v0) {
		e := h.edge[ev.Key]
		flag := true
		if len(e.Vtx) >= len(vtx) {
			for _, k := range vtx {
				if _, ok := e.Vtx[k]; !ok {
					flag = false
					break
				}
			}
		}
		if flag {
			if exact && len(e.Vtx) != len(vtx) {
				continue
			}
			res = append(res, e)
		}
	}
	return res, true
}

func (h *hypergraph[K, W]) GetEdgeByKey(key K) (HyperEdge[K, W], bool) {
	i, ok := h.eIdx[key]
	if !ok {
		return HyperEdge[K, W]{}, false
	}
	return h.edge[i], true
}

func (h *hypergraph[K, W]) GetVertexesByLabel(labels Labels) []Vertex[K, W] {
	var ves []Vertex[K, W]
	if labels != nil {
		for _, u := range h.vtx {
			if u.Labels != nil {
				match := true
				for k, v := range labels {
					l, ok := u.Labels[k]
					if !ok || l != v {
						match = false
						break
					}
				}
				if match {
					ves = append(ves, u)
				}
			}
		}
	}
	return ves
}

func (h *hypergraph[K, W]) GetEdgesByLabel(labels Labels) []HyperEdge[K, W] {
	var edges []HyperEdge[K, W]
	if labels != nil {
		for _, e := range h.edge {
			if e.Labels != nil {
				match := true
				for k, v := range labels {
					l, ok := e.Labels[k]
					if !ok || l != v {
						match = false
						break
					}
				}
				if match {
					edges = append(edges, e)
				}
			}
		}
	}
	return edges
}

func (h *hypergraph[K, W]) SetVertexValue(key K, value any) bool {
	i, ok := h.vIdx[key]
	if !ok {
		return false
	}
	v := h.vtx[i]
	v.Value = value
	h.vtx[i] = v
	return true
}

func (h *hypergraph[K, W]) SetVertexLabel(key K, labelKey string, labelVal any) bool {
	i, ok := h.vIdx[key]
	if !ok {
		return false
	}
	v := h.vtx[i]
	if v.Labels == nil {
		v.Labels = make(map[string]any)
	}
	v.Labels[labelKey] = labelVal
	h.vtx[i] = v
	return true
}

func (h *hypergraph[K, W]) DeleteVertexLabel(key K, labelKey string) bool {
	i, ok := h.vIdx[key]
	if !ok {
		return false
	}
	v := h.vtx[i]
	delete(v.Labels, labelKey)
	h.vtx[i] = v
	return true
}

func (h *hypergraph[K, W]) SetVertexWeight(key K, weight W) bool {
	i, ok := h.vIdx[key]
	if !ok {
		return false
	}
	v := h.vtx[i]
	v.Weight = weight
	h.vtx[i] = v
	return true
}

func (h *hypergraph[K, W]) SetEdgeWeight(key K, weight W) bool {
	i, ok := h.eIdx[key]
	if !ok {
		return false
	}
	e := h.edge[i]
	e.Weight = weight
	h.edge[i] = e
	return true
}

func (h *hypergraph[K, W]) SetEdgeValueByKey(key K, value any) bool {
	i, ok := h.eIdx[key]
	if !ok {
		return false
	}
	e := h.edge[i]
	e.Value = value
	h.edge[i] = e
	return true
}

func (h *hypergraph[K, W]) SetEdgeLabelByKey(key K, labelKey string, labelVal any) bool {
	i, ok := h.eIdx[key]
	if !ok {
		return false
	}
	e := h.edge[i]
	if e.Labels == nil {
		e.Labels = make(map[string]any)
	}
	e.Labels[labelKey] = labelVal
	h.edge[i] = e
	return true
}

func (h *hypergraph[K, W]) DeleteEdgeLabelByKey(key K, labelKey string) bool {
	i, ok := h.eIdx[key]
	if !ok {
		return false
	}
	e := h.edge[i]
	delete(e.Labels, labelKey)
	h.edge[i] = e
	return true
}

func (h *hypergraph[K, W]) SetEdgeValue(vtx []K, value any, exact bool) bool {
	es, ok := h.GetEdge(vtx, exact)
	if !ok {
		return false
	}
	for _, e := range es {
		i := h.eIdx[e.Key]
		e.Value = value
		h.edge[i] = e
	}
	return true
}

func (h *hypergraph[K, W]) SetEdgeLabel(vtx []K, labelKey string, labelVal any, exact bool) bool {
	es, ok := h.GetEdge(vtx, exact)
	if !ok {
		return false
	}
	for _, e := range es {
		i := h.eIdx[e.Key]
		if e.Labels == nil {
			e.Labels = make(map[string]any)
		}
		e.Labels[labelKey] = labelVal
		h.edge[i] = e
	}
	return true
}

func (h *hypergraph[K, W]) DeleteEdgeLabel(vtx []K, labelKey string, exact bool) bool {
	es, ok := h.GetEdge(vtx, exact)
	if !ok {
		return false
	}
	for _, e := range es {
		i := h.eIdx[e.Key]
		delete(e.Labels, labelKey)
		h.edge[i] = e
	}
	return true
}

func (h *hypergraph[K, W]) Clone() HyperGraph[K, W] {
	bi := h.bi.Clone()
	b, _ := bi.(Bipartite[int, int])

	nh := &hypergraph[K, W]{
		bi:   b,
		key:  h.key,
		vIdx: make(map[K]int),
		eIdx: make(map[K]int),
		vtx:  make(map[int]Vertex[K, W]),
		edge: make(map[int]HyperEdge[K, W]),
	}
	for k, v := range h.vIdx {
		nh.vIdx[k] = v
	}
	for k, v := range h.eIdx {
		nh.eIdx[k] = v
	}
	for k, v := range h.vtx {
		nh.vtx[k] = v.Clone()
	}
	for k, e := range h.edge {
		nh.edge[k] = e.Clone()
	}
	return nh
}

func (h *hypergraph[K, W]) RandomVertex() (Vertex[K, W], bool) {
	if len(h.vIdx) == 0 {
		return Vertex[K, W]{}, false
	}
	n := rand.Intn(len(h.vIdx))
	for _, i := range h.vIdx {
		n--
		if n < 0 {
			return h.vtx[i], true
		}
	}
	return Vertex[K, W]{}, false
}

func (h *hypergraph[K, W]) RandomEdge() (HyperEdge[K, W], bool) {
	if len(h.eIdx) == 0 {
		return HyperEdge[K, W]{}, false
	}
	n := rand.Intn(len(h.eIdx))
	for _, i := range h.eIdx {
		n--
		if n < 0 {
			return h.edge[i], true
		}
	}
	return HyperEdge[K, W]{}, false
}

func (h *hypergraph[K, W]) NeighbourEdgesByKey(key K) ([]HyperEdge[K, W], bool) {
	i, ok := h.eIdx[key]
	if !ok {
		return nil, false
	}
	e := h.edge[i]
	mp := make(map[int]struct{})
	for v := range e.Vtx {
		j := h.vIdx[v]
		for _, ne := range h.incidentEdges(j) {
			mp[ne.Key] = struct{}{}
		}
	}
	var res []HyperEdge[K, W]
	for i := range mp {
		res = append(res, h.edge[i])
	}
	return res, true
}

func (h *hypergraph[K, W]) NeighbourEdges(vtx []K) ([]HyperEdge[K, W], bool) {
	mp := make(map[int]struct{})
	for _, v := range vtx {
		j := h.vIdx[v]
		for _, ne := range h.incidentEdges(j) {
			mp[ne.Key] = struct{}{}
		}
	}
	var res []HyperEdge[K, W]
	for i := range mp {
		res = append(res, h.edge[i])
	}
	return res, true
}

func (h *hypergraph[K, W]) IncidentEdges(vertex K) ([]HyperEdge[K, W], bool) {
	v, ok := h.vIdx[vertex]
	if !ok {
		return nil, false
	}
	var res []HyperEdge[K, W]
	for _, ev := range h.incidentEdges(v) {
		e, ok := h.edge[ev.Key]
		if !ok {
			return nil, false
		}
		res = append(res, e)
	}
	return res, true
}
