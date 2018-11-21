/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

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
