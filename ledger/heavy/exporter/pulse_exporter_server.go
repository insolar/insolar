// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package exporter

import (
	"context"
	"time"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/node"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/insmetrics"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/pulse"
)

type PulseServer struct {
	pulses    insolarPulse.Calculator
	jetKeeper executor.JetKeeper
	nodes     node.Accessor
}

func NewPulseServer(pulses insolarPulse.Calculator, jetKeeper executor.JetKeeper, nodeAccessor node.Accessor) *PulseServer {
	return &PulseServer{
		pulses:    pulses,
		jetKeeper: jetKeeper,
		nodes:     nodeAccessor,
	}
}

func (p *PulseServer) Export(getPulses *GetPulses, stream PulseExporter_ExportServer) error {
	ctx := stream.Context()

	exportStart := time.Now()
	defer func(ctx context.Context) {
		stats.Record(
			insmetrics.InsertTag(ctx, TagHeavyExporterMethodName, "pulse-export"),
			HeavyExporterMethodTiming.M(float64(time.Since(exportStart).Nanoseconds())/1e6),
		)
	}(ctx)

	logger := inslogger.FromContext(ctx)

	if getPulses.Count == 0 {
		return ErrNilCount
	}

	read := uint32(0)
	if getPulses.PulseNumber == 0 {
		getPulses.PulseNumber = pulse.MinTimePulse
		err := stream.Send(&Pulse{
			PulseNumber:    pulse.MinTimePulse,
			Entropy:        insolar.GenesisPulse.Entropy,
			PulseTimestamp: insolar.GenesisPulse.PulseTimestamp,
		})
		if err != nil {
			logger.Error(err)
			return err
		}
		read++
	}
	currentPN := getPulses.PulseNumber
	for read < getPulses.Count {
		topPulse := p.jetKeeper.TopSyncPulse()
		if currentPN >= topPulse {
			return nil
		}

		pulse, err := p.pulses.Forwards(ctx, currentPN, 1)
		if err != nil {
			logger.Error(err)
			return err
		}
		nodes, err := p.nodes.All(pulse.PulseNumber)
		if err != nil {
			logger.Error(err)
			return err
		}
		err = stream.Send(&Pulse{
			PulseNumber:    pulse.PulseNumber,
			Entropy:        pulse.Entropy,
			PulseTimestamp: pulse.PulseTimestamp,
			Nodes:          nodes,
		})
		if err != nil {
			logger.Error(err)
			return err
		}

		read++
		currentPN = pulse.PulseNumber
	}

	return nil
}

func (p *PulseServer) TopSyncPulse(ctx context.Context, _ *GetTopSyncPulse) (*TopSyncPulseResponse, error) {
	exportStart := time.Now()
	defer func(ctx context.Context) {
		stats.Record(
			insmetrics.InsertTag(ctx, TagHeavyExporterMethodName, "pulse-top-sync-pulse"),
			HeavyExporterMethodTiming.M(float64(time.Since(exportStart).Nanoseconds())/1e6),
		)
	}(ctx)

	return &TopSyncPulseResponse{
		PulseNumber: p.jetKeeper.TopSyncPulse().AsUint32(),
	}, nil
}
