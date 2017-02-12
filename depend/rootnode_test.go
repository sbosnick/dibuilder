// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/types"
	"testing"

	"github.com/cheekybits/is"
	"github.com/stretchr/testify/assert"
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
