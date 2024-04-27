package graphlib

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

type PropertyName int

type GraphProperty[T any] struct {
	Name  PropertyName
	Value T
}

const (
	PropertyDigraph PropertyName = iota
	PropertyAcyclic
	PropertySimple
	PropertyRegular
	PropertyConnected
	PropertyForest
	PropertyLoop
	PropertyCompleted
	PropertyTree
	PropertyNegativeWeight
	PropertyGraphName
	PropertyOrder
	PropertySize
	PropertyMaxDegree
	PropertyMinDegree
	PropertyAvgDegree
)

// Graph [K, V, W] represents the graph object,
// K represents the identification type of vertices and edges,
// V represents the data type of vertices,
// and W represents the weight data type of edges.
//
// In a mathematician's terminology, a graph is a collection of points
// and lines connecting some (possibly empty) subset of them.
// The points of a graph are most commonly known as graph vertices,
// but may also be called "nodes" or simply "points."
// Similarly, the lines connecting the vertices of a graph
// are most commonly known as graph edges, but may also be called "arcs" or "lines."
//
// The concept of graph can be referenced:
// https://mathworld.wolfram.com/Graph.html
type Graph[K comparable, V any, W number] interface {
	//
	// The name of current graph object.
	Name() string
	//
	// Change the name of graph.
	SetName(name string)
	//
	// The number of vertexes in the graph.
	Order() int
	//
	// The number of edges in the graph.
	Size() int
	//
	// Is it a directed graph.
	IsDigraph() bool
	//
	// Query the properties of the current graph.
	// The current default Graph implementation in Graphlib provides
	// calculations for the following properties
	//
	// PropertyGraph: Is it a directed graph
	// PropertyAcyclic: Is it an acyclic graph
	//
	// PropertySimple: Is it a simple graph?
	// For the definition of a simple graph,
	// please refer to https://mathworld.wolfram.com/SimpleGraph.html
	//
	// PropertyRegular: Is it a regular graph?
	// For the definition of a regular graph,
	// please refer to https://mathworld.wolfram.com/RegularGraph.html
	//
	// PropertyConnected: Is it a connected graph?
	// For the definition of connectivity,
	// please refer to https://mathworld.wolfram.com/ConnectedGraph.html
	//
	// PropertyForest: Is it a forest?
	// For the definition of forest,
	// please refer to https://mathworld.wolfram.com/Forest.html
	//
	// PropertyCompleted: Is it a complete graph?
	// For the definition of a complete graph,
	// please refer to https://mathworld.wolfram.com/CompleteGraph.html
	//
	// PropertyTree: Is it a tree?
	// For the definition of a tree,
	// please refer to https://mathworld.wolfram.com/Tree.html
	//
	// PropertyLoop: Does it include a loop.
	//
	// PropertyNegativeWeight: Does it contain negative weight edges
	//
	// PropertyGraphName: Graph name
	//
	// PropertyOrder: The number of vertices in the graph
	//
	// PropertySize: The number of edges in a graph
	//
	// PropertyMaxDegree: Maximum Degree
	//
	// PropertyMinDegree: Minimum Read
	//
	// PropertyAvgDegree: Average degree (float64)
	Property(p PropertyName) (GraphProperty[any], error)
	//
	// The unordered set of all vertices in a graph.
	AllVertexes() ([]Vertex[K, V], error)
	//
	// The unordered set of all edges in the graph.
	AllEdges() ([]Edge[K, W], error)
	//
	// Add vertices to the graph.
	AddVertex(vertex Vertex[K, V]) error
	//
	// Delete a vertex, and all edges corresponding to that vertex
	// will also be deleted. If the vertex does not exist, return an error.
	RemoveVertex(key K) error
	//
	// Add new edge, if the corresponding vertex of the
	// edge does not exist, return an error.
	AddEdge(edge Edge[K, W]) error
	//
	// Delete specified edge.
	RemoveEdgeByKey(key K) error
	//
	// Delete edges with endpoints ednpoint1 and endpoint2.
	// If it is a directed graph, delete all arcs in the 'endpoint1->endpoint2'
	// and 'endpoint2->endpoint1' directions simultaneously.
	RemoveEdge(endpoint1, endpoint2 K) error
	//
	// Calculate the degree of vertices.
	// If it is a directed graph, calculate the sum of in degree and out degree.
	// If the vertex does not exist, an error is returned.
	Degree(vertex K) (int, error)
	//
	// Query the adjacent vertices of a specified vertex.
	Neighbours(vertex K) ([]Vertex[K, V], error)
	//
	// Query specified vertex.
	GetVertex(key K) (Vertex[K, V], error)
	//
	// Query all edges with endpoints 1 and 2 as their respective endpoints.
	GetEdge(endpoint1, endpoint2 K) ([]Edge[K, W], error)
	//
	// Query specified edge.
	GetEdgeByKey(key K) (Edge[K, W], error)
	//
	// Filter vertices based on label information,
	// and eligible vertices need to include all label items in
	// the label parameter simultaneously.
	GetVertexesByLabel(labels map[string]string) ([]Vertex[K, V], error)
	//
	// Filter edges based on label information,
	// and eligible edges need to include all label items
	// in the label parameter simultaneously.
	GetEdgesByLabel(labels map[string]string) ([]Edge[K, W], error)
	//
	// Update vertex data.
	SetVertexValue(key K, value V) error
	//
	// Update vertex label.
	SetVertexLabel(key K, labelKey, labelVal string) error
	//
	// Remove vertex label.
	DeleteVertexLabel(key K, labelKey string) error
	//
	// Update edge data.
	SetEdgeValueByKey(key K, value any) error
	//
	// Update dege label.
	SetEdgeLabelByKey(key K, labelKey, labelVal string) error
	//
	// Remove edge label.
	DeleteEdgeLabelByKey(key K, labelKey string) error
	//
	// Update edge data. If there are multiple edges associated with
	// endpoints1 and endpoint2 simultaneously,
	// the data of these edges will be updated simultaneously.
	SetEdgeValue(endpoint1, endpoint2 K, value any) error
	//
	// Update edge label. If there are multiple edges associated with
	// endpoints1 and endpoint2 simultaneously,
	// the label of these edges will be updated simultaneously.
	SetEdgeLabel(endpoint1, endpoint2 K, labelKey, labelVal string) error
	//
	// Delete edge label. If there are multiple edges associated with
	// endpoints1 and endpoint2 simultaneously,
	// the label of these edges will be updated simultaneously.
	DeleteEdgeLabel(endpoint1, endpoint2 K, labelKey string) error
	//
	// Copy the current graph.
	Clone() (Graph[K, V, W], error)
}

type Vertex[K comparable, V any] struct {
	// The unique identifier of this vertex.
	Key K `json:"key" yaml:"key"`
	// The data object of this vertex .
	Value V `json:"value" yaml:"value"`
	// The label of this vertex.
	Labels map[string]string `json:"labels" yaml:"labels"`
}

func (v *Vertex[K, V]) Clone() *Vertex[K, V] {
	vv := &Vertex[K, V]{
		Key:   v.Key,
		Value: v.Value,
	}
	if v.Labels != nil {
		vv.Labels = make(map[string]string)
		for k, l := range v.Labels {
			vv.Labels[k] = l
		}
	}
	return vv
}

type Edge[K comparable, W number] struct {
	// The unique identifier of this edge.
	// use a key to distinguish edge,
	// because maybe exists multiedge in graph
	Key K `json:"key" yaml:"key"`
	// One endpoint of an edge.
	Head K `json:"head" yaml:"head"`
	// The other endpoint of an edge.
	// For a directed graph, the direction of the edge is head ->tail.
	Tail K `json:"tail" yaml:"tail"`
	// Edge weight.
	Weight W `json:"weight" yaml:"weight"`
	// Edge data.
	Value any `json:"value" yaml:"value"`
	// Edge labels.
	Labels map[string]string `json:"labels" yaml:"labels"`
}

func (e *Edge[K, W]) Clone() *Edge[K, W] {
	ee := &Edge[K, W]{
		Key:   e.Key,
		Head:  e.Head,
		Tail:  e.Tail,
		Value: e.Value,
	}
	if e.Labels != nil {
		ee.Labels = make(map[string]string)
		for k, l := range e.Labels {
			ee.Labels[k] = l
		}
	}
	return ee
}
