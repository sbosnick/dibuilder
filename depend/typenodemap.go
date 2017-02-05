// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/types"

	"golang.org/x/tools/go/types/typeutil"
)

type typeNodeMap struct {
	typeMap typeutil.Map
}

func newTypeNodeMap(hasher typeutil.Hasher) *typeNodeMap {
	tnm := typeNodeMap{}
	tnm.typeMap.SetHasher(hasher)
	return &tnm
}

func (m *typeNodeMap) AddNode(typ types.Type, n commonNode) {
	if m == nil {
		panic("AddNode: attempt to add to a nil typeNodeMap")
	}

	result := m.typeMap.At(typ)
	var nodes []commonNode
	if result != nil {
		nodes = result.([]commonNode)
	}

	nodes = append(nodes, n)
	m.typeMap.Set(typ, nodes)
}

func (m *typeNodeMap) Nodes(typ types.Type) []commonNode {
	if m == nil {
		var ret []commonNode
		return ret
	}

	result := m.typeMap.At(typ)
	var nodes []commonNode
	if result != nil {
		nodes = result.([]commonNode)
	}

	return nodes
}

func (m *typeNodeMap) Types() []types.Type {
	if m == nil {
		var ret []types.Type
		return ret
	}

	return m.typeMap.Keys()
}
