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

// TestNamespaceCRUD test cases
func TestNamespaceCRUD(t *testing.T) {
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

	// TestGetNamespaces
	t.Run("TestGetNamespaces", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/cluster/production/namespace",
			`{"namespaces": [{"name": "default","uid": "f03ea2f1-bc1c-4563-b9c7-4413dffc18db","status": "active"},{"name": "kube-node-lease","uid": "398c907f-d888-455d-871d-145752f9ca73","status": "active"}]}`,
		)

		defer srv.Close()

		result, err := GetNamespaces(context.TODO(), httpClient, srv.URL, "production", "")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, result, model.Namespaces{
			Namespaces: []model.Namespace{
				model.Namespace{Name: "default", UID: "f03ea2f1-bc1c-4563-b9c7-4413dffc18db", Status: "active"},
				model.Namespace{Name: "kube-node-lease", UID: "398c907f-d888-455d-871d-145752f9ca73", Status: "active"},
			},
		})
	})

	// TestGetNamespace
	t.Run("TestGetNamespace", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/cluster/production/namespace/default",
			`{"name":"default","status":"active","uid":"f03ea2f1-bc1c-4563-b9c7-4413dffc18db"}`,
		)

		defer srv.Close()

		result, err := GetNamespace(context.TODO(), httpClient, srv.URL, "production", "default", "")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, result, model.Namespace{Name: "default", UID: "f03ea2f1-bc1c-4563-b9c7-4413dffc18db", Status: "active"})
	})
}
