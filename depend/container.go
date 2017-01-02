// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

// Package depend provides the core data types for the buildtime dependancy injection system.
package depend

import (
	"go/types"

	"github.com/gonum/graph"
)

// A Container exposes the dependancies implicit in set of constructors or
// static factories as a directed graph.Container implements the
// graph.Directed interface from github.com/gonum/graph.
type Container struct {
	dummy int // temporarily needed so that two different Container's have different addresses
}

func (c *Container) Has(node graph.Node) bool {
	switch n := node.(type) {
	case *missingNode:
		return n.container == c
	default:
		return false
	}
}

func (c *Container) Nodes() []graph.Node {
	panic("not implemented")
}

func (c *Container) From(graph.Node) []graph.Node {
	panic("not implemented")
}

func (c *Container) To(graph.Node) []graph.Node {
	panic("not implemented")
}

func (c *Container) HasEdgeBetween(x graph.Node, y graph.Node) bool {
	panic("not implemented")
}

func (c *Container) HasEdgeFromTo(u graph.Node, v graph.Node) bool {
	panic("not implemented")
}

func (c *Container) Edge(u graph.Node, v graph.Node) graph.Edge {
	panic("not implemented")
}

func (c *Container) SetRoot(root types.Type) {
	panic("Not implemented")
}

func (c *Container) Root() (graph.Node, error) {
	panic("Not implemented")
}

func (c *Container) AddFunc(f types.Func) {
	panic("Not implemented")
}

// An edge expreses the dependancy relationship between two nodes in a Container.
// The relationship is from a node that provides a type to a node that requires
// that type. The name of the edge is common to all edges for a given type that
// orginate from the same node. The name is designed to be used as a variable
// name in a generated builder function.
type edge interface {
	graph.Edge
	name() string
	edgeType() types.Type
}

// A node is an element in a Container that can generate a code fragment to
// produce named instances of specific types but requires that named instances
// of other types be produced first. The generated code uses they required, named
// instances of the other types to provide the named instatances of the specific types.
type node interface {
	graph.Node
	Generate()
	requires() []edge
	provides() []edge
}
