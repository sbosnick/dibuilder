// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import "go/types"

// The fixed ID's used for the single-node-per-container notes. These must
// be negative.
const (
	rootNodeID    int = -2
	missingNodeID int = -1
)

// A rootNode generates a code fragment to return the instance of its one
// required type from the builder function. This type is the anchor of the
// Container and it is expected that all other (useful) nodes will be a part
// of the transitive closure of the requirement of this node. There should be
// at most one rootNode in a given Container.
type rootNode struct {
	container *Container
}

func (r rootNode) ID() int {
	return rootNodeID
}

func (r rootNode) Generate() {
	panic("not implemented")
}

func (r rootNode) requires() []types.Type {
	return []types.Type{r.container.root}
}

func (r rootNode) provides() []types.Type {
	return nil
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
}

func (m missingNode) ID() int {
	return missingNodeID
}

func (m missingNode) Generate() {
	panic("not implemented")
}

func (m missingNode) requires() []types.Type {
	return nil
}

func (m missingNode) provides() []types.Type {
	return nil
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

var _ node = missingNode{}
var _ node = rootNode{}
var _ node = funcNode{}
