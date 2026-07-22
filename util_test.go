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
	"math"
	"testing"
)

/*
func TestCostQueue(t *testing.T) {
	p := newCostQueue[int]()
	c1 := &item[int]{
		key:   1,
		value: 3.0,
	}
	c2 := &item[int]{
		key:   2,
		value: 1.0,
	}
	c3 := &item[int]{
		key:   3,
		value: -3.0,
	}
	c4 := &item[int]{
		key:   4,
		value: 13.0,
	}
	c5 := &item[int]{
		key:   5,
		value: math.MaxFloat64,
	}
	c6 := &item[int]{
		key:   6,
		value: 0.0,
	}
	p.Push(c1)
	p.Push(c2)
	p.Push(c3)
	p.Push(c4)
	p.Push(c5)
	p.Push(c6)
	if p.Len() != 6 {
		fmt.Println("queue length error")
		return
	}
	for p.Len() > 0 {
		v := p.Pop()
		fmt.Printf("key:%d,value:%f\n", v.key, v.value)
	}

	fmt.Println("=========")
	p.Push(c1)
	p.Push(c2)
	p.Push(c3)
	p.Push(c4)
	p.Update(1, 100.0)
	p.Update(4, -100.0)
	for p.Len() > 0 {
		v := p.Pop()
		fmt.Printf("key:%d,value:%f\n", v.key, v.value)
	}
}
*/

func TestPriorityQueue(t *testing.T) {
	p := newPriorityQueue[int, int, float64](func(p1, p2 float64) bool { return p1 < p2 })
	c1 := &item[int]{
		key:   1,
		value: 3.0,
	}
	c2 := &item[int]{
		key:   2,
		value: 1.0,
	}
	c3 := &item[int]{
		key:   3,
		value: -3.0,
	}
	c4 := &item[int]{
		key:   4,
		value: 13.0,
	}
	c5 := &item[int]{
		key:   5,
		value: math.MaxFloat64,
	}
	c6 := &item[int]{
		key:   6,
		value: 0.0,
	}
	p.Push(c1.key, 0, c1.value)
	p.Push(c2.key, 0, c2.value)
	p.Push(c3.key, 0, c3.value)
	p.Push(c4.key, 0, c4.value)
	p.Push(c5.key, 0, c5.value)
	p.Push(c6.key, 0, c6.value)
	if p.Len() != 6 {
		fmt.Println("queue length error")
		return
	}
	for p.Len() > 0 {
		k, v, p, _ := p.Pop()
		fmt.Printf("key:%d,value:%d priority:%f\n", k, v, p)
	}

	fmt.Println("=========")
	p.Push(c1.key, 0, c1.value)
	p.Push(c2.key, 0, c2.value)
	p.Push(c4.key, 0, c4.value)
	p.Push(c5.key, 0, c5.value)
	p.Update(1, 100.0)
	p.Update(4, -100.0)
	for p.Len() > 0 {
		k, v, p, _ := p.Pop()
		fmt.Printf("key:%d,value:%d priority:%f\n", k, v, p)
	}
}

func TestLCM(t *testing.T) {
	tree := NewThinTree[int]()
	es := []Edge[int, int]{
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
		{Head: 9, Tail: 17},
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
		if err := tree.AddVertex(Vertex[int, any]{
			Key: i,
		}); err != nil {
			panic(err.Error())
		}
	}
	for i, e := range es {
		e.Key = i + 1
		if err := tree.AddEdge(e); err != nil {
			panic(fmt.Sprintf("edge:(%d,%d) err:%s", e.Head, e.Tail, err.Error()))
		}
	}
	tree.SetRoot(1)

	fmt.Printf("Tree vettex=%d,edges=%d\n", tree.Order(), tree.Size())

	q := [][3]int{
		{2, 3, 1},
		{2, 9, 1},
		{6, 10, 1},
		{13, 7, 2},
		{11, 12, 5},
		{24, 26, 9},
		{22, 16, 7},
		{13, 30, 2},
		{25, 26, 9},
		{30, 28, 1},
	}

	for _, u := range q {
		v, _ := tree.LeastCommonAncestor(u[0], u[1])
		if v != u[2] {
			panic(fmt.Sprintf("{%d,%d} lcm should be %d,but get %d", u[0], u[1], u[2], v))
		}
	}
}
