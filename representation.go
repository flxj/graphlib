package graphlib

import (
	"fmt"
)

// Create weight matrix for graph.
func NewWeightMatrix[K comparable, W number](g Graph[K, any, W]) (*WeightMatrix[K, W], error) {
	p, err := g.Property(PropertySimple)
	if err != nil {
		return nil, err
	}
	if !p.Value.(bool) {
		return nil, errNotSimple
	}

	var (
		vs []Vertex[K, any]
		es []Edge[K, W]
	)

	if vs, err = g.AllVertexes(); err != nil {
		return nil, err
	}
	if es, err = g.AllEdges(); err != nil {
		return nil, err
	}

	none := any(MaxFloatDistance).(W)
	wm := &WeightMatrix[K, W]{
		none:     none,
		vertexes: make([]K, len(vs)),
		data:     make([][]W, len(vs)),
	}
	idx := make(map[K]int)
	for i, v := range vs {
		idx[v.Key] = i
		wm.vertexes[i] = v.Key
		wm.data[i] = make([]W, len(vs))
		for j := 0; j < len(vs); j++ {
			wm.data[i][j] = none
			if i == j {
				wm.data[i][j] = any(0).(W)
			}
		}
	}
	//
	for _, e := range es {
		i := idx[e.Head]
		j := idx[e.Tail]
		wm.data[i][j] = e.Weight
		if !g.IsDigraph() {
			wm.data[j][i] = e.Weight
		}
	}

	return wm, nil
}

// Create adjacency matrix for graph.
func NewAdjacencytMatrix[K comparable, W number](g Graph[K, any, W]) (*AdjacencyMatrix[K], error) {
	var (
		err error
		vs  []Vertex[K, any]
		es  []Edge[K, W]
	)
	if vs, err = g.AllVertexes(); err != nil {
		return nil, err
	}
	if es, err = g.AllEdges(); err != nil {
		return nil, err
	}

	am := &AdjacencyMatrix[K]{
		vertexes: make([]K, len(vs)),
		data:     make([][]int, len(vs)),
	}
	idx := make(map[K]int)
	for i, v := range vs {
		idx[v.Key] = i
		am.vertexes[i] = v.Key
		am.data[i] = make([]int, len(vs))
	}
	//
	for _, e := range es {
		i := idx[e.Head]
		j := idx[e.Tail]
		am.data[i][j] = 1
		if !g.IsDigraph() {
			am.data[j][i] = 1
		}
	}

	return am, nil
}

func NewDegreeMatrix[K comparable, W number](g Graph[K, any, W]) (*DegreeMatrix[K], error) {
	return nil, errNotImplement
}

type AdjacencyMatrix[K comparable] struct {
	vertexes []K
	data     [][]int
}

func (m *AdjacencyMatrix[K]) Matrix() [][]int {
	return m.data
}

func (m *AdjacencyMatrix[K]) Columns() []K {
	return m.vertexes
}

type DegreeMatrix[K comparable] struct {
	vertexes []K
	data     [][]int
}

func (m *DegreeMatrix[K]) Weight() [][]int {
	return m.data
}

func (m *DegreeMatrix[K]) Columns() []K {
	return m.vertexes
}

type WeightMatrix[K comparable, W number] struct {
	none     W
	vertexes []K
	data     [][]W
}

func (m *WeightMatrix[K, W]) Weight() [][]W {
	return m.data
}

func (m *WeightMatrix[K, W]) Distance(infinite float64) [][]float64 {
	w := make([][]float64, len(m.data))
	for i, d := range m.data {
		ds := make([]float64, len(d))
		for j, p := range d {
			if p == m.none {
				ds[j] = infinite
			} else {
				ds[j] = any(p).(float64)
			}
		}
		w[i] = ds
	}
	return w
}

func (m *WeightMatrix[K, W]) Columns() []K {
	return m.vertexes
}

type endpoint[K comparable, W number] struct {
	key    K // vertex key
	edge   K // edge key
	weight W
	next   *endpoint[K, W]
	//prev   *endpoint[K, W]
}

type edge[K comparable, W number] struct {
	key    K
	head   K
	tail   K
	weight W
}

type adjacencyList[K comparable, W number] struct {
	digraph bool
	outAdj  map[K]*endpoint[K, W] // adjacency list
	inAdj   map[K]*endpoint[K, W] // contrary adjacency list
}

func newAdjacencyLis[K comparable, W number](digraph bool) (*adjacencyList[K, W], error) {
	adj := &adjacencyList[K, W]{
		digraph: digraph,
		outAdj:  make(map[K]*endpoint[K, W]),
	}
	if adj.digraph {
		adj.inAdj = make(map[K]*endpoint[K, W])
	}
	return adj, nil
}

func newAdjacencyListFromGraph[K comparable, V any, W number](g Graph[K, V, W]) (*adjacencyList[K, W], error) {
	var (
		err error
		vs  []Vertex[K, V]
		es  []Edge[K, W]
		adj *adjacencyList[K, W]
	)
	if vs, err = g.AllVertexes(); err != nil {
		return nil, err
	}
	if es, err = g.AllEdges(); err != nil {
		return nil, err
	}
	if adj, err = newAdjacencyLis[K, W](g.IsDigraph()); err != nil {
		return nil, err
	}

	for _, v := range vs {
		if err = adj.addVertexes(v.Key); err != nil {
			return nil, err
		}
	}
	for _, e := range es {
		if err = adj.addEdge(e.Head, e.Tail, e.Key, e.Weight); err != nil {
			return nil, err
		}
	}
	return adj, nil
}

func (l *adjacencyList[K, W]) addVertexes(vs ...K) error {
	for _, v := range vs {
		if _, ok := l.outAdj[v]; !ok {
			l.outAdj[v] = nil
		}
		if l.digraph {
			if _, ok := l.inAdj[v]; !ok {
				l.inAdj[v] = nil
			}
		}
	}
	return nil
}

func (l *adjacencyList[K, W]) delVertex(v K) error {
	del := func(v K, adj map[K]*endpoint[K, W]) {
		delete(adj, v)
		for k, p := range adj {
			var head = p
			var prev = &endpoint[K, W]{next: head}
			//prev.next = head

			for q := head; q != nil; {
				if q.key == v {
					if q == head {
						// remove head element
						prev = q
						q = q.next
						head = q
						prev.next = nil
					} else {
						//
						prev.next = q.next
						q = q.next
					}
				} else {
					prev = q
					q = q.next
				}
			}
			if head != p {
				adj[k] = head
			}
		}
	}
	del(v, l.outAdj)
	if l.digraph {
		del(v, l.inAdj)
	}
	return nil
}

func (l *adjacencyList[K, W]) delVertexes(vs ...K) error {
	for _, v := range vs {
		if _, ok := l.outAdj[v]; !ok {
			return fmt.Errorf("vertex %v not exists", v)
		}
	}
	for _, v := range vs {
		if err := l.delVertex(v); err != nil {
			return err
		}
	}
	return nil
}

func (l *adjacencyList[K, W]) addEdge(head, tail, key K, weight W) error {
	insert := func(v1, v2, edge K, w W, adj map[K]*endpoint[K, W]) error {
		p, ok := adj[v1]
		if !ok {
			return fmt.Errorf("vertex %v not exists", v1)
		}
		var exists bool
		for q := p; q != nil; q = q.next {
			if q.key == v2 && q.edge == edge {
				q.weight = w
				exists = true
				break
			}
		}
		if !exists {
			q := &endpoint[K, W]{
				key:    v2,
				edge:   edge,
				weight: w,
			}
			if p != nil {
				q.next = p
			}
			adj[v1] = q
		}
		return nil
	}
	// insert to outAdj
	if err := insert(head, tail, key, weight, l.outAdj); err != nil {
		return err
	}
	if l.digraph {
		// insert to inAdj
		if err := insert(tail, head, key, weight, l.inAdj); err != nil {
			return err
		}
	} else {
		if err := insert(tail, head, key, weight, l.outAdj); err != nil {
			return err
		}
	}
	return nil
}

func (l *adjacencyList[K, W]) delEdge(head, tail, key K) error {
	del := func(v1, v2, edge K, adj map[K]*endpoint[K, W]) error {
		p, ok := adj[v1]
		if !ok {
			return fmt.Errorf("vertex %v not exists", v1)
		}
		if p == nil {
			return fmt.Errorf("edge %v not exists", edge)
		}
		var prev = &endpoint[K, W]{next: p}
		var q *endpoint[K, W]
		for e := p; e != nil; e = e.next {
			if e.key == v2 && e.edge == edge {
				q = e
				break
			}
			prev = e
		}
		if q != nil {
			// remove head element of list.
			if q == p {
				adj[v1] = q.next
			} else {
				prev.next = q.next
				q.next = nil
			}
			return nil
		}
		return fmt.Errorf("edge %v not exists", edge)
	}
	//
	if err := del(head, tail, key, l.outAdj); err != nil {
		return err
	}
	if l.digraph {
		if err := del(tail, head, key, l.inAdj); err != nil {
			return err
		}
	} else {
		if err := del(tail, head, key, l.outAdj); err != nil {
			return err
		}
	}
	return nil
}

func (l *adjacencyList[K, W]) addEdges(es ...*edge[K, W]) error {
	for _, e := range es {
		if _, ok := l.outAdj[e.head]; !ok {
			return fmt.Errorf("vertex %v not exists", e.head)
		}
		if _, ok := l.outAdj[e.tail]; !ok {
			return fmt.Errorf("vertex %v not exists", e.tail)
		}
	}
	for _, e := range es {
		if err := l.addEdge(e.head, e.tail, e.key, e.weight); err != nil {
			return err
		}
	}
	return nil
}

func (l *adjacencyList[K, W]) delEdges(es ...*edge[K, W]) error {
	for _, e := range es {
		if err := l.delEdge(e.head, e.tail, e.key); err != nil {
			return err
		}
	}
	return nil
}

func (l *adjacencyList[K, W]) degree(v K) (int, error) {
	d, err := l.outDegree(v)
	if err != nil {
		return 0, err
	}
	if l.digraph {
		in, err := l.inDegree(v)
		if err != nil {
			return 0, err
		}
		d += in
	}
	return d, nil
}

func (l *adjacencyList[K, W]) outDegree(v K) (int, error) {
	p, ok := l.outAdj[v]
	if !ok {
		return 0, fmt.Errorf("vertex %v not exists", v)
	}
	var d int
	for q := p; q != nil; q = q.next {
		d++
	}
	return d, nil
}

func (l *adjacencyList[K, W]) inDegree(v K) (int, error) {
	var adj map[K]*endpoint[K, W]
	if l.digraph {
		adj = l.inAdj
	} else {
		adj = l.outAdj
	}
	p, ok := adj[v]
	if !ok {
		return 0, fmt.Errorf("vertex %v not exists", v)
	}
	var d int
	for q := p; q != nil; q = q.next {
		d++
	}
	return d, nil
}

func (l *adjacencyList[K, W]) neighbours(v K) ([]K, error) {
	ks := make(map[K]struct{})
	p, ok := l.outAdj[v]
	if !ok {
		return nil, fmt.Errorf("vertex %v not exists", v)
	}
	//
	for q := p; q != nil; q = q.next {
		ks[q.key] = struct{}{}
	}
	if l.digraph {
		p, ok = l.inAdj[v]
		if !ok {
			return nil, fmt.Errorf("vertex %v not exists", v)
		}
		for q := p; q != nil; q = q.next {
			ks[q.key] = struct{}{}
		}
	}
	var ns []K
	for k := range ks {
		ns = append(ns, k)
	}
	return ns, nil
}

func (l *adjacencyList[K, W]) inNeighbours(v K) ([]K, error) {
	var adj map[K]*endpoint[K, W]
	if l.digraph {
		adj = l.inAdj
	} else {
		adj = l.outAdj
	}
	ks := make(map[K]struct{})
	p, ok := adj[v]
	if !ok {
		return nil, fmt.Errorf("vertex %v not exists", v)
	}
	for q := p; q != nil; q = q.next {
		ks[q.key] = struct{}{}
	}
	var ns []K
	for k := range ks {
		ns = append(ns, k)
	}
	return ns, nil
}

func (l *adjacencyList[K, W]) outNeighbours(v K) ([]K, error) {
	ks := make(map[K]struct{})
	p, ok := l.outAdj[v]
	if !ok {
		return nil, fmt.Errorf("vertex %v not exists", v)
	}
	for q := p; q != nil; q = q.next {
		ks[q.key] = struct{}{}
	}
	var ns []K
	for k := range ks {
		ns = append(ns, k)
	}
	return ns, nil
}

func (l *adjacencyList[K, W]) inEdges(v K) ([]K, error) {
	var adj map[K]*endpoint[K, W]
	if l.digraph {
		adj = l.inAdj
	} else {
		adj = l.outAdj
	}
	var ks []K
	p, ok := adj[v]
	if !ok {
		return nil, fmt.Errorf("vertex %v not exists", v)
	}
	for q := p; q != nil; q = q.next {
		ks = append(ks, q.edge)
	}
	return ks, nil
}

func (l *adjacencyList[K, W]) outEdges(v K) ([]K, error) {
	p, ok := l.outAdj[v]
	if !ok {
		return nil, fmt.Errorf("vertex %v not exists", v)
	}
	var ks []K
	for q := p; q != nil; q = q.next {
		ks = append(ks, q.edge)
	}
	return ks, nil
}

func (l *adjacencyList[K, W]) sources() ([]K, error) {
	if !l.digraph {
		return nil, errNotDigraph
	}
	var vs []K
	for k, v := range l.inAdj {
		if v == nil {
			vs = append(vs, k)
		}
	}
	return vs, nil
}

func (l *adjacencyList[K, W]) sinks() ([]K, error) {
	if !l.digraph {
		return nil, errNotDigraph
	}
	var vs []K
	for k, v := range l.outAdj {
		if v == nil {
			vs = append(vs, k)
		}
	}
	return vs, nil
}

func (l *adjacencyList[K, W]) minDegree() (int, error) {
	minD := len(l.outAdj)
	for v := range l.outAdj {
		d, err := l.degree(v)
		if err != nil {
			return 0, err
		}
		if d < minD {
			minD = d
		}
	}
	return minD, nil
}

func (l *adjacencyList[K, W]) maxDegree() (int, error) {
	maxD := len(l.outAdj)
	for v := range l.outAdj {
		d, err := l.degree(v)
		if err != nil {
			return 0, err
		}
		if d > maxD {
			maxD = d
		}
	}
	return maxD, nil
}

func (l *adjacencyList[K, W]) avgDegree() (float64, error) {
	if len(l.outAdj) == 0 {
		return 0.0, nil
	}
	var sumD int
	for v := range l.outAdj {
		d, err := l.degree(v)
		if err != nil {
			return 0, err
		}
		sumD += d
	}
	return float64(sumD) / float64(len(l.outAdj)), nil
}

func (l *adjacencyList[K, W]) isAcyclic() (bool, error) {
	if len(l.outAdj) == 0 {
		return true, nil
	}
	neighbour := l.neighbours
	if l.digraph {
		neighbour = l.outNeighbours
	}

	var start K
	for k := range l.outAdj {
		start = k
		break
	}
	//
	visited := make(map[K]bool)
	prev := make(map[K]K)

	stack := newStack[K]()
	stack.push(start)
	for !stack.empty() {
		v, _ := stack.pop()
		if _, ok := visited[v]; !ok {
			visited[v] = true
		}
		vs, err := neighbour(v)
		if err != nil {
			return false, err
		}
		for _, k := range vs {
			// loop
			if k == v {
				return false, nil
			}
			// exclude the parent vertex that visited just now.
			if k != prev[v] {
				// find a back edge
				if _, ok := visited[k]; ok {
					return false, nil
				} else {
					stack.push(k)
					prev[k] = v
				}
			}
		}
		// to dfs another components.
		if stack.empty() && len(visited) < len(l.outAdj) {
			for k := range l.outAdj {
				if _, ok := visited[k]; !ok {
					stack.push(k)
					break
				}
			}
		}
	}
	return true, nil
}

func (l *adjacencyList[K, W]) isConnected() (bool, error) {
	if len(l.outAdj) == 0 {
		return false, nil
	}
	// bfs
	var start K
	for k := range l.outAdj {
		start = k
		break
	}
	visited := make(map[K]bool)
	que := newFIFO[K]()
	que.push(start)

	for !que.empty() {
		v, _ := que.pop()
		if _, ok := visited[v]; !ok {
			visited[v] = true
		}
		vs, err := l.neighbours(v)
		if err != nil {
			return false, err
		}
		for _, v := range vs {
			if _, ok := visited[v]; !ok {
				que.push(v)
			}
		}
	}
	if len(visited) != len(l.outAdj) {
		return false, nil
	}
	return true, nil
}

func (l *adjacencyList[K, W]) isSimple() (bool, error) {
	if l.digraph {
		for k, v := range l.outAdj {
			tails := make(map[K]int)
			for p := v; p != nil; p = p.next {
				// loop
				if p.key == k {
					return false, nil
				}
				//
				t := tails[p.key]
				if t >= 1 {
					return false, nil
				} else {
					tails[p.key] = t + 1
					in := l.inAdj[k]
					for q := in; q != nil; q = q.next {
						if q.key == p.key {
							return false, nil
						}
					}
				}
			}
		}
		return true, nil
	}
	//
	for k, v := range l.outAdj {
		vs := make(map[K]struct{})
		for p := v; p != nil; p = p.next {
			if p.key == k {
				return false, nil
			}
			if _, ok := vs[p.key]; ok {
				return false, nil
			}
			vs[p.key] = struct{}{}
		}
	}
	return true, nil
}

func (l *adjacencyList[K, W]) isRegular() (bool, error) {
	d := -1
	for k := range l.outAdj {
		n, err := l.degree(k)
		if err != nil {
			return false, err
		}
		if d >= 0 {
			if n != d {
				return false, nil
			}
		} else {
			d = n
		}
	}
	return true, nil
}

func (l *adjacencyList[K, W]) isForest() (bool, error) {
	return l.isAcyclic()
}

func (l *adjacencyList[K, W]) hasLoop() (bool, error) {
	for k, v := range l.outAdj {
		for p := v; p != nil; p = p.next {
			if p.key == k {
				return true, nil
			}
		}
	}
	return false, nil
}

func (l *adjacencyList[K, W]) hasNegativeWeight() (bool, error) {
	for _, v := range l.outAdj {
		for p := v; p != nil; p = p.next {
			if p.weight < 0 {
				return true, nil
			}
		}
	}
	if l.digraph {
		for _, v := range l.inAdj {
			for p := v; p != nil; p = p.next {
				if p.weight < 0 {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

func (l *adjacencyList[K, W]) property(p int) (property[bool], error) {
	var r bool
	var err error
	switch p {
	case acyclic:
		r, err = l.isAcyclic()
	case connected:
		r, err = l.isConnected()
	case simple:
		r, err = l.isSimple()
	case regular:
		r, err = l.isRegular()
	case forest:
		r, err = l.isForest()
	case negativeWeight:
		r, err = l.hasNegativeWeight()
	case loop:
		r, err = l.hasLoop()
	default:
		err = errUnknownProperty
	}
	if err != nil {
		return property[bool]{}, err
	}
	return property[bool]{
		name:  p,
		value: r,
	}, nil
}
