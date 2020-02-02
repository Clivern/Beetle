// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package main

import (
	"github.com/nbio/st"
	"testing"
)

// TestPkgName test cases
func TestPkgName(t *testing.T) {
	st.Expect(t, "Beetle", "Beetle")
}
