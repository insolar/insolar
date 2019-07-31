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
	"fmt"
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

type Config struct {
	TmpDirectory          string
	TargetDirectory       string
	BackupInfoFile        string
	BackupConfirmFile     string
	BackupDirNameTemplate string
}

type BackupMaker struct {
	lock                  sync.Mutex
	lastBackupedTimestamp uint64
	lastBackupedPulse     insolar.PulseNumber
	backuper              store.Backuper
	config                Config
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

func checkConfig(config Config) error {
	if err := isPathExists(config.TmpDirectory); err != nil {
		return errors.Wrap(err, "checkDirectory returns error")
	}
	if err := isPathExists(config.TargetDirectory); err != nil {
		return errors.Wrap(err, "checkDirectory returns error")
	}
	if len(config.BackupConfirmFile) == 0 {
		return errors.New("BackupConfirmFile can't be empty")
	}
	if len(config.BackupInfoFile) == 0 {
		return errors.New("BackupInfoFile can't be empty")
	}
	if len(config.BackupDirNameTemplate) == 0 {
		return errors.New("BackupDirNameTemplate can't be empty")
	}

	return nil
}

func NewBackupMaker(backuper store.Backuper, config Config, lastBackupedPulse insolar.PulseNumber) (*BackupMaker, error) {
	if err := checkConfig(config); err != nil {
		return nil, errors.Wrap(err, "bad config")
	}

	return &BackupMaker{
		backuper:          backuper,
		config:            config,
		lastBackupedPulse: lastBackupedPulse,
	}, nil
}

func move(what string, toDirectory string) error {
	fmt.Println("move: ", what, " -> ", toDirectory)
	err := os.Rename(what, toDirectory)
	if err != nil {
		return errors.Wrapf(err, "can't move %s to %s", what, toDirectory)
	}

	return nil
}

// waitForBackup waits for file filePath appearance
func waitForBackup(ctx context.Context, filePath string, numIterations int) error {
	fmt.Println("waiting for ", filePath)
	inslogger.FromContext(ctx).Debug("waiting for ", filePath)
	for i := 0; i < numIterations; i++ {
		fmt.Println("WAITING: iteration: ", i)
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

	return errors.New("no backup confirmation for pulse")
}

func writeBackupInfoFile(hash string, pulse insolar.PulseNumber, currentBT uint64, to string) error {
	type backupInfo struct {
		MD5                   string
		Pulse                 insolar.PulseNumber
		LastBackupedTimestamp uint64
	}

	bi := backupInfo{MD5: hash, Pulse: pulse, LastBackupedTimestamp: currentBT}

	rawInfo, err := json.MarshalIndent(bi, "", "    ")
	if err != nil {
		return errors.Wrap(err, "can't marshal backup info")
	}

	err = ioutil.WriteFile(to, rawInfo, 0600)
	return errors.Wrapf(err, "can't write file %s", to)
}

// prepareBackup make incremental backup and write auxiliary file with meta info
func (b *BackupMaker) prepareBackup(ctx context.Context, dirHolder *tmpDirHolder, pulse insolar.PulseNumber) (uint64, error) {
	currentBT, err := b.backuper.Backup(dirHolder.tmpFile, b.lastBackupedTimestamp)
	if err != nil {
		return 0, errors.Wrap(err, "Backup return error")
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>: ", currentBT)

	if err := dirHolder.reopenFile(ctx); err != nil {
		return 0, errors.Wrap(err, "reopenFile return error")
	}

	hasher := md5.New()
	if _, err := io.Copy(hasher, dirHolder.tmpFile); err != nil {
		return 0, errors.Wrap(err, "io.Copy return error")
	}
	md5sum := fmt.Sprintf("%x", hasher.Sum(nil))

	metaInfoFile := filepath.Join(dirHolder.tmpDir, b.config.BackupInfoFile)
	err = writeBackupInfoFile(md5sum, pulse, currentBT, metaInfoFile)
	if err != nil {
		return 0, errors.Wrap(err, "writeBackupInfoFile return error")
	}

	return currentBT, nil
}

type tmpDirHolder struct {
	tmpDir     string
	tmpFile    *os.File
	fileClosed bool
}

func (t *tmpDirHolder) closeAll(ctx context.Context) {

	err := t.tmpFile.Close()
	if err != nil {
		inslogger.FromContext(ctx).Error("can't close backup file: ", t.tmpFile, err)
	}

	err = os.RemoveAll(t.tmpDir)
	if err != nil {
		inslogger.FromContext(ctx).Error("can't remove backup file: ", t.tmpDir, err)
	}
}

func (t *tmpDirHolder) reopenFile(ctx context.Context) error {
	if err := t.tmpFile.Close(); err != nil {
		return errors.Wrapf(err, "can't close file %s", t.tmpFile.Name())
	}

	reopenedFile, err := os.OpenFile(t.tmpFile.Name(), os.O_RDONLY, 0600)
	if err != nil {
		return errors.Wrapf(err, "can't open file %s", t.tmpFile.Name())
	}

	t.tmpFile = reopenedFile
	return nil
}

func (t *tmpDirHolder) create(where string, pulse insolar.PulseNumber) (func(context.Context), error) {
	tmpDir, err := ioutil.TempDir(where, "tmp-bkp-"+pulse.String()+"-")
	if err != nil {
		return nil, errors.Wrapf(err, "can't create tmp dir: %s", where)
	}

	file, err := ioutil.TempFile(tmpDir, pulse.String())
	if err != nil {
		return nil, errors.Wrapf(err, "can't create tmp file. dir: %s, pattern: %s", tmpDir, pulse.String())
	}

	t.tmpDir = tmpDir
	t.tmpFile = file

	return t.closeAll, nil
}

func (b *BackupMaker) doBackup(ctx context.Context, lastFinalizedPulse insolar.PulseNumber) (uint64, error) {

	dirHolder := &tmpDirHolder{}
	closer, err := dirHolder.create(b.config.TmpDirectory, lastFinalizedPulse)
	defer closer(ctx)

	currentBkpTs, err := b.prepareBackup(ctx, dirHolder, lastFinalizedPulse)
	if err != nil {
		return 0, errors.Wrap(err, "prepareBackup returns error")
	}

	currentBkpDirName := fmt.Sprintf(b.config.BackupDirNameTemplate, lastFinalizedPulse)
	currentBkpDirPath := filepath.Join(b.config.TargetDirectory, currentBkpDirName)
	err = move(dirHolder.tmpDir, currentBkpDirPath)
	if err != nil {
		return 0, errors.Wrap(err, "move returns error")
	}

	err = waitForBackup(ctx, filepath.Join(currentBkpDirPath, b.config.BackupConfirmFile), 60)
	if err != nil {
		return 0, errors.Wrap(err, "waitForBackup returns error")
	}

	return currentBkpTs, nil
}

func (b *BackupMaker) Start(ctx context.Context, lastFinalizedPulse insolar.PulseNumber) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if lastFinalizedPulse <= b.lastBackupedPulse {
		return errors.Errorf("given pulse %d must more then last backuped %d", lastFinalizedPulse, b.lastBackupedPulse)
	}

	currentBkpTs, err := b.doBackup(ctx, lastFinalizedPulse)
	if err != nil {
		return errors.Wrap(err, "doBackup return error")
	}

	b.lastBackupedPulse = lastFinalizedPulse
	b.lastBackupedTimestamp = currentBkpTs
	return nil
}
