/*
	Copyright (C) 2023 flxj(https:=//github.com/flxj)

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http:=//www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package tree

import (
	"flag"
	"math/rand"
	"testing"
)

func testSplayTree(n int, tt *testing.T) {
	tt.Log("> test splay tree")
	t := NewSplayTree[int, string](func(a, b int) int {
		if a > b {
			return 1
		} else if a == b {
			return 0
		}
		return -1
	})
	tt.Log("> 0. init data")
	keys, vals := generateIntStr(0, n, true, 5)
	for i, k := range keys {
		t.Insert(k, vals[i])
	}
	if t.Len() != n {
		tt.Errorf("init error,expect %d data,but actual %d ", n, t.Len())
	}
	tt.Log("> 1. test read")
	for i := 0; i < n/2; i++ {
		j := rand.Intn(n)
		v, ok := t.Search(keys[j])
		if !ok || v != vals[j] {
			tt.Errorf("i:=%d (key:=%d,val:=%s) but get (ok%v,v%s)\n", i, keys[j], vals[j], ok, v)
		}
	}
	tt.Log("> 2. test insert")
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
		tt.Errorf("insert error,expect %d data,but actual %d ", 2*n, t.Len())
	}
	tt.Log("> 3. test update")
	v := "value"
	for i := 0; i < n/2; i++ {
		k := keys2[rand.Intn(n)]
		t.Insert(k, v)
		vv, ok := t.Search(k)
		if !ok || v != vv {
			tt.Errorf("update erroe,key:=%d, expect value:=%s but actual %s", k, v, vv)
		}
	}
	tt.Log("> 4. test delete")
	for i := 0; i < n; i++ {
		v, ok := t.Delete(keys[i])
		if !ok || v != vals[i] {
			tt.Errorf("delete erroe,key:=%d", keys[i])
		}
	}
	if t.Len() != n {
		tt.Errorf("delete error,expect %d data,but actual %d ", n, t.Len())
	}
}

func testTreapRW(n int, randKey bool, tt *testing.T) {
	tt.Log("> test Treap ReadWrite")
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
	tt.Logf("> 0 init date, size%d\n", oldLen)
	tt.Log("> 1 test random read...")
	// random read
	for i := 0; i < n/2; i++ {
		// read
		j := rand.Intn(n)
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] %d'th read key%d,expected_value%s, actual_value%s", i, k[j], v[j], val)
		}
	}
	for i := 0; i < n/4; i++ {
		val, ok := t.Search(k[i])
		if !ok || val != v[i] {
			tt.Errorf("[ERROR] read key%d,expected_value%s, actual_value%s", k[i], v[i], val)
		}
	}
	tt.Log("> 2 test insert...")
	minK, _, _ := t.Min()
	maxK, _, _ := t.Max()

	n1k, _, _ := t.Nth(1)
	n2k, _, _ := t.Nth(n)
	if minK != n1k || maxK != n2k {
		tt.Error("Nth error")
	}
	kk, vv := generateIntStr(maxK+1, n, !randKey, 10)
	for i := 0; i < n/2; i++ {
		kk[i] = minK - kk[i]
	}
	for i := 0; i < n; i++ {
		t.Insert(kk[i], vv[i])
	}
	if t.Len() != oldLen+n {
		tt.Errorf("[ERROR] insert err size%d,expected_size%d", t.Len(), oldLen+n)
	}
	tt.Logf("> insert %d date, now size%d\n", n, t.Len())

	tt.Log("> 3 test update...")
	for i := 0; i < n/2; i++ {
		j := rand.Intn(n)
		v[j] = seqStr("value-update-", i)
		t.Insert(k[j], v[j])
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] update failure key%d, expected_value%s, but actual_value%s", k[j], v[j], val)
		}
	}
	oldLen = t.Len()
	tt.Logf("> update %d date, now size%d\n", n/2, oldLen)

	tt.Log("> 4 test delete1...")
	for i := 0; i < n/4; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			tt.Logf("delete1 cannot del %d'th, key%d\n", i, k[i])
			tt.Error("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen-n/4 {
		tt.Errorf("[ERROR] after delete size%d, expected_size%d", t.Len(), oldLen-n/4)
	}
	oldLen = t.Len()
	tt.Logf("> delete %d date, now size%d\n", n/4, oldLen)

	tt.Log("> 5 test delete2...")
	for i := n / 4; i < n; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			tt.Error("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen+n/4-n {
		tt.Errorf("[ERROR] after delete size%d, expected_size%d", t.Len(), n)
	}
	tt.Logf("> delete %d date, now size%d\n", n-n/4, t.Len())

	tt.Log("> test complete")
}

func testSGTRW(n int, randKey bool, tt *testing.T) {
	tt.Log("> test SGT ReadWrite")
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
	tt.Logf("> 0 init date, size%d\n", oldLen)
	tt.Log("> 1 test random read...")
	// random read
	for i := 0; i < n/2; i++ {
		// read
		j := rand.Intn(n)
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] %d'th read key%d,expected_value%s, actual_value%s", i, k[j], v[j], val)
		}
	}
	for i := 0; i < n/4; i++ {
		val, ok := t.Search(k[i])
		if !ok || val != v[i] {
			tt.Errorf("[ERROR] read key%d,expected_value%s, actual_value%s", k[i], v[i], val)
		}
	}
	tt.Log("> 2 test insert...")
	minK, _, ok := t.Min()
	if !ok {
		tt.Error("min error")
	}
	maxK, _, ok := t.Max()
	if !ok {
		tt.Error("max error")
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
		tt.Errorf("[ERROR] insert err size%d,expected_size%d", t.Len(), oldLen+n)
	}
	tt.Logf("> insert %d date, now size%d\n", n, t.Len())

	tt.Log("> 3 test update...")
	for i := 0; i < n/2; i++ {
		j := rand.Intn(n)
		v[j] = seqStr("value-update-", i)
		t.Insert(k[j], v[j])
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] update failure key%d, expected_value%s, but actual_value%s", k[j], v[j], val)
		}
	}
	oldLen = t.Len()
	tt.Logf("> update %d date, now size%d\n", n/2, oldLen)

	tt.Log("> 4 test delete1...")
	for i := 0; i < n/4; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			tt.Logf("delete1 cannot del %d'th, key%d\n", i, k[i])
			tt.Error("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen-n/4 {
		tt.Errorf("[ERROR] after delete size%d, expected_size%d", t.Len(), oldLen-n/4)
	}
	oldLen = t.Len()
	tt.Logf("> delete %d date, now size%d\n", n/4, oldLen)

	tt.Log("> 5 test delete2...")
	for i := n / 4; i < n; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			tt.Error("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen+n/4-n {
		tt.Errorf("[ERROR] after delete size%d, expected_size%d", t.Len(), n)
	}
	tt.Logf("> delete %d date, now size%d\n", n-n/4, t.Len())

	tt.Log("> test complete")
}

func testRedBlackTreeRW(n int, randKey bool, tt *testing.T) {
	tt.Log("> test RedBlackTree ReadWrite")
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
	tt.Logf("> 0 init date, size%d\n", oldLen)
	tt.Log("> 1 test random read...")
	// random read
	for i := 0; i < n/2; i++ {
		// read
		j := rand.Intn(n)
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] %d'th read key%d,expected_value%s, actual_value%s", i, k[j], v[j], val)
		}
	}
	for i := 0; i < n/4; i++ {
		val, ok := t.Search(k[i])
		if !ok || val != v[i] {
			tt.Errorf("[ERROR] read key%d,expected_value%s, actual_value%s", k[i], v[i], val)
		}
	}
	tt.Log("> 2 test insert...")
	minK, _, ok := t.Min()
	if !ok {
		tt.Error("min error")
	}
	maxK, _, ok := t.Max()
	if !ok {
		tt.Error("max error")
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
		tt.Errorf("[ERROR] insert err size%d,expected_size%d", t.Len(), oldLen+n)
	}
	tt.Logf("> insert %d date, now size%d\n", n, t.Len())

	tt.Log("> 3 test update...")
	for i := 0; i < n/2; i++ {
		j := rand.Intn(n)
		v[j] = seqStr("value-update-", i)
		t.Insert(k[j], v[j])
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] update failure key%d, expected_value%s, but actual_value%s", k[j], v[j], val)
		}
	}
	oldLen = t.Len()
	tt.Logf("> update %d date, now size%d\n", n/2, oldLen)

	tt.Log("> 4 test delete1...")
	for i := 0; i < n/4; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			tt.Logf("delete1 cannot del %d'th, key%d\n", i, k[i])
			tt.Error("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen-n/4 {
		tt.Errorf("[ERROR] after delete size%d, expected_size%d", t.Len(), oldLen-n/4)
	}
	oldLen = t.Len()
	tt.Logf("> delete %d date, now size%d\n", n/4, oldLen)

	tt.Log("> 5 test delete2...")
	for i := n / 4; i < n; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			tt.Error("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen+n/4-n {
		tt.Errorf("[ERROR] after delete size%d, expected_size%d", t.Len(), n)
	}
	tt.Logf("> delete %d date, now size%d\n", n-n/4, t.Len())

	tt.Log("> test complete")
}

func TestBST(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "splay":
		testSplayTree(100, t)
	case "treap":
		testTreapRW(100, true, t)
	case "sgt":
		testSGTRW(100, true, t)
	case "rbt":
		testRedBlackTreeRW(100, true, t)
	default:
	}
}
