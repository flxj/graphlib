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

func TestBasicOp(t *testing.T) {
	g, err := NewGraph[int, int, int](false, "test-g")
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
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}
	//
	fmt.Println("=================>[0] init property")
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	ps, err := g.Property(PropertySimple)
	if err != nil {
		fmt.Printf("get property simple error:%v\n", err)
		return
	}
	fmt.Printf("simple:%v\n", ps.Value)
	pc, err := g.Property(PropertyConnected)
	if err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", pc.Value)
	pa, err := g.Property(PropertyAcyclic)
	if err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", pa.Value)

	fmt.Println("=====================>[1] delete vertex")

	gs = `
   v2
  /
 /   
v3     v4-----v5----v6
`

	if err := g.RemoveVertex(1); err != nil {
		fmt.Printf("delete vertex error:%v\n", err)
		return
	}
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	ps, err = g.Property(PropertySimple)
	if err != nil {
		fmt.Printf("get property simple error:%v\n", err)
		return
	}
	fmt.Printf("simple:%v\n", ps.Value)
	pc, err = g.Property(PropertyConnected)
	if err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", pc.Value)
	pa, err = g.Property(PropertyAcyclic)
	if err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", pa.Value)

	fmt.Println("=====================>[2] add vertex")

	gs = `
    v2
   /
  /   
v3    v4-----v5----v6  v7
`

	v := Vertex[int, int]{Key: 7, Value: 7}
	if err := g.AddVertex(v); err != nil {
		fmt.Printf("add vertex error:%v\n", err)
		return
	}
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	if ps, err = g.Property(PropertySimple); err != nil {
		fmt.Printf("get property simple error:%v\n", err)
		return
	}
	fmt.Printf("simple:%v\n", ps.Value)
	if pc, err = g.Property(PropertyConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", pc.Value)
	if pa, err = g.Property(PropertyAcyclic); err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", pa.Value)

	fmt.Println("=====================>[3] add edges")
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
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	if ps, err = g.Property(PropertySimple); err != nil {
		fmt.Printf("get property simple error:%v\n", err)
		return
	}
	fmt.Printf("simple:%v\n", ps.Value)
	if pc, err = g.Property(PropertyConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", pc.Value)
	if pa, err = g.Property(PropertyAcyclic); err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", pa.Value)

	fmt.Println("=====================>[4] delete edge v3-v4")
	gs = `
v2---v5----v6---v7
|    |
|    |
v3   v4 
`
	if err := g.RemoveEdge(3, 4); err != nil {
		fmt.Printf("delete edge error:%v\n", err)
		return
	}
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	if ps, err = g.Property(PropertySimple); err != nil {
		fmt.Printf("get property simple error:%v\n", err)
		return
	}
	fmt.Printf("simple:%v\n", ps.Value)
	if pc, err = g.Property(PropertyConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", pc.Value)
	if pa, err = g.Property(PropertyAcyclic); err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", pa.Value)

	fmt.Println("=====================>[4] add edge v4-v7,v4-v5")
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
	if err := g.AddEdge(es[0]); err != nil {
		fmt.Printf("add edge error:%v\n", err)
		return
	}
	if err := g.AddEdge(es[1]); err != nil {
		fmt.Printf("add edge error:%v\n", err)
		return
	}
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	if ps, err = g.Property(PropertySimple); err != nil {
		fmt.Printf("get property simple error:%v\n", err)
		return
	}
	fmt.Printf("simple:%v\n", ps.Value)
	if pc, err = g.Property(PropertyConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", pc.Value)
	if pa, err = g.Property(PropertyAcyclic); err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", pa.Value)

}

func TestConnected(t *testing.T) {
	g, err := NewGraph[int, int, int](false, "test-g")
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

	//
	fmt.Println("=================>[0] init property")
	fmt.Printf("name:%s\n", g.Name())
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())
	pc, err := g.Property(PropertyConnected)
	if err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", pc.Value)

	fmt.Println("=================>[1] add edges")

	es := []Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2},
		{Key: 2, Head: 1, Tail: 3},
		{Key: 3, Head: 2, Tail: 3},
		{Key: 4, Head: 4, Tail: 5},
		{Key: 5, Head: 5, Tail: 6},
		{Key: 6, Head: 4, Tail: 3},
	}

	for _, e := range es {
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}

	fmt.Printf("name:%s\n", g.Name())
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())
	if pc, err = g.Property(PropertyConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", pc.Value)

}
