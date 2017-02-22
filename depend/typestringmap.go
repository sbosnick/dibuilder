// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/types"

	"golang.org/x/tools/go/types/typeutil"
)

type typeStringMap struct {
	typeMap typeutil.Map
}

func newTypeStringMap(hasher typeutil.Hasher) *typeStringMap {
	tsm := typeStringMap{}
	tsm.typeMap.SetHasher(hasher)
	return &tsm
}

func (m *typeStringMap) Set(key types.Type, value string) {
	if m == nil {
		panic("Set called on a nil typeStringMap.")
	}

	m.typeMap.Set(key, value)
}

func (m *typeStringMap) Get(key types.Type) string {
	if m == nil {
		return ""
	}

	value := m.typeMap.At(key)
	if value == nil {
		return ""
	}

	return value.(string)
}
