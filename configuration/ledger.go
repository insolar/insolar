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

package configuration

// Storage configures Ledger's storage.
type Storage struct {
	// DataDirectory is a directory where database's files live.
	DataDirectory string
	// BadgerValueLogGCDiscardRatio controls badger's value log GC behaviour.
	// Compaction on value log file happens only if data would be compacted to at least 1-BadgerValueLogGCDiscardRatio ratio.
	BadgerValueLogGCDiscardRatio float64
}

// JetSplit holds configuration for jet split.
type JetSplit struct {
	// RecordsCountThreshold is a drop threshold in records to perform split for jet.
	ThresholdRecordsCount int
	// ThresholdOverflowCount is a how many times in row ThresholdRecordsCount should be surpassed.
	ThresholdOverflowCount int
	// DepthLimit limits jet tree depth (maximum possible jets = 2^DepthLimit)
	DepthLimit uint8
}

// Ledger holds configuration for ledger.
type Ledger struct {
	// Storage defines storage configuration.
	Storage Storage
	// JetSplit holds jet split configuration.
	JetSplit JetSplit

	// common/sharable values:

	// LightChainLimit is maximum pulse difference (NOT number of pulses)
	// between current and the latest replicated on heavy.
	//
	// IMPORTANT: It should be the same on ALL nodes.
	LightChainLimit int

	// Backup holds configuration of BackupMaker
	Backup Backup

	// CleanerDelay holds value of pulses, that should happen before end of LightChainLimit and start
	// of LME's data cleaning
	CleanerDelay int

	// MaxNotificationsPerPulse holds the limit for abandoned requests notifications limit
	MaxNotificationsPerPulse uint

	// FilamentCacheLimit holds the limit for cache items for an object
	FilamentCacheLimit int
}

// Backup holds configuration for backuping.
type Backup struct {
	// Enabled switches on backuping
	Enabled bool

	// TmpDirectory is directory for tmp storage of backup data. Must be created
	TmpDirectory string

	// TargetDirectory is directory where backups will be moved to
	TargetDirectory string

	// MetaInfoFile contains meta info about current incremental backup. It will be in json format
	MetaInfoFile string

	// ConfirmFile: we wait this file being created when backup was saved on remote host
	ConfirmFile string

	// BackupFile is file with incremental backup data
	BackupFile string

	// DirNameTemplate is template for saving current incremental backup. Should be like "pulse-%d"
	DirNameTemplate string

	// BackupWaitPeriod - how much time we will wait for appearing of file ConfirmFile
	BackupWaitPeriod uint

	// LastBackupInfoFile contains info about last backup. Must exist and lie inside Storage.DataDirectory
	// Is's a full path.
	LastBackupInfoFile string

	// PostProcessBackupCmd - command which will be invoked after creating backup. It might be used to
	// send backup to remote node and do some external checks. If everything is ok, this command must create ConfirmFile
	// It will be invoked with environment variable 'INSOLAR_CURRENT_BACKUP_DIR'
	// PostProcessBackupCmd[0] is interpreted as command, and PostProcessBackupCmd[1:] as arguments
	PostProcessBackupCmd []string

	// Paths:
	// Every incremental backup live in  TargetDirectory/"DirNameTemplate"%<pulse_number>
	// and it contains:
	// incr.bkp - backup file
	// MetaInfoFile - meta info about current backup
	// ConfirmFile - must be set from outside. When it appear it means that we successfully saved the backup
}

// NewLedger creates new default Ledger configuration.
func NewLedger() Ledger {
	dataDir := "./data"
	return Ledger{
		Storage: Storage{
			DataDirectory:                dataDir,
			BadgerValueLogGCDiscardRatio: 0.4,
		},

		JetSplit: JetSplit{
			// TODO: find best default values
			ThresholdRecordsCount:  100,
			ThresholdOverflowCount: 3,
			DepthLimit:             5, // limit to 32 jets
		},
		LightChainLimit: 5, // 5 pulses

		Backup: Backup{
			Enabled:            false,
			DirNameTemplate:    "pulse-%d",
			BackupWaitPeriod:   60,
			MetaInfoFile:       "meta.json",
			BackupFile:         "incr.bkp",
			ConfirmFile:        "BACKUPED",
			LastBackupInfoFile: dataDir + "/last_backup_info.json",
		},

		CleanerDelay:             3,    // 3 pulses
		MaxNotificationsPerPulse: 100,  // 100 objects
		FilamentCacheLimit:       3000, // 3000 records for every object
	}
}
