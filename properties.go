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

func isConnected[K comparable, V any, W number](g Graph[K, V, W], vertex K, edges map[K]Edge[K, W]) (bool, error) {
	var nilK K
	var v K
	for {
		rv, err := g.RandomVertex()
		if err != nil {
			return false, err
		}
		if rv.Key != vertex {
			v = rv.Key
			break
		}
	}
	//
	order := g.Order()
	if vertex != nilK {
		order--
	}
	visited := make(map[K]bool)
	que := newFIFO[K]()
	que.push(v)

	for !que.empty() {
		p, _ := que.pop()
		if _, ok := visited[p]; !ok {
			visited[p] = true
		} else {
			continue
		}
		//
		es, err := g.IncidentEdges(p)
		if err != nil {
			return false, err
		}

		vs := make(map[K]bool)
		for _, e := range es {
			if e.Tail != vertex && e.Head != vertex {
				_, ok := edges[e.Key]
				if !ok {
					vs[e.Tail] = true
					vs[e.Head] = true
				}
			}
		}
		for k := range vs {
			if k != p {
				if _, ok := visited[k]; !ok {
					que.push(v)
				}
			}
		}
	}

	return len(visited) == order, nil
}

//
func IsCutvertex[K comparable, V any, W number](g Graph[K, V, W], vertex K) (bool, error) {
	vs, err := g.Neighbours(vertex)
	if err != nil {
		return false, err
	}
	if len(vs) == 0 || len(vs) == 1 {
		return false, nil
	}

	es := make(map[K]Edge[K, W])
	for _, v := range vs {
		ee, err := g.GetEdge(vertex, v.Key)
		if err != nil {
			return false, err
		}
		for _, e := range ee {
			es[e.Key] = e
		}
	}

	ok, err := isConnected(g, vertex, es)
	if err != nil {
		return false, err
	}
	return !ok, nil
}

func IsBridge[K comparable, V any, W number](g Graph[K, V, W], edge K) (bool, error) {
	e, err := g.GetEdgeByKey(edge)
	if err != nil {
		return false, err
	}
	es := map[K]Edge[K, W]{e.Key: e}

	var k K

	ok, err := isConnected(g, k, es)
	if err != nil {
		return false, err
	}
	return !ok, nil
}
