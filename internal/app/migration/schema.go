// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package migration

var (
	// MigrationTableStatement mysql statement to create migrations table
	MigrationTableStatement = `DROP TABLE IF EXISTS migrations
	CREATE TABLE IF NOT EXISTS migrations (
	  id int(11) NOT NULL AUTO_INCREMENT,
	  file varchar(150) NOT NULL,
	  run_at datetime(6) NOT NULL,
	  PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1`

	// JobTableStatement mysql statement to create jobs table
	JobTableStatement = `DROP TABLE IF EXISTS jobs
	CREATE TABLE IF NOT EXISTS jobs (
	  id int(11) NOT NULL AUTO_INCREMENT,
	  run_at datetime(6) NOT NULL,
	  PRIMARY KEY (id)
	) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1`

	// DatabaseSchema database schema
	DatabaseSchema = map[string]string{
		"create_migrations_table": MigrationTableStatement,
		"create_jobs_table":       JobTableStatement,
	}
)
