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

	"golang.org/x/tools/go/types/typeutil"
)

type varNamer struct {
	baseNamer   varBasenameGen
	basenameMap map[string][]types.Type
	varNames    *typeStringMap
}

func newVarNamer(hasher typeutil.Hasher) *varNamer {
	return &varNamer{
		basenameMap: make(map[string][]types.Type),
		varNames:    newTypeStringMap(hasher),
	}
}

func (v *varNamer) Name(typ types.Type, instance int) string {
	name := v.varNames.Get(typ)

	if name == "" {
		basename := v.baseNamer.getBasename(typ)
		idx, found := findType(v.basenameMap[basename], typ)
		if !found {
			idx = len(v.basenameMap[basename])
			v.basenameMap[basename] = append(v.basenameMap[basename], typ)
		}

		name = buildTypeName(basename, idx)
		v.varNames.Set(typ, name)
	}

	return buildFullName(name, instance)
}

type varBasenameGen uint

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
		varname = generateVarName(uint(*v))
		*v++
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

func findType(typs []types.Type, typ types.Type) (int, bool) {
	for i := range typs {
		if types.Identical(typ, typs[i]) {
			return i, true
		}
	}

	return 0, false
}

func buildTypeName(basename string, i int) string {
	if isKeyword(basename) {
		i++
	}

	if i == 0 {
		return basename
	}
	i--

	var suffix []rune
	for ; i >= 26; i = i / 26 {
		suffix = append(suffix, suffixMap[i%26])
	}
	suffix = append(suffix, suffixMap[i%26])

	var out bytes.Buffer
	out.WriteString(basename)
	out.WriteRune('_')
	for j := len(suffix) - 1; j >= 0; j-- {
		out.WriteRune(suffix[j])
	}

	return out.String()
}

var suffixMap = [...]rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

func buildFullName(name string, i int) string {
	if i == 0 {
		return name
	}

	var out bytes.Buffer
	out.WriteString(name)
	out.WriteRune('_')
	out.WriteString(strconv.Itoa(i))

	return out.String()
}

func isKeyword(name string) bool {
	switch name {
	case "break", "case", "chan", "const", "continue", "default", "defer", "else":
		fallthrough
	case "fallthrough", "for", "func", "go", "goto", "if", "import", "interface":
		fallthrough
	case "map", "package", "range", "return", "select", "struct", "switch", "type", "var":
		return true
	}

	return false
}
