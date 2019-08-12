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
	"sync"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/log"
)

// Pulsar is a base struct for pulsar's node
// It contains all the stuff, which is needed for working of a pulsar
type Pulsar struct {
	Config configuration.Pulsar

	EntropyGenerator entropygenerator.EntropyGenerator

	Certificate                certificate.Certificate
	CryptographyService        insolar.CryptographyService
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme
	KeyProcessor               insolar.KeyProcessor
	PulseDistributor           insolar.PulseDistributor

	lastPNMutex sync.RWMutex
	lastPN      insolar.PulseNumber
}

// NewPulsar creates a new pulse with using of custom GeneratedEntropy Generator
func NewPulsar(
	configuration configuration.Pulsar,
	cryptographyService insolar.CryptographyService,
	scheme insolar.PlatformCryptographyScheme,
	keyProcessor insolar.KeyProcessor,
	pulseDistributor insolar.PulseDistributor,
	entropyGenerator entropygenerator.EntropyGenerator,
) (*Pulsar, error) {

	log.Info("[NewPulsar]")

	pulsar := &Pulsar{
		CryptographyService:        cryptographyService,
		PlatformCryptographyScheme: scheme,
		KeyProcessor:               keyProcessor,
		PulseDistributor:           pulseDistributor,
		Config:                     configuration,
		EntropyGenerator:           entropyGenerator,
	}

	return pulsar, nil
}

func (p *Pulsar) Send(ctx context.Context, pulseNumber insolar.PulseNumber) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("before sending new pulseNumber: %v", pulseNumber)

	entropy, _, err := p.generateNewEntropyAndSign()
	if err != nil {
		logger.Error(err)
		return err
	}

	pulseForSending := insolar.Pulse{
		PulseNumber:      pulseNumber,
		Entropy:          entropy,
		NextPulseNumber:  pulseNumber + insolar.PulseNumber(p.Config.NumberDelta),
		PrevPulseNumber:  p.lastPN,
		EpochPulseNumber: int(pulseNumber),
		OriginID:         [16]byte{206, 41, 229, 190, 7, 240, 162, 155, 121, 245, 207, 56, 161, 67, 189, 0},
		PulseTimestamp:   time.Now().UnixNano(),
	}

	logger.Debug("Start a process of sending pulse")
	go func() {
		logger.Debug("Before sending to network")
		p.PulseDistributor.Distribute(ctx, pulseForSending)
	}()

	p.lastPNMutex.Lock()
	p.lastPN = pulseNumber
	p.lastPNMutex.Unlock()
	logger.Infof("set latest pulse: %v", pulseForSending.PulseNumber)

	stats.Record(ctx, statPulseGenerated.M(1))
	return nil
}

func (p *Pulsar) LastPN() insolar.PulseNumber {
	p.lastPNMutex.RLock()
	defer p.lastPNMutex.RUnlock()

	return p.lastPN
}

func (p *Pulsar) generateNewEntropyAndSign() (insolar.Entropy, []byte, error) {
	e := p.EntropyGenerator.GenerateEntropy()

	sign, err := p.CryptographyService.Sign(e[:])
	if err != nil {
		return insolar.Entropy{}, nil, err
	}

	return e, sign.Bytes(), nil
}
