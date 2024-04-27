package graphlib

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

// GraphInfo represents the basic information of a graph,
// used for serialization of graph objects.
type GraphInfo[K comparable, V any, W number] struct {
	Name     string         `json:"name" yaml:"name"`
	Digraph  bool           `json:"digraph" yaml:"digraph"`
	Vertexes []Vertex[K, V] `json:"vertexes" yaml:"vertexes"`
	Edges    []Edge[K, W]   `json:"edges" yaml:"edges"`
}

// Serialize Graph in JSON format.
func MarshalGraphToJSON[K comparable, V any, W number](g Graph[K, V, W]) ([]byte, error) {
	var (
		err error
		vs  []Vertex[K, V]
		es  []Edge[K, W]
	)
	if vs, err = g.AllVertexes(); err != nil {
		return nil, err
	}
	if es, err = g.AllEdges(); err != nil {
		return nil, err
	}

	gi := GraphInfo[K, V, W]{
		Name:     g.Name(),
		Digraph:  g.IsDigraph(),
		Vertexes: vs,
		Edges:    es,
	}

	return json.Marshal(gi)
}

// Serialize Graph in yaml format.
func MarshalGraphToYaml[K comparable, V any, W number](g Graph[K, V, W]) ([]byte, error) {
	var (
		err error
		vs  []Vertex[K, V]
		es  []Edge[K, W]
	)
	if vs, err = g.AllVertexes(); err != nil {
		return nil, err
	}
	if es, err = g.AllEdges(); err != nil {
		return nil, err
	}

	gi := GraphInfo[K, V, W]{
		Name:     g.Name(),
		Digraph:  g.IsDigraph(),
		Vertexes: vs,
		Edges:    es,
	}

	return yaml.Marshal(gi)
}

func UnmarshalGraph[K comparable, V any, W number](s []byte) (Graph[K, V, W], error) {
	gi := GraphInfo[K, V, W]{}
	if json.Valid(s) {
		if err := json.Unmarshal(s, &gi); err != nil {
			return nil, err
		}
	} else {
		if err := yaml.Unmarshal(s, &gi); err != nil {
			return nil, err
		}
	}
	g, err := NewGraph[K, V, W](gi.Digraph, gi.Name)
	if err != nil {
		return nil, err
	}
	for _, v := range gi.Vertexes {
		if err = g.AddVertex(v); err != nil {
			return nil, err
		}
	}
	for _, e := range gi.Edges {
		if err = g.AddEdge(e); err != nil {
			return nil, err
		}
	}
	return g, nil
}

func UnmarshalDigraph[K comparable, V any, W number](s []byte) (Digraph[K, V, W], error) {
	gi := GraphInfo[K, V, W]{}
	if json.Valid(s) {
		if err := json.Unmarshal(s, &gi); err != nil {
			return nil, err
		}
	} else {
		if err := yaml.Unmarshal(s, &gi); err != nil {
			return nil, err
		}
	}
	g, err := NewDigraph[K, V, W](gi.Name)
	if err != nil {
		return nil, err
	}
	for _, v := range gi.Vertexes {
		if err = g.AddVertex(v); err != nil {
			return nil, err
		}
	}
	for _, e := range gi.Edges {
		if err = g.AddEdge(e); err != nil {
			return nil, err
		}
	}
	return g, nil
}
