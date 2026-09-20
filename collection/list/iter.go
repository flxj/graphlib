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

package list

type alCursor[T any] struct {
	arr *ArrayList[T]
	cur int
}

func (c *alCursor[T]) Open() error {
	return nil
}

// release cursor resources.
func (c *alCursor[T]) Close() {}

// The Seek(key) method locates the cursor at the key.
// If the key does not exist, it locates at the next key and returns it
func (c *alCursor[T]) Seek(i int) (j int, v T, ok bool) {
	c.cur = -1
	if c.arr == nil || i < 0 || i >= c.arr.Len() {
		return
	}
	c.cur = i
	return i, c.arr.Get(i), ok
}

// The First method locates the cursor at the minimum element of the set.
// If there is no minimum element (the set is empty), it returns false
func (c *alCursor[T]) First() (i int, v T, ok bool) {
	c.cur = 0
	if c.arr == nil || c.arr.Len() == 0 {
		return
	}
	i, v, ok = 0, c.arr.First(), true
	return
}

// The Last method locates the cursor at the maximum element of the set.
// If there is no maximum element (the set is empty), it returns false.
func (c *alCursor[T]) Last() (i int, v T, ok bool) {
	if c.arr == nil || c.arr.Len() == 0 {
		return
	}
	c.cur = c.arr.Len() - 1
	i, v, ok = c.cur, c.arr.Last(), true
	return
}

// HasNext returns whether the next element exists relative to the current cursor position.
func (c *alCursor[T]) HasNext() bool {
	if c.arr == nil || c.arr.Len() == 0 {
		return false
	}
	return c.cur < c.arr.Len()
}

// Next() moves the cursor backwards and returns the element.
// If the element does not exist, it returns a type zero value.
func (c *alCursor[T]) Next() (i int, v T) {
	if c.arr == nil {
		return
	}
	v = c.arr.Get(c.cur + 1)
	c.cur++
	return c.cur, v
}

// HasPrev() returns whether the previous element exists relative to the current cursor position.
func (c *alCursor[T]) HasPrev() bool {
	if c.arr == nil || c.arr.Len() == 0 {
		return false
	}
	return c.cur > 0
}

// Prev() moves the cursor forward and returns the element.
// If the element does not exist, it returns a type value of zero.
func (c *alCursor[T]) Prev() (i int, v T) {
	if c.arr == nil {
		return
	}
	v = c.arr.Get(c.cur - 1)
	c.cur--
	return c.cur, v
}
