package artifactmanager

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/ledger/heavy"
)

func Test_ErrToReply(t *testing.T) {

	err := func() error { return heavy.ErrSyncInProgress }()
	reply := heavyerrreply(err)
	fmt.Printf("err=%+v (%T), repl=%#v (%T)\n",
		err, err, reply, reply)
}
