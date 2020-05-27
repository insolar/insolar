package load

import (
	"context"
	"io"

	"github.com/spf13/viper"

	"github.com/insolar/insolar/log"

	"github.com/insolar/insolar/ledger/heavy/exporter"
	"github.com/insolar/loadgen"
)

type GetRecordsAttack struct {
	loadgen.WithRunner
	Exporter
}

func (a *GetRecordsAttack) Setup(hc loadgen.RunnerConfig) error {
	a.Exporter = NewExporter(viper.GetString("generator.target"))
	return nil
}

func (a *GetRecordsAttack) Do(ctx context.Context) loadgen.DoResult {
	tsp, err := GetTopSyncPulse(ctx, a.pulseClient)
	if err != nil {
		return loadgen.DoResult{
			RequestLabel: GetRecordsLabel,
			Error:        err,
		}
	}
	request := &exporter.GetRecords{Count: 2000, PulseNumber: tsp, RecordNumber: 0}
	stream, err := a.recordClient.Export(ctx, request)
	if err != nil {
		a.R.L.Infof("failed to create stream")
		return loadgen.DoResult{
			Error:        err,
			RequestLabel: GetPulsesLabel,
		}
	}

	for {
		log.Infof("iterating on pulse: %d", tsp)
		resp, err := stream.Recv()
		if resp != nil {
			a.R.L.Infof("record fetched: %s", resp.RecordNumber)
		}
		if err == io.EOF {
			a.R.L.Info("EOF received")
			break
		}
		if err != nil {
			a.R.L.Infof("error: %s", err)
			return loadgen.DoResult{
				RequestLabel: GetRecordsLabel,
				Error:        err,
			}
		}
		a.R.L.Infof("pulse from msg: %d", resp.Record.ID.Pulse())
		if resp.Record.ID.Pulse() != tsp {
			a.R.L.Infof("next pulse, skipping")
			break
		}
	}
	return loadgen.DoResult{
		RequestLabel: GetRecordsLabel,
		Error:        err,
	}
}

func (a *GetRecordsAttack) Clone(r *loadgen.Runner) loadgen.Attack {
	return &GetRecordsAttack{WithRunner: loadgen.WithRunner{R: r}}
}
