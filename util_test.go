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
		value: MaxFloatDistance,
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
		value: MaxFloatDistance,
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
