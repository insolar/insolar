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
	"github.com/insolar/insolar/application/proxy/member"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

const (
	_            int = 0
	HandlerError int = -1
	BadRequest   int = -2
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

func makeHandlerMarshalErrorJSON() []byte {
	jsonErr := writeError("Invalid data from handler", HandlerError)
	serJSON, err := json.Marshal(jsonErr)
	if err != nil {
		log.Fatal("Can't marshal base error")
	}
	return serJSON
}

var handlerMarshalErrorJSON = makeHandlerMarshalErrorJSON()

func processQueryType(rh *RequestHandler, qTypeStr string) map[string]interface{} {
	qtype := QTypeFromString(qTypeStr)
	var answer map[string]interface{}

	var hError error
	switch qtype {
	case CreateMember:
		answer, hError = rh.ProcessCreateMember()
	case DumpUserInfo:
		answer, hError = rh.ProcessDumpUsers(false)
	case DumpAllUsers:
		answer, hError = rh.ProcessDumpUsers(true)
	case GetBalance:
		answer, hError = rh.ProcessGetBalance()
	case SendMoney:
		answer, hError = rh.ProcessSendMoney()
	case RegisterNode:
		answer, hError = rh.ProcessRegisterNode()
	case IsAuth:
		answer, hError = rh.ProcessIsAuthorized()
	case GetSeed:
		answer, hError = rh.ProcessGetSeed()
	default:
		msg := fmt.Sprintf("Wrong query parameter 'query_type' = '%s'", qTypeStr)
		answer = writeError(msg, BadRequest)
		log.Warnf("[QID=%s] %s\n", rh.qid, msg)
		return answer
	}
	if hError != nil {
		errMsg := "Handler error: " + hError.Error()
		log.Errorf("[QID=%s] %s\n", rh.qid, errMsg)
		answer = writeError(errMsg, HandlerError)
	}

	return answer
}

const qidQueryParam = "qid"

// PreprocessRequest extracts params from requests
func PreprocessRequest(req *http.Request) (*Params, error) {
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

	if len(params.QID) == 0 {
		qid, err := uuid.NewV1()
		if err == nil {
			params.QID = qid.String()
		}
	}

	log.Infof("[QID=%s] Query: %s. Url: %s\n", params.QID, string(body), req.URL)

	return &params, nil
}

func wrapAPIV1Handler(runner *Runner, rootDomainReference core.RecordRef) func(w http.ResponseWriter, r *http.Request) {
	sm := seedmanager.New()
	return func(response http.ResponseWriter, req *http.Request) {
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
			answer[qidQueryParam] = params.QID
			serJSON, err := json.MarshalIndent(answer, "", "    ")
			if err != nil {
				serJSON = handlerMarshalErrorJSON
			}
			response.Header().Add("Content-Type", "application/json")
			var newLine byte = '\n'
			_, err = response.Write(append(serJSON, newLine))
			if err != nil {
				log.Errorf("[QID=%s] Can't write response\n", params.QID)
			}
			log.Infof("[QID=%s] Request completed. Total time: %s\n", params.QID, time.Since(startTime))
		}()

		params, err := PreprocessRequest(req)
		if err != nil {
			answer = writeError("Bad request", BadRequest)
			log.Errorf("[QID=] Can't parse input request: %s, error: %s\n", req.RequestURI, err)
			return
		}
		rh := NewRequestHandler(params, runner.messageBus, runner.netCoordinator, rootDomainReference, sm)

		answer = processQueryType(rh, params.QueryType)
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

func (ar *Runner) reloadMessageBus(c core.Components) {
	if c.MessageBus == nil {
		log.Warn("Working in demo mode: without MessageBus")
	} else {
		ar.messageBus = c.MessageBus
	}
}

// Start runs api server
func (ar *Runner) Start(c core.Components) error {

	ar.reloadMessageBus(c)

	rootDomainReference := c.Bootstrapper.GetRootDomainRef()
	ar.netCoordinator = c.NetworkCoordinator

	fw := wrapAPIV1Handler(ar, *rootDomainReference)
	http.HandleFunc(ar.cfg.Location, fw)
	http.HandleFunc(ar.cfg.Info, ar.infoHandler(c))
	log.Info("Starting ApiRunner ...")
	log.Info("Config: ", ar.cfg)
	go func() {
		if err := ar.server.ListenAndServe(); err != nil {
			log.Error("Httpserver: ListenAndServe() error: ", err)
		}
	}()
	return nil
}

// Stop stops api server
func (ar *Runner) Stop() error {
	const timeOut = 5
	log.Infof("Shutting down server gracefully ...(waiting for %d seconds)", timeOut)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	err := ar.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop API server")
	}

	return nil
}

func (ar *Runner) getMemberPubKey(ref string) string { //nolint
	ar.cacheLock.RLock()
	key, ok := ar.keyCache[ref]
	ar.cacheLock.RUnlock()
	if ok {
		return key
	}
	key = member.GetObject(core.NewRefFromBase58(ref)).GetPublicKey()
	ar.cacheLock.Lock()
	ar.keyCache[ref] = key
	ar.cacheLock.Unlock()
	return key
}
