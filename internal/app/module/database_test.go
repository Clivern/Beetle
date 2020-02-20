// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/pkg"

	"github.com/drone/envsubst"
	"github.com/spf13/viper"
)

var testingConfig = "config.testing.yml"

// TestDatabase test cases
func TestDatabase(t *testing.T) {
	// LoadConfigFile
	t.Run("LoadConfigFile", func(t *testing.T) {
		fs := FileSystem{}

		dir, _ := os.Getwd()
		configFile := fmt.Sprintf("%s/%s", dir, testingConfig)

		for {
			if fs.FileExists(configFile) {
				break
			}
			dir = filepath.Dir(dir)
			configFile = fmt.Sprintf("%s/%s", dir, testingConfig)
		}

		t.Logf("Load Config File %s", configFile)

		configUnparsed, _ := ioutil.ReadFile(configFile)
		configParsed, _ := envsubst.EvalEnv(string(configUnparsed))
		viper.SetConfigType("yaml")
		viper.ReadConfig(bytes.NewBuffer([]byte(configParsed)))
	})

	// TestDatabaseConnection
	t.Run("TestDatabaseConnection", func(t *testing.T) {
		db := Database{}
		err := db.Connect(model.DSN{
			Driver:   viper.GetString("app.database.driver"),
			Username: viper.GetString("app.database.username"),
			Password: viper.GetString("app.database.password"),
			Hostname: viper.GetString("app.database.host"),
			Port:     viper.GetInt("app.database.port"),
			Name:     viper.GetString("app.database.name"),
		})
		pkg.Expect(t, nil, err)

		defer db.Close()
		pkg.Expect(t, true, db.Rollback())
		pkg.Expect(t, true, db.Migrate())
		pkg.Expect(t, true, db.HasTable("jobs"))
	})

	// TestJobCRUD
	t.Run("TestJobCRUD", func(t *testing.T) {
		db := Database{}
		err := db.Connect(model.DSN{
			Driver:   viper.GetString("app.database.driver"),
			Username: viper.GetString("app.database.username"),
			Password: viper.GetString("app.database.password"),
			Hostname: viper.GetString("app.database.host"),
			Port:     viper.GetInt("app.database.port"),
			Name:     viper.GetString("app.database.name"),
		})
		pkg.Expect(t, nil, err)

		defer db.Close()
		pkg.Expect(t, true, db.Rollback())
		pkg.Expect(t, true, db.Migrate())
		pkg.Expect(t, true, db.HasTable("jobs"))

		// Delete the job if it exists
		db.DeleteJobByID(1)
		db.DeleteJobByUUID("dddde755-5f99-4e51-a517-77878986a07e")

		// Create the job
		job := db.CreateJob(&model.Job{
			UUID: "dddde755-5f99-4e51-a517-77878986a07e",
		})

		pkg.Expect(t, 1, job.ID)
		pkg.Expect(t, "dddde755-5f99-4e51-a517-77878986a07e", job.UUID)

		job1 := db.GetJobByID(1)
		job2 := db.GetJobByUUID("dddde755-5f99-4e51-a517-77878986a07e")

		pkg.Expect(t, job1.ID, job2.ID)
		pkg.Expect(t, job1.UUID, job2.UUID)
	})
}
