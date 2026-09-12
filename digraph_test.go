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

func TestDigraph1(t *testing.T) {
	g := NewDigraph[int, int]("test-g")
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
		{Key: 2, Head: 2, Tail: 3},
		{Key: 3, Head: 5, Tail: 6},
		{Key: 4, Head: 4, Tail: 5},
		{Key: 5, Head: 2, Tail: 5},
	}

	for _, e := range es {
		_ = g.AddEdge(e)
	}
	gs := `
V1---> V2 ---> V3
       |
       v
V4---> V5 ---> V6
`
	t.Log(">[0] init property")
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())
	p, _ := g.Property(ProConnected)
	t.Logf("connected:%v", p.Value)
	p, _ = g.Property(ProUnilateralConnected)
	t.Logf("unidirectional connected:%v", p.Value)
	p, _ = g.Property(ProAcyclic)
	t.Logf("acyclic:%v", p.Value)

	t.Log(">[1] delete vertrx v4")
	_, _ = g.RemoveVertex(4)
	gs = `
V1---> V2 ---> V3
       |
       v
       V5 ---> V6
`
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())
	p, _ = g.Property(ProConnected)
	t.Logf("connected:%v", p.Value)
	p, _ = g.Property(ProUnilateralConnected)
	t.Logf("unidirectional connected:%v", p.Value)
	p, _ = g.Property(ProAcyclic)
	t.Logf("acyclic:%v", p.Value)

	t.Log(">[2] add edge v5->v1")
	ed := Edge[int, int]{Key: 10, Head: 5, Tail: 1}
	_ = g.AddEdge(ed)
	gs = `
V1---> V2 ---> V3
^      |
|      v
 \---  V5 ---> V6
`
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())
	p, _ = g.Property(ProConnected)
	t.Logf("connected:%v", p.Value)
	p, _ = g.Property(ProUnilateralConnected)
	t.Logf("unidirectional connected:%v", p.Value)
	p, _ = g.Property(ProAcyclic)
	t.Logf("acyclic:%v", p.Value)

	t.Log(">[2] add edge v3->v6")
	ed = Edge[int, int]{Key: 11, Head: 3, Tail: 6}
	_ = g.AddEdge(ed)
	gs = `
V1---> V2 ---> V3
^      |       |
|      v       v
 \---  V5 ---> V6
`
	t.Log(gs)
	t.Logf("order:%d", g.Order())
	t.Logf("size:%d", g.Size())
	p, _ = g.Property(ProConnected)
	t.Logf("connected:%v", p.Value)
	p, _ = g.Property(ProUnilateralConnected)
	t.Logf("unidirectional connected:%v", p.Value)
	p, _ = g.Property(ProAcyclic)
	t.Logf("acyclic:%v", p.Value)
}
