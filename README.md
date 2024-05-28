# graphlib

`Graphlib` is a graph data structure generic library implemented in Golang, providing definitions and basic operations for undirected/directed graphs, as well as built-in common graph algorithms. Additionally, as a feature, graphlib also comes with a DAG based goroutine workflow engine (ExecGraph/Workflow).😀

### Features

✔️ **Basic operation of the graph:**

* Create undirected/directed graphs
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

✔️ **Visualization:**

* Support graphical display of Graph objects (based on [D3](https://d3js.org))


### Getting started

```shell
go get github.com/flxj/graphlib
```

> Currently, Graphlib is in the process of development and testing, and some features are not yet fully developed. Please do not use for production environments now.

Create an undirected graph using the following example 👇


v1---v2
|   /
|  /   
v3      v4-----v5----v6


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
		_ = g.AddVertex(v)
	}

	es := []graphlib.Edge[int, int]{
		{Key: 1, Head: "v1", Tail: "v2"},
		{Key: 2, Head: "v1", Tail: "v3"},
		{Key: 3, Head: "v2", Tail: "v3"},
		{Key: 4, Head: "v4", Tail: "v5"},
		{Key: 5, Head: "v5", Tail: "v6"},
	}
	for _, e := range es {
		_ = g.AddEdge(e)
	}
	
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	ps, _ := g.Property(graphlib.PropertySimple)
	fmt.Printf("simple:%v\n", ps.Value)

	pc, _ := g.Property(graphlib.PropertyConnected)
	fmt.Printf("connected:%v\n", pc.Value)

	pa, _ := g.Property(graphlib.PropertyAcyclic)
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
		_ = g.AddVertex(v)
	}
	
	es := []Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2},
		{Key: 2, Head: 2, Tail: 3},
		{Key: 3, Head: 5, Tail: 6},
		{Key: 4, Head: 4, Tail: 5},
		{Key: 5, Head: 2, Tail: 5},
	}
	for _, e := range es {
		_ = g.AddEdge(e)
	}

	fmt.Println(gs)
	fmt.Printf("order:%d\n", g.Order())
	fmt.Printf("size:%d\n", g.Size())

	p, _ := g.Property(PropertyConnected)
	fmt.Printf("connected:%v\n", p.Value)

	p, _ = g.Property(PropertyUnilateralConnected)
	fmt.Printf("unidirectional connected:%v\n", p.Value)

	p, _ = g.Property(PropertyAcyclic)
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



+----------------------------+
|  job1---->job2 --.         |
|                   \        |
|                    \       |
|                     V      |
|  job3---->job4---> job5    |
|                            |
+----------------------------+


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
		_ =g.AddJob(k,j)
	}

	deps:=[][]int{
		{1,2},
		{3,4},
		{2,5},
		{4,5},
	}
	for _,d:=range deps {
		_ =g.AddDependency(d[0],d[1])
	}

	v1 = 100
	v2 = 200 
	var val = 2*(v1+100) + 3*v2-10

	_=g.Start()

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

### Visualization

Support graphical display of Graph objects (based on [D3](https://d3js.org)). 

Calling the RenderHTML method will generate an HTML file about the given Graph in the specified directory, with the file name consistent with the Graph name. Open the file in the browser to view the graphical graph object.

```golang
import(
	"fmt"
    
	"github.com/flxj/graphlib"
)

func main(){
    g, err := graphlib.NewGraph[int, int, int](false,"test-g")
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
	_ = g.SetVertexLabel(1,"color","green")
	_ = g.SetVertexLabel(6,"color","red")
	
	es := []graphlib.Edge[int, int]{
		{Key: 1, Head: 1, Tail: 2,Weight:5},
		{Key: 2, Head: 2, Tail: 3,Weight:6},
		{Key: 3, Head: 5, Tail: 6,Weight:7},
		{Key: 4, Head: 4, Tail: 5,Weight:8},
		{Key: 5, Head: 2, Tail: 5,Weight:9},
	}
	for _, e := range es {
		_ = g.AddEdge(e)
	}
	_ = g.SetEdgeLabelByKey(3,"color","red")

	file,err:=RenderHTML(g,false,"/tmp")
	if err!=nil{
		fmt.Printf("draw error:%v\n", err)
		return
	}
	fmt.Println(file)
}
```

