// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import "go/types"

// A rootNode generates a code fragment to return the instance of its one
// required type from the builder function. This type is the anchor of the
// Container and it is expected that all other (useful) nodes will be a part
// of the transitive closure of the requirement of this node. There should be
// at most one rootNode in a given Container.
type rootNode struct {
	container *Container
	id        int
	root      types.Type
}

func newRootNode(container *Container, id int, root types.Type) *rootNode {
	return &rootNode{
		container: container,
		id:        id,
		root:      root,
	}
}

func (r rootNode) ID() int {
	if r.id < 0 {
		panic("Root node cannot have a negative id.")
	}
	return r.id
}

func (r rootNode) Generate() {
	panic("not implemented")
}

func (r rootNode) requires() []types.Type {
	return []types.Type{r.root}
}

func (r rootNode) provides() []types.Type {
	return nil
}

func (r rootNode) getContainer() *Container {
	return r.container
}

// A missingNode is a placeholder for another type of node that has not yet
// been added to the Container. It allows the requirements of a node to be
// expressed in the Container before the provider that can satisfy those
// requirements has been added. A missingNode does not itself have any requirements.
// An attempt to generate a code fragment for a missingNode is an error which
// indicates that a requirement for some other node has not been met. There
// should be exactly one missingNode in a given Container. A Container in which
// all nodes requirements have been met will not have any edges from its
// missingNode to any other nodes.
type missingNode struct {
	container *Container
	id        int
}

func newMissingNode(container *Container, id int) *missingNode {
	return &missingNode{container: container, id: id}
}

func (m missingNode) ID() int {
	return m.id
}

func (m missingNode) Generate() {
	panic("not implemented")
}

func (m missingNode) requires() []types.Type {
	return nil
}

func (m missingNode) provides() []types.Type {
	if m.container.requiredBy == nil || m.container.providedBy == nil {
		return nil
	}

	var result []types.Type

	for _, typ := range m.container.requiredBy.Types() {
		if len(m.container.providedBy.Nodes(typ)) == 0 {
			result = append(result, typ)
		}
	}

	return result
}

func (m missingNode) getContainer() *Container {
	return m.container
}

// A funcNode generates a code fragment to produce instances of the provided
// types by calling a function (a constructor or other static factory). Its
// required types are the parameters to the function and its provided types
// are the (non-error) results of the function.
type funcNode struct {
	container *Container
	id        int
	function  *types.Func
}

func newFuncNode(container *Container, id int, function *types.Func) (*funcNode, error) {
	sig := function.Type().(*types.Signature)

	// Check for a method.
	if sig.Recv() != nil {
		return nil, newInvalidFuncError(function, "cannot add methods to a Container")
	}

	// Check for an error return type that is not the last return type
	if tupleHasEarlyError(sig.Results()) {
		return nil, newInvalidFuncError(function, "error return type must be last return type")
	}

	node := &funcNode{
		container: container,
		id:        id,
		function:  function,
	}
	return node, nil
}

func (f funcNode) ID() int {
	if f.id < 0 {
		panic("Non singleton nodes cannot have a negative id.")
	}
	return f.id
}

func (f funcNode) Generate() {
	panic("not implemented")
}

func (f funcNode) requires() []types.Type {
	sig := f.function.Type().(*types.Signature)

	return extractTypesForTuple(sig.Params(), false)
}

func (f funcNode) provides() []types.Type {
	sig := f.function.Type().(*types.Signature)

	return extractTypesForTuple(sig.Results(), true)
}

func (f funcNode) getContainer() *Container {
	return f.container
}

func extractTypesForTuple(tuple *types.Tuple, excludeError bool) []types.Type {
	var result []types.Type
	errType := types.Universe.Lookup("error").Type()

	for i := 0; i < tuple.Len(); i++ {
		typ := tuple.At(i).Type()
		if !excludeError || !types.Identical(typ, errType) {
			result = append(result, tuple.At(i).Type())
		}
	}

	return result
}

func tupleHasEarlyError(tuple *types.Tuple) bool {
	errType := types.Universe.Lookup("error").Type()

	// Note: "tuple.Len() - 1" is correct because an error
	// as the last return type should return false
	for i := 0; i < tuple.Len()-1; i++ {
		if types.Identical(tuple.At(i).Type(), errType) {
			return true
		}
	}

	return false
}

var _ commonNode = missingNode{}
var _ commonNode = rootNode{}
var _ commonNode = funcNode{}
