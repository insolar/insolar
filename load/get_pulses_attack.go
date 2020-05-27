package load

import (
	"context"
	"io"

	"github.com/spf13/viper"

	"github.com/insolar/insolar/ledger/heavy/exporter"
	"github.com/insolar/insolar/log"
	"github.com/skudasov/loadgen"
)

type GetPulsesAttack struct {
	loadgen.WithRunner
	Exporter
}

func (a *GetPulsesAttack) Setup(hc loadgen.RunnerConfig) error {
	a.Exporter = NewExporter(viper.GetString("generator.target"))
	return nil
}

func (a *GetPulsesAttack) Do(ctx context.Context) loadgen.DoResult {
	tsp, err := GetTopSyncPulse(ctx, a.pulseClient)
	if err != nil {
		return loadgen.DoResult{
			Error:        err,
			RequestLabel: GetPulsesLabel,
		}
	}
	a.R.L.Infof("top sync pulse: %d", tsp)
	// as observer we are always fetching the last pulse
	request := &exporter.GetPulses{Count: 1, PulseNumber: tsp - 10}
	stream, err := a.pulseClient.Export(ctx, request)
	if err != nil {
		log.Infof("failed to create stream")
		return loadgen.DoResult{
			Error:        err,
			RequestLabel: GetPulsesLabel,
		}
	}

	for {
		resp, err := stream.Recv()
		a.R.L.Infof("resp: %s, err: %s\n", resp, err)
		if err == io.EOF {
			a.R.L.Info("EOF received")
			break
		}
		if err != nil {
			return loadgen.DoResult{
				Error:        err,
				RequestLabel: GetPulsesLabel,
			}
		}

	}
	return loadgen.DoResult{
		Error:        nil,
		RequestLabel: GetPulsesLabel,
	}
}

func (a *GetPulsesAttack) Clone(r *loadgen.Runner) loadgen.Attack {
	return &GetPulsesAttack{WithRunner: loadgen.WithRunner{R: r}}
}
