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
	"flag"
	"fmt"
	"testing"

	"github.com/flxj/graphlib"
)

func testDraw1() {
	g := graphlib.NewGraph[int, int](false, "test-g")

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

func testDraw2() {
	g := graphlib.NewGraph[int, int](true, "test-g")

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

func testTexTree(di bool) {
	f := graphlib.NewForest[int, int]()
	es := []graphlib.Edge[int, int]{
		{Head: 1, Tail: 2},
		{Head: 1, Tail: 3},
		{Head: 1, Tail: 4},
		{Head: 1, Tail: 5},
		{Head: 2, Tail: 6},
		{Head: 2, Tail: 7},
		{Head: 2, Tail: 8},
		{Head: 3, Tail: 9},
		{Head: 4, Tail: 10},
		{Head: 5, Tail: 11},
		{Head: 5, Tail: 12},
		{Head: 6, Tail: 13},
		{Head: 7, Tail: 14},
		{Head: 7, Tail: 15},
		{Head: 7, Tail: 16},
		//{Head: 9, Tail: 17},
		{Head: 9, Tail: 18},
		{Head: 9, Tail: 19},
		{Head: 12, Tail: 20},
		{Head: 12, Tail: 21},
		{Head: 14, Tail: 22},
		{Head: 14, Tail: 23},
		{Head: 17, Tail: 24},
		{Head: 17, Tail: 25},
		{Head: 18, Tail: 26},
		{Head: 20, Tail: 27},
		{Head: 20, Tail: 28},
		{Head: 23, Tail: 29},
		{Head: 23, Tail: 30},
	}
	for i := 1; i <= 30; i++ {
		if err := f.AddVertex(graphlib.Vertex[int, int]{Key: i}); err != nil {
			panic(err.Error())
		}
	}
	for i, e := range es {
		e.Key = i + 1
		//e.Tail, e.Head = e.Head, e.Tail
		if err := f.AddEdge(e); err != nil {
			panic(fmt.Sprintf("edge:(%d,%d) err:%s", e.Head, e.Tail, err.Error()))
		}
	}
	f.SetRoot(1)
	if di {
		f.SetDirected(1)
	}
	str, err := TikzDrawForest(f)
	if err != nil {
		fmt.Println(err)
		panic("draw error")
	}
	fmt.Println(str)
}

func testTexD() {
	g := graphlib.RandomTournament(5)
	fmt.Println("order=", g.Order(), " size=", g.Size())
	str, err := TikzDrawDigraph(g)
	if err != nil {
		fmt.Println(err)
		panic("draw error")
	}
	fmt.Println(str)
}

func testTexG() {
	//g := graphlib.CompleteGraph(5)
	g := graphlib.PetersenGraph()
	str, err := TikzDrawGraph(g)
	if err != nil {
		fmt.Println(err)
		panic("draw error")
	}
	fmt.Println(str)
}

func TestDraw(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "draw":
		testDraw1()
		testDraw2()
	case "tex-tree":
		testTexTree(false)
	case "tex-dig":
		testTexD()
	case "tex-g":
		testTexG()
	default:
	}
}
