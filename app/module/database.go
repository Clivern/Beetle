// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"fmt"
	"github.com/clivern/beetle/app/migration"
	"github.com/clivern/beetle/app/model"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

// Database struct
type Database struct {
	Connection *gorm.DB
}

// Connect connects to a MySQL database
func (db *Database) Connect(dsn model.DSN) error {
	var err error

	// Reuse db connections http://go-database-sql.org/surprises.html
	if db.Ping() == nil {
		return nil
	}

	db.Connection, err = gorm.Open(dsn.Driver, dsn.ToString())

	if err != nil {
		return err
	}

	return nil
}

// Ping check the db connection
func (db *Database) Ping() error {

	if db.Connection == nil {
		return fmt.Errorf("No DB Connections Found")
	}

	err := db.Connection.DB().Ping()

	if err != nil {
		return err
	}

	// Cleanup stale connections http://go-database-sql.org/surprises.html
	db.Connection.DB().SetMaxOpenConns(5)
	db.Connection.DB().SetConnMaxLifetime(time.Duration(10) * time.Second)
	dbStats := db.Connection.DB().Stats()

	log.WithFields(log.Fields{
		"dbStats.maxOpenConnections": int(dbStats.MaxOpenConnections),
		"dbStats.openConnections":    int(dbStats.OpenConnections),
		"dbStats.inUse":              int(dbStats.InUse),
		"dbStats.idle":               int(dbStats.Idle),
	}).Debug(`Open DB Connection`)

	return nil
}

// AutoConnect connects to a MySQL database using loaded configs
func (db *Database) AutoConnect() error {
	var err error

	// Reuse db connections http://go-database-sql.org/surprises.html
	if db.Ping() == nil {
		return nil
	}

	dsn := model.DSN{
		Driver:   viper.GetString("app.database.driver"),
		Username: viper.GetString("app.database.username"),
		Password: viper.GetString("app.database.password"),
		Hostname: viper.GetString("app.database.host"),
		Port:     viper.GetInt("app.database.port"),
		Name:     viper.GetString("app.database.name"),
	}

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
	db.Connection.DropTableIfExists(&migration.Job{})
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

// JobExistByID check if job exists
func (db *Database) JobExistByID(id int) bool {
	job := model.Job{}

	db.Connection.Where("id = ?", id).First(&job)

	return job.ID > 0
}

// GetJobByID gets a job by id
func (db *Database) GetJobByID(id int) model.Job {
	job := model.Job{}

	db.Connection.Where("id = ?", id).First(&job)

	return job
}

// GetJobs gets jobs
func (db *Database) GetJobs() []model.Job {
	jobs := []model.Job{}

	db.Connection.Select("*").Find(&jobs)

	return jobs
}

// JobExistByUUID check if job exists
func (db *Database) JobExistByUUID(uuid string) bool {
	job := model.Job{}

	db.Connection.Where("uuid = ?", uuid).First(&job)

	return job.ID > 0
}

// GetJobByUUID gets a job by uuid
func (db *Database) GetJobByUUID(uuid string) model.Job {
	job := model.Job{}

	db.Connection.Where("uuid = ?", uuid).First(&job)

	return job
}

// GetPendingJobByType gets a job by uuid
func (db *Database) GetPendingJobByType(jobType string) model.Job {
	job := model.Job{}

	db.Connection.Where("status = ? AND type = ?", model.JobPending, jobType).First(&job)

	return job
}

// CountJobs count jobs by status
func (db *Database) CountJobs(status string) int {
	count := 0

	db.Connection.Model(&model.Job{}).Where("status = ?", status).Count(&count)

	return count
}

// DeleteJobByID deletes a job by id
func (db *Database) DeleteJobByID(id int) {
	db.Connection.Unscoped().Where("id=?", id).Delete(&migration.Job{})
}

// DeleteJobByUUID deletes a job by uuid
func (db *Database) DeleteJobByUUID(uuid string) {
	db.Connection.Unscoped().Where("uuid=?", uuid).Delete(&migration.Job{})
}

// UpdateJobByID updates a job by ID
func (db *Database) UpdateJobByID(job *model.Job) {
	db.Connection.Save(&job)
}

// Close closes MySQL database connection
func (db *Database) Close() error {
	return db.Connection.Close()
}

// ReleaseChildJobs count jobs by status
func (db *Database) ReleaseChildJobs(parentID int) {
	db.Connection.Model(&model.Job{}).Where(
		"parent = ? AND status = ?",
		parentID,
		model.JobOnHold,
	).Update("status", model.JobPending)
}
