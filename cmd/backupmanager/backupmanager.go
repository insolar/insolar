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

type BadgerLogger struct {
	insolar.Logger
}

func (b BadgerLogger) Warningf(fmt string, args ...interface{}) {
	b.Warnf(fmt, args...)
}

func exit(ctx context.Context, err error) {
	if err != nil {
		inslogger.FromContext(ctx).Error(err.Error())
		Exit(1)
	}
	Exit(0)
}

func dbIsEmpty(bdb *badger.DB) bool {
	badgerTables := bdb.Tables(false)
	if len(badgerTables) != 0 {
		return false
	}

	lsmSize, vlogSize := bdb.Size()
	return lsmSize == 0 && vlogSize == 0
}

func dbClose(ctx context.Context, bdb *badger.DB) {
	err := bdb.Close()
	if err != nil {
		inslogger.FromContext(ctx).Errorf("Failed to close database: %s", err.Error())
	}
}

func dbStop(ctx context.Context, bdb *store.BadgerDB) {
	err := bdb.Stop(ctx)
	if err != nil {
		inslogger.FromContext(ctx).Errorf("Failed to stop database: %s", err.Error())
	}
}

func defaultBadgerOpts(ctx context.Context, dbPath string) badger.Options {
	_, logger := inslogger.WithField(ctx, "component", "badger")
	badgerLogger := BadgerLogger{Logger: logger}

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = badgerLogger
	return opts
}

func merge(ctx context.Context, targetDBPath string, backupFileName string, numberOfWorkers int) error {
	logger := inslogger.FromContext(ctx)
	logger.WithFields(map[string]interface{}{
		"targetDBPath":    targetDBPath,
		"backupFileName":  backupFileName,
		"numberOfWorkers": numberOfWorkers,
	}).Info("merge started")

	bdb, err := badger.Open(defaultBadgerOpts(ctx, targetDBPath))
	if err != nil {
		return errors.Wrap(err, "failed to open database")
	}
	logger.Info("database is opened")
	defer dbClose(ctx, bdb)

	if dbIsEmpty(bdb) {
		return errors.New("database must not be empty")
	}
	logger.Info("database is not empty")

	bkpFile, err := os.Open(backupFileName)
	if err != nil {
		return errors.Wrap(err, "failed to open backup file")
	}
	logger.Info("Backup file is opened")

	err = bdb.Load(bkpFile, numberOfWorkers)
	if err != nil {
		return errors.Wrap(err, "failed to load backup file")
	}
	logger.Info("Successfully merged")

	return nil
}

func parseMergeParams(ctx context.Context) *cobra.Command {
	var (
		targetDBPath    string
		backupFileName  string
		numberOfWorkers int
	)

	var mergeCmd = &cobra.Command{
		Use:   "merge",
		Short: "merge incremental backup to existing db",
		Run: func(cmd *cobra.Command, args []string) {
			exit(ctx, merge(ctx, targetDBPath, backupFileName, numberOfWorkers))
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

	err := cobra.MarkFlagRequired(mergeFlags, targetDBFlagName)
	if err != nil {
		exit(ctx, errors.Wrap(err, "failed to set required param: "+targetDBFlagName))
	}
	err = cobra.MarkFlagRequired(mergeFlags, bkpFileName)
	if err != nil {
		exit(ctx, errors.Wrap(err, "failed to set required param: "+bkpFileName))
	}

	return mergeCmd
}

func createEmptyBadger(ctx context.Context, dbDir string) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("createEmptyBadger. dbDir: ", dbDir)

	bdb, err := badger.Open(defaultBadgerOpts(ctx, dbDir))
	if err != nil {
		return errors.Wrap(err, "failed to open database")
	}
	logger.Info("database is opened")
	defer dbClose(ctx, bdb)

	if !dbIsEmpty(bdb) {
		return errors.New("database must be empty")
	}
	logger.Info("database is empty")

	timeValue := time.Now()
	timeMarshalled, err := timeValue.MarshalBinary()
	if err != nil {
		return errors.New("failed to marshal time: " + err.Error())
	}
	var key executor.DBInitializedKey
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	err = bdb.Update(func(txn *badger.Txn) error {
		return txn.Set(fullKey, timeMarshalled)
	})
	if err != nil {
		return errors.New("failed to set DBInitializedKey")
	}
	logger.Info("DBInitializedKey is set: ", timeValue.String())

	return nil
}

func parseCreateParams(ctx context.Context) *cobra.Command {
	var dbDir string
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "create new empty DB",
		Run: func(cmd *cobra.Command, args []string) {
			exit(ctx, createEmptyBadger(ctx, dbDir))
		},
	}

	dbDirFlagName := "db-dir"
	createCmd.Flags().StringVarP(
		&dbDir, dbDirFlagName, "d", "", "directory where new DB will be created (required)")

	err := cobra.MarkFlagRequired(createCmd.Flags(), dbDirFlagName)
	if err != nil {
		exit(ctx, errors.Wrap(err, "failed to set required param: "+dbDirFlagName))
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
	logger := inslogger.FromContext(ctx)
	pulsesDB := pulse.NewDB(bdb)

	jetKeeper := executor.NewJetKeeper(jet.NewDBStore(bdb), bdb, pulsesDB)
	logger.Info("Current top sync pulse: ", jetKeeper.TopSyncPulse().String())

	it := bdb.NewIterator(executor.BackupStartKey(math.MaxUint32), true)
	if !it.Next() {
		return 0, errors.New("no backup start keys")
	}

	pulseNumber := insolar.NewPulseNumber(it.Key())
	logger.Info("Found last backup start key: ", pulseNumber.String())

	if pulseNumber < jetKeeper.TopSyncPulse() {
		return 0, errors.New("Found last backup start key must be grater or equal to top sync pulse")
	}

	if !jetKeeper.HasAllJetConfirms(ctx, pulseNumber) {
		return 0, errors.New("data is inconsistent. pulse " + pulseNumber.String() + " must have all confirms")
	}

	logger.Info("All jets confirmed for pulse: ", pulseNumber.String())
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
func prepareBackup(ctx context.Context, dbDir string, lastBackupedVersionFile string) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("prepareBackup. dbDir: ", dbDir, ", lastBackupedVersionFile: ", lastBackupedVersionFile)

	bdb, err := store.NewBadgerDB(defaultBadgerOpts(ctx, dbDir))
	if err != nil {
		return errors.Wrap(err, "failed to open database")
	}
	logger.Info("database is opened")
	defer dbStop(ctx, bdb)

	topSyncPulse, err := finalizeLastPulse(ctx, bdb)
	if err != nil {
		return errors.Wrap(err, "failed to finalizeLastPulse")
	}

	lastVersion, err := bdb.Backup(nopWriter{}, 0)
	if err != nil {
		return errors.Wrap(err, "failed to calculate last backuped version")
	}
	logger.Info("Got last backup version: ", lastVersion)

	if err := writeLastBackupFile(lastBackupedVersionFile, lastVersion); err != nil {
		return errors.Wrap(err, "failed to writeLastBackupFile")
	}
	logger.Info("Write last backup version file: ", lastBackupedVersionFile)
	logger.Info("New top sync pulse: ", topSyncPulse.String())

	return nil
}

func parsePrepareBackupParams(ctx context.Context) *cobra.Command {
	var (
		dbDir                   string
		lastBackupedVersionFile string
	)
	var prepareBackupCmd = &cobra.Command{
		Use:   "prepare_backup",
		Short: "prepare backup for usage",
		Run: func(cmd *cobra.Command, args []string) {
			exit(ctx, prepareBackup(ctx, dbDir, dbDir+"/"+lastBackupedVersionFile))
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
		exit(ctx, errors.Wrap(err, "failed to set required param: "+dbDirFlagName))
	}

	err = cobra.MarkFlagRequired(prepareBackupCmd.Flags(), lastBackupFileFlagName)
	if err != nil {
		exit(ctx, errors.Wrap(err, "failed to set required param: "+lastBackupFileFlagName))
	}

	return prepareBackupCmd
}

func parseInputParams(ctx context.Context) error {
	var rootCmd = &cobra.Command{
		Use:   "backupmanager",
		Short: "backupmanager is the command line client for managing backups",
	}

	rootCmd.AddCommand(parseMergeParams(ctx))
	rootCmd.AddCommand(parseCreateParams(ctx))
	rootCmd.AddCommand(parsePrepareBackupParams(ctx))

	return rootCmd.Execute()
}

func initLogger() context.Context {
	err := log.SetLevel("Debug")
	if err != nil {
		exit(context.Background(), errors.Wrap(err, "failed to set log level"))
	}

	cfg := configuration.NewLog()
	cfg.Level = "Debug"
	cfg.Formatter = "text"

	ctx, _ := inslogger.InitNodeLogger(context.Background(), cfg, "", "", "backuper")
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
	exit(ctx, parseInputParams(ctx))
}
