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

package art

import (
	"bytes"
	"flag"
	"math/rand"
	"sort"
	"testing"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randBytes(n int) []byte {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return b
}

func generateStrInt(n int, prefix []byte) ([][]byte, []int) {
	k := make([][]byte, n)
	v := make([]int, n)

	for i := 0; i < n; i++ {
		v[i] = rand.Intn(1000000007)

		num := rand.Intn(10) + 5
		kk := randBytes(num)
		if len(prefix) > 0 {
			b := make([]byte, len(prefix)+len(kk))
			copy(b, prefix)
			copy(b[len(prefix):], kk)
			k[i] = b
		} else {
			if num%2 == 1 {
				k[i] = append([]byte("abcdef"), kk...)
			} else {
				k[i] = kk
			}
		}
	}
	return k, v
}

func testARTRW(tt *testing.T, n int) {
	tt.Log("> test ART ReadWrite")
	k, v := generateStrInt(n, nil)
	t := NewTree[int](false)
	for i := 0; i < n; i++ {
		t.Insert(k[i], v[i])
	}
	oldLen := t.Len()
	tt.Logf("> 0 init date, size=%d\n", oldLen)

	tt.Log("> 1 test random read...")
	for i := 0; i < n/2; i++ {
		j := rand.Intn(n)
		//tt.Logf("i=%d search key=%s\n", i, string(k[j]))
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] key=%s,expected_value=%d, actual_value=%d", string(k[j]), v[j], val)
		}
	}
	for i := 0; i < n; i++ {
		//tt.Logf("i=%d search2 key=%s\n", i, string(k[i]))
		val, ok := t.Search(k[i])
		if !ok || val != v[i] {
			tt.Errorf("[ERROR] key=%s,expected_value=%d, actual_value=%d", string(k[i]), v[i], val)
		}
	}

	tt.Log("> 2 test scan...")
	var cnt int
	fn := func(k Key, v int) error {
		cnt++
		return nil
	}
	t.Scan(fn, false)
	if cnt != t.Len() {
		tt.Errorf("[ERROR] scan count %d elements, but expected value=%d", cnt, t.Len())
	}

	tt.Log("> 3 test insert...")
	kk, vv := generateStrInt(n, []byte("7654321"))
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
		v[j] = rand.Intn(100)
		t.Insert(k[j], v[j])
		val, ok := t.Search(k[j])
		if !ok || val != v[j] {
			tt.Errorf("[ERROR] update failure key=%s, expected_value=%d, but actual_value=%d", string(k[j]), v[j], val)
		}
	}

	tt.Log("> 5 test scan...")
	cnt = 0
	t.PrefixScan([]byte("7654321"), fn, false)
	if cnt != n {
		tt.Errorf("[ERROR] scan count %d elements, but expected value=%d", cnt, n)
	}

	tt.Log("> 6 test delete1...")
	for i := 0; i < n/4; i++ {
		val, ok := t.Delete(k[i])
		if !ok || val != v[i] {
			tt.Fatalf("[Err] delete key='%s' failed, expected_value=%d, but actual_value=%d", string(k[i]), v[i], val)
		}
	}
	if t.Len() != oldLen-n/4 {
		tt.Errorf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), oldLen-n/4)
	}
	tt.Logf("> delete %d date, now size=%d\n", n/4, t.Len())

	tt.Log("> 7 test delete2...")
	for i := n / 4; i < n; i++ {
		val, ok := t.Delete(kk[i])
		if !ok || val != vv[i] {
			tt.Fatalf("[Err] delete key='%s' failed, expected_value=%d, but actual_value=%d", string(kk[i]), vv[i], val)
		}
	}
	if t.Len() != n {
		tt.Errorf("[ERROR] after delete size=%d, expected_size=%d", t.Len(), n)
	}
	tt.Logf("> delete %d date, now size=%d\n", n-n/4, t.Len())

	tt.Log("> test complete")
}

func testARTCursor(tt *testing.T, n int) {
	tt.Log("> test ART Cursor")
	ks, vs := generateStrInt(n, nil)
	t := NewTree[int](false)
	for i := 0; i < n; i++ {
		t.Insert(ks[i], vs[i])
	}
	tt.Logf("> insert date, size=%d\n", t.Len())
	k := make([][]byte, n)
	copy(k, ks)
	sort.Slice(k, func(i, j int) bool {
		for p := 0; p < min(len(k[i]), len(k[j])); p++ {
			if k[i][p] < k[j][p] {
				return true
			} else if k[i][p] > k[j][p] {
				return false
			}
		}
		return len(k[i]) < len(k[j])
	})

	cur := t.Cursor()
	tt.Log("> 1 test asc...")
	kk, _, ok := cur.First()
	if !ok {
		tt.Fatal("[ERROR] First() failed")
	}
	if cur.HasPrev() {
		tt.Fatal("cursor at first,cannot has prev")
	}
	if !bytes.Equal(kk, k[0]) {
		tt.Errorf("[ERROR] get first_key=%s, but expected_key=%s", string(kk), string(k[0]))
	}
	for i := 1; i < len(k) && cur.HasNext(); i++ {
		kk, _ = cur.Next()
		if !bytes.Equal(kk, k[i]) {
			tt.Errorf("[ERROR] get the %d'th key, current_key=%s,expected_key=%s", i, string(kk), string(k[i]))
		}
	}

	tt.Log("> 2 test desc...")
	kk, _, ok = cur.Last()
	if !ok {
		tt.Fatal("[ERROR] Last() failed")
	}
	if cur.HasNext() {
		tt.Fatal("cursor at last,cannot has next")
	}
	if !bytes.Equal(kk, k[n-1]) {
		tt.Errorf("[ERROR] get last_key=%s, but expected_key=%s", string(kk), string(k[n-1]))
	}
	for i := n - 2; i >= 0 && cur.HasPrev(); i-- {
		kk, _ = cur.Prev()
		if !bytes.Equal(kk, k[i]) {
			tt.Errorf("[ERROR] get %d'th key,current_key=%s, expected-key=%s", i+1, string(kk), string(k[i]))
		}
	}

	tt.Log("> 3 test seek+forward...")
	kk, _, ok = cur.Seek(k[n/2])
	if !ok {
		tt.Errorf("[ERROR] seek key %s failed", string(k[n/2]))
	}
	if !bytes.Equal(kk, k[n/2]) {
		tt.Errorf("[ERROR] seek key %s but get %s", string(k[n/2]), string(kk))
	}
	for i := n/2 + 1; i < n; i++ {
		kk, _ = cur.Next()
		if !bytes.Equal(kk, k[i]) {
			tt.Errorf("[ERROR] seek current_key=%s, expected_key=%s", string(kk), string(k[i]))
		}
	}

	tt.Log("> 5 test seek+backward...")
	kk, _, ok = cur.Seek(k[n/2+1])
	if !ok {
		tt.Errorf("[ERROR] seek key %d failed", k[n/2])
	}
	if !bytes.Equal(kk, k[n/2+1]) {
		tt.Errorf("[ERROR] seek key %s but get %s", string(k[n/2+1]), string(kk))
	}
	for i := n / 2; i >= 0; i-- {
		kk, _ = cur.Prev()
		if !bytes.Equal(kk, k[i]) {
			tt.Errorf("[ERROR] seek current_key=%s, expected_key=%s", string(kk), string(k[i]))
		}
	}

	tt.Log("> test complete")
}

func TestART(t *testing.T) {
	args := flag.Args()
	switch args[0] {
	case "rw":
		testARTRW(t, 100)
	case "iter":
		testARTCursor(t, 20)
	default:
	}
}
