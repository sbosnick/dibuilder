// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

// Package depend provides the core data types for the buildtime dependency injection system.
package depend

import (
	"go/types"

	"github.com/gonum/graph"
)

// A Container exposes the dependencies implicit in set of constructors or
// static factories as a directed graph.Container implements the
// graph.Directed interface from github.com/gonum/graph.
type Container struct {
	root types.Type
}

// Has returns whether a node exists within the Container.
func (c *Container) Has(node graph.Node) bool {
	switch n := node.(type) {
	case missingNode:
		return n.container == c
	case rootNode:
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
func (c *Container) From(node graph.Node) []graph.Node {
	var nodes []graph.Node

	switch node.(type) {
	case missingNode:
		if c.hasRoot() {
			nodes = append(nodes, rootNode{container: c})
		}
	}

	return nodes
}

// To returns all nodes that can reach directly to the given node.
func (c *Container) To(node graph.Node) []graph.Node {
	var nodes []graph.Node

	switch node.(type) {
	case rootNode:
		if c.hasRoot() {
			nodes = append(nodes, missingNode{container: c})
		}
	}

	return nodes
}

// HasEdgeBetween returns whether an edge exists between nodes x and y without considering
// the direction.
func (c *Container) HasEdgeBetween(x graph.Node, y graph.Node) bool {
	return c.HasEdgeFromTo(x, y) || c.HasEdgeFromTo(y, x)
}

// HasEdgeFromTo returns whether an edge exists in the Container from u to v.
func (c *Container) HasEdgeFromTo(u graph.Node, v graph.Node) bool {
	switch u.(type) {
	case missingNode:
		if _, ok := v.(rootNode); ok && c.hasRoot() {
			return true
		}
	}

	return false
}

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (c *Container) Edge(u graph.Node, v graph.Node) graph.Edge {
	if c.HasEdgeFromTo(u, v) {
		return &edgeImpl{from: u, to: v}
	}

	return nil
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

// An edge expresses the dependency relationship between two nodes in a Container.
// The relationship is from a node that provides a type to a node that requires
// that type. The name of the edge is common to all edges for a given type that
// originate from the same node. The name is designed to be used as a variable
// name in a generated builder function.
type edge interface {
	graph.Edge
	name() string
	edgeType() types.Type
}

// A node is an element in a Container that can generate a code fragment to
// produce named instances of specific types but requires that named instances
// of other types be produced first. The generated code uses they required, named
// instances of the other types to provide the named instances of the specific types.
type node interface {
	graph.Node
	Generate()
	requires() []edge
	provides() []edge
}
