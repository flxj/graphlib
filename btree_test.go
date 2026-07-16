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
	"slices"
	"strconv"
	"testing"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStr(n int, s rand.Source) string {
	//s := rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func seqStr(perfix string, n int) string {
	return perfix + strconv.Itoa(n)
}

func generateIntStr(start, n int, randKey bool, maxGap int) ([]int, []string) {
	k := make([]int, n)
	v := make([]string, n)
	s := rand.NewSource(time.Now().UnixNano())

	prev := start
	for i := 0; i < n; i++ {
		if randKey {
			k[i] = prev + rand.Intn(maxGap) + 1
		} else {
			k[i] = prev + i
		}
		prev = k[i]
		v[i] = randStr(10, s)
	}
	rand.Shuffle(n, func(i, j int) {
		k[i], k[j] = k[j], k[i]
	})
	return k, v
}

func testBTreeReadWrite(n, degree int, randKey bool) {
	fmt.Println("===========> testBTreeReadWrite")
	k, v := generateIntStr(0, n, randKey, 20)
	t := NewBTree[int, string](&BTreeConfig{MinDegree: degree}, func(a, b int) int {
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
	fmt.Printf("==========> 0 init date, size=%d high=%d\n", oldLen, t.High())
	fmt.Println("=========> 1 test random read...")
	// random read
	for i := 0; i < n/2; i++ {
		// read
		j := rand.Intn(n)
		val, err := t.Search(k[j])
		if err != nil {
			panic(fmt.Sprintf("[ERROR] key=%d,err=%s", k[j], err.Error()))
		}
		if val != v[j] {
			panic(fmt.Sprintf("[ERROR] key=%d,expected_value=%s, actual_value=%s", k[j], v[j], val))
		}
	}

	fmt.Println("=========> 2 test scan...")
	var cnt int
	_ = t.Scan(false, func(a int, b string) error {
		cnt++
		return nil
	})
	if cnt != t.Len() {
		panic(fmt.Sprintf("[ERROR] scan count %d elements, but expected value=%d", cnt, t.Len()))
	}

	fmt.Println("===========> 3 test insert...")
	minK, _, _ := t.First()
	maxK, _, _ := t.Last()
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
	oldLen += n
	fmt.Printf("===========> insert %d date, now size=%d high=%d\n", n, t.Len(), t.High())

	fmt.Println("===========> 4 test update...")
	for i := 1; i <= n/2; i++ {
		j := rand.Intn(n)
		//oldV := v[j]
		v[j] = seqStr("value-update-", i)
		//fmt.Printf("update key=%d, old_value=%s, new_value=%s\n", k[j], oldV, v[j])

		t.Insert(k[j], v[j])

		val, err := t.Search(k[j])
		if err != nil {
			panic(fmt.Sprintf("[ERROR] %s", err.Error()))
		}
		if val != v[j] {
			panic(fmt.Sprintf("[ERROR] update failure key=%d, expected_value=%s, but actual_value=%s", k[j], v[j], val))
		}
	}

	fmt.Println("===========> 5 test delete1...")
	for i := 0; i < n/4; i++ {
		ok, err := t.Delete(k[i])
		if err != nil {
			panic(err.Error())
		}
		if !ok {
			panic("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen-n/4 {
		panic(fmt.Sprintf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), oldLen-n/4))
	}
	fmt.Printf("===========> delete %d date, now size=%d high=%d\n", n/4, t.Len(), t.High())

	fmt.Println("===========> 6 test delete2...")
	for i := n / 4; i < n; i++ {
		ok, err := t.Delete(k[i])
		if err != nil {
			panic(err.Error())
		}
		if !ok {
			panic("[ERROR] delete failure")
		}
	}
	if t.Len() != n {
		panic(fmt.Sprintf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), n))
	}
	fmt.Printf("===========> delete %d date, now size=%d high=%d\n", n-n/4, t.Len(), t.High())

	fmt.Println("==========> test complete")
}

func testBTreeCursor(n, degree int, randKey bool) {
	fmt.Printf("========> testBTreeCursor")
	ks, vs := generateIntStr(0, n, randKey, 5)
	t := NewBTree[int, string](&BTreeConfig{MinDegree: degree}, func(a, b int) int {
		if a > b {
			return 1
		} else if a == b {
			return 0
		}
		return -1
	})
	for i := 0; i < n; i++ {
		t.Insert(ks[i], vs[i])
	}
	fmt.Printf("============> insert date, size=%d high=%d\n", t.Len(), t.High())
	k := make([]int, n)
	copy(k, ks)
	slices.Sort(k)
	cur := t.Cursor()
	fmt.Println("============> 1 test asc...")
	kk, _, ok := cur.First()
	if !ok {
		panic("[ERROR] First() failed")
	}
	if kk != k[0] {
		panic(fmt.Sprintf("[ERROR] get first_key=%d, but expected_key=%d", kk, k[0]))
	}
	for i := 1; i < len(k) && cur.HasNext(); i++ {
		kk, _ = cur.Next()
		if kk != k[i] {
			panic(fmt.Sprintf("[ERROR] get the %d'th key, current_key=%d,expected_key=%d", i, kk, k[i]))
		}
	}

	_, _, _ = cur.First()
	if cur.HasPrev() {
		panic("cursor at first,cannot has prev")
	}

	fmt.Println("============> 2 test desc...")
	kk, _, ok = cur.Last()
	if !ok {
		panic("[ERROR] Last() failed")
	}
	if kk != k[n-1] {
		panic(fmt.Sprintf("[ERROR] get last_key=%d, but expected_key=%d", kk, k[n-1]))
	}
	for i := n - 2; i >= 0 && cur.HasPrev(); i-- {
		kk, _ = cur.Prev()
		if kk != k[i] {
			panic(fmt.Sprintf("[ERROR] get %d'th key,current_key=%d, expected-key=%d", i+1, kk, k[i]))
		}
	}

	_, _, _ = cur.Last()
	if cur.HasNext() {
		panic("cursor at last,cannot has next")
	}

	fmt.Println("============> 3 test seek+random...")
	for i := 0; i < n/4; i++ {
		j := rand.Intn(n)
		kk, vv, _ := cur.Seek(ks[j])
		if kk != ks[j] || vv != vs[j] {
			panic(fmt.Sprintf("[ERROR] seek current_key=%d, expected_key=%d", kk, ks[j]))
		}
	}

	fmt.Println("============> 4 test seek+forward...")
	kk, _, ok = cur.Seek(k[n/2])
	if !ok {
		panic(fmt.Sprintf("[ERROR] seek key %d failed", k[n/2]))
	}
	if kk != k[n/2] {
		panic(fmt.Sprintf("[ERROR] seek key %d but get %d", k[n/2], kk))
	}
	for i := n/2 + 1; i < n; i++ {
		kk, _ = cur.Next()
		if kk != k[i] {
			panic(fmt.Sprintf("[ERROR] seek current_key=%d, expected_key=%d", kk, k[i]))
		}
	}

	fmt.Println("============> 5 test seek+backward...")
	kk, _, ok = cur.Seek(k[n/2])
	if !ok {
		panic(fmt.Sprintf("[ERROR] seek key %d failed", k[n/2]))
	}
	if kk != k[n/2] {
		panic(fmt.Sprintf("[ERROR] seek key %d but get %d", k[n/2], kk))
	}
	for i := n/2 - 1; i >= 0; i-- {
		kk, _ = cur.Prev()
		//fmt.Printf("seek %d'th element,key=%d\n", i+1, kk)
		if kk != k[i] {
			panic(fmt.Sprintf("[ERROR] seek current_key=%d, expected_key=%d", kk, k[i]))
		}
	}

	fmt.Println("==========> test complete")
}

func TestBTree(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "rw":
		testBTreeReadWrite(100, 8, true)
	case "iter":
		testBTreeCursor(100, 16, true)
	default:
	}
}

func BenchmarkBTree(b *testing.B) {
	// TODO
}
