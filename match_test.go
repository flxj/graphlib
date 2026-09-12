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

func bwmTestGraph(maximum bool) Bipartite[int, int] {
	g := NewBipartite[int, int](false, "g1")
	for i := 1; i < 9; i++ {
		if i <= 4 {
			_ = g.AddVertexTo(Vertex[int, int]{Key: i}, true)
		} else {
			_ = g.AddVertexTo(Vertex[int, int]{Key: i}, false)
		}
	}
	E := []Edge[int, int]{
		{Key: 0, Head: 1, Tail: 5, Weight: 82},
		{Key: 1, Head: 2, Tail: 5, Weight: 77},
		{Key: 2, Head: 3, Tail: 5, Weight: 11},
		{Key: 3, Head: 4, Tail: 5, Weight: 8},

		{Key: 4, Head: 1, Tail: 6, Weight: 83},
		{Key: 5, Head: 2, Tail: 6, Weight: 37},
		{Key: 6, Head: 3, Tail: 6, Weight: 69},
		{Key: 7, Head: 4, Tail: 6, Weight: 9},

		{Key: 8, Head: 1, Tail: 7, Weight: 69},
		{Key: 9, Head: 2, Tail: 7, Weight: 49},
		{Key: 10, Head: 3, Tail: 7, Weight: 5},
		{Key: 11, Head: 4, Tail: 7, Weight: 78},

		{Key: 12, Head: 1, Tail: 8, Weight: 92},
		{Key: 13, Head: 2, Tail: 8, Weight: 92},
		{Key: 14, Head: 3, Tail: 8, Weight: 86},
		{Key: 15, Head: 4, Tail: 8, Weight: 23},
	}
	for _, e := range E {
		if maximum {
			e.Weight = -1 * e.Weight
		}
		_ = g.AddEdge(e)
	}
	return g
}

func bwmG1(t *testing.T) {
	g := bwmTestGraph(false)
	t.Log("g1 order=", g.Order(), " size=", g.Size())
	res, err := BipartiteWeightedMatching(g, false)
	if err != nil {
		t.Fatal(err.Error())
	}
	var s int
	for _, e := range res {
		t.Logf("select edge:%d weight:%d", e.Key, e.Weight)
		s += e.Weight
	}
	if s != 140 {
		t.Fatal("g1 bwm wrong")
	}
	t.Log("> g1 pass")
}

func bwmG2(t *testing.T) {
	g := bwmTestGraph(true)
	t.Log("g2 order=", g.Order(), " size=", g.Size())
	res, err := BipartiteWeightedMatching(g, true)
	if err != nil {
		t.Fatal(err.Error())
	}
	var s int
	for _, e := range res {
		t.Logf("select edge:%d weight:%d", e.Key, e.Weight)
		s += e.Weight
	}
	if s != -140 {
		t.Fatal("g2 bwm wrong")
	}
	t.Log("> g2 pass")
}

func mvcG1(t *testing.T) {
	g := NewBipartite[int, int](false, "g1")
	A := []Vertex[int, int]{
		{Key: 1},
		{Key: 2},
		{Key: 3},
		{Key: 4},
	}
	B := []Vertex[int, int]{
		{Key: 5},
		{Key: 6},
		{Key: 7},
		{Key: 8},
	}
	for _, v := range A {
		_ = g.AddVertexTo(v, true)
	}
	for _, v := range B {
		_ = g.AddVertexTo(v, false)
	}
	E := []Edge[int, int]{
		{Key: 0, Head: 1, Tail: 5},
		{Key: 1, Head: 1, Tail: 7},
		{Key: 2, Head: 2, Tail: 5},
		{Key: 3, Head: 3, Tail: 5},
		{Key: 4, Head: 3, Tail: 6},
		{Key: 5, Head: 4, Tail: 7},
		{Key: 6, Head: 4, Tail: 8},
	}
	M := []Edge[int, int]{
		{Head: 1, Tail: 7},
		{Head: 2, Tail: 5},
		{Head: 3, Tail: 6},
		{Head: 4, Tail: 8},
	}
	for _, e := range E {
		_ = g.AddEdge(e)
	}
	C, err := bipartiteMVC(A, B, M, g)
	if err != nil {
		t.Fatal(err.Error())
	}
	for v := range C {
		t.Log("g1 mvc: ", v)
	}
	if len(C) != len(M) {
		t.Fatal("g1 mvc wrong")
	}
	t.Log("> g1 pass")
}

func mvcG2(t *testing.T) {
	g := NewBipartite[string, int](false, "g2")
	A := []Vertex[string, int]{
		{Key: "r1"},
		{Key: "r2"},
		{Key: "r3"},
		{Key: "r4"},
		{Key: "r5"},
		{Key: "r6"},
		{Key: "r7"},
		{Key: "r8"},
	}
	B := []Vertex[string, int]{
		{Key: "c1"},
		{Key: "c2"},
		{Key: "c3"},
		{Key: "c4"},
		{Key: "c5"},
		{Key: "c6"},
		{Key: "c7"},
		{Key: "c8"},
	}
	for _, v := range A {
		_ = g.AddVertexTo(v, true)
	}
	for _, v := range B {
		_ = g.AddVertexTo(v, false)
	}
	E := []Edge[string, int]{
		{Key: "11", Head: "r2", Tail: "c6"},
		{Key: "12", Head: "r2", Tail: "c7"},
		{Key: "13", Head: "r2", Tail: "c1"},
		{Key: "14", Head: "r2", Tail: "c3"},
		{Key: "15", Head: "r2", Tail: "c2"},
		{Key: "16", Head: "r6", Tail: "c6"},
		{Key: "17", Head: "r6", Tail: "c7"},
		{Key: "18", Head: "r6", Tail: "c1"},
		{Key: "19", Head: "r6", Tail: "c3"},
		{Key: "100", Head: "r6", Tail: "c2"},
		{Key: "101", Head: "r6", Tail: "c5"},
		{Key: "102", Head: "r4", Tail: "c2"},
		{Key: "103", Head: "r4", Tail: "c8"},
		{Key: "104", Head: "r5", Tail: "c2"},
		{Key: "105", Head: "r5", Tail: "c8"},
		{Key: "106", Head: "r7", Tail: "c2"},
		{Key: "107", Head: "r7", Tail: "c8"},
		{Key: "108", Head: "r1", Tail: "c2"},
		{Key: "109", Head: "r1", Tail: "c8"},
		{Key: "110", Head: "r1", Tail: "c5"},
		{Key: "111", Head: "r1", Tail: "c4"},
		{Key: "112", Head: "r8", Tail: "c2"},
		{Key: "113", Head: "r8", Tail: "c4"},
		{Key: "114", Head: "r3", Tail: "c5"},
		{Key: "115", Head: "r3", Tail: "c4"},
	}
	M := []Edge[string, int]{
		{Head: "r2", Tail: "c3"},
		{Head: "r6", Tail: "c1"},
		{Head: "r4", Tail: "c8"},
		{Head: "r5", Tail: "c2"},
		{Head: "r1", Tail: "c4"},
		{Head: "r3", Tail: "c5"},
	}
	for _, e := range E {
		_ = g.AddEdge(e)
	}
	C, err := bipartiteMVC(A, B, M, g)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(C) != len(M) {
		t.Fatal("g2 mvc wrong")
	}
	for v := range C {
		t.Log("g2 mvc: ", v)
	}
	t.Log("> g2 pass")
}

func testMaxMatching(t *testing.T) {
	g := PetersenGraph()
	m, err := MaxMatchingBlossom(g)
	if err != nil {
		t.Log(g.Name())
		t.Fatal(err.Error())
	}
	if len(m) != 5 {
		t.Fatalf("%s get mm=%d, but expected=%d", g.Name(), len(m), 5)
	}
}

func TestMatching(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "mvc":
		mvcG1(t)
		mvcG2(t)
	case "bwm_min":
		bwmG1(t)
	case "bwm_max":
		bwmG2(t)
	case "max":
		testMaxMatching(t)
	default:
	}
}
