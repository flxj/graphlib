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

func TestDFS1(t *testing.T) {
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
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	vis := func(v Vertex[int, int]) error {
		fmt.Printf("[Visit] vertex is %v\n", v.Key)
		return nil
	}

	if err := DFS(g, 5, vis); err != nil {
		fmt.Println("[ERR] ", err)
		return
	}
}

func TestDFS2(t *testing.T) {
	g, err := NewGraph[int, int, int](true, "test-g")
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
		{Key: 6, Head: 3, Tail: 4},
	}
	gs := `
v1--->v2
|   /
| /   
v
v3---->v4---->v5--->v6
`

	for _, e := range es {
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	vis := func(v Vertex[int, int]) error {
		fmt.Printf("[Visit] vertex is %v\n", v.Key)
		return nil
	}

	if err := DFS(g, 3, vis); err != nil {
		fmt.Println("[ERR] ", err)
		return
	}
	fmt.Println("===============")
	if err := DFS(g, 2, vis); err != nil {
		fmt.Println("[ERR] ", err)
		return
	}
	fmt.Println("===============")
	if err := DFS(g, 1, vis); err != nil {
		fmt.Println("[ERR] ", err)
		return
	}
}

func TestBFS1(t *testing.T) {
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
		{Key: 6, Head: 2, Tail: 5},
	}
	gs := `
v1---v2
|   /  \
|  /    \
v3      v5-----v6
        |
        |
        v4
`

	for _, e := range es {
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	vis := func(v Vertex[int, int]) error {
		fmt.Printf("[Visit] vertex is %v\n", v.Key)
		return nil
	}

	if err := BFS(g, 1, vis); err != nil {
		fmt.Println("[ERR] ", err)
		return
	}
}

func TestBFS2(t *testing.T) {
	g, err := NewGraph[int, int, int](true, "test-g")
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
		{Key: 4, Head: 5, Tail: 4},
		{Key: 5, Head: 5, Tail: 6},
		{Key: 6, Head: 2, Tail: 5},
	}
	gs := `
v1--->v2
|   /   \
| /      \
v         v
v3        v5---->v6
          |
          |
          v4
`

	for _, e := range es {
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}
	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	vis := func(v Vertex[int, int]) error {
		fmt.Printf("[Visit] vertex is %v\n", v.Key)
		return nil
	}

	if err := BFS(g, 2, vis); err != nil {
		fmt.Println("[ERR] ", err)
		return
	}
}
