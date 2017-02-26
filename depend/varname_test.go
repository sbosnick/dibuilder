// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/token"
	"go/types"
	"testing"

	"golang.org/x/tools/go/types/typeutil"

	"github.com/cheekybits/is"
)

func TestGetVarPrefix(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		expected string
		typ      types.Type
	}{
		{"", makeNamedType("MyName", types.Typ[types.Int])},
		{"b", types.Typ[types.Int]},
		{"", types.NewPointer(types.Typ[types.Int])},
		{"", types.NewArray(types.Typ[types.Int], 3)},
		{"", types.NewSlice(types.Typ[types.Int])},
		{"m", types.NewMap(types.Typ[types.Int], types.Typ[types.Float64])},
		{"", types.NewChan(types.SendRecv, types.Typ[types.Int])},
		{"s", types.NewStruct(nil, nil)},
		{"sig", types.NewSignature(nil, nil, nil, false)},
		{"int", types.NewInterface(nil, nil)},
	}

	for _, test := range tests {
		result := getVarPrefix(test.typ)

		is.Equal(result, test.expected)
	}
}

func TestGetBasename(t *testing.T) {
	is := is.New(t)
	basic := types.Typ[types.Int]
	basicfield := types.NewField(token.NoPos, nil, "doit", basic, false)
	named := makeNamedType("MyName", basic)
	namedfield := types.NewField(token.NoPos, nil, "doit", named, false)

	tests := []struct {
		expected string
		typ      types.Type
	}{
		{"myName", named},
		{"var0", basic},
		{"myName", types.NewPointer(named)},
		{"var0", types.NewPointer(basic)},
		{"myName", types.NewArray(named, 3)},
		{"var0", types.NewArray(basic, 3)},
		{"myName", types.NewSlice(named)},
		{"var0", types.NewSlice(basic)},
		// maps tested below
		{"myName", types.NewChan(types.SendRecv, named)},
		{"var0", types.NewChan(types.SendRecv, basic)},
		{"myName", types.NewStruct([]*types.Var{namedfield}, nil)},
		{"var0", types.NewStruct([]*types.Var{basicfield}, nil)},
		{"var0", types.NewSignature(nil, nil, nil, false)},
		{"var0", types.NewInterface(nil, nil)},
	}

	for _, test := range tests {
		var sut varBasenameGen
		result := sut.getBasename(test.typ)

		is.Equal(result, test.expected)
	}
}

func TestGetBasenameForMaps(t *testing.T) {
	is := is.New(t)
	basic1 := types.Typ[types.Int]
	basic2 := types.Typ[types.Uint]
	named1 := makeNamedType("MyType1", basic1)
	named2 := makeNamedType("MyType2", basic2)

	tests := []struct {
		expected string
		key      types.Type
		value    types.Type
	}{
		{"myType1ToMyType2", named1, named2},
		{"var0", named1, basic2},
		{"var0", basic1, named2},
		{"var0", basic1, basic2},
	}

	for _, test := range tests {
		var sut varBasenameGen
		typ := types.NewMap(test.key, test.value)
		result := sut.getBasename(typ)

		is.Equal(result, test.expected)
	}
}

func TestGetBasenameIncrementsGeneratedNames(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		expected string
		typ      types.Type
	}{
		{"var0", types.Typ[types.Int]},
		{"var1", types.Typ[types.Uint]},
	}

	var sut varBasenameGen
	for _, test := range tests {
		result := sut.getBasename(test.typ)

		is.Equal(result, test.expected)
	}
}

func TestGetBasenameDisallowsUnderscores(t *testing.T) {
	is := is.New(t)
	typ := makeNamedType("My_Int", types.Typ[types.Int])

	var sut varBasenameGen
	result := sut.getBasename(typ)

	is.Equal("var0", result)
}

func TestVarNamerGivesExpectedNames(t *testing.T) {
	is := is.New(t)
	basic := types.Typ[types.Int]
	named := makeNamedType("MyInt", basic)
	other := makeNamedType("MyInt", basic)
	tests := []struct {
		expected string
		typ      types.Type
		instance int
	}{
		{"var0", basic, 0},
		{"myInt", named, 0},
		{"var0_1", basic, 1},
		{"myInt_1", named, 1},
		{"myInt_A", other, 0},
		{"myInt_A_1", other, 1},
	}

	// note that the test cases are not independent
	// this is testing the accumulation of state in sut
	sut := newVarNamer(typeutil.MakeHasher())
	for _, test := range tests {
		result := sut.Name(test.typ, test.instance)

		is.Equal(result, test.expected)
	}
}

func TestBuildTypeName(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		expected string
		basename string
		i        int
	}{
		{"base", "base", 0},
		{"base_A", "base", 1},
		{"base_B", "base", 2},
		{"base_Y", "base", 25},
		{"base_Z", "base", 26},
		{"base_BA", "base", 27},
		{"base_BZ", "base", 52},
		{"base_CA", "base", 53},
		{"func_A", "func", 0},
		{"return_A", "return", 0},
	}

	for _, test := range tests {
		result := buildTypeName(test.basename, test.i)

		is.Equal(result, test.expected)
	}
}
