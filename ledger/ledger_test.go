package ledger

import (
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestNewLedger(t *testing.T) {
	ledger, err := newLedger()
	if err != nil {
		fmt.Println("error:", err.Error())
	}

	s := &leveldb.DBStats{}
	ledger.db.Stats(s)
	spew.Dump(s)
}
