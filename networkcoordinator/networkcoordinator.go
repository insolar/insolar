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

package networkcoordinator

import (
	"context"

	"github.com/insolar/insolar/network"

	"github.com/insolar/insolar/insolar"
)

// NetworkCoordinator encapsulates logic of network configuration
type NetworkCoordinator struct {
	CertificateManager insolar.CertificateManager  `inject:""`
	ContractRequester  insolar.ContractRequester   `inject:""`
	MessageBus         insolar.MessageBus          `inject:""`
	Gatewayer          network.Gatewayer           `inject:""`
	CS                 insolar.CryptographyService `inject:""`

	realCoordinator Coordinator
	zeroCoordinator Coordinator
	isStarted       bool
}

// New creates new NetworkCoordinator
func New() (*NetworkCoordinator, error) {
	return &NetworkCoordinator{}, nil
}

// Start implements interface of Component
func (nc *NetworkCoordinator) Start(ctx context.Context) error {
	nc.MessageBus.MustRegister(insolar.TypeNodeSignRequest, nc.signCertHandler)

	nc.zeroCoordinator = newZeroNetworkCoordinator()
	nc.realCoordinator = newRealNetworkCoordinator(
		nc.CertificateManager,
		nc.ContractRequester,
		nc.MessageBus,
		nc.CS,
	)
	nc.isStarted = true
	return nil
}

func (nc *NetworkCoordinator) getCoordinator() Coordinator {
	if nc.Gatewayer.Gateway().GetState() == insolar.CompleteNetworkState {
		return nc.realCoordinator
	}
	return nc.zeroCoordinator
}

// IsStarted returns true if component was started and false in other way
func (nc *NetworkCoordinator) IsStarted() bool {
	return nc.isStarted
}

// GetCert method returns node certificate by requesting sign from discovery nodes
func (nc *NetworkCoordinator) GetCert(ctx context.Context, registeredNodeRef *insolar.Reference) (insolar.Certificate, error) {
	return nc.getCoordinator().GetCert(ctx, registeredNodeRef)
}

// ValidateCert validates node certificate
func (nc *NetworkCoordinator) ValidateCert(ctx context.Context, certificate insolar.AuthorizationCertificate) (bool, error) {
	return nc.CertificateManager.VerifyAuthorizationCertificate(certificate)
}

// signCertHandler is MsgBus handler that signs certificate for some node with node own key
func (nc *NetworkCoordinator) signCertHandler(ctx context.Context, p insolar.Parcel) (insolar.Reply, error) {
	return nc.getCoordinator().signCertHandler(ctx, p)
}

// SetPulse writes pulse data on local storage
func (nc *NetworkCoordinator) SetPulse(ctx context.Context, pulse insolar.Pulse) error {
	return nc.getCoordinator().SetPulse(ctx, pulse)
}
