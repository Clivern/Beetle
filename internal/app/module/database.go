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
func (db *Database) Migrate() bool {
	status := true
	db.Connection.AutoMigrate(&migration.Job{})
	status = status && db.Connection.HasTable(&migration.Job{})
	return status
}

// Rollback drop tables
func (db *Database) Rollback() bool {
	status := true
	db.Connection.DropTable(&migration.Job{})
	status = status && !db.Connection.HasTable(&migration.Job{})
	return status
}

// HasTable checks if table exists
func (db *Database) HasTable(table string) bool {
	return db.Connection.HasTable(table)
}

// CreateJob creates a new job
func (db *Database) CreateJob(job *model.Job) *model.Job {
	db.Connection.Create(job)
	return job
}

// GetJobByID gets a job by id
func (db *Database) GetJobByID(id int) model.Job {
	job := model.Job{}
	db.Connection.Where("id = ?", id).First(&job)
	return job
}

// GetJobByUUID gets a job by uuid
func (db *Database) GetJobByUUID(uuid string) model.Job {
	job := model.Job{}
	db.Connection.Where("uuid = ?", uuid).First(&job)
	return job
}

// DeleteJobByID deletes a job by id
func (db *Database) DeleteJobByID(id int) {
	db.Connection.Where("id=?", id).Delete(&migration.Job{})
}

// DeleteJobByUUID deletes a job by uuid
func (db *Database) DeleteJobByUUID(uuid string) {
	db.Connection.Where("uuid=?", uuid).Delete(&migration.Job{})
}

// UpdateJobByID updates a job by ID
func (db *Database) UpdateJobByID(job *model.Job) {
	db.Connection.Save(&job)
}

// Close closes MySQL database connection
func (db *Database) Close() error {
	return db.Connection.Close()
}
