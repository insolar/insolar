//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/instrumenter"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/rpc/v2"
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
func (s *NodeCertService) get(ctx context.Context, _ *http.Request, args *NodeCertArgs, _ *rpc.RequestBody, reply *NodeCertReply) error {
	nodeRef, err := insolar.NewReferenceFromBase58(args.Ref)
	if err != nil {
		return errors.Wrap(err, "failed to parse args.Ref")
	}
	cert, err := s.runner.CertificateGetter.GetCert(ctx, nodeRef)
	if err != nil {
		return errors.Wrap(err, "failed to get certificate")
	}

	reply.Cert = cert.(*certificate.Certificate)
	return nil
}

func (s *NodeCertService) Get(r *http.Request, args *NodeCertArgs, requestBody *rpc.RequestBody, reply *NodeCertReply) error {
	ctx, instr := instrumenter.NewMethodInstrument("NodeCertService.get")
	defer instr.End()

	msg := fmt.Sprint("Incoming request: ", r.RequestURI)
	instr.Annotate(msg)

	logger := inslogger.FromContext(ctx)
	logger.Info("[ NodeCertService.get ] ", msg)

	err := s.get(ctx, r, args, requestBody, reply)
	if err != nil {
		instr.SetError(err, InternalErrorShort)
		return errors.Wrap(err, "failed to execute NodeCertService.get")
	}

	return err
}
