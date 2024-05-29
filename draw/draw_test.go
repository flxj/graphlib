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

package draw

import (
	"fmt"
	"testing"

	"github.com/flxj/graphlib"
)

func TestDraw(t *testing.T) {
	g, err := graphlib.NewGraph[int, int, int](false, "test-g")
	if err != nil {
		fmt.Printf("new graph error:%v\n", err)
		return
	}

	vs := []graphlib.Vertex[int, int]{
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
	_ = g.SetVertexLabel(1, "color", "green")
	_ = g.SetVertexLabel(6, "color", "red")

	es := []graphlib.Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2, Weight: 5},
		{Key: 2, Head: 2, Tail: 3, Weight: 6},
		{Key: 3, Head: 5, Tail: 6, Weight: 7},
		{Key: 4, Head: 4, Tail: 5, Weight: 8},
		{Key: 5, Head: 2, Tail: 5, Weight: 9},
	}
	for _, e := range es {
		_ = g.AddEdge(e)
	}
	_ = g.SetEdgeLabelByKey(3, "color", "red")

	file, err := RenderHTML(g, true, "/tmp")
	if err != nil {
		fmt.Printf("draw error:%v\n", err)
		return
	}
	fmt.Println(file)
}

func TestDraw2(t *testing.T) {
	g, err := graphlib.NewGraph[int, int, int](true, "test-g")
	if err != nil {
		fmt.Printf("new graph error:%v\n", err)
		return
	}

	vs := []graphlib.Vertex[int, int]{
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
	_ = g.SetVertexLabel(1, "color", "green")
	_ = g.SetVertexLabel(6, "color", "red")

	es := []graphlib.Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2, Weight: 100},
		{Key: 2, Head: 2, Tail: 3, Weight: -6},
		{Key: 3, Head: 5, Tail: 6, Weight: 70},
		{Key: 4, Head: 4, Tail: 5, Weight: 200},
		{Key: 5, Head: 2, Tail: 5, Weight: 9},
	}
	for _, e := range es {
		fmt.Println("Weight ", e.Weight)
		_ = g.AddEdge(e)
	}
	_ = g.SetEdgeLabelByKey(3, "color", "red")

	file, err := RenderSVG(g, "", []string{}, true, "/tmp")
	if err != nil {
		fmt.Printf("draw error:%v\n", err)
		return
	}
	fmt.Println(file)
}
