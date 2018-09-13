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
	"log"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type ResponseCode int

const (
	Ok           ResponseCode = 0
	HandlerError              = -1
	BadRequest                = -2
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

	var handlerError error
	switch qtype {
	case CreateMember:
		answer, handlerError = rh.ProcessCreateMember()
	case DumpUserInfo:
		answer, handlerError = rh.ProcessDumpUsers(false)
	case DumpAllUsers:
		answer, handlerError = rh.ProcessDumpUsers(true)
	case GetBalance:
		answer, handlerError = rh.ProcessGetBalance()
	case SendMoney:
		answer, handlerError = rh.ProcessSendMoney()
	default:
		msg := fmt.Sprintf("Wrong query parameter 'query_type' = '%s'", qTypeStr)
		answer = writeError(msg, BadRequest)
		logrus.Printf("[QID=%s] %s\n", rh.qid, msg)
		return answer
	}
	if handlerError != nil {
		errMsg := "Handler error: " + handlerError.Error()
		log.Printf("[QID=%s] %s\n", rh.qid, errMsg)
		answer = writeError(errMsg, HandlerError)
	}

	return answer
}

const QIDQueryParam = "qid"

func preprocessRequest(req *http.Request) (*Params, error) {
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

	logrus.Printf("[QID=%s] Query: %s\n", params.QID, string(body))

	return &params, nil
}

func wrapAPIV1Handler(router core.MessageRouter) func(w http.ResponseWriter, r *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {
		answer := make(map[string]interface{})
		var params *Params
		defer func() {
			if answer == nil {
				answer = make(map[string]interface{})
			}
			if params == nil {
				params = &Params{}
			}
			answer[QIDQueryParam] = params.QID
			serJSON, err := json.MarshalIndent(answer, "", "    ")
			if err != nil {
				serJSON = handlerMarshalErrorJSON
			}
			response.Header().Add("Content-Type", "application/json")
			var newLine byte = '\n'
			_, err = response.Write(append(serJSON, newLine))
			if err != nil {
				logrus.Printf("[QID=%s] Can't write response\n", params.QID)
			}
			logrus.Infof("[QID=%s] Request completed\n", params.QID)
		}()

		params, err := preprocessRequest(req)
		if err != nil {
			answer = writeError("Bad request", BadRequest)
			logrus.Errorf("[QID=]Can't parse input request: %s\n", err, req.RequestURI)
			return
		}
		rh := NewRequestHandler(params, router)

		answer = processQueryType(rh, params.QType)
	}
}

type APIRunner struct {
	messageRouter core.MessageRouter
	server        *http.Server
	cfg           *configuration.APIRunner
}

func NewAPIRunner(cfg *configuration.APIRunner) (*APIRunner, error) {
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
	ar := APIRunner{
		server: &http.Server{Addr: ":" + portStr},
		cfg:    cfg,
	}

	return &ar, nil
}

// Start runs api server
func (ar *APIRunner) Start(c core.Components) error {

	// TODO: init message router
	_, ok := c["core.MessageRouter"]
	if !ok {
		logrus.Warnln("Working in demo mode: without MessageRouter")
	} else {
		ar.messageRouter = c["core.MessageRouter"].(core.MessageRouter)
	}

	fw := wrapAPIV1Handler(ar.messageRouter)
	http.HandleFunc(ar.cfg.Location, fw)
	logrus.Info("Starting ApiRunner ...")
	logrus.Info("Config: ", ar.cfg)
	go func() {
		if err := ar.server.ListenAndServe(); err != nil {
			logrus.Errorln("Httpserver: ListenAndServe() error: ", err)
		}
	}()
	return nil
}

// Stop stops api server
func (ar *APIRunner) Stop() error {
	const timeOut = 5
	logrus.Infof("Shutting down server gracefully ...(waiting for %d seconds)", timeOut)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeOut)*time.Second)
	defer cancel()
	err := ar.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "Can't gracefully stop API server")
	}

	return nil
}
