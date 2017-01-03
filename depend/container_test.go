// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

// Package depend provides the core data types for the buildtime dependancy injection system.
package depend

import (
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
