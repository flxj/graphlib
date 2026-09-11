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

func TestColour(t *testing.T) {
	g := PetersenGraph()
	col, x, err := GreedyVertexColouring(g)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("colour=", x)
	for v, c := range col {
		fmt.Println("v=", v, " c=", c)
	}
}

func TestMCQ(t *testing.T) {
	//
	k5 := CompleteGraph(5)

	for v := 5; v <= 11; v++ {
		_ = k5.AddVertex(Vertex[int, int]{Key: v})
	}
	edges := []Edge[int, int]{
		{Head: 1, Tail: 5},
		{Head: 5, Tail: 6},
		{Head: 6, Tail: 7},
		{Head: 5, Tail: 7},

		{Head: 0, Tail: 8},
		{Head: 0, Tail: 9},
		{Head: 0, Tail: 10},
		{Head: 8, Tail: 9},
		{Head: 8, Tail: 10},
		{Head: 9, Tail: 10},

		{Head: 11, Tail: 1},
		{Head: 11, Tail: 2},
		{Head: 11, Tail: 3},
		{Head: 11, Tail: 4},
	}
	ek := 10
	for _, e := range edges {
		e.Key = ek
		_ = k5.AddEdge(e)
		ek++
	}
	c, err := MaximumClique(k5)
	if err != nil {
		t.Error(err)
	}
	if len(c) != 5 {
		t.Errorf("find %d, but expect 5", len(c))
	}
	for _, k := range c {
		fmt.Println(k)
	}
}
