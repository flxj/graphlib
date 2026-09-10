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

import "github.com/flxj/graphlib/collection"

// Record the current search path.
type path[V any] struct {
	node node[V]
	nth  uint16
}

// Implement art cursor.
type artCursor[V any] struct {
	tree *art[V]
	root node[V]
	stk  *collection.Stack[*path[V]]
}

func newCursor[V any](tree *art[V], root node[V]) *artCursor[V] {
	return &artCursor[V]{
		tree: tree,
		root: root,
		stk:  collection.NewStack[*path[V]](),
	}
}

func (c *artCursor[V]) Open() error {
	if c.tree != nil && c.tree.locked {
		c.tree.mu.RLock()
	}
	return nil
}

func (c *artCursor[V]) Close() {
	if c.tree != nil && c.tree.locked {
		c.tree.mu.RUnlock()
	}
}

// Starting from the root node of art, search for the specified element.
// If the element exists, move the cursor to that element.
func (c *artCursor[V]) Seek(key Key) (k Key, v V, ok bool) {
	if c.tree == nil {
		return
	}
	c.stk.Clean()
	c.root = c.tree.root
	depth := 0
	for n := c.root; n != nil; depth++ {
		if n.kind() == leaf {
			c.stk.Push(&path[V]{node: n, nth: 0})
			break
		}
		if checkPrefix(n, key, depth) != n.prefixLen() {
			return
		}
		depth += n.prefixLen()
		if depth < len(key) {
			nxt, r := n.rank(key[depth])
			if nxt == nil {
				return
			}
			c.stk.Push(&path[V]{node: n, nth: r})
			n = nxt
		} else {
			if leafMatches(n.end(), key, depth) {
				c.stk.Push(&path[V]{node: n, nth: 0})
				c.stk.Push(&path[V]{node: n.end(), nth: 0})
				break
			}
			return
		}
	}
	p := c.stk.Top()
	if p != nil {
		k, v, ok = p.node.key(), p.node.value(), true
	}
	return
}

func (c *artCursor[V]) First() (k Key, v V, ok bool) {
	c.stk.Clean()
	for n := c.root; n != nil; {
		if n.kind() == leaf {
			c.stk.Push(&path[V]{node: n, nth: 0})
			break
		} else {
			if n.end() != nil {
				c.stk.Push(&path[V]{node: n, nth: 0})
				c.stk.Push(&path[V]{node: n.end(), nth: 0})
				break
			}
			c.stk.Push(&path[V]{node: n, nth: 1})
			n = n.nthChild(1)
		}
	}
	p := c.stk.Top()
	if p != nil {
		k, v, ok = p.node.key(), p.node.value(), true
	}
	return
}

func (c *artCursor[V]) Last() (k Key, v V, ok bool) {
	c.stk.Clean()
	for n := c.root; n != nil; {
		if n.kind() == leaf {
			c.stk.Push(&path[V]{node: n, nth: 0})
			break
		} else {
			c.stk.Push(&path[V]{node: n, nth: n.count()})
			n = n.nthChild(n.count())
		}
	}
	p := c.stk.Top()
	if p != nil {
		k, v, ok = p.node.key(), p.node.value(), true
	}
	return
}

func (c *artCursor[V]) HasNext() bool {
	for {
		p := c.stk.Top()
		if p == nil {
			break
		}
		if p.node.kind() != leaf {
			if p.node.count() > p.nth {
				return true
			}
		}
		_, _ = c.stk.Pop()
	}
	return false
}

func (c *artCursor[V]) Next() (k Key, v V) {
	for {
		p := c.stk.Top()
		if p == nil {
			break
		}
		if p.node.kind() != leaf {
			if p.node.count() > p.nth {
				p.nth++
				for n := p.node.nthChild(p.nth); n != nil; {
					if n.end() != nil {
						c.stk.Push(&path[V]{node: n, nth: 0})
						c.stk.Push(&path[V]{node: n.end(), nth: 0})
						break
					}
					c.stk.Push(&path[V]{node: n, nth: 1})
					n = n.nthChild(1)
				}
				break
			}
		}
		_, _ = c.stk.Pop()
	}
	p := c.stk.Top()
	if p != nil {
		k, v = p.node.key(), p.node.value()
	}
	return
}

func (c *artCursor[V]) HasPrev() bool {
	for {
		p := c.stk.Top()
		if p == nil {
			break
		}
		if p.node.kind() != leaf {
			if p.nth > 1 || p.node.end() != nil {
				return true
			}
		}
		_, _ = c.stk.Pop()
	}
	return false
}

func (c *artCursor[V]) Prev() (k Key, v V) {
	for {
		p := c.stk.Top()
		if p == nil {
			break
		}
		if p.node.kind() != leaf {
			if p.nth > 1 {
				p.nth--
				for n := p.node.nthChild(p.nth); n != nil; {
					c.stk.Push(&path[V]{node: n, nth: n.count()})
					n = n.nthChild(n.count())
				}
				break
			} else if p.node.end() != nil {
				c.stk.Push(&path[V]{node: p.node.end(), nth: 0})
			}
		}
		_, _ = c.stk.Pop()
	}
	p := c.stk.Top()
	if p != nil {
		k, v = p.node.key(), p.node.value()
	}
	return
}
