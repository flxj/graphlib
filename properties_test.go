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
	"testing"
)

func testBridgeG(t *testing.T) Graph[int, int] {
	g := NewGraph[int, int](false, "")
	for i := 0; i < 8; i++ {
		_ = g.AddVertex(Vertex[int, int]{Key: i})
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
		_ = g.AddEdge(e)
	}
	t.Log("size=", g.Size(), " order=", g.Order())
	return g
}

func testIsBridge(t *testing.T) {
	g := testBridgeG(t)
	ok, err := IsBridge(g, 3)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !ok {
		t.Fatal("edge 3 is bridge")
	}

	if ok, err = IsBridge(g, 9); err != nil {
		t.Fatal(err.Error())
	}
	if !ok {
		t.Fatal("edge 9 is bridge")
	}
	//
	if ok, err = IsBridge(g, 0); err != nil {
		t.Fatal(err.Error())
	}
	if ok {
		t.Fatal("edge 0 is not bridge")
	}

	if ok, err = IsBridge(g, 5); err != nil {
		t.Fatal(err.Error())
	}
	if ok {
		t.Fatal("edge 5 is not bridge")
	}
	t.Log("> test bridge pass")
}

func testFindBridge(t *testing.T) {
	g := testBridgeG(t)

	b, err := FindBridges(g)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(b) != 2 {
		t.Fatal("find bridge wrong1")
	}
	if (b[0].Key == 3 && b[1].Key == 9) || (b[0].Key == 9 && b[1].Key == 3) {
		t.Log("=======> test find bridges pass")
	} else {
		t.Logf("key:%d (%d,%d)", b[0].Key, b[0].Head, b[0].Tail)
		t.Logf("key:%d (%d,%d)", b[1].Key, b[1].Head, b[1].Tail)
		t.Fatal("find bridge wrong2")
	}
}

func TestBridge(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "is":
		testIsBridge(t)
	case "find":
		testFindBridge(t)
	default:
	}
}
