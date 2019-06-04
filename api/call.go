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
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	insolarJose "github.com/insolar/go-jose"
	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/metrics"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

// Request is a representation of request struct to api
type Request struct {
	PublicKey string  `json:"jwk"`
	Signature string  `json:"jws"`
	LogLevel  *string `json:"logLevel,omitempty"`
}

// Data which is signed
type SignedPayload struct {
	Reference string `json:"reference"` // contract reference
	Method    string

	// method name
	Params string `json:"params"` // json object
	Seed   string `json:"seed"`
}

type answer struct {
	Error   string      `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	TraceID string      `json:"traceID,omitempty"`
}

// UnmarshalRequest unmarshal request to api
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

func (ar *Runner) checkSeedString(_seed string) error {
	decoded, err := base64.StdEncoding.DecodeString(_seed)
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

func (ar *Runner) makeCall(ctx context.Context, signedPayload SignedPayload, request Request) (interface{}, error) {
	ctx, span := instracer.StartSpan(ctx, "SendRequest "+signedPayload.Method)
	defer span.End()

	reference, err := insolar.NewReferenceFromBase58(signedPayload.Reference)
	if err != nil {
		return nil, errors.Wrap(err, "[ makeCall ] failed to parse signedPayload.Reference")
	}

	type PostParams = map[string]interface{}
	args, err := insolar.MarshalArgs(string(request.PublicKey), request.Signature)
	if err != nil {
		return nil, errors.Wrap(err, "[ makeCall ] failed to marshal arguments")
	}

	res, err := ar.ContractRequester.SendRequest(
		ctx,
		reference,
		"Call",
		[]interface{}{*ar.CertificateManager.GetCertificate().GetRootDomainReference(), args},
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

func processError(err error, extraMsg string, resp *answer, insLog insolar.Logger) {
	resp.Error = err.Error()
	insLog.Error(errors.Wrapf(err, "[ CallHandler ] %s", extraMsg))
}

func (ar *Runner) callHandler() func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, req *http.Request) {
		traceID := utils.RandTraceID()
		ctx, insLog := inslogger.WithTraceField(context.Background(), traceID)

		ctx, span := instracer.StartSpan(ctx, "callHandler")
		defer span.End()

		// unmarshal jws
		request := Request{}
		resp := answer{}

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

		_, err := UnmarshalRequest(req, &request)
		if err != nil {
			processError(err, "Can't unmarshal request", &resp, insLog)
			return
		}

		signedPayload, err := verify(request.PublicKey, request.Signature)
		if err != nil {
			processError(err, "Can't verify signature", &resp, insLog)
			return
		}

		startTime := time.Now()
		defer func() {
			success := "success"
			if resp.Error != "" {
				success = "fail"
			}
			metrics.APIContractExecutionTime.WithLabelValues(signedPayload.Method, success).Observe(time.Since(startTime).Seconds())
		}()

		resp.TraceID = traceID

		insLog.Infof("[ callHandler ] Incoming request: %s", req.RequestURI)

		if request.LogLevel != nil {
			logLevelNumber, err := insolar.ParseLevel(*request.LogLevel)
			if err != nil {
				processError(err, "Can't parse logLevel", &resp, insLog)
				return
			}
			ctx = inslogger.WithLoggerLevel(ctx, logLevelNumber)
		}

		err = ar.checkSeedString(signedPayload.Seed)
		if err != nil {
			processError(err, "Can't checkSeed", &resp, insLog)
			return
		}

		var result interface{}
		ch := make(chan interface{}, 1)
		go func() {
			result, err = ar.makeCall(ctx, *signedPayload, request)
			ch <- nil
		}()
		select {

		case <-ch:
			if err != nil {
				processError(err, "Can't makeCall", &resp, insLog)
				return
			}
			resp.Result = result

		case <-time.After(time.Duration(ar.cfg.Timeout) * time.Second):
			resp.Error = "Messagebus timeout exceeded"
			return

		}

		resp.Result = result
	}
}

// verify and get payload
// return parsed payload, public key, signature and error
func verify(publicKey string, _signature string) (*SignedPayload, error) {
	type rawJSONWebKey struct {
		Crv string `json:"crv,omitempty"`
	}
	var raw rawJSONWebKey
	err := json.Unmarshal([]byte(publicKey), &raw)
	if err != nil {
		return nil, err
	}
	switch raw.Crv {
	case "P-256K":
		{
			// unmarshal public key
			public := insolarJose.JSONWebKey{}
			err = public.UnmarshalJSON([]byte(publicKey))
			if err != nil {
				return nil, err
			}
			// parse jws token
			signature, err := insolarJose.ParseSigned(_signature)
			if err != nil {
				return nil, err
			}
			var payloadRequest = SignedPayload{}

			payload, err := signature.Verify(&public)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(payload, &payloadRequest)
			if err != nil {
				return nil, err
			}

			return &payloadRequest, nil
		}
	case "P-256":
		{
			// unmarshal public key
			public := jose.JSONWebKey{}
			err = public.UnmarshalJSON([]byte(publicKey))
			if err != nil {
				return nil, err
			}
			// parse jws token
			signature, err := jose.ParseSigned(_signature)
			if err != nil {
				return nil, err
			}
			var payloadRequest = SignedPayload{}

			payload, err := signature.Verify(&public)
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(payload, &payloadRequest)
			if err != nil {
				return nil, err
			}

			return &payloadRequest, nil

		}
	default:
		return nil, errors.Errorf("Unsupported key format")
	}

}
