package graphlib

type number interface {
	int | int64
}

type Graph[K comparable, V any, W number] interface {
	IsDigraph() bool
	IsSimple() bool
	IsRegular() bool
	IsMulti() bool
	IsAcyclic() bool
	IsConnected() bool
	IsCompleted() bool
	IsTree() bool 
	IsForest() bool 
	HasLoop()bool

	Order() int
	Size() int
	MinDegree() int
	MaxDegree() int
	AvgDegree() int

	AllVertexes()([]Vertex[K,V],error)
	AllEdges()([]Edge[K,W],error)

	AddVertex(vertex Vertex[K, V]) error
	RemoveVertex(key K) error
	AddEdge(edge Edge[K, W]) error
	RemoveEdgeByKey(key K) error
	RemoveEdge(endpoint1, endpoint2 K) error

	Degree(vertex K) int
	Neighbours(vertex K) []Vertex[K, V]
	GetVertex(key K) (Vertex[K, V], error)
	GetEdge(key K) (Edge[K, W], error)
	GetVerticesByLabel(labels map[string]string) ([]Vertex[K, V], error)
	GetEdgesByLabel(labels map[string]string) ([]Vertex[K, W], error)

	SetVertexLabel(key K, labelKey, labelVal string) error
	DeleteVertexLabel(key K, labelKey string) error
	SetEdgeLabel(key K, labelKey, labelVal string) error
	DeleteEdgeLabel(key K, labelKey string) error

	//Diameter()int
	//Radius()int

	Clone() (Graph[K, V, W], error)
}

type Vertex[K comparable, V any] struct {
	Key    K
	Value  V
	Lables map[string]string
}

type Edge[K comparable, W number] struct {
	Key    K
	Head   K
	Tail   K
	Weight W
	Labels map[string]string
}
