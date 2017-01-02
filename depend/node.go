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
	container Container
	root      edge
}

func (r rootNode) ID() int {
	panic("not implemented")
}

func (r rootNode) Generate() {
	panic("not implemented")
}

func (r rootNode) requires() []edge {
	panic("not implemented")
}

func (r rootNode) provides() []edge {
	panic("not implemented")
}

// A missingNode is a placeholder for another type of node that has not yet
// been added to the Container. It allows the requirements of a node to be
// expressed in the Container before the provider that can satify those
// requirments has been added. A missingNode does not itself have any requirments.
// An attempt to generate a code fragment for a missingNode is an error which
// indicates that a requirement for some other node has not been met. There
// should be exactly one missingNode in a given Container. A Container in which
// all nodes requirements have been met will not have any edges from its
// missingNode to any other nodes.
type missingNode struct {
	container Container
}

func (m missingNode) ID() int {
	panic("not implemented")
}

func (m missingNode) Generate() {
	panic("not implemented")
}

func (m missingNode) requires() []edge {
	panic("not implemented")
}

func (m missingNode) provides() []edge {
	panic("not implemented")
}

// A funcNode generates a code fragment to produce instances of the provided
// types by calling a function (a constructor or other static factory). Its
// required types are the parameters to the function and its provided types
// are the (non-error) results of the function.
type funcNode struct {
	container Container
	function  types.Func
}

func (f funcNode) ID() int {
	panic("not implemented")
}

func (f funcNode) Generate() {
	panic("not implemented")
}

func (f funcNode) requires() []edge {
	panic("not implemented")
}

func (f funcNode) provides() []edge {
	panic("not implemented")
}
