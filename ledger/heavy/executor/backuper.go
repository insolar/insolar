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
	"fmt"
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
	// Pulse is number of backuped pulse
	Pulse insolar.PulseNumber
	// LastBackupedVersion is last backaped badger's version\timestamp
	LastBackupedVersion uint64
	// Since is badger's version\timestamp from which we started backup
	Since uint64
}

// BackupMakerDefault is component which does incremental backups by consequent invoke MakeBackup()
type BackupMakerDefault struct {
	lock              sync.RWMutex
	lastBackupedPulse insolar.PulseNumber
	backuper          store.Backuper
	config            configuration.Backup
	db                store.DB
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

func checkConfig(config configuration.Ledger) error {
	backupConfig := config.Backup
	if err := isPathExists(backupConfig.TmpDirectory); err != nil {
		return errors.Wrap(err, "check TmpDirectory returns error")
	}
	if err := isPathExists(backupConfig.TargetDirectory); err != nil {
		return errors.Wrap(err, "check TargetDirectory returns error")
	}
	if len(backupConfig.ConfirmFile) == 0 {
		return errors.New("ConfirmFile can't be empty")
	}
	if len(backupConfig.MetaInfoFile) == 0 {
		return errors.New("MetaInfoFile can't be empty")
	}
	if len(backupConfig.DirNameTemplate) == 0 {
		return errors.New("DirNameTemplate can't be empty")
	}
	if backupConfig.BackupWaitPeriod == 0 {
		return errors.New("BackupWaitPeriod can't be 0")
	}
	if len(backupConfig.BackupFile) == 0 {
		return errors.New("BackupFile can't be empty")
	}
	if len(backupConfig.PostProcessBackupCmd) == 0 {
		return errors.New("PostProcessBackupCmd can't be empty")
	}

	return nil
}

type DBInitializedKey byte

func (k DBInitializedKey) Scope() store.Scope {
	return store.ScopeDBInit
}

func (k DBInitializedKey) ID() []byte {
	return []byte{1}
}

func setDBInitialized(db store.DB) error {
	var key DBInitializedKey
	_, err := db.Get(key)
	if err != nil && err != store.ErrNotFound {
		return errors.Wrap(err, "failed to get db initialized key")
	}
	if err == store.ErrNotFound {
		value, err := time.Now().MarshalBinary()
		if err != nil {
			panic("failed to marshal time: " + err.Error())
		}
		err = db.Set(key, value)
		return errors.Wrap(err, "failed to set db initialized key")
	}

	return nil
}

func NewBackupMaker(ctx context.Context,
	backuper store.Backuper,
	config configuration.Ledger,
	lastBackupedPulse insolar.PulseNumber,
	db store.DB,
) (*BackupMakerDefault, error) {
	backupConfig := config.Backup
	if backupConfig.Enabled {
		if err := checkConfig(config); err != nil {
			return nil, errors.Wrap(err, "bad config")
		}

		if err := setDBInitialized(db); err != nil {
			return nil, errors.Wrap(err, "failed to setDBInitialized")
		}

	} else {
		inslogger.FromContext(ctx).Info("Backup is disabled")
	}

	return &BackupMakerDefault{
		backuper:          backuper,
		config:            backupConfig,
		lastBackupedPulse: lastBackupedPulse,
		db:                db,
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
	logger.Info("invokeBackupPostProcessCommand starts")
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

type BackupStartKey insolar.PulseNumber

func (k BackupStartKey) Scope() store.Scope {
	return store.ScopeBackupStart
}

func (k BackupStartKey) ID() []byte {
	return insolar.PulseNumber(k).Bytes()
}

func NewBackupStartKey(raw []byte) BackupStartKey {
	key := BackupStartKey(insolar.NewPulseNumber(raw))
	return key
}

// prepareBackup just set key with new pulse
func (b *BackupMakerDefault) prepareBackup(pulse insolar.PulseNumber) error {
	err := b.db.Set(BackupStartKey(pulse), []byte{})
	if err != nil {
		return errors.Wrap(err, "Failed to set start backup key")
	}

	return nil
}

func (b *BackupMakerDefault) doBackup(ctx context.Context, lastFinalizedPulse insolar.PulseNumber) error {

	err := b.prepareBackup(lastFinalizedPulse)
	if err != nil {
		return errors.Wrap(err, "prepareBackup returns error")
	}

	currentBkpDirName := fmt.Sprintf(b.config.DirNameTemplate, lastFinalizedPulse)
	currentBkpDirPath := filepath.Join(b.config.TargetDirectory, currentBkpDirName)

	confirmFile := filepath.Join(currentBkpDirPath, b.config.ConfirmFile)

	err = os.MkdirAll(filepath.Dir(confirmFile), 0777)
	if err != nil {
		return errors.Wrapf(err, "can't create target dir")
	}

	err = invokeBackupPostProcessCommand(ctx, b.config.PostProcessBackupCmd, currentBkpDirPath)
	if err != nil {
		return errors.Wrapf(err, "failed to invoke PostProcessBackupCmd. pulse: %d", lastFinalizedPulse)
	}

	err = waitForFile(ctx, confirmFile, b.config.BackupWaitPeriod)
	if err != nil {
		return errors.Wrapf(err, "waitForBackup returns error. pulse: %d", lastFinalizedPulse)
	}

	return nil
}

func (b *BackupMakerDefault) MakeBackup(ctx context.Context, lastFinalizedPulse insolar.PulseNumber) error {
	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"last_finalized_pulse": lastFinalizedPulse,
	})

	logger.Info("MakeBackup: before acquiring the lock")
	b.lock.Lock()
	defer b.lock.Unlock()
	logger.Info("MakeBackup: lock acquired!")

	if lastFinalizedPulse <= b.lastBackupedPulse {
		logger.Info("MakeBackup: backup already done")
		return ErrAlreadyDone
	}

	if !b.config.Enabled {
		logger.Info("MakeBackup: backup disabled")
		b.lastBackupedPulse = lastFinalizedPulse
		return ErrBackupDisabled
	}

	err := b.doBackup(ctx, lastFinalizedPulse)
	if err != nil {
		logger.Infof("MakeBackup: doBackup() returned an error %v", err)
		return errors.Wrap(err, "failed to doBackup")
	}

	b.lastBackupedPulse = lastFinalizedPulse

	logger.Infof("Done!")
	return nil
}

func (b *BackupMakerDefault) TruncateHead(ctx context.Context, from insolar.PulseNumber) error {
	b.lock.Lock()
	defer b.lock.Unlock()

	it := b.db.NewIterator(BackupStartKey(from), false)
	defer it.Close()

	var hasKeys bool
	for it.Next() {
		hasKeys = true
		key := NewBackupStartKey(it.Key())
		err := b.db.Delete(&key)
		if err != nil {
			return errors.Wrapf(err, "can't delete key: %+v", key)
		}

		inslogger.FromContext(ctx).Debugf("Erased key. Pulse number: %s", key)
	}

	if !hasKeys {
		inslogger.FromContext(ctx).Infof("No records. Nothing done. Pulse number: %s", from.String())
	}

	return nil
}
