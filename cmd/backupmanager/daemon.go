package main

import (
	"context"

	"github.com/insolar/insolar/log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func parseDaemonParams(ctx context.Context) *cobra.Command {
	var (
		daemonHost   string
		daemonPort   int
		targetDBPath string
	)

	var daemonCmd = &cobra.Command{
		Use:   "daemon",
		Short: "run merge daemon",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("Starting merge daemon, host = %s, port = %d, target-db = %s", daemonHost, daemonPort, targetDBPath)
			// daemon(daemonHost, daemonPort, targetDBPath) // AALEKSEEV TODO
		},
	}
	mergeFlags := daemonCmd.Flags()
	targetDBFlagName := "target-db"
	mergeFlags.StringVarP(
		&targetDBPath, targetDBFlagName, "t", "", "directory where backup will be roll to (required)")
	mergeFlags.StringVarP(
		&daemonHost, "address", "a", "localhost", "listen address or host")
	mergeFlags.IntVarP(
		&daemonPort, "port", "p", 8080, "listen port")

	err := cobra.MarkFlagRequired(mergeFlags, targetDBFlagName)
	if err != nil {
		err := errors.Wrap(err, "failed to set required param: "+targetDBFlagName)
		exitWithError(err)
	}

	return daemonCmd
}
