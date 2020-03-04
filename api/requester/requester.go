// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package requester

import (
	"bytes"
	"context"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	mathrand "math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"strconv"
	"time"

	crypto "github.com/insolar/x-crypto"
	"github.com/insolar/x-crypto/ecdsa"
	"github.com/insolar/x-crypto/rand"
	"github.com/insolar/x-crypto/sha256"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

var httpClient *http.Client

const (
	RequestTimeout = 32 * time.Second
	Digest         = "Digest"
	Signature      = "Signature"
	ContentType    = "Content-Type"
	JSONRPCVersion = "2.0"
)

func init() {
	httpClient = createHTTPClient()
}

func SetTimeout(timeout uint) {
	if timeout > 0 {
		httpClient.Timeout = time.Duration(timeout) * time.Second
	} else {
		httpClient.Timeout = RequestTimeout
	}
}

// createHTTPClient for connection re-use
func createHTTPClient() *http.Client {
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: &http.Transport{},
		Timeout:   RequestTimeout,
		Jar:       jar,
	}

	return client
}

// verbose switches on verbose mode
var verbose = false

func verboseInfo(ctx context.Context, msg string) {
	if verbose {
		inslogger.FromContext(ctx).Info(msg)
	}
}

// SetVerbose switches on verbose mode
func SetVerbose(verb bool) {
	verbose = verb
}

// PlatformRequest represents params struct
type PlatformRequest struct {
	Request
	PlatformParams interface{} `json:"params,omitempty"`
	LogLevel       string      `json:"logLevel,omitempty"`
}

// ContractRequest is a representation of request struct to api
type ContractRequest struct {
	Request
	Params Params `json:"params,omitempty"`
}

type Request struct {
	Version string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
}

type Params struct {
	Seed       string      `json:"seed"`
	CallSite   string      `json:"callSite"`
	CallParams interface{} `json:"callParams,omitempty"`
	Reference  string      `json:"reference"`
	PublicKey  string      `json:"publicKey"`
	LogLevel   interface{} `json:"logLevel,omitempty"`
	Test       string      `json:"test,omitempty"`
}

// GetResponseBodyContract makes request to contract and extracts body
func GetResponseBodyContract(url string, postP ContractRequest, signature string) ([]byte, error) {
	req, err := MakeContractRequest(url, postP, signature)
	if err != nil {
		return nil, err
	}
	return doReq(context.Background(), req)
}

func MakeContractRequest(url string, postP ContractRequest, signature string) (*http.Request, error) {
	req, jsonValue, err := prepareReq(url, postP)
	if err != nil {
		return nil, errors.Wrap(err, "problem with preparing contract request")
	}

	h := sha256.New()
	_, err = h.Write(jsonValue)
	if err != nil {
		return nil, errors.Wrap(err, "[ GetResponseBodyContract ] Cant get hash")
	}
	sha := base64.StdEncoding.EncodeToString(h.Sum(nil))
	req.Header.Set(Digest, "SHA-256="+sha)
	req.Header.Set(Signature, "keyId=\"member-pub-key\", algorithm=\"ecdsa\", headers=\"digest\", signature="+signature)

	return req, nil
}

// GetResponseBodyPlatform makes request to platform and extracts body
func GetResponseBodyPlatform(url string, method string, params interface{}) ([]byte, error) {
	request := PlatformRequest{
		Request: Request{
			Version: JSONRPCVersion,
			ID:      uint64(mathrand.Int63()),
			Method:  method,
		},
		PlatformParams: params,
	}

	req, _, err := prepareReq(url, request)
	if err != nil {
		return nil, errors.Wrap(err, "problem with preparing platform request")
	}

	return doReq(context.Background(), req)
}

func prepareReq(url string, postP interface{}) (*http.Request, []byte, error) {
	jsonValue, err := json.Marshal(postP)
	if err != nil {
		return nil, nil, errors.Wrap(err, "problem with marshaling params")
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, nil, errors.Wrap(err, "problem with creating request")
	}
	req.Header.Set(ContentType, "application/json")

	return req, jsonValue, nil
}

func doReq(ctx context.Context, req *http.Request) ([]byte, error) {
	if verbose {
		requestDump, err := httputil.DumpRequest(req, true)
		if err != nil {
			verboseInfo(ctx, "Could not dump HTTP request")
			verboseInfo(ctx, err.Error())
		} else {
			verboseInfo(ctx, fmt.Sprintf("\n-----> %s\n--> END\n", requestDump))
		}
	}
	postResp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "problem with sending request")
	}

	if postResp == nil {
		return nil, errors.New("response is nil")
	}
	if verbose {
		requestDump, err := httputil.DumpResponse(postResp, true)
		if err != nil {
			verboseInfo(ctx, "Could not dump HTTP response")
			verboseInfo(ctx, err.Error())
		} else {
			verboseInfo(ctx, fmt.Sprintf("\n<----- %s\n<-- END\n\n", requestDump))
		}
	}

	defer postResp.Body.Close()
	if http.StatusOK != postResp.StatusCode {
		return nil, errors.New("bad http response code: " + strconv.Itoa(postResp.StatusCode))
	}

	body, err := ioutil.ReadAll(postResp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "problem with reading body")
	}

	return body, nil
}

// GetSeed makes rpc request to node.getSeed method and extracts it
func GetSeed(url string) (string, error) {
	body, err := GetResponseBodyPlatform(url, "node.getSeed", nil)
	if err != nil {
		return "", errors.Wrap(err, "[ GetSeed ] seed request")
	}

	seedResp := rpcSeedResponse{}

	err = json.Unmarshal(body, &seedResp)
	if err != nil {
		return "", errors.Wrap(err, "[ GetSeed ] Can't unmarshal")
	}
	if seedResp.Error != nil {
		return "", seedResp.Error
	}
	if len(seedResp.Result.Seed) == 0 {
		return "", errors.New("[ GetSeed ] Field seed is empty")
	}

	return seedResp.Result.Seed, nil
}

// SendWithSeed sends request with known seed
func SendWithSeed(ctx context.Context, url string, userCfg *UserConfigJSON, params *Params, seed string) ([]byte, error) {
	req, err := MakeRequestWithSeed(ctx, url, userCfg, params, seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ SendWithSeed ] Problem with creating target request")
	}
	b, err := doReq(ctx, req)
	return b, errors.Wrap(err, "[ SendWithSeed ] Problem with sending target request")
}

// MakeRequestWithSeed creates request with provided url, user config, params and seed.
func MakeRequestWithSeed(ctx context.Context, url string, userCfg *UserConfigJSON, params *Params, seed string) (*http.Request, error) {
	if userCfg == nil || params == nil {
		return nil, errors.New("configs must be initialized")
	}

	params.Reference = userCfg.Caller
	params.Seed = seed

	request := &ContractRequest{
		Request: Request{
			Version: JSONRPCVersion,
			ID:      uint64(mathrand.Int63()),
			Method:  "contract.call",
		},
		Params: *params,
	}

	verboseInfo(ctx, "Signing request ...")
	dataToSign, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(err, "config request marshaling failed")
	}
	signature, err := Sign(userCfg.privateKeyObject, dataToSign)
	if err != nil {
		return nil, errors.Wrap(err, "problem with signing request")
	}
	verboseInfo(ctx, "Signing request completed")

	return MakeContractRequest(url, *request, signature)
}

func Sign(privateKey crypto.PrivateKey, data []byte) (string, error) {
	hash := sha256.Sum256(data)

	r, s, err := ecdsa.Sign(rand.Reader, privateKey.(*ecdsa.PrivateKey), hash[:])
	if err != nil {
		return "", errors.Wrap(err, "[ sign ] Cant sign data")
	}

	return marshalSig(r, s)
}

// marshalSig encodes ECDSA signature to ASN.1.
func marshalSig(r, s *big.Int) (string, error) {
	var ecdsaSig struct {
		R, S *big.Int
	}
	ecdsaSig.R, ecdsaSig.S = r, s

	asnSig, err := asn1.Marshal(ecdsaSig)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(asnSig), nil
}

// Send first gets seed and after that makes target request
func Send(ctx context.Context, url string, userCfg *UserConfigJSON, params *Params) ([]byte, error) {
	verboseInfo(ctx, "Sending GETSEED request ...")
	seed, err := GetSeed(url)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with getting seed")
	}
	verboseInfo(ctx, "GETSEED request completed. seed: "+seed)

	response, err := SendWithSeed(ctx, url, userCfg, params, seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ]")
	}

	return response, nil
}

// Status makes rpc request to node.getStatus method and extracts it
func Status(url string) (*StatusResponse, error) {
	body, err := GetResponseBodyPlatform(url, "node.getStatus", nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ Status ]")
	}

	statusResp := rpcStatusResponse{}

	err = json.Unmarshal(body, &statusResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ Status ] Can't unmarshal")
	}
	if statusResp.Error != nil {
		return nil, statusResp.Error
	}

	return &statusResp.Result, nil
}
