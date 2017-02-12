// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import "go/types"

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

var _ commonNode = funcNode{}
