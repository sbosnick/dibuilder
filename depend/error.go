// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"bytes"
	"errors"
	"go/token"
	"go/types"
)

var (
	// ErrNoRoot is the error used to indicate a root has not been
	// set on a Container for an operation that requires such a root.
	ErrNoRoot = errors.New("no root set for container")

	// ErrRootAlreadySet is the error used to indicate that a root has
	// already been set on a Container for an operation that only allows
	// the root to be set once.
	ErrRootAlreadySet = errors.New("root already set for container")
)

// An Error represents an error with an associated position in an
// implicit token.FileSet.
type Error interface {
	error

	// Pos returns the position associated with this error.
	Pos() token.Pos

	// ErrorWithPosition returns the error message for the error preceded with
	// the string representation of the position as given by the token.FileSet.
	ErrorWithPosition(fileSet *token.FileSet) string
}

// InvalidFuncError records an error with an attempt to add an invalid types.Func
// to a Container. InvalidFuncError implements Error.
type InvalidFuncError struct {
	pos      token.Pos
	funcName string
	reason   string
}

func (ife *InvalidFuncError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString("Invalid func for Container (")
	buffer.WriteString(ife.funcName)
	buffer.WriteString("): ")
	buffer.WriteString(ife.reason)
	return buffer.String()
}

func (ife *InvalidFuncError) Pos() token.Pos {
	return ife.pos
}

func (ife *InvalidFuncError) ErrorWithPosition(fileSet *token.FileSet) string {
	var buffer bytes.Buffer
	buffer.WriteString(fileSet.Position(ife.pos).String())
	buffer.WriteString(": Invalid func for Container (")
	buffer.WriteString(ife.funcName)
	buffer.WriteString("): ")
	buffer.WriteString(ife.reason)
	return buffer.String()
}

func newInvalidFuncError(function *types.Func, reason string) *InvalidFuncError {
	return &InvalidFuncError{
		pos:      function.Pos(),
		funcName: function.Name(),
		reason:   reason,
	}
}

var _ Error = &InvalidFuncError{}
