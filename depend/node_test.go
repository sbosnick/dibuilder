// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/cheekybits/is"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestRootNodeProvidesNothing(t *testing.T) {
	sut, _ := createRootedContainer()
	root, _ := sut.Root()
	rootnode := root.(rootNode)

	assert.Len(t, rootnode.provides(), 0, "Root node unexpectedly provides some types")
}

func TestRootNodeRequiresTypeSetOnContainer(t *testing.T) {
	is := is.New(t)
	expected := types.Typ[types.Int]
	container := &Container{}
	container.SetRoot(expected)

	sut, _ := container.Root()
	sutnode := sut.(rootNode)
	requires := sutnode.requires()

	is.Equal(len(requires), 1)
	is.OK(containsType(requires, expected))
}

func containsType(list []types.Type, expectedItem types.Type) bool {
	for _, item := range list {
		if item == expectedItem {
			return true
		}
	}
	return false
}

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

func makeFunc(param, ret types.Type, returnsErr bool) *types.Func {
	var paramTuple *types.Tuple
	if param == nil {
		paramTuple = types.NewTuple()
	} else {
		paramTuple = types.NewTuple(types.NewVar(token.NoPos, nil, "", param))
	}

	resultVar := types.NewVar(token.NoPos, nil, "", ret)
	var retTuple *types.Tuple
	if returnsErr {
		errObj := types.Universe.Lookup("error")
		errVar := types.NewVar(token.NoPos, nil, "", errObj.Type())
		retTuple = types.NewTuple(resultVar, errVar)
	} else {
		retTuple = types.NewTuple(resultVar)
	}

	sig := types.NewSignature(nil, paramTuple, retTuple, false)
	return types.NewFunc(token.NoPos, nil, "myfunc", sig)
}

func TestFuncNodeRequiresTypesOfFuncParameters(t *testing.T) {
	param := types.Typ[types.Int]
	ret := types.Typ[types.Bool]
	function := makeFunc(param, ret, false)

	sut := funcNode{function: function}
	requires := sut.requires()

	assert.Len(t, requires, 1, "Unexpected number of required types on funcNode")
	assert.Contains(t, requires, param, "funcNode did not require the expected type")
}

func TestFuncNodeProvidesTypesOfFuncReturns(t *testing.T) {
	param := types.Typ[types.Int]
	ret := types.Typ[types.Bool]
	function := makeFunc(param, ret, false)

	sut := funcNode{function: function}
	provides := sut.provides()

	assert.Len(t, provides, 1, "Unexpected number of provided types on funcNode")
	assert.Contains(t, provides, ret, "funcNode did not provide the expected type")
}

func TestFuncNodeDoesNotProvideErrorType(t *testing.T) {
	param := types.Typ[types.Int]
	ret := types.Typ[types.Bool]
	function := makeFunc(param, ret, true)

	sut := funcNode{function: function}
	provides := sut.provides()

	assert.Len(t, provides, 1, "Unexpected number of provided types on funcNode")
	assert.NotContains(t, provides, types.Universe.Lookup("error").Type(),
		"Unexpected type provided by funcNode")
}

func TestNewFuncNodeWithMethodIsError(t *testing.T) {
	pkg := types.NewPackage("github.com/sbosnick/mytestpkg", "mytestpkg")
	name := types.NewTypeName(token.NoPos, pkg, "MyType", nil)
	typ := types.NewNamed(name, types.Typ[types.Int], nil)
	receiver := types.NewParam(token.NoPos, pkg, "m", typ)
	sig := types.NewSignature(receiver, types.NewTuple(), types.NewTuple(), false)
	function := types.NewFunc(token.NoPos, pkg, "MyFunc", sig)

	_, err := newFuncNode(nil, 0, function)

	require.Error(t, err, "Expected error was not returned")
	assert.IsType(t, &InvalidFuncError{}, err, "Error return not of expected type")
}

func TestNewFuncNodeWithEarlyErrorReturnIsError(t *testing.T) {
	errtyp := types.Universe.Lookup("error").Type()
	inttyp := types.Typ[types.Int]
	sig := types.NewSignature(nil, types.NewTuple(),
		types.NewTuple(
			types.NewParam(token.NoPos, nil, "", errtyp),
			types.NewParam(token.NoPos, nil, "", inttyp)),
		false)
	function := types.NewFunc(token.NoPos, nil, "MyFunc", sig)

	_, err := newFuncNode(nil, 0, function)

	require.Error(t, err, "Expected error was not returned")
	assert.IsType(t, &InvalidFuncError{}, err, "Error return not of expected type")
}

func TestNewFuncNodeWithNoErrorReturnGivesNode(t *testing.T) {
	function := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], false)

	node, err := newFuncNode(nil, 0, function)

	require.NoError(t, err, "Unexpected error in newFuncNode call")
	require.NotNil(t, node, "Returned funcNode was nil in newFuncNode call")
	assert.Equal(t, function, node.function)
}

func TestNewFuncNodeWtihLastErrorReturnGivesNode(t *testing.T) {
	function := makeFunc(types.Typ[types.Int], types.Typ[types.Bool], true)

	node, err := newFuncNode(nil, 0, function)

	require.NoError(t, err, "Unexpected error in newFuncNode call")
	require.NotNil(t, node, "Returned funcNode was nil in newFuncNode call")
	assert.Equal(t, function, node.function)
}
