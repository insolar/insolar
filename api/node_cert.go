// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"github.com/insolar/rpc/v2"

	"github.com/insolar/insolar/api/instrumenter"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
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
	nodeRef, err := insolar.NewReferenceFromString(args.Ref)
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
