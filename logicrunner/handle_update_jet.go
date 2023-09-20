package logicrunner

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
)

type HandleUpdateJet struct {
	dep *Dependencies

	meta payload.Meta
}

func (h *HandleUpdateJet) Present(ctx context.Context, _ flow.Flow) error {
	pl := payload.UpdateJet{}
	err := pl.Unmarshal(h.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal payload")
	}

	err = h.dep.JetStorage.Update(ctx, pl.Pulse, true, pl.JetID)
	if err != nil {
		return errors.Wrap(err, "failed to update jets")
	}

	return nil
}
