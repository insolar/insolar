package api

import (
	"context"
	"net/http"

	"github.com/insolar/insolar/core"
)

type ExporterArgs struct {
	From uint32
	Size int
}

type ExporterReply = core.ExportResult

type ExporterService struct {
	runner *Runner
}

func NewExporterService(runner *Runner) *ExporterService {
	return &ExporterService{runner: runner}
}

func (s *ExporterService) Export(r *http.Request, args *ExporterArgs, reply *ExporterReply) error {
	exp := s.runner.Exporter
	ctx := context.TODO()
	result, err := exp.Export(ctx, core.PulseNumber(args.From), args.Size)
	if err != nil {
		return err
	}

	reply.Data = result.Data
	reply.Size = result.Size
	reply.NextFrom = result.NextFrom

	return nil
}
