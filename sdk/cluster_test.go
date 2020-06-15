// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package sdk

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/clivern/beetle/internal/app/model"
	"github.com/clivern/beetle/internal/app/module"
	"github.com/clivern/beetle/pkg"

	"github.com/drone/envsubst"
	"github.com/spf13/viper"
)

// TestClusterCMD test cases
func TestClusterCMD(t *testing.T) {
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

	// TestGetClusters
	t.Run("TestGetClusters", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/cluster",
			`{"clusters": [{"name": "staging", "health": false},{"name": "production", "health": true}]}`,
		)

		defer srv.Close()

		result, err := GetClusters(context.TODO(), httpClient, srv.URL, "")

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
		)

		defer srv.Close()

		result, err := GetCluster(context.TODO(), httpClient, srv.URL, "staging", "")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, result, model.Cluster{Name: "staging", Health: false})
	})
}
