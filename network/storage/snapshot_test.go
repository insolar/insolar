package storage

import (
	"context"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
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

	ks := platformpolicy.NewKeyProcessor()
	p1, err := ks.GeneratePrivateKey()
	n := node.NewNode(testutils.RandomRef(), insolar.StaticRoleVirtual, ks.ExtractPublicKey(p1), "127.0.0.1:22", "ver2")

	nodes := make(map[insolar.Reference]insolar.NetworkNode)
	nodes[testutils.RandomRef()] = n

	pulse := insolar.Pulse{PulseNumber: 15}
	snap := node.NewSnapshot(pulse.PulseNumber, nodes)

	err = ss.Append(ctx, pulse.PulseNumber, snap)
	assert.NoError(t, err)

	snapshot2, err := ss.ForPulseNumber(ctx, pulse.PulseNumber)
	assert.NoError(t, err)

	assert.True(t, snap.Equal(snapshot2))

	err = cm.Stop(ctx)
}
