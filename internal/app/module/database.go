// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Database struct
type Database struct {
	Connection *sql.DB
}

// DSN struct
type DSN struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

// ToString gets the dsn string
func (dsn *DSN) ToString() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		dsn.Username,
		dsn.Password,
		dsn.Hostname,
		dsn.Port,
		dsn.Database,
	)
}

// Connect connects to a MySQL database
func (db *Database) Connect(dsn DSN) (bool, error) {
	var err error

	db.Connection, err = sql.Open("mysql", dsn.ToString())

	if err != nil {
		return false, err
	}

	return true, nil
}

// Close closes MySQL database connection
func (db *Database) Close() error {
	return db.Connection.Close()
}
