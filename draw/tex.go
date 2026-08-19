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

package draw

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flxj/graphlib"
)

const (
	gdHead   = `\usetikzlibrary{graphs,graphdrawing}`
	treeHead = `\usegdlibrary{trees}
\tikz [tree layout, sibling distance=8mm]
\graph [nodes={circle, draw, inner sep=1.5pt}]`
)

func TikzDrawGraph[K comparable, W number](g graphlib.Graph[K, W]) (string, error) {
	return "", errors.New("")
}

func TikzDrawDigraph[K comparable, W number](g graphlib.Digraph[K, W]) (string, error) {
	return "", errors.New("")
}

func treeTexStr[K comparable, W number](t *graphlib.Forest[K, W], root K) (string, error) {
	visited := make(map[K]struct{})
	var bu strings.Builder
	var dfs func(K) error
	dfs = func(v K) error {
		_, _ = fmt.Fprintf(&bu, "%v", v)
		visited[v] = struct{}{}
		ns, err := t.Neighbours(v)
		if err != nil {
			return err
		}
		var child []K
		for _, n := range ns {
			_, ok := visited[n.Key]
			if ok {
				continue
			} else {
				child = append(child, n.Key)
			}
		}
		if len(child) > 1 {
			_, _ = bu.WriteString("-- {")
			for i, u := range child {
				if err := dfs(u); err != nil {
					return err
				}
				if i != len(child)-1 {
					_ = bu.WriteByte(',')
				}
			}
			_ = bu.WriteByte('}')
		} else if len(child) == 1 {
			_, _ = bu.WriteString("--")
			if err := dfs(child[0]); err != nil {
				return err
			}
		}
		return nil
	}
	if err := dfs(root); err != nil {
		return "", err
	}
	return bu.String(), nil
}

func ditreeTexStr[K comparable, W number](t *graphlib.Forest[K, W], root K) (string, error) {
	visited := make(map[K]struct{})
	var bu strings.Builder
	var dfs func(K) error
	dfs = func(v K) error {
		_, _ = fmt.Fprintf(&bu, "%v", v)
		visited[v] = struct{}{}
		es, err := t.IncidentEdges(v)
		if err != nil {
			return err
		}
		var out, in []K
		for _, e := range es {
			var o bool
			var u K
			if e.Head == v {
				u = e.Tail
			} else {
				u = e.Head
				o = true
			}
			_, ok := visited[u]
			if ok {
				continue
			} else {
				if o {
					out = append(out, u)
				} else {
					in = append(in, u)
				}
			}
		}
		if len(out)+len(in) == 0 {
			return nil
		}
		if len(out) == 0 {
			if len(in) > 1 {
				_, _ = bu.WriteString("<- {")
				for i, u := range in {
					if err := dfs(u); err != nil {
						return err
					}
					if i != len(in)-1 {
						_ = bu.WriteByte(',')
					}
				}
				_ = bu.WriteByte('}')
			} else if len(in) == 1 {
				_, _ = bu.WriteString("<-")
				if err := dfs(in[0]); err != nil {
					return err
				}
			}
		} else {
			if len(in) == 0 {
				if len(out) > 1 {
					_, _ = bu.WriteString("-> {")
					for i, u := range out {
						if err := dfs(u); err != nil {
							return err
						}
						if i != len(out)-1 {
							_ = bu.WriteByte(',')
						}
					}
					_ = bu.WriteByte('}')
				} else if len(out) == 1 {
					_, _ = bu.WriteString("->")
					if err := dfs(out[0]); err != nil {
						return err
					}
				}
			} else { //TODO
				// out != 0 , in != 0

			}
		}
		return nil
	}
	if err := dfs(root); err != nil {
		return "", err
	}
	return bu.String(), nil
}

func tikzDrawTree[K comparable, W number](t *graphlib.Forest[K, W], root K, bu *strings.Builder, ignoreDirect bool) error {
	var tr string
	var err error
	if t.IsDigraph() && !ignoreDirect {
		tr, err = ditreeTexStr(t, root)
	} else {
		tr, err = treeTexStr(t, root)
	}
	if err != nil {
		return err
	}
	_, _ = bu.WriteString(gdHead)
	_ = bu.WriteByte('\n')
	_, _ = bu.WriteString(treeHead)
	_, _ = bu.WriteString("{\n")
	_, _ = bu.WriteString(tr)
	_, _ = bu.WriteString("\n};\n")
	return nil
}

func TikzDrawTree[K comparable, W number](t *graphlib.Forest[K, W], ignoreDirect bool) (string, error) {
	if t == nil || t.Order() == 0 || t.IsTree() == false {
		return "", nil
	}
	var root K
	r := t.AllRoots()
	if len(r) == 0 {
		v, err := t.RandomVertex()
		if err != nil {
			return "", err
		}
		root = v.Key
	} else {
		root = r[0]
	}
	var bu strings.Builder
	if err := tikzDrawTree(t, root, &bu, ignoreDirect); err != nil {
		return "", err
	}
	return bu.String(), nil
}

func TikzDrawForest[K comparable, W number](f *graphlib.Forest[K, W], ignoreDirect bool) (string, error) {
	if f == nil || f.Order() == 0 {
		return "", nil
	}
	var bu strings.Builder
	for _, r := range f.AllRoots() {
		if err := tikzDrawTree(f, r, &bu, ignoreDirect); err != nil {
			return "", err
		}
		_, _ = bu.WriteString("\\quad\n")
	}
	return bu.String(), errors.New("")
}
