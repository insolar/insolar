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

package executor

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/heavy/executor.BackupMaker -o ./ -s _gen_mock.go -g

// BackupMaker is interface for doing backups
type BackupMaker interface {
	// MakeBackup starts process of incremental backups
	MakeBackup(ctx context.Context, lastFinalizedPulse insolar.PulseNumber) error
}

var (
	// ErrAlreadyDone is returned when you try to do backup for pulse less then lastBackupedPulse
	ErrAlreadyDone = errors.New("backup already done for this pulse")
	// ErrBackupDisabled is returned when backups are disabled
	ErrBackupDisabled = errors.New("backup disabled")
)

// BackupInfo contains meta information about current incremental backup
type BackupInfo struct {
	// SHA256 is hash of backup file
	SHA256 string
	// Pulse is number of backuped pulse
	Pulse insolar.PulseNumber
	// LastBackupedVersion is last backaped badger's version\timestamp
	LastBackupedVersion uint64
	// Since is badger's version\timestamp from which we started backup
	Since uint64
}

// BackupMakerDefault is component which does incremental backups by consequent invoke MakeBackup()
type BackupMakerDefault struct {
	lock                sync.RWMutex
	lastBackupedVersion uint64
	lastBackupedPulse   insolar.PulseNumber
	backuper            store.Backuper
	config              configuration.Backup
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

func checkConfig(config configuration.Backup) error {
	if err := isPathExists(config.TmpDirectory); err != nil {
		return errors.Wrap(err, "check TmpDirectory returns error")
	}
	if err := isPathExists(config.TargetDirectory); err != nil {
		return errors.Wrap(err, "check TargetDirectory returns error")
	}
	if len(config.ConfirmFile) == 0 {
		return errors.New("ConfirmFile can't be empty")
	}
	if len(config.MetaInfoFile) == 0 {
		return errors.New("MetaInfoFile can't be empty")
	}
	if len(config.DirNameTemplate) == 0 {
		return errors.New("DirNameTemplate can't be empty")
	}
	if config.BackupWaitPeriod == 0 {
		return errors.New("BackupWaitPeriod can't be 0")
	}
	if len(config.BackupFile) == 0 {
		return errors.New("BackupFile can't be empty")
	}
	if len(config.PostProcessBackupCmd) == 0 {
		return errors.New("PostProcessBackupCmd can't be empty")
	}

	return nil
}

func NewBackupMaker(ctx context.Context, backuper store.Backuper, config configuration.Backup, lastBackupedPulse insolar.PulseNumber) (*BackupMakerDefault, error) {
	if config.Enabled {
		if err := checkConfig(config); err != nil {
			return nil, errors.Wrap(err, "bad config")
		}
	} else {
		inslogger.FromContext(ctx).Info("Backup is disabled")
	}

	return &BackupMakerDefault{
		backuper:          backuper,
		config:            config,
		lastBackupedPulse: lastBackupedPulse,
	}, nil
}

func move(ctx context.Context, what string, toDirectory string) error {
	inslogger.FromContext(ctx).Debugf("backuper. move %s -> %s", what, toDirectory)
	err := os.Rename(what, toDirectory)

	return errors.Wrapf(err, "can't move %s to %s", what, toDirectory)
}

// waitForFile waits for file filePath appearance
func waitForFile(ctx context.Context, filePath string, numIterations uint) error {
	inslogger.FromContext(ctx).Debug("waiting for ", filePath)
	for i := uint(0); i < numIterations; i++ {
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

func writeBackupInfoFile(hash string, pulse insolar.PulseNumber, since uint64, upto uint64, to string) error {
	bi := BackupInfo{
		SHA256:              hash,
		Pulse:               pulse,
		LastBackupedVersion: upto,
		Since:               since,
	}

	rawInfo, err := json.MarshalIndent(bi, "", "    ")
	if err != nil {
		return errors.Wrap(err, "can't marshal backup info")
	}

	err = ioutil.WriteFile(to, rawInfo, 0600)
	return errors.Wrapf(err, "can't write file %s", to)
}

func calculateFileHash(f *os.File) (string, error) {
	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", errors.Wrap(err, "io.Copy return error")
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

type logWrapper struct {
	logger insolar.Logger
	isInfo bool
}

func (lw *logWrapper) Write(p []byte) (n int, err error) {
	if lw.isInfo {
		lw.logger.Info(string(p))
	} else {
		lw.logger.Error(string(p))
	}
	return len(p), nil
}

func invokeBackupPostProcessCommand(ctx context.Context, command []string, currentBkpDirPath string) error {
	logger := inslogger.FromContext(ctx)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "INSOLAR_CURRENT_BACKUP_DIR="+currentBkpDirPath)
	cmd.Stdout = &logWrapper{logger: logger, isInfo: true}
	cmd.Stderr = &logWrapper{logger: logger, isInfo: false}

	err := cmd.Start()
	if err != nil {
		return errors.Wrap(err, "failed to start post process command")
	}

	err = cmd.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait for completion of post process command")
	}

	return nil
}

// prepareBackup make incremental backup and write auxiliary file with meta info
func (b *BackupMakerDefault) prepareBackup(dirHolder *tmpDirHolder, pulse insolar.PulseNumber) (currentBackupedVersion uint64, err error) {
	currentBT, err := b.backuper.Backup(dirHolder.tmpFile, b.lastBackupedVersion)
	if err != nil {
		return 0, errors.Wrap(err, "Backup return error")
	}

	if err := dirHolder.reopenTmpFile(); err != nil {
		return 0, errors.Wrap(err, "reopenFile return error")
	}

	fileHash, err := calculateFileHash(dirHolder.tmpFile)
	if err != nil {
		return 0, errors.Wrap(err, "calculateFileHash return error")
	}

	metaInfoFile := filepath.Join(dirHolder.tmpDir, b.config.MetaInfoFile)
	err = writeBackupInfoFile(fileHash, pulse, b.lastBackupedVersion, currentBT, metaInfoFile)
	if err != nil {
		return 0, errors.Wrap(err, "writeBackupInfoFile return error")
	}

	return currentBT, nil
}

type tmpDirHolder struct {
	tmpDir  string
	tmpFile *os.File
}

func (t *tmpDirHolder) release(ctx context.Context) {
	err := t.tmpFile.Close()
	if err != nil {
		inslogger.FromContext(ctx).Fatal("can't close backup file: ", t.tmpFile, err)
	}

	err = os.RemoveAll(t.tmpDir)
	if err != nil {
		inslogger.FromContext(ctx).Fatal("can't remove backup file: ", t.tmpDir, err)
	}
}

func (t *tmpDirHolder) reopenTmpFile() error {
	if err := t.tmpFile.Close(); err != nil {
		return errors.Wrapf(err, "can't close file %s", t.tmpFile.Name())
	}

	reopenedFile, err := os.OpenFile(t.tmpFile.Name(), os.O_RDONLY, 0)
	if err != nil {
		return errors.Wrapf(err, "can't open file %s", t.tmpFile.Name())
	}

	t.tmpFile = reopenedFile
	return nil
}

func (t *tmpDirHolder) create(where string, pulse insolar.PulseNumber) error {
	tmpDir, err := ioutil.TempDir(where, "tmp-bkp-"+pulse.String()+"-")
	if err != nil {
		return errors.Wrapf(err, "can't create tmp dir: %s", where)
	}

	file, err := os.OpenFile(tmpDir+"/incr.bkp", os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0600)
	if err != nil {
		return errors.Wrapf(err, "can't create tmp file. dir: %s", tmpDir)
	}

	t.tmpDir = tmpDir
	t.tmpFile = file

	return nil
}

func (b *BackupMakerDefault) doBackup(ctx context.Context, lastFinalizedPulse insolar.PulseNumber) (uint64, error) {

	dirHolder := &tmpDirHolder{}
	err := dirHolder.create(b.config.TmpDirectory, lastFinalizedPulse)
	if err != nil {
		return 0, errors.Wrap(err, "can't create tmp dir")
	}
	defer dirHolder.release(ctx)

	currentBkpVersion, err := b.prepareBackup(dirHolder, lastFinalizedPulse)
	if err != nil {
		return 0, errors.Wrap(err, "prepareBackup returns error")
	}

	currentBkpDirName := fmt.Sprintf(b.config.DirNameTemplate, lastFinalizedPulse)
	currentBkpDirPath := filepath.Join(b.config.TargetDirectory, currentBkpDirName)
	err = move(ctx, dirHolder.tmpDir, currentBkpDirPath)
	if err != nil {
		return 0, errors.Wrap(err, "move returns error")
	}

	err = invokeBackupPostProcessCommand(ctx, b.config.PostProcessBackupCmd, currentBkpDirPath)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to invoke PostProcessBackupCmd. pulse: %d", lastFinalizedPulse)
	}

	err = waitForFile(ctx, filepath.Join(currentBkpDirPath, b.config.ConfirmFile), b.config.BackupWaitPeriod)
	if err != nil {
		return 0, errors.Wrapf(err, "waitForBackup returns error. pulse: %d", lastFinalizedPulse)
	}

	return currentBkpVersion, nil
}

func (b *BackupMakerDefault) MakeBackup(ctx context.Context, lastFinalizedPulse insolar.PulseNumber) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if lastFinalizedPulse <= b.lastBackupedPulse {
		return ErrAlreadyDone
	}

	if !b.config.Enabled {
		inslogger.FromContext(ctx).Info("Trying to do backup, but it's disabled. Do nothing")
		b.lastBackupedPulse = lastFinalizedPulse
		return ErrBackupDisabled
	}

	currentBkpVersion, err := b.doBackup(ctx, lastFinalizedPulse)
	if err != nil {
		return errors.Wrap(err, "doBackup return error")
	}

	b.lastBackupedPulse = lastFinalizedPulse
	b.lastBackupedVersion = currentBkpVersion
	inslogger.FromContext(ctx).Infof("Pulse %d successfully backuped", lastFinalizedPulse)
	return nil
}
