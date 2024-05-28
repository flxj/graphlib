package draw

import (
	"fmt"
	"testing"

	"github.com/flxj/graphlib"
)

func TestDraw(t *testing.T) {
	g, err := graphlib.NewGraph[int, int, int](false, "test-g")
	if err != nil {
		fmt.Printf("new graph error:%v\n", err)
		return
	}

	vs := []graphlib.Vertex[int, int]{
		{Key: 1, Value: 1},
		{Key: 2, Value: 2},
		{Key: 3, Value: 3},
		{Key: 4, Value: 4},
		{Key: 5, Value: 5},
		{Key: 6, Value: 6},
	}
	for _, v := range vs {
		_ = g.AddVertex(v)
	}
	_ = g.SetVertexLabel(1, "color", "green")
	_ = g.SetVertexLabel(6, "color", "red")

	es := []graphlib.Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2, Weight: 5},
		{Key: 2, Head: 2, Tail: 3, Weight: 6},
		{Key: 3, Head: 5, Tail: 6, Weight: 7},
		{Key: 4, Head: 4, Tail: 5, Weight: 8},
		{Key: 5, Head: 2, Tail: 5, Weight: 9},
	}
	for _, e := range es {
		_ = g.AddEdge(e)
	}
	_ = g.SetEdgeLabelByKey(3, "color", "red")

	file, err := RenderHTML(g, true, "/tmp")
	if err != nil {
		fmt.Printf("draw error:%v\n", err)
		return
	}
	fmt.Println(file)
}
