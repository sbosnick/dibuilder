// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/types"
	"testing"

	"golang.org/x/tools/go/types/typeutil"

	"github.com/stretchr/testify/assert"
)

func TestNilTypeNodeMapGivesEmptyNodes(t *testing.T) {
	var sut *typeNodeMap
	nodes := sut.Nodes(types.Typ[types.Int])

	assert.Len(t, nodes, 0, "Nodes was not empty")
}

func TestNilTypeNodeMapPanicsOnAddNode(t *testing.T) {
	var sut *typeNodeMap
	test := func() {
		sut.AddNode(types.Typ[types.Int], &funcNode{})
	}

	assert.Panics(t, test, "Expected AddNode to panic")
}

func TestTypeNodeMapNodesIncludesAddedNode(t *testing.T) {
	node := &funcNode{}
	typ := types.Typ[types.Int]

	sut := newTypeNodeMap(typeutil.MakeHasher())
	sut.AddNode(typ, node)
	nodes := sut.Nodes(typ)

	assert.Contains(t, nodes, node, "Expected node not present in Nodes()")
}

func TestTypeNodeMapNodesIncludesTwoAddedNode(t *testing.T) {
	node1 := &funcNode{}
	node2 := &funcNode{}
	typ := types.Typ[types.Int]

	sut := newTypeNodeMap(typeutil.MakeHasher())
	sut.AddNode(typ, node1)
	sut.AddNode(typ, node2)
	nodes := sut.Nodes(typ)

	assert.Contains(t, nodes, node1, "Expected node1 not present in Nodes()")
	assert.Contains(t, nodes, node2, "Expected node2 not present in Nodes()")
}

func TestTypeNodeMapNodesIncludesTwoAddedTypes(t *testing.T) {
	node := &funcNode{}
	typ1 := types.Typ[types.Int]
	typ2 := types.Typ[types.Bool]

	sut := newTypeNodeMap(typeutil.MakeHasher())
	sut.AddNode(typ1, node)
	sut.AddNode(typ2, node)
	nodes1 := sut.Nodes(typ1)
	nodes2 := sut.Nodes(typ2)

	assert.Contains(t, nodes1, node, "Expected node not present in Nodes()")
	assert.Contains(t, nodes2, node, "Expected node not present in Nodes()")
}

func TestTypeNodeMapGivesEmptyNodesForAbsentType(t *testing.T) {
	sut := newTypeNodeMap(typeutil.MakeHasher())
	nodes := sut.Nodes(types.Typ[types.Int])

	assert.Len(t, nodes, 0, "Node unexpectedly not 0 length")
}
