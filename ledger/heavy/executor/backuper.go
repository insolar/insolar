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

package executor

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/pkg/errors"
)

type BackupMaker struct {
	lastBackupedTimestamp uint64
	lastBackupedPulse     insolar.PulseNumber
	backuper              store.Backuper
	tmpDirectory          string
	targetDirectory       string
	backInfoFile          string
	lock                  sync.Mutex
}

func isPathExists(dirName string) error {
	if _, err := os.Stat(dirName); err != nil {
		if os.IsNotExist(err) {
			return err
		}
		return errors.Wrapf(err, "can't check existence of directory %s ", dirName)
	}

	return nil
}

func NewBackupMaker(backuper store.Backuper, tmpDir string, targetDir string) (*BackupMaker, error) {
	if err := isPathExists(tmpDir); err != nil {
		return nil, errors.Wrap(err, "checkDirectory returns error")
	}
	if err := isPathExists(targetDir); err != nil {
		return nil, errors.Wrap(err, "checkDirectory returns error")
	}

	return &BackupMaker{
		backuper:        backuper,
		tmpDirectory:    tmpDir,
		targetDirectory: targetDir,
	}, nil
}

func move(what string, toDirectory string) error {

	err := os.Rename(what, filepath.Join(toDirectory, filepath.Base(what)))
	if err != nil {
		return errors.Wrapf(err, "can't move %s to %s", what, toDirectory)
	}

	return nil
}

func waitForBackup(ctx context.Context, filePath string, numIterations int) error {
	for i := 0; i < numIterations; i++ {
		if err := isPathExists(filePath); err != nil {
			if os.IsNotExist(err) {
				inslogger.FromContext(ctx).Debugf("backup confirmation ( %s ) still doesn't exists. Sleep second.", filePath)
				time.Sleep(time.Second)
				continue
			}
			return errors.Wrap(err, "isPathExists return error")
		}
		return nil
	}

	return nil
}

func createTmpDirectory(where string, pulse insolar.PulseNumber) (string, *os.File, error) {
	tmpDir, err := ioutil.TempDir(where, "tmp-bkp-"+pulse.String()+"-")
	if err != nil {
		return "", nil, errors.Wrapf(err, "can't create tmp dir: %s", where)
	}

	file, err := ioutil.TempFile(tmpDir, pulse.String())
	if err != nil {
		return "", nil, errors.Wrapf(err, "can't create tmp file. dir: %s, pattern: ", tmpDir, pulse.String())
	}

	return tmpDir, file, nil
}

func writeBackupInfoFile(hash []byte, pulse insolar.PulseNumber, currentBT uint64, to string) error {
	type backupInfo struct {
		Hash                  []byte
		Pulse                 insolar.PulseNumber
		LastBackupedTimestamp uint64
	}

	bi := backupInfo{Hash: hash, Pulse: pulse, LastBackupedTimestamp: currentBT}

	rawInfo, err := json.MarshalIndent(bi, "", "    ")
	if err != nil {
		return errors.Wrap(err, "can't marshal backup info")
	}

	err = ioutil.WriteFile(to, rawInfo, 0600)
	return errors.Wrapf(err, "can't write file %s", to)
}

func (b *BackupMaker) prepareBackup(file *os.File, tmpDir string, pulse insolar.PulseNumber) (uint64, error) {
	currentBT, err := b.backuper.Backup(file, b.lastBackupedTimestamp)
	if err != nil {
		return 0, errors.Wrap(err, "Backup return error")
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return 0, errors.Wrap(err, "io.Copy return error")
	}

	err = writeBackupInfoFile(hash.Sum(nil), pulse, currentBT, filepath.Join(tmpDir, b.backInfoFile))
	if err != nil {
		return 0, errors.Wrap(err, "writeBackupInfoFile return error")
	}

	return currentBT, nil
}

func (b *BackupMaker) Start(ctx context.Context, lastFinalizedPulse insolar.PulseNumber) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if lastFinalizedPulse <= b.lastBackupedPulse {
		return errors.Errorf("given pulse %d must more then last backuped %d", lastFinalizedPulse, b.lastBackupedPulse)
	}

	bkpDir, bkpFile, err := createTmpDirectory(b.tmpDirectory, lastFinalizedPulse)
	if err != nil {
		return errors.Wrap(err, "createTmpDirectory returns error")
	}
	var fileClosed bool
	defer func() {
		if !fileClosed {
			err := bkpFile.Close()
			if err != nil {
				inslogger.FromContext(ctx).Error("can't close backup file: ", bkpFile, err)
			}
		}

		err = os.RemoveAll(bkpDir)
		if err != nil {
			inslogger.FromContext(ctx).Error("can't remove backup file: ", bkpFile, err)
		}
	}()

	currentBkpTs, err := b.prepareBackup(bkpFile, bkpDir, lastFinalizedPulse)
	if err != nil {
		return errors.Wrap(err, "prepareBackup returns error")
	}

	err = bkpFile.Close()
	if err != nil {
		inslogger.FromContext(ctx).Error("can't close backup file: ", bkpFile, err)
	}
	fileClosed = true

	err = move(bkpDir, b.targetDirectory)
	if err != nil {
		return errors.Wrap(err, "move returns error")
	}

	err = waitForBackup(ctx, "", 60)
	if err != nil {
		return errors.Wrap(err, "waitForBackup returns error")
	}

	b.lastBackupedPulse = lastFinalizedPulse
	b.lastBackupedTimestamp = currentBkpTs

	return nil
}
