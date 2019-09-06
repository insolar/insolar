///
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
///

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func closeRawDB(bdb *badger.DB, err error) {
	closeError := bdb.Close()
	if closeError != nil || err != nil {
		printError("failed to close db", closeError)
		printError("", err)
		os.Exit(1)
	}
}

func merge(targetDBPath string, backupFileName string, numberOfWorkers int) {
	ops := badger.DefaultOptions(targetDBPath)
	bdb, err := badger.Open(ops)
	if err != nil {
		printError("failed to open badger", err)
		os.Exit(1)
	}

	if err := isDBEmpty(bdb); err == nil {
		closeRawDB(bdb, errors.New("db must not be empty"))
		return
	}
	log.Info("DB is not empty")

	bkpFile, err := os.Open(backupFileName)
	if err != nil {
		closeRawDB(bdb, err)
		return
	}

	err = bdb.Load(bkpFile, numberOfWorkers)
	if err != nil {
		closeRawDB(bdb, err)
		return
	}
	log.Info("Successfully merged")
	closeRawDB(bdb, nil)
}

func parseMergeParams() *cobra.Command {
	var (
		targetDBPath    string
		backupFileName  string
		numberOfWorkers int
	)

	var mergeCmd = &cobra.Command{
		Use:   "merge",
		Short: "merge incremental backup to existing db",
		Run: func(cmd *cobra.Command, args []string) {
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

	err := cobra.MarkFlagRequired(mergeFlags, targetDBFlagName)
	if err != nil {
		printError("failed to set required param: "+targetDBFlagName, err)
		os.Exit(1)
	}
	err = cobra.MarkFlagRequired(mergeFlags, bkpFileName)
	if err != nil {
		printError("failed to set required param: "+bkpFileName, err)
		os.Exit(1)
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
		println("lsm: ", lsm, ", vlog: ", vlog)
		return errors.New("lsm ot vlog are not empty")
	}

	return nil
}

func createEmptyBadger(dbDir string) {
	ops := badger.DefaultOptions(dbDir)
	var err error
	bdb, err := badger.Open(ops)
	if err != nil {
		printError("failed to open badger", err)
		os.Exit(1)
	}

	err = isDBEmpty(bdb)
	if err != nil {
		closeRawDB(bdb, errors.Wrap(err, "DB must be empty"))
		return
	}
	log.Info("DB is empty")

	value, err := time.Now().MarshalBinary()
	if err != nil {
		panic("failed to marshal time: " + err.Error())
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

	t := time.Time{}
	err = t.UnmarshalBinary(value)
	if err != nil {
		closeRawDB(bdb, err)
		return
	}
	log.Info("Set db initialized key: ", t.String())

	closeRawDB(bdb, nil)
}

func parseCreateParams() *cobra.Command {
	var dbDir string
	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "create new empty badger",
		Run: func(cmd *cobra.Command, args []string) {
			createEmptyBadger(dbDir)
		},
	}

	dbDirFlagName := "db-dir"
	createCmd.Flags().StringVarP(
		&dbDir, dbDirFlagName, "d", "", "directory where new badger will be created (required)")

	err := cobra.MarkFlagRequired(createCmd.Flags(), dbDirFlagName)
	if err != nil {
		printError("failed to set required param: "+dbDirFlagName, err)
		os.Exit(1)
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
		return 0, errors.New("data is inconsistent. pulse " + pulseNumber.String() + " must has all confirms")
	}

	log.Info("All jet confirmed for pulse: ", pulseNumber.String())
	err := jetKeeper.AddBackupConfirmation(ctx, pulseNumber)
	if err != nil {
		return 0, errors.New("failed to add backup confirmation for pulse" + pulseNumber.String())
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
	ops := badger.DefaultOptions(dbDir)
	bdb, err := store.NewBadgerDB(ops)
	if err != nil {
		printError("failed to open badger", err)
		os.Exit(1)
	}
	ctx := context.Background()
	closeDB := func(err error) {
		errStop := bdb.Stop(ctx)
		if err != nil || errStop != nil {
			printError("", err)
			printError("failed to close db", errStop)
			os.Exit(1)
		}
	}

	topSyncPulse, err := finalizeLastPulse(ctx, bdb)
	if err != nil {
		closeDB(errors.Wrap(err, "failed to finalizeLastPulse"))
		return
	}

	lastVersion, err := bdb.Backup(nopWriter{}, 0)
	if err != nil {
		closeDB(errors.Wrap(err, "failed to calculate last backuped version"))
		return
	}
	log.Info("Get last backup version: ", lastVersion)

	if err := writeLastBackupFile(lastBackupedVersionFile, lastVersion); err != nil {
		closeDB(errors.Wrap(err, "failed to writeLastBackupFile"))
		return
	}
	log.Info("Write last backup version file: ", lastBackupedVersionFile)

	closeDB(nil)
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
		&dbDir, dbDirFlagName, "d", "", "directory where new badger will be created (required)")
	lastBackupFileFlagName := "last-backup-info"
	prepareBackupCmd.Flags().StringVarP(
		&lastBackupedVersionFile, lastBackupFileFlagName, "l", "", "file where last backup info will be stored (required)")

	err := cobra.MarkFlagRequired(prepareBackupCmd.Flags(), dbDirFlagName)
	if err != nil {
		printError("failed to set required param: "+dbDirFlagName, err)
		os.Exit(1)
	}

	err = cobra.MarkFlagRequired(prepareBackupCmd.Flags(), lastBackupFileFlagName)
	if err != nil {
		printError("failed to set required param: "+lastBackupFileFlagName, err)
		os.Exit(1)
	}

	return prepareBackupCmd
}

func parseInputParams() {

	var rootCmd = &cobra.Command{
		Use:   "backupmanager",
		Short: "backupmanager is the command line client for managing backups",
	}

	rootCmd.AddCommand(parseMergeParams())
	rootCmd.AddCommand(parseCreateParams())
	rootCmd.AddCommand(parsePrepareBackupParams())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printError(message string, err error) {
	if err == nil {
		return
	}
	println(errors.Wrap(err, "ERROR "+message).Error())
}

func main() {
	err := log.SetLevel("Debug")
	if err != nil {
		printError("failed to set log level", err)
		os.Exit(1)
	}

	cfg := configuration.NewLog()
	cfg.Level = "Debug"
	cfg.Formatter = "text"
	l, _ := log.NewLog(cfg)

	log.SetGlobalLogger(l)
	parseInputParams()
}
