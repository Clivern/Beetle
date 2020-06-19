// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package util

import (
	"testing"

	"github.com/clivern/beetle/app/model"
	"github.com/clivern/beetle/pkg"
)

// TestInArray test cases
func TestInArray(t *testing.T) {
	// TestInArray
	t.Run("TestInArray", func(t *testing.T) {
		pkg.Expect(t, InArray("A", []string{"A", "B", "C", "D"}), true)
		pkg.Expect(t, InArray("B", []string{"A", "B", "C", "D"}), true)
		pkg.Expect(t, InArray("H", []string{"A", "B", "C", "D"}), false)
		pkg.Expect(t, InArray(1, []int{2, 3, 1}), true)
		pkg.Expect(t, InArray(9, []int{2, 3, 1}), false)

		payload := []model.PatchStringValue{
			model.PatchStringValue{
				Op:    "replace",
				Path:  "/spec/template/spec/containers/0/image",
				Value: "clivern/toad:release-0.2.4",
			},
		}

		data, err := ConvertToJSON(payload)

		pkg.Expect(t, data, `[{"op":"replace","path":"/spec/template/spec/containers/0/image","value":"clivern/toad:release-0.2.4"}]`)
		pkg.Expect(t, err, nil)
	})
}
