# graphlib

`Graphlib` is a graph data structure generic library implemented by Golang, providing definitions and basic operations for undirected/directed graphs, as well as built-in common graph algorithms. Additionally, as a feature, graphlib also comes with a DAG based goroutine workflow engine (ExecGraph).😀

### Features

✔️ **Basic operation of the graph:**

* Create undirected/directed graphs (supports simple and multiple graphs)
* Serialization and Deserialization of Graph Objects (JSON/YAML Format)
* Dynamically adjust vertex (increase/decrease/modify attributes)
* Dynamically adjust edge (increase/decrease/modify attributes)
* Basic properties of computational graphs: connectivity, acyclicity, etc

✔️ **Graph calculation:**
   
* Generate induced subgraph
* Generate spanning subgraph
* Minimum Spanning Tree
* Calculate strongly connected components
* Algebraic operations on graphs (intersection/union/difference/sum/product)
* Construct matrix representations of graphs (adjacency matrix, degree matrix, weight matrix)

✔️ **Graph algorithm:**

* Graph traversal (BFS, DFS)
* Shortest path (single source, multiple sources, negative weight)
* Calculate maximum flow
* Bipartite matching
* Topological sorting
* Vertex colouring/edge colouring
  
✔️ **DAG:**

* Support for a directed acyclic graph based goroutine workflow engine


### Getting started

```shell
go get github.com/flxj/graphlib
```

> Currently, Graphlib is in the process of development and testing, and some features are not yet fully developed. Please do not use versions below 1.0 for production environments

Create an undirected graph using the following example 👇

```shell
v1---v2
|   /
|  /   
v3         v4-----v5----v6
```

```golang
import(
	"fmt"
    
	"github.com/flxj/graphlib"
)

func main() {
    g, err := graphlib.NewGraph[string, int, int](false, "graph")
	if err != nil {
		fmt.Println(err)
		return
	}

	vs := []graphlib.Vertex[int, int]{
		{Key: "v1", Value: 1},
		{Key: "v2", Value: 2},
		{Key: "v3", Value: 3},
		{Key: "v4", Value: 4},
		{Key: "v5", Value: 5},
		{Key: "v6", Value: 6},
	}

	for _, v := range vs {
		if err := g.AddVertex(v); err != nil {
			fmt.Printf("add vertex error:%v\n", err)
			return
		}
	}

	es := []graphlib.Edge[int, int]{
		{Key: 1, Head: "v1", Tail: "v2"},
		{Key: 2, Head: "v1", Tail: "v3"},
		{Key: 3, Head: "v2", Tail: "v3"},
		{Key: 4, Head: "v4", Tail: "v5"},
		{Key: 5, Head: "v5", Tail: "v6"},
	}

	for _, e := range es {
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}
	
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	ps, err := g.Property(graphlib.PropertySimple)
	if err != nil {
		fmt.Printf("get property simple error:%v\n", err)
		return
	}
	fmt.Printf("simple:%v\n", ps.Value)
	pc, err := g.Property(graphlib.PropertyConnected)
	if err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", pc.Value)
	pa, err := g.Property(graphlib.PropertyAcyclic)
	if err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", pa.Value)
}
```
The output is as follows:
```shell
order:6
size:5
simple:true
connected:false
acyclic:false
```


Create a directed graph using the following example 👇

```shell
1----> 2 ---> 3
       |
       v
4----> 5 ---> 6
```

```golang
import(
	"fmt"
    
	"github.com/flxj/graphlib"
)

func main(){
    g, err := graphlib.NewDigraph[int, int, int]("g")
	if err != nil {
		fmt.Printf("new graph error:%v\n", err)
		return
	}

	vs := []Vertex[int, int]{
		{Key: 1, Value: 1},
		{Key: 2, Value: 2},
		{Key: 3, Value: 3},
		{Key: 4, Value: 4},
		{Key: 5, Value: 5},
		{Key: 6, Value: 6},
	}
	for _, v := range vs {
		if err := g.AddVertex(v); err != nil {
			fmt.Printf("add vertex error:%v\n", err)
			return
		}
	}
	
	es := []Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2},
		{Key: 2, Head: 2, Tail: 3},
		{Key: 3, Head: 5, Tail: 6},
		{Key: 4, Head: 4, Tail: 5},
		{Key: 5, Head: 2, Tail: 5},
	}
	for _, e := range es {
		if err := g.AddEdge(e); err != nil {
			fmt.Printf("add edge error:%v\n", err)
			return
		}
	}

	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())
	p, err := g.Property(PropertyConnected)
	if err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("connected:%v\n", p.Value)
	if p, err = g.Property(PropertyUnilateralConnected); err != nil {
		fmt.Printf("get property connected error:%v\n", err)
		return
	}
	fmt.Printf("unidirectional connected:%v\n", p.Value)
	if p, err = g.Property(PropertyAcyclic);err != nil {
		fmt.Printf("get property acyclic error:%v\n", err)
		return
	}
	fmt.Printf("acyclic:%v\n", p.Value)
}
```

```shell
order:6
size:5
connected:true
unidirectional connected:false
acyclic:true
```


### Workflow 

Graphlib also provides an ExecGraph object that represents a goroutine (Job) execution process arranged in a directed acyclic graph logic, conceptually similar to any other workflow system, such as [argo-workflows](https://argo-workflows.readthedocs.io/en/latest/).

Users can add tasks to the ExecGraph object and set dependencies between tasks. Users can manage the entire workflow declaration cycle through the ExecGraph interface method.

The following example shows how to create, run, and wait for ExecGraph:

```shell
job1---->job2--.
                \
                 \
                  v
job3---->job4--->job5
```

```golang
import(
	"fmt"
    
	"github.com/flxj/graphlib"
)

func main() {
	g,err:= graphlib.NewExecGraph[int,graphlib.Job]("exec")
	if err!=nil{
		fmt.Printf("[ERR] create exec graph error: %v\n",err)
		return 
	}
    var (
		v1 int
		v2 int 
		v3 int 
	)
	// input:  v1 <- x, v2 <- y
	// output: v3 <- 2*(x+100) + 3*x-10

	job1:=func() error {
		v1 += 100
		return nil 
	}
	job2:=func() error {
		v1 = 2*v1 
		return nil 
	}
	job3:=func() error {
		v2 = 3*v2
		return nil 
	}
	job4:=func() error {
		v2 = v2-10
		return nil 
	}
	job5:=func() error {
		v3 = v1+v2
		return nil
	}

	jobs:=map[int]Job{
		1:job1,
		2:job2,
		3:job3,
		4:job4,
		5:job5,
	}
	for k,j:=range jobs {
		if err:=g.AddJob(k,j);err!=nil{
			fmt.Printf("[ERR] add job error: %v\n",err)
			return 
		}
	}

	deps:=[][]int{
		{1,2},
		{3,4},
		{2,5},
		{4,5},
	}
	for _,d:=range deps {
		if err:=g.AddDependency(d[0],d[1]);err!=nil{
			fmt.Printf("[ERR] add dep error: %v\n",err)
			return 
		}
	}

	v1 = 100
	v2 = 200 
	var val = 2*(v1+100) + 3*v2-10

	if err:=g.Start();err!=nil{
		fmt.Printf("[ERR] start graph error: %v\n",err)
		return 
	}

	if err:=g.Wait();err!=nil{
		fmt.Printf("[ERR] wait graph error: %v\n",err)
		return 
	}

	if v3 != val {
		fmt.Printf("exec err: expect %d, actual get %d\n",val,v3)
	}else{
		fmt.Println("success")
	}
}
```