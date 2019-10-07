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

func parseDaemonMergeParams(ctx context.Context) *cobra.Command {
	var (
		daemonHost     string
		daemonPort     int
		backupFileName string
	)

	var daemonMergeCmd = &cobra.Command{
		Use:   "daemon-merge",
		Short: "merge incremental backup using merge daemon",
		Run: func(cmd *cobra.Command, args []string) {
			log.Infof("Starting daemon-merge, host = %s, port = %d, bkp-name = %s", daemonHost, daemonPort, backupFileName)
			// daemonMerge(daemonHost, daemonPort, backupFileName) // AALEKSEEV TODO
		},
	}
	mergeFlags := daemonMergeCmd.Flags()
	bkpFileName := "bkp-name"
	mergeFlags.StringVarP(
		&backupFileName, bkpFileName, "n", "", "file name if incremental backup (required)")
	mergeFlags.StringVarP(
		&daemonHost, "address", "a", "localhost", "merge daemon listen address or host")
	mergeFlags.IntVarP(
		&daemonPort, "port", "p", 8080, "merge daemon listen port")

	err := cobra.MarkFlagRequired(mergeFlags, bkpFileName)
	if err != nil {
		err := errors.Wrap(err, "failed to set required param: "+bkpFileName)
		exitWithError(err)
	}

	return daemonMergeCmd
}
