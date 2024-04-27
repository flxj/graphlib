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
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func exportTestGraph1() (Graph[int, int, int], bool) {
	g, err := NewGraph[int, int, int](false, "test-g")
	if err != nil {
		fmt.Printf("new graph error:%v\n", err)
		return nil, false
	}
	vs := []Vertex[int, int]{
		{Key: 1, Value: 1},
		{Key: 2, Value: 2},
		{Key: 3, Value: 3},
		{Key: 4, Value: 4},
		{Key: 5, Value: 5},
		{Key: 6, Value: 6},
	}
	es := []Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2, Value: "e1"},
		{Key: 2, Head: 1, Tail: 3, Value: "e2"},
		{Key: 3, Head: 2, Tail: 3, Value: "e3"},
		{Key: 4, Head: 4, Tail: 5, Value: "e4"},
		{Key: 5, Head: 5, Tail: 6, Value: "e5"},
	}
	for i, v := range vs {
		if err := g.AddVertex(v); err != nil {
			fmt.Printf("add vertex error:%v\n", err)
			return nil, false
		}
		if err := g.SetVertexLabel(v.Key, "name", fmt.Sprintf("vertex-%d", i)); err != nil {
			fmt.Printf("add vertex label error:%v\n", err)
			return nil, false
		}
	}
	for i, e := range es {
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return nil, false
		}
		if err := g.SetEdgeLabelByKey(e.Key, "name", fmt.Sprintf("edge-%d", i)); err != nil {
			fmt.Printf("add edge label error:%v\n", err)
			return nil, false
		}
	}
	return g, true
}

func TestMarshalJSON(t *testing.T) {
	g, ok := exportTestGraph1()
	if !ok {
		return
	}
	//
	fmt.Printf("name:%s\n", g.Name())
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	fmt.Println("==================> marshal")
	s, err := MarshalGraphToJSON[int, int, int](g)
	if err != nil {
		fmt.Printf("marshal graph error:%v\n", err)
		return
	}
	fmt.Println("==================> json")
	var bf bytes.Buffer
	if err := json.Indent(&bf, s, "", "  "); err != nil {
		fmt.Printf("output graph error:%v\n", err)
		return
	}
	fmt.Println(bf.String())
}

func TestUnmarshalJSON(t *testing.T) {
	g, ok := exportTestGraph1()
	if !ok {
		return
	}
	fmt.Printf("name:%s\n", g.Name())
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())
	fmt.Println("==================> marshal")
	s, err := MarshalGraphToJSON[int, int, int](g)
	if err != nil {
		fmt.Printf("marshal graph error:%v\n", err)
		return
	}
	fmt.Println("==================> unmarshal")
	g2, err := UnmarshalGraph[int, int, int](s)
	if err != nil {
		fmt.Printf("unmarshal graph error:%v\n", err)
		return
	}
	fmt.Printf("name:%s\n", g2.Name())
	fmt.Printf("order:%d\n", g2.Order())
	fmt.Printf("size:%d\n", g2.Size())
}
