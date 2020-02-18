// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"testing"

	"github.com/clivern/beetle/pkg"
)

// TestDatabase test cases
func TestDatabase(t *testing.T) {
	t.Run("TestDatabase", func(t *testing.T) {
		pkg.Expect(t, true, true)
	})
}
