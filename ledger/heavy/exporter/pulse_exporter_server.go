// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package exporter

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go.opencensus.io/stats"
	"golang.org/x/crypto/sha3"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/node"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/pulse"
)

type PulseServer struct {
	pulses    insolarPulse.Calculator
	jetKeeper executor.JetKeeper
	nodes     node.Accessor
	authCfg   configuration.Auth
}

func NewPulseServer(
	pulses insolarPulse.Calculator,
	jetKeeper executor.JetKeeper,
	nodeAccessor node.Accessor,
	authCfg configuration.Auth,
) *PulseServer {
	return &PulseServer{
		pulses:    pulses,
		jetKeeper: jetKeeper,
		nodes:     nodeAccessor,
		authCfg:   authCfg,
	}
}

func (p *PulseServer) Export(getPulses *GetPulses, stream PulseExporter_ExportServer) error {
	ctx := stream.Context()
	ctxWithTags := addTagsForExporterMethodTiming(p.authCfg.Required, ctx, "pulse-export")

	exportStart := time.Now()
	defer func(ctx context.Context) {
		stats.Record(
			ctxWithTags,
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

	stats.Record(
		ctxWithTags,
		HeavyExporterLastExportedPulse.M(int64(currentPN)),
	)
	return nil
}

func (p *PulseServer) TopSyncPulse(ctx context.Context, _ *GetTopSyncPulse) (*TopSyncPulseResponse, error) {
	exportStart := time.Now()
	defer func(ctx context.Context) {
		stats.Record(
			addTagsForExporterMethodTiming(p.authCfg.Required, ctx, "pulse-top-sync-pulse"),
			HeavyExporterMethodTiming.M(float64(time.Since(exportStart).Nanoseconds())/1e6),
		)
	}(ctx)

	return &TopSyncPulseResponse{
		PulseNumber: p.jetKeeper.TopSyncPulse().AsUint32(),
	}, nil
}

func (p *PulseServer) NextFinalizedPulse(ctx context.Context, gnfp *GetNextFinalizedPulse) (*FullPulse, error) {
	pn := gnfp.GetPulseNo()
	logger := inslogger.FromContext(ctx)

	if pn == 0 {
		pu, err := p.pulses.Forwards(ctx, p.jetKeeper.TopSyncPulse(), 0)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		return p.makeFullPulse(ctx, pu, p.jetKeeper.Storage())
	} else if pn < pulse.MinTimePulse {
		pu, err := p.pulses.Forwards(ctx, pulse.MinTimePulse, 0)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		return p.makeFullPulse(ctx, pu, p.jetKeeper.Storage())
	}

	pu, err := p.pulses.Forwards(ctx, insolar.PulseNumber(pn), 1)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	pu.PrevPulseNumber = insolar.PulseNumber(pn)

	return p.makeFullPulse(ctx, pu, p.jetKeeper.Storage())
}

type JetData struct {
	PulseNumber insolar.PulseNumber
	JetID       string
}

// JetIDToString returns the string representation of JetID
func JetIDToString(id insolar.JetID) string {
	depth, prefix := id.Depth(), id.Prefix()
	if depth == 0 {
		return ""
	}
	res := strings.Builder{}
	for i := uint8(0); i < depth; i++ {
		bytePos, bitPos := i/8, 7-i%8

		byteValue := prefix[bytePos]
		bitValue := byteValue >> uint(bitPos) & 0x01
		bitString := strconv.Itoa(int(bitValue))
		res.WriteString(bitString)
	}
	return res.String()
}

func (p *PulseServer) makeFullPulse(ctx context.Context, pu insolar.Pulse, js jet.Storage) (*FullPulse, error) {
	logger := inslogger.FromContext(ctx)
	jets := js.All(ctx, pu.PulseNumber)
	var res []JetDropContinue
	prevJetDrops := js.All(ctx, pu.PrevPulseNumber)
	for _, j := range jets {
		rawData, err := json.Marshal(JetData{PulseNumber: pu.PulseNumber, JetID: JetIDToString(j)})
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		hash := sha3.Sum224(rawData)

		var prevDropHashes [][]byte
		jetIDLeft, jetIDRight := jet.Siblings(j)
		for _, prevJetDrop := range prevJetDrops {
			if prevJetDrop == j || prevJetDrop == jet.Parent(j) || prevJetDrop == jetIDLeft || prevJetDrop == jetIDRight {
				rawData, err := json.Marshal(JetData{PulseNumber: pu.PrevPulseNumber, JetID: JetIDToString(prevJetDrop)})
				if err != nil {
					logger.Error(err)
					return nil, err
				}
				prevDropHash := sha3.Sum224(rawData)
				prevDropHashes = append(prevDropHashes, prevDropHash[:])
			}
		}

		res = append(res, JetDropContinue{
			JetID:          j,
			Hash:           hash[:],
			PrevDropHashes: prevDropHashes,
		})
	}
	return &FullPulse{
		PulseNumber:      pu.PulseNumber,
		PrevPulseNumber:  pu.PrevPulseNumber,
		NextPulseNumber:  pu.NextPulseNumber,
		Entropy:          pu.Entropy,
		PulseTimestamp:   pu.PulseTimestamp,
		EpochPulseNumber: pu.EpochPulseNumber,
		Jets:             res,
	}, nil
}
