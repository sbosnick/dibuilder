// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

// Package depend provides the core data types for the buildtime dependency injection system.
package depend

import (
	"go/types"

	"golang.org/x/tools/go/types/typeutil"

	"github.com/gonum/graph"
)

// A Container exposes the dependencies implicit in set of constructors or
// static factories as a directed graph.Container implements the
// graph.Directed interface from github.com/gonum/graph.
type Container struct {
	rootnode   *rootNode
	nodes      []commonNode
	providedBy *typeNodeMap
	requiredBy *typeNodeMap
}

// Has returns whether a node exists within the Container.
func (c *Container) Has(node graph.Node) bool {
	switch n := node.(type) {
	case missingNode:
		return n.container == c
	case commonNode:
		return n.getContainer() == c && n.ID() < c.nextID()
	default:
		return false
	}
}

// Nodes returns all of the nodes within the Container.
func (c *Container) Nodes() []graph.Node {
	var nodes []graph.Node

	nodes = append(nodes, missingNode{container: c})

	for _, node := range c.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// From returns all nodes that can be reached directly from the given node.
func (c *Container) From(node graph.Node) []graph.Node {
	var nodes []graph.Node

	switch node := node.(type) {
	case missingNode:
		for _, typ := range c.requiredBy.Types() {
			if len(c.providedBy.Nodes(typ)) == 0 {
				for _, reqnode := range c.requiredBy.Nodes(typ) {
					nodes = append(nodes, reqnode)
				}
			}
		}

	case commonNode:
		for _, provide := range node.provides() {
			for _, requirer := range c.requiredBy.Nodes(provide) {
				nodes = append(nodes, requirer)
			}
		}
	}

	return nodes
}

// To returns all nodes that can reach directly to the given node.
func (c *Container) To(node graph.Node) []graph.Node {
	var nodes []graph.Node

	switch node.(type) {
	case *rootNode:
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
		if _, ok := v.(*rootNode); ok && c.hasRoot() {
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
func (c *Container) SetRoot(root types.Type) error {
	if c.hasRoot() {
		return ErrRootAlreadySet
	}

	c.rootnode = newRootNode(c, c.nextID(), root)
	c.addNode(c.rootnode)
	return nil
}

// Root returns the root node of the container or ErrNoRoot is a root
// has not been set.
func (c *Container) Root() (graph.Node, error) {
	if c.hasRoot() {
		return c.rootnode, nil
	}

	return nil, ErrNoRoot

}

// AddFunc adds function to the Container. function should be a constructor
// or other static factory. The non-error return types of function are made
// available as components that can satisfy the components required by other
// functions added to the container or by SetRoot. The parameters to the function
// are required to be satisfied by components in the Container for the Container
// to be complete. function can have an error return type as its last return type.
// AddFunc will return an InvalidFuncError for a function with an error return type
// in any position except the last. It will also return an InvalidFuncError if a
// method is passed in as function.
func (c *Container) AddFunc(function *types.Func) error {
	// create a new node
	node, err := newFuncNode(c, c.nextID(), function)
	if err != nil {
		return err
	}

	c.addNode(node)

	return nil
}

func (c *Container) hasRoot() bool {
	return c.rootnode != nil
}

func (c *Container) ensureMaps() {
	if c.requiredBy != nil && c.providedBy != nil {
		return
	}

	hasher := typeutil.MakeHasher()
	c.requiredBy = newTypeNodeMap(hasher)
	c.providedBy = newTypeNodeMap(hasher)
}

func (c *Container) nextID() int {
	return len(c.nodes)
}

func (c *Container) addNode(newNode commonNode) {
	// add the node to the nodes slice
	c.nodes = append(c.nodes, newNode)

	// add the node to the appropriate maps for its provides() and requires() types
	c.ensureMaps()
	for _, typ := range newNode.provides() {
		c.providedBy.AddNode(typ, newNode)
	}
	for _, typ := range newNode.requires() {
		c.requiredBy.AddNode(typ, newNode)
	}
}

// A node is an element in a Container that can generate a code fragment to
// produce instances of specific types but requires that instances  of other types
// be produced first. The generated code uses the required instances of
// the other types to provide the instances of the specific types.
type commonNode interface {
	graph.Node
	Generate()
	requires() []types.Type
	provides() []types.Type
	getContainer() *Container
}
