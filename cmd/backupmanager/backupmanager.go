//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/log"
)

func exitWithError(err error) {
	log.Error(err.Error())
	Exit(1)
}

func exit(err error) {
	if err == nil {
		Exit(0)
	}
	exitWithError(err)
}

func stopDB(ctx context.Context, bdb *store.BadgerDB, originalError error) {
	closeError := bdb.Stop(ctx)
	if closeError != nil {
		err := errors.Wrap(originalError, "failed to stop store.BadgerDB")
		log.Error(err.Error())
	}
	if originalError != nil {
		log.Error(originalError.Error())
	}
	if originalError != nil || closeError != nil {
		Exit(1)
	}
}

func closeRawDB(bdb *badger.DB, originalError error) {
	closeError := bdb.Close()
	if closeError != nil {
		closeError = errors.Wrap(originalError, "failed to close badger.DB")
		log.Error(closeError.Error())
	}
	if originalError != nil {
		log.Error(originalError.Error())
	}
	if originalError != nil || closeError != nil {
		Exit(1)
	}
}

type BadgerLogger struct {
	insolar.Logger
}

func (b BadgerLogger) Warningf(fmt string, args ...interface{}) {
	b.Warnf(fmt, args...)
}

var (
	badgerLogger BadgerLogger
)

func merge(targetDBPath string, backupFileName string, numberOfWorkers int) {
	log.GlobalLogger().WithFields(map[string]interface{}{
		"targetDBPath":    targetDBPath,
		"backupFileName":  backupFileName,
		"numberOfWorkers": numberOfWorkers,
	}).Info("merge started")

	ops := badger.DefaultOptions(targetDBPath)
	ops.Logger = badgerLogger
	bdb, err := badger.Open(ops)
	if err != nil {
		err = errors.Wrap(err, "failed to open DB")
		exitWithError(err)
	}
	log.Info("DB is opened")

	err = isDBEmpty(bdb)
	if err == nil {
		// will exit
		err = errors.New("DB must not be empty")
		closeRawDB(bdb, err)
	}
	log.Info("DB is not empty")

	bkpFile, err := os.Open(backupFileName)
	if err != nil {
		err = errors.Wrap(err, "failed to open incremental backup file")
		closeRawDB(bdb, err)
	}
	defer bkpFile.Close()
	log.Info("Backup file is opened")

	err = bdb.Load(bkpFile, numberOfWorkers)
	if err != nil {
		// will exit
		err = errors.Wrap(err, "failed to load incremental backup file")
		closeRawDB(bdb, err)
	}

	log.Info("Successfully loaded, running GC")
	err = bdb.RunValueLogGC(0.7)
	if err == nil {
		log.Info("GC done, closing DB")
	} else {
		log.Warn("Failed to run GC: " + err.Error())
	}

	closeRawDB(bdb, nil)
	log.Info("DB closed")
}

func parseMergeParams() *cobra.Command {
	var (
		targetDBPath    string
		backupFileName  string
		numberOfWorkers int
		pprofFlag       bool
	)

	var mergeCmd = &cobra.Command{
		Use:   "merge",
		Short: "merge incremental backup to existing db",
		Run: func(cmd *cobra.Command, args []string) {
			if pprofFlag {
				f, err := os.Create(string(time.Now().Unix()) + ".pprof")
				if err != nil {
					log.Fatal("could not create CPU profile: ", err)
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					log.Fatal("could not start CPU profile: ", err)
				}
				defer pprof.StopCPUProfile()
			}

			merge(targetDBPath, backupFileName, numberOfWorkers)
		},
	}
	mergeFlags := mergeCmd.Flags()
	targetDBFlagName := "target-db"
	bkpFileName := "bkp-name"
	mergeFlags.StringVarP(
		&targetDBPath, targetDBFlagName, "t", "", "directory where backup will be roll to (required)")
	mergeFlags.StringVarP(
		&backupFileName, bkpFileName, "n", "", "file name if incremental backup (required)")
	mergeFlags.IntVarP(
		&numberOfWorkers, "workers-num", "w", 1, "number of workers to read backup file")
	mergeFlags.BoolVarP(
		&pprofFlag, "pprof", "p", false, "run merge with cpu profile")

	err := cobra.MarkFlagRequired(mergeFlags, targetDBFlagName)
	if err != nil {
		err := errors.Wrap(err, "failed to set required param: "+targetDBFlagName)
		exitWithError(err)
	}
	err = cobra.MarkFlagRequired(mergeFlags, bkpFileName)
	if err != nil {
		err := errors.Wrap(err, "failed to set required param: "+bkpFileName)
		exitWithError(err)
	}

	return mergeCmd
}

func isDBEmpty(bdb *badger.DB) error {
	tableInfo := bdb.Tables(true)
	if len(tableInfo) != 0 {
		return errors.New("tableInfo is not empty")
	}

	lsm, vlog := bdb.Size()
	if lsm != 0 || vlog != 0 {
		log.Infof("lsm: %zd, vlog: %zd", lsm, vlog)
		return errors.New("lsm or vlog are not empty")
	}

	return nil
}

func createEmptyBadger(dbDir string) {
	log.Info("createEmptyBadger. dbDir: ", dbDir)

	ops := badger.DefaultOptions(dbDir)
	ops.Logger = badgerLogger
	bdb, err := badger.Open(ops)
	if err != nil {
		err = errors.Wrap(err, "failed to open DB")
		exitWithError(err)
	}
	log.Info("DB is opened")

	err = isDBEmpty(bdb)
	if err != nil {
		// will exit
		closeRawDB(bdb, errors.Wrap(err, "DB must be empty"))
	}
	log.Info("DB is empty")

	value, err := time.Now().MarshalBinary()
	if err != nil {
		log.Panic("failed to marshal time: ", err.Error())
	}
	var key executor.DBInitializedKey
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	err = bdb.Update(func(txn *badger.Txn) error {
		return txn.Set(fullKey, value)
	})
	if err != nil {
		closeRawDB(bdb, err)
		return
	}

	log.Info("DBInitializedKey is set")

	t := time.Time{}
	err = t.UnmarshalBinary(value)
	if err != nil {
		// will exit
		closeRawDB(bdb, err)
	}
	log.Info("Set db initialized key: ", t.String())

	closeRawDB(bdb, nil)
}

func parseCreateParams() *cobra.Command {
	var dbDir string
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "create new empty DB",
		Run: func(cmd *cobra.Command, args []string) {
			createEmptyBadger(dbDir)
		},
	}

	dbDirFlagName := "db-dir"
	createCmd.Flags().StringVarP(
		&dbDir, dbDirFlagName, "d", "", "directory where new DB will be created (required)")

	err := cobra.MarkFlagRequired(createCmd.Flags(), dbDirFlagName)
	if err != nil {
		err := errors.Wrap(err, "failed to set required param: "+dbDirFlagName)
		exitWithError(err)
	}

	return createCmd
}

func writeLastBackupFile(to string, lastBackupedVersion uint64) error {
	backupInfo := executor.LastBackupInfo{
		LastBackupedVersion: lastBackupedVersion,
	}
	rawInfo, err := json.MarshalIndent(backupInfo, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to MarshalIndent")
	}

	err = ioutil.WriteFile(to, rawInfo, 0600)
	return errors.Wrap(err, "failed to write to file")
}

func finalizeLastPulse(ctx context.Context, bdb store.DB) (insolar.PulseNumber, error) {
	pulsesDB := pulse.NewDB(bdb)

	jetKeeper := executor.NewJetKeeper(jet.NewDBStore(bdb), bdb, pulsesDB)
	log.Info("Current top sync pulse: ", jetKeeper.TopSyncPulse().String())

	it := bdb.NewIterator(executor.BackupStartKey(math.MaxUint32), true)
	if !it.Next() {
		return 0, errors.New("no backup start keys")
	}

	pulseNumber := insolar.NewPulseNumber(it.Key())
	log.Info("Found last backup start key: ", pulseNumber.String())

	if pulseNumber < jetKeeper.TopSyncPulse() {
		return 0, errors.New("Found last backup start key must be grater or equal to top sync pulse")
	}

	if !jetKeeper.HasAllJetConfirms(ctx, pulseNumber) {
		return 0, errors.New("data is inconsistent. pulse " + pulseNumber.String() + " must have all confirms")
	}

	log.Info("All jets confirmed for pulse: ", pulseNumber.String())
	err := jetKeeper.AddBackupConfirmation(ctx, pulseNumber)
	if err != nil {
		return 0, errors.Wrap(err, "failed to add backup confirmation for pulse "+pulseNumber.String())
	}

	if jetKeeper.TopSyncPulse() != pulseNumber {
		return 0, errors.New("new top sync pulse must be equal to last backuped")

	}

	return jetKeeper.TopSyncPulse(), nil
}

type nopWriter struct{}

func (nopWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// prepareBackup does:
// 1. finalize last pulse, since it comes not finalized ( since we set finalization after success of backup )
// 2. gets last backuped version
// 3. write 2. to file
func prepareBackup(dbDir string, lastBackupedVersionFile string) {
	log.Info("prepareBackup. dbDir: ", dbDir, ", lastBackupedVersionFile: ", lastBackupedVersionFile)

	ops := badger.DefaultOptions(dbDir)
	ops.Logger = badgerLogger
	bdb, err := store.NewBadgerDB(ops)
	if err != nil {
		err := errors.Wrap(err, "failed to open DB")
		exitWithError(err)
	}
	log.Info("DB is opened")
	ctx := context.Background()

	topSyncPulse, err := finalizeLastPulse(ctx, bdb)
	if err != nil {
		err = errors.Wrap(err, "failed to finalizeLastPulse")
		stopDB(ctx, bdb, err)
	}

	lastVersion, err := bdb.Backup(nopWriter{}, 0)
	if err != nil {
		err = errors.Wrap(err, "failed to calculate last backuped version")
		stopDB(ctx, bdb, err)
	}
	log.Info("Got last backup version: ", lastVersion)

	if err := writeLastBackupFile(lastBackupedVersionFile, lastVersion); err != nil {
		err = errors.Wrap(err, "failed to writeLastBackupFile")
		stopDB(ctx, bdb, err)
	}
	log.Info("Write last backup version file: ", lastBackupedVersionFile)

	stopDB(ctx, bdb, nil)
	log.Info("New top sync pulse: ", topSyncPulse.String())
}

func parsePrepareBackupParams() *cobra.Command {
	var (
		dbDir                   string
		lastBackupedVersionFile string
	)
	var prepareBackupCmd = &cobra.Command{
		Use:   "prepare_backup",
		Short: "prepare backup for usage",
		Run: func(cmd *cobra.Command, args []string) {
			prepareBackup(dbDir, dbDir+"/"+lastBackupedVersionFile)
		},
	}

	dbDirFlagName := "db-dir"
	prepareBackupCmd.Flags().StringVarP(
		&dbDir, dbDirFlagName, "d", "", "directory where new DB will be created (required)")
	lastBackupFileFlagName := "last-backup-info"
	prepareBackupCmd.Flags().StringVarP(
		&lastBackupedVersionFile, lastBackupFileFlagName, "l", "", "file where last backup info will be stored (required)")

	err := cobra.MarkFlagRequired(prepareBackupCmd.Flags(), dbDirFlagName)
	if err != nil {
		err = errors.Wrap(err, "failed to set required param: "+dbDirFlagName)
		exitWithError(err)
	}

	err = cobra.MarkFlagRequired(prepareBackupCmd.Flags(), lastBackupFileFlagName)
	if err != nil {
		err = errors.Wrap(err, "failed to set required param: "+lastBackupFileFlagName)
		exitWithError(err)
	}

	return prepareBackupCmd
}

func parseInputParams() {
	var rootCmd = &cobra.Command{
		Use:   "backupmanager",
		Short: "backupmanager is the command line client for managing backups",
	}

	rootCmd.AddCommand(parseCreateParams())
	rootCmd.AddCommand(parseMergeParams())
	rootCmd.AddCommand(parseDaemonParams())
	rootCmd.AddCommand(parseDaemonMergeParams())
	rootCmd.AddCommand(parsePrepareBackupParams())

	exit(rootCmd.Execute())
}

func initLogger() context.Context {
	err := log.SetLevel("Debug")
	if err != nil {
		err = errors.Wrap(err, "failed to set log level")
		exitWithError(err)
	}

	cfg := configuration.NewLog()
	cfg.Level = "Debug"
	cfg.Formatter = "text"

	ctx, logger := inslogger.InitNodeLogger(context.Background(), cfg, "", "", "backuper")
	badgerLogger.Logger = logger.WithField("component", "badger")

	return ctx
}

func initExit(ctx context.Context) {
	InitExitContext(inslogger.FromContext(ctx))
	AtExit("logger-flusher", func() error {
		log.Flush()
		return nil
	})
}

func main() {
	ctx := initLogger()
	initExit(ctx)
	parseInputParams()
}
