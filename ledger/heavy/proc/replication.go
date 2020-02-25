// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package proc

import (
	"context"
	"fmt"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
)

type Replication struct {
	message payload.Meta
	cfg     configuration.Ledger

	dep struct {
		replicator executor.HeavyReplicator
	}
}

func NewReplication(msg payload.Meta, cfg configuration.Ledger) *Replication {
	return &Replication{
		message: msg,
		cfg:     cfg,
	}
}

func (p *Replication) Dep(
	replicator executor.HeavyReplicator,
) {
	p.dep.replicator = replicator
}

func (p *Replication) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx)
	logger.Info("got replication msg")

	pl, err := payload.Unmarshal(p.message.Payload)
	if err != nil {
		logger.Error(err)
		return errors.Wrap(err, "failed to unmarshal payload")
	}
	msg, ok := pl.(*payload.Replication)
	if !ok {
		logger.Error(err)
		return fmt.Errorf("unexpected payload %T", pl)
	}

	logger.Debugf("notify heavy replicator about jetID:%v, pn:%v", msg.JetID.DebugString(), msg.Pulse)
	go p.dep.replicator.NotifyAboutMessage(ctx, msg)

	stats.Record(ctx, statReceivedHeavyPayloadCount.M(1))

	logger.Info("finish replication msg processing")

	return nil
}
