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
	"testing"
)

func TestBasicOp(t *testing.T) {
	g := NewGraph[int, int](false, "test-g")

	vs := []Vertex[int, int]{
		{Key: 1, Value: 1},
		{Key: 2, Value: 2},
		{Key: 3, Value: 3},
		{Key: 4, Value: 4},
		{Key: 5, Value: 5},
		{Key: 6, Value: 6},
	}

	for _, v := range vs {
		_ = g.AddVertex(v)
	}

	es := []Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2},
		{Key: 2, Head: 1, Tail: 3},
		{Key: 3, Head: 2, Tail: 3},
		{Key: 4, Head: 4, Tail: 5},
		{Key: 5, Head: 5, Tail: 6},
	}
	gs := `
v1---v2
|   /
|  /   
v3     v4-----v5----v6
`

	for _, e := range es {
		_ = g.AddEdge(e)
	}
	//
	t.Log(">[0] init Pro")
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())

	ps, _ := g.Property(ProSimple)

	t.Logf("simple:%v", ps.Value)
	pc, _ := g.Property(ProConnected)

	t.Logf("connected:%v", pc.Value)
	pa, _ := g.Property(ProAcyclic)

	t.Logf("acyclic:%v", pa.Value)

	t.Log(">[1] delete vertex")

	gs = `
   v2
  /
 /   
v3     v4-----v5----v6
`

	_, _ = g.RemoveVertex(1)
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())

	ps, _ = g.Property(ProSimple)

	t.Logf("simple:%v", ps.Value)
	pc, _ = g.Property(ProConnected)

	t.Logf("connected:%v", pc.Value)
	pa, _ = g.Property(ProAcyclic)

	t.Logf("acyclic:%v", pa.Value)

	t.Log(">[2] add vertex")

	gs = `
    v2
   /
  /   
v3    v4-----v5----v6  v7
`

	v := Vertex[int, int]{Key: 7, Value: 7}
	_ = g.AddVertex(v)
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())

	ps, _ = g.Property(ProSimple)
	t.Logf("simple:%v", ps.Value)
	pc, _ = g.Property(ProConnected)
	t.Logf("connected:%v", pc.Value)
	pa, _ = g.Property(ProAcyclic)
	t.Logf("acyclic:%v", pa.Value)

	t.Log(">[3] add edges")
	gs = `
v2---v5----v6---v7
|    |
|    |
v3---v4 
`
	es = []Edge[int, int]{
		{Key: 6, Head: 2, Tail: 5},
		{Key: 7, Head: 3, Tail: 4},
		{Key: 8, Head: 7, Tail: 6},
	}
	for _, e := range es {
		_ = g.AddEdge(e)
	}
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())

	ps, _ = g.Property(ProSimple)
	t.Logf("simple:%v", ps.Value)
	pc, _ = g.Property(ProConnected)
	t.Logf("connected:%v", pc.Value)
	pa, _ = g.Property(ProAcyclic)
	t.Logf("acyclic:%v", pa.Value)

	t.Log(">[4] delete edge v3-v4")
	gs = `
v2---v5----v6---v7
|    |
|    |
v3   v4 
`
	_, _ = g.RemoveEdge(3, 4)
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())

	ps, _ = g.Property(ProSimple)
	t.Logf("simple:%v", ps.Value)
	pc, _ = g.Property(ProConnected)
	t.Logf("connected:%v", pc.Value)
	pa, _ = g.Property(ProAcyclic)
	t.Logf("acyclic:%v", pa.Value)

	t.Log(">[4] add edge v4-v7,v4-v5")
	gs = `
v2---v5----v6---v7
|    ||         /
|    ||        /
v3   v4------/ 
`
	es = []Edge[int, int]{
		{Key: 100, Head: 4, Tail: 7},
		{Key: 101, Head: 4, Tail: 5},
	}
	_ = g.AddEdge(es[0])
	_ = g.AddEdge(es[1])
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())

	ps, _ = g.Property(ProSimple)
	t.Logf("simple:%v", ps.Value)
	pc, _ = g.Property(ProConnected)
	t.Logf("connected:%v", pc.Value)
	pa, _ = g.Property(ProAcyclic)
	t.Logf("acyclic:%v", pa.Value)

}

func TestConnected(t *testing.T) {
	g := NewGraph[int, int](false, "test-g")

	vs := []Vertex[int, int]{
		{Key: 1, Value: 1},
		{Key: 2, Value: 2},
		{Key: 3, Value: 3},
		{Key: 4, Value: 4},
		{Key: 5, Value: 5},
		{Key: 6, Value: 6},
	}

	for _, v := range vs {
		_ = g.AddVertex(v)
	}

	//
	t.Log(">[0] init Pro")
	t.Logf("name:%s", g.Name())
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())
	pc, _ := g.Property(ProConnected)
	t.Logf("connected:%v", pc.Value)

	t.Log(">[1] add edges")

	es := []Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2},
		{Key: 2, Head: 1, Tail: 3},
		{Key: 3, Head: 2, Tail: 3},
		{Key: 4, Head: 4, Tail: 5},
		{Key: 5, Head: 5, Tail: 6},
		{Key: 6, Head: 4, Tail: 3},
	}

	for _, e := range es {
		_ = g.AddEdge(e)
	}

	t.Logf("name:%s", g.Name())
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())
	pc, _ = g.Property(ProConnected)
	t.Logf("connected:%v", pc.Value)

}
