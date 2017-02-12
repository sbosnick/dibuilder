// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import "go/types"

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

var _ commonNode = missingNode{}
