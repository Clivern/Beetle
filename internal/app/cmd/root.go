// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

var rootCmd = &cobra.Command{
	Use: "beetle",
	Short: `Work seamlessly with Beetle from the command line.

Beetle is in early stages of development, and we'd love to hear your
feedback at <https://github.com/Clivern/Beetle>`,
}

// Execute runs cmd tool
func Execute(version, commit, date, builtBy string) {
	// Move build ldflags to cmd pkg
	version = version
	commit = commit
	date = date
	builtBy = builtBy

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
