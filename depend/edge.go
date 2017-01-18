// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import "github.com/gonum/graph"

type edgeImpl struct {
	from, to graph.Node
}

func (e *edgeImpl) From() graph.Node {
	return e.from
}

func (e *edgeImpl) To() graph.Node {
	return e.to
}

func (e *edgeImpl) Weight() float64 {
	return 1.0
}
