// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"bytes"
	"go/types"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type varBasenameGen struct {
	next uint
}

func (v *varBasenameGen) getBasename(typ types.Type) string {
	var named *types.Named
	var varname string

	switch typ := typ.(type) {
	case *types.Map:
		varname = makeMapVarName(typ.Key(), typ.Elem())
	case *types.Named:
		named = typ
	case elemProvider:
		elem := typ.Elem()
		if elem, ok := elem.(*types.Named); ok {
			named = elem
		}
	case *types.Struct:
		if typ.NumFields() > 0 {
			if fieldtype, ok := typ.Field(0).Type().(*types.Named); ok {
				named = fieldtype
			}
		}
	}

	if named != nil {
		name := named.Obj().Name()
		if !hasUnderscore(name) {
			varname = toLowercaseLeading(name)
		}
	}

	if varname == "" {
		varname = generateVarName(v.next)
		v.next++
	}

	return varname
}

// This interface is satisfied by types.Array, types.Chan, type.Map, types.Pointer, and types.Slice
type elemProvider interface {
	Elem() types.Type
}

func getVarPrefix(typ types.Type) string {
	switch typ.(type) {
	case *types.Basic:
		return "b"
	case *types.Map:
		return "m"
	case *types.Struct:
		return "s"
	case *types.Signature:
		return "sig"
	case *types.Interface:
		return "int"
	}
	return ""
}

func toLowercaseLeading(str string) string {
	var in *strings.Reader = strings.NewReader(str)
	var out bytes.Buffer

	first, _, err := in.ReadRune()
	if err != nil || !unicode.IsLetter(first) {
		return ""
	}

	_, err = out.WriteRune(unicode.ToLower(first))
	if err != nil {
		return ""
	}

	_, err = io.Copy(&out, in)
	if err != nil {
		return ""
	}

	return out.String()
}

func generateVarName(next uint) string {
	var out bytes.Buffer

	out.WriteString("var")
	out.WriteString(strconv.FormatUint(uint64(next), 10))

	return out.String()
}

func makeMapVarName(keytype types.Type, valtype types.Type) string {
	var keyname string
	var valname string

	if keytype, ok := keytype.(*types.Named); ok {
		keyname = toLowercaseLeading(keytype.Obj().Name())
	}

	if valtype, ok := valtype.(*types.Named); ok {
		valname = valtype.Obj().Name()
	}

	if keyname == "" || valname == "" || hasUnderscore(keyname) || hasUnderscore(valname) {
		return ""
	}

	var out bytes.Buffer
	out.WriteString(keyname)
	out.WriteString("To")
	out.WriteString(valname)
	return out.String()
}

func hasUnderscore(s string) bool {
	return strings.ContainsRune(s, '_')
}
