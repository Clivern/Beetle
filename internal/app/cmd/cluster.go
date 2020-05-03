// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
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
			panic("Error! beetle url is missing (eg. $ export REMOTE_BEETLE_URL=http://127.0.0.1")
		}

		if len(args) > 0 {
			err, result = getCluster(args[0], url, token)
		} else {
			err, result = getClusters(url, token)
		}

		if err != nil {
			panic(err.Error())
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
func getClusters(beetleURL, token string) (error, [][]string) {
	return nil, [][]string{
		{"staging", "down"},
		{"production", "up"},
	}
}

// getCluster Get Cluster
func getCluster(cluster string, beetleURL, token string) (error, [][]string) {
	return nil, [][]string{
		{"staging", "down"},
	}
}
