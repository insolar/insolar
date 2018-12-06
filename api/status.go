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

	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// StatusReply is reply for Status service requests.
type StatusReply struct {
	NetworkState string
}

// StatusService is a service that provides API for getting status of node.
type StatusService struct {
	runner *Runner
}

// NewStatusService creates new StatusService instance.
func NewStatusService(runner *Runner) *StatusService {
	return &StatusService{runner: runner}
}

// Get returns status info
func (s *StatusService) Get(r *http.Request, args *interface{}, reply *StatusReply) error {
	traceID := utils.RandTraceID()
	_, inslog := inslogger.WithTraceField(context.Background(), traceID)

	inslog.Infof("[ StatusService.Get ] Incoming request: %s", r.RequestURI)

	reply.NetworkState = s.runner.NetworkSwitcher.GetState().String()

	return nil
}
