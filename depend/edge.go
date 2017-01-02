// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/types"

	"github.com/gonum/graph"
)

type edgeImpl struct{}

func newEdge(from node, to node, name string, typ types.Type) edgeImpl {
	panic("not implemented")
}

func (e edgeImpl) From() graph.Node {
	panic("not implemented")
}

func (e edgeImpl) To() graph.Node {
	panic("not implemented")
}

func (e edgeImpl) Weight() float64 {
	panic("not implemented")
}

func (e edgeImpl) name() string {
	panic("not implemented")
}

func (e edgeImpl) edgeType() types.Type {
	panic("not implemented")
}
