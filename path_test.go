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

func TestPath1(t *testing.T) {
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
		{Key: 1, Head: 1, Tail: 2, Weight: 10},
		{Key: 2, Head: 1, Tail: 3, Weight: 4},
		{Key: 3, Head: 2, Tail: 3, Weight: 2},
		{Key: 4, Head: 2, Tail: 4, Weight: 8},
		{Key: 5, Head: 2, Tail: 5, Weight: 6},
		{Key: 6, Head: 3, Tail: 4, Weight: 15},
		{Key: 7, Head: 3, Tail: 5, Weight: 6},
		{Key: 8, Head: 4, Tail: 5, Weight: 1},
		{Key: 9, Head: 4, Tail: 6, Weight: 5},
		{Key: 10, Head: 5, Tail: 6, Weight: 12},
	}

	for _, e := range es {
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}

	paths, err := ShortestPaths[int, int, int](g, 1)
	if err != nil {
		fmt.Println("[Err] ", err)
		return
	}
	for _, p := range paths {
		fmt.Printf("source:%d target:%d  weight:%v\n", p.Source, p.Target, p.Weight)
	}

}
