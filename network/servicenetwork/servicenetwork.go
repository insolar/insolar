// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package servicenetwork

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	component "github.com/insolar/component-manager"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller"
	"github.com/insolar/insolar/network/gateway"
	"github.com/insolar/insolar/network/gateway/bootstrap"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/routing"
	"github.com/insolar/insolar/network/storage"
	"github.com/insolar/insolar/network/termination"
	"github.com/insolar/insolar/network/transport"
)

// ServiceNetwork is facade for network.
type ServiceNetwork struct {
	cfg configuration.Configuration
	cm  *component.Manager

	// dependencies
	CertificateManager  insolar.CertificateManager         `inject:""`
	PulseManager        insolar.PulseManager               `inject:""`
	CryptographyService insolar.CryptographyService        `inject:""`
	CryptographyScheme  insolar.PlatformCryptographyScheme `inject:""`
	KeyProcessor        insolar.KeyProcessor               `inject:""`
	ContractRequester   insolar.ContractRequester          `inject:""`

	// watermill support interfaces
	Pub message.Publisher `inject:""`

	// subcomponents
	RPC                controller.RPCController   `inject:"subcomponent"`
	TransportFactory   transport.Factory          `inject:"subcomponent"`
	PulseAccessor      storage.PulseAccessor      `inject:"subcomponent"`
	PulseAppender      storage.PulseAppender      `inject:"subcomponent"`
	NodeKeeper         network.NodeKeeper         `inject:"subcomponent"`
	TerminationHandler network.TerminationHandler `inject:"subcomponent"`

	HostNetwork network.HostNetwork

	Gatewayer   network.Gatewayer
	BaseGateway *gateway.Base
}

// NewServiceNetwork returns a new ServiceNetwork.
func NewServiceNetwork(conf configuration.Configuration, rootCm *component.Manager) (*ServiceNetwork, error) {
	serviceNetwork := &ServiceNetwork{cm: component.NewManager(rootCm), cfg: conf}
	return serviceNetwork, nil
}

// Init implements component.Initer
func (n *ServiceNetwork) Init(ctx context.Context) error {
	hostNetwork, err := hostnetwork.NewHostNetwork(n.CertificateManager.GetCertificate().GetNodeRef().String())
	if err != nil {
		return errors.Wrap(err, "failed to create hostnetwork")
	}
	n.HostNetwork = hostNetwork

	options := network.ConfigureOptions(n.cfg)

	cert := n.CertificateManager.GetCertificate()

	nodeNetwork, err := nodenetwork.NewNodeNetwork(n.cfg.Host.Transport, cert)
	if err != nil {
		return errors.Wrap(err, "failed to create NodeNetwork")
	}

	n.BaseGateway = &gateway.Base{Options: options}
	n.Gatewayer = gateway.NewGatewayer(n.BaseGateway.NewGateway(ctx, insolar.NoNetworkState))

	table := &routing.Table{}

	n.cm.Inject(n,
		table,
		cert,
		transport.NewFactory(n.cfg.Host.Transport),
		hostNetwork,
		nodeNetwork,
		controller.NewRPCController(options),
		controller.NewPulseController(),
		bootstrap.NewRequester(options),
		storage.NewMemoryStorage(),
		n.BaseGateway,
		n.Gatewayer,
		storage.NewMemoryStorage(),
		termination.NewHandler(n),
	)

	err = n.cm.Init(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to init internal components")
	}

	return nil
}

// Start implements component.Starter
func (n *ServiceNetwork) Start(ctx context.Context) error {
	err := n.cm.Start(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to start component manager")
	}

	bootstrapPulse := gateway.GetBootstrapPulse(ctx, n.PulseAccessor)
	n.Gatewayer.Gateway().Run(ctx, bootstrapPulse)
	n.RPC.RemoteProcedureRegister(deliverWatermillMsg, n.processIncoming)

	return nil
}

func (n *ServiceNetwork) Leave(ctx context.Context, eta insolar.PulseNumber) {
	logger := inslogger.FromContext(ctx)
	logger.Info("Gracefully stopping service network")

	// TODO: fix leave
	// n.consensusController.Leave(0)
}

func (n *ServiceNetwork) GracefulStop(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	// node leaving from network
	// all components need to do what they want over net in gracefulStop

	logger.Info("ServiceNetwork.GracefulStop wait for accepting leaving claim")
	n.TerminationHandler.Leave(ctx, 0)
	logger.Info("ServiceNetwork.GracefulStop - leaving claim accepted")

	return nil
}

// Stop implements insolar.Component
func (n *ServiceNetwork) Stop(ctx context.Context) error {
	return n.cm.Stop(ctx)
}

// HandlePulse process pulse from PulseController
func (n *ServiceNetwork) HandlePulse(ctx context.Context, pulse insolar.Pulse, originalPacket network.ReceivedPacket) {
	n.Gatewayer.Gateway().OnPulseFromPulsar(ctx, pulse, originalPacket)
}

func (n *ServiceNetwork) GetOrigin() insolar.NetworkNode {
	return n.NodeKeeper.GetOrigin()
}

func (n *ServiceNetwork) GetAccessor(p insolar.PulseNumber) network.Accessor {
	return n.NodeKeeper.GetAccessor(p)
}

func (n *ServiceNetwork) GetCert(ctx context.Context, ref *insolar.Reference) (insolar.Certificate, error) {
	return n.Gatewayer.Gateway().Auther().GetCert(ctx, ref)
}
