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
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/stretchr/testify/require"
)

type testKey struct{}

func (t *testKey) ID() []byte {
	return make([]byte, 10)
}

func (t *testKey) Scope() store.Scope {
	return store.ScopeJetDrop
}

func TestBackuper_BadConfig(t *testing.T) {
	existingDir, err := os.Getwd()
	require.NoError(t, err)

	testPulse := insolar.GenesisPulse.PulseNumber

	cfg := executor.Config{TmpDirectory: "-----"}
	_, err = executor.NewBackupMaker(nil, cfg, testPulse)
	require.Contains(t, err.Error(), "checkDirectory returns error: stat -----: no such file or directory")

	cfg = executor.Config{TmpDirectory: existingDir, TargetDirectory: "+_+_+_+"}
	_, err = executor.NewBackupMaker(nil, cfg, testPulse)
	require.Contains(t, err.Error(), "checkDirectory returns error: stat +_+_+_+: no such file or directory")

	cfg.TargetDirectory = existingDir
	_, err = executor.NewBackupMaker(nil, cfg, testPulse)
	require.Contains(t, err.Error(), "BackupConfirmFile can't be empty")

	cfg.BackupConfirmFile = "Test"
	_, err = executor.NewBackupMaker(nil, cfg, testPulse)
	require.Contains(t, err.Error(), "BackupInfoFile can't be empty")

	cfg.BackupInfoFile = "Test2"
	_, err = executor.NewBackupMaker(nil, cfg, testPulse)
	require.Contains(t, err.Error(), "BackupDirNameTemplate can't be empty")

	cfg.BackupDirNameTemplate = "Test3"
	_, err = executor.NewBackupMaker(nil, cfg, testPulse)
	require.NoError(t, err)

}

func makeBackuperConfig() executor.Config {
	return executor.Config{
		BackupConfirmFile:     "BACKUPED",
		BackupInfoFile:        "META.json",
		TargetDirectory:       "/tmp/BKP/TARGET",
		TmpDirectory:          "/tmp/BKP/TMP",
		BackupDirNameTemplate: "pulse-%d",
	}
}

func TestBackuper_Backup_OldPulse(t *testing.T) {
	cfg := makeBackuperConfig()
	testPulse := insolar.GenesisPulse.PulseNumber
	bm, err := executor.NewBackupMaker(nil, cfg, testPulse)
	require.NoError(t, err)

	err = bm.Start(context.Background(), testPulse)
	require.Contains(t, err.Error(), "given pulse 65537 must more then last backuped 65537")

	err = bm.Start(context.Background(), testPulse-1)
	require.Contains(t, err.Error(), "given pulse 65536 must more then last backuped 65537")
}

func TestBackuper(t *testing.T) {
	cfg := makeBackuperConfig()

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	db, err := store.NewBadgerDB(tmpdir)
	require.NoError(t, err)

	bm, err := executor.NewBackupMaker(db, cfg, insolar.GenesisPulse.PulseNumber)
	require.NoError(t, err)

	for i := 0; i < 200; i++ {
		err = db.Set(&testKey{}, make([]byte, i))
		require.NoError(t, err)
	}

	testPulse := insolar.GenesisPulse.PulseNumber + insolar.PulseNumber(rand.Int()%20000+1)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		err = bm.Start(context.Background(), testPulse)
		require.NoError(t, err)
		wg.Done()
	}()

	time.Sleep(2 * time.Second)
	currentBkpDirPath := filepath.Join(cfg.TargetDirectory, fmt.Sprintf(cfg.BackupDirNameTemplate, testPulse), cfg.BackupConfirmFile)

	for true {
		fff, err := os.Create(currentBkpDirPath)
		if err != nil && strings.Contains(err.Error(), "no such file or directory") {
			time.Sleep(time.Millisecond * 200)
			fmt.Printf("%s not created yet\n", currentBkpDirPath)
			continue
		}
		require.NoError(t, err)
		require.NoError(t, fff.Close())
		break

	}

	wg.Wait()
}
