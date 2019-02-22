package rules

import (
	"testing"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/testutils/network"
	"github.com/stretchr/testify/assert"
)

func TestRules_CheckMajorityRule(t *testing.T) {
	cm := component.Manager{}
	r := NewRules()
	certManager := testutils.NewCertificateManagerMock(t)
	cert := testutils.NewCertificateMock(t)

	certManager.GetCertificateMock.Set(func() (r core.Certificate) {
		return cert
	})

	cert.GetDiscoveryNodesMock.Set(func() (r []core.DiscoveryNode) {
		_, nodes := getDiscoveryNodes(5)
		return nodes
	})

	cert.GetMajorityRuleMock.Set(func() (r int) {
		return 4
	})

	nodeKeeper := network.NewNodeKeeperMock(t)
	nodeKeeper.GetActiveNodesMock.Set(func() (r []core.Node) {
		nodes, _ := getDiscoveryNodes(5)
		nodes = append(nodes, newNode(250))
		return nodes
	})

	cm.Inject(r, certManager, nodeKeeper)

	result, count := r.CheckMajorityRule()
	assert.True(t, result)
	assert.Equal(t, 5, count)
}

func getDiscoveryNodes(count int) ([]core.Node, []core.DiscoveryNode) {
	result1 := make([]core.Node, count)
	result2 := make([]core.DiscoveryNode, count)

	for i := 0; i < count; i++ {
		n := newNode(i)
		d := certificate.NewBootstrapNode(nil, "", "127.0.0.1:3000", n.ID().String())
		result1[i] = n
		result2[i] = d
	}

	return result1, result2
}

func newNode(id int) core.Node {
	recordRef := core.RecordRef{byte(id)}
	node := nodenetwork.NewNode(recordRef, core.StaticRoleVirtual, nil, "127.0.0.1:3000", "")
	return node
}
