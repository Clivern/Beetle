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

	"github.com/clivern/beetle/core/model"
	"github.com/clivern/beetle/core/module"
	"github.com/clivern/beetle/pkg"

	"github.com/drone/envsubst"
	"github.com/spf13/viper"
)

// TestClusterCRUD test cases
func TestClusterCRUD(t *testing.T) {
	testingConfig := "config.testing.yml"

	httpClient := Client{}
	httpClient.SetHTTPClient(module.NewHTTPClient(20))
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

	// TestGetClusters
	t.Run("TestGetClusters", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/cluster",
			`{"clusters": [{"name": "staging", "health": false},{"name": "production", "health": true}]}`,
			http.StatusOK,
		)

		defer srv.Close()

		httpClient.SetAPIURL(srv.URL)
		result, err := httpClient.GetClusters(context.TODO())

		pkg.Expect(t, nil, err)
		pkg.Expect(t, result, model.Clusters{
			Clusters: []model.Cluster{
				model.Cluster{Name: "staging", Health: false},
				model.Cluster{Name: "production", Health: true},
			},
		})
	})

	// TestGetCluster
	t.Run("TestGetCluster", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/cluster/staging",
			`{"name": "staging", "health": false}`,
			http.StatusOK,
		)

		defer srv.Close()

		httpClient.SetAPIURL(srv.URL)
		result, err := httpClient.GetCluster(context.TODO(), "staging")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, result, model.Cluster{Name: "staging", Health: false})
	})
}
