// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
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
	result := sut.Has(&missing)

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

func TestRootedContainerRootIsRootNode(t *testing.T) {
	pkg := types.NewPackage("path", "mypackage")
	name := types.NewTypeName(token.NoPos, pkg, "MyIntType", types.Typ[types.Int])

	sut := &Container{}
	sut.SetRoot(name.Type())
	root, err := sut.Root()

	assert.NoError(t, err, "unexpected error in getting the container root")
	assert.IsType(t, rootNode{}, root, "unexpected node type for the root node")
}

func TestRootedContainerHasRootNode(t *testing.T) {
	pkg := types.NewPackage("path", "mypackage")
	name := types.NewTypeName(token.NoPos, pkg, "MyIntType", types.Typ[types.Int])

	sut := &Container{}
	sut.SetRoot(name.Type())
	root := rootNode{container: sut}
	result := sut.Has(&root)

	assert.True(t, result, "Container did not have a rootNode")
}

func TestRootedContainerDoesNotHasRootNodeForOtherContainer(t *testing.T) {
	pkg := types.NewPackage("path", "mypackage")
	name := types.NewTypeName(token.NoPos, pkg, "MyIntType", types.Typ[types.Int])
	other := &Container{}
	root := rootNode{container: other}

	sut := &Container{}
	sut.SetRoot(name.Type())
	result := sut.Has(&root)

	assert.False(t, result, "Container did not had rootNode for other Container")
}

func TestRootedContainerReturnsRootNode(t *testing.T) {
	pkg := types.NewPackage("path", "mypackage")
	name := types.NewTypeName(token.NoPos, pkg, "MyIntType", types.Typ[types.Int])

	sut := &Container{}
	sut.SetRoot(name.Type())
	nodes := sut.Nodes()

	assert.Contains(t, nodes, rootNode{container: sut}, "Root node is missing.")
}
