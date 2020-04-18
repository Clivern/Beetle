// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package kubernetes

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/clivern/beetle/internal/app/module"
	"github.com/clivern/beetle/pkg"

	"github.com/drone/envsubst"
	"github.com/spf13/viper"
)

// TestCluster test cases
func TestCluster(t *testing.T) {
	testingConfig := "config.testing.yml"

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
		clusters, err := GetClusters()

		pkg.Expect(t, nil, err)
		pkg.Expect(t, clusters[0].Name, "production")
		pkg.Expect(t, clusters[0].Kubeconfig, "/app/configs/production-cluster-kubeconfig.yaml")
		pkg.Expect(t, clusters[0].ConfigMapName, "beetle-configs")

		pkg.Expect(t, clusters[1].Name, "staging")
		pkg.Expect(t, clusters[1].Kubeconfig, "/app/configs/staging-cluster-kubeconfig.yaml")
		pkg.Expect(t, clusters[1].ConfigMapName, "beetle-configs")
	})

	// TestGetCluster
	t.Run("TestGetCluster", func(t *testing.T) {
		cluster, err := GetCluster("production")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, cluster.Name, "production")
		pkg.Expect(t, cluster.Kubeconfig, "/app/configs/production-cluster-kubeconfig.yaml")
		pkg.Expect(t, cluster.ConfigMapName, "beetle-configs")

		cluster, err = GetCluster("not-found")

		pkg.Expect(t, fmt.Errorf("Unable to find cluster not-found"), err)
		pkg.Expect(t, cluster.Name, "")
		pkg.Expect(t, cluster.Kubeconfig, "")
		pkg.Expect(t, cluster.ConfigMapName, "")
	})
}
