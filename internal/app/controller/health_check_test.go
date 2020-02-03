// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package controller

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/clivern/beetle/internal/app/module"

	"github.com/drone/envsubst"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var testingConfig = "config.testing.yml"

// TestHealthCheck test cases
func TestHealthCheck(t *testing.T) {
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

		got := viper.GetString("app.mode")
		want := "test"
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	// TestHealthCheckController
	t.Run("TestHealthCheckController", func(t *testing.T) {
		router := gin.Default()

		router.GET("/_healthcheck", HealthCheck)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/_healthcheck", nil)
		router.ServeHTTP(w, req)

		got1 := w.Code
		want1 := 200

		if got1 != want1 {
			t.Errorf("got %v, want %v", got1, want1)
		}

		got2 := strings.TrimSpace(w.Body.String())
		want2 := `{"status":"ok"}`

		if got2 != want2 {
			t.Errorf("got %v, want %v", got2, want2)
		}
	})
}
