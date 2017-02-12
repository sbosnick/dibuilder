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

var _ commonNode = rootNode{}
