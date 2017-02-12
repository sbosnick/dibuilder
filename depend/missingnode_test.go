// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"testing"

	"github.com/cheekybits/is"
	"github.com/stretchr/testify/assert"
)

func TestMissingNodeRequiresNothing(t *testing.T) {
	container := &Container{}

	sut := missingNode{container: container}

	assert.Len(t, sut.requires(), 0, "missingNode unexpectedly requires some types")
}

func TestMissingNodeProvidesNothing(t *testing.T) {
	container := &Container{}

	sut := missingNode{container: container}

	assert.Len(t, sut.provides(), 0, "missingNode unexpectedly provides some types")
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
