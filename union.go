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

// The Union-Find algorithm, also known as the Disjoint-Set data structure,
// is a powerful tool used to manage a collection of disjoint sets.
// It supports two primary operations: Find and Union.
// This data structure is particularly useful in scenarios where we need to
// determine whether two elements belong to the same set or to merge two sets.
type UnionFind struct {
	count  int
	parent []int
	size   []int
}

func NewUnionFind(n int) *UnionFind {
	u := &UnionFind{
		count:  n,
		parent: make([]int, n),
		size:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		u.parent[i] = i
		u.size[i] = 1
	}
	return u
}

// The Find operation determines the representative or root of the set containing a particular element.
func (u *UnionFind) Find(x int) int {
	if x >= len(u.parent) || x < 0 {
		return -1
	}
	for x != u.parent[x] {
		u.parent[x] = u.parent[u.parent[x]]
		x = u.parent[x]
	}
	return x
}

// The Union operation merges two sets by connecting their roots.
func (u *UnionFind) Union(x, y int) {
	rx := u.Find(x)
	ry := u.Find(y)
	if rx < 0 || ry < 0 || rx == ry {
		return
	}
	if u.size[rx] > u.size[ry] {
		u.parent[ry] = rx
		u.size[rx] += u.size[ry]
	} else {
		u.parent[rx] = ry
		u.size[ry] += u.size[rx]
	}
	u.count--
}

// Check if two elements are connected.
func (u *UnionFind) Connected(x, y int) bool {
	return u.Find(x) == u.Find(y)
}

// return the subtree size with root as x.
func (u *UnionFind) Size(x int) int {
	if x < 0 || x >= len(u.parent) {
		return -1
	}
	return u.size[x]
}

// The current number of connected components.
func (u *UnionFind) Component() int {
	return u.count
}

// TODO
type dynamicUnionFind struct {
	count  int
	parent []int
	size   []int
}

func newDynamicUnionFind(n int) *dynamicUnionFind {
	u := &dynamicUnionFind{
		count:  n,
		parent: make([]int, n),
		size:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		u.parent[i] = i
		u.size[i] = 1
	}
	return u
}

func (u *dynamicUnionFind) Find(x int) int {
	if x >= len(u.parent) || x < 0 {
		return -1
	}
	for x != u.parent[x] {
		u.parent[x] = u.parent[u.parent[x]]
		x = u.parent[x]
	}
	return x
}

func (u *dynamicUnionFind) Union(x, y int) {
	rx := u.Find(x)
	ry := u.Find(y)
	if rx < 0 || ry < 0 || rx == ry {
		return
	}
	if u.size[rx] > u.size[ry] {
		u.parent[ry] = rx
		u.size[rx] += u.size[ry]
	} else {
		u.parent[rx] = ry
		u.size[ry] += u.size[rx]
	}
	u.count--
}

func (u *dynamicUnionFind) Connected(x, y int) bool {
	return u.Find(x) == u.Find(y)
}

// return the subtree size with root as x.
func (u *dynamicUnionFind) Size(x int) int {
	if x < 0 || x >= len(u.parent) {
		return -1
	}
	return u.size[x]
}

func (u *dynamicUnionFind) Component() int {
	return u.count
}

func (u *dynamicUnionFind) GetParent(x int) int {
	if x < len(u.parent) && x >= 0 {
		return u.parent[x]
	}
	return -1
}

func (u *dynamicUnionFind) SetParent(x int, y int) {
	if x < 0 || x >= len(u.parent) || y < 0 || y >= len(u.parent) {
		return
	}
	u.parent[x] = y
	for u.parent[y] != y {
		u.size[y] += u.size[x]
		y = u.parent[y]
	}
}

func (u *dynamicUnionFind) Reset() {
	for i := 0; i < len(u.parent); i++ {
		u.parent[i] = i
		u.size[i] = 1
	}
	u.count = len(u.parent)
}

func (u *dynamicUnionFind) Add(n int) {
	l := len(u.parent)
	for i := 0; i < n; i++ {
		u.parent = append(u.parent, l+i)
		u.size = append(u.size, 1)
	}
	u.count += n
}

func (u *dynamicUnionFind) Cut(x int) {
	if x < 0 || x >= len(u.parent) {
		return
	}
	if u.parent[x] == x {
		return
	}
	u.parent[x] = x
	for y := u.parent[x]; ; {
		u.size[y] -= u.size[x]
		if y == u.parent[y] {
			break
		}
		y = u.parent[y]
	}
}
