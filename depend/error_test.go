// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"go/token"
	"go/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidFuncErrorIncludesFuncNameAndReason(t *testing.T) {
	name := "MyFunction"
	reason := "No good reason."
	sig := types.NewSignature(nil, types.NewTuple(), types.NewTuple(), false)
	function := types.NewFunc(token.NoPos, nil, name, sig)

	sut := newInvalidFuncError(function, reason)
	result := sut.Error()

	assert.Contains(t, result, name, "Error string did not include the expected function name")
	assert.Contains(t, result, reason, "Error string did not include the expected reason")
}

func TestInvalidFuncErrorRecordsPosOfFunc(t *testing.T) {
	fileset := token.NewFileSet()
	expectedPos := fileset.AddFile("myfile.go", -1, 30).Pos(10)
	sig := types.NewSignature(nil, types.NewTuple(), types.NewTuple(), false)
	function := types.NewFunc(expectedPos, nil, "MyFunction", sig)

	sut := newInvalidFuncError(function, "No good reason.")
	pos := sut.Pos()

	assert.Equal(t, expectedPos, pos, "Unexpected position recorded in InvalidFuncError")
}

func TestInvaidFuncErrorIncludeFilenameInErrorWithPositon(t *testing.T) {
	filename := "myfile.go"
	fileset := token.NewFileSet()
	pos := fileset.AddFile(filename, -1, 30).Pos(10)
	sig := types.NewSignature(nil, types.NewTuple(), types.NewTuple(), false)
	function := types.NewFunc(pos, nil, "MyFunction", sig)

	sut := newInvalidFuncError(function, "No good reason.")
	result := sut.ErrorWithPosition(fileset)

	assert.Contains(t, result, filename, "Error string did not include the expected filename")
}
