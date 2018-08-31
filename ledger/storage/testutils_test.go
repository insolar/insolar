package storage_test

import (
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/badgerdb/badgertestutils"
	"github.com/insolar/insolar/ledger/storage/leveldb/leveltestutils"
)

func _tmpstore(t *testing.T) (storage.Store, func()) {
	return leveltestutils.TmpDB(t, "")
}

func tmpstore(t *testing.T) (storage.Store, func()) {
	return badgertestutils.TmpDB(t, "")
}

func zerohash() []byte {
	b := make([]byte, record.HashSize)
	return b
}

func randhash() []byte {
	b := zerohash()
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func hexhash(hash string) []byte {
	b := zerohash()
	if len(hash)%2 == 1 {
		hash = "0" + hash
	}
	h, err := hex.DecodeString(hash)
	if err != nil {
		panic(err)
	}
	_ = copy(b, h)
	return b
}

func referenceWithHashes(domainhash, recordhash string) record.Reference {
	dh := hexhash(domainhash)
	rh := hexhash(recordhash)

	return record.Reference{
		Domain: record.ID{Hash: dh},
		Record: record.ID{Hash: rh},
	}
}
