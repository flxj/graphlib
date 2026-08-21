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
	"strings"

	"github.com/flxj/graphlib"
)

const (
	gdHead = `\usetikzlibrary{graphs,graphdrawing}`

	treeHead = `\usegdlibrary{trees}
\tikz [tree layout, sibling distance=8mm]
\graph [nodes={circle, draw, inner sep=1.5pt}]`

	ditreeHead = `\usegdlibrary {trees,layered,force}
\tikz \graph [tree layout,nodes={draw,circle}]`
	diHead = `\usegdlibrary {force}
\graph [random seed=10, spring layout]`
	gHead = `\usegdlibrary {force}
\graph [random seed=10, spring layout]`
)

func TikzDrawGraph[K comparable, W number](g graphlib.Graph[K, W]) (string, error) {
	tr, err := graphlib.FmtGraphTex(g)
	if err != nil {
		return "", err
	}
	var bu strings.Builder
	_, _ = bu.WriteString(gdHead)
	_ = bu.WriteByte('\n')
	_, _ = bu.WriteString(gHead)

	_, _ = bu.WriteString(tr)

	return bu.String(), nil
}

func TikzDrawDigraph[K comparable, W number](g graphlib.Digraph[K, W]) (string, error) {
	tr, err := graphlib.FmtDigraphTex(g)
	if err != nil {
		return "", err
	}
	var bu strings.Builder
	_, _ = bu.WriteString(gdHead)
	_ = bu.WriteByte('\n')
	_, _ = bu.WriteString(diHead)

	bu.WriteString(tr)

	return bu.String(), nil
}

func TikzDrawForest[K comparable, W number](f *graphlib.Forest[K, W]) (string, error) {
	if f == nil || f.Order() == 0 {
		return "", nil
	}
	var bu strings.Builder

	tr, err := graphlib.FmtForestTex(f)
	if err != nil {
		return "", err
	}
	_, _ = bu.WriteString(gdHead)
	_ = bu.WriteByte('\n')
	_, _ = bu.WriteString(treeHead)
	_, _ = bu.WriteString(tr)
	return bu.String(), nil
}
