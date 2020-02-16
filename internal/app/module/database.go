// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

// Database struct
type Database struct{}

// Connect Ccnnects to a MySQL database
func (db *Database) Connect() (bool, error) {
	return true, nil
}
