// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package handle

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/proc"
)

type HotObjects struct {
	dep  *proc.Dependencies
	meta payload.Meta
}

func NewHotObjects(dep *proc.Dependencies, meta payload.Meta) *HotObjects {
	return &HotObjects{
		dep:  dep,
		meta: meta,
	}
}

func (s *HotObjects) Present(ctx context.Context, f flow.Flow) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("start hotObjects msg processing")

	msg, err := payload.Unmarshal(s.meta.Payload)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal message")
	}

	hots, ok := msg.(*payload.HotObjects)
	if !ok {
		return errors.New("received wrong message")
	}
	logger = logger.WithFields(map[string]interface{}{
		"pulse":  hots.Pulse,
		"jet_id": hots.JetID.DebugString(),
	})

	notificationLimit := s.dep.Config().MaxNotificationsPerPulse
	hdProc := proc.NewHotObjects(s.meta, hots.Pulse, hots.JetID, hots.Drop, hots.Indexes, notificationLimit)
	s.dep.HotObjects(hdProc)
	if err := f.Procedure(ctx, hdProc, false); err != nil {
		return err
	}

	logger.Info("finish hotObjects msg processing")
	return nil
}
