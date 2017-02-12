// Copyright Steven Bosnick 2017. All rights reserved.
// Use of this source code is governed by the GNU General Public License version 3.
// See the file COPYING for your rights under that license.

package depend

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingNodeRequiresNothing(t *testing.T) {
	container := &Container{}

	sut := missingNode{container: container}

	assert.Len(t, sut.requires(), 0, "missingNode unexpectedly requires some types")
}

func TestMissingNodeProvidesNothing(t *testing.T) {
	container := &Container{}

	sut := missingNode{container: container}

	assert.Len(t, sut.provides(), 0, "missingNode unexpectedly provides some types")
}
