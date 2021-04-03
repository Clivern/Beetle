// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/clivern/beetle/core/module"
	"github.com/clivern/beetle/sdk"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var (
	// Beetle API Server URL
	apiURL string

	// Beetle API Server API Key
	apiKey string

	// The Kubernetes Cluster
	cluster string

	// The Kubernetes Cluster Namespace
	namespace string
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get resources",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`You must specify the type of resource to get. Current supported resources are (apps).`)
	},
}

var appsCmd = &cobra.Command{
	Use:   "apps",
	Short: "Get a list of applications with cluster id and namespace",
	Run: func(cmd *cobra.Command, aras []string) {
		// Usage
		// $ ./beetle get apps -u "http://localhost:8080" -k "" -c "production" -n "default"

		data := [][]string{}

		client := sdk.Client{}
		client.SetHTTPClient(module.NewHTTPClient(20))
		client.SetAPIURL(apiURL)
		client.SetAPIKey(apiKey)

		apps, err := client.GetApplications(context.TODO(), cluster, namespace)

		if err != nil {
			data = append(data, []string{
				fmt.Sprintf("Error: %s", err.Error()),
				"",
				"",
				"",
			})
		} else {
			for _, app := range apps.Applications {
				version := "N/A"

				if len(app.Containers) > 0 {
					version = app.Containers[0].Version
				}

				data = append(data, []string{
					app.ID,
					app.Name,
					fmt.Sprintf("%d", len(app.Containers)),
					version,
				})
			}
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"ID", "Name", "Containers", "Version"})
		table.SetAutoWrapText(false)
		table.SetAutoFormatHeaders(true)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetHeaderLine(false)
		table.SetBorder(false)
		table.SetTablePadding("\t")
		table.SetNoWhiteSpace(true)
		table.AppendBulk(data)
		table.Render()
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	appsCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "The Kubernetes Cluster Namespace (eg. default)")
	appsCmd.MarkFlagRequired("namespace")

	appsCmd.Flags().StringVarP(&cluster, "cluster", "c", "", "The Kubernetes Cluster (eg. production)")
	appsCmd.MarkFlagRequired("cluster")

	appsCmd.Flags().StringVarP(&apiKey, "api_key", "k", "", "API Key of the Beetle API Server")
	appsCmd.MarkFlagRequired("api_key")

	appsCmd.Flags().StringVarP(&apiURL, "api_url", "u", "", "Beetle API Server URL (eg. https://example.com/)")
	appsCmd.MarkFlagRequired("api_url")

	getCmd.AddCommand(appsCmd)
}
