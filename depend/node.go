// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/types"

	"github.com/gonum/graph"
)

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
