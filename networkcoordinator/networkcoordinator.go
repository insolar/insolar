/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification,
 * are permitted (subject to the limitations in the disclaimer below) provided that
 * the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *  * Neither the name of Insolar Technologies nor the names of its contributors
 *    may be used to endorse or promote products derived from this software
 *    without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
 * BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND
 * CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING,
 * BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
 * HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
 * SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
 * DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
 * THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package networkcoordinator

import (
	"context"

	"github.com/insolar/insolar/core"
)

// NetworkCoordinator encapsulates logic of network configuration
type NetworkCoordinator struct {
	CertificateManager core.CertificateManager  `inject:""`
	NetworkSwitcher    core.NetworkSwitcher     `inject:""`
	ContractRequester  core.ContractRequester   `inject:""`
	MessageBus         core.MessageBus          `inject:""`
	CS                 core.CryptographyService `inject:""`
	PS                 core.PulseStorage        `inject:""`

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
	nc.MessageBus.MustRegister(core.TypeNodeSignRequest, nc.signCertHandler)

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
	if nc.NetworkSwitcher.GetState() == core.CompleteNetworkState {
		return nc.realCoordinator
	}
	return nc.zeroCoordinator
}

// IsStarted returns true if component was started and false in other way
func (nc *NetworkCoordinator) IsStarted() bool {
	return nc.isStarted
}

// GetCert method returns node certificate by requesting sign from discovery nodes
func (nc *NetworkCoordinator) GetCert(ctx context.Context, registeredNodeRef *core.RecordRef) (core.Certificate, error) {
	return nc.getCoordinator().GetCert(ctx, registeredNodeRef)
}

// ValidateCert validates node certificate
func (nc *NetworkCoordinator) ValidateCert(ctx context.Context, certificate core.AuthorizationCertificate) (bool, error) {
	return nc.CertificateManager.VerifyAuthorizationCertificate(certificate)
}

// signCertHandler is MsgBus handler that signs certificate for some node with node own key
func (nc *NetworkCoordinator) signCertHandler(ctx context.Context, p core.Parcel) (core.Reply, error) {
	return nc.getCoordinator().signCertHandler(ctx, p)
}

// SetPulse writes pulse data on local storage
func (nc *NetworkCoordinator) SetPulse(ctx context.Context, pulse core.Pulse) error {
	return nc.getCoordinator().SetPulse(ctx, pulse)
}
