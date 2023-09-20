package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/misbehavior"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/errors"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/packetdispatch"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/population"
	"github.com/insolar/insolar/pulse"
)

// hides embedded pointer from external access
type hLocker interface {
	sync.Locker
	RLock()
	RUnlock()
}

type coreRealm struct {
	/* Provided externally at construction. Don't need mutex */
	hLocker

	roundContext    context.Context
	strategy        RoundStrategy
	config          api.LocalNodeConfiguration
	initialCensus   census.Operational
	controlFeeder   api.ConsensusControlFeeder
	candidateFeeder api.CandidateControlFeeder
	ephemeralFeeder api.EphemeralControlFeeder

	pollingWorker coreapi.PollingWorker

	/* Derived from the ones provided externally - set at init() or start(). Don't need mutex */
	signer            cryptkit.DigestSigner
	digest            transport.ConsensusDigestFactory
	assistant         transport.CryptographyAssistant
	stateMachine      api.RoundStateCallback
	roundStartedAt    time.Time
	postponedPacketFn packetdispatch.PostponedPacketFunc

	expectedPopulationSize member.Index
	nbhSizes               transport.NeighbourhoodSizes

	self *population.NodeAppearance /* Special case - this field is set twice, by start() of PrepRealm and FullRealm */

	requestedPowerFlag bool

	/*
		Other fields - need mutex during PrepRealm, unless accessed by start() of PrepRealm
		FullRealm doesnt need a lock to read them
	*/
	pulseData     pulse.Data
	originalPulse proofs.OriginalPulsarPacket
	unsafeRound   bool // nolint: structcheck
}

func (r *coreRealm) initBefore(hLocker hLocker, strategy RoundStrategy, transport transport.Factory,
	config api.LocalNodeConfiguration, initialCensus census.Operational,
	controlFeeder api.ConsensusControlFeeder, candidateFeeder api.CandidateControlFeeder,
	ephemeralFeeder api.EphemeralControlFeeder) {

	r.hLocker = hLocker

	r.strategy = strategy
	r.config = config
	r.initialCensus = initialCensus

	r.assistant = transport.GetCryptographyFactory()
	r.digest = r.assistant.GetDigestFactory()

	sks := config.GetSecretKeyStore()
	r.signer = r.assistant.CreateNodeSigner(sks)

	r.controlFeeder = controlFeeder
	r.candidateFeeder = candidateFeeder
	r.ephemeralFeeder = ephemeralFeeder
}

func (r *coreRealm) initBeforePopulation(nbhSizes transport.NeighbourhoodSizes) {

	r.nbhSizes = nbhSizes
	pop := r.initialCensus.GetOnlinePopulation()

	/*
		Here we initialize self like for PrepRealm. This is not perfect, but it enables access to Self earlier.
	*/
	profile := pop.GetLocalProfile()

	if profile.GetOpMode().IsEvicted() {
		/*
			Previous round has provided an incorrect population, as eviction of local node must force one-node population of self
		*/
		panic("illegal state")
	}

	pn := r.initialCensus.GetExpectedPulseNumber()

	/* Will only be used during PrepRealm */
	selfNodeHookTmp := population.NewHook(profile,
		population.NewPanicDispatcher("updates of stub-self are not allowed"),
		population.NewSharedNodeContextByPulseNumber(r.assistant, pn, 0, r.getEphemeralMode(),
			func(report misbehavior.Report) interface{} {
				inslogger.FromContext(r.roundContext).Warnf("Got Report: %+v", report)
				r.initialCensus.GetMisbehaviorRegistry().AddReport(report)
				return nil
			},
		))

	powerRequest := r.controlFeeder.GetRequiredPowerLevel()
	r.requestedPowerFlag = !powerRequest.IsEmpty()
	selfNode := population.NewNodeAppearanceAsSelf(profile, powerRequest, &selfNodeHookTmp)

	r.self = &selfNode

	r.expectedPopulationSize = member.AsIndex(pop.GetIndexedCapacity())
}

func (r *coreRealm) getEphemeralMode() api.EphemeralMode {
	if r.ephemeralFeeder == nil {
		return api.EphemeralNotAllowed
	}
	return api.EphemeralAllowed
}

func (r *coreRealm) GetStrategy() RoundStrategy {
	return r.strategy
}

func (r *coreRealm) GetVerifierFactory() cryptkit.SignatureVerifierFactory {
	return r.assistant
}

func (r *coreRealm) GetDigestFactory() transport.ConsensusDigestFactory {
	return r.digest
}

func (r *coreRealm) GetSigner() cryptkit.DigestSigner {
	return r.signer
}

func (r *coreRealm) GetSignatureVerifier(pks cryptkit.PublicKeyStore) cryptkit.SignatureVerifier {
	return r.assistant.CreateSignatureVerifierWithPKS(pks)
}

func (r *coreRealm) GetStartedAt() time.Time {
	// r.Lock()
	// defer r.Unlock()

	if r.roundStartedAt.IsZero() {
		panic("illegal state")
	}
	return r.roundStartedAt
}

func (r *coreRealm) AdjustedAfter(d time.Duration) time.Duration {
	return time.Until(r.GetStartedAt().Add(d))
}

func (r *coreRealm) GetRoundContext() context.Context {
	return r.roundContext
}

func (r *coreRealm) GetLocalConfig() api.LocalNodeConfiguration {
	return r.config
}

// polling fn must be fast, and it will remain in polling until it returns false
func (r *coreRealm) AddPoll(fn api.MaintenancePollFunc) {
	r.pollingWorker.AddPoll(fn)
}

func (r *coreRealm) VerifyPacketAuthenticity(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound, sourceNode packetdispatch.MemberPacketReceiver, unverifiedFlag coreapi.PacketVerifyFlags,
	pd population.PacketDispatcher, verifyFlags coreapi.PacketVerifyFlags) (coreapi.PacketVerifyFlags, error) {

	if verifyFlags&(coreapi.SkipVerify|coreapi.SuccessfullyVerified) != 0 {
		return 0, nil
	}

	var err error
	strict := verifyFlags&coreapi.RequireStrictVerify != 0

	switch {
	case sourceNode != nil:
		err = sourceNode.VerifyPacketAuthenticity(packet.GetPacketSignature(), from, strict)
		if err != nil {
			return 0, err
		}
		verifyFlags |= coreapi.SuccessfullyVerified
	case pd != nil && pd.HasCustomVerifyForHost(from, verifyFlags):
		// skip default behavior
	default:
		sourceID := packet.GetSourceID()

		nr := coreapi.FindHostProfile(sourceID, from, r.initialCensus)
		if nr != nil {
			sf := r.assistant.CreateSignatureVerifierWithPKS(nr.GetPublicKeyStore())
			err := coreapi.VerifyPacketAuthenticityBy(packet.GetPacketSignature(), nr, sf, from, strict)
			if err != nil {
				return 0, err
			}
			verifyFlags |= coreapi.SuccessfullyVerified

		} else {
			pt := packet.GetPacketType()

			if strict || verifyFlags&coreapi.AllowUnverified == 0 {
				return 0, fmt.Errorf("unable to verify sender for packet type (%v): from=%v", pt, from)
			}
			inslogger.FromContext(ctx).Errorf("unable to verify sender for packet type (%v): from=%v", pt, from)
			verifyFlags |= unverifiedFlag
		}
	}
	return verifyFlags, nil
}

func (r *coreRealm) VerifyPacketPulseNumber(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	filterPN, nextPN pulse.Number, details string) (bool, error) {

	pn := packet.GetPulseNumber()
	if filterPN == pn || filterPN.IsUnknown() || pn.IsUnknown() {
		return true, nil
	}

	sourceID := packet.GetSourceID()
	localID := r.self.GetNodeID()

	return false, errors.NewPulseRoundMismatchErrorDef(pn, filterPN, localID, sourceID, details)
}
