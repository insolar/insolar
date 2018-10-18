/*
 *    Copyright 2018 INS Ecosystem
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
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/insolar/insolar/application/contract/member/signer"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/log"
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
	Error  string      `json:"error,omitempty"`
	Result interface{} `json:"result,omitempty"`
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

func (ar *Runner) verifySignature(params request) error {
	key, err := ar.getMemberPubKey(params.Reference)
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

	verified, err := ecdsa.Verify(args, params.Signature, key)
	if err != nil {
		return errors.Wrap(err, "[ VerifySignature ] Can't verify signature")
	}
	if !verified {
		return errors.New("[ VerifySignature ] Incorrect signature")
	}
	return nil
}

func (ar *Runner) callHandler(c core.Components) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {

		params := request{}
		resp := answer{}

		defer func() {
			res, err := json.MarshalIndent(resp, "", "    ")
			if err != nil {
				res = []byte(`{"error": "can't marshal answer to json'"}`)
			}
			response.Header().Add("Content-Type", "application/json")
			_, err = response.Write(res)
			if err != nil {
				log.Errorf("Can't write response\n")
			}
		}()

		_, err := unmarshalRequest(req, &params)
		if err != nil {
			resp.Error = err.Error()
			log.Error(errors.Wrap(err, "[ CallHandler ] Can't unmarshal request"))
			return
		}

		if !ar.seedmanager.Exists(ar.seedmanager.SeedFromBytes(params.Seed)) {
			resp.Error = "Incorrect seed"
			log.Error("[ CallHandler ] ", resp.Error)
			return
		}

		err = ar.verifySignature(params)
		if err != nil {
			resp.Error = err.Error()
			log.Error(errors.Wrap(err, "[ CallHandler ] "))
			return
		}

		args, err := core.MarshalArgs(*c.Bootstrapper.GetRootDomainRef(), params.Method, params.Params, params.Seed, params.Signature)
		if err != nil {
			resp.Error = err.Error()
			log.Error(err)
			return
		}
		res, err := ar.messageBus.Send(&message.CallMethod{
			ObjectRef: core.NewRefFromBase58(params.Reference),
			Method:    "Call",
			Arguments: args,
		})
		if err != nil {
			resp.Error = err.Error()
			log.Error(err)
			return
		}

		var result interface{}
		var contractErr *foundation.Error
		err = signer.UnmarshalParams(res.(*reply.CallMethod).Result, &result, &contractErr)
		if err != nil {
			resp.Error = err.Error()
			log.Error(err)
			return
		}

		resp.Result = result
		if contractErr != nil {
			resp.Error = contractErr.S
		}
	}
}
