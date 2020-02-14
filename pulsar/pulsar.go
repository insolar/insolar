// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package pulsar

import (
	"context"
	"sync"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/log"
)

// Pulsar is a base struct for pulsar's node
// It contains all the stuff, which is needed for working of a pulsar
type Pulsar struct {
	Config       configuration.Pulsar
	PublicKeyRaw string

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
) *Pulsar {

	log.Info("[NewPulsar]")

	pulsar := &Pulsar{
		CryptographyService:        cryptographyService,
		PlatformCryptographyScheme: scheme,
		KeyProcessor:               keyProcessor,
		PulseDistributor:           pulseDistributor,
		Config:                     configuration,
		EntropyGenerator:           entropyGenerator,
	}

	pubKey, err := cryptographyService.GetPublicKey()
	if err != nil {
		log.Fatal(err)
	}
	pubKeyRaw, err := keyProcessor.ExportPublicKeyPEM(pubKey)
	if err != nil {
		log.Fatal(err)
	}
	pulsar.PublicKeyRaw = string(pubKeyRaw)

	return pulsar
}

func (p *Pulsar) Send(ctx context.Context, pulseNumber insolar.PulseNumber) error {
	logger := inslogger.FromContext(ctx)
	logger.Infof("before sending new pulseNumber: %v", pulseNumber)

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
		EpochPulseNumber: pulseNumber.AsEpoch(),
		OriginID:         [16]byte{206, 41, 229, 190, 7, 240, 162, 155, 121, 245, 207, 56, 161, 67, 189, 0},
		PulseTimestamp:   time.Now().UnixNano(),
		Signs:            map[string]insolar.PulseSenderConfirmation{},
	}

	payload := PulseSenderConfirmationPayload{PulseSenderConfirmation: insolar.PulseSenderConfirmation{
		ChosenPublicKey: p.PublicKeyRaw,
		Entropy:         entropy,
		PulseNumber:     pulseNumber,
	}}
	hasher := platformpolicy.NewPlatformCryptographyScheme().IntegrityHasher()
	hash, err := payload.Hash(hasher)
	if err != nil {
		return err
	}
	signature, err := p.CryptographyService.Sign(hash)
	if err != nil {
		return err
	}

	pulseForSending.Signs[p.PublicKeyRaw] = insolar.PulseSenderConfirmation{
		ChosenPublicKey: p.PublicKeyRaw,
		Signature:       signature.Bytes(),
		Entropy:         entropy,
		PulseNumber:     pulseNumber,
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

	stats.Record(ctx, statPulseGenerated.M(1), statCurrentPulse.M(int64(pulseNumber.AsUint32())))
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

// PulseSenderConfirmationPayload is a struct with info about pulse's confirmations
type PulseSenderConfirmationPayload struct {
	insolar.PulseSenderConfirmation
}

// Hash calculates hash of payload
func (ps *PulseSenderConfirmationPayload) Hash(hashProvider insolar.Hasher) ([]byte, error) {
	_, err := hashProvider.Write(ps.PulseNumber.Bytes())
	if err != nil {
		return nil, err
	}
	_, err = hashProvider.Write([]byte(ps.ChosenPublicKey))
	if err != nil {
		return nil, err
	}
	_, err = hashProvider.Write(ps.Entropy[:])
	if err != nil {
		return nil, err
	}
	return hashProvider.Sum(nil), nil
}

/*

if currentPulsar.isStandalone() {
currentPulsar.SetCurrentSlotEntropy(currentPulsar.GetGeneratedEntropy())
currentPulsar.CurrentSlotPulseSender = currentPulsar.PublicKeyRaw

payload := PulseSenderConfirmationPayload{insolar.PulseSenderConfirmation{
ChosenPublicKey: currentPulsar.CurrentSlotPulseSender,
Entropy:         *currentPulsar.GetCurrentSlotEntropy(),
PulseNumber:     currentPulsar.ProcessingPulseNumber,
}}
hashProvider := currentPulsar.PlatformCryptographyScheme.IntegrityHasher()
hash, err := payload.Hash(hashProvider)
if err != nil {
currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
return
}
signature, err := currentPulsar.CryptographyService.Sign(hash)
if err != nil {
currentPulsar.StateSwitcher.SwitchToState(ctx, Failed, err)
return
}

currentPulsar.currentSlotSenderConfirmationsLock.Lock()
currentPulsar.CurrentSlotSenderConfirmations[currentPulsar.PublicKeyRaw] = insolar.PulseSenderConfirmation{
ChosenPublicKey: currentPulsar.CurrentSlotPulseSender,
Signature:       signature.Bytes(),
Entropy:         *currentPulsar.GetCurrentSlotEntropy(),
PulseNumber:     currentPulsar.ProcessingPulseNumber,
}
currentPulsar.currentSlotSenderConfirmationsLock.Unlock()

currentPulsar.StateSwitcher.SwitchToState(ctx, SendingPulse, nil)

return
}*/
