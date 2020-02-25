// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package adapters

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/gcpv2/api"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

func ConsensusContext(ctx context.Context) context.Context {
	return inslogger.UpdateLogger(ctx, func(logger insolar.Logger) (insolar.Logger, error) {
		return logger.Copy().WithFields(map[string]interface{}{
			"component":  "consensus",
			"LowLatency": true,
		}).WithMetrics(insolar.LogMetricsWriteDelayField).BuildLowLatency()
	})
}

func PacketEarlyLogger(ctx context.Context, senderAddr string) (context.Context, insolar.Logger) {
	ctx = ConsensusContext(ctx)

	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"sender_address": senderAddr,
	})

	return ctx, logger
}

func PacketLateLogger(ctx context.Context, parser transport.PacketParser) (context.Context, insolar.Logger) {
	ctx, logger := inslogger.WithFields(ctx, map[string]interface{}{
		"sender_id":    parser.GetSourceID(),
		"packet_type":  parser.GetPacketType().String(),
		"packet_pulse": parser.GetPulseNumber(),
	})

	return ctx, logger
}

func ReportContext(report api.UpstreamReport) context.Context {
	return network.NewPulseContext(context.Background(), uint32(report.PulseNumber))
}
