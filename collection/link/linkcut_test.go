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

package link

import (
	"strconv"
	"testing"
)

func TestLCT(tt *testing.T) {
	tt.Log("> test link-cut-tree")
	t := NewLinkCutTree[int, string, int](func(a, b int) int {
		if a > b {
			return 1
		} else if a == b {
			return 0
		} else {
			return -1
		}
	})
	tt.Log("> 1. build tree T")
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
		tt.Error("add node error")
	}
	edges := [][2]int{
		{0, 1}, {0, 2}, {0, 3}, {0, 4}, {1, 5}, {1, 6}, {2, 7}, {3, 8},
		{3, 9}, {3, 10}, {4, 11}, {7, 12}, {11, 13}, {11, 14}, {11, 15},
	}
	for i, e := range edges {
		tt.Logf("add edge (%d,%d)\n", e[0], e[1])
		if !t.Link(e[0], e[1]) {
			tt.Error("link error")
		}
		if t.Component() != n-i-1 {
			tt.Logf("component %d, n-i-i%d\n", t.Component(), n-i-1)
			tt.Error("component error")
		}
	}
	tt.Log("> 2. test Connected")
	for _, e := range [][2]int{{1, 12}, {3, 8}, {4, 7}, {6, 10}, {2, 4}} {
		if !t.Connected(e[0], e[1]) {
			tt.Error("connected error")
		}
	}

	tt.Log("> 3. test MakeRoot")
	if !t.MakeRoot(3) {
		tt.Error("makeRoot 3 error")
	}
	if !t.IsRoot(3) {
		tt.Error("isRoot error")
	}
	for _, v := range []int{15, 5, 2, 4} {
		r, _, ok := t.FindRoot(v)
		if !ok || r != 3 {
			tt.Log("v", v, " findRoot", r)
			tt.Error("findRoot error")
		}
	}
	if !t.MakeRoot(0) {
		tt.Error("makeRoot 0 error")
	}
	tt.Log("> 4. test PathSum")
	s := [][2]int{{12, 21}, {6, 7}, {7, 9}, {8, 11}, {11, 15}, {14, 29}}
	for _, p := range s {
		r, _ := t.PathSum(p[0])
		if r != p[1] {
			tt.Log("v", p[0], " sum", r)
			tt.Error("pathAgg error")
		}
	}
	tt.Log("> 5. test PathPathAggregate")
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
			tt.Log("key", p[0], " res", p[1], " but get", pp)
			tt.Error("pathAgg error")
		}
	}
	tt.Log("> 6. test MakeTree")
	//   17------16-----18-----21
	//                   | \
	//                   19 20
	//
	n2 := 6
	for i := 0; i < n2; i++ {
		if !t.MakeTree(i+n, "val-"+strconv.Itoa(i+n), i+n) {
			tt.Error("makeTree error")
		}
	}
	edges2 := [][2]int{{16, 17}, {18, 19}, {18, 20}, {18, 21}, {18, 16}}
	for _, e := range edges2 {
		tt.Logf("add edge (%d,%d)\n", e[0], e[1])
		if !t.Link(e[0], e[1]) {
			tt.Error("link error 2")
		}
	}
	if t.Component() != 2 {
		tt.Error("component error 2")
	}
	tt.Log("> 7. test Cut")
	for _, v := range []int{7, 3, 11, 16} {
		if !t.Cut(v) {
			tt.Error("cut error")
		}
	}
	if t.Component() != 6 {
		tt.Logf("components %d, but expeect 6\n", t.Component())
		tt.Error("cut component error")
	}
	if t.Connected(6, 11) || !t.Connected(5, 2) {
		tt.Error("cut connected error")
	}
	tt.Log("> 8. test Link")
	if !t.Link(18, 2) {
		tt.Error("link 18-2 error")
	}
	if !t.Link(11, 0) {
		tt.Error("link 11-0 error")
	}
	if t.Component() != 4 {
		tt.Error("link component error 2")
	}
	if !t.Connected(6, 21) || t.Connected(8, 13) {
		tt.Error("link connected error 3")
	}
	tt.Log("> test pass")
}
