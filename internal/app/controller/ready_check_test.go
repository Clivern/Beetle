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
	"github.com/clivern/beetle/pkg"

	"github.com/drone/envsubst"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// TestReadyCheck test cases
func TestReadyCheck(t *testing.T) {
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

	// TestReadyCheckController
	t.Run("TestReadyCheckController", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = ioutil.Discard
		gin.DisableConsoleColor()

		router := gin.Default()

		router.GET("/_ready", ReadyCheck)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/_ready", nil)
		router.ServeHTTP(w, req)

		pkg.Expect(t, viper.GetString("app.mode"), "test")
		pkg.Expect(t, w.Code, 200)
		pkg.Expect(t, strings.TrimSpace(w.Body.String()), `{"status":"ok"}`)
	})
}
