// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var jobCmd = &cobra.Command{
	Use:   "job",
	Short: "# job",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("job")
	},
}

func init() {
	rootCmd.AddCommand(jobCmd)
}
