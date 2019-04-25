package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/pkg/errors"
)

type CheckJet func(ctx context.Context, target insolar.ID, pn insolar.PulseNumber) (insolar.JetID, bool, error)

func NewCheckJet(fetcher jet.Fetcher, coord jet.Coordinator) CheckJet {
	return func(ctx context.Context, target insolar.ID, pn insolar.PulseNumber) (insolar.JetID, bool, error) {
		// Special case for genesis pulse. No one was executor at that time, so anyone can fetch data from it.
		if pn <= insolar.FirstPulseNumber {
			return *insolar.NewJetID(0, nil), true, nil
		}

		jetID, err := fetcher.Fetch(ctx, target, pn)
		if err != nil {
			return insolar.JetID(*jetID), false, errors.Wrap(err, "failed to fetch jet")
		}
		executor, err := coord.LightExecutorForJet(ctx, *jetID, pn)
		if err != nil {
			return insolar.JetID(*jetID), false, errors.Wrap(err, "failed to calculate executor for jet")
		}
		if *executor != coord.Me() {
			return insolar.JetID(*jetID), false, nil
		}

		return insolar.JetID(*jetID), true, nil
	}
}
