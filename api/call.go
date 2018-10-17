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

		key, err := ar.getMemberPubKey(params.Reference)
		if err != nil {
			resp.Error = err.Error()
			log.Error(err)
			return
		}
		if key == "" {
			resp.Error = "Not found public key for this member"
			log.Error("[ CallHandler ] Not found public key for this member")
			return
		}

		from := core.NewRefFromBase58(params.Reference)
		args, err := core.MarshalArgs(
			from,
			params.Method,
			params.Params,
			params.Seed)
		if err != nil {
			resp.Error = err.Error()
			log.Error(errors.Wrap(err, "[ CallHandler ] Can't marshal arguments for verify signature"))
			return
		}

		verified, err := ecdsa.Verify(args, params.Signature, key)
		if err != nil {
			resp.Error = err.Error()
			log.Error(errors.Wrap(err, "[ CallHandler ] Can't verify signature"))
			return
		}
		if !verified {
			resp.Error = "Incorrect signature"
			log.Error("[ CallHandler ] Incorrect signature")
			return
		}

		args, err = core.MarshalArgs(*c.Bootstrapper.GetRootDomainRef(), params.Method, params.Params, params.Seed, params.Signature)
		if err != nil {
			resp.Error = err.Error()
			log.Error(err)
			return
		}
		res, err := ar.messageBus.Send(&message.CallMethod{
			ObjectRef: from,
			Method:    "Call",
			Arguments: args,
		})

		var typeHolder1 interface{}
		var typeHolder2 *foundation.Error
		refOrig, err := core.UnMarshalResponse(res.(*reply.CallMethod).Result, []interface{}{typeHolder1, typeHolder2})
		if err != nil {
			resp.Error = err.Error()
			log.Error(err)
			return
		}

		resp.Result = refOrig[0]
		if refOrig[1] != nil {
			resp.Error = refOrig[1].(*foundation.Error).S
		}
	}
}
