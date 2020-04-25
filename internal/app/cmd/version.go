// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/clivern/beetle/internal/app/module"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get current and latest version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(
			fmt.Sprintf(
				`Current Beetle Version %v Commit %v, Built @%v`,
				version,
				commit,
				date,
			),
		)

		latest, err := module.GetLatestRelease()

		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			return
		}

		fmt.Printf(
			"Latest Release [%s], Latest Tag [%s]",
			latest.Name,
			latest.TagName,
		)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
