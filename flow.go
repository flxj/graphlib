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

// page 163
func MaxNetworkFlow[K comparable, W number](g Digraph[K, any, W]) (W, error) {
	return 0, errNotImplement
}

func MaxMatching[K comparable, W number](g Bipartite[K, any, W]) ([]Edge[K, W], error) {
	return nil, errNotImplement
}

func PerfectMatching[K comparable, W number](g Bipartite[K, any, W]) ([]Edge[K, W], error) {
	return nil, errNotImplement
}
