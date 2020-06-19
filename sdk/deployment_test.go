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

	"github.com/clivern/beetle/app/model"
	"github.com/clivern/beetle/app/module"
	"github.com/clivern/beetle/pkg"

	"github.com/drone/envsubst"
	"github.com/spf13/viper"
)

// TestDeploymentCRUD test cases
func TestDeploymentCRUD(t *testing.T) {
	testingConfig := "config.testing.yml"

	httpClient := Client{}
	httpClient.SetHTTPClient(module.NewHTTPClient())
	httpClient.SetAPIKey("")

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

	// TestCreateDeployment
	t.Run("TestCreateDeployment", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/cluster/production/namespace/default/app/toad",
			`{"id":1,"uuid":"4f540ab1-2c29-47e6-b900-675312b784d8","status":"pending","type":"deployment.update","created_at":"2020-06-16T18:20:35Z"}`,
			http.StatusAccepted,
		)

		defer srv.Close()

		deploymentRequest := model.DeploymentRequest{
			Cluster:     "production",
			Namespace:   "default",
			Application: "toad",
			Version:     "1.0.0",
			Strategy:    "recreate",
		}

		httpClient.SetAPIURL(srv.URL)
		result, err := httpClient.CreateDeployment(context.TODO(), deploymentRequest)

		pkg.Expect(t, err, nil)
		pkg.Expect(t, 1, result.ID)
		pkg.Expect(t, "4f540ab1-2c29-47e6-b900-675312b784d8", result.UUID)
		pkg.Expect(t, "pending", result.Status)
		pkg.Expect(t, "deployment.update", result.Type)
	})
}
