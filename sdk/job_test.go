// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/clivern/beetle/internal/app/module"
	"github.com/clivern/beetle/pkg"

	"github.com/drone/envsubst"
	"github.com/spf13/viper"
)

// TestJobCRUD test cases
func TestJobCRUD(t *testing.T) {
	testingConfig := "config.testing.yml"
	httpClient := module.NewHTTPClient()

	// LoadConfigFile
	t.Run("LoadConfigFile", func(t *testing.T) {
		fs := module.FileSystem{}

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

	// TestGetJobs
	t.Run("TestGetJobs", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/job",
			`{"jobs": [{"id":1,"uuid":"4f540ab1-2c29-47e6-b900-675312b784d8","payload":"{}","status":"pending","type":"deployment.update","result":"","retry":0,"parent":0,"run_at":null,"created_at":"2020-06-16T18:20:35Z","updated_at":"2020-06-16T18:20:35Z"}]}`,
			http.StatusOK,
		)

		defer srv.Close()

		result, err := GetJobs(context.TODO(), httpClient, srv.URL, "")

		pkg.Expect(t, nil, err)

		pkg.Expect(t, result.Jobs[0].ID, 1)
		pkg.Expect(t, result.Jobs[0].UUID, "4f540ab1-2c29-47e6-b900-675312b784d8")
		pkg.Expect(t, result.Jobs[0].Status, "pending")
		pkg.Expect(t, result.Jobs[0].Type, "deployment.update")
	})

	// TestGetJob
	t.Run("TestGetJob", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/job/4f540ab1-2c29-47e6-b900-675312b784d8",
			`{"id":1,"uuid":"4f540ab1-2c29-47e6-b900-675312b784d8","payload":"{}","status":"pending","type":"deployment.update","result":"","retry":0,"parent":0,"run_at":null,"created_at":"2020-06-16T18:20:35Z","updated_at":"2020-06-16T18:20:35Z"}`,
			http.StatusOK,
		)

		defer srv.Close()

		result, err := GetJob(context.TODO(), httpClient, srv.URL, "4f540ab1-2c29-47e6-b900-675312b784d8", "")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, result.ID, 1)
		pkg.Expect(t, result.UUID, "4f540ab1-2c29-47e6-b900-675312b784d8")
		pkg.Expect(t, result.Status, "pending")
		pkg.Expect(t, result.Type, "deployment.update")
	})

	// TestDeleteJob
	t.Run("TestDeleteJob", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/job/4f540ab1-2c29-47e6-b900-675312b784d8",
			``,
			http.StatusNoContent,
		)

		defer srv.Close()

		result, err := DeleteJob(context.TODO(), httpClient, srv.URL, "4f540ab1-2c29-47e6-b900-675312b784d8", "")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, true, result)
	})
}
