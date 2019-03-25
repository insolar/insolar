package storage

import (
	"context"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/node"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestNewSnapshotStorage(t *testing.T) {

	tmpdir, err := ioutil.TempDir("", "bdb-test-")
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)

	ctx := context.Background()
	cm := component.NewManager(nil)
	badgerDB, err := NewBadgerDB(configuration.ServiceNetwork{CacheDirectory: tmpdir})
	ss := NewSnapshotStorage()

	cm.Register(badgerDB, ss)
	cm.Inject()

	pulse := core.Pulse{PulseNumber: 15}
	snap := node.NewSnapshot(pulse.PulseNumber, nil)

	err = ss.Append(ctx, pulse.PulseNumber, snap)
	assert.NoError(t, err)

	snapshot2, err := ss.ForPulseNumber(ctx, pulse.PulseNumber)
	assert.NoError(t, err)

	assert.Equal(t, snap, snapshot2)

	err = cm.Stop(ctx)
}
