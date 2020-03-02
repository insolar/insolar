// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"github.com/spf13/cobra"
)

func (app *appCtx) fixCommand() *cobra.Command {
	var fixCmd = &cobra.Command{
		Use:   "fix",
		Short: "opens and closes badger database. Could fix 'Database was not properly closed' error.",
		Run: func(_ *cobra.Command, _ []string) {
			_, close := openDB(app.dataDir)
			close()
		},
	}
	return fixCmd
}
