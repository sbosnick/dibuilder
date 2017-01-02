// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import "github.com/stretchr/testify/assert"
import "testing"

func TestRootNodeIDIsNegative(t *testing.T) {
	sut := rootNode{}
	id := sut.ID()

	assert.Condition(t, func() bool { return id < 0 }, "Non-negative ID()")
}

func TestMissingNodeIDIsNegative(t *testing.T) {
	sut := missingNode{}
	id := sut.ID()

	assert.Condition(t, func() bool { return id < 0 }, "Non-negative ID()")
}

func TestSingletonNodesIDAreDifferent(t *testing.T) {
	root := rootNode{}
	missing := missingNode{}
	id1 := root.ID()
	id2 := missing.ID()

	assert.NotEqual(t, id1, id2, "rootNode.ID() == missingNode.ID()")
}

func TestFuncNodeWithNonNegativeIDReturnsExpectedID(t *testing.T) {
	expected := 1

	sut := funcNode{id: expected}
	id := sut.ID()

	assert.Equal(t, expected, id)
}

func TestFuncNodeWithNegativeIDPanicsOnID(t *testing.T) {
	sut := funcNode{id: -1}

	assert.Panics(t, func() { sut.ID() }, "Negative ID did not panic")
}
