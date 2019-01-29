/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package servicenetwork

import (
	"context"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network/pulsenetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/pkg/errors"
)

type TestPulsar interface {
	component.Starter
	component.Stopper
}

func NewTestPulsar(pulseTimeMs int, bootstrapHosts []string) (TestPulsar, error) {
	transportCfg := configuration.Transport{}
	tp, err := transport.NewTransport(transportCfg, relay.NewProxy())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create distributor transport")
	}

	distributorCfg := configuration.PulseDistributor{}
	distributor, err := pulsenetwork.NewDistributor(distributorCfg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create pulse distributor")
	}
	return &testPulsar{
		transport:      tp,
		distributor:    distributor,
		generator:      &entropygenerator.StandardEntropyGenerator{},
		pulseTimeMs:    pulseTimeMs,
		bootstrapHosts: bootstrapHosts,
	}, nil
}

type testPulsar struct {
	transport   transport.Transport
	distributor core.PulseDistributor
	generator   entropygenerator.EntropyGenerator
	cm          *component.Manager

	pulseTimeMs    int
	bootstrapHosts []string
}

func (tp *testPulsar) Start(ctx context.Context) error {
	tp.cm = &component.Manager{}
	tp.cm.Inject(tp.transport, tp.distributor)

	if err := tp.cm.Init(ctx); err != nil {
		return errors.Wrap(err, "Failed to init test pulsar components")
	}
	if err := tp.cm.Start(ctx); err != nil {
		return errors.Wrap(err, "Failed to start test pulsar components")
	}
	go tp.distribute()
	return nil
}

func (tp *testPulsar) distribute() {
	// TODO: distribute pulse
}

func (tp *testPulsar) Stop(ctx context.Context) error {
	if err := tp.cm.Stop(ctx); err != nil {
		return errors.Wrap(err, "Failed to stop test pulsar components")
	}
	return nil
}
