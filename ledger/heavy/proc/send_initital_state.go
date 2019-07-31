package proc

import (
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/bus"
	"context"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/ledger/heavy/executor"
)

type SendInitialState struct {
	meta payload.Meta

	dep struct {
		startPulse     pulse.StartPulse
		jetKeeper      executor.JetKeeper
		jetTree        jet.Storage
		jetCoordinator jet.Coordinator
		dropDB         *drop.DB
		pulseAccessor pulse.Accessor
		sender bus.Sender
	}
}

func (p *SendInitialState) Dep(
	startPulse pulse.StartPulse,
	jetKeeper  executor.JetKeeper,
	jetTree    jet.Storage,
	jetCoordinator jet.Coordinator,
	dropDB          *drop.DB,
	pulseAccessor pulse.Accessor,
	sender bus.Sender,
) {
	p.dep.startPulse = startPulse
	p.dep.jetKeeper = jetKeeper
	p.dep.jetTree = jetTree
	p.dep.jetCoordinator = jetCoordinator
	p.dep.dropDB = dropDB
	p.dep.pulseAccessor = pulseAccessor
	p.dep.sender = sender
}

func NewSendInitialState(meta payload.Meta) *SendInitialState {
	return &SendInitialState{
		meta: meta,
	}
}

func (p *SendInitialState) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	startPulse, err := p.dep.startPulse.PulseNumber()
	if err != nil {

		logger.Fatal("Couldn't get start pulse", err)
	}
	msg, err := payload.Unmarshal(p.meta.Payload)
	if err != nil {
		logger.Fatal("Couldn't unmarshall request", err)
	}
	req := msg.(*payload.GetLightInitialState)

	if req.Pulse == startPulse {
		topSyncPulseNumber := p.dep.jetKeeper.TopSyncPulse()
		var IDs []insolar.JetID
		var drops [][]byte
		for _, id := range p.dep.jetTree.All(ctx, topSyncPulseNumber) {
			light, err := p.dep.jetCoordinator.LightExecutorForJet(ctx, insolar.ID(id), req.Pulse)
			if err != nil {
				logger.Fatal("Couldn't receive light executor for jet: ", id, " ", err)
			}
			if light.Equal(p.meta.Sender) {
				IDs = append(IDs, id)
				dr, err := p.dep.dropDB.ForPulse(ctx, id, topSyncPulseNumber)
				if err != nil {
					logger.Fatal("Couldn't get drops for jet: ", id, " ", err)
				}
				drops = append(drops, drop.MustEncode(&dr))
			}
		}


		topSyncPulse, err := p.dep.pulseAccessor.ForPulseNumber(ctx, topSyncPulseNumber)
		if err != nil {
			logger.Fatal("Couldn't get pulse for topSyncPulse: ", topSyncPulseNumber, " ", err)
		}
		msg, err := payload.NewMessage(&payload.LightInitialState{
			JetIDs: IDs,
			Drops: drops,
			Pulse: pulse.ToProto(&topSyncPulse),
		})
		if err != nil {
			logger.Fatal("Couldn't make message", err)
		}
		p.dep.sender.Reply(ctx, p.meta, msg)
	} else if req.Pulse > startPulse {
		msg, err := payload.NewMessage(&payload.LightInitialState{})
		if err != nil {
			logger.Fatal("Couldn't make message", err)
		}
		p.dep.sender.Reply(ctx, p.meta, msg)
	} else {
		logger.Fatal("impossible situation")
	}
	return nil
}
