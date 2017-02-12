// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/token"
	"go/types"

	"github.com/gonum/graph"
)

func containsNode(nodes []graph.Node, expected graph.Node) bool {
	for _, node := range nodes {
		if node.ID() == expected.ID() {
			return true
		}
	}

	return false
}

func getRootNode(nodes []graph.Node) *rootNode {
	for _, node := range nodes {
		if n, ok := node.(*rootNode); ok {
			return n
		}
	}
	return nil
}

func containsMissingNode(nodes []graph.Node) bool {
	for _, node := range nodes {
		if _, ok := node.(*missingNode); ok {
			return true
		}
	}
	return false
}

func getNodeIDs(nodes []graph.Node) []int {
	var ids []int

	for _, node := range nodes {
		ids = append(ids, node.ID())
	}

	return ids
}

func createRootedContainer() (*Container, types.Type) {
	pkg := types.NewPackage("path", "mypackage")
	name := types.NewTypeName(token.NoPos, pkg, "MyIntType", nil)
	typ := types.NewNamed(name, types.Typ[types.Int], nil)

	container := &Container{}
	_ = container.SetRoot(typ)

	return container, typ
}

func findFuncNodeForFunction(nodes []graph.Node, function *types.Func) graph.Node {
	for _, node := range nodes {
		if fn, ok := node.(*funcNode); ok {
			if fn.function == function {
				return fn
			}
		}
	}

	return nil
}

func findMissingNode(nodes []graph.Node) graph.Node {
	for _, node := range nodes {
		if missing, ok := node.(*missingNode); ok {
			return missing
		}
	}
	return nil
}

func findRootNode(nodes []graph.Node) graph.Node {
	for _, node := range nodes {
		if root, ok := node.(*rootNode); ok {
			return root
		}
	}
	return nil
}

func hasFuncNodeForFunction(nodes []graph.Node, function *types.Func) func() bool {
	return func() bool {
		for _, node := range nodes {
			if fn, ok := node.(*funcNode); ok {
				if fn.function == function {
					return true
				}
			}
		}

		return false
	}
}

func containsType(list []types.Type, expectedItem types.Type) bool {
	for _, item := range list {
		if item == expectedItem {
			return true
		}
	}
	return false
}

func makeFunc(param, ret types.Type, returnsErr bool) *types.Func {
	var paramTuple *types.Tuple
	if param == nil {
		paramTuple = types.NewTuple()
	} else {
		paramTuple = types.NewTuple(types.NewVar(token.NoPos, nil, "", param))
	}

	resultVar := types.NewVar(token.NoPos, nil, "", ret)
	var retTuple *types.Tuple
	if returnsErr {
		errObj := types.Universe.Lookup("error")
		errVar := types.NewVar(token.NoPos, nil, "", errObj.Type())
		retTuple = types.NewTuple(resultVar, errVar)
	} else {
		retTuple = types.NewTuple(resultVar)
	}

	sig := types.NewSignature(nil, paramTuple, retTuple, false)
	return types.NewFunc(token.NoPos, nil, "myfunc", sig)
}
