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

type DifferenceArray[N number] struct {
	diff []N
}

func NewDifferenceArray[N number](nums []N) *DifferenceArray[N] {
	d := &DifferenceArray[N]{
		diff: make([]N, len(nums)),
	}
	if len(nums) > 0 {
		d.diff[0] = nums[0]
		for i := 1; i < len(nums); i++ {
			d.diff[i] = nums[i] - nums[i-1]
		}
	}
	return d
}

func (d *DifferenceArray[N]) Append(n N) {
	if len(d.diff) == 0 {
		d.diff = append(d.diff, n)
	} else {
		t := d.diff[0]
		for i := 1; i < len(d.diff); i++ {
			t = t + d.diff[i]
		}
		d.diff = append(d.diff, n-t)
	}
}

func (d *DifferenceArray[N]) Len() int {
	return len(d.diff)
}

// closed interval [left,right].
func (d *DifferenceArray[N]) Add(n N, left, right int) bool {
	if left < 0 || right >= len(d.diff) || left > right {
		return false
	}
	d.diff[left] += n
	if right+1 < len(d.diff) {
		d.diff[right+1] -= n
	}
	return true
}

// sum of closed interval [left,right].
func (d *DifferenceArray[N]) Sum(left, right int) (N, bool) {
	var a, b N
	if left < 0 || right >= len(d.diff) || left > right {
		return a, false
	}
	for i := 0; i <= right; i++ {
		if i < left {
			a += d.diff[i]
		}
		b += d.diff[i]
	}
	return b - a, true
}

func (d *DifferenceArray[N]) Array() []N {
	arr := make([]N, len(d.diff))
	if len(d.diff) > 0 {
		arr[0] = d.diff[0]
	}
	for i := 1; i < len(d.diff); i++ {
		arr[i] = arr[i-1] + d.diff[i]
	}
	return arr
}

func (d *DifferenceArray[N]) Get(i int) (N, bool) {
	var n N
	if i < 0 || i >= len(d.diff) {
		return n, false
	}
	for j := 0; j <= i; j++ {
		n += d.diff[j]
	}
	return n, true
}
