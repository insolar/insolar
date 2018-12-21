/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package fakepulsar

import (
	"context"

	"github.com/insolar/insolar/core"

	"github.com/insolar/insolar/configuration"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/transport/packet"
	"github.com/insolar/insolar/network/transport/packet/types"
)

type FakePulsarManager interface {
	Start(ctx context.Context)
	Stop(ctx context.Context)
}

func NewFakePulsarManager(callback callbackOnPulse, timeoutMs int32, cfg configuration.Configuration) FakePulsarManager {
	mngr := &FPManager{
		pulsar: NewFakePulsar(callback, timeoutMs),
	}
	mngr.transport = hostnetwork.NewFakePulsarTransport(mngr.messageHandler, cfg)
	return mngr
}

type FPManager struct {
	CertificateManager core.CertificateManager `inject:""`
	pulsar             *FakePulsar
	transport          hostnetwork.FakePulsarTransport
}

func (fpm *FPManager) Start(ctx context.Context) {
	fpm.transport.Start(ctx)
	if len(fpm.CertificateManager.GetCertificate().GetDiscoveryNodes()) == 0 {
		fpm.pulsar.Start(ctx)
		return
	}

}

func (fpm *FPManager) Stop(ctx context.Context) {
	fpm.transport.Stop(ctx)
	fpm.pulsar.Stop(ctx)
}

func (fpm *FPManager) messageHandler(msg *packet.Packet) {
	switch msg.Type {
	case types.FakePulsarResponse:
		p := msg.Data.(*FakePulsarResponse)
		if p == nil {
			log.Errorf("[ messageHandler ] failed to parse a message")
			return
		}
		fpm.pulsar.Start(p.firstPulseTime, p.pulseNum)
	case types.FakePulsarRequest:
		p := FakePulsarResponse{
			pulseNum:       fpm.pulsar.GetPulseNum(),
			firstPulseTime: fpm.pulsar.GetFirstPulseTime(),
		}
		builder := fpm.transport.NewBuilder()
		request := builder.Type(types.FakePulsarResponse).Data(p).Build()
		err := fpm.transport.SendRequest(request, msg.Sender)
		if err != nil {
			log.Errorf("[ messageHandler ] failed to send a request")
		}
	default:
		log.Errorf("[ FPManager.messageHandler ] received a wrong packet type from: %s", msg.Sender.String())
	}
}
