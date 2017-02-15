// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"testing"

	"github.com/cheekybits/is"
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

func TestContainerAddFuncDetectsAndAddsRootTypeProvidedByFunc(t *testing.T) {
	is := is.New(t)
	roottype := makeRunnableType("MyIntType")
	function := makeFunc(nil, roottype, false)

	sut := &Container{}
	sut.AddFunc(function)

	root, err := sut.Root()
	is.NoErr(err)
	is.OK(root)
}

func TestContainerAddFuncIsErrorWithSecondRootProvider(t *testing.T) {
	is := is.New(t)
	roottype1 := makeRunnableType("MyFirstType")
	roottype2 := makeRunnableType("MySecondType")
	function1 := makeFunc(nil, roottype1, false)
	function2 := makeFunc(nil, roottype2, false)

	sut := &Container{}
	sut.AddFunc(function1)
	err := sut.AddFunc(function2)

	is.Err(err)
}
