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
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/network/storage"
	"github.com/pkg/errors"
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
	require.Contains(t, err.Error(), "checkDirectory returns error: stat -----: no such file or directory")

	cfg = configuration.Backup{TmpDirectory: existingDir, TargetDirectory: "+_+_+_+", Enabled: true}
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "checkDirectory returns error: stat +_+_+_+: no such file or directory")

	cfg.TargetDirectory = existingDir
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "BackupConfirmFile can't be empty")

	cfg.BackupConfirmFile = "Test"
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "BackupInfoFile can't be empty")

	cfg.BackupInfoFile = "Test2"
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "BackupDirNameTemplate can't be empty")

	cfg.BackupDirNameTemplate = "Test3"
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.Contains(t, err.Error(), "BackupWaitPeriod can't be 0")

	cfg.BackupWaitPeriod = 20
	_, err = executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.NoError(t, err)

}

func makeBackuperConfig(t *testing.T, prefix string) configuration.Backup {

	cfg := configuration.Backup{
		BackupConfirmFile:     "BACKUPED",
		BackupInfoFile:        "META.json",
		TargetDirectory:       "/tmp/BKP/TARGET/" + prefix,
		TmpDirectory:          "/tmp/BKP/TMP",
		BackupDirNameTemplate: "pulse-%d",
		BackupWaitPeriod:      60,
		Enabled:               true,
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

	err = bm.Do(context.Background(), 1)
	require.Equal(t, err, executor.ErrAlreadyDone)
}

func TestBackuper_BackupWaitPeriodExpired(t *testing.T) {
	cfg := makeBackuperConfig(t, t.Name())
	defer clearData(t, cfg)

	cfg.BackupWaitPeriod = 1
	testPulse := insolar.GenesisPulse.PulseNumber + 1

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	require.NoError(t, err)
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, testPulse)
	require.NoError(t, err)

	err = bm.Do(context.Background(), testPulse+1)
	require.Contains(t, err.Error(), "no backup confirmation")
}

func TestBackuper_CantMoveToTargetDir(t *testing.T) {
	cfg := makeBackuperConfig(t, t.Name())
	defer clearData(t, cfg)

	testPulse := insolar.GenesisPulse.PulseNumber

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	require.NoError(t, err)
	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, 0)
	require.NoError(t, err)
	// Create dir to fail move operation
	_, err = os.Create(filepath.Join(cfg.TargetDirectory, fmt.Sprintf(cfg.BackupDirNameTemplate, testPulse)))
	require.NoError(t, err)

	err = bm.Do(context.Background(), testPulse)
	require.Contains(t, err.Error(), "can't move")
}

func TestBackuper_Backup_OldPulse(t *testing.T) {
	cfg := makeBackuperConfig(t, t.Name())
	defer clearData(t, cfg)

	testPulse := insolar.GenesisPulse.PulseNumber
	bm, err := executor.NewBackupMaker(context.Background(), nil, cfg, testPulse)
	require.NoError(t, err)

	err = bm.Do(context.Background(), testPulse)
	require.Equal(t, err, executor.ErrAlreadyDone)

	err = bm.Do(context.Background(), testPulse-1)
	require.Equal(t, err, executor.ErrAlreadyDone)
}

func makeCurrentBkpDir(cfg configuration.Backup, pulse insolar.PulseNumber) string {
	return filepath.Join(cfg.TargetDirectory, fmt.Sprintf(cfg.BackupDirNameTemplate, pulse))
}

func calculateFileHash(t *testing.T, fileName string) string {
	f, err := os.Open(fileName)
	require.NoError(t, err)
	defer f.Close()
	hasher := md5.New()
	_, err = io.Copy(hasher, f)
	require.NoError(t, err)

	return fmt.Sprintf("%x", hasher.Sum(nil))
}

const backupFileName = "incr.bkp"

func checkBackupMetaInfo(t *testing.T, cfg configuration.Backup, numIterations int, testPulse insolar.PulseNumber) {
	for i := 0; i < numIterations+1; i++ {
		currentPulse := testPulse + insolar.PulseNumber(i)
		currentBkpDir := makeCurrentBkpDir(cfg, currentPulse)
		metaInfo := filepath.Join(currentBkpDir, cfg.BackupInfoFile)
		raw, err := ioutil.ReadFile(metaInfo)
		require.NoError(t, err)

		bi := executor.BackupInfo{}
		err = json.Unmarshal(raw, &bi)
		require.NoError(t, err)

		// check file hash
		bkpFile := filepath.Join(currentBkpDir, backupFileName)
		md5sum := calculateFileHash(t, bkpFile)
		require.Equal(t, md5sum, bi.MD5)

		// check pulse
		require.Equal(t, currentPulse, bi.Pulse)
	}
}

func TestBackuperM(t *testing.T) {
	cfg := makeBackuperConfig(t, t.Name())
	defer clearData(t, cfg)

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	require.NoError(t, err)

	bm, err := executor.NewBackupMaker(context.Background(), db, cfg, insolar.GenesisPulse.PulseNumber)
	require.NoError(t, err)

	savedKeys := make(map[store.Key]insolar.PulseNumber, 0)

	var stopWriting uint32
	sgWriteStopped := sync.WaitGroup{}
	sgWriteStopped.Add(1)

	testPulse := insolar.GenesisPulse.PulseNumber + insolar.PulseNumber(rand.Int()%20000+1)
	// writing data to db
	go func() {
		for i := 0; i < 2000000; i++ {
			if atomic.LoadUint32(&stopWriting) != 0 {
				break
			}
			key := &testKey{id: uint64(i)}
			value := testPulse + insolar.PulseNumber(i)
			err := db.Set(key, value.Bytes())
			require.NoError(t, err)
			savedKeys[key] = value
			require.NoError(t, err)
			time.Sleep(time.Duration(rand.Int()%10) * time.Millisecond)
		}
		sgWriteStopped.Done()
	}()

	wgBackup := sync.WaitGroup{}
	numIterations := 5

	wgBackup.Add(numIterations)
	// doing backups
	go func() {
		for i := 0; i < numIterations; i++ {
			err := bm.Do(context.Background(), testPulse+insolar.PulseNumber(i))
			require.NoError(t, err)
			wgBackup.Done()
			time.Sleep(time.Duration(rand.Int()%1000) * time.Millisecond)
		}
	}()

	// creating backup confirmation files
	go func() {
		for i := 0; i < numIterations+1; i++ {
			time.Sleep(2 * time.Second)

			backupConfirmFile := filepath.Join(makeCurrentBkpDir(cfg, testPulse+insolar.PulseNumber(i)), cfg.BackupConfirmFile)
			for true {
				fff, err := os.Create(backupConfirmFile)
				if err != nil && strings.Contains(err.Error(), "no such file or directory") {
					time.Sleep(time.Millisecond * 200)
					fmt.Printf("%s not created yet\n", backupConfirmFile)
					continue
				}
				require.NoError(t, err)
				require.NoError(t, fff.Close())
				break
			}
		}
	}()

	// wait for all backups done
	wgBackup.Wait()
	// stop writing to db
	atomic.StoreUint32(&stopWriting, 1)
	// wait for stopping
	sgWriteStopped.Wait()

	// final backup to collect all rest records
	err = bm.Do(context.Background(), testPulse+insolar.PulseNumber(numIterations))
	require.NoError(t, err)

	// check backup hashes
	checkBackupMetaInfo(t, cfg, numIterations, testPulse)

	// load all backups and check all records
	{
		recovTmpDir, err := ioutil.TempDir("", "bdb-test-")
		defer os.RemoveAll(recovTmpDir)
		require.NoError(t, err)
		recoveredDB, err := makeRawBadger(recovTmpDir)
		require.NoError(t, err)

		for i := 0; i < numIterations+1; i++ {
			bkpFileName := filepath.Join(
				cfg.TargetDirectory,
				fmt.Sprintf(cfg.BackupDirNameTemplate, testPulse+insolar.PulseNumber(i)),
				backupFileName,
			)
			bkpFile, err := os.Open(bkpFileName)
			require.NoError(t, err)
			err = recoveredDB.Load(bkpFile, 2)
			require.NoError(t, err)
		}

		require.NotEqual(t, 0, len(savedKeys))

		for k, v := range savedKeys {
			gotRawValue, err := getFromDB(recoveredDB, k)
			require.NoError(t, err)
			gotPulseNumber := insolar.NewPulseNumber(gotRawValue)
			require.Equal(t, v, gotPulseNumber)
		}
	}

}

func makeRawBadger(dir string) (*badger.DB, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	ops := badger.DefaultOptions(dir)
	bdb, err := badger.Open(ops)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open badger")
	}

	return bdb, nil
}

func getFromDB(db *badger.DB, key store.Key) (value []byte, err error) {
	fullKey := append(key.Scope().Bytes(), key.ID()...)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(fullKey)
		if err != nil {
			return err
		}
		value, err = item.ValueCopy(nil)
		return err
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return
}
