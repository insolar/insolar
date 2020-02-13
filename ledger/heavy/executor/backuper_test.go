// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package executor_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/stretchr/testify/require"
)

func TestBackuper_BadConfig(t *testing.T) {
	existingDir, err := os.Getwd()
	require.NoError(t, err)

	testPulse := insolar.GenesisPulse.PulseNumber

	cfg := configuration.Backup{TmpDirectory: "-----", Enabled: true}

	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "check TmpDirectory returns error: stat -----: no such file or directory")

	cfg = configuration.Backup{TmpDirectory: existingDir, TargetDirectory: "+_+_+_+", Enabled: true}
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "check TargetDirectory returns error: stat +_+_+_+: no such file or directory")

	cfg.TargetDirectory = existingDir
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "ConfirmFile can't be empty")

	cfg.ConfirmFile = "Test"
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "MetaInfoFile can't be empty")

	cfg.MetaInfoFile = "Test2"
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "DirNameTemplate can't be empty")

	cfg.DirNameTemplate = "Test3"
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "BackupWaitPeriod can't be 0")

	cfg.BackupWaitPeriod = 20
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "BackupFile can't be empty")

	cfg.BackupFile = "Test"
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, nil)
	require.Contains(t, err.Error(), "PostProcessBackupCmd can't be empty")

	db := store.NewDBMock(t)
	db.GetMock.Return([]byte{}, nil)

	cfg.PostProcessBackupCmd = []string{"some command"}
	_, err = executor.NewBackupMaker(context.Background(), nil, configuration.Ledger{Backup: cfg}, testPulse, db)
	require.NoError(t, err)
}

func makeBackuperConfig(t *testing.T, prefix string, badgerDir string) (configuration.Ledger, string) {

	tmpDir, err := ioutil.TempDir("", "bdb-test-")
	require.NoError(t, err)
	cfg := configuration.Backup{
		ConfirmFile:          "BACKUPED",
		MetaInfoFile:         "META.json",
		TargetDirectory:      tmpDir + "/TARGET/" + prefix,
		TmpDirectory:         tmpDir + "/TMP",
		DirNameTemplate:      "pulse-%d",
		BackupWaitPeriod:     60,
		BackupFile:           "incr.bkp",
		Enabled:              true,
		PostProcessBackupCmd: []string{"ls"},
	}

	err = os.MkdirAll(cfg.TargetDirectory, 0777)
	require.NoError(t, err)
	err = os.MkdirAll(cfg.TmpDirectory, 0777)
	require.NoError(t, err)

	return configuration.Ledger{
		Backup: cfg,
		Storage: configuration.Storage{
			DataDirectory: badgerDir,
		},
	}, tmpDir
}

func clearData(t *testing.T, tmpDir string) {
	err := os.RemoveAll(tmpDir)
	require.NoError(t, err)
}

func TestBackuper_Disabled(t *testing.T) {
	cfg, tmpDir := makeBackuperConfig(t, t.Name(), os.TempDir())
	cfg.Backup.Enabled = false
	defer clearData(t, tmpDir)
	bm, err := executor.NewBackupMaker(context.Background(), nil, cfg, 0, nil)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), 1)
	require.Equal(t, err, executor.ErrBackupDisabled)
}

func TestBackuper_PostProcessCmdReturnError(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	cfg, tmpDir := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, tmpDir)

	cfg.Backup.BackupWaitPeriod = 1

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())

	cfg.Backup.PostProcessBackupCmd = []string{""}
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse, db)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse+1)
	require.Contains(t, err.Error(), "failed to start post process command")
}

func TestBackuper_HappyPath(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	cfg, tmpDir := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, tmpDir)

	cfg.Backup.BackupWaitPeriod = 1

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	confirmFile := filepath.Join(cfg.Backup.TargetDirectory, fmt.Sprintf(cfg.Backup.DirNameTemplate, testPulse+1), cfg.Backup.ConfirmFile)
	cfg.Backup.PostProcessBackupCmd = []string{"touch", confirmFile}
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse, db)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse+1)
	require.NoError(t, err)
}

func TestBackuper_BackupWaitPeriodExpired(t *testing.T) {
	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	cfg, tmpDir := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, tmpDir)
	cfg.Backup.BackupWaitPeriod = 1

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse, db)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse+1)
	require.Contains(t, err.Error(), "no backup confirmation")
}

func TestBackuper_Backup_OldPulse(t *testing.T) {
	cfg, tmpDir := makeBackuperConfig(t, t.Name(), os.TempDir())
	defer clearData(t, tmpDir)

	db := store.NewDBMock(t)
	db.GetMock.Return([]byte{}, nil)

	testPulse := insolar.GenesisPulse.PulseNumber
	bm, err := executor.NewBackupMaker(context.Background(), nil, cfg, testPulse, db)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse)
	require.Equal(t, err, executor.ErrAlreadyDone)

	err = bm.MakeBackup(context.Background(), testPulse-1)
	require.Equal(t, err, executor.ErrAlreadyDone)
}

func TestBackuper_TruncateHead(t *testing.T) {

	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	cfg, tmpDir := makeBackuperConfig(t, t.Name(), tmpdir)
	defer clearData(t, tmpDir)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse, db)
	require.NoError(t, err)

	numElements := 10

	for i := 0; i < numElements; i++ {
		err = db.Set(executor.BackupStartKey(testPulse+insolar.PulseNumber(i)), []byte{})
		require.NoError(t, err)
	}

	numLeftElements := numElements / 2

	err = bm.TruncateHead(context.Background(), testPulse+insolar.PulseNumber(numLeftElements))
	require.NoError(t, err)

	for i := 0; i < numLeftElements; i++ {
		_, err = db.Get(executor.BackupStartKey(testPulse + insolar.PulseNumber(i)))
		require.NoError(t, err)
	}

	for i := numElements - 1; i >= numLeftElements; i-- {
		_, err = db.Get(executor.BackupStartKey(testPulse + insolar.PulseNumber(i)))
		require.EqualError(t, err, store.ErrNotFound.Error())
	}
}
