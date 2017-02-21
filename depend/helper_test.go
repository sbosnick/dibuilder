// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/token"
	"go/types"

	"github.com/gonum/graph"
)

func createRootedContainer() (*Container, types.Type) {
	pkg := types.NewPackage("path", "mypackage")
	name := types.NewTypeName(token.NoPos, pkg, "MyIntType", nil)
	typ := types.NewNamed(name, types.Typ[types.Int], nil)

	container := &Container{}
	_ = container.setRoot(typ)

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

func containsType(list []types.Type, expectedItem types.Type) bool {
	for _, item := range list {
		if item == expectedItem {
			return true
		}
	}
	return false
}

func makeSignature(param, ret types.Type, returnsErr bool) *types.Signature {
	var paramTuple *types.Tuple
	if param == nil {
		paramTuple = types.NewTuple()
	} else {
		paramTuple = types.NewTuple(types.NewVar(token.NoPos, nil, "", param))
	}

	var resultsVars []*types.Var
	if ret != nil {
		resultsVars = append(resultsVars, types.NewVar(token.NoPos, nil, "", ret))
	}
	if returnsErr {
		errObj := types.Universe.Lookup("error")
		resultsVars = append(resultsVars, types.NewVar(token.NoPos, nil, "", errObj.Type()))
	}
	retTuple := types.NewTuple(resultsVars...)

	return types.NewSignature(nil, paramTuple, retTuple, false)
}

func makeFunc(param, ret types.Type, returnsErr bool) *types.Func {
	sig := makeSignature(param, ret, returnsErr)
	return types.NewFunc(token.NoPos, nil, "myfunc", sig)
}

func makeRunnableType(name string) types.Type {
	named := makeNamedType(name, types.Typ[types.Int])
	sig := types.NewSignature(
		types.NewParam(token.NoPos, nil, "m", named),
		types.NewTuple(),
		types.NewTuple(),
		false)
	function := types.NewFunc(token.NoPos, nil, "Run", sig)
	named.AddMethod(function)

	return named
}

func makeNamedType(name string, underlying types.Type) *types.Named {
	typename := types.NewTypeName(token.NoPos, nil, name, nil)
	return types.NewNamed(typename, underlying, nil)
}
