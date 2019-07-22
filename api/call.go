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

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar/genesisrefs"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
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

func (ar *Runner) makeCall(ctx context.Context, method string, params requester.Params, rawBody []byte, signature string, pulseTimeStamp int64, seedPulse insolar.PulseNumber) (interface{}, error) {
	ctx, span := instracer.StartSpan(ctx, "SendRequest "+method)
	defer span.End()

	reference, err := insolar.NewReferenceFromBase58(params.Reference)
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func setRootReferenceIfNeeded(params *requester.Params) {
	if params.Reference != "" {
		return
	}
	methods := []string{"member.create", "member.migrationCreate", "member.get"}
	if contains(methods, params.CallSite) {
		params.Reference = genesisrefs.ContractRootMember.String()
	}
}

func validateRequestHeaders(digest string, richSignature string, body []byte) (string, error) {
	// Digest = "SHA-256=<hashString>"
	// Signature = "keyId="member-pub-key", algorithm="ecdsa", headers="digest", signature=<signatureString>"
	if len(digest) < 15 || strings.Count(digest, "=") < 2 || len(richSignature) == 15 ||
		strings.Count(richSignature, "=") < 4 || len(body) == 0 {
		return "", errors.Errorf("invalid input data length digest: %d, signature: %d, body: %d", len(digest),
			len(richSignature), len(body))
	}
	h := sha256.New()
	_, err := h.Write(body)
	if err != nil {
		return "", errors.Wrap(err, "Cant get hash")
	}
	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))
	if sha == digest[strings.IndexByte(digest, '=')+1:] {
		sig := richSignature[strings.Index(richSignature, "signature=")+10:]
		if len(sig) == 0 {
			return "", errors.New("empty signature")
		}
		return sig, nil

	}
	return "", errors.New("cant get signature from header")
}
