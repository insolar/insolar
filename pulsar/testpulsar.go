//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package pulsar

import (
	"context"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/pulse"
)

type TestPulsar struct {
	distributor   insolar.PulseDistributor
	generator     entropygenerator.EntropyGenerator
	configuration configuration.Pulsar
}

func NewTestPulsar(
	configuration configuration.Pulsar,
	distributor insolar.PulseDistributor,
	generator entropygenerator.EntropyGenerator,
) *TestPulsar {
	return &TestPulsar{
		distributor:   distributor,
		generator:     generator,
		configuration: configuration,
	}
}

func (p *TestPulsar) SendPulse(ctx context.Context) error {
	timeNow := time.Now()
	pulseNumber := insolar.PulseNumber(pulse.OfTime(timeNow))

	pls := insolar.Pulse{
		PulseNumber:      pulseNumber,
		Entropy:          p.generator.GenerateEntropy(),
		NextPulseNumber:  pulseNumber + insolar.PulseNumber(p.configuration.NumberDelta),
		PrevPulseNumber:  pulseNumber - insolar.PulseNumber(p.configuration.NumberDelta),
		EpochPulseNumber: int(pulseNumber),
		OriginID:         [16]byte{206, 41, 229, 190, 7, 240, 162, 155, 121, 245, 207, 56, 161, 67, 189, 0},
	}

	var err error
	pls.Signs, err = getPSC(pls)
	if err != nil {
		log.Errorf("[ distribute ]", err)
		return err
	}

	pls.PulseTimestamp = time.Now().UnixNano()

	p.distributor.Distribute(ctx, pls)

	return nil
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

	payload := PulseSenderConfirmationPayload{PulseSenderConfirmation: psc}
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
