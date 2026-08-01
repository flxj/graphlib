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
	"math/rand"
	"testing"
)

func testSplayTree(n int) {
	fmt.Println("=====> test splay tree")
	t := NewSplayTree[int, string](func(a, b int) int {
		if a > b {
			return 1
		} else if a == b {
			return 0
		}
		return -1
	})
	fmt.Println("=====> 0. init data")
	keys, vals := generateIntStr(0, n, true, 5)
	for i, k := range keys {
		t.Insert(k, vals[i])
	}
	if t.Len() != n {
		panic(fmt.Sprintf("init error,expect %d data,but actual %d ", n, t.Len()))
	}
	fmt.Println("=====> 1. test read")
	for i := 0; i < n/2; i++ {
		j := rand.Intn(n)
		v, ok := t.Search(keys[j])
		if !ok || v != vals[j] {
			fmt.Printf("i:%d (key:%d,val:%s) but get (ok=%v,v=%s)\n", i, keys[j], vals[j], ok, v)
			panic("read error")
		}
	}
	fmt.Println("=====> 2. test insert")
	var maxK int
	for _, k := range keys {
		if k > maxK {
			maxK = k
		}
	}
	keys2, vals2 := generateIntStr(maxK+1, n, true, 10)
	for i, k := range keys2 {
		t.Insert(k, vals2[i])
	}
	if t.Len() != 2*n {
		panic(fmt.Sprintf("insert error,expect %d data,but actual %d ", 2*n, t.Len()))
	}
	fmt.Println("=====> 3. test update")
	v := "value"
	for i := 0; i < n/2; i++ {
		k := keys2[rand.Intn(n)]
		t.Insert(k, v)
		vv, ok := t.Search(k)
		if !ok || v != vv {
			panic(fmt.Sprintf("update erroe,key:%d, expect value:%s but actual %s", k, v, vv))
		}
	}
	fmt.Println("=====> 4. test delete")
	for i := 0; i < n; i++ {
		v, ok := t.Delete(keys[i])
		if !ok || v != vals[i] {
			panic(fmt.Sprintf("delete erroe,key:%d", keys[i]))
		}
	}
	if t.Len() != n {
		panic(fmt.Sprintf("delete error,expect %d data,but actual %d ", n, t.Len()))
	}
}

func TestBST(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "splay":
		testSplayTree(100)
	default:
	}
}
