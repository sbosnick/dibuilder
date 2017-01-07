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
	root types.Type
}

// Has returns whether a node exists within the Container.
func (c *Container) Has(node graph.Node) bool {
	switch n := node.(type) {
	case *missingNode:
		return n.container == c
	case *rootNode:
		return c.hasRoot() && n.container == c
	default:
		return false
	}
}

// Nodes returns all of the nodes within the Container.
func (c *Container) Nodes() []graph.Node {
	var nodes []graph.Node

	nodes = append(nodes, missingNode{container: c})

	if c.hasRoot() {
		nodes = append(nodes, rootNode{container: c})
	}

	return nodes
}

// From returns all nodes that can be reached directly from the given node.
func (c *Container) From(graph.Node) []graph.Node {
	var nodes []graph.Node

	return nodes
}

// To returns all nodes that can reach directly to the given node.
func (c *Container) To(graph.Node) []graph.Node {
	var nodes []graph.Node

	return nodes
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

// SetRoot sets the root type for the Container. A Container for which a root
// type has been set has a root node.
func (c *Container) SetRoot(root types.Type) {
	c.root = root
}

// Root returns the root node of the container or ErrNoRoot is a root
// has not been set.
func (c *Container) Root() (graph.Node, error) {
	if c.hasRoot() {
		return rootNode{container: c}, nil
	}

	return nil, ErrNoRoot

}

func (c *Container) AddFunc(f types.Func) {
	panic("Not implemented")
}

func (c *Container) hasRoot() bool {
	return c.root != nil
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
