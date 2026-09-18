package rtree

import (
	"math/rand"
	"testing"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStr(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
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

func generateIntStr(start, n int, randKey bool, maxGap int) ([]int, []string) {
	k := make([]int, n)
	v := make([]string, n)

	prev := start
	for i := 0; i < n; i++ {
		if randKey {
			k[i] = prev + rand.Intn(maxGap) + 1
		} else {
			k[i] = prev + i
		}
		prev = k[i]
		v[i] = randStr(10)
	}
	rand.Shuffle(n, func(i, j int) {
		k[i], k[j] = k[j], k[i]
	})
	return k, v
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
