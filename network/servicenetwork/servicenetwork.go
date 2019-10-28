//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package servicenetwork

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/api"
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
	CryptographyService insolar.CryptographyService        `inject:""`
	CryptographyScheme  insolar.PlatformCryptographyScheme `inject:""`
	KeyProcessor        insolar.KeyProcessor               `inject:""`

	// todo: remove
	// PulseManager        insolar.PulseManager               `inject:""`
	// ContractRequester   insolar.ContractRequester          `inject:""`

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

	n.cm.Inject(
		api.NewApiService(n.cfg.AdminAPIRunner),
		n,
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
