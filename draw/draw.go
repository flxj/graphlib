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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/flxj/graphlib"
)

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type d3Node struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type d3Link struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Weight string `json:"weight"`
	Color  string `json:"color"`
}

type d3NetworkData struct {
	Nodes []*d3Node `json:"nodes"`
	Links []*d3Link `json:"links"`
}

type d3Network[K comparable, V any, W number] struct {
	Digraph     bool                    `json:"digraph"`
	RandomColor bool                    `json:"random_color"`
	ShowWeight  bool                    `json:"show_weight"`
	Data        *d3NetworkData          `json:"data"`
	Vertexes    []graphlib.Vertex[K, V] `json:"vertexes"`
	Edges       []graphlib.Edge[K, W]   `json:"edges"`
}

func getHTMLTemplate(digraph bool) (string, error) {
	if !digraph {
		return graphHTML, nil
	}
	if digraphHTML != "" {
		return digraphHTML, nil
	}
	f, err := os.OpenFile("./digraph.tpl", os.O_RDWR, 0666)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	//
	tpl, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	digraphHTML = string(tpl)

	return digraphHTML, nil
}

func RenderHTML[K comparable, V any, W number](g graphlib.Graph[K, V, W], showWeight bool, dir string) (string, error) {
	var (
		err  error
		vs   []graphlib.Vertex[K, V]
		es   []graphlib.Edge[K, W]
		data d3NetworkData
	)

	if vs, err = g.AllVertexes(); err != nil {
		return "", err
	}
	if es, err = g.AllEdges(); err != nil {
		return "", err
	}
	//
	for _, v := range vs {
		k := fmt.Sprintf("%v", v.Key)
		node := &d3Node{
			ID:    k,
			Name:  k,
			Color: "",
		}
		if v.Labels != nil {
			node.Color = v.Labels["color"]
		}
		data.Nodes = append(data.Nodes, node)
	}
	//
	for _, e := range es {
		l := &d3Link{
			ID:     fmt.Sprintf("%v", e.Key),
			Source: fmt.Sprintf("%v", e.Head),
			Target: fmt.Sprintf("%v", e.Tail),
			Weight: fmt.Sprintf("%v", e.Weight),
		}
		if e.Labels != nil {
			l.Color = e.Labels["color"]
		}
		data.Links = append(data.Links, l)
	}
	//
	net := &d3Network[K, V, W]{
		Digraph:    g.IsDigraph(),
		ShowWeight: showWeight,
		Data:       &data,
		//Vertexes:vs,
		//Edges:es,
	}

	var (
		bs      []byte
		f       *os.File
		tpl     *template.Template
		htmlTpl string
	)
	if htmlTpl, err = getHTMLTemplate(net.Digraph); err != nil {
		return "", err
	}
	if tpl, err = template.New(g.Name()).Parse(htmlTpl); err != nil {
		return "", err
	}
	if f, err = os.OpenFile(fmt.Sprintf("%s/%s.html", dir, g.Name()), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	if bs, err = json.Marshal(net); err != nil {
		return "", err
	}
	if err = tpl.Execute(f, string(bs)); err != nil {
		return "", err
	}

	_ = f.Sync()

	return f.Name(), nil
}

type graphDOT struct {
	Digraph  bool
	Attr     []string
	NodeAttr []string
	EdgeAttr []string
	Nodes    []string
	Edges    []string
}

func getDOT[K comparable, V any, W number](g graphlib.Graph[K, V, W], vertexShape string, attr []string, showWeight bool) ([]byte, error) {
	var (
		err error
		vs  []graphlib.Vertex[K, V]
		es  []graphlib.Edge[K, W]
	)

	if vs, err = g.AllVertexes(); err != nil {
		return nil, err
	}
	if es, err = g.AllEdges(); err != nil {
		return nil, err
	}

	dot := &graphDOT{
		Digraph: g.IsDigraph(),
		Attr:    attr,
	}
	for _, v := range vs {
		attrs := []string{
			"shape=ellipse",
			fmt.Sprintf("label=%v", v.Key),
		}
		if vertexShape != "" {
			attrs[0] = "shape=" + vertexShape
		}
		if v.Labels != nil {
			c := v.Labels["color"]
			if c != "" {
				attrs = append(attrs, "color="+c)
			}
		}
		node := fmt.Sprintf("%v [%s]", v.Key, strings.Join(attrs, ","))
		dot.Nodes = append(dot.Nodes, node)
	}
	for _, e := range es {
		attrs := []string{}
		if e.Labels != nil {
			c := e.Labels["color"]
			if c != "" {
				attrs = append(attrs, "color="+c)
			}
		}
		if showWeight {
			fmt.Println("Weight: ", e.Weight)
			attrs = append(attrs, fmt.Sprintf("label=%v", e.Weight))
		}
		var edge string
		if dot.Digraph {
			edge = fmt.Sprintf("%v->%v [%s]", e.Head, e.Tail, strings.Join(attrs, ","))
		} else {
			edge = fmt.Sprintf("%v--%v [%s]", e.Head, e.Tail, strings.Join(attrs, ","))
		}
		dot.Edges = append(dot.Edges, edge)
	}

	tpl, err := template.New(g.Name()).Parse(dotTemplate)
	if err != nil {
		return nil, err
	}

	var bs []byte
	var buf = bytes.NewBuffer(bs)
	if err = tpl.Execute(buf, dot); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func GetDOT[K comparable, V any, W number](g graphlib.Graph[K, V, W], vertexShape string, attr []string, showWeight bool) ([]byte, error) {
	return getDOT(g, vertexShape, attr, showWeight)
}

func RenderSVG[K comparable, V any, W number](g graphlib.Graph[K, V, W], vertexShape string, attr []string, showEdgeWeight bool, dir string) (string, error) {
	var (
		f   *os.File
		err error
		dot []byte
	)
	if dot, err = getDOT(g, vertexShape, attr, showEdgeWeight); err != nil {
		return "", err
	}

	if f, err = os.OpenFile(fmt.Sprintf("%s/%s.dot", dir, g.Name()), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	if _, err = f.Write(dot); err != nil {
		return "", err
	}
	_ = f.Sync()

	svg := fmt.Sprintf("%s/%s.svg", dir, g.Name())
	cmd := exec.Command("bash", "-c", "dot", "-Tsvg", f.Name(), ">", svg)

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return svg, nil
}
