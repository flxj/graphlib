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

package array

func BuildSA(s string) []int {
	n := len(s)
	if n == 0 {
		return []int{}
	}

	sa := make([]int, n)
	rank := make([]int, n)
	tmp := make([]int, n)
	tmpSa := make([]int, n)
	cnt := make([]int, max(n, 256)+1)

	for i := 0; i < n; i++ {
		sa[i] = i
		rank[i] = int(s[i]) + 1
	}

	for k := 1; k < n; k <<= 1 {
		for i := range cnt {
			cnt[i] = 0
		}
		for i := 0; i < n; i++ {
			key := 0
			if sa[i]+k < n {
				key = rank[sa[i]+k]
			}
			cnt[key]++
		}
		for i := 1; i < len(cnt); i++ {
			cnt[i] += cnt[i-1]
		}
		for i := n - 1; i >= 0; i-- {
			key := 0
			if sa[i]+k < n {
				key = rank[sa[i]+k]
			}
			cnt[key]--
			tmpSa[cnt[key]] = sa[i]
		}

		for i := range cnt {
			cnt[i] = 0
		}
		for i := 0; i < n; i++ {
			cnt[rank[tmpSa[i]]]++
		}
		for i := 1; i < len(cnt); i++ {
			cnt[i] += cnt[i-1]
		}
		for i := n - 1; i >= 0; i-- {
			key := rank[tmpSa[i]]
			cnt[key]--
			sa[cnt[key]] = tmpSa[i]
		}

		tmp[sa[0]] = 1
		for i := 1; i < n; i++ {
			p, c := sa[i-1], sa[i]
			pk1, ck1 := rank[p], rank[c]
			pk2, ck2 := 0, 0
			if p+k < n {
				pk2 = rank[p+k]
			}
			if c+k < n {
				ck2 = rank[c+k]
			}
			if pk1 == ck1 && pk2 == ck2 {
				tmp[c] = tmp[p]
			} else {
				tmp[c] = tmp[p] + 1
			}
		}
		rank, tmp = tmp, rank
		if rank[sa[n-1]] == n {
			break
		}
	}
	return sa
}

func BuildLCP(s string, sa []int) []int {
	n := len(s)
	if n == 0 {
		return []int{}
	}

	rank := make([]int, n)
	for i, p := range sa {
		rank[p] = i
	}

	height := make([]int, n)
	h := 0
	for i := 0; i < n; i++ {
		if rank[i] > 0 {
			j := sa[rank[i]-1]
			for i+h < n && j+h < n && s[i+h] == s[j+h] {
				h++
			}
			height[rank[i]] = h
			if h > 0 {
				h--
			}
		}
	}
	return height
}

/*
For a string s, its suffix array sa[i] represents
the starting position of the suffix ranked at position
i when all suffixes are sorted lexicographically.

For example, s = "banana", all suffixes:

	0: banana
	1: anana
	2: nana
	3: ana
	4: na
	5: a

After sorting lexicographically:

	5: a
	3: ana
	1: anana
	0: banana
	4: na
	2: nana

So sa=[5, 3, 1, 0, 4, 2].

Usually combined with a rank array (also known as rk):
rank[i] represents the ranking of suffixes starting from position i.
The two are mutually inverse operations: sa[rank[i]]=i.
*/
type SuffixArray struct {
	S    string
	SA   []int
	Rank []int
	// The suffix array is often used in conjunction
	// with LCP (Longest Common prefix).
	// Height[i]=the longest common prefix length
	// of the two suffixes sa[i] and sa[i-1]
	Height []int
}

func NewSuffixArray(s string) *SuffixArray {
	sa := BuildSA(s)
	height := BuildLCP(s, sa)

	n := len(s)
	rank := make([]int, n)
	for i, p := range sa {
		rank[p] = i
	}

	return &SuffixArray{
		S:      s,
		SA:     sa,
		Rank:   rank,
		Height: height,
	}
}

// LCP returns the longest common prefix
// length between any two suffixes i and j.
func (sa *SuffixArray) LCP(i, j int) int {
	if i == j {
		return len(sa.S) - i
	}
	ri, rj := sa.Rank[i], sa.Rank[j]
	if ri > rj {
		ri, rj = rj, ri
	}
	res := 1 << 30
	for k := ri + 1; k <= rj; k++ {
		if sa.Height[k] < res {
			res = sa.Height[k]
		}
	}
	return res
}

// returns the number of different substrings.
func (sa *SuffixArray) CountDistinctSubstrings() int64 {
	n := len(sa.S)
	total := int64(n) * int64(n+1) / 2
	for _, h := range sa.Height {
		total -= int64(h)
	}
	return total
}

// Return the longest repeated substring (which can overlap)
func (sa *SuffixArray) LongestRepeatedSubstring() string {
	best, pos := 0, 0
	for i := 1; i < len(sa.Height); i++ {
		if sa.Height[i] > best {
			best = sa.Height[i]
			pos = sa.SA[i]
		}
	}
	return sa.S[pos : pos+best]
}
