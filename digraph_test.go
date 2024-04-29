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
	"testing"
)

func TestDigraph1(t *testing.T) {
	g, err := NewDigraph[int, int, int]("test-g")
	if err != nil {
		fmt.Printf("new graph error:%v\n", err)
		return
	}

	vs := []Vertex[int, int]{
		{Key: 1, Value: 1},
		{Key: 2, Value: 2},
		{Key: 3, Value: 3},
		{Key: 4, Value: 4},
		{Key: 5, Value: 5},
		{Key: 6, Value: 6},
	}

	for _, v := range vs {
		if err := g.AddVertex(v); err != nil {
			fmt.Printf("add vertex error:%v\n", err)
			return
		}
	}

	es := []Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2},
		{Key: 2, Head: 2, Tail: 3},
		{Key: 3, Head: 5, Tail: 6},
		{Key: 4, Head: 4, Tail: 5},
		{Key: 5, Head: 2, Tail: 5},
	}

	for _, e := range es {
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}
	gs := `
V1---> V2 ---> V3
       |
       v
V4---> V5 ---> V6
`
	fmt.Println("=================>[0] init property")
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())
	p, err := g.Property(PropertyConnected)
	if err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", p.Value)
	if p, err = g.Property(PropertyUnilateralConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("unidirectional connected:%v\n", p.Value)
	if p, err = g.Property(PropertyAcyclic); err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", p.Value)

	fmt.Println("===================>[1] delete vertrx v4")
	if err := g.RemoveVertex(4); err != nil {
		fmt.Printf("delete edge error:%v\n", err)
		return
	}
	gs = `
V1---> V2 ---> V3
       |
       v
       V5 ---> V6
`
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())
	if p, err = g.Property(PropertyConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", p.Value)
	if p, err = g.Property(PropertyUnilateralConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("unidirectional connected:%v\n", p.Value)
	if p, err = g.Property(PropertyAcyclic); err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", p.Value)

	fmt.Println("===================>[2] add edge v5->v1")
	ed := Edge[int, int]{Key: 10, Head: 5, Tail: 1}
	if err := g.AddEdge(ed); err != nil {
		fmt.Printf("add edge error:%v\n", err)
		return
	}
	gs = `
V1---> V2 ---> V3
^      |
|      v
 \---  V5 ---> V6
`
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())
	if p, err = g.Property(PropertyConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", p.Value)
	if p, err = g.Property(PropertyUnilateralConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("unidirectional connected:%v\n", p.Value)
	if p, err = g.Property(PropertyAcyclic); err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", p.Value)

	fmt.Println("===================>[2] add edge v3->v6")
	ed = Edge[int, int]{Key: 11, Head: 3, Tail: 6}
	if err := g.AddEdge(ed); err != nil {
		fmt.Printf("add edge error:%v\n", err)
		return
	}
	gs = `
V1---> V2 ---> V3
^      |       |
|      v       v
 \---  V5 ---> V6
`
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())
	if p, err = g.Property(PropertyConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", p.Value)
	if p, err = g.Property(PropertyUnilateralConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("unidirectional connected:%v\n", p.Value)
	if p, err = g.Property(PropertyAcyclic); err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", p.Value)

}
