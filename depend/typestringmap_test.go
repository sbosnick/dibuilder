// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/types"
	"testing"

	"golang.org/x/tools/go/types/typeutil"

	"github.com/cheekybits/is"
)

func TestGetNilTypeStringMapIsEmpty(t *testing.T) {
	is := is.New(t)

	var sut *typeStringMap
	result := sut.Get(types.Typ[types.Int])

	is.Equal(result, "")
}

func TestSetNilTypeStringMapPanics(t *testing.T) {
	is := is.New(t)

	var sut *typeStringMap
	is.Panic(func() { sut.Set(types.Typ[types.Int], "doit") })
}

func TestGetEmptyTypeStringMapIsEmpty(t *testing.T) {
	is := is.New(t)

	sut := newTypeStringMap(typeutil.MakeHasher())
	result := sut.Get(types.Typ[types.Int])

	is.Equal(result, "")
}

func TestGetSetTypeStringMapIsSetValue(t *testing.T) {
	is := is.New(t)
	typ := types.Typ[types.Int]
	expected := "myname"

	sut := newTypeStringMap(typeutil.MakeHasher())
	sut.Set(typ, expected)
	result := sut.Get(typ)

	is.Equal(result, expected)
}
