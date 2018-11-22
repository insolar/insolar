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

	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/cryptography"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/pkg/errors"
)

type request struct {
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

func unmarshalRequest(req *http.Request, params interface{}) ([]byte, error) {
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

func (ar *Runner) verifySignature(ctx context.Context, params request) error {
	key, err := ar.getMemberPubKey(ctx, params.Reference)
	if err != nil {
		return err
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

	cs := cryptography.NewKeyBoundCryptographyService(nil)
	verified := cs.Verify(key, core.SignatureFromBytes(params.Signature), args)
	if !verified {
		return errors.New("[ VerifySignature ] Incorrect signature")
	}
	return nil
}

func (ar *Runner) callHandler() func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {

		params := request{}
		resp := answer{}

		traceid := utils.RandTraceID()
		ctx, inslog := inslogger.WithTraceField(context.Background(), traceid)
		resp.TraceID = traceid

		defer func() {
			res, err := json.MarshalIndent(resp, "", "    ")
			if err != nil {
				res = []byte(`{"error": "can't marshal answer to json'"}`)
			}
			response.Header().Add("Content-Type", "application/json")
			_, err = response.Write(res)
			if err != nil {
				inslog.Errorf("Can't write response\n")
			}
		}()

		_, err := unmarshalRequest(req, &params)
		if err != nil {
			resp.Error = err.Error()
			inslog.Error(errors.Wrap(err, "[ CallHandler ] Can't unmarshal request"))
			return
		}

		seed := seedmanager.SeedFromBytes(params.Seed)
		if seed == nil {
			resp.Error = "[ CallHandler ] Bad seed param"
			inslog.Error(resp.Error)
			return
		}

		if !ar.seedmanager.Exists(*seed) {
			resp.Error = "[ CallHandler ] Incorrect seed"
			inslog.Error(resp.Error)
			return
		}

		err = ar.verifySignature(ctx, params)
		if err != nil {
			resp.Error = err.Error()
			inslog.Error(errors.Wrap(err, "[ CallHandler ] Can't verify signature"))
			return
		}

		reference := core.NewRefFromBase58(params.Reference)
		res, err := ar.ContractRequester.SendRequest(
			ctx,
			&reference,
			"Call",
			[]interface{}{*ar.Certificate.GetRootDomainReference(), params.Method, params.Params, params.Seed, params.Signature},
		)
		if err != nil {
			resp.Error = err.Error()
			inslog.Error(errors.Wrap(err, "[ CallHandler ] Can't send request"))
			return
		}

		var result interface{}
		var contractErr *foundation.Error
		err = signer.UnmarshalParams(res.(*reply.CallMethod).Result, &result, &contractErr)
		if err != nil {
			resp.Error = err.Error()
			inslog.Error(errors.Wrap(err, "[ CallHandler ] Can't unmarshal params"))
			return
		}

		resp.Result = result
		if contractErr != nil {
			resp.Error = contractErr.S
			inslog.Error(errors.Wrap(errors.New(contractErr.S), "[ CallHandler ] Error in called method"))
		}
	}
}
