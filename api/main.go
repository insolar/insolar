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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/inscontext"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
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

func makeHandlerMarshalErrorJSON(ctx core.Context) []byte {
	jsonErr := writeError("Invalid data from handler", HandlerError)
	serJSON, err := json.Marshal(jsonErr)
	if err != nil {
		ctx.Log().Fatal("Can't marshal base error")
	}
	return serJSON
}

var handlerMarshalErrorJSON = makeHandlerMarshalErrorJSON(inscontext.WithTraceID(context.Background(), "handlerMarshalErrorJSON"))

func processQueryType(ctx core.Context, rh *RequestHandler, qTypeStr string) map[string]interface{} {
	qtype := QTypeFromString(qTypeStr)
	var answer map[string]interface{}

	var hError error
	switch qtype {
	case RegisterNode:
		answer, hError = rh.ProcessRegisterNode(ctx)
	case IsAuth:
		answer, hError = rh.ProcessIsAuthorized(ctx)
	case GetSeed:
		answer, hError = rh.ProcessGetSeed(ctx)
	default:
		msg := fmt.Sprintf("Wrong query parameter 'query_type' = '%s'", qTypeStr)
		answer = writeError(msg, BadRequest)
		ctx.Log().Warnf("[ processQueryType ] %s\n", msg)
		return answer
	}
	if hError != nil {
		errMsg := "Handler error: " + hError.Error()
		ctx.Log().Errorf("[ processQueryType ] %s\n", errMsg)
		answer = writeError(errMsg, HandlerError)
	}

	return answer
}

const TraceIDQueryParam = "traceID"

// PreprocessRequest extracts params from requests
func PreprocessRequest(ctx core.Context, req *http.Request) (*Params, error) {
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

	ctx.Log().Infof("[ PreprocessRequest ] Query: %s. Url: %s\n", string(body), req.URL)

	return &params, nil
}

func wrapAPIV1Handler(runner *Runner, rootDomainReference core.RecordRef) func(w http.ResponseWriter, r *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {
		ctx := inscontext.WithTraceID(inscontext.Background(), RandTraceID())

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
			answer[TraceIDQueryParam] = ctx.TraceID()
			serJSON, err := json.MarshalIndent(answer, "", "    ")
			if err != nil {
				serJSON = handlerMarshalErrorJSON
			}
			response.Header().Add("Content-Type", "application/json")
			var newLine byte = '\n'
			_, err = response.Write(append(serJSON, newLine))
			if err != nil {
				ctx.Log().Errorf("[ wrapAPIV1Handler ] Can't write response\n")
			}
			ctx.Log().Infof("[ wrapAPIV1Handler ] Request completed. Total time: %s\n", time.Since(startTime))
		}()

		params, err := PreprocessRequest(ctx, req)
		if err != nil {
			answer = writeError("Bad request", BadRequest)
			ctx.Log().Errorf("[ wrapAPIV1Handler ] Can't parse input request: %s, error: %s\n", req.RequestURI, err)
			return
		}
		rh := NewRequestHandler(params, runner.messageBus, runner.netCoordinator, rootDomainReference, runner.seedmanager)

		answer = processQueryType(ctx, rh, params.QueryType)
	}
}

// Runner implements Component for API
type Runner struct {
	messageBus     core.MessageBus
	server         *http.Server
	cfg            *configuration.APIRunner
	netCoordinator core.NetworkCoordinator
	keyCache       map[string]string
	cacheLock      *sync.RWMutex
	seedmanager    *seedmanager.SeedManager
}

// NewRunner is C-tor for API Runner
func NewRunner(cfg *configuration.APIRunner) (*Runner, error) {
	if cfg == nil {
		return nil, errors.New("[ NewAPIRunner ] config is nil")
	}
	if cfg.Port == 0 {
		return nil, errors.New("[ NewAPIRunner ] Port must not be 0")
	}
	if len(cfg.Location) == 0 {
		return nil, errors.New("[ NewAPIRunner ] Location must exist")
	}

	portStr := fmt.Sprint(cfg.Port)
	ar := Runner{
		server:    &http.Server{Addr: ":" + portStr},
		cfg:       cfg,
		keyCache:  make(map[string]string),
		cacheLock: &sync.RWMutex{},
	}

	return &ar, nil
}

func (ar *Runner) reloadMessageBus(ctx core.Context, c core.Components) {
	if c.MessageBus == nil {
		ctx.Log().Warn("Working in demo mode: without MessageBus")
	} else {
		ar.messageBus = c.MessageBus
	}
}

// Start runs api server
func (ar *Runner) Start(ctx core.Context, c core.Components) error {
	ar.reloadMessageBus(ctx, c)

	rootDomainReference := c.Bootstrapper.GetRootDomainRef()
	ar.netCoordinator = c.NetworkCoordinator

	ar.seedmanager = seedmanager.New()

	fw := wrapAPIV1Handler(ar, *rootDomainReference)
	http.HandleFunc(ar.cfg.Location, fw)
	http.HandleFunc(ar.cfg.Info, ar.infoHandler(c))
	http.HandleFunc(ar.cfg.Call, ar.callHandler(c))
	ctx.Log().Info("Starting ApiRunner ...")
	ctx.Log().Info("Config: ", ar.cfg)
	go func() {
		if err := ar.server.ListenAndServe(); err != nil {
			ctx.Log().Error("Httpserver: ListenAndServe() error: ", err)
		}
	}()
	return nil
}

// Stop stops api server
func (ar *Runner) Stop(ctx core.Context) error {
	const timeOut = 5

	ctx.Log().Infof("Shutting down server gracefully ...(waiting for %d seconds)", timeOut)
	ctxWithTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeOut)*time.Second)
	defer cancel()
	err := ar.server.Shutdown(ctxWithTimeout)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop API server")
	}

	return nil
}

func (ar *Runner) getMemberPubKey(ctx core.Context, ref string) (string, error) { //nolint
	ar.cacheLock.RLock()
	key, ok := ar.keyCache[ref]
	ar.cacheLock.RUnlock()
	if ok {
		return key, nil
	}
	args, err := core.MarshalArgs()
	if err != nil {
		return "", err
	}
	res, err := ar.messageBus.Send(
		ctx,
		&message.CallMethod{
			ObjectRef: core.NewRefFromBase58(ref),
			Method:    "GetPublicKey",
			Arguments: args,
		},
	)
	if err != nil {
		return "", err
	}

	var contractErr error
	err = signer.UnmarshalParams(res.(*reply.CallMethod).Result, &key, &contractErr)
	if err != nil {
		return "", err
	}
	if contractErr != nil {
		return "", contractErr
	}

	ar.cacheLock.Lock()
	ar.keyCache[ref] = key
	ar.cacheLock.Unlock()
	return key, nil
}

func RandTraceID() string {
	qid, err := uuid.NewV4()
	if err != nil {
		return "createRandomTraceIDFailed:" + err.Error()
	}
	return qid.String()
}
