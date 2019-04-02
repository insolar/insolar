package network2

import (
	"context"
	"sync"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
)

type Network struct {
	cfg configuration.Configuration
	cm  *component.Manager

	// dependencies
	CertificateManager  insolar.CertificateManager  `inject:""`
	PulseManager        insolar.PulseManager        `inject:""`
	PulseStorage        insolar.PulseStorage        `inject:""`
	CryptographyService insolar.CryptographyService `inject:""`
	NetworkCoordinator  insolar.NetworkCoordinator  `inject:""`
	NodeKeeper          network.NodeKeeper          `inject:""`
	NetworkSwitcher     insolar.NetworkSwitcher     `inject:""`
	TerminationHandler  insolar.TerminationHandler  `inject:""`

	gateway Gateway
	gwmutex sync.Mutex
}

func GW(n *Network) Gateway {
	return n.gateway
}

func (n *Network) Init(ctx context.Context) error {
	return nil
}

func (n *Network) Start(ctx context.Context) error {
	return nil
}

func (n *Network) Stop(ctx context.Context) error {
	return nil
}

func (n *Network) switchGateway(g Gateway) {
	n.gwmutex.Lock()
	defer n.gwmutex.Unlock()
	n.gateway = g
}
