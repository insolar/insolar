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
	"github.com/pkg/errors"
)

// StorageExporterArgs is arguments that StorageExporter service accepts.
type StorageExporterArgs struct {
	From uint32
	Size int
}

// StorageExporterReply is reply for StorageExporter service requests.
type StorageExporterReply = core.StorageExportResult

// StorageExporterService is a service that provides API for exporting storage data.
type StorageExporterService struct {
	runner *Runner
}

// NewStorageExporterService creates new StorageExporter service instance.
func NewStorageExporterService(runner *Runner) *StorageExporterService {
	return &StorageExporterService{runner: runner}
}

// Export returns data view from storage.
//
//   Request structure:
//   {
//     "jsonrpc": "2.0",
//     "method": "exporter.Export",
//     "params": {
//       // Pulse number from which data load should start.
//       // If less than first pulse, the load will start from the first pulse (e.i. use "0" to load from the beginning).
//       "From": int,
//       // Number of pulses to load.
//       "Size": int
//       },
//     "id": str|int|null
//   }
//
//   Response structure:
//   {
//     "Data": {
//       [pulse number]: {
//         "Records": {
//           [record ID]: {
//             "Type": str, // Constant record type.
//             "Data": { ... }, // Structured record data.
//             "Payload": { ... }|null // Additional data related to the record (e.g. Object's memory).
//         }
//       },
//         "Pulse": {
//           "PulseNumber": int, // Pulse number. Same as parent key.
//         }
//       },
//       "NextFrom": int|null, // Pulse number from which to start next batch. Put it as "From" param for next incremental fetch.
//       "Size": int // Number of returned pulses (length of the "Data" dictionary).
//   }
//
func (s *StorageExporterService) Export(r *http.Request, args *StorageExporterArgs, reply *StorageExporterReply) error {
	exp := s.runner.StorageExporter
	ctx := context.TODO()
	result, err := exp.Export(ctx, core.PulseNumber(args.From), args.Size)
	if err != nil {
		return errors.Wrap(err, "[ Export ]")
	}

	reply.Data = result.Data
	reply.Size = result.Size
	reply.NextFrom = result.NextFrom

	return nil
}
