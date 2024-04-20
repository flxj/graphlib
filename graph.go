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
	PropertyNegativeWeight
	PropertyGraphName
	PropertyOrder
	PropertySize
	PropertyMaxDegree
	PropertyMinDegree
	PropertyAvgDegree
)

type Graph[K comparable, V any, W number] interface {
	Name() string
	SetName(name string)

	Order() int
	Size() int
	MinDegree() int
	MaxDegree() int
	AvgDegree() float64
	IsDigraph() bool
	IsSimple() bool
	IsRegular() bool
	IsAcyclic() bool
	IsConnected() bool
	IsCompleted() bool
	IsTree() bool
	IsForest() bool
	HasLoop() bool
	Property(p PropertyName) (GraphProperty[any], error)

	AllVertexes() ([]Vertex[K, V], error)
	AllEdges() ([]Edge[K, W], error)

	AddVertex(vertex Vertex[K, V]) error
	RemoveVertex(key K) error
	AddEdge(edge Edge[K, W]) error
	RemoveEdgeByKey(key K) error
	RemoveEdge(endpoint1, endpoint2 K) error

	Degree(vertex K) (int, error)
	Neighbours(vertex K) ([]Vertex[K, V], error)
	GetVertex(key K) (Vertex[K, V], error)
	GetEdge(endpoint1, endpoint2 K) ([]Edge[K, W], error)
	GetEdgeByKey(key K) (Edge[K, W], error)
	GetVertexesByLabel(labels map[string]string) ([]Vertex[K, V], error)
	GetEdgesByLabel(labels map[string]string) ([]Edge[K, W], error)

	SetVertexValue(key K, value V) error
	SetVertexLabel(key K, labelKey, labelVal string) error
	DeleteVertexLabel(key K, labelKey string) error

	SetEdgeValueByKey(key K, value V) error
	SetEdgeLabelByKey(key K, labelKey, labelVal string) error
	DeleteEdgeLabelByKey(key K, labelKey string) error
	SetEdgeValue(endpoint1, endpoint2 K, value V) error
	SetEdgeLabel(endpoint1, endpoint2 K, labelKey, labelVal string) error
	DeleteEdgeLabel(endpoint1, endpoint2 K, labelKey string) error

	Clone() (Graph[K, V, W], error)
}

type Vertex[K comparable, V any] struct {
	Key    K
	Value  V
	Labels map[string]string
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
	Key    K // use a key to distinguish edge, because maybe exists multiedge in graph
	Head   K
	Tail   K
	Weight W
	Value  any
	Labels map[string]string
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
