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

// map operation
func Map[T, V any](arr *ArrayList[T], op func(T) V) *ArrayList[V] {
	if arr == nil {
		return nil
	}
	res := &ArrayList[V]{
		head: make([]V, len(arr.head)),
		tail: make([]V, len(arr.tail)),
	}
	for i := len(res.head) - 1; i >= 0; i-- {
		res.head[i] = op(arr.head[i])
	}
	for i := 0; i < len(res.tail); i++ {
		res.tail[i] = op(arr.tail[i])
	}
	return res
}

// reduce operation
func Reduce[T any](arr *ArrayList[T], op func(T, T) T) (res T) {
	if arr == nil || arr.Len() == 0 {
		panic("arraylist: array size is zero,not support reduce operation")
	}
	if len(arr.head) > 0 {
		res = arr.head[len(arr.head)-1]
		for i := len(arr.head) - 2; i >= 0; i-- {
			res = op(res, arr.head[i])
		}
		for _, v := range arr.tail {
			res = op(res, v)
		}
	} else {
		res = arr.tail[0]
		for i := 1; i < len(arr.tail); i++ {
			res = op(res, arr.tail[i])
		}
	}
	return
}

// fold operation
func Fold[T any](arr *ArrayList[T], initVal T, op func(T, T) T) T {
	return FoldLeft(arr, initVal, op)
}

// foldLeft operation
func FoldLeft[T, V any](arr *ArrayList[T], initVal V, op func(V, T) V) (res V) {
	if arr == nil {
		panic("arraylist: nil array pointer")
	}
	res = initVal
	_ = arr.ForEach(func(v T) bool {
		res = op(res, v)
		return true
	})
	return
}

// foldRight operation
func FoldRight[T, V any](arr *ArrayList[T], initVal V, op func(T, V) V) (res V) {
	if arr == nil {
		panic("arraylist: nil array pointer")
	}
	res = initVal
	_ = arr.ReverseForEach(func(v T) bool {
		res = op(v, res)
		return true
	})
	return
}

// Copy the subarray elements of source[si: sj] to target[ti: tj]
func Copy[T any](target *ArrayList[T], ti, tj int, source *ArrayList[T], si, sj int) {
	if source == nil {
		return
	} else if target == nil {
		panic("arraylist: copy target array is nil")
	}
	target.checkRange(ti, tj)
	source.checkRange(si, sj)
	for i := 0; i < min(tj-ti, sj-si); i++ {
		target.Set(ti+i, source.Get(si+i))
	}
}
