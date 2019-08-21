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

// +build networktest

package tests

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/pulsenetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulse"
)

type TestPulsar interface {
	Start(ctx context.Context, bootstrapHosts []string) error
	Pause()
	Continue()
	component.Stopper
}

func NewTestPulsar(requestsTimeoutMs, pulseDelta int32) (TestPulsar, error) {

	return &testPulsar{
		generator:         &entropygenerator.StandardEntropyGenerator{},
		reqTimeoutMs:      requestsTimeoutMs,
		pulseDelta:        pulseDelta,
		cancellationToken: make(chan struct{}),
	}, nil
}

type testPulsar struct {
	distributor insolar.PulseDistributor
	generator   entropygenerator.EntropyGenerator
	cm          *component.Manager

	activityMutex sync.Mutex

	reqTimeoutMs int32
	pulseDelta   int32

	cancellationToken chan struct{}
}

func (tp *testPulsar) Start(ctx context.Context, bootstrapHosts []string) error {

	distributorCfg := configuration.PulseDistributor{
		BootstrapHosts:      bootstrapHosts,
		PulseRequestTimeout: tp.reqTimeoutMs,
	}

	var err error
	tp.distributor, err = pulsenetwork.NewDistributor(distributorCfg)
	if err != nil {
		return errors.Wrap(err, "Failed to create pulse distributor")
	}

	tp.cm = &component.Manager{}
	if UseFakeTransport {
		tp.cm.Register(transport.NewFakeFactory(configuration.NewHostNetwork().Transport))
	} else {
		tp.cm.Register(transport.NewFactory(configuration.NewHostNetwork().Transport))
	}
	tp.cm.Inject(tp.distributor)

	if err = tp.cm.Init(ctx); err != nil {
		return errors.Wrap(err, "Failed to init test pulsar components")
	}
	if err = tp.cm.Start(ctx); err != nil {
		return errors.Wrap(err, "Failed to start test pulsar components")
	}

	go tp.distribute(ctx)
	return nil
}

func (tp *testPulsar) Pause() {
	tp.activityMutex.Lock()
}

func (tp *testPulsar) Continue() {
	tp.activityMutex.Unlock()
}

func (tp *testPulsar) distribute(ctx context.Context) {
	timeNow := time.Now()
	pulseNumber := insolar.PulseNumber(pulse.OfTime(timeNow))

	pls := insolar.Pulse{
		PulseNumber:      pulseNumber,
		Entropy:          tp.generator.GenerateEntropy(),
		NextPulseNumber:  pulseNumber + insolar.PulseNumber(tp.pulseDelta),
		PrevPulseNumber:  pulseNumber - insolar.PulseNumber(tp.pulseDelta),
		EpochPulseNumber: int(pulseNumber),
		OriginID:         [16]byte{206, 41, 229, 190, 7, 240, 162, 155, 121, 245, 207, 56, 161, 67, 189, 0},
	}

	var err error
	pls.Signs, err = getPSC(pls)
	if err != nil {
		log.Errorf("[ distribute ]", err)
	}

	for {
		select {
		case <-time.After(time.Duration(tp.pulseDelta) * time.Second):
			go func(pulse insolar.Pulse) {
				tp.activityMutex.Lock()
				defer tp.activityMutex.Unlock()

				pulse.PulseTimestamp = time.Now().UnixNano()

				tp.distributor.Distribute(ctx, pulse)
			}(pls)

			pls = tp.incrementPulse(pls)
		case <-tp.cancellationToken:
			return
		}
	}
}

func (tp *testPulsar) incrementPulse(pulse insolar.Pulse) insolar.Pulse {
	newPulseNumber := pulse.PulseNumber + insolar.PulseNumber(tp.pulseDelta)
	newPulse := insolar.Pulse{
		PulseNumber:      newPulseNumber,
		Entropy:          tp.generator.GenerateEntropy(),
		NextPulseNumber:  newPulseNumber + insolar.PulseNumber(tp.pulseDelta),
		PrevPulseNumber:  pulse.PulseNumber,
		EpochPulseNumber: pulse.EpochPulseNumber,
		OriginID:         pulse.OriginID,
		PulseTimestamp:   time.Now().UnixNano(),
		Signs:            pulse.Signs,
	}
	var err error
	newPulse.Signs, err = getPSC(newPulse)
	if err != nil {
		log.Errorf("[ incermentPulse ]", err)
	}
	return newPulse
}

func getPSC(pulse insolar.Pulse) (map[string]insolar.PulseSenderConfirmation, error) {
	proc := platformpolicy.NewKeyProcessor()
	key, err := proc.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	pem, err := proc.ExportPublicKeyPEM(proc.ExtractPublicKey(key))
	if err != nil {
		return nil, err
	}
	result := make(map[string]insolar.PulseSenderConfirmation)
	psc := insolar.PulseSenderConfirmation{
		PulseNumber:     pulse.PulseNumber,
		ChosenPublicKey: string(pem),
		Entropy:         pulse.Entropy,
	}

	payload := pulsar.PulseSenderConfirmationPayload{PulseSenderConfirmation: psc}
	hasher := platformpolicy.NewPlatformCryptographyScheme().IntegrityHasher()
	hash, err := payload.Hash(hasher)
	if err != nil {
		return nil, err
	}
	service := cryptography.NewKeyBoundCryptographyService(key)
	sign, err := service.Sign(hash)
	if err != nil {
		return nil, err
	}

	psc.Signature = sign.Bytes()
	result[string(pem)] = psc

	return result, nil
}

func (tp *testPulsar) Stop(ctx context.Context) error {
	if err := tp.cm.Stop(ctx); err != nil {
		return errors.Wrap(err, "Failed to stop test pulsar components")
	}
	close(tp.cancellationToken)
	return nil
}
