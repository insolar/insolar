package bootstrap

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/packets"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	mock "github.com/insolar/insolar/testutils/network"
)

func TestRequester_Authorize(t *testing.T) {

	var cert certificate.Certificate

	d := cert.GetDiscoveryNodes()[0]
	h, err := host.NewHostN(d.GetHost(), *d.GetNodeRef())
	assert.NoError(t, err)

	options := common.ConfigureOptions(configuration.NewConfiguration())

	r := NewRequester(options)
	resp, err := r.Authorize(context.Background(), h, &cert)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestRequester_Bootstrap(t *testing.T) {
	options := common.ConfigureOptions(configuration.NewConfiguration())

	hn := mock.NewHostNetworkMock(t)
	hn.SendRequestToHostMock.Set(func(p context.Context, p1 types.PacketType, p2 interface{}, p3 *host.Host) (r network.Future, r1 error) {
		return nil, errors.New("123")
	})

	p := &packet.Permit{}
	claim := &packets.NodeJoinClaim{}
	r := NewRequester(options)
	// inject HostNetwork
	r.(*requester).HostNetwork = hn

	resp, err := r.Bootstrap(context.Background(), p, claim, 0)
	assert.Nil(t, resp)
	assert.Error(t, err)
}
