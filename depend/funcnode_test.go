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
