package load

import (
	"context"

	"github.com/insolar/insolar/log"
	"google.golang.org/grpc"

	"github.com/insolar/insolar/insolar"

	"github.com/insolar/insolar/ledger/heavy/exporter"
)

type Exporter struct {
	recordClient exporter.RecordExporterClient
	pulseClient  exporter.PulseExporterClient
}

func NewExporter(target string) Exporter {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	return Exporter{
		pulseClient:  exporter.NewPulseExporterClient(conn),
		recordClient: exporter.NewRecordExporterClient(conn),
	}
}

func GetTopSyncPulse(ctx context.Context, a exporter.PulseExporterClient) (insolar.PulseNumber, error) {
	request := &exporter.GetTopSyncPulse{}
	res, err := a.TopSyncPulse(ctx, request)
	if err != nil {
		return 0, err
	}
	return insolar.PulseNumber(res.PulseNumber), nil
}
