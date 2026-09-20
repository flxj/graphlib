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

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/flxj/graphlib/collection"
)

// ArrayList are dynamic generic arrays implemented by Golang,
// (ArrayList is non concurrent secure),
// similar to Java's ArrayList and Scala's ArrayBuffer.
type ArrayList[T any] struct {
	head []T
	tail []T
}

// Create a new ArrayList with length and capacity parameters
// that have the same semantic meaning as the golang slice.
func NewArrayList[T any](length int, capacity int) *ArrayList[T] {
	if length < 0 {
		panic(fmt.Sprintf("arraylist: negative length parameter %d", length))
	}
	if capacity < length {
		panic(fmt.Sprintf("arraylist: length parameter %d greater than capacity %d", length, capacity))
	}
	arr := &ArrayList[T]{}
	arr.head = make([]T, length/2, capacity/2)
	arr.tail = make([]T, length-length/2, capacity-capacity/2)
	return arr
}

// The current length of the array
func (l *ArrayList[T]) Len() int {
	return len(l.head) + len(l.tail)
}

// The current capacity of the array
func (l *ArrayList[T]) Cap() int {
	return cap(l.head) + cap(l.tail)
}

// Retrieve elements based on subscripts.
func (l *ArrayList[T]) Get(i int) T {
	l.checkIndex(i)
	if i < len(l.head) {
		return l.head[len(l.head)-i-1]
	} else {
		return l.tail[i-len(l.head)]
	}
}

// Get the first element of the array
func (l *ArrayList[T]) First() (v T) {
	if l.Len() > 0 {
		if len(l.head) > 0 {
			v = l.head[len(l.head)-1]
		} else {
			v = l.tail[0]
		}
	}
	return
}

// Get the last element of the array
func (l *ArrayList[T]) Last() (v T) {
	if l.Len() > 0 {
		if len(l.tail) > 0 {
			v = l.tail[len(l.tail)-1]
		} else {
			v = l.head[0]
		}
	}
	return
}

// Update the element corresponding to the specified index.
func (l *ArrayList[T]) Set(i int, elem T) {
	l.checkIndex(i)
	if i < len(l.head) {
		l.head[len(l.head)-i-1] = elem
	} else {
		l.tail[i-len(l.head)] = elem
	}
}

// Insert a new element at the specified index
func (l *ArrayList[T]) Insert(i int, elem T) {
	l.checkIndex(i)
	if i < len(l.head) {
		i = len(l.head) - i - 1
		l.head = append(l.head, elem)
		copy(l.head[i+1:], l.head[i:])
		l.head[i] = elem
	} else {
		i -= len(l.head)
		l.tail = append(l.tail, elem)
		copy(l.tail[i+1:], l.tail[i:])
		l.tail[i] = elem
	}
}

// Insert several new elements at the specified index
func (l *ArrayList[T]) InsertAll(i int, elems ...T) {
	if len(elems) == 0 {
		return
	}
	l.checkIndex(i)
	if i < len(l.head) {
		old := len(l.head)
		i = old - i - 1
		l.head = append(l.head, elems...)
		copy(l.head[i+len(elems):], l.head[i:old])
		for j := 0; j < len(elems); j++ {
			l.head[i+len(elems)-j-1] = elems[j]
		}
	} else {
		old := len(l.tail)
		i -= old
		l.tail = append(l.tail, elems...)
		copy(l.tail[i+len(elems):], l.tail[i:old])
		copy(l.tail[i:i+len(elems)], elems)
	}
	l.rebalance()
}

// add a new element at the head of the array.
func (l *ArrayList[T]) Prepend(elem T) {
	l.head = append(l.head, elem)
	l.rebalance()
}

// add several new elements at the head of the array,
// note that after PrependAll, the last element of elems
// will be at the head of the array.
func (l *ArrayList[T]) PrependAll(elems ...T) {
	l.head = append(l.head, elems...)
	l.rebalance()
}

// Add a new element at the end of the array.
func (l *ArrayList[T]) Append(elem T) {
	l.tail = append(l.tail, elem)
	l.rebalance()
}

// Add several new elements at the end of the array.
func (l *ArrayList[T]) AppendAll(elems ...T) {
	l.tail = append(l.tail, elems...)
	l.rebalance()
}

// Pop up the last element of the array,
// subtract one from the length of the array,
// and behave similarly to a stack or queue.
func (l *ArrayList[T]) Back() (v T) {
	if len(l.tail) > 0 {
		v = l.tail[len(l.tail)-1]
		l.tail = l.tail[:len(l.tail)-1]
	} else if len(l.head) > 0 {
		v = l.head[0]
		l.head = l.head[1:]
	}
	return
}

// Pop up the first element of the array,
// subtract one from the length of the array,
// and behave similarly to a stack or queue.
func (l *ArrayList[T]) Front() (v T) {
	if len(l.head) > 0 {
		v = l.head[len(l.head)-1]
		l.head = l.head[:len(l.head)-1]
	} else if len(l.tail) > 0 {
		v = l.tail[0]
		l.tail = l.tail[1:]
	}
	return
}

// Delete the element at the specified index,
// reducing the length of the array by one
func (l *ArrayList[T]) Remove(i int) (v T) {
	l.checkIndex(i)
	if i < len(l.head) {
		j := len(l.head) - i - 1
		v = l.head[j]
		if i > 0 {
			copy(l.head[j:], l.head[j+1:])
		}
		l.head = l.head[:len(l.head)-1]
	} else {
		i -= len(l.head)
		v = l.tail[i]
		if i < len(l.tail)-1 {
			copy(l.tail[i:], l.tail[i+1:])
		}
		l.tail = l.tail[:len(l.tail)-1]
	}
	l.rebalance()
	return
}

// Delete consecutive count elements starting from the specified index,
// and the array length will decrease by count. Note that if i+count is
// greater than the length of the array,
// all remaining elements starting from i will be deleted.
func (l *ArrayList[T]) RemoveAll(i int, count int) {
	l.checkIndex(i)
	if count < 0 {
		panic(fmt.Sprintf("araylist: call RemoveAll with negative count %d", count))
	}
	if count == 0 {
		return
	}
	if i < len(l.head) {
		j := len(l.head) - i - 1
		k := max(0, j-count+1)
		// remove head[k:j+1]
		if i > 0 {
			copy(l.head[k:], l.head[j+1:])
		}
		l.head = l.head[:len(l.head)+k-j-1]
		count -= (j - k + 1)
		if count > 0 {
			// remove elements from t.tail
			if count >= len(l.tail) {
				l.tail = l.tail[0:0]
			} else {
				copy(l.tail[:], l.tail[count:])
				l.tail = l.tail[:len(l.tail)-count]
			}
		}
	} else {
		i -= len(l.head)
		if i+count >= len(l.tail) {
			l.tail = l.tail[:i]
		} else {
			copy(l.tail[i:], l.tail[i+count:])
			l.tail = l.tail[:len(l.tail)-count]
		}
	}
	l.rebalance()
}

// in-place sort.
func (l *ArrayList[T]) SortInPlace(less func(T, T) bool) {
	if l.Len() < 2 {
		return
	}
	sort.Slice(l.head, func(i, j int) bool { return less(l.head[i], l.head[j]) })
	sort.Slice(l.tail, func(i, j int) bool { return less(l.tail[i], l.tail[j]) })
	if len(l.head) == 0 || len(l.tail) == 0 {
		if len(l.tail) == 0 {
			slices.Reverse(l.head)
		}
		return
	}
	arr := make([]T, l.Len())
	var i, j, k int
	for i < len(l.head) && j < len(l.tail) {
		for i < len(l.head) && less(l.head[i], l.tail[j]) {
			arr[k] = l.head[i]
			k++
			i++
		}
		if i == len(l.head) {
			break
		}
		for j < len(l.tail) && less(l.tail[j], l.head[i]) {
			arr[k] = l.tail[j]
			k++
			j++
		}
	}
	for ; i < len(l.head); i++ {
		arr[k] = l.head[i]
		k++
	}
	for ; j < len(l.tail); j++ {
		arr[k] = l.tail[j]
		k++
	}
	l.head = arr[: len(arr)/2 : len(arr)/2]
	l.tail = arr[len(arr)/2:]
	slices.Reverse(l.head)
}

// Check if the array is ordered.
func (l *ArrayList[T]) Sorted(less func(T, T) bool) bool {
	if l.Len() < 2 {
		return true
	}
	var flag int
	if len(l.head) > 0 {
		prev := l.head[len(l.head)-1]
		for i := len(l.head) - 2; i >= 0; i-- {
			if less(prev, l.head[i]) {
				flag |= 2
			} else {
				flag |= 1
			}
			if flag == 3 {
				return false
			}
			prev = l.head[i]
		}
		for _, v := range l.tail {
			if less(prev, v) {
				flag |= 2
			} else {
				flag |= 1
			}
			if flag == 3 {
				return false
			}
			prev = v
		}
	} else {
		prev := l.tail[0]
		for i := 1; i < len(l.tail); i++ {
			if less(prev, l.tail[i]) {
				flag |= 2
			} else {
				flag |= 1
			}
			if flag == 3 {
				return false
			}
			prev = l.tail[i]
		}
	}
	return flag != 3
}

// Delete the first n elements of the array,
// The length of the array will be reduced by n
func (l *ArrayList[T]) Trim(n int) {
	l.checkLen(n)
	if n == 0 {
		return
	}
	l.head = l.head[0:max(len(l.head)-n, 0)]
	if n > len(l.head) {
		l.tail = l.tail[n-len(l.head):]
	}
	l.rebalance()
}

// Delete the last n elements of the array,
// The length of the array will be reduced by n
func (l *ArrayList[T]) TrimRight(n int) {
	l.checkLen(n)
	if n == 0 {
		return
	}
	l.tail = l.tail[0:max(len(l.tail)-n, 0)]
	if n > len(l.tail) {
		l.head = l.head[n-len(l.tail):]
	}
	l.rebalance()
}

// Copy a current array.
func (l *ArrayList[T]) Colne() *ArrayList[T] {
	arr := &ArrayList[T]{}
	arr.head = make([]T, len(l.head))
	arr.tail = make([]T, len(l.tail))
	copy(arr.head, l.head)
	copy(arr.tail, l.tail)
	return arr
}

// Call the operation function on all elements in sequence
// and count the number of times it returns true
func (l *ArrayList[T]) Count(op func(T) bool) int {
	var cnt int
	l.ForEach(func(v T) bool {
		if op(v) {
			cnt++
		}
		return true
	})
	return cnt
}

// Return the head subarray of the array,
// which includes all elements except for the last one.
func (l *ArrayList[T]) Head() *ArrayList[T] {
	arr := &ArrayList[T]{}
	if l.Len() > 0 {
		arr.head = l.head[0:len(l.head):len(l.head)]
		if len(l.tail) > 0 {
			arr.tail = l.tail[0 : len(l.tail)-1 : len(l.tail)-1]
		}
	}
	return arr
}

// Return the tail subarray of the array, that is,
// all elements except for the first element.
func (l *ArrayList[T]) Tail() *ArrayList[T] {
	arr := &ArrayList[T]{}
	if l.Len() > 0 {
		arr.tail = l.head[0:len(l.tail):len(l.tail)]
		if len(l.head) > 0 {
			arr.head = l.head[0 : len(l.head)-1 : len(l.head)-1]
		}
	}
	return arr
}

// Extract the first n elements of the array
func (l *ArrayList[T]) Take(n int) *ArrayList[T] {
	l.checkLen(n)
	arr := &ArrayList[T]{}
	if n == 0 {
		return arr
	}
	i := max(len(l.head)-n, 0)
	arr.head = l.head[i : len(l.head) : len(l.head)-i]
	if n > len(l.head) {
		arr.tail = l.tail[0 : n-len(l.head) : n-len(l.head)]
	}
	return arr
}

// Extract the last n elements of the array
func (l *ArrayList[T]) TakeRight(n int) *ArrayList[T] {
	l.checkLen(n)
	arr := &ArrayList[T]{}
	if n == 0 {
		return arr
	}
	i := max(len(l.tail)-n, 0)
	arr.tail = l.tail[i : len(l.tail) : len(l.tail)-i]
	if n > len(l.tail) {
		arr.head = l.head[0 : n-len(l.tail) : n-len(l.tail)]
	}
	return arr
}

// Traverse array elements in sequence.
func (l *ArrayList[T]) ForEach(op func(T) bool) bool {
	for i := len(l.head) - 1; i >= 0; i-- {
		if !op(l.head[i]) {
			return false
		}
	}
	for _, v := range l.tail {
		if !op(v) {
			return false
		}
	}
	return true
}

// Traverse array elements in reverse order
func (l *ArrayList[T]) ReverseForEach(op func(T) bool) bool {
	for i := len(l.tail) - 1; i >= 0; i-- {
		if !op(l.tail[i]) {
			return false
		}
	}
	for _, v := range l.head {
		if !op(v) {
			return false
		}
	}
	return true
}

// Check if all elements in the array meet a certain condition.
func (l *ArrayList[T]) All(op func(T) bool) bool {
	if l.Len() == 0 {
		return false
	}
	return l.ForEach(op)
}

// Find the first element that meets the condition, and return -1 if it does not exist
func (l *ArrayList[T]) Find(op func(T) bool) (idx int, v T) {
	idx = -1
	for i := len(l.head) - 1; i >= 0; i-- {
		if op(l.head[i]) {
			return len(l.head) - i - 1, l.head[i]
		}
	}
	for i, v := range l.tail {
		if op(v) {
			return i + len(l.head), v
		}
	}
	return
}

// Find the last element that meets the condition, and return -1 if it does not exist
func (l *ArrayList[T]) FindLast(op func(T) bool) (idx int, v T) {
	idx = -1
	for i := len(l.tail) - 1; i >= 0; i-- {
		if op(l.tail[i]) {
			return len(l.head) + i, l.tail[i]
		}
	}
	for i, v := range l.head {
		if op(v) {
			return len(l.head) - i - 1, v
		}
	}
	return
}

// Find the first element that satisfies the condition starting from the
// specified index, and return -1 if it does not exist
func (l *ArrayList[T]) FindFrom(op func(T) bool, from int) (idx int, v T) {
	l.checkIndex(from)
	idx = -1
	if from < len(l.head) {
		for i := len(l.head) - from - 1; i >= 0; i-- {
			if op(l.head[i]) {
				return len(l.head) - i - 1, l.head[i]
			}
		}
		for i, v := range l.tail {
			if op(v) {
				return i + len(l.head), v
			}
		}
	} else {
		for i := from - len(l.head); i < len(l.tail); i++ {
			if op(l.tail[i]) {
				return i + len(l.head), l.tail[i]
			}
		}
	}
	return
}

// Filter all elements that meet the criteria
func (l *ArrayList[T]) Filter(op func(T) bool) *ArrayList[T] {
	arr := &ArrayList[T]{}
	for i := len(l.head) - 1; i >= 0; i-- {
		if op(l.head[i]) {
			arr.head = append(arr.head, l.head[i])
		}
	}
	for _, v := range l.tail {
		if op(v) {
			arr.tail = append(arr.tail, v)
		}
	}
	slices.Reverse(arr.head)
	return arr
}

// Update all elements
func (l *ArrayList[T]) MapInPlace(op func(T) T) {
	for i := len(l.head) - 1; i >= 0; i-- {
		l.head[i] = op(l.head[i])
	}
	for i, v := range l.tail {
		l.tail[i] = op(v)
	}
}

// Invert the array in place
func (l *ArrayList[T]) Reverse() {
	l.head, l.tail = l.tail, l.head
}

// Array slicing operation returns a subarray with index range [from, until).
func (l *ArrayList[T]) Range(from, until int) *ArrayList[T] {
	l.checkRange(from, until)
	arr := &ArrayList[T]{}
	if from == until {
		return arr
	}
	if from < len(l.head) {
		if until <= len(l.head) {
			arr.head = l.head[len(l.head)-until : len(l.head)-from : until-from]
		} else {
			arr.head = l.head[0 : len(l.head)-from : len(l.head)-from]
			until -= len(l.head)
			arr.tail = l.tail[0:min(len(l.tail), until):min(len(l.tail), until)]
		}
	} else {
		from -= len(l.head)
		until -= len(l.head)
		arr.tail = l.tail[from : min(len(l.tail), until) : min(len(l.tail), until)-from]
	}
	return arr
}

// Copy and return the underlying Golang slice.
func (l *ArrayList[T]) ToSlice() []T {
	res := make([]T, l.Len())
	copy(res, l.head)
	copy(res[len(l.head):], l.tail)
	for i := 0; i < len(l.head)/2; i++ {
		res[i], res[len(l.head)-i-1] = res[len(l.head)-i-1], res[i]
	}
	return res
}

// Traverse with subscript
func (l *ArrayList[T]) IndexForEach(op func(int, T) bool) bool {
	for i := len(l.head) - 1; i >= 0; i-- {
		if !op(len(l.head)-i-1, l.head[i]) {
			return false
		}
	}
	for i, v := range l.tail {
		if !op(len(l.head)+i, v) {
			return false
		}
	}
	return true
}

// Traverse with subscript in reverse order
func (l *ArrayList[T]) ReverseIndexForEach(op func(int, T) bool) bool {
	for i := len(l.tail) - 1; i >= 0; i-- {
		if !op(len(l.head)+i, l.tail[i]) {
			return false
		}
	}
	for i, v := range l.head {
		if !op(len(l.head)-i-1, v) {
			return false
		}
	}
	return true
}

// Clear the array and free up space
func (l *ArrayList[T]) Clear() {
	l.head = []T{}
	l.tail = []T{}
}

// Return string representation
func (l *ArrayList[T]) String(conv func(T) string) string {
	var b strings.Builder
	b.WriteByte('[')
	_ = l.IndexForEach(func(i int, v T) bool {
		b.WriteString(conv(v))
		if i != l.Len()-1 {
			b.WriteByte(',')
		}
		return true
	})
	b.WriteByte(']')
	return b.String()
}

// Create an Iterator
func (l *ArrayList[T]) Cursor() collection.Cursor[int, T] {
	return &alCursor[T]{arr: l, cur: -1}
}

func (l *ArrayList[T]) rebalance() {
	if len(l.head) > 64 && len(l.head) > 2*len(l.tail) {
		// move 1/4 elements from l.head to l.tail
		slices.Reverse(l.tail)
		l.tail = append(l.tail, l.head[:len(l.head)/4]...)
		slices.Reverse(l.tail)
		l.head = l.head[len(l.head)/4:]
	} else if len(l.tail) > 64 && len(l.tail) > 2*len(l.head) {
		slices.Reverse(l.head)
		l.head = append(l.head, l.tail[:len(l.tail)/4]...)
		slices.Reverse(l.head)
		l.tail = l.tail[len(l.tail)/4:]
	}
}

func (l *ArrayList[T]) checkIndex(i int) {
	if i < 0 || i >= l.Len() {
		panic(fmt.Sprintf("arraylist: index %d out of range [0,%d)", i, l.Len()))
	}
}

func (l *ArrayList[T]) checkLen(n int) {
	if n < 0 || n > l.Len() {
		panic(fmt.Sprintf("arraylist:n %d out of length range:%d", n, l.Len()))
	}
}

func (l *ArrayList[T]) checkRange(i, j int) {
	if i < 0 || i >= l.Len() {
		panic(fmt.Sprintf("arraylist: index %d out of range [0,%d)", i, l.Len()))
	}
	if i > j {
		panic(fmt.Sprintf("arraylist: wrong range, 'from' %d greater than 'until' %d", i, j))
	}
	if j > l.Len() {
		panic(fmt.Sprintf("arraylist: index %d out of range [0,%d)", j, l.Len()))
	}
}
