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

type rbNode[K any, V any] struct {
	red bool
	key K
	val V
	p   *rbNode[K, V]
	l   *rbNode[K, V]
	r   *rbNode[K, V]
}

func (t *rbNode[K, V]) left() bstNode[K, V] {
	return t.l
}

func (t *rbNode[K, V]) right() bstNode[K, V] {
	return t.r
}

func (t *rbNode[K, V]) getKey() K {
	return t.key
}

func (t *rbNode[K, V]) getVal() V {
	return t.val
}

// A Red-Black Tree is a self-balancing binary search tree with a height limit of O(logN),
// enabling efficient search, insertion, and deletion operations in O(logN) time,
// unlike standard binary search trees which can take O(N) time.
type RedBlackTree[K any, V any] struct {
	comp  CompareFunc[K]
	count int
	root  *rbNode[K, V]
}

// Create a red black tree.
func NewRedBlackTree[K any, V any](comp CompareFunc[K]) *RedBlackTree[K, V] {
	return &RedBlackTree[K, V]{comp: comp}
}

func (t *RedBlackTree[K, V]) Len() int {
	return t.count
}

func (t *RedBlackTree[K, V]) Compare(k1, k2 K) int {
	return t.comp(k1, k2)
}

// 1. The root of the tree is colored Black.
// 2. A Red node can have only Black children.
// 3. Every path from the root to a leaf contains the same number of Black nodes
func (t *RedBlackTree[K, V]) rebalance(x *rbNode[K, V]) {
	if x.p == nil { // x is root
		x.red = false
		t.root = x
		return
	}
	y := x.p
	// 1.Parent is Black:
	// In this case we can safely color the problematic node Red.
	// It is easy to see that the invariant is not violated. X cannot be the root, so Red is a valid color.
	if !y.red {
		x.red = true
		return
	} else {
		// 2. Parent is Red:
		// This case is little more involved, because we need to know the color of z, x’s uncle, or y’s sibling
		// (y is x’s parent). p cannot be the root of the tree, because it is Red, soCexists.
		// If not recall that NIL pointers are taken as Blacks. Again, there are two subcases.
		var gl bool
		var z *rbNode[K, V]
		g := y.p
		if y == g.l {
			z = g.r
			gl = true
		} else {
			z = g.l
		}
		// (a) Uncle is Red
		if z != nil && z.red {
			// In this case, we color X Red.before doing so, we have to make sure that its parent has become Black.
			y.red, z.red = false, false
			x.red = true
			t.rebalance(g)
		} else {
			// (b) Uncle is Black
			if gl && y.l == x {
				// x is g's left-left child, right rotate
				g.l = y.r
				if y.r != nil {
					y.r.p = g
				}
				y.p = g.p
				if g.p != nil {
					if g == g.p.l {
						g.p.l = y
					} else {
						g.p.r = y
					}
				}
				y.r = g
				g.p = y
				x.red, g.red = true, true
				t.rebalance(y)
			} else if !gl && y.r == x {
				// x is g's right-right child, left rotate
				g.r = y.l
				if y.l != nil {
					y.l.p = g
				}
				y.p = g.p
				if g.p != nil {
					if g == g.p.l {
						g.p.l = y
					} else {
						g.p.r = y
					}
				}
				y.l = g
				g.p = y
				x.red, g.red = true, true
				t.rebalance(y)
			} else if gl && y.r == x {
				// x is g's left-right child, left+right rotate
				y.r = x.l
				if x.l != nil {
					x.l.p = y
				}
				x.l = y
				y.p = x

				g.l = x.r
				if x.r != nil {
					x.r.p = g
				}
				x.p = g.p
				if g.p != nil {
					if g == g.p.l {
						g.p.l = x
					} else {
						g.p.r = x
					}
				}
				x.r = g
				g.p = x
				y.red, g.red = true, true
				t.rebalance(x)
			} else {
				y.l = x.r
				if x.r != nil {
					x.r.p = y
				}
				x.r = y
				y.p = x

				g.r = x.l
				if x.l != nil {
					x.l.p = g
				}
				x.p = g.p
				if g.p != nil {
					if g == g.p.l {
						g.p.l = x
					} else {
						g.p.r = x
					}
				}
				x.l = g
				g.p = x
				y.red, g.red = true, true
				t.rebalance(x)
			}
		}
	}
}

// Insert key value data, and if the key already exists, update its value in place.
func (t *RedBlackTree[K, V]) Insert(k K, v V) {
	var p *rbNode[K, V]
	for q := t.root; q != nil; {
		if t.comp(q.key, k) == 0 {
			q.val = v
			return
		} else if t.comp(q.key, k) > 0 {
			if q.l == nil {
				// insert a new node as p's left child
				p = &rbNode[K, V]{red: true, key: k, val: v, p: q}
				q.l = p
				break
			} else {
				q = q.l
			}
		} else {
			if q.r == nil {
				p = &rbNode[K, V]{red: true, key: k, val: v, p: q}
				q.r = p
				break
			} else {
				q = q.r
			}
		}
	}
	t.count++
	if p == nil {
		t.root = &rbNode[K, V]{key: k, val: v, p: p}
		return
	}
	t.rebalance(p)
}

// Use key to query elements. If the element does not exist, return null data value and false flag.
func (t *RedBlackTree[K, V]) Search(k K) (v V, ok bool) {
	for p := t.root; p != nil; {
		if t.comp(p.key, k) == 0 {
			return p.val, true
		} else if t.comp(p.key, k) < 0 {
			p = p.r
		} else {
			p = p.l
		}
	}
	return
}

// Query the minimum element of the key.
func (t *RedBlackTree[K, V]) Min() (k K, v V, ok bool) {
	for p := t.root; p != nil; p = p.l {
		if p.l == nil {
			return p.key, p.val, true
		}
	}
	return
}

// Query the maximum element of the key.
func (t *RedBlackTree[K, V]) Max() (k K, v V, ok bool) {
	for p := t.root; p != nil; p = p.r {
		if p.r == nil {
			return p.key, p.val, true
		}
	}
	return
}

// Delete key. If the key exists, delete it and return its value. If the key does not exist, return a false flag.
func (t *RedBlackTree[K, V]) Delete(k K) (v V, ok bool) {
	var p *rbNode[K, V]
	for q := t.root; q != nil; {
		if t.comp(q.key, k) == 0 {
			v, ok = q.val, true
			for p = q.r; p != nil && p.l != nil; p = p.l {
			}
			if p != nil {
				q.key, q.val = p.key, p.val
			} else {
				for p = q.l; p != nil && p.r != nil; p = p.r {
				}
				if p != nil {
					q.key, q.val = p.key, p.val
				} else {
					p = q
				}
			}
			break
		} else if t.comp(q.key, k) < 0 {
			q = q.r
		} else {
			q = q.r
		}
	}
	if p != nil {
		t.count--
		t.del(p)
	}
	return
}

func (t *RedBlackTree[K, V]) del(x *rbNode[K, V]) {
	z := x.l
	if x.r != nil {
		z = x.r
	}
	// If X is Red, then we just remove it. The invariant is not affected,
	// since X cannot be the root, and its removal does not change the number of Black nodes in any path.
	if x.red {
		if x == x.p.l {
			x.p.l = z
		} else {
			x.p.r = z
		}
		if z != nil {
			z.p = x.p
		}
		x.p, x.l, x.r = nil, nil, nil
	} else {
		// If X is Black and its single child is Red, then we can safely remove X and color the child Black.
		if z != nil && z.red {
			if x.p == nil { // x is root
				z.p = nil
				z.red = false
				t.root = z
				return
			} else {
				if x == x.p.l {
					x.p.l = z
				} else {
					x.p.r = z
				}
				z.p = x.p
				z.red = false
				x.p, x.l, x.r = nil, nil, nil
				return
			}
		}
		// So, the only case that is really the problem is when the node to be deleted is Black and it has
		// a Black child or no children at all (remember that NIL pointers are treated as Black).
		y := x.p
		if y == nil { // x is root
			if z != nil {
				z.p = nil
				z.red = false
			}
			t.root = z
			return
		}
		// 1. Parent is Red (and Sibling is Black)
		// If the parent of X is Red, then X must have a sibling, otherwise the invariant would be violated
		// (the path through X wouldhaveat leastonemoreBlack node, that is X, than the path that
		// would terminate at the NIL pointer of X’s parent).
		// Moreover, the sibling has to be Black,given that it is a son of a Red node.
		if y.red {
			s := y.l
			sl := true
			if x == y.l {
				s = y.r
				sl = false
			}
			// (a) Close Nephew is Black
			// The idea is to pull the Black sibling s to the top to compensate for the loss of the Black X on the path (y −x −T1).
			// This is done by a single rotation.
			if !sl && (s.l == nil || !s.l.red) {
				y.l = z
				if z != nil {
					z.p = y
				}
				x.p, x.l, x.r = nil, nil, nil
				y.r = s.l
				if s.l != nil {
					s.l.p = y
				}
				s.l = y
				s.p = y.p
				if y.p.l == y { // y is red so its not root, y.p != nil
					y.p.l = s
				} else {
					y.p.r = s
				}
				y.p = s
				return
			}
			if sl && (s.r == nil || !s.r.red) {
				y.r = z
				if z != nil {
					z.p = y
				}
				x.p, x.l, x.r = nil, nil, nil
				y.l = s.r
				if s.r != nil {
					s.r.p = y
				}
				s.r = y
				s.p = y.p
				if y.p.l == y {
					y.p.l = s
				} else {
					y.p.r = s
				}
				y.p = s
				return
			}
			// (b) Close Nephew is Red
			// The problem is fully restored by a double rotation and some recoloring.
			if !sl && s.l != nil && s.l.red {
				y.l = z
				if z != nil {
					z.p = y
				}
				w := s.l
				s.l = w.r
				if w.r != nil {
					w.r.p = s
				}
				w.r = s
				s.p = w
				y.r = w.l
				if w.l != nil {
					w.l.p = y
				}
				w.l = y
				w.p = y.p
				if y.p != nil {
					if y == y.p.l {
						y.p.l = w
					} else {
						y.p.r = w
					}
				}
				y.p = w
				y.red = false
				return
			}
			if sl && s.r != nil && s.r.red {
				y.r = z
				if z != nil {
					z.p = y
				}
				w := s.r
				s.r = w.l
				if w.l != nil {
					w.l.p = s
				}
				w.l = s
				s.p = w
				y.l = w.r
				w.r = y
				w.p = y.p
				if y.p != nil {
					if y == y.p.l {
						y.p.l = w
					} else {
						y.p.r = w
					}
				}
				y.p = w
				y.red = false
				return
			}
		} else {
			// 2.Parent is Black
			// X must have a sibling, as in the previous case, otherwise the invariant would be violated (the
			// path through X wouldhaveat leastonemoreBlack node, that is X, than the path that would
			// terminate at the NIL pointer of X’s parent). However, we cannot determine the color of X’s sibling
			// given this situation. So, we distinguish among two cases:
			s := y.l
			sl := true
			if x == y.l {
				s = y.r
				sl = false
			}
			// (a) Sibling is Red
			// The problem is attacked by a single rotation that,interestingly,
			// pushes the problem down! However, notice that the resulting situation is the
			// one where the parent is Red and the sibling is Black, which we just described and can be
			// fully resolved locally.
			if s.red {
				if sl {
					y.l = s.r
					if s.r != nil {
						s.r.p = y
					}
					s.r = y
					s.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = s
						} else {
							y.p.r = s
						}
					} else {
						t.root = s
					}
					y.p = s
					y.red, s.red = true, false
					t.del(x)
				} else {
					y.r = s.l
					if s.l != nil {
						s.l.p = y
					}
					s.l = y
					s.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = s
						} else {
							y.p.r = s
						}
					} else {
						t.root = s
					}
					y.p = s
					y.red, s.red = true, false
					t.del(x)
				}
			} else {
				// (b) Sibling is Black
				// There are several cases, depending on the coloring of the nephews of X.
				// Notice that the nephews do not necessarily exist. However, the case where some
				// nephew does not exist is equivalent to that where this nephew is colored Black.

				// b.i. Far Nephew is Red
				// In this case, the problem is resolved with a single rotation
				if !sl && s.r != nil && s.r.red {
					y.l = z
					if z != nil {
						z.p = y
					}
					x.p, x.l, x.r = nil, nil, nil
					y.r = s.l
					if s.l != nil {
						s.l.p = y
					}
					s.l = y
					s.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = s
						} else {
							y.p.r = s
						}
					} else {
						// y is root
						t.root = s
					}
					s.r.red = false
					return
				}
				if sl && s.l != nil && s.l.red {
					y.r = z
					if z != nil {
						z.p = y
					}
					x.p, x.l, x.r = nil, nil, nil
					y.l = s.r
					if s.r != nil {
						s.r.p = y
					}
					s.r = y
					s.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = s
						} else {
							y.p.r = s
						}
					} else {
						t.root = s
					}
					s.l.red = false
					return
				}
				// b.ii. Far Nephew is Black
				// This case requires checking the color of the other nephew.
				// Depending on its color there are two cases:

				// b.ii.A Close Nephew is Red
				if !sl && s.l != nil && s.l.red {
					y.l = z
					if z != nil {
						z.p = y
					}
					w := s.l
					s.l = w.r
					if w.r != nil {
						w.r.p = s
					}
					w.r = s
					s.p = w
					y.r = w.l
					if w.l != nil {
						w.l.p = y
					}
					w.l = y
					w.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = w
						} else {
							y.p.r = w
						}
					} else {
						t.root = w
					}
					y.p = w
					w.red = false
					return
				}
				if sl && s.r != nil && s.r.red {
					y.r = z
					if z != nil {
						z.p = y
					}
					w := s.r
					s.r = w.l
					if w.l != nil {
						w.l.p = s
					}
					w.l = s
					s.p = w
					y.l = w.r
					if w.r != nil {
						w.r.p = y
					}
					w.r = y
					w.p = y.p
					if y.p != nil {
						if y == y.p.l {
							y.p.l = w
						} else {
							y.p.r = w
						}
					} else {
						t.root = w
					}
					y.p = w
					w.red = false
					return
				}
				// b.ii.B. Close Nephew is Black
				// This last case is tricky, because the existence of only Black nodes in the neighborhood,
				// does not provide for removal of X and local recoloring and/or restructuring. The only
				// thing we can do in case, is to push the problem up in the tree.
				if !sl && (s.l == nil || !s.l.red) {
					s.red = true
					y.l = z
					if z != nil {
						z.p = y
					}
					x.p, x.l, x.r = nil, nil, nil
					if y.p != nil {
						x.p = y.p
						if y == y.p.l {
							y.p.l = x
						} else {
							y.p.r = x
						}
						y.p = x
						x.r = y
						t.del(x)
					} else {
						t.root = y
						return
					}
				}
				if sl && (s.r == nil || !s.r.red) {
					s.red = true
					y.r = z
					if z != nil {
						z.p = y
					}
					x.p, x.l, x.r = nil, nil, nil
					if y.p != nil {
						x.p = y.p
						if y == y.p.l {
							y.p.l = x
						} else {
							y.p.r = x
						}
						y.p = x
						x.l = y
						t.del(x)
					} else {
						t.root = y
						return
					}
				}
			}
		}
	}
}

func (t *RedBlackTree[K, V]) Clean() {
	t.count = 0
	t.root = nil
}

func (t *RedBlackTree[K, V]) Cursor() Cursor[K, V] {
	return newBSTCursor(t.root, t.comp)
}
