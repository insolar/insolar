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

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// NodeCertArgs is arguments that NodeCert service accepts.
type NodeCertArgs struct {
	Ref string
}

// NodeCertReply is reply for NodeCert service requests.
type NodeCertReply struct {
	Cert *certificate.Certificate `json:"cert"`
}

// NodeCertService is a service that provides cert for node.
type NodeCertService struct {
	runner *Runner
}

// NewNodeCertService creates new NodeCert service instance.
func NewNodeCertService(runner *Runner) *NodeCertService {
	return &NodeCertService{runner: runner}
}

// Get returns certificate for node with given reference.
func (s *NodeCertService) Get(r *http.Request, args *NodeCertArgs, reply *NodeCertReply) error {
	ctx, inslog := inslogger.WithTraceField(context.Background(), utils.RandTraceID())

	inslog.Infof("[ NodeCertService.Get ] Incoming request: %s", r.RequestURI)

	cert, err := s.runner.NetworkCoordinator.GetCert(ctx, core.NewRefFromBase58(args.Ref))
	if err != nil {
		return errors.Wrap(err, "[ NodeCertService.Get] ")
	}

	reply.Cert = cert.(*certificate.Certificate)
	return nil
}
