package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

type FetchJet struct {
	JetAccessor jet.Accessor
	Coordinator insolar.JetCoordinator
	JetUpdater  *jetTreeUpdater

	Parcel insolar.Parcel

	Res struct {
		JetID insolar.JetID
		Miss  bool
	}
}

func (p *FetchJet) Proceed(ctx context.Context) error {
	msg := p.Parcel.Message()
	if msg.DefaultTarget() == nil {
		return errors.New("unexpected message")
	}

	// Hack to temporary allow any genesis request.
	if p.Parcel.Pulse() <= insolar.FirstPulseNumber {
		p.Res.JetID = *insolar.NewJetID(0, nil)
		return nil
	}

	// Check token jet.
	token := p.Parcel.DelegationToken()
	if token != nil {
		// Calculate jet for target pulse.
		target := *msg.DefaultTarget().Record()
		pulse := target.Pulse()
		switch tm := msg.(type) {
		case *message.GetObject:
			pulse = tm.State.Pulse()
		case *message.GetChildren:
			if tm.FromChild == nil {
				return errors.New("fetching children without child pointer is forbidden")
			}
			pulse = tm.FromChild.Pulse()
		case *message.GetRequest:
			pulse = tm.Request.Pulse()
		}
		jetID, actual := p.JetAccessor.ForID(ctx, pulse, target)
		if !actual {
			inslogger.FromContext(ctx).WithFields(map[string]interface{}{
				"msg":   msg.Type().String(),
				"jet":   jetID.DebugString(),
				"pulse": pulse,
			}).Error("jet is not actual")
		}

		p.Res.JetID = jetID
		return nil
	}

	// Calculate jet for current pulse.
	var jetID insolar.ID
	if msg.DefaultTarget().Record().Pulse() == insolar.PulseNumberJet {
		jetID = *msg.DefaultTarget().Record()
	} else {
		j, err := p.JetUpdater.fetchJet(ctx, *msg.DefaultTarget().Record(), p.Parcel.Pulse())
		if err != nil {
			return errors.Wrap(err, "failed to fetch jet tree")
		}

		jetID = *j
	}

	// Check if jet is ours.
	node, err := p.Coordinator.LightExecutorForJet(ctx, jetID, p.Parcel.Pulse())
	if err != nil {
		return errors.Wrap(err, "failed to calculate executor for jet")
	}

	if *node != p.Coordinator.Me() {
		p.Res.Miss = true
		p.Res.JetID = insolar.JetID(jetID)
		return nil
	}

	ctx = addJetIDToLogger(ctx, jetID)

	p.Res.JetID = insolar.JetID(jetID)
	return nil
}

type WaitHot struct {
	Waiter HotDataWaiter

	Parcel insolar.Parcel
	JetID  insolar.JetID

	Res struct {
		Timeout bool
	}
}

func (p *WaitHot) Proceed(ctx context.Context) error {
	parcel := p.Parcel
	// Hack is needed for genesis:
	// because we don't have hot data on first pulse and without this we would stale.
	if parcel.Pulse() <= insolar.FirstPulseNumber {
		return nil
	}

	// If the call is a call in redirect-chain
	// skip waiting for the hot records
	if parcel.DelegationToken() != nil {
		return nil
	}

	err := p.Waiter.Wait(ctx, insolar.ID(p.JetID))
	if err != nil {
		p.Res.Timeout = true
	}

	return nil
}
