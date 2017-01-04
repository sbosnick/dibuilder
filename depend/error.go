// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import "errors"

var (
	// ErrNoRoot is the error used to indicate a root has not been
	// set on a Container for an operation that requires such a root.
	ErrNoRoot = errors.New("no root set for container")
)
