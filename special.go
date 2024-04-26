package graphlib

type Bipartite[K comparable, V any, W number] struct {
	g *graph[K,V,W]
	partA map[K]bool 
	partB map[K]bool 
}

func NewBipartite[K comparable, V any, W number](digraph bool,name string)(*Bipartite[K,V,W],error){
	g,err:= newGraph[K,V,W](digraph,name)
	if err!=nil{
		return nil, err 
	}
	return &Bipartite[K,V,W]{
		g:g,
		partA:make(map[K]bool),
		partB:make(map[K]bool),
	},nil 
}


