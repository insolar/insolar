package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network/consensus/common/args"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/census"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
	"github.com/insolar/insolar/network/consensus/gcpv2/core/coreapi"
	gcpErrors "github.com/insolar/insolar/network/consensus/gcpv2/core/errors"
	"github.com/insolar/insolar/pulse"
)

type RoundStrategyFactory interface {
	CreateRoundStrategy(online census.OnlinePopulation, config api.LocalNodeConfiguration) (RoundStrategy, PhaseControllersBundle)
}

type RoundStrategy interface {
	GetBaselineWeightForNeighbours() uint32
	ShuffleNodeSequence(n int, swap func(i, j int))

	ConfigureRoundContext(ctx context.Context, expectedPulse pulse.Number, self profiles.LocalNode) context.Context
}

var _ api.RoundController = &PhasedRoundController{}

type PhasedRoundController struct {
	rw sync.RWMutex

	/* Derived from the provided externally - set at init() or start(). Don't need mutex */
	chronicle api.ConsensusChronicles
	bundle    PhaseControllersBundle

	// fullCancel     context.CancelFunc /* cancels prepareCancel as well */
	prepareCancel context.CancelFunc

	roundWorker RoundStateMachineWorker

	/* Other fields - need mutex */
	prepR *PrepRealm
	realm FullRealm
}

func NewPhasedRoundController(strategy RoundStrategy, chronicle api.ConsensusChronicles, bundle PhaseControllersBundle,
	transport transport.Factory, config api.LocalNodeConfiguration,
	controlFeeder api.ConsensusControlFeeder, candidateFeeder api.CandidateControlFeeder, ephemeralFeeder api.EphemeralControlFeeder,
) *PhasedRoundController {

	r := &PhasedRoundController{chronicle: chronicle, bundle: bundle}

	latestCensus, _ := chronicle.GetLatestCensus()
	r.realm.coreRealm.initBefore(&r.rw, strategy, transport, config, latestCensus,
		controlFeeder, candidateFeeder, ephemeralFeeder)

	nbhSizes := r.realm.initBefore(transport)
	r.realm.coreRealm.initBeforePopulation(nbhSizes)

	return r
}

func (r *PhasedRoundController) PrepareConsensusRound(upstream api.UpstreamController) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.realm.coreRealm.roundContext = r.roundWorker.preInit(
		r.realm.coreRealm.strategy.ConfigureRoundContext(
			r.realm.config.GetParentContext(),
			r.realm.initialCensus.GetExpectedPulseNumber(),
			r.realm.GetLocalProfile(),
		), upstream, r.realm.coreRealm.controlFeeder, time.Second*2) // TODO parameterize the constant

	r.realm.coreRealm.stateMachine = &r.roundWorker

	r.realm.coreRealm.postponedPacketFn = func(packet transport.PacketParser, from endpoints.Inbound, verifyFlags coreapi.PacketVerifyFlags) bool {
		// There is no real context for delayed reprocessing, so we use the round context
		ctx := r.realm.coreRealm.roundContext
		inslogger.FromContext(ctx).Debugf("replayPacket %v", packet)
		_, err := r.handlePacket(ctx, packet, from, verifyFlags)
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
		}
		return true
	}

	inslogger.FromContext(r.realm.roundContext).Infof(
		"Starting consensus round: self={%v}, ephemeral=%v, bundle=%v, census=%+v", r.realm.GetLocalProfile(),
		r.realm.ephemeralFeeder != nil, r.bundle, r.realm.initialCensus)

	preps := r.bundle.CreatePrepPhaseControllers()
	if len(preps) == 0 {
		panic("illegal state - no prep realm")
	}

	prep := PrepRealm{coreRealm: &r.realm.coreRealm}
	prep.init(
		func(successful bool, startedAt time.Time) {
			// RUNS under lock
			if r.prepR == nil {
				return
			}
			defer r.prepR.stop() // initiates handover from PrepRealm
			r.prepR = nil

			r.realm.roundStartedAt = startedAt
			r.roundWorker.Start() // ensures that worker was started
			r._startFullRealm(successful)
		})

	var prepCtx context.Context
	// r.prepareCancel will be cancelled through r.fullCancel()
	prepCtx, r.prepareCancel = context.WithCancel(r.realm.roundContext)

	r.prepR = &prep
	r.prepR.beforeStart(prepCtx, preps)

	r.roundWorker.init(func() {
		// requires r.roundWorker.StartXXX to happen under lock
		r._setStartedAt()
		if r.prepR != nil { // PrepRealm can be finished before starting
			r.prepR._startWorkers(prepCtx, preps)
		}
	},
		// both further handlers MUST not use round's lock inside
		r.onConsensusStopper,
		nil, // r.onConsensusFinished,
	)

	r.realm.coreRealm.pollingWorker.Start(r.realm.roundContext, 100*time.Millisecond)
	r.prepR.prepareEphemeralPolling(prepCtx)
}

func (r *PhasedRoundController) onConsensusStopper() {

	latest, isExpected := r.chronicle.GetLatestCensus()
	var expt interface{}
	failed := false

	switch {
	case latest == r.realm.census:
		expt = "<nil>"
		failed = true
	case !isExpected:
		expt = "<unknown>"
	default:
		expected := latest.(census.Expected)
		if expected.GetPrevious() == r.realm.census {
			expt = expected
		} else {
			expt = "<unknown>"
		}
	}

	inslogger.FromContext(r.realm.roundContext).Infof(
		"Stopping consensus round: self={%v}, ephemeral=%v, bundle=%v, census=%+v, expected=%+v", r.realm.GetLocalProfile(),
		r.realm.ephemeralFeeder != nil, r.bundle, r.realm.census, expt)

	if failed {
		inslogger.FromContext(r.realm.roundContext).Panic("DEBUG FAIL-FAST: consensus didn't finish")
	}

	// TODO print purgatory
}

func (r *PhasedRoundController) _setStartedAt() {
	if r.realm.roundStartedAt.IsZero() { // can be called a few times
		r.realm.roundStartedAt = time.Now()
	}
}

func (r *PhasedRoundController) StopConsensusRound() {
	r.roundWorker.Stop()
}

func (r *PhasedRoundController) IsRunning() bool {
	return r.roundWorker.IsRunning()
}

func (r *PhasedRoundController) beforeHandlePacket() (prep *PrepRealm, current pulse.Number, possibleNext pulse.Number, ephemeralFeeder api.EphemeralControlFeeder) {

	r.rw.RLock()
	defer r.rw.RUnlock()
	if r.prepR != nil {
		if r.prepR.disableEphemeral {
			return r.prepR, r.realm.coreRealm.initialCensus.GetExpectedPulseNumber(), 0, nil
		}
		return r.prepR, r.realm.coreRealm.initialCensus.GetExpectedPulseNumber(), 0, r.realm.ephemeralFeeder
	}
	return nil, r.realm.GetPulseNumber(), r.realm.GetNextPulseNumber(), r.realm.ephemeralFeeder
}

/*
RUNS under lock.
Can be called from a polling function (for ephemeral), and happen BEFORE PrepRealm start
*/
func (r *PhasedRoundController) _startFullRealm(prepWasSuccessful bool) {

	if !prepWasSuccessful {
		r.roundWorker.OnPrepRoundFailed()
		return
	}

	r.roundWorker.OnFullRoundStarting()

	chronicle := r.chronicle
	lastCensus, isLastExpected := chronicle.GetLatestCensus()
	pd := r.realm.pulseData

	var active census.Active
	if lastCensus.GetCensusState() == census.PrimingCensus {
		/* This is the priming census */
		priming := lastCensus.GetMandateRegistry().GetPrimingCloudHash()
		active = lastCensus.(census.Prime).BuildCopy(pd, priming, priming).MakeExpected().MakeActive(pd)
	} else {
		// TODO restore to exact equality for expected population!!!!!
		if !lastCensus.GetPulseNumber().IsUnknownOrEqualTo(pd.PulseNumber) {
			// TODO inform control feeder when our pulse is less
			inslogger.FromContext(r.realm.roundContext).Panicf(
				"illegal state - pulse number of expected census (%v) and of the realm (%v) are mismatched for %v",
				lastCensus.GetPulseNumber(), pd.PulseNumber, r.realm.GetSelfNodeID())
		}
		if !isLastExpected {
			if lastCensus.GetOnlinePopulation().GetLocalProfile().IsJoiner() {
				inslogger.FromContext(r.realm.roundContext).Panic("DEBUG FAIL-FAST: local remains as joiner")
			}
			inslogger.FromContext(r.realm.roundContext).Panic("DEBUG FAIL-FAST: previous consensus didn't finish")
			// r.realm.unsafeRound = true
			// active = lastCensus.(census.Active)
		} else {
			/* Auto-activation of the prepared lastCensus */
			expCensus := lastCensus.(census.Expected)
			if !r.realm.unsafeRound {
				unsafe := true
				switch {
				case expCensus.GetPulseNumber() != pd.PulseNumber:
					inslogger.FromContext(r.realm.roundContext).Debugf("Unsafe round: expected=%d, pn=%d", expCensus.GetPulseNumber(), pd.PulseNumber)
				case !expCensus.GetPrevious().GetExpectedPulseNumber().IsUnknownOrEqualTo(pd.PulseNumber):
					inslogger.FromContext(r.realm.roundContext).Debugf("Unsafe round: prev.expected=%d, pn=%d", expCensus.GetPrevious().GetExpectedPulseNumber(), pd.PulseNumber)
				case !expCensus.GetOnlinePopulation().IsClean():
					inslogger.FromContext(r.realm.roundContext).Debugf("Unsafe round: population.clean=false, pn=%d", pd.PulseNumber)
				default:
					unsafe = false
				}
				r.realm.unsafeRound = unsafe
			}
			active = expCensus.MakeActive(pd)
		}
	}

	if r.realm.ephemeralFeeder != nil && !active.GetPulseData().IsFromEphemeral() {
		r.realm.ephemeralFeeder.OnEphemeralCancelled()
		r.realm.ephemeralFeeder = nil
		r.realm.unsafeRound = true
	}

	r.realm.start(active, active.GetOnlinePopulation(), r.bundle)

	endOf := r.realm.roundStartedAt.Add(r.realm.timings.EndOfConsensus)
	r.roundWorker.SetTimeout(endOf)

	inslogger.FromContext(r.realm.roundContext).Infof(
		"Starting consensus full realm: self={%v}, ephemeral=%v, unsafe=%v, startedAt=%v, endOf=%v, census=%+v, timings=%s", r.realm.GetLocalProfile(),
		r.realm.ephemeralFeeder != nil, r.realm.unsafeRound,
		args.LazyTimeFmt("15:04:05.000000", r.realm.GetStartedAt()),
		args.LazyTimeFmt("15:04:05.000000", endOf), active,
		r.realm.timings,
	)
}

func (r *PhasedRoundController) ensureStarted() bool {

	isStarted, isRunning := r.roundWorker.IsStartedAndRunning()
	if isStarted {
		return isRunning
	}

	r.rw.Lock() // ensure that starting closure will run under lock
	defer r.rw.Unlock()
	return r.roundWorker.SafeStartAndGetIsRunning()
}

func (r *PhasedRoundController) HandlePacket(ctx context.Context, packet transport.PacketParser,
	from endpoints.Inbound) (api.RoundControlCode, error) {

	inslogger.FromContext(ctx).Debugf("processPacket %v", packet)
	return r.handlePacket(ctx, packet, from, coreapi.DefaultVerify)
}

func (r *PhasedRoundController) handlePacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	verifyFlags coreapi.PacketVerifyFlags) (api.RoundControlCode, error) {

	isHandled, prep, err := r._handlePacket(ctx, packet, from, verifyFlags)
	if !isHandled || err == nil {
		return api.KeepRound, err
	}

	isPulse, pn := gcpErrors.IsMismatchPulseError(err)
	if !isPulse {
		return api.KeepRound, err
	}

	return r.handlePulseChange(ctx, pn, prep, err)
}

func (r *PhasedRoundController) _handlePacket(ctx context.Context, packet transport.PacketParser, from endpoints.Inbound,
	verifyFlags coreapi.PacketVerifyFlags) (bool, *PrepRealm, error) {

	pn := packet.GetPulseNumber()
	/* a separate method with lock is to ensure that further packet processing is not connected to a lock */
	prep, filterPN, _, ephemeralFeeder := r.beforeHandlePacket()

	// TODO HACK - network doesnt have information about pulsars to validate packets, hackIgnoreVerification must be removed when fixed
	const defaultOptions = coreapi.SkipVerify // coreapi.DefaultVerify

	if ephemeralFeeder != nil && !packet.GetPacketType().IsEphemeralPacket() { // TODO need fix, too ugly
		_, err := r.realm.VerifyPacketAuthenticity(ctx, packet, from, nil, coreapi.DefaultVerify, nil, defaultOptions)
		if err == nil {
			err = r.realm.ephemeralFeeder.OnNonEphemeralPacket(ctx, packet, from)
		}
		return false, nil, err
	}

	if prep != nil {
		// NB! Round may NOT be running yet here - ensure it is working before calling the state machine
		r.ensureStarted()

		if !pn.IsUnknown() && (filterPN.IsUnknown() || filterPN == pn) && r.roundWorker.IsRunning() /* can be already stopped */ {
			r.roundWorker.OnPulseDetected()
		}

		return true, prep, prep.dispatchPacket(ctx, packet, from, defaultOptions) // prep realm can't inherit flags
	}

	return true, nil, r.realm.dispatchPacket(ctx, packet, from, verifyFlags|defaultOptions)
}

func (r *PhasedRoundController) handlePulseChange(ctx context.Context, pn pulse.Number, prep *PrepRealm, origErr error) (api.RoundControlCode, error) {

	var pulseControl api.PulseControlFeeder
	if r.realm.ephemeralFeeder != nil && (prep == nil || !prep.disableEphemeral) {
		pulseControl = r.realm.ephemeralFeeder
	} else {
		pulseControl = r.realm.controlFeeder
	}

	var epn pulse.Number
	var c census.Operational

	if prep != nil {
		c = prep.initialCensus
		epn = c.GetPulseNumber()
	} else {
		c = r.realm.census
		epn = c.GetExpectedPulseNumber()
	}

	expected := r.chronicle.GetExpectedCensus()
	if expected != nil {
		if expected.GetPrevious() != c {
			inslogger.FromContext(ctx).Warnf("unable to switch a past round/population")
			return api.KeepRound, origErr
		}
		epn = expected.GetPulseNumber()
	}

	switch {
	case epn.IsUnknownOrEqualTo(epn):
		break

	case pn < epn:
		r.roundWorker.onUnexpectedPulse(pn)
		return api.KeepRound, origErr

	case c.GetCensusState() == census.PrimingCensus:
		panic(fmt.Sprintf("unable to fast-forward a priming census: %s", origErr.Error()))

	default:
		_, pd := c.GetNearestPulseData()
		if !pulseControl.CanFastForwardPulse(epn, pn, pd) {
			r.roundWorker.onUnexpectedPulse(pn)
			return api.KeepRound, origErr
		}
	}

	switch {
	case !r.roundWorker.IsRunning():
		latest, _ := r.chronicle.GetLatestCensus()
		if c == latest && !c.GetOnlinePopulation().IsValid() {
			return api.NextRoundTerminate, fmt.Errorf("current population is invalid and an expected population is missing: %v", origErr.Error())
		}
		if expected == nil {
			return api.KeepRound, origErr
		}

		inslogger.FromContext(ctx).Debug("switch to a next round by changed pulse: ", origErr)
	default:
		endOfConsensus := r.realm.GetStartedAt().Add(r.realm.timings.EndOfConsensus)
		if time.Now().Before(endOfConsensus) && !pulseControl.CanStopOnHastyPulse(pn, endOfConsensus) {
			return api.KeepRound, fmt.Errorf("too early: %v", origErr)
		}
		if expected == nil {
			return api.KeepRound, origErr
		}

		inslogger.FromContext(ctx).Debug("stopping round by changed pulse: ", origErr)
	}

	if !expected.GetPulseNumber().IsUnknownOrEqualTo(pn) {
		expected.Rebuild(pn).MakeExpected()
	}
	r.roundWorker.onNextPulse(pn)

	return api.StartNextRound, nil
}
