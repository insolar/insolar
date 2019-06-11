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
	"crypto"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/insolar/insolar/network"

	"github.com/gorilla/rpc/v2"
	jsonrpc "github.com/gorilla/rpc/v2/json2"
	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/seedmanager"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/platformpolicy"
)

// Runner implements Component for API
type Runner struct {
	CertificateManager  insolar.CertificateManager  `inject:""`
	ContractRequester   insolar.ContractRequester   `inject:""`
	GenesisDataProvider insolar.GenesisDataProvider `inject:""`
	NodeNetwork         insolar.NodeNetwork         `inject:""`
	Gatewayer           network.Gatewayer           `inject:""`
	ServiceNetwork      insolar.Network             `inject:""`
	PulseAccessor       pulse.Accessor              `inject:""`
	ArtifactManager     artifacts.Client            `inject:""`
	JetCoordinator      jet.Coordinator             `inject:""`
	server              *http.Server
	rpcServer           *rpc.Server
	cfg                 *configuration.APIRunner
	keyCache            map[string]crypto.PublicKey
	cacheLock           *sync.RWMutex
	SeedManager         *seedmanager.SeedManager
	SeedGenerator       seedmanager.SeedGenerator
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
	if cfg.Timeout == 0 {
		return errors.New("[ checkConfig ] Timeout must not be null")
	}

	return nil
}

func (ar *Runner) registerServices(rpcServer *rpc.Server) error {
	err := rpcServer.RegisterService(NewSeedService(ar), "seed")
	if err != nil {
		return errors.Wrap(err, "[ registerServices ] Can't RegisterService: seed")
	}

	err = rpcServer.RegisterService(NewInfoService(ar), "info")
	if err != nil {
		return errors.Wrap(err, "[ registerServices ] Can't RegisterService: info")
	}

	err = rpcServer.RegisterService(NewStatusService(ar), "status")
	if err != nil {
		return errors.Wrap(err, "[ registerServices ] Can't RegisterService: status")
	}

	err = rpcServer.RegisterService(NewNodeCertService(ar), "cert")
	if err != nil {
		return errors.Wrap(err, "[ registerServices ] Can't RegisterService: cert")
	}

	err = rpcServer.RegisterService(NewContractService(ar, ar.cfg.EnableContractUploader), "contract")
	if err != nil {
		return errors.Wrap(err, "[ registerServices ] Can't RegisterService: contract")
	}

	return nil
}

// NewRunner is C-tor for API Runner
func NewRunner(cfg *configuration.APIRunner) (*Runner, error) {

	if err := checkConfig(cfg); err != nil {
		return nil, errors.Wrap(err, "[ NewAPIRunner ] Bad config")
	}

	addrStr := fmt.Sprint(cfg.Address)
	rpcServer := rpc.NewServer()
	ar := Runner{
		server:    &http.Server{Addr: addrStr},
		rpcServer: rpcServer,
		cfg:       cfg,
		keyCache:  make(map[string]crypto.PublicKey),
		cacheLock: &sync.RWMutex{},
	}

	rpcServer.RegisterCodec(jsonrpc.NewCodec(), "application/json")

	if err := ar.registerServices(rpcServer); err != nil {
		return nil, errors.Wrap(err, "[ NewAPIRunner ] Can't register services:")
	}

	return &ar, nil
}

// IsAPIRunner is implementation of APIRunner interface for component manager
func (ar *Runner) IsAPIRunner() bool {
	return true
}

// Start runs api server
func (ar *Runner) Start(ctx context.Context) error {
	hc := NewHealthChecker(ar.CertificateManager, ar.NodeNetwork)
	http.HandleFunc("/healthcheck", hc.CheckHandler)
	ar.SeedManager = seedmanager.New()
	http.HandleFunc(ar.cfg.Call, ar.callHandler())
	http.Handle(ar.cfg.RPC, ar.rpcServer)
	inslog := inslogger.FromContext(ctx)
	inslog.Info("Starting ApiRunner ...")
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

// Stop stops api server
func (ar *Runner) Stop(ctx context.Context) error {
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

func (ar *Runner) getMemberPubKey(ctx context.Context, ref string) (crypto.PublicKey, error) { //nolint
	ar.cacheLock.RLock()
	publicKey, ok := ar.keyCache[ref]
	ar.cacheLock.RUnlock()
	if ok {
		return publicKey, nil
	}

	reference, err := insolar.NewReferenceFromBase58(ref)
	if err != nil {
		return nil, errors.Wrap(err, "[ getMemberPubKey ] Can't parse ref")
	}
	res, err := ar.ContractRequester.SendRequest(ctx, reference, "GetPublicKey", []interface{}{})
	if err != nil {
		return nil, errors.Wrap(err, "[ getMemberPubKey ] Can't get public key")
	}

	publicKeyString, err := extractor.PublicKeyResponse(res.(*reply.CallMethod).Result)
	if err != nil {
		return nil, errors.Wrap(err, "[ getMemberPubKey ] Can't extract response")
	}

	kp := platformpolicy.NewKeyProcessor()
	publicKey, err = kp.ImportPublicKeyPEM([]byte(publicKeyString))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert public key")
	}

	ar.cacheLock.Lock()
	ar.keyCache[ref] = publicKey
	ar.cacheLock.Unlock()
	return publicKey, nil
}
