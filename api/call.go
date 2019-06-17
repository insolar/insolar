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
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/insolar/insolar/api/requester"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/metrics"
)

const (
	TimeoutError = -215
	ResultError  = -217
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

func (ar *Runner) checkSeed(paramsSeed string) error {
	decoded, err := base64.StdEncoding.DecodeString(paramsSeed)
	if err != nil {
		return errors.New("[ checkSeed ] Decode error")
	}
	seed := seedmanager.SeedFromBytes(decoded)
	if seed == nil {
		return errors.New("[ checkSeed ] Bad seed param")
	}

	if !ar.SeedManager.Exists(*seed) {
		return errors.New("[ checkSeed ] Incorrect seed")
	}

	return nil
}

func (ar *Runner) makeCall(ctx context.Context, request requester.Request, rawBody []byte, signature string, pulseTimeStamp int64) (interface{}, error) {
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

	res, err := ar.ContractRequester.SendRequest(
		ctx,
		reference,
		"Call",
		[]interface{}{*ar.CertificateManager.GetCertificate().GetRootDomainReference(), requestArgs},
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
	resp.Error.Message = err.Error()
	resp.Error.Code = ResultError
	resp.Error.TraceID = traceID
	insLog.Error(errors.Wrapf(err, "[ CallHandler ] %s", extraMsg))
}

func (ar *Runner) callHandler() func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {
		traceID := utils.RandTraceID()
		ctx, insLog := inslogger.WithTraceField(context.Background(), traceID)

		ctx, span := instracer.StartSpan(ctx, "callHandler")
		defer span.End()

		request := requester.Request{}
		resp := requester.ContractAnswer{}

		startTime := time.Now()
		defer func() {
			success := "success"
			if resp.Error.Message != "" {
				success = "fail"
			}
			metrics.APIContractExecutionTime.WithLabelValues(request.Method, success).Observe(time.Since(startTime).Seconds())
		}()

		insLog.Infof("[ callHandler ] Incoming request: %s", req.RequestURI)

		defer func() {
			res, err := json.MarshalIndent(resp, "", "    ")
			if err != nil {
				res = []byte(`{"error": "can't marshal ContractAnswer to json'"}`)
			}
			response.Header().Add("Content-Type", "application/json")
			_, err = response.Write(res)
			if err != nil {
				insLog.Errorf("Can't write response\n")
			}
		}()

		rawBody, err := UnmarshalRequest(req, &request)
		if err != nil {
			processError(err, "Can't unmarshal request", &resp, insLog, traceID)
			return
		}

		resp.JSONRPC = request.JSONRPC
		resp.ID = request.ID

		digest := req.Header.Get("Digest")
		richSignature := req.Header.Get("Signature")

		signature, err := validateRequestHeaders(digest, richSignature, rawBody)
		if err != nil {
			processError(err, "Can't validate request", &resp, insLog, traceID)
			return
		}

		if len(request.LogLevel) > 0 {
			logLevelNumber, err := insolar.ParseLevel(request.LogLevel)
			if err != nil {
				processError(err, "Can't parse logLevel", &resp, insLog, traceID)
				return
			}
			ctx = inslogger.WithLoggerLevel(ctx, logLevelNumber)
		}

		err = ar.checkSeed(request.Params.Seed)
		if err != nil {
			processError(err, "Can't checkSeed", &resp, insLog, traceID)
			return
		}

		pulse, err := ar.PulseAccessor.Latest(ctx)
		if err != nil {
			processError(err, "Can't get last pulse", &resp, insLog, traceID)
			return
		}

		var result interface{}
		ch := make(chan interface{}, 1)
		go func() {
			result, err = ar.makeCall(ctx, request, rawBody, signature, pulse.PulseTimestamp)
			ch <- nil
		}()
		select {

		case <-ch:
			if err != nil {
				processError(err, "Can't makeCall", &resp, insLog, traceID)
				return
			}
			resp.Result.ContractResult = result

		case <-time.After(time.Duration(ar.cfg.Timeout) * time.Second):
			resp.Error.Message = "Messagebus timeout exceeded"
			resp.Error.Code = TimeoutError
			resp.Error.TraceID = traceID
			return
		}

		resp.Result.ContractResult = result
		resp.Result.TraceID = traceID
	}
}

func validateRequestHeaders(digest string, richSignature string, body []byte) (string, error) {
	if len(digest) == 0 || len(richSignature) == 0 || len(body) == 0 {
		return "", errors.New("Invalid input data")
	}
	h := sha256.New()
	_, err := h.Write(body)
	if err != nil {
		return "", errors.Wrap(err, "Cant get hash")
	}
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	if sha == digest[strings.IndexByte(digest, '=')+1:] {
		sig := richSignature[strings.Index(richSignature, "signature=")+10:]
		if len(sig) == 0 {
			return "", errors.New("Empty signature")
		}
		return sig, nil

	}
	return "", errors.New("Cant get signature from header")
}
