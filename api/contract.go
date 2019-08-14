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
	"net/http"
	"strings"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/application/extractor"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/rpc/v2"
	"github.com/pkg/errors"
)

// ContractService is a service that provides API for working with smart contracts.
type ContractService struct {
	runner *Runner
}

// NewContractService creates new Contract service instance.
func NewContractService(runner *Runner) *ContractService {
	return &ContractService{runner: runner}
}

func (cs *ContractService) Call(req *http.Request, args *requester.Params, requestBody *rpc.RequestBody, result *requester.ContractResult) error {
	traceID := utils.RandTraceID()
	ctx, insLog := inslogger.WithTraceField(context.Background(), traceID)

	ctx, span := instracer.StartSpan(ctx, "Call")
	defer span.End()

	insLog.Infof("[ ContractService.Call ] Incoming request: %s", req.RequestURI)

	if args.Test != "" {
		insLog.Infof("ContractRequest related to %s", args.Test)
	}

	signature, err := validateRequestHeaders(req.Header.Get(requester.Digest), req.Header.Get(requester.Signature), requestBody.Raw)
	if err != nil {
		return err
	}

	seedPulse, err := cs.runner.checkSeed(args.Seed)
	if err != nil {
		return err
	}

	setRootReferenceIfNeeded(args)

	callResult, _, err := cs.runner.makeCall(ctx, "contract.call", *args, requestBody.Raw, signature, 0, seedPulse)
	if err != nil {
		return err
	}

	result.CallResult = callResult
	result.TraceID = traceID

	return nil
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

func (ar *Runner) makeCall(ctx context.Context, method string, params requester.Params, rawBody []byte, signature string, pulseTimeStamp int64, seedPulse insolar.PulseNumber) (interface{}, *insolar.Reference, error) {
	ctx, span := instracer.StartSpan(ctx, "SendRequest "+method)
	defer span.End()

	reference, err := insolar.NewReferenceFromBase58(params.Reference)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ makeCall ] failed to parse params.Reference")
	}

	requestArgs, err := insolar.Serialize([]interface{}{rawBody, signature, pulseTimeStamp})
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ makeCall ] failed to marshal arguments")
	}

	res, ref, err := ar.ContractRequester.SendRequestWithPulse(
		ctx,
		reference,
		"Call",
		[]interface{}{requestArgs},
		seedPulse,
	)

	if err != nil {
		return nil, ref, errors.Wrap(err, "[ makeCall ] Can't send request")
	}

	result, contractErr, err := extractor.CallResponse(res.(*reply.CallMethod).Result)

	if err != nil {
		return nil, ref, errors.Wrap(err, "[ makeCall ] Can't extract response")
	}

	if contractErr != nil {
		return nil, ref, errors.Wrap(errors.New(contractErr.S), "[ makeCall ] Error in called method")
	}

	return result, ref, nil
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
