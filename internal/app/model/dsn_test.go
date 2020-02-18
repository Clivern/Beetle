// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"testing"

	"github.com/clivern/beetle/pkg"
)

// TestDsnToString test cases
func TestDsnToString(t *testing.T) {
	t.Run("TestDsnToString", func(t *testing.T) {
		dsn := DSN{
			Username: "root",
			Password: "root",
			Hostname: "127.0.0.1",
			Port:     3306,
			Database: "beetle",
		}
		pkg.Expect(t, "root:root@tcp(127.0.0.1:3306)/beetle", dsn.ToString())
	})
}
