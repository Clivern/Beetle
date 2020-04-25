// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"strings"
	"testing"

	"github.com/clivern/beetle/pkg"
)

// TestRemote test cases
func TestRemote(t *testing.T) {
	t.Run("TestRemote", func(t *testing.T) {
		result, err := GetLatestRelease()
		pkg.Expect(t, true, strings.Contains(result.Name, "."))
		pkg.Expect(t, true, strings.Contains(result.TagName, "."))
		pkg.Expect(t, nil, err)
	})
}
