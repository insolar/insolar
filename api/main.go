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
	"crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/rpc/v2"
	jsonrpc "github.com/gorilla/rpc/v2/json2"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
)

const (
	_ int = 0
	// HandlerError is error in handler
	HandlerError int = -1
	// BadRequest is bad formed request
	BadRequest int = -2
)

func writeError(message string, code int) map[string]interface{} {
	errJSON := map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"code":    code,
		},
	}
	return errJSON
}

func makeHandlerMarshalErrorJSON(ctx context.Context) []byte {
	jsonErr := writeError("Invalid data from handler", HandlerError)
	serJSON, err := json.Marshal(jsonErr)
	if err != nil {
		inslogger.FromContext(ctx).Fatal("Can't marshal base error")
	}
	return serJSON
}

var handlerMarshalErrorJSON = makeHandlerMarshalErrorJSON(inslogger.ContextWithTrace(context.Background(), "handlerMarshalErrorJSON"))

func processQueryType(ctx context.Context, rh *RequestHandler, qTypeStr string) map[string]interface{} {
	qtype := QTypeFromString(qTypeStr)
	var answer map[string]interface{}

	var hError error
	switch qtype {
	case GetSeed:
		answer, hError = rh.ProcessGetSeed(ctx)
	default:
		msg := fmt.Sprintf("Wrong query parameter 'query_type' = '%s'", qTypeStr)
		answer = writeError(msg, BadRequest)
		inslogger.FromContext(ctx).Warnf("[ processQueryType ] %s\n", msg)
		return answer
	}
	if hError != nil {
		errMsg := "Handler error: " + hError.Error()
		inslogger.FromContext(ctx).Errorf("[ processQueryType ] %s\n", errMsg)
		answer = writeError(errMsg, HandlerError)
	}

	return answer
}

const traceIDQueryParam = "traceID"

// PreprocessRequest extracts params from requests
func PreprocessRequest(ctx context.Context, req *http.Request) (*Params, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.Wrap(err, "[ PreprocessRequest ] Can't read body. So strange")
	}
	if len(body) == 0 {
		return nil, errors.New("[ PreprocessRequest ] Empty body")
	}

	var params Params
	err = json.Unmarshal(body, &params)
	if err != nil {
		return nil, errors.Wrap(err, "[ PreprocessRequest ] Can't parse input params")
	}

	inslogger.FromContext(ctx).Infof("[ PreprocessRequest ] Query: %s. Url: %s\n", string(body), req.URL)

	return &params, nil
}

func wrapAPIV1Handler(runner *Runner, rootDomainReference core.RecordRef) func(w http.ResponseWriter, r *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {
		traceid := utils.RandTraceID()
		ctx, inslog := inslogger.WithTraceField(context.Background(), traceid)
		startTime := time.Now()
		answer := make(map[string]interface{})
		var params *Params
		defer func() {
			if answer == nil {
				answer = make(map[string]interface{})
			}
			if params == nil {
				params = &Params{}
			}
			answer[traceIDQueryParam] = traceid
			serJSON, err := json.MarshalIndent(answer, "", "    ")
			if err != nil {
				serJSON = handlerMarshalErrorJSON
			}
			response.Header().Add("Content-Type", "application/json")
			var newLine byte = '\n'
			_, err = response.Write(append(serJSON, newLine))
			if err != nil {
				inslog.Errorf("[ wrapAPIV1Handler ] Can't write response\n")
			}
			inslog.Infof("[ wrapAPIV1Handler ] Request completed. Total time: %s\n", time.Since(startTime))
		}()

		params, err := PreprocessRequest(ctx, req)
		if err != nil {
			answer = writeError("Bad request", BadRequest)
			inslog.Errorf("[ wrapAPIV1Handler ] Can't parse input request: %s, error: %s\n", req.RequestURI, err)
			return
		}
		rh := NewRequestHandler(params, runner.MessageBus, runner.NetworkCoordinator, rootDomainReference, runner.seedmanager)

		answer = processQueryType(ctx, rh, params.QueryType)
	}
}

// Runner implements Component for API
type Runner struct {
	MessageBus         core.MessageBus         `inject:""`
	Genesis            core.Genesis            `inject:""`
	NetworkCoordinator core.NetworkCoordinator `inject:""`
	server             *http.Server
	rpcServer          *rpc.Server
	cfg                *configuration.APIRunner
	keyCache           map[string]crypto.PublicKey
	cacheLock          *sync.RWMutex
	seedmanager        *seedmanager.SeedManager
	StorageExporter    core.StorageExporter `inject:""`
}

// NewRunner is C-tor for API Runner
func NewRunner(cfg *configuration.APIRunner) (*Runner, error) {
	if cfg == nil {
		return nil, errors.New("[ NewAPIRunner ] config is nil")
	}
	if cfg.Address == "" {
		return nil, errors.New("[ NewAPIRunner ] Address must not be empty")
	}
	if len(cfg.Location) == 0 {
		return nil, errors.New("[ NewAPIRunner ] Location must exist")
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
	err := rpcServer.RegisterService(NewStorageExporterService(&ar), "exporter")
	if err != nil {
		return nil, err
	}

	return &ar, nil
}

func (ar *Runner) IsAPIRunner() bool {
	return true
}

// Start runs api server
func (ar *Runner) Start(ctx context.Context) error {
	rootDomainReference := ar.Genesis.GetRootDomainRef()

	ar.seedmanager = seedmanager.New()

	fw := wrapAPIV1Handler(ar, *rootDomainReference)
	http.HandleFunc(ar.cfg.Location, fw)
	http.HandleFunc(ar.cfg.Info, ar.infoHandler())
	http.HandleFunc(ar.cfg.Call, ar.callHandler())
	http.Handle(ar.cfg.RPC, ar.rpcServer)
	inslog := inslogger.FromContext(ctx)
	inslog.Info("Starting ApiRunner ...")
	inslog.Info("Config: ", ar.cfg)
	go func() {
		if err := ar.server.ListenAndServe(); err != nil {
			inslog.Error("Httpserver: ListenAndServe() error: ", err)
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
	args, err := core.MarshalArgs()
	if err != nil {
		return nil, errors.Wrap(err, "Can't marshal empty args")
	}
	res, err := ar.MessageBus.Send(
		ctx,
		&message.CallMethod{
			ObjectRef: core.NewRefFromBase58(ref),
			Method:    "GetPublicKey",
			Arguments: args,
		},
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Can't get public key")
	}

	var publicKeyString string
	var contractErr error
	err = signer.UnmarshalParams(res.(*reply.CallMethod).Result, &publicKeyString, &contractErr)
	if err != nil {
		return nil, errors.Wrap(err, "Can't unmarshal public key")
	}
	if contractErr != nil {
		return nil, errors.Wrap(contractErr, "Error in get public key")
	}

	kp := platformpolicy.NewKeyProcessor()
	publicKey, err = kp.ImportPublicKey([]byte(publicKeyString))
	if err != nil {
		return nil, errors.Wrap(err, "Failed to convert public key")
	}

	ar.cacheLock.Lock()
	ar.keyCache[ref] = publicKey
	ar.cacheLock.Unlock()
	return publicKey, nil
}
