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

// TODO implement some algebra operations about graph

func Union[K comparable, V any, W number](g1, g2 Graph[K, V, W]) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

func Intersection[K comparable, V any, W number](g1, g2 Graph[K, V, W]) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

func CartesianProduct[K comparable, V any, W number](g1, g2 Graph[K, V, W]) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

// page 55
func Identify[K comparable, V any, W number](g Graph[K, V, W], v1, v2 K, newVertex K) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

// page 55
func Contract[K comparable, V any, W number](g Graph[K, V, W], v1, v2 K, newVertex K) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

// page 55
func Split[K comparable, V any, W number](g Graph[K, V, W], vertex K, v1, v2, edge K) (Graph[K, V, W], error) {
	return nil, errNotImplement
}

// page 55
func Subdivide[K comparable, V any, W number](g Graph[K, V, W], edge K, newVertex K) (Graph[K, V, W], error) {
	return nil, errNotImplement
}
