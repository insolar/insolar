/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package leveldb

import (
	"path/filepath"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/insolar/insolar/ledger/record"
)

const (
	dbDirPath = "_db"
)

// LevelLedger represents ledger's LevelDB storage.
type LevelLedger struct {
	// LDB contains LevelDB database instance.
	LDB *leveldb.DB
}

// InitDB returns LevelLedger with LevelDB initialized with default settings.
func InitDB() (*LevelLedger, error) {
	// Options struct doc: https://godoc.org/github.com/syndtr/goleveldb/leveldb/opt#Options.
	opts := &opt.Options{
		AltFilters:  nil,
		BlockCacher: opt.LRUCacher,
		// BlockCacheCapacity increased to 32MiB from default 8 MiB.
		// BlockCacheCapacity defines the capacity of the 'sorted table' block caching.
		BlockCacheCapacity:                    32 * 1024 * 1024,
		BlockRestartInterval:                  16,
		BlockSize:                             4 * 1024,
		CompactionExpandLimitFactor:           25,
		CompactionGPOverlapsFactor:            10,
		CompactionL0Trigger:                   4,
		CompactionSourceLimitFactor:           1,
		CompactionTableSize:                   2 * 1024 * 1024,
		CompactionTableSizeMultiplier:         1.0,
		CompactionTableSizeMultiplierPerLevel: nil,
		// CompactionTotalSize increased to 32MiB from default 10 MiB.
		// CompactionTotalSize limits total size of 'sorted table' for each level.
		// The limits for each level will be calculated as:
		//   CompactionTotalSize * (CompactionTotalSizeMultiplier ^ Level)
		CompactionTotalSize:                   32 * 1024 * 1024,
		CompactionTotalSizeMultiplier:         10.0,
		CompactionTotalSizeMultiplierPerLevel: nil,
		Comparer:                     comparer.DefaultComparer,
		Compression:                  opt.DefaultCompression,
		DisableBufferPool:            false,
		DisableBlockCache:            false,
		DisableCompactionBackoff:     false,
		DisableLargeBatchTransaction: false,
		ErrorIfExist:                 false,
		ErrorIfMissing:               false,
		Filter:                       nil,
		IteratorSamplingRate:         1 * 1024 * 1024,
		NoSync:                       false,
		NoWriteMerge:                 false,
		OpenFilesCacher:              opt.LRUCacher,
		OpenFilesCacheCapacity:       500,
		ReadOnly:                     false,
		Strict:                       opt.DefaultStrict,
		WriteBuffer:                  16 * 1024 * 1024, // Default is 4 MiB
		WriteL0PauseTrigger:          12,
		WriteL0SlowdownTrigger:       8,
	}

	absPath, err := filepath.Abs(dbDirPath)
	if err != nil {
		return nil, err
	}
	db, err := leveldb.OpenFile(absPath, opts)
	if err != nil {
		return nil, err
	}

	return &LevelLedger{
		LDB: db,
	}, nil
}

// Get returns record from leveldb by timeslot and hash passed in record.Key
func (ll *LevelLedger) Get(k record.Key) (rec record.Record, found bool) {
	return nil, false
}

// Set stores record in leveldb
func (ll *LevelLedger) Set(record record.Record) error {
	return nil
}
