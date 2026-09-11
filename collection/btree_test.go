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

package collection

import (
	"flag"
	"math/rand"
	"slices"
	"strconv"
	"testing"
	"time"
)

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

func generateIntArr(start, n int, randKey bool, maxGap int) []int {
	k := make([]int, n)
	prev := start
	for i := 0; i < n; i++ {
		if randKey {
			k[i] = prev + rand.Intn(maxGap) + 1
		} else {
			k[i] = prev + i
		}
		prev = k[i]
	}
	rand.Shuffle(n, func(i, j int) {
		k[i], k[j] = k[j], k[i]
	})
	return k
}

func testBTreeReadWrite(tt *testing.T, n, degree int, randKey bool) {
	tt.Log("> testBTreeReadWrite")
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
	tt.Logf("> 0 init date, size=%d high=%d\n", oldLen, t.High())
	tt.Log("> 1 test random read...")
	// random read
	for i := 0; i < n/2; i++ {
		// read
		j := rand.Intn(n)
		val, err := t.Search(k[j])
		if err != nil {
			tt.Errorf("[ERROR] key=%d,err=%s", k[j], err.Error())
		}
		if val != v[j] {
			tt.Errorf("[ERROR] key=%d,expected_value=%s, actual_value=%s", k[j], v[j], val)
		}
	}

	tt.Log("> 2 test scan...")
	var cnt int
	_ = t.Scan(false, func(a int, b string) error {
		cnt++
		return nil
	})
	if cnt != t.Len() {
		tt.Errorf("[ERROR] scan count %d elements, but expected value=%d", cnt, t.Len())
	}

	tt.Log("> 3 test insert...")
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
		tt.Errorf("[ERROR] insert err size=%d,expected_size=%d", t.Len(), oldLen+n)
	}
	oldLen += n
	tt.Logf("> insert %d date, now size=%d high=%d\n", n, t.Len(), t.High())

	tt.Log("> 4 test update...")
	for i := 1; i <= n/2; i++ {
		j := rand.Intn(n)
		//oldV := v[j]
		v[j] = seqStr("value-update-", i)
		//tt.Logf("update key=%d, old_value=%s, new_value=%s\n", k[j], oldV, v[j])

		t.Insert(k[j], v[j])

		val, err := t.Search(k[j])
		if err != nil {
			tt.Errorf("[ERROR] %s", err.Error())
		}
		if val != v[j] {
			tt.Errorf("[ERROR] update failure key=%d, expected_value=%s, but actual_value=%s", k[j], v[j], val)
		}
	}

	tt.Log("> 5 test delete1...")
	for i := 0; i < n/4; i++ {
		ok, err := t.Delete(k[i])
		if err != nil {
			tt.Fatal(err.Error())
		}
		if !ok {
			tt.Fatal("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen-n/4 {
		tt.Errorf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), oldLen-n/4)
	}
	tt.Logf("> delete %d date, now size=%d high=%d\n", n/4, t.Len(), t.High())

	tt.Log("> 6 test delete2...")
	for i := n / 4; i < n; i++ {
		ok, err := t.Delete(k[i])
		if err != nil {
			tt.Fatal(err.Error())
		}
		if !ok {
			tt.Fatal("[ERROR] delete failure")
		}
	}
	if t.Len() != n {
		tt.Errorf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), n)
	}
	tt.Logf("> delete %d date, now size=%d high=%d\n", n-n/4, t.Len(), t.High())

	tt.Log("> test complete")
}

func testBTreeCursor(tt *testing.T, n, degree int, randKey bool) {
	tt.Logf("=======> testBTreeCursor")
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
	tt.Logf("=> insert date, size=%d high=%d\n", t.Len(), t.High())
	k := make([]int, n)
	copy(k, ks)
	slices.Sort(k)
	cur := t.Cursor()
	tt.Log("=> 1 test asc...")
	kk, _, ok := cur.First()
	if !ok {
		tt.Fatal("[ERROR] First() failed")
	}
	if kk != k[0] {
		tt.Errorf("[ERROR] get first_key=%d, but expected_key=%d", kk, k[0])
	}
	for i := 1; i < len(k) && cur.HasNext(); i++ {
		kk, _ = cur.Next()
		if kk != k[i] {
			tt.Errorf("[ERROR] get the %d'th key, current_key=%d,expected_key=%d", i, kk, k[i])
		}
	}

	_, _, _ = cur.First()
	if cur.HasPrev() {
		tt.Fatal("cursor at first,cannot has prev")
	}

	tt.Log("=> 2 test desc...")
	kk, _, ok = cur.Last()
	if !ok {
		tt.Fatal("[ERROR] Last() failed")
	}
	if kk != k[n-1] {
		tt.Errorf("[ERROR] get last_key=%d, but expected_key=%d", kk, k[n-1])
	}
	for i := n - 2; i >= 0 && cur.HasPrev(); i-- {
		kk, _ = cur.Prev()
		if kk != k[i] {
			tt.Errorf("[ERROR] get %d'th key,current_key=%d, expected-key=%d", i+1, kk, k[i])
		}
	}

	_, _, _ = cur.Last()
	if cur.HasNext() {
		tt.Fatal("cursor at last,cannot has next")
	}

	tt.Log("=> 3 test seek+random...")
	for i := 0; i < n/4; i++ {
		j := rand.Intn(n)
		kk, vv, _ := cur.Seek(ks[j])
		if kk != ks[j] || vv != vs[j] {
			tt.Errorf("[ERROR] seek current_key=%d, expected_key=%d", kk, ks[j])
		}
	}

	tt.Log("=> 4 test seek+forward...")
	kk, _, ok = cur.Seek(k[n/2])
	if !ok {
		tt.Errorf("[ERROR] seek key %d failed", k[n/2])
	}
	if kk != k[n/2] {
		tt.Errorf("[ERROR] seek key %d but get %d", k[n/2], kk)
	}
	for i := n/2 + 1; i < n; i++ {
		kk, _ = cur.Next()
		if kk != k[i] {
			tt.Errorf("[ERROR] seek current_key=%d, expected_key=%d", kk, k[i])
		}
	}

	tt.Log("=> 5 test seek+backward...")
	kk, _, ok = cur.Seek(k[n/2])
	if !ok {
		tt.Errorf("[ERROR] seek key %d failed", k[n/2])
	}
	if kk != k[n/2] {
		tt.Errorf("[ERROR] seek key %d but get %d", k[n/2], kk)
	}
	for i := n/2 - 1; i >= 0; i-- {
		kk, _ = cur.Prev()
		//tt.Logf("seek %d'th element,key=%d\n", i+1, kk)
		if kk != k[i] {
			tt.Errorf("[ERROR] seek current_key=%d, expected_key=%d", kk, k[i])
		}
	}

	tt.Log("> test complete")
}

func TestBTree(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "rw":
		testBTreeReadWrite(t, 100, 8, true)
	case "iter":
		testBTreeCursor(t, 100, 16, true)
	default:
		testBTreeReadWrite(t, 100, 8, true)
	}
}

func BenchmarkBTree(b *testing.B) {
	// TODO
}

func testSkipListRW(tt *testing.T, n int, randKey bool) {
	tt.Log("> testSkipLisReadWrite")
	k, v := generateIntStr(0, n, randKey, 20)
	t := NewSkipList[int, string](&SkipListConfig{}, func(a, b int) int {
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
	tt.Logf("> 0 init date, size=%d\n", oldLen)
	tt.Log("> 1 test random read...")
	// random read
	for i := 0; i < n/2; i++ {
		// read
		j := rand.Intn(n)
		_, val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] key=%d,expected_value=%s, actual_value=%s", k[j], v[j], val)
		}
	}

	tt.Log("> 2 test scan...")
	var cnt int
	_ = t.Scan(func(a int, b string) error {
		cnt++
		return nil
	})
	if cnt != t.Len() {
		tt.Errorf("[ERROR] scan count %d elements, but expected value=%d", cnt, t.Len())
	}

	tt.Log("> 3 test insert...")
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
		tt.Errorf("[ERROR] insert err size=%d,expected_size=%d", t.Len(), oldLen+n)
	}
	oldLen += n
	tt.Logf("> insert %d date, now size=%d\n", n, t.Len())

	tt.Log("> 4 test update...")
	for i := 1; i <= n/2; i++ {
		j := rand.Intn(n)
		v[j] = seqStr("value-update-", i)
		if ok := t.Update(k[j], v[j]); !ok {
			tt.Fatal("update failure")
		}
		_, val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] update failure key=%d, expected_value=%s, but actual_value=%s", k[j], v[j], val)
		}
	}

	tt.Log("> 5 test delete1...")
	for i := 0; i < n/4; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			tt.Fatal("[ERROR] delete failure")
		}
	}
	if t.Len() != oldLen-n/4 {
		tt.Errorf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), oldLen-n/4)
	}
	tt.Logf("> delete %d date, now size=%d\n", n/4, t.Len())

	tt.Log("> 6 test delete2...")
	for i := n / 4; i < n; i++ {
		_, ok := t.Delete(k[i])
		if !ok {
			tt.Fatal("[ERROR] delete failure")
		}
	}
	if t.Len() != n {
		tt.Errorf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), n)
	}
	tt.Logf("> delete %d date, now size=%d\n", n-n/4, t.Len())

	tt.Log("> test complete")
}

func testSkipListCursor(tt *testing.T, n int, randKey bool) {
	tt.Logf("=======> testSkiplistCursor")
	ks, vs := generateIntStr(0, n, randKey, 5)
	t := NewSkipList[int, string](&SkipListConfig{}, func(a, b int) int {
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
	tt.Logf("=> insert date, size=%d\n", t.Len())
	k := make([]int, n)
	copy(k, ks)
	slices.Sort(k)
	cur := t.Cursor()
	tt.Log("=> 1 test asc...")
	kk, _, ok := cur.First()
	if !ok {
		tt.Fatal("[ERROR] First() failed")
	}
	if kk != k[0] {
		tt.Errorf("[ERROR] get first_key=%d, but expected_key=%d", kk, k[0])
	}
	for i := 1; i < len(k) && cur.HasNext(); i++ {
		kk, _ = cur.Next()
		if kk != k[i] {
			tt.Errorf("[ERROR] get the %d'th key, current_key=%d,expected_key=%d", i, kk, k[i])
		}
	}

	_, _, _ = cur.First()
	if cur.HasPrev() {
		tt.Fatal("cursor at first,cannot has prev")
	}

	tt.Log("=> 2 test desc...")
	kk, _, ok = cur.Last()
	if !ok {
		tt.Fatal("[ERROR] Last() failed")
	}
	if kk != k[n-1] {
		tt.Errorf("[ERROR] get last_key=%d, but expected_key=%d", kk, k[n-1])
	}
	for i := n - 2; i >= 0 && cur.HasPrev(); i-- {
		kk, _ = cur.Prev()
		if kk != k[i] {
			tt.Errorf("[ERROR] get %d'th key,current_key=%d, expected-key=%d", i+1, kk, k[i])
		}
	}

	_, _, _ = cur.Last()
	if cur.HasNext() {
		tt.Fatal("cursor at last,cannot has next")
	}

	tt.Log("=> 3 test seek+random...")
	for i := 0; i < n/4; i++ {
		j := rand.Intn(n)
		kk, vv, _ := cur.Seek(ks[j])
		if kk != ks[j] || vv != vs[j] {
			tt.Errorf("[ERROR] seek current_key=%d, expected_key=%d", kk, ks[j])
		}
	}

	tt.Log("=> 4 test seek+forward...")
	kk, _, ok = cur.Seek(k[n/2])
	if !ok {
		tt.Errorf("[ERROR] seek key %d failed", k[n/2])
	}
	if kk != k[n/2] {
		tt.Errorf("[ERROR] seek key %d but get %d", k[n/2], kk)
	}
	for i := n/2 + 1; i < n; i++ {
		kk, _ = cur.Next()
		if kk != k[i] {
			tt.Errorf("[ERROR] seek current_key=%d, expected_key=%d", kk, k[i])
		}
	}

	tt.Log("=> 5 test seek+backward...")
	kk, _, ok = cur.Seek(k[n/2])
	if !ok {
		tt.Errorf("[ERROR] seek key %d failed", k[n/2])
	}
	if kk != k[n/2] {
		tt.Errorf("[ERROR] seek key %d but get %d", k[n/2], kk)
	}
	for i := n/2 - 1; i >= 0; i-- {
		kk, _ = cur.Prev()
		//tt.Logf("seek %d'th element,key=%d\n", i+1, kk)
		if kk != k[i] {
			tt.Errorf("[ERROR] seek current_key=%d, expected_key=%d", kk, k[i])
		}
	}

	tt.Log("> test complete")
}

func TestSkipList(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "rw":
		testSkipListRW(t, 100, true)
	case "iter":
		testSkipListCursor(t, 100, true)
	default:
		testSkipListRW(t, 100, true)
	}
}

func testRTreeRW(tt *testing.T, n int) {
	tt.Log("> testRTreeReadWrite")
	t := NewRTree[string](32, false, BoxDist[int])
	r1 := generateIntArr(0, n, true, 10)
	r2 := generateIntArr(0, n, true, 5)
	r3 := generateIntArr(n, n, true, 20)
	r4, d := generateIntStr(n, n, true, 30)
	tt.Log("> insert...")
	for i := 0; i < n; i++ {
		if r3[i] < r1[i] {
			r3[i], r1[i] = r1[i], r3[i]
		}
		if r4[i] < r2[i] {
			r4[i], r2[i] = r2[i], r4[i]
		}
		rect := Rectangle[int]{r1[i], r2[i], r3[i], r4[i]}
		t.Insert(d[i], rect)
	}
	oldLen := t.Len()
	tt.Logf("> 0 init date, size=%d\n", oldLen)
	tt.Log("> 1 test random read...")
	// random read
	for i := 0; i < n/2; i++ {
		// read
		j := rand.Intn(n)
		rect := Rectangle[int]{r1[j], r2[j], r3[j], r4[j]}
		rec, val := t.Search(rect)
		if len(rec) == 0 {
			tt.Log("rect=", rect)
			tt.Log("data=", d[j])
			tt.Log(rec)
			tt.Log(val)
			tt.Fatal("read error")
		}
	}

	tt.Log("> 2 test scan...")
	var cnt int
	fn := func(_ Rectangle[int], _ string) error {
		cnt++
		return nil
	}
	t.Scan(fn)
	if cnt != t.Len() {
		tt.Errorf("[ERROR] scan count %d elements, but expected value=%d", cnt, t.Len())
	}

	tt.Log("> 3 test delete1...")
	for i := 0; i < n/4; i++ {
		rect := Rectangle[int]{r1[i], r2[i], r3[i], r4[i]}
		t.Delete(rect)
	}
	if t.Len() != oldLen-n/4 {
		tt.Errorf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), oldLen-n/4)
	}

	tt.Log("> 6 test delete2...")
	for i := n / 4; i < n; i++ {
		rect := Rectangle[int]{r1[i], r2[i], r3[i], r4[i]}
		t.Delete(rect)
	}
	if t.Len() != 0 {
		tt.Errorf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), 0)
	}
	cnt = 0
	t.Scan(fn)
	if cnt != t.Len() {
		tt.Errorf("[ERROR] scan count %d elements, but expected value=%d", cnt, t.Len())
	}
	tt.Log("> test complete")
}

func TestRTree(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "r":
		testRTreeRW(t, 100)
	case "r+":
	case "r*":
	default:
		testRTreeRW(t, 100)
	}
}
