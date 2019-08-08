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
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar/genesisrefs"
	"github.com/insolar/insolar/metrics"

	"github.com/insolar/insolar/api/requester"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
)

const (
	TimeoutError = 215
	ResultError  = 217
)

// UnmarshalRequest unmarshals request to api
func UnmarshalRequest(req *http.Request, params interface{}) ([]byte, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, errors.Wrap(err, "[ UnmarshalRequest ] Can't read body. So strange")
	}
	if len(body) == 0 {
		return nil, errors.New("[ UnmarshalRequest ] Empty body")
	}

	err = json.Unmarshal(body, &params)
	if err != nil {
		return body, errors.Wrap(err, "[ UnmarshalRequest ] Can't unmarshal input params")
	}
	return body, nil
}

func (ar *Runner) checkSeed(paramsSeed string) (insolar.PulseNumber, error) {
	decoded, err := base64.StdEncoding.DecodeString(paramsSeed)
	if err != nil {
		return 0, errors.New("[ checkSeed ] Failed to decode seed from string")
	}
	seed := seedmanager.SeedFromBytes(decoded)
	if seed == nil {
		return 0, errors.New("[ checkSeed ] Bad seed param")
	}

	if pulse, ok := ar.SeedManager.Pop(*seed); ok {
		return pulse, nil
	}

	return 0, errors.New("[ checkSeed ] Incorrect seed")
}

func (ar *Runner) makeCall(ctx context.Context, request requester.Request, rawBody []byte, signature string, pulseTimeStamp int64, seedPulse insolar.PulseNumber) (interface{}, error) {
	ctx, span := instracer.StartSpan(ctx, "SendRequest "+request.Method)
	defer span.End()

	reference, err := insolar.NewReferenceFromBase58(request.Params.Reference)
	if err != nil {
		return nil, errors.Wrap(err, "[ makeCall ] failed to parse params.Reference")
	}

	requestArgs, err := insolar.MarshalArgs(rawBody, signature, pulseTimeStamp)
	if err != nil {
		return nil, errors.Wrap(err, "[ makeCall ] failed to marshal arguments")
	}

	res, err := ar.ContractRequester.SendRequestWithPulse(
		ctx,
		reference,
		"Call",
		[]interface{}{requestArgs},
		seedPulse,
	)

	if err != nil {
		return nil, errors.Wrap(err, "[ makeCall ] Can't send request")
	}

	result, contractErr, err := extractor.CallResponse(res.(*reply.CallMethod).Result)

	if err != nil {
		return nil, errors.Wrap(err, "[ makeCall ] Can't extract response")
	}

	if contractErr != nil {
		return nil, errors.Wrap(errors.New(contractErr.S), "[ makeCall ] Error in called method")
	}

	return result, nil
}

func processError(err error, extraMsg string, resp *requester.ContractAnswer, insLog insolar.Logger, traceID string) {
	errResponse := &requester.Error{Message: extraMsg, Code: ResultError, Data: requester.Data{TraceID: traceID}}
	resp.Error = errResponse
	insLog.Error(errors.Wrapf(err, "[ CallHandler ] %s", extraMsg))
}

func writeResponse(insLog assert.TestingT, response http.ResponseWriter, contractAnswer *requester.ContractAnswer) {
	res, err := json.MarshalIndent(*contractAnswer, "", "    ")
	if err != nil {
		hcContractAnswer := requester.ContractAnswer{
			JSONRPC: "2.0",
			ID:      contractAnswer.ID,
			Error:   &requester.Error{Message: fmt.Sprintf("can't marshal ContractAnswer to json; error: '%v'", err.Error())},
		}
		res, _ = json.MarshalIndent(hcContractAnswer, "", "    ")
	}
	response.Header().Add("Content-Type", "application/json")
	_, err = response.Write(res)
	if err != nil {
		insLog.Errorf("Can't write response\n")
	}
}

func observeResultStatus(requestMethod string, contractAnswer *requester.ContractAnswer, startTime time.Time) {
	success := "success"
	if contractAnswer.Error != nil {
		success = "fail"
	}
	metrics.APIContractExecutionTime.WithLabelValues(requestMethod, success).Observe(time.Since(startTime).Seconds())
}

func processRequest(ctx context.Context,
	req *http.Request, contractRequest *requester.Request, contractAnswer *requester.ContractAnswer) (context.Context, []byte, error) {

	rawBody, err := UnmarshalRequest(req, contractRequest)
	if err != nil {
		return ctx, nil, errors.Wrap(err, "failed to unmarshal request")
	}

	contractAnswer.JSONRPC = contractRequest.JSONRPC
	contractAnswer.ID = contractRequest.ID

	if len(contractRequest.LogLevel) > 0 {
		logLevelNumber, err := insolar.ParseLevel(contractRequest.LogLevel)
		if err != nil {
			return ctx, nil, errors.Wrap(err, "failed to parse logLevel")
		}
		ctx = inslogger.WithLoggerLevel(ctx, logLevelNumber)
	}

	return ctx, rawBody, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func setRootReferenceIfNeeded(request *requester.Request) {
	if request.Params.Reference != "" {
		return
	}
	methods := []string{"member.create", "member.migrationCreate", "member.get"}
	if contains(methods, request.Params.CallSite) {
		request.Params.Reference = genesisrefs.ContractRootMember.String()
	}
}

func (ar *Runner) callHandler() func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {
		traceID := utils.RandTraceID()
		ctx, insLog := inslogger.WithTraceField(context.Background(), traceID)

		ctx, span := instracer.StartSpan(ctx, "ApiCallHandler.callHandler")
		defer span.End()

		contractRequest := &requester.Request{}
		contractAnswer := &requester.ContractAnswer{}
		defer writeResponse(insLog, response, contractAnswer)

		startTime := time.Now()
		defer observeResultStatus(contractRequest.Method, contractAnswer, startTime)

		info := fmt.Sprintf("[ callHandler ] Incoming contractRequest: %s", req.RequestURI)
		insLog.Infof(info)
		span.Annotate(nil, info)

		ctx, rawBody, err := processRequest(ctx, req, contractRequest, contractAnswer)
		if err != nil {
			instracer.AddError(span, err)
			processError(err, err.Error(), contractAnswer, insLog, traceID)
			return
		}

		if contractRequest.Test != "" {
			insLog.Infof("Request related to %s", contractRequest.Test)
		}

		if contractRequest.Method != "api.call" {
			err := errors.New("rpc method does not exist")
			instracer.AddError(span, err)
			processError(err, err.Error(), contractAnswer, insLog, traceID)
			return
		}

		signature, err := validateRequestHeaders(req.Header.Get(requester.Digest), req.Header.Get(requester.Signature), rawBody)
		if err != nil {
			instracer.AddError(span, err)
			processError(err, err.Error(), contractAnswer, insLog, traceID)
			return
		}

		seedPulse, err := ar.checkSeed(contractRequest.Params.Seed)
		if err != nil {
			instracer.AddError(span, err)
			processError(err, err.Error(), contractAnswer, insLog, traceID)
			return
		}

		setRootReferenceIfNeeded(contractRequest)

		var result interface{}
		ch := make(chan interface{}, 1)
		go func() {
			result, err = ar.makeCall(ctx, *contractRequest, rawBody, signature, 0, seedPulse)
			ch <- nil
		}()
		select {

		case <-ch:
			if err != nil {
				instracer.AddError(span, err)
				processError(err, err.Error(), contractAnswer, insLog, traceID)
				return
			}
			contractResult := &requester.Result{ContractResult: result, TraceID: traceID}
			contractAnswer.Result = contractResult
			return

		case <-time.After(ar.timeout):
			instracer.AddError(span, errors.New("API timeout exceeded"))
			errResponse := &requester.Error{Message: "API timeout exceeded", Code: TimeoutError, Data: requester.Data{TraceID: traceID}}
			contractAnswer.Error = errResponse
			return
		}
	}
}

func validateRequestHeaders(digest string, signature string, body []byte) (string, error) {
	// Digest = "SHA-256=<hashString>"
	// Signature = "keyId="member-pub-key", algorithm="ecdsa", headers="digest", signature=<signatureString>"
	if len(digest) < 15 || strings.Count(digest, "=") < 2 || len(signature) == 15 ||
		strings.Count(signature, "=") < 4 || len(body) == 0 {
		return "", errors.Errorf("invalid input data length digest: %d, signature: %d, body: %d", len(digest),
			len(signature), len(body))
	}
	h := sha256.New()
	_, err := h.Write(body)
	if err != nil {
		return "", errors.Wrap(err, "cant calculate hash")
	}
	calculatedHash := h.Sum(nil)
	digest, err = parseDigest(digest)
	if err != nil {
		return "", err
	}
	incomingHash, err := base64.StdEncoding.DecodeString(digest)
	if err != nil {
		return "", errors.Wrap(err, "cant decode digest")
	}

	if !bytes.Equal(calculatedHash, incomingHash) {
		return "", errors.New("incorrect digest")
	}

	signature, err = parseSignature(signature)
	if err != nil {
		return "", err
	}
	return signature, nil
}

func parseDigest(digest string) (string, error) {
	index := strings.IndexByte(digest, '=')
	if index < 1 || (index+1) >= len(digest) {
		return "", errors.New("invalid digest")
	}

	return digest[index+1:], nil
}

func parseSignature(signature string) (string, error) {
	index := strings.Index(signature, "signature=")
	if index < 1 || (index+10) >= len(signature) {
		return "", errors.New("invalid signature")
	}

	return signature[index+10:], nil
}
