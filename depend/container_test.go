// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/types"
	"testing"

	"github.com/cheekybits/is"
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
	is := is.New(t)

	sut := Container{}
	missing := findMissingNode(sut.Nodes())
	result := sut.Has(missing)

	is.OK(result)
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
	is := is.New(t)

	sut := &Container{}
	nodes := sut.Nodes()

	is.OK(findMissingNode(nodes))
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

func TestRootedContainerRootIsRootNode(t *testing.T) {
	is := is.New(t)

	sut, _ := createRootedContainer()
	root, err := sut.Root()

	is.NoErr(err)
	is.OK(root, func() { _ = root.(*rootNode) })
}

func TestRootedContainerHasErrorOnSetRoot(t *testing.T) {
	is := is.New(t)

	sut, _ := createRootedContainer()
	err := sut.SetRoot(types.Typ[types.Uint8])

	is.Err(err)
}

func TestRootedContainerHasRootNode(t *testing.T) {
	sut, _ := createRootedContainer()
	root := rootNode{container: sut}
	result := sut.Has(root)

	assert.True(t, result, "Container did not have a rootNode")
}

func TestRootedContainerDoesNotHasRootNodeForOtherContainer(t *testing.T) {
	other := &Container{}
	root := rootNode{container: other}

	sut, _ := createRootedContainer()
	result := sut.Has(&root)

	assert.False(t, result, "Container did not had rootNode for other Container")
}

func TestRootedContainerReturnsRootNode(t *testing.T) {
	is := is.New(t)

	sut, typ := createRootedContainer()
	nodes := sut.Nodes()

	rootnode := findRootNode(nodes)
	otherrootnode, _ := sut.Root()
	is.OK(rootnode, rootnode.(*rootNode).root == typ)
	is.Equal(rootnode, otherrootnode)
}

func TestRootedContainerHasNoNodesFromRoot(t *testing.T) {
	sut, _ := createRootedContainer()
	root, _ := sut.Root()
	fromNodes := sut.From(root)

	assert.Empty(t, fromNodes)
}

func TestRootedContinerHasRootNodeFromMissingNode(t *testing.T) {
	is := is.New(t)

	sut, _ := createRootedContainer()
	missing := findMissingNode(sut.Nodes())
	fromNodes := sut.From(missing)

	is.OK(findRootNode(fromNodes))
}

func TestRootedContainerHasMissingNodeToRootNode(t *testing.T) {
	is := is.New(t)

	sut, _ := createRootedContainer()
	root, _ := sut.Root()
	toNodes := sut.To(root)

	is.OK(findMissingNode(toNodes))
}

func TestRootedContarainerHasNoNodesToMissingNode(t *testing.T) {
	sut, _ := createRootedContainer()
	missing := missingNode{container: sut}
	toNodes := sut.To(missing)

	assert.Empty(t, toNodes, "Unexpected edge to the missingNode")
}

func TestRootedContainerHasEdgeFromMissingNodeToRootNode(t *testing.T) {
	is := is.New(t)

	sut, _ := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()
	hasedge := sut.HasEdgeFromTo(missing, root)

	is.OK(hasedge)
}

func TestRootedContainerDoesNotHaveEdgeFromRootNoodToMissingNode(t *testing.T) {
	sut, _ := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()

	assert.False(t, sut.HasEdgeFromTo(root, missing), "No edge from rootNode to missingNode")
}

func TestRootedContainerHasEdgeBetweenMissingNodeAndRootNode(t *testing.T) {
	sut, _ := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()

	assert.True(t, sut.HasEdgeBetween(missing, root), "No edge between missingNode and rootNode")
	assert.True(t, sut.HasEdgeBetween(root, missing), "No edge between missingNode and rootNode")
}

func TestRootedContainerEdgeFromMissingNodeToRootNodeNotNil(t *testing.T) {
	sut, _ := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()
	edge := sut.Edge(missing, root)

	assert.NotNil(t, edge, "Edge from missingNode to rootNode was nil")
}

func TestRootedContainerEdgeFromRootNodeToMissingNodeIsNil(t *testing.T) {
	sut, _ := createRootedContainer()
	missing := missingNode{container: sut}
	root, _ := sut.Root()
	edge := sut.Edge(root, missing)

	assert.Nil(t, edge, "Edge from rootNode to missingNode was not nil")
}

func TestRootedContainerEdgeHoldsExpectedFromAndTo(t *testing.T) {
	sut, _ := createRootedContainer()
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
	is := is.New(t)
	function := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	_ = sut.AddFunc(function)
	nodes := sut.Nodes()

	is.OK(findFuncNodeForFunction(nodes, function))

}

func TestRootedContainerWithFuncProvidingRootDoesNotHaveMissingEdge(t *testing.T) {
	sut, typ := createRootedContainer()
	sut.AddFunc(makeFunc(nil, typ, false))
	nodes := sut.From(missingNode{container: sut})

	assert.Len(t, nodes, 0, "Unexpected nodes from missingNode")
}

func TestContainerWithUnsatifiedFuncRequirmentHasMissingEdge(t *testing.T) {
	is := is.New(t)
	function := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	_ = sut.AddFunc(function)
	nodes := sut.From(missingNode{container: sut})

	is.OK(findFuncNodeForFunction(nodes, function))
}

func TestContainerWithFuncProvidingRootHasEdgeFromFunc(t *testing.T) {
	is := is.New(t)

	sut, typ := createRootedContainer()
	function := makeFunc(nil, typ, false)
	sut.AddFunc(function)
	node := findFuncNodeForFunction(sut.Nodes(), function)
	nodes := sut.From(node)

	rootnode := findRootNode(nodes)
	is.OK(rootnode)
}

func TestContainerWithFuncProvidingRequiredFuncHasEdgeFromFunc(t *testing.T) {
	is := is.New(t)
	function1 := makeFunc(nil, types.Typ[types.Int], false)
	function2 := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	sut.AddFunc(function1)
	sut.AddFunc(function2)
	node1 := findFuncNodeForFunction(sut.Nodes(), function1)
	is.OK(node1)
	nodes := sut.From(node1)

	is.OK(findFuncNodeForFunction(nodes, function2))
}

func TestContainerWithUnsatifiedFuncRequirmentHasMissingEdgeToFunc(t *testing.T) {
	is := is.New(t)
	function := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	sut.AddFunc(function)
	node := findFuncNodeForFunction(sut.Nodes(), function)
	is.OK(node)
	nodes := sut.To(node)

	is.OK(findMissingNode(nodes))
}

func TestContainerWithFuncProvidingRequiredFuncHasEdgeToFunc(t *testing.T) {
	is := is.New(t)
	function1 := makeFunc(nil, types.Typ[types.Int], false)
	function2 := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	sut.AddFunc(function1)
	sut.AddFunc(function2)
	node2 := findFuncNodeForFunction(sut.Nodes(), function2)
	is.OK(node2)
	nodes := sut.To(node2)

	is.OK(findFuncNodeForFunction(nodes, function1))
}

func TestContainerWithFuncProvidingRootHasEdgeToRoot(t *testing.T) {
	is := is.New(t)

	sut, typ := createRootedContainer()
	function := makeFunc(nil, typ, false)
	sut.AddFunc(function)
	root, _ := sut.Root()
	nodes := sut.To(root)

	is.OK(findFuncNodeForFunction(nodes, function))
}

func TestContainerWithUnsatifiedFuncRequirementHasEdgeFromMissingToFunc(t *testing.T) {
	is := is.New(t)
	function := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	sut.AddFunc(function)
	u := findMissingNode(sut.Nodes())
	v := findFuncNodeForFunction(sut.Nodes(), function)
	result := sut.HasEdgeFromTo(u, v)

	is.OK(u, v, result)
}

func TestContainerWithFuncProvidingRootHasEdgeFromFuncToRoot(t *testing.T) {
	is := is.New(t)

	sut, typ := createRootedContainer()
	function := makeFunc(nil, typ, false)
	sut.AddFunc(function)
	u := findFuncNodeForFunction(sut.Nodes(), function)
	v := findRootNode(sut.Nodes())
	result := sut.HasEdgeFromTo(u, v)

	is.OK(u, v, result)
}

func TestContainerWithFuncProvidingRequiredFuncHasEdgeFromFuncToFunc(t *testing.T) {
	is := is.New(t)
	function1 := makeFunc(nil, types.Typ[types.Int], false)
	function2 := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	sut := &Container{}
	sut.AddFunc(function1)
	sut.AddFunc(function2)
	u := findFuncNodeForFunction(sut.Nodes(), function1)
	v := findFuncNodeForFunction(sut.Nodes(), function2)
	result := sut.HasEdgeFromTo(u, v)

	is.OK(u, v, result)
}
