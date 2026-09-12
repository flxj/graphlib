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
	"slices"
	"sort"
)

/*
MCQ:
An efficient branch-and-bound algorithm for finding a maximum clique,
Etsuji Tomita and Tomokazu Seki,
Proceedings of Discrete Mathematics and Theoretical Computer Science. LNCS 2731, pp. 278–289, 2003.
[doi:10.1007/3-540-45066-1_22](https://doi.org/10.1007/3-540-45066-1_22)

MCR:
An efficient branch-and-bound algorithm for finding a maximum clique with computational experiments,
Etsuji Tomita and Toshikatsu Kameda,
Journal of Global Optimization, **37**:95--111, 2007.
[doi:10.1007/s10898-006-9039-7](https://doi.org/10.1007/s10898-006-9039-7)

MCS:
A simple and faster branch-and-bound algorithm for finding a maximum clique,
Etsuji Tomita, Yoichi Sutani, Takanori Higashi, Shinya Takahashi, and Mitsuo Wakatsuki,
Proceedings of WALCOM 2010, LNCS 5942, pp. 191–203, 2010.
[doi:10.1007/978-3-642-20662-7_31](https://doi.org/10.1007/978-3-642-20662-7_31)
*/

/*
"Our algorithm begins with a small clique, and continues finding larger and larger
cliques until one is found that can be verified to have the maximum size."
*/
func mcq[K comparable, W number](g Graph[K, W]) []K {
	vtx := g.AllVertexes()
	if len(vtx) == 0 {
		return []K{}
	}
	sort.Slice(vtx, func(i, j int) bool {
		di, _ := g.Degree(vtx[i].Key)
		dj, _ := g.Degree(vtx[j].Key)
		return di > dj
	})
	D, _ := g.Degree(vtx[0].Key)
	// in order to prune unnecessary searching,
	// we make use of approximate coloring of vertices.
	// We assign in advance for each p ∈ R
	// a positive integer N[p] called the Number or Color of p.
	Num := make(map[K]int)
	for i := 1; i <= D; i++ {
		Num[vtx[i-1].Key] = i
	}
	for i := D + 1; i <= len(vtx); i++ {
		Num[vtx[i-1].Key] = D + 1
	}
	// we maintain global variables Q, Qmax,
	// where Q consists of vertices of a current clique,
	// Qmax consists of vertices of the largest clique found so far.
	var Qmax []K
	Q := make(map[K]struct{})
	// We select a certain vertex p from R and add p to Q (Q := Q ∪ {p}).
	var expand func([]K, map[K]int)
	expand = func(R []K, N map[K]int) {
		for len(R) != 0 {
			// p := the vertex in R such that N[p] = Max{N[q] | q ∈ R};
			// {i.e., the last vertex in R}
			p := R[len(R)-1]
			if len(Q)+N[p] > len(Qmax) {
				Q[p] = struct{}{}
				// compute rp := R ∩ Γ(p) as the new set of candidate vertices.
				rp := make(map[K]struct{})
				ns, _ := g.Neighbours(p)
				for _, u := range ns {
					if slices.Contains(R, u.Key) {
						rp[u.Key] = struct{}{}
					}
				}
				// This procedure(EXPAND) is applied recursively, while Rp ̸= ∅
				if len(rp) != 0 {
					Rp, Np := numberSort(g, rp)
					expand(Rp, Np)
				} else if len(Q) > len(Qmax) {
					// When Rp = ∅ is reached, Q constitutes a maximal clique.
					Qmax = make([]K, len(Q))
					var i int
					for k := range Q {
						Qmax[i] = k
						i++
					}
				}
				// We then backtrack by removing p from Q and R
				delete(Q, p)
			} else {
				return
			}
			// We select a new vertex p from the resulting R
			// and continue the same procedure until R = ∅
			R = R[:len(R)-1]
		}
	}
	// Let R⊆V consist of vertices (candidates) which may be added to Q.
	R := make([]K, len(vtx))
	for i, v := range vtx {
		R[i] = v.Key
	}
	expand(R, Num)

	return Qmax
}

func numberSort[K comparable, W number](g Graph[K, W], R map[K]struct{}) ([]K, map[K]int) {
	Np := make(map[K]int)
	color := make(map[int][]K)
	maxno := 1
	for len(R) != 0 {
		// p := the first vertex in R;
		var p K
		for v := range R {
			p = v
			break
		}
		k := 1
		// while C_k ∩ Γ(p) ̸= ∅ do k = k +1 ;
		for {
			ck := color[k]
			if len(ck) == 0 {
				break
			}
			ns, _ := g.Neighbours(p)
			var ok bool
			for _, u := range ns {
				if ok = slices.Contains(ck, u.Key); ok {
					break
				}
			}
			if !ok {
				break
			}
			k++
		}
		if k > maxno {
			maxno = k
		}
		Np[p] = k
		color[k] = append(color[k], p)
		delete(R, p)
	}
	var Rp []K
	for k := 1; k <= maxno; k++ {
		Rp = append(Rp, color[k]...)
	}
	return Rp, Np
}

// TODO: mcr, mcs
