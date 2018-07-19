package ledger

import (
	"path/filepath"

	"github.com/insolar/insolar/ledger/record"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	dbDirPath = "_db"
)

type levelLedger struct {
	db *leveldb.DB
}

func newLedger() (*levelLedger, error) {
	opts := &opt.Options{
		AltFilters:                            nil,
		BlockCacher:                           opt.LRUCacher,
		BlockCacheCapacity:                    32 * 1024 * 1024, // Default is 8 MiB
		BlockRestartInterval:                  16,
		BlockSize:                             4 * 1024,
		CompactionExpandLimitFactor:           25,
		CompactionGPOverlapsFactor:            10,
		CompactionL0Trigger:                   4,
		CompactionSourceLimitFactor:           1,
		CompactionTableSize:                   2 * 1024 * 1024,
		CompactionTableSizeMultiplier:         1.0,
		CompactionTableSizeMultiplierPerLevel: nil,
		CompactionTotalSize:                   32 * 1024 * 1024, // Default is 10 MiB
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

	return &levelLedger{
		db: db,
	}, nil
}

func (ll *levelLedger) Get(id record.RecordHash) (bool, record.Record) {
	return false, nil
}

func (ll *levelLedger) Set(record record.Record) error {
	return nil
}

func (ll *levelLedger) Update(id record.RecordHash, record record.Record) error {
	return nil
}
