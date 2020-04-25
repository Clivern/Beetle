// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Interact with running beetle server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("shell", args)
	},
}

func init() {
	rootCmd.AddCommand(shellCmd)
}
