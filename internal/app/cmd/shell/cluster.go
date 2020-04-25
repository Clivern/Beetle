// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "# cluster",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("cluster", args) },
}

var clusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "# cluster list",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("cluster list", args) },
}

var clusterGetCmd = &cobra.Command{
	Use:   "get",
	Short: "# cluster get",
	Run:   func(cmd *cobra.Command, args []string) { fmt.Println("cluster get", args) },
}

func init() {
	clusterCmd.AddCommand(clusterListCmd, clusterGetCmd)
	shellCmd.AddCommand(clusterCmd)
}
