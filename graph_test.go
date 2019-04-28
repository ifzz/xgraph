package xgraph // import "github.com/orkestr8/xgraph"

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/graph/topo"
)

func TestGonumGraph(t *testing.T) {

	g := simple.NewDirectedGraph()

	// AddNode must be called right after NewNode to ensure the ID is properly assigned and registered
	// in the graph, or we'd get ID collision panic.
	a := g.NewNode()
	require.Nil(t, g.Node(a.ID()))
	g.AddNode(a)
	require.NotNil(t, g.Node(a.ID()))

	b := g.NewNode()
	g.AddNode(b)

	aLikesB := g.NewEdge(a, b)
	require.Nil(t, g.Edge(a.ID(), b.ID()))
	g.SetEdge(aLikesB)

	cycle := topo.DirectedCyclesIn(g)
	require.Equal(t, 0, len(cycle))

	// Calling ReversedEdge doesn't actually reverses the edge in the graph.
	reversed := aLikesB.ReversedEdge()
	require.Nil(t, g.Edge(b.ID(), a.ID()))
	require.NotNil(t, g.Edge(a.ID(), b.ID()))

	// Now an edge exists.  For this DAG we have a loop now.
	g.SetEdge(reversed)
	require.NotNil(t, g.Edge(b.ID(), a.ID()))
	require.NotNil(t, g.Edge(a.ID(), b.ID()))

	_, err := topo.SortStabilized(g, nil)
	require.Error(t, err)

	cycle = topo.DirectedCyclesIn(g)
	require.Equal(t, 1, len(cycle))
	t.Log(cycle)

	c := g.NewNode()
	g.AddNode(c)
	g.SetEdge(g.NewEdge(a, c))
	g.SetEdge(g.NewEdge(c, a))
	cycle = topo.DirectedCyclesIn(g)
	require.Equal(t, 2, len(cycle))
	t.Log(cycle)
}

type nodeT struct {
	id string
}

func (n *nodeT) NodeKey() NodeKey {
	return NodeKey(n.id)
}

func (n *nodeT) String() string {
	return n.id
}

func TestAdd(t *testing.T) {

	A := &nodeT{id: "A"}
	B := &nodeT{id: "B"}
	C := &nodeT{id: "C"}
	plus := &nodeT{id: "+"}
	minus := &nodeT{id: "-"}

	g := New(Options{})
	require.NoError(t, g.Add(A, B, C, plus, minus))

	require.NoError(t, g.Add(A), "Idempotent: same node by identity")
	require.NoError(t, g.Add(&nodeT{id: "A"}), "OK for duplicate key when struct identity fails")

	for _, n := range []Node{plus, minus, A, B, C} {
		require.True(t, g.Has(n))
	}
}

func TestAssociate(t *testing.T) {

	A := &nodeT{id: "A"}
	B := &nodeT{id: "B"}
	C := &nodeT{id: "C"}
	D := &nodeT{id: "D"}

	g := New(Options{})
	require.NoError(t, g.Add(A, B, C))

	require.True(t, g.Has(A))
	require.True(t, g.Has(B))
	require.True(t, g.Has(C))
	require.False(t, g.Has(D))

	likes := EdgeKind(1)
	shares := EdgeKind(2)

	_, err := g.Associate(A, likes, B)
	require.NoError(t, err)
	require.True(t, g.Edge(A, likes, B))

	_, err = g.Associate(D, likes, A)
	require.Error(t, err, "Expects error because D was not added to the graph.")
	require.False(t, g.Edge(D, likes, A), "Expects false because C is not part of the graph.")

	_, err = g.Associate(A, likes, C)
	require.NoError(t, err, "No error because A and C are members of the graph.")
	require.True(t, g.Edge(A, likes, C), "A likes C.")
	require.False(t, g.Edge(C, shares, A), "Shares is not an association kind between A and B.")

}
