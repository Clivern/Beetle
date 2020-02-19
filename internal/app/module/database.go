// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"github.com/clivern/beetle/internal/app/migration"
	"github.com/clivern/beetle/internal/app/model"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Database struct
type Database struct {
	Connection *gorm.DB
}

// Connect connects to a MySQL database
func (db *Database) Connect(dsn model.DSN) error {
	var err error

	db.Connection, err = gorm.Open(dsn.Driver, dsn.ToString())

	if err != nil {
		return err
	}

	return nil
}

// Migrate migrates the database
func (db *Database) Migrate() {
	db.Connection.AutoMigrate(&migration.Job{})
}

// CreateJob creates a new job
func (db *Database) CreateJob(_ model.Job) (model.Job, error) {
	return model.Job{}, nil
}

// GetJobByID gets a job by id
func (db *Database) GetJobByID(_ int) (model.Job, error) {
	return model.Job{}, nil
}

// GetJobByUUID gets a job by id
func (db *Database) GetJobByUUID(_ string) (model.Job, error) {
	return model.Job{}, nil
}

// GetJobs get a list of jobs
func (db *Database) GetJobs(_ int, _ int) ([]model.Job, error) {
	return []model.Job{}, nil
}

// DeleteJob deletes a job by id
func (db *Database) DeleteJob(_ int) (bool, error) {
	return true, nil
}

// UpdateJobByID updates a job by ID
func (db *Database) UpdateJobByID(_ model.Job) (bool, error) {
	return true, nil
}

// UpdateJobByUUID updates a job by UUID
func (db *Database) UpdateJobByUUID(_ model.Job) (bool, error) {
	return true, nil
}

// Close closes MySQL database connection
func (db *Database) Close() error {
	return db.Connection.Close()
}
