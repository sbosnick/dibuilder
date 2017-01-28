// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/gonum/graph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockNode struct{}

func (m mockNode) ID() int {
	return 1
}

func TestContainerDoesNotHasForeignNodeType(t *testing.T) {
	sut := Container{}
	result := sut.Has(mockNode{})

	assert.False(t, result, "Container has a foreign node")
}

func TestContainerHasMissingNode(t *testing.T) {
	sut := Container{}
	missing := missingNode{container: &sut}
	result := sut.Has(missing)

	assert.True(t, result, "Container did not have a missingNode")
}

func TestContainerDoesNotHasMissingNodeForOtherContainer(t *testing.T) {
	other := Container{}
	missing := missingNode{container: &other}

	sut := Container{}
	result := sut.Has(&missing)

	assert.False(t, result, "Container has missing node for a different container")
}

func TestZeroContainerReturnsOneNode(t *testing.T) {
	sut := &Container{}
	nodes := sut.Nodes()

	assert.Len(t, nodes, 1, "Unexpected number of  nodes")
}

func TestZeroContainerReturnsMissingNode(t *testing.T) {
	sut := &Container{}
	nodes := sut.Nodes()

	assert.IsType(t, missingNode{}, nodes[0], "Unexpected node type")
}

func TestNodeOfZeroContainerHasNoFromNodes(t *testing.T) {
	sut := &Container{}
	node := sut.Nodes()[0]
	fromNodes := sut.From(node)

	assert.Empty(t, fromNodes)
}

func TestNodeOfZeroContainerHsNoToNodes(t *testing.T) {
	sut := &Container{}
	node := sut.Nodes()[0]
	toNodes := sut.To(node)

	assert.Empty(t, toNodes)
}

func TestZeroContainerRootIsError(t *testing.T) {
	sut := &Container{}
	_, err := sut.Root()

	assert.Error(t, err, "expected error not returned")
}

func createRootedContainer() *Container {
	pkg := types.NewPackage("path", "mypackage")
	name := types.NewTypeName(token.NoPos, pkg, "MyIntType", types.Typ[types.Int])

	container := &Container{}
	container.SetRoot(name.Type())

	return container
}

func TestRootedContainerRootIsRootNode(t *testing.T) {
	sut := createRootedContainer()
	root, err := sut.Root()

	assert.NoError(t, err, "unexpected error in getting the container root")
	assert.IsType(t, rootNode{}, root, "unexpected node type for the root node")
}

func TestRootedContainerHasRootNode(t *testing.T) {
	sut := createRootedContainer()
	root := rootNode{container: sut}
	result := sut.Has(root)

	assert.True(t, result, "Container did not have a rootNode")
}

func TestRootedContainerDoesNotHasRootNodeForOtherContainer(t *testing.T) {
	other := &Container{}
	root := rootNode{container: other}

	sut := createRootedContainer()
	result := sut.Has(&root)

	assert.False(t, result, "Container did not had rootNode for other Container")
}

func TestRootedContainerReturnsRootNode(t *testing.T) {
	sut := createRootedContainer()
	nodes := sut.Nodes()

	assert.Contains(t, nodes, rootNode{container: sut}, "Root node is missing.")
}

func TestRootedContainerHasNoNodesFromRoot(t *testing.T) {
	sut := createRootedContainer()
	root, _ := sut.Root()
	fromNodes := sut.From(root)

	assert.Empty(t, fromNodes)
}

func containsNode(nodes []graph.Node, expected graph.Node) assert.Comparison {
	return func() bool {
		for _, node := range nodes {
			if node.ID() == expected.ID() {
				return true
			}
		}

		return false
	}
}

func getNodeIDs(nodes []graph.Node) []int {
	var ids []int

	for _, node := range nodes {
		ids = append(ids, node.ID())
	}

	return ids
}

func TestRootedContinerHasRootNodeFromMissingNode(t *testing.T) {
	sut := createRootedContainer()
	missing := missingNode{container: sut}
	fromNodes := sut.From(missing)

	root, _ := sut.Root()
	assert.Condition(t, containsNode(fromNodes, root),
		"No edge from missingNode to rootNode: %v does not contain %v",
		getNodeIDs(fromNodes), root.ID())
}

func TestRootedContainerHasMissingNodeToRootNode(t *testing.T) {
	sut := createRootedContainer()
	root, _ := sut.Root()
	toNodes := sut.To(root)

	missing := missingNode{container: sut}
	assert.Condition(t, containsNode(toNodes, missing),
		"No edge from missingNode to rootNode: %v does not contain %v",
		getNodeIDs(toNodes), missing.ID())
}

func TestRootedContarainerHasNoNodesToMissingNode(t *testing.T) {
	sut := createRootedContainer()
	missing := missingNode{container: sut}
	toNodes := sut.To(missing)

	assert.Empty(t, toNodes, "Unexpected edge to the missingNode")
}

func TestRootedContainerHasEdgeFromMissingNodeToRootNode(t *testing.T) {
	sut := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()

	assert.True(t, sut.HasEdgeFromTo(missing, root), "No edge from missingNode to rootNode")
}

func TestRootedContainerDoesNotHaveEdgeFromRootNoodToMissingNode(t *testing.T) {
	sut := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()

	assert.False(t, sut.HasEdgeFromTo(root, missing), "No edge from rootNode to missingNode")
}

func TestRootedContainerHasEdgeBetweenMissingNodeAndRootNode(t *testing.T) {
	sut := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()

	assert.True(t, sut.HasEdgeBetween(missing, root), "No edge between missingNode and rootNode")
	assert.True(t, sut.HasEdgeBetween(root, missing), "No edge between missingNode and rootNode")
}

func TestRootedContainerEdgeFromMissingNodeToRootNodeNotNil(t *testing.T) {
	sut := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()
	edge := sut.Edge(missing, root)

	assert.NotNil(t, edge, "Edge from missingNode to rootNode was nil")
}

func TestRootedContainerEdgeFromRootNodeToMissingNodeIsNil(t *testing.T) {
	sut := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()
	edge := sut.Edge(root, missing)

	assert.Nil(t, edge, "Edge from rootNode to missingNode was not nil")
}

func TestRootedContainerEdgeHoldsExpectedFromAndTo(t *testing.T) {
	sut := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()
	edge := sut.Edge(missing, root)

	require.NotNil(t, edge.From(), "Unexpected nil From node in the Edge")
	require.NotNil(t, edge.To(), "Unexpected nil To node in the Edge")
	assert.Equal(t, missing.ID(), edge.From().ID(), "Unexpected From node in the Edge")
	assert.Equal(t, root.ID(), edge.To().ID(), "Unexpected To node in the Edge")
}

func TestContainerAddFuncAddsFuncNode(t *testing.T) {
	function := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	err := sut.AddFunc(function)

	require.NoError(t, err, "Unexpected error returned from AddFunc")
	require.Len(t, sut.nodes, 1, "AddFunc did not add node to Container")
	assert.IsType(t, &funcNode{}, sut.nodes[0], "Node added by AddFunc had unexpected type")
}

func TestConainerHasFuncNodeAfterAddFunc(t *testing.T) {
	function := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	_ = sut.AddFunc(function)
	result := sut.Has(sut.nodes[0])

	assert.True(t, result, "Container did not Has() node added by AddFunc")
}

func TestContainerNodesIncludesFuncNodeAfterAddFunc(t *testing.T) {
	function := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	_ = sut.AddFunc(function)
	nodes := sut.Nodes()

	assert.Condition(t, hasFuncNodeForFunction(nodes, function),
		"Nodes() did not contain a funcNode for the expected function")

}

func hasFuncNodeForFunction(nodes []graph.Node, function *types.Func) func() bool {
	return func() bool {
		for _, node := range nodes {
			if fn, ok := node.(*funcNode); ok {
				if fn.function == function {
					return true
				}
			}
		}

		return false
	}
}
