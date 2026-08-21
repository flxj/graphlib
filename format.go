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

import (
	"fmt"
	"strings"
)

func fmtGraphTex[K comparable, W number](g Graph[K, W], F func(v K) ([]Edge[K, W], error)) (string, error) {
	var bu strings.Builder
	_, _ = bu.WriteString("{\n")
	var bfs func(K) error
	vis := make(map[K]struct{})
	bfs = func(v K) error {
		_, _ = fmt.Fprintf(&bu, "%v", v)
		out, err := F(v)
		if err != nil {
			return err
		}
		var filter []Edge[K, W]
		for _, e := range out {
			if _, ok := vis[e.Key]; !ok {
				if !g.IsDigraph() && e.Head == v {
					e.Tail, e.Head = e.Head, e.Tail
				}
				filter = append(filter, e)
			}
		}
		if len(filter) > 1 {
			if g.IsDigraph() {
				_, _ = bu.WriteString("-> {")
			} else {
				_, _ = bu.WriteString("-- {")
			}
			for i, e := range filter {
				vis[e.Key] = struct{}{}
				_, _ = fmt.Fprintf(&bu, "%v", e.Head)
				if i != len(filter)-1 {
					_ = bu.WriteByte(',')
				}
			}
			_ = bu.WriteByte('}')
		} else if len(filter) == 1 {
			if g.IsDigraph() {
				_, _ = bu.WriteString("->")
			} else {
				_, _ = bu.WriteString("--")
			}
			vis[filter[0].Key] = struct{}{}
			_, _ = fmt.Fprintf(&bu, "%v", filter[0].Head)
		}
		return nil
	}
	edges := g.AllEdges()
	for len(vis) != g.Size() {
		for _, e := range edges {
			if _, ok := vis[e.Key]; ok {
				continue
			}
			_, _ = bu.WriteString("\n{")
			if err := bfs(e.Tail); err != nil {
				return "", err
			}
			bu.WriteString("};\n")
		}
	}
	_, _ = bu.WriteString("\n};\n")
	return bu.String(), nil
}

func FmtGraphTex[K comparable, W number](g Graph[K, W]) (string, error) {
	if g == nil || g.Order() == 0 {
		return "", nil
	}
	return fmtGraphTex(g, g.IncidentEdges)
}

func FmtDigraphTex[K comparable, W number](g Digraph[K, W]) (string, error) {
	if g == nil || g.Order() == 0 {
		return "", nil
	}
	return fmtGraphTex(g, g.OutEdges)
}

func fmtForestTex[K comparable, W number](f *Forest[K, W]) (string, error) {
	var digraph bool
	var bu strings.Builder
	var dfs func(K) error
	_, _ = bu.WriteString("{\n")
	vis := make(map[K]struct{})
	dfs = func(u K) error {
		_, _ = fmt.Fprintf(&bu, "%v", u)
		es, err := f.IncidentEdges(u)
		if err != nil {
			return err
		}
		var out []Edge[K, W]
		if digraph {
			for _, e := range es {
				if e.Head == u {
					continue
				}
				if _, ok := vis[e.Key]; !ok {
					out = append(out, e)
				}
			}
		} else {
			for _, e := range es {
				if _, ok := vis[e.Key]; !ok {
					if e.Head == u {
						e.Head, e.Tail = e.Tail, e.Head
					}
					out = append(out, e)
				}
			}
		}
		if len(out) > 1 {
			if digraph {
				_, _ = bu.WriteString("-> {")
			} else {
				_, _ = bu.WriteString("-- {")
			}
			for i, e := range out {
				vis[e.Key] = struct{}{}
				if err := dfs(e.Head); err != nil {
					return err
				}
				if i != len(out)-1 {
					_ = bu.WriteByte(',')
				}
			}
			_ = bu.WriteByte('}')
		} else if len(out) == 1 {
			if digraph {
				_, _ = bu.WriteString("->")
			} else {
				_, _ = bu.WriteString("--")
			}
			vis[out[0].Key] = struct{}{}
			if err := dfs(out[0].Head); err != nil {
				return err
			}
		}
		return nil
	}
	for _, r := range f.AllRoots() {
		digraph = f.IsDirected(r)
		bu.WriteString("{")
		if err := dfs(r); err != nil {
			return "", err
		}
		bu.WriteString("};\n")
	}
	digraph = false
	edges := f.AllEdges()
	var root K
	for len(vis) != len(edges) {
		for _, e := range edges {
			if _, ok := vis[e.Key]; !ok {
				root = e.Tail
				break
			}
		}
		bu.WriteString("{")
		if err := dfs(root); err != nil {
			return "", err
		}
		bu.WriteString("};\n")
	}
	_, _ = bu.WriteString("\n};\n")
	return bu.String(), nil
}

func FmtForestTex[K comparable, W number](f *Forest[K, W]) (string, error) {
	if f == nil || f.Order() == 0 {
		return "", nil
	}
	return fmtForestTex(f)
}
