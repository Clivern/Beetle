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

// TestApplicationCRUD test cases
func TestApplicationCRUD(t *testing.T) {
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

	// TestGetApplications
	t.Run("TestGetApplications", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/cluster/production/namespace/default/app",
			`{"applications":[{"id":"toad","name":"Toad App","format":"clivern/toad:release-[.Release]","containers":[{"name":"toad","image":"clivern/toad:release-0.2.3","version":"0.2.3","deployment":{"name":"toad-deployment","uid":"0f77903a-ce69-4aa5-a025-cad4b4a3209e"}}]}]}`,
			http.StatusOK,
		)

		defer srv.Close()

		httpClient.SetAPIURL(srv.URL)
		result, err := httpClient.GetApplications(context.TODO(), "production", "default")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, result, model.Applications{
			Applications: []model.Application{
				model.Application{
					ID:     "toad",
					Name:   "Toad App",
					Format: "clivern/toad:release-[.Release]",
					Containers: []model.Container{
						model.Container{
							Name:    "toad",
							Image:   "clivern/toad:release-0.2.3",
							Version: "0.2.3",
							Deployment: model.Deployment{
								Name: "toad-deployment",
								UID:  "0f77903a-ce69-4aa5-a025-cad4b4a3209e",
							},
						},
					},
				},
			},
		})
	})

	// TestGetApplication
	t.Run("TestGetApplication", func(t *testing.T) {
		srv := pkg.ServerMock(
			"/api/v1/cluster/production/namespace/default/app/toad",
			`{"id":"toad","name":"Toad App","format":"clivern/toad:release-[.Release]","containers":[{"name":"toad","image":"clivern/toad:release-0.2.3","version":"0.2.3","deployment":{"name":"toad-deployment","uid":"0f77903a-ce69-4aa5-a025-cad4b4a3209e"}}]}`,
			http.StatusOK,
		)

		defer srv.Close()

		httpClient.SetAPIURL(srv.URL)
		result, err := httpClient.GetApplication(context.TODO(), "production", "default", "toad")

		pkg.Expect(t, nil, err)
		pkg.Expect(t, result, model.Application{
			ID:     "toad",
			Name:   "Toad App",
			Format: "clivern/toad:release-[.Release]",
			Containers: []model.Container{
				model.Container{
					Name:    "toad",
					Image:   "clivern/toad:release-0.2.3",
					Version: "0.2.3",
					Deployment: model.Deployment{
						Name: "toad-deployment",
						UID:  "0f77903a-ce69-4aa5-a025-cad4b4a3209e",
					},
				},
			},
		})
	})
}
