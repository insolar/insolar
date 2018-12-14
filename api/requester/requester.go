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

package requester

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

// verbose switches on verbose mode
var verbose = false
var scheme = platformpolicy.NewPlatformCryptographyScheme()

func verboseInfo(ctx context.Context, msg string) {
	if verbose {
		inslogger.FromContext(ctx).Infoln(msg)
	}
}

// SetVerbose switches on verbose mode
func SetVerbose(verb bool) {
	verbose = verb
}

// PostParams represents params struct
type PostParams = map[string]interface{}

// GetResponseBody makes request and extracts body
func GetResponseBody(url string, postP PostParams) ([]byte, error) {
	jsonValue, err := json.Marshal(postP)
	if err != nil {
		return nil, errors.Wrap(err, "[ getResponseBody ] Problem with marshaling params")
	}

	postResp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, errors.Wrap(err, "[ getResponseBody ] Problem with sending request")
	}
	if http.StatusOK != postResp.StatusCode {
		return nil, errors.New("[ getResponseBody ] Bad http response code: " + strconv.Itoa(postResp.StatusCode))
	}

	body, err := ioutil.ReadAll(postResp.Body)
	defer postResp.Body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "[ getResponseBody ] Problem with reading body")
	}

	return body, nil
}

// GetSeed makes rpc request to seed.Get method and extracts it
func GetSeed(url string) ([]byte, error) {
	body, err := GetResponseBody(url+"/rpc", PostParams{
		"jsonrpc": "2.0",
		"method":  "seed.Get",
		"id":      "",
	})
	if err != nil {
		return nil, errors.Wrap(err, "[ getSeed ]")
	}

	seedResp := rpcSeedResponse{}

	err = json.Unmarshal(body, &seedResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ getSeed ] Can't unmarshal")
	}
	if seedResp.Error != nil {
		return nil, errors.New("[ getSeed ] Field 'error' is not nil: " + fmt.Sprint(seedResp.Error))
	}
	res := &seedResp.Result
	if res == nil {
		return nil, errors.New("[ getSeed ] Field 'result' is nil")
	}

	return res.Seed, nil
}

func constructParams(params []interface{}) ([]byte, error) {
	args, err := core.MarshalArgs(params...)
	if err != nil {
		return nil, errors.Wrap(err, "[ constructParams ]")
	}
	return args, nil
}

// SendWithSeed sends request with known seed
func SendWithSeed(ctx context.Context, url string, userCfg *UserConfigJSON, reqCfg *RequestConfigJSON, seed []byte) ([]byte, error) {
	if userCfg == nil || reqCfg == nil {
		return nil, errors.New("[ Send ] Configs must be initialized")
	}

	params, err := constructParams(reqCfg.Params)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with serializing params")
	}

	callerRef, err := core.NewRefFromBase58(userCfg.Caller)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Failed to parse userCfg.Caller")
	}

	serRequest, err := core.MarshalArgs(
		*callerRef,
		reqCfg.Method,
		params,
		seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with serializing request")
	}

	verboseInfo(ctx, "Signing request ...")
	cs := scheme.Signer(userCfg.privateKeyObject)
	signature, err := cs.Sign(serRequest)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with signing request")
	}
	verboseInfo(ctx, "Signing request completed")

	body, err := GetResponseBody(url, PostParams{
		"params":    params,
		"method":    reqCfg.Method,
		"reference": userCfg.Caller,
		"seed":      seed,
		"signature": signature.Bytes(),
	})

	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with sending target request")
	}

	return body, nil
}

// Send first gets seed and after that makes target request
func Send(ctx context.Context, url string, userCfg *UserConfigJSON, reqCfg *RequestConfigJSON) ([]byte, error) {
	verboseInfo(ctx, "Sending GETSEED request ...")
	seed, err := GetSeed(url)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with getting seed")
	}
	verboseInfo(ctx, "GETSEED request completed. seed: "+string(seed))

	response, err := SendWithSeed(ctx, url+"/call", userCfg, reqCfg, seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ]")
	}

	return response, nil
}

func getDefaultRPCParams(method string) PostParams {
	return PostParams{
		"jsonrpc": "2.0",
		"id":      "",
		"method":  method,
	}
}

// Info makes rpc request to info.Get method and extracts it
func Info(url string) (*InfoResponse, error) {
	params := getDefaultRPCParams("info.Get")

	body, err := GetResponseBody(url+"/rpc", params)
	if err != nil {
		return nil, errors.Wrap(err, "[ Info ]")
	}

	infoResp := rpcInfoResponse{}

	err = json.Unmarshal(body, &infoResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ Info ] Can't unmarshal")
	}
	if infoResp.Error != nil {
		return nil, errors.New("[ Info ] Field 'error' is not nil: " + fmt.Sprint(infoResp.Error))
	}
	res := &infoResp.Result
	if res == nil {
		return nil, errors.New("[ Info ] Field 'result' is nil")
	}

	return res, nil
}

// Status makes rpc request to info.Status method and extracts it
func Status(url string) (*StatusResponse, error) {
	params := getDefaultRPCParams("status.Get")

	body, err := GetResponseBody(url+"/rpc", params)
	if err != nil {
		return nil, errors.Wrap(err, "[ Status ]")
	}

	statusResp := rpcStatusResponse{}

	err = json.Unmarshal(body, &statusResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ Status ] Can't unmarshal")
	}
	if statusResp.Error != nil {
		return nil, errors.New("[ Status ] Field 'error' is not nil: " + fmt.Sprint(statusResp.Error))
	}
	res := &statusResp.Result
	if res == nil {
		return nil, errors.New("[ Status ] Field 'result' is nil")
	}

	return res, nil
}
