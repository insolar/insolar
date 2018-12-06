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
	"io/ioutil"
	"net/http"

	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/platformpolicy"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

var scheme = platformpolicy.NewPlatformCryptographyScheme()

// Request is a representation of request struct to api
type Request struct {
	Reference string `json:"reference"`
	Method    string `json:"method"`
	Params    []byte `json:"params"`
	Seed      []byte `json:"seed"`
	Signature []byte `json:"signature"`
}

type answer struct {
	Error   string      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	TraceID string      `json:"traceID,omitempty"`
}

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

func (ar *Runner) verifySignature(ctx context.Context, params Request) error {
	key, err := ar.getMemberPubKey(ctx, params.Reference)
	if err != nil {
		return errors.Wrap(err, "[ VerifySignature ] Can't getMemberPubKey")
	}
	if key == "" {
		return errors.New("[ VerifySignature ] Not found public key for this member")
	}

	args, err := core.MarshalArgs(
		core.NewRefFromBase58(params.Reference),
		params.Method,
		params.Params,
		params.Seed)
	if err != nil {
		return errors.Wrap(err, "[ VerifySignature ] Can't marshal arguments for verify signature")
	}
	verifier := scheme.Verifier(key)
	verified := verifier.Verify(core.SignatureFromBytes(params.Signature), args)
	if !verified {
		return errors.New("[ VerifySignature ] Incorrect signature")
	}
	return nil
}

func (ar *Runner) checkSeed(paramsSeed []byte) error {
	seed := seedmanager.SeedFromBytes(paramsSeed)
	if seed == nil {
		return errors.New("[ checkSeed ] Bad seed param")
	}

	if !ar.SeedManager.Exists(*seed) {
		return errors.New("[ checkSeed ] Incorrect seed")
	}

	return nil
}

func (ar *Runner) makeCall(ctx context.Context, params Request) (interface{}, error) {
	reference := core.NewRefFromBase58(params.Reference)
	res, err := ar.ContractRequester.SendRequest(
		ctx,
		&reference,
		"Call",
		[]interface{}{*ar.CertificateManager.GetCertificate().GetRootDomainReference(), params.Method, params.Params, params.Seed, params.Signature},
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

func processError(err error, extraMsg string, resp *answer, insLog core.Logger) {
	resp.Error = err.Error()
	insLog.Error(errors.Wrapf(err, "[ CallHandler ] %s", extraMsg))
}

func (ar *Runner) callHandler() func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {

		params := Request{}
		resp := answer{}

		traceID := utils.RandTraceID()
		ctx, insLog := inslogger.WithTraceField(context.Background(), traceID)
		resp.TraceID = traceID

		insLog.Info("[ callHandler ] Incoming request: %s", req.RequestURI)

		defer func() {
			res, err := json.MarshalIndent(resp, "", "    ")
			if err != nil {
				res = []byte(`{"error": "can't marshal answer to json'"}`)
			}
			response.Header().Add("Content-Type", "application/json")
			_, err = response.Write(res)
			if err != nil {
				insLog.Errorf("Can't write response\n")
			}
		}()

		_, err := UnmarshalRequest(req, &params)
		if err != nil {
			processError(err, "Can't unmarshal request", &resp, insLog)
			return
		}

		err = ar.checkSeed(params.Seed)
		if err != nil {
			processError(err, "Can't checkSeed", &resp, insLog)
			return
		}

		err = ar.verifySignature(ctx, params)
		if err != nil {
			processError(err, "Can't verify signature", &resp, insLog)
			return
		}

		result, err := ar.makeCall(ctx, params)
		if err != nil {
			processError(err, "Can't makeCall", &resp, insLog)
			return
		}

		resp.Result = result
	}
}
