// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/clivern/beetle/core/model"
	"github.com/clivern/beetle/core/module"
	"github.com/clivern/beetle/sdk"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora/v3"
	"github.com/spf13/cobra"
)

var (
	// Application ID
	application string

	// Application Version
	version string

	// Deployment Strategy
	strategy string

	// Ramped Strategy MaxSurge
	maxSurge string

	// Ramped Strategy MaxUnavailable
	maxUnavailable string

	// Whether to watch the deployment
	watch bool
)

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new application version",
	Run: func(cmd *cobra.Command, aras []string) {
		// Usage
		// $ ./beetle deploy -u "http://localhost:8080" -k "" -c "production" -n "default" -a "toad" -s "recreate" -v "0.2.3" -w

		client := sdk.Client{}
		client.SetHTTPClient(module.NewHTTPClient(20))
		client.SetAPIURL(apiURL)
		client.SetAPIKey(apiKey)

		spin := spinner.New(spinner.CharSets[26], 100*time.Millisecond)
		spin.Color("green")
		spin.Start()

		job, err := client.CreateDeployment(context.TODO(), model.DeploymentRequest{
			Cluster:        cluster,
			Namespace:      namespace,
			Application:    application,
			Version:        version,
			Strategy:       strategy,
			MaxSurge:       maxSurge,
			MaxUnavailable: maxUnavailable,
		})

		if err != nil {
			fmt.Println(aurora.Red(fmt.Sprintf("Error: %s", err.Error())))
			spin.Stop()
			return
		}

		if watch {
			for {
				job, err := client.GetJob(context.TODO(), job.UUID)

				if err != nil {
					fmt.Println(aurora.Red(fmt.Sprintf("Error: %s", err.Error())))
					spin.Stop()
					return
				}

				if job.Status == model.JobFailed {
					fmt.Println(aurora.Red(fmt.Sprintf(
						"Deployment Request %s Failed!",
						job.UUID,
					)))

					spin.Stop()
					return
				}

				if job.Status == model.JobSuccess {
					fmt.Println(aurora.Green(fmt.Sprintf(
						"Deployment Request %s Succeeded!",
						job.UUID,
					)))

					spin.Stop()
					return
				}

				time.Sleep(2 * time.Second)
			}
		} else {
			fmt.Println(aurora.Green(fmt.Sprintf(
				"Deployment Request %s Submitted Successfully!",
				job.UUID,
			)))
		}

		spin.Stop()
	},
}

func init() {
	deployCmd.Flags().StringVarP(&application, "application", "a", "", "The Application ID")
	deployCmd.MarkFlagRequired("application")

	deployCmd.Flags().StringVarP(&version, "version", "v", "", "The Application Version")
	deployCmd.MarkFlagRequired("version")

	deployCmd.Flags().StringVarP(&strategy, "strategy", "s", "recreate", "The Deployment Strategy (recreate, ramped, canary or blue_green)")
	deployCmd.MarkFlagRequired("strategy")

	deployCmd.Flags().StringVarP(&maxSurge, "max_surge", "g", "50%", "Deployment Strategy MaxSurge")

	deployCmd.Flags().StringVarP(&maxUnavailable, "max_unavailable", "b", "50%", "Deployment Strategy MaxUnavailable")

	deployCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "The Kubernetes Cluster Namespace (eg. default)")
	deployCmd.MarkFlagRequired("namespace")

	deployCmd.Flags().StringVarP(&cluster, "cluster", "c", "", "The Kubernetes Cluster (eg. production)")
	deployCmd.MarkFlagRequired("cluster")

	deployCmd.Flags().StringVarP(&apiKey, "api_key", "k", "", "API Key of the Beetle API Server")
	deployCmd.MarkFlagRequired("api_key")

	deployCmd.Flags().StringVarP(&apiURL, "api_url", "u", "", "Beetle API Server URL (eg. https://example.com/)")
	deployCmd.MarkFlagRequired("api_url")

	deployCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch the deployment")

	rootCmd.AddCommand(deployCmd)
}
