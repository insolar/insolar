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

package stubs

import (
	"context"
	"crypto"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/insolar/rpc/v2"
	jsonrpc "github.com/insolar/rpc/v2/json2"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/hostnetwork/resolver"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/version"
)

// NewRunner is C-tor for API Runner.
func NewRunner(cfg *configuration.APIRunner, cfgTransport *configuration.Transport, cert insolar.Certificate) (insolar.APIRunner, error) {

	if err := checkConfig(cfg); err != nil {
		return nil, errors.Wrap(err, "[ NewAPIRunner ] Bad APIRunner config")
	}

	if cfgTransport == nil {
		return nil, errors.New("[ NewAPIRunner ] Bad Transport config")
	}

	if cert == nil {
		return nil, errors.New("[ NewAPIRunner ] Bad Certificate")
	}

	addrStr := fmt.Sprint(cfg.Address)
	rpcServer := rpc.NewServer()
	ar := apiRunnerStub{
		server:    &http.Server{Addr: addrStr},
		rpcServer: rpcServer,
		cfg:       cfg,
		timeout:   30 * time.Second,
		keyCache:  make(map[string]crypto.PublicKey),
		cacheLock: &sync.RWMutex{},
	}

	rpcServer.RegisterCodec(jsonrpc.NewCodec(), "application/json")
	err := rpcServer.RegisterService(&nodeServiceStub{configuration: cfgTransport, certificate: cert}, "node")
	if err != nil {
		return nil, errors.Wrap(err, "[ registerServices ] Can't RegisterService: node")
	}

	return &ar, nil
}

type apiRunnerStub struct {
	server    *http.Server
	rpcServer *rpc.Server
	cfg       *configuration.APIRunner
	keyCache  map[string]crypto.PublicKey
	cacheLock *sync.RWMutex
	timeout   time.Duration
}

// IsAPIRunner is implementation of APIRunner interface for component manager.
func (ar *apiRunnerStub) IsAPIRunner() bool {
	return true
}

// Start runs api server.
func (ar *apiRunnerStub) Start(ctx context.Context) error {

	router := http.NewServeMux()
	ar.server.Handler = router

	router.HandleFunc("/healthcheck", ar.healthCheck)
	router.Handle(ar.cfg.RPC, ar.rpcServer)

	inslog := inslogger.FromContext(ctx)
	inslog.Info("Starting Observer ApiRunner ...")
	inslog.Info("Config: ", ar.cfg)
	listener, err := net.Listen("tcp", ar.server.Addr)
	if err != nil {
		return errors.Wrap(err, "Can't start listening")
	}
	go func() {
		if err := ar.server.Serve(listener); err != http.ErrServerClosed {
			inslog.Error("Http server: ListenAndServe() error: ", err)
		}
	}()
	return nil
}

// Stop stops api server.
func (ar *apiRunnerStub) Stop(ctx context.Context) error {
	const timeOut = 5

	inslogger.FromContext(ctx).Infof("Shutting down server gracefully ...(waiting for %d seconds)", timeOut)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeOut)*time.Second)
	defer cancel()
	err := ar.server.Shutdown(ctxWithTimeout)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop API server")
	}

	return nil
}

func checkConfig(cfg *configuration.APIRunner) error {
	if cfg == nil {
		return errors.New("[ checkConfig ] config is nil")
	}
	if cfg.Address == "" {
		return errors.New("[ checkConfig ] Address must not be empty")
	}
	if len(cfg.Call) == 0 {
		return errors.New("[ checkConfig ] Call must exist")
	}
	if len(cfg.RPC) == 0 {
		return errors.New("[ checkConfig ] RPC must exist")
	}

	return nil
}

func (ar *apiRunnerStub) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

type nodeServiceStub struct {
	configuration *configuration.Transport
	certificate   insolar.Certificate
}

// Get returns status info.
func (s *nodeServiceStub) GetStatus(r *http.Request, args *interface{}, reply *api.StatusReply) error {
	traceID := utils.RandTraceID()
	_, inslog := inslogger.WithTraceField(context.Background(), traceID)

	inslog.Infof("[ NodeService.GetStatus ] Incoming request: %s", r.RequestURI)

	reply.ActiveListSize = 0
	reply.WorkingListSize = 0

	addr, err := net.ResolveTCPAddr("tcp", s.configuration.Address)
	if err != nil {
		return errors.Wrap(err, "Failed to resolve public address")
	}
	publicAddress, err := resolver.Resolve(s.configuration.FixedPublicAddress, addr.String())
	if err != nil {
		return errors.Wrap(err, "Failed to resolve public address")
	}

	role := s.certificate.GetRole()
	if role == insolar.StaticRoleUnknown {
		log.Info("[ createOrigin ] Use insolar.StaticRoleLightMaterial, since no role in certificate")
		role = insolar.StaticRoleLightMaterial
	}

	origin := node.NewNode(
		*s.certificate.GetNodeRef(),
		role,
		s.certificate.GetPublicKey(),
		publicAddress,
		version.Version,
	)

	reply.Origin = api.Node{
		Reference: origin.ID().String(),
		Role:      origin.Role().String(),
		IsWorking: origin.GetState() == insolar.NodeReady,
	}

	reply.PulseNumber = uint32(insolar.GenesisPulse.PulseNumber)
	reply.Version = version.Version

	reply.NetworkState = insolar.CompleteNetworkState.String()
	reply.NodeState = origin.GetState().String()

	return nil
}
