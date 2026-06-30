package graphlib

type UnionFindInt struct {
	count  int
	parent []int
	size   []int
}

func NewUnionFindInt(n int) *UnionFindInt {
	u := &UnionFindInt{
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

func (u *UnionFindInt) Find(x int) int {
	if x >= len(u.parent) || x < 0 {
		return -1
	}
	for x != u.parent[x] {
		u.parent[x] = u.parent[u.parent[x]]
		x = u.parent[x]
	}
	return x
}

func (u *UnionFindInt) Union(x, y int) {
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

func (u *UnionFindInt) Connected(x, y int) bool {
	return u.Find(x) == u.Find(y)
}

// return the subtree size with root as x.
func (u *UnionFindInt) Size(x int) int {
	if x < 0 || x >= len(u.parent) {
		return -1
	}
	return u.size[x]
}

func (u *UnionFindInt) Component() int {
	return u.count
}
