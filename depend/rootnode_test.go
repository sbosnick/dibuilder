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

func TestRootNodeWithNonNegativeIDReturnsExpectedID(t *testing.T) {
	is := is.New(t)
	expected := 1

	sut := newRootNode(nil, expected, nil)

	is.Equal(sut.ID(), expected)
}

func TestRootNodeWithNegativeIDPanicsOnID(t *testing.T) {
	is := is.New(t)

	sut := newRootNode(nil, -2, nil)

	is.Panic(func() { sut.ID() })
}

func TestRootNodeProvidesNothing(t *testing.T) {
	sut, _ := createRootedContainer()
	root, _ := sut.Root()
	rootnode := root.(*rootNode)

	assert.Len(t, rootnode.provides(), 0, "Root node unexpectedly provides some types")
}

func TestRootNodeRequiresTypeSetOnContainer(t *testing.T) {
	is := is.New(t)
	expected := types.Typ[types.Int]
	container := &Container{}
	_ = container.SetRoot(expected)

	sut, _ := container.Root()
	sutnode := sut.(*rootNode)
	requires := sutnode.requires()

	is.Equal(len(requires), 1)
	is.OK(containsType(requires, expected))
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
