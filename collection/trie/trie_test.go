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

package trie

import (
	"math/rand"
	"testing"
)

func testTrieRW(tt *testing.T, n int) {
	str := []string{
		"a",
		"b",
		"abc",
		"abcd",
		"abcdefg",
		"bcde",
		"bcdef",
		"bcdefghi",
		"bgty",
		"bgtyhn",
		"nhy",
		"nhyujm", //12
		"abcdedtgb",
		"abcdolp",
		"abcdyhnvv",
	}
	t := &Trie[int]{}
	for i, s := range str {
		t.Insert([]byte(s), i)
	}
	tt.Log("trie len", t.Len()) // 15

	for i := 0; i < 10; i++ {
		j := rand.Intn(len(str))
		if _, ok := t.Search([]byte(str[j])); !ok {
			tt.Error("search error")
		}
	}
	for _, s := range []string{"abcdxxx", "bcdehyu", "nhyj", "qaz"} {
		if _, ok := t.Search([]byte(s)); ok {
			tt.Error("search error 2")
		}
	}
	// "abcd" --> 5
	res, _ := t.Prefix([]byte("abcd"))
	if len(res) != 5 {
		tt.Error("prefix search error")
	}

	var st []string
	fn := func(k []byte, _ int) error {
		st = append(st, string(k))
		return nil
	}
	_ = t.Scan(fn)
	for _, s := range st {
		tt.Log(s)
	}

	_ = t.DeleteByPrefix([]byte("abc"))
	tt.Log("after delete ,trie len", t.Len()) // 9
}

func TestTrie(t *testing.T) {
	testTrieRW(t, 100)
}
