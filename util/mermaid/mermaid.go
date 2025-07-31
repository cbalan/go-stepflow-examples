package mermaid

import (
	"github.com/cbalan/go-stepflow"
	"github.com/cbalan/go-stepflow-examples/util/graph"
	"github.com/cbalan/go-stepflow/core"
	"maps"
	"slices"
	"strings"
)

func FromSteps(stepsSpec *stepflow.StepsSpec) (string, error) {
	scope, transitions, err := stepflow.Transitions(stepsSpec)
	if err != nil {
		return "", err
	}

	return FromTransitions(scope, transitions), nil
}

func FromTransitions(root core.Scope, transitions []core.Transition) string {
	scopeTree, eventsGraph := NewGraphFromTransitions(root, transitions)
	return MermaidDiagramFromGraph(scopeTree, eventsGraph)
}

func NewGraphFromTransitions(root core.Scope, transitions []core.Transition) (*graph.Tree[core.Scope], *graph.Graph[core.Event, string]) {
	scopeTree := graph.NewTree[core.Scope](root.Name(), root)
	eventsGraph := graph.NewGraph[core.Event, string]()

	for _, t := range transitions {
		source := t.Source()
		destinations := t.PossibleDestinations()

		if source == nil {
			continue
		}

		if source.Scope().Parent() != nil {
			scopeTree.AddChild(source.Scope().Parent().Name(), source.Scope().Name(), source.Scope())
		}

		eventsGraph.AddNode(eventString(source), source)

		for _, destination := range destinations {
			eventsGraph.AddNode(eventString(destination.Event()), destination.Event())
			eventsGraph.AddEdge(eventString(source), eventString(destination.Event()), destination.Reason())

			if destination.Event().Scope().Parent() != nil {
				scopeTree.AddChild(destination.Event().Scope().Parent().Name(), destination.Event().Scope().Name(), destination.Event().Scope())
			}
		}
	}

	return scopeTree, eventsGraph
}

func MermaidDiagramFromGraph(scopesTree *graph.Tree[core.Scope], eventsGraph *graph.Graph[core.Event, string]) string {
	rootScope := scopesTree.RootNode().Data
	var sb strings.Builder

	sb.WriteString("stateDiagram-v2\n")

	// wire in top level start and end states.
	sb.WriteString("[*] --> " + eventString(core.StartCommand(rootScope)) + "\n")
	sb.WriteString(eventString(core.CompletedEvent(rootScope)) + " --> [*]\n")

	// subgraphs
	subGraph(&sb, scopesTree, scopesTree.RootNode())

	for _, source := range slices.Sorted(maps.Keys(eventsGraph.Edges)) {
		for _, destination := range slices.Sorted(maps.Keys(eventsGraph.Edges[source])) {
			label := eventsGraph.Edges[source][destination]
			sb.WriteString(source + " --> " + destination + ": " + label + "\n")
		}
	}

	return sb.String()
}

func subGraph(sb *strings.Builder, tree *graph.Tree[core.Scope], node *graph.Node[core.Scope]) {
	sb.WriteString("state " + node.Name + " {\n")

	sortedChildren := slices.SortedFunc(tree.Children(node.Name), func(n *graph.Node[core.Scope], n2 *graph.Node[core.Scope]) int {
		return strings.Compare(n.Name, n2.Name)
	})

	// write children subgraph
	for _, child := range sortedChildren {
		subGraph(sb, tree, child)
	}

	// TODO: list scope events in place of hardcoding these
	sb.WriteString(eventString(core.StartCommand(node.Data)) + "\n")
	sb.WriteString(eventString(core.CompletedEvent(node.Data)) + "\n")

	sb.WriteString("}\n")
}

func eventString(event core.Event) string {
	return event.Name() + "&colon;" + event.Scope().Name()
}
