package graph

import "iter"

type Node[T any] struct {
	Name string
	Data T
}

type Graph[N, E any] struct {
	Nodes map[string]*Node[N]
	Edges map[string]map[string]E
}

func NewGraph[N, E any]() *Graph[N, E] {
	return &Graph[N, E]{Nodes: make(map[string]*Node[N]), Edges: make(map[string]map[string]E)}
}

func (g *Graph[N, E]) AddNode(name string, data N) {
	// Ignore node if already set.
	if _, ok := g.Nodes[name]; ok {
		return
	}

	g.Nodes[name] = &Node[N]{Name: name, Data: data}
}

func (g *Graph[N, E]) AddEdge(source string, destination string, data E) {
	if _, ok := g.Edges[source]; !ok {
		g.Edges[source] = make(map[string]E)
	}

	g.Edges[source][destination] = data
}

type Tree[N any] struct {
	Root  string
	Graph *Graph[N, bool]
}

func NewTree[N any](name string, node N) *Tree[N] {
	graph := NewGraph[N, bool]()
	graph.AddNode(name, node)
	return &Tree[N]{Root: name, Graph: graph}
}

func (t *Tree[N]) RootNode() *Node[N] {
	return t.Graph.Nodes[t.Root]
}

func (t *Tree[N]) AddChild(parent string, name string, data N) {
	t.Graph.AddNode(name, data)
	t.Graph.AddEdge(parent, name, true)
}

func (t *Tree[T]) Children(parent string) iter.Seq[*Node[T]] {
	return func(yield func(node *Node[T]) bool) {
		for child := range t.Graph.Edges[parent] {
			yield(t.Graph.Nodes[child])
		}
	}
}
