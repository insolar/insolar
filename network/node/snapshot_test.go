package node

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSnapshotEncodeDecode(t *testing.T) {

	ks := platformpolicy.NewKeyProcessor()
	p1, err := ks.GeneratePrivateKey()
	p2, err := ks.GeneratePrivateKey()
	assert.NoError(t, err)

	n1 := newMutableNode(testutils.RandomRef(), insolar.StaticRoleVirtual, ks.ExtractPublicKey(p1), "127.0.0.1:22", "ver2")
	n2 := newMutableNode(testutils.RandomRef(), insolar.StaticRoleHeavyMaterial, ks.ExtractPublicKey(p2), "127.0.0.1:33", "ver5")

	s := Snapshot{}
	s.pulse = 22
	s.state = insolar.CompleteNetworkState
	s.nodeList[ListLeaving] = []insolar.NetworkNode{n1, n2}
	s.nodeList[ListJoiner] = []insolar.NetworkNode{n2}

	buff, err := s.Encode()
	assert.NoError(t, err)
	assert.NotEmptyf(t, buff, "should not be empty")

	s2 := Snapshot{}
	err = s2.Decode(buff)
	assert.NoError(t, err)
	assert.True(t, s.Equal(&s2))
}
