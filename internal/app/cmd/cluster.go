// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/clivern/beetle/internal/app/module"

	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Display clusters names and health",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var result [][]string

		url := os.Getenv("REMOTE_BEETLE_URL")
		token := os.Getenv("REMOTE_BEETLE_TOKEN")

		if url == "" {
			module.DrawTable(
				[]string{"Cluster", "Health"},
				[][]string{{"Error! beetle url is missing ($ export REMOTE_BEETLE_URL=http[s]://remote_url) is required", ""}},
			)
			return
		}

		httpClient := module.NewHTTPClient()

		if len(args) > 0 {
			err, result = getCluster(httpClient, args[0], url, token)
		} else {
			err, result = getClusters(httpClient, url, token)
		}

		if err != nil {
			module.DrawTable(
				[]string{"Cluster", "Health"},
				[][]string{{fmt.Sprintf("Error! %s", err.Error()), ""}},
			)
			return
		}

		module.DrawTable(
			[]string{"Cluster", "Health"},
			result,
		)
	},
}

func init() {
	shellCmd.AddCommand(clusterCmd)
}

// getClusters Get Clusters List
func getClusters(httpClient *module.HTTPClient, beetleURL, token string) (error, [][]string) {
	response, err := httpClient.Get(
		context.Background(),
		fmt.Sprintf("%s/api/v1/cluster", beetleURL),
		map[string]string{},
		map[string]string{"X-AUTH-TOKEN": token},
	)

	if httpClient.GetStatusCode(response) != http.StatusOK || err != nil {
		return fmt.Errorf("Unable to fetch remote data"), [][]string{}
	}

	_, err = httpClient.ToString(response)

	if err != nil {
		return fmt.Errorf("Invalid response"), [][]string{}
	}

	// @TODO
	// convert json response to struct
	// then into [][]string{}

	return nil, [][]string{
		{"staging", "down"},
		{"production", "up"},
	}
}

// getCluster Get Cluster
func getCluster(httpClient *module.HTTPClient, cluster, beetleURL, token string) (error, [][]string) {
	response, err := httpClient.Get(
		context.Background(),
		fmt.Sprintf("%s/api/v1/cluster/%s", beetleURL, cluster),
		map[string]string{},
		map[string]string{"X-AUTH-TOKEN": token},
	)

	if httpClient.GetStatusCode(response) != http.StatusOK || err != nil {
		return fmt.Errorf("Unable to fetch remote data"), [][]string{}
	}

	_, err = httpClient.ToString(response)

	if err != nil {
		return fmt.Errorf("Invalid response"), [][]string{}
	}

	// @TODO
	// convert json response to struct
	// then into [][]string{}

	return nil, [][]string{
		{"staging", "down"},
	}
}
