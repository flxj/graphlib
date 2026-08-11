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
	"strconv"
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

func testLCT() {
	fmt.Println("=====> test link-cut-tree")
	t := NewLinkCutTree[int, string, int](func(a, b int) int {
		if a > b {
			return 1
		} else if a == b {
			return 0
		} else {
			return -1
		}
	})
	fmt.Println("=========> 1. build tree T")
	//    5----1------0----4------11----15
	//        / |     | \         | \
	//	      5 6     2  3___     13 14
	//                |  | \ \
	//                7  8  9 10
	//                |
	//                12
	n := 16
	for i := 0; i < n; i++ {
		t.AddOrUpdate(i, "val-"+strconv.Itoa(i), i)
	}
	if t.Len() != n || t.Component() != n {
		panic("add node error")
	}
	edges := [][2]int{
		{0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 5}, {1, 6}, {2, 7}, {3, 8},
		{3, 9}, {3, 10}, {4, 11}, {7, 12}, {11, 13}, {11, 14}, {11, 15},
	}
	for i, e := range edges {
		fmt.Printf("add edge (%d,%d)\n", e[0], e[1])
		if !t.Link(e[0], e[1]) {
			panic("link error")
		}
		if t.Component() != n-i-1 {
			fmt.Printf("component %d, n-i-i=%d\n", t.Component(), n-i-1)
			panic("component error")
		}
	}
	fmt.Println("=========> 2. test Connected")
	for _, e := range [][2]int{{1, 12}, {3, 8}, {4, 7}, {6, 10}, {2, 4}} {
		if !t.Connected(e[0], e[1]) {
			panic("connected error")
		}
	}

	fmt.Println("=========> 3. test MakeRoot")
	if !t.MakeRoot(3) {
		panic("makeRoot 3 error")
	}
	if !t.IsRoot(3) {
		panic("isRoot error")
	}
	for _, v := range []int{15, 5, 2, 4} {
		r, _, ok := t.FindRoot(v)
		if !ok || r != 3 {
			fmt.Println("v=", v, " findRoot=", r)
			panic("findRoot error")
		}
	}
	if !t.MakeRoot(0) {
		panic("makeRoot 0 error")
	}
	fmt.Println("=========> 4. test PathSum")
	s := [][2]int{{12, 21}, {6, 7}, {7, 9}, {8, 11}, {11, 15}, {14, 29}}
	for _, p := range s {
		r, _ := t.PathSum(p[0])
		if r != p[1] {
			fmt.Println("v=", p[0], " sum=", r)
			panic("pathAgg error")
		}
	}
	fmt.Println("=========> 5. test PathPathAggregate")
	arr := [][2]int{{12, (0 | 2 | 7 | 12)}, {6, (0 | 1 | 6)}, {7, (0 | 2 | 7)},
		{8, (0 | 3 | 8)}, {11, (0 | 4 | 11)}, {14, (0 | 4 | 11 | 14)}}
	var pp int
	or := func(v int, _ string, _ int) {
		pp |= v
	}
	for _, p := range arr {
		pp = 0
		_ = t.PathAggregate(p[0], or)
		if pp != p[1] {
			fmt.Println("key=", p[0], " res=", p[1], " but get=", pp)
			panic("pathAgg error")
		}
	}
	fmt.Println("=========> 6. test MakeTree")
	//   17------16-----18-----21
	//                   | \
	//                   19 20
	//
	n2 := 6
	for i := 0; i < n2; i++ {
		if !t.MakeTree(i+n, "val-"+strconv.Itoa(i+n), i+n) {
			panic("makeTree error")
		}
	}
	edges2 := [][2]int{{16, 17}, {18, 19}, {18, 20}, {18, 21}, {18, 16}}
	for _, e := range edges2 {
		fmt.Printf("add edge (%d,%d)\n", e[0], e[1])
		if !t.Link(e[0], e[1]) {
			panic("link error 2")
		}
	}
	if t.Component() != 2 {
		panic("component error 2")
	}
	fmt.Println("=========> 7. test Cut")
	for _, v := range []int{7, 3, 11, 16} {
		if !t.Cut(v) {
			panic("cut error")
		}
	}
	if t.Component() != 6 {
		fmt.Printf("components %d, but expeect 6\n", t.Component())
		panic("cut component error")
	}
	if t.Connected(6, 11) || !t.Connected(5, 2) {
		panic("cut connected error")
	}
	fmt.Println("=========> 8. test Link")
	if !t.Link(18, 2) {
		panic("link 18-2 error")
	}
	if !t.Link(11, 0) {
		panic("link 11-0 error")
	}
	if t.Component() != 4 {
		panic("link component error 2")
	}
	if !t.Connected(6, 21) || t.Connected(8, 13) {
		panic("link connected error 3")
	}
	fmt.Println("=========> test pass")
}

func testTreapRW(n int, randKey bool) {
	fmt.Println("===========> testTreapReadWrite")
	k, v := generateIntStr(0, n, randKey, 20)
	t := NewTreap[int, string](func(a, b int) int {
		if a > b {
			return 1
		} else if a == b {
			return 0
		}
		return -1
	})
	for i := 0; i < n; i++ {
		t.Insert(k[i], v[i])
	}
	oldLen := t.Len()
	fmt.Printf("==========> 0 init date, size=%d\n", oldLen)
	fmt.Println("=========> 1 test random read...")
	// random read
	for i := 0; i < n/2; i++ {
		// read
		j := rand.Intn(n)
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			panic(fmt.Sprintf("[ERROR] %d'th read key=%d,expected_value=%s, actual_value=%s", i, k[j], v[j], val))
		}
	}
	for i := 0; i < n/4; i++ {
		val, ok := t.Search(k[i])
		if !ok || val != v[i] {
			panic(fmt.Sprintf("[ERROR] read key=%d,expected_value=%s, actual_value=%s", k[i], v[i], val))
		}
	}
	fmt.Println("===========> 2 test insert...")
	minK, _, _ := t.Min()
	maxK, _, _ := t.Max()

	n1k, _, _ := t.Nth(1)
	n2k, _, _ := t.Nth(n)
	if minK != n1k || maxK != n2k {
		panic("Nth error")
	}
	kk, vv := generateIntStr(maxK+1, n, !randKey, 10)
	for i := 0; i < n/2; i++ {
		kk[i] = minK - kk[i]
	}
	for i := 0; i < n; i++ {
		t.Insert(kk[i], vv[i])
	}
	if t.Len() != oldLen+n {
		panic(fmt.Sprintf("[ERROR] insert err size=%d,expected_size=%d", t.Len(), oldLen+n))
	}
	fmt.Printf("===========> insert %d date, now size=%d\n", n, t.Len())

	fmt.Println("===========> 3 test update...")
	for i := 0; i < n/2; i++ {
		j := rand.Intn(n)
		v[j] = seqStr("value-update-", i)
		t.Insert(k[j], v[j])
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			panic(fmt.Sprintf("[ERROR] update failure key=%d, expected_value=%s, but actual_value=%s", k[j], v[j], val))
		}
	}
	oldLen = t.Len()
	fmt.Printf("===========> update %d date, now size=%d\n", n/2, oldLen)

	fmt.Println("===========> 4 test delete1...")
	for i := 0; i < n/4; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			fmt.Printf("delete1 cannot del %d'th, key=%d\n", i, k[i])
			panic("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen-n/4 {
		panic(fmt.Sprintf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), oldLen-n/4))
	}
	oldLen = t.Len()
	fmt.Printf("===========> delete %d date, now size=%d\n", n/4, oldLen)

	fmt.Println("===========> 5 test delete2...")
	for i := n / 4; i < n; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			panic("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen+n/4-n {
		panic(fmt.Sprintf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), n))
	}
	fmt.Printf("===========> delete %d date, now size=%d\n", n-n/4, t.Len())

	fmt.Println("==========> test complete")
}

func testSGTRW(n int, randKey bool) {
	fmt.Println("===========> testSGTReadWrite")
	k, v := generateIntStr(0, n, randKey, 20)
	t := NewScapegoatTree[int, string](0.75, func(a, b int) int {
		if a > b {
			return 1
		} else if a == b {
			return 0
		}
		return -1
	})
	for i := 0; i < n; i++ {
		t.Insert(k[i], v[i])
	}
	oldLen := t.Len()
	fmt.Printf("==========> 0 init date, size=%d\n", oldLen)
	fmt.Println("=========> 1 test random read...")
	// random read
	for i := 0; i < n/2; i++ {
		// read
		j := rand.Intn(n)
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			panic(fmt.Sprintf("[ERROR] %d'th read key=%d,expected_value=%s, actual_value=%s", i, k[j], v[j], val))
		}
	}
	for i := 0; i < n/4; i++ {
		val, ok := t.Search(k[i])
		if !ok || val != v[i] {
			panic(fmt.Sprintf("[ERROR] read key=%d,expected_value=%s, actual_value=%s", k[i], v[i], val))
		}
	}
	fmt.Println("===========> 2 test insert...")
	minK, _, ok := t.Min()
	if !ok {
		panic("min error")
	}
	maxK, _, ok := t.Max()
	if !ok {
		panic("max error")
	}
	M := -1
	for _, x := range k {
		if x > M {
			M = x
		}
	}
	kk, vv := generateIntStr(maxK+1, n, !randKey, 10)
	for i := 0; i < n/2; i++ {
		kk[i] = minK - kk[i]
	}
	for i := 0; i < n; i++ {
		t.Insert(kk[i], vv[i])
	}
	if t.Len() != oldLen+n {
		panic(fmt.Sprintf("[ERROR] insert err size=%d,expected_size=%d", t.Len(), oldLen+n))
	}
	fmt.Printf("===========> insert %d date, now size=%d\n", n, t.Len())

	fmt.Println("===========> 3 test update...")
	for i := 0; i < n/2; i++ {
		j := rand.Intn(n)
		v[j] = seqStr("value-update-", i)
		t.Insert(k[j], v[j])
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			panic(fmt.Sprintf("[ERROR] update failure key=%d, expected_value=%s, but actual_value=%s", k[j], v[j], val))
		}
	}
	oldLen = t.Len()
	fmt.Printf("===========> update %d date, now size=%d\n", n/2, oldLen)

	fmt.Println("===========> 4 test delete1...")
	for i := 0; i < n/4; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			fmt.Printf("delete1 cannot del %d'th, key=%d\n", i, k[i])
			panic("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen-n/4 {
		panic(fmt.Sprintf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), oldLen-n/4))
	}
	oldLen = t.Len()
	fmt.Printf("===========> delete %d date, now size=%d\n", n/4, oldLen)

	fmt.Println("===========> 5 test delete2...")
	for i := n / 4; i < n; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			panic("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen+n/4-n {
		panic(fmt.Sprintf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), n))
	}
	fmt.Printf("===========> delete %d date, now size=%d\n", n-n/4, t.Len())

	fmt.Println("==========> test complete")
}

func TestBST(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "splay":
		testSplayTree(100)
	case "lct":
		testLCT()
	case "treap":
		testTreapRW(100, true)
	case "sgt":
		testSGTRW(100, true)
	default:
	}
}
