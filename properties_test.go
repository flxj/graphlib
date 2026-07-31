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
	"flag"
	"fmt"
	"testing"
)

func testBridgeG() Graph[int, any, int] {
	g, _ := NewGraph[int, any, int](false, "")
	for i := 0; i < 8; i++ {
		err := g.AddVertex(Vertex[int, any]{Key: i})
		if err != nil {
			panic(err.Error())
		}
	}
	E := []Edge[int, int]{
		{Head: 0, Tail: 1},
		{Head: 0, Tail: 2},
		{Head: 1, Tail: 2},
		{Head: 2, Tail: 3}, // 3 bridge
		{Head: 3, Tail: 4},
		{Head: 3, Tail: 5},
		{Head: 3, Tail: 6},
		{Head: 4, Tail: 5},
		{Head: 5, Tail: 6},
		{Head: 5, Tail: 7}, // 9 bridge
	}
	for i, e := range E {
		e.Key = i
		err := g.AddEdge(e)
		if err != nil {
			panic(err.Error())
		}
	}
	fmt.Println("size=", g.Size(), " order=", g.Order())
	return g
}

func testIsBridge() {
	g := testBridgeG()
	ok, err := IsBridge(g, 3)
	if err != nil {
		panic(err.Error())
	}
	if !ok {
		panic("edge 3 is bridge")
	}

	if ok, err = IsBridge(g, 9); err != nil {
		panic(err.Error())
	}
	if !ok {
		panic("edge 9 is bridge")
	}
	//
	if ok, err = IsBridge(g, 0); err != nil {
		panic(err.Error())
	}
	if ok {
		panic("edge 0 is not bridge")
	}

	if ok, err = IsBridge(g, 5); err != nil {
		panic(err.Error())
	}
	if ok {
		panic("edge 5 is not bridge")
	}
	fmt.Println("========> test bridge pass")
}

func testFindBridge() {
	g := testBridgeG()

	b, err := FindBridges(g)
	if err != nil {
		panic(err.Error())
	}
	if len(b) != 2 {
		panic("find bridge wrong1")
	}
	if (b[0].Key == 3 && b[1].Key == 9) || (b[0].Key == 9 && b[1].Key == 3) {
		fmt.Println("=======> test find bridges pass")
	} else {
		fmt.Printf("key:%d (%d,%d)\n", b[0].Key, b[0].Head, b[0].Tail)
		fmt.Printf("key:%d (%d,%d)\n", b[1].Key, b[1].Head, b[1].Tail)
		panic("find bridge wrong2")
	}
}

func TestBridge(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "is":
		testIsBridge()
	case "find":
		testFindBridge()
	default:
	}
}
