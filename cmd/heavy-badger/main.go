// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"os"

	"github.com/spf13/cobra"
)

type appCtx struct {
	dataDir string
}

func main() {
	arg0 := os.Args[0]

	app := appCtx{}
	var rootCmd = &cobra.Command{
		Use: arg0,
		Run: func(_ *cobra.Command, _ []string) {
			fatalf("bye!")
		},
	}
	dirFlagName := "dir"
	rootCmd.PersistentFlags().StringVarP(&app.dataDir, dirFlagName, "d", "", "badger data dir")
	if err := rootCmd.MarkPersistentFlagRequired(dirFlagName); err != nil {
		fatalf("cobra error: %v", err)
	}

	rootCmd.AddCommand(
		scopesListCommand(),
		app.fixCommand(),
		app.valueHexDumpCommand(),
		app.scanCommand(),
	)

	err := rootCmd.Execute()
	if err != nil {
		fatalf("%v execution failed: %v", arg0, err)
	}
}
