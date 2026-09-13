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

import (
	"math"
	"math/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randStr(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func maxValue[N number](n N) N {
	switch any(n).(type) {
	case int:
		return any(math.MaxInt).(N)
	case int8:
		return any(math.MaxInt8).(N)
	case int16:
		return any(math.MaxInt16).(N)
	case int32:
		return any(math.MaxInt32).(N)
	case int64:
		return any(math.MaxInt64).(N)
	case uint:
		return any(math.MaxInt).(N)
	case uint8:
		return any(math.MaxUint8).(N)
	case uint16:
		return any(math.MaxUint16).(N)
	case uint32:
		return any(math.MaxUint32).(N)
	case uint64:
		return any(math.MaxInt64).(N)
	case float32:
		return any(math.MaxFloat32).(N)
	case float64:
		return any(math.MaxFloat64).(N)
	default:
		return n
	}
}

func minValue[N number](n N) N {
	switch any(n).(type) {
	case int:
		return any(math.MinInt).(N)
	case int8:
		return any(math.MinInt8).(N)
	case int16:
		return any(math.MinInt16).(N)
	case int32:
		return any(math.MinInt32).(N)
	case int64:
		return any(math.MinInt64).(N)
	case uint:
		return any(math.MinInt).(N)
	case uint8:
		return any(uint8(0)).(N)
	case uint16:
		return any(uint16(0)).(N)
	case uint32:
		return any(uint32(0)).(N)
	case uint64:
		return any(uint64(0)).(N)
	case float32:
		return any(math.SmallestNonzeroFloat32).(N)
	case float64:
		return any(math.SmallestNonzeroFloat64).(N)
	default:
		return n
	}
}

type stack[K comparable] struct {
	elems []K
	idx   int
}

func newStack[K comparable]() *stack[K] {
	return &stack[K]{}
}

func (s *stack[K]) size() int {
	return s.idx
}

func (s *stack[K]) empty() bool {
	return s.idx == 0
}

func (s *stack[K]) push(k K) {
	if s.idx < len(s.elems) {
		s.elems[s.idx] = k
	} else {
		s.elems = append(s.elems, k)
	}
	s.idx++
}

func (s *stack[K]) pop() (K, bool) {
	var k K
	if !s.empty() {
		k = s.elems[s.idx-1]
		s.idx--
		return k, true
	}
	return k, false
}

func (s *stack[K]) contains(k K) bool {
	for i := 0; i < s.idx; i++ {
		if s.elems[i] == k {
			return true
		}
	}
	return false
}

func (s *stack[K]) top() (k K) {
	if s.idx > 0 {
		return s.elems[s.idx-1]
	}
	return
}

func (s *stack[K]) clean() {
	s.idx = 0
}
