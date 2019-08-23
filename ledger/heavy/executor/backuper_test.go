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

package executor_test

import (
	"context"
	"encoding/binary"
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

type testKey struct {
	id uint64
}

func (t *testKey) ID() []byte {
	bs := make([]byte, 8)
	binary.PutUvarint(bs, t.id)
	return bs
}

func (t *testKey) Scope() store.Scope {
	return store.ScopeJetDrop
}

func TestBackuper_BadConfig(t *testing.T) {
	existingDir, err := os.Getwd()
	require.NoError(t, err)

	testPulse := insolar.GenesisPulse.PulseNumber

	cfg := configuration.Backup{TmpDirectory: "-----", Enabled: true}
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "check TmpDirectory returns error: stat -----: no such file or directory")

	cfg = configuration.Backup{TmpDirectory: existingDir, TargetDirectory: "+_+_+_+", Enabled: true}
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "check TargetDirectory returns error: stat +_+_+_+: no such file or directory")

	cfg.TargetDirectory = existingDir
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "ConfirmFile can't be empty")

	cfg.ConfirmFile = "Test"
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "MetaInfoFile can't be empty")

	cfg.MetaInfoFile = "Test2"
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "DirNameTemplate can't be empty")

	cfg.DirNameTemplate = "Test3"
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "BackupWaitPeriod can't be 0")

	cfg.BackupWaitPeriod = 20
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "BackupFile can't be empty")

	cfg.BackupFile = "Test"
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.NoError(t, err)

}

func makeBackuperConfig(t *testing.T, prefix string) configuration.Backup {

	cfg := configuration.Backup{
		ConfirmFile:      "BACKUPED",
		MetaInfoFile:     "META.json",
		TargetDirectory:  "/tmp/BKP/TARGET/" + prefix,
		TmpDirectory:     "/tmp/BKP/TMP",
		DirNameTemplate:  "pulse-%d",
		BackupWaitPeriod: 60,
		BackupFile:       "incr.bkp",
		Enabled:          true,
	}

	err := os.MkdirAll(cfg.TargetDirectory, 0777)
	require.NoError(t, err)
	err = os.MkdirAll(cfg.TmpDirectory, 0777)
	require.NoError(t, err)

	return cfg
}

func clearData(t *testing.T, cfg configuration.Backup) {
	err := os.RemoveAll(cfg.TargetDirectory)
	require.NoError(t, err)
}

func TestBackuper_Disabled(t *testing.T) {
	cfg := configuration.Backup{Enabled: false}
	bm, err := executor.NewBackupMaker(context.Background(), nil, cfg, 0)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), 1)
	require.Equal(t, err, executor.ErrBackupDisabled)
}

func TestBackuper_BackupWaitPeriodExpired(t *testing.T) {
	cfg := makeBackuperConfig(t, t.Name())
	defer clearData(t, cfg)

	cfg.BackupWaitPeriod = 1
	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse+1)
	require.Contains(t, err.Error(), "no backup confirmation")
}

func TestBackuper_CantMoveToTargetDir(t *testing.T) {
	cfg := makeBackuperConfig(t, t.Name())
	defer clearData(t, cfg)

	testPulse := insolar.GenesisPulse.PulseNumber

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ops := BadgerDefaultOptions(tmpdir)
	db, err := store.NewBadgerDB(ops)
	require.NoError(t, err)
	defer db.Stop(context.Background())
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, 0)
	require.NoError(t, err)
	// Create dir to fail move operation
	_, err = os.Create(filepath.Join(cfg.TargetDirectory, fmt.Sprintf(cfg.DirNameTemplate, testPulse)))
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse)
	require.Contains(t, err.Error(), "can't move")
}

func TestBackuper_Backup_OldPulse(t *testing.T) {
	cfg := makeBackuperConfig(t, t.Name())
	defer clearData(t, cfg)

	testPulse := insolar.GenesisPulse.PulseNumber
	bm, err := executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.NoError(t, err)

	err = bm.MakeBackup(context.Background(), testPulse)
	require.Equal(t, err, executor.ErrAlreadyDone)

	err = bm.MakeBackup(context.Background(), testPulse-1)
	require.Equal(t, err, executor.ErrAlreadyDone)
}
