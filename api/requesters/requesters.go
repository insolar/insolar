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

package requesters

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/insolar/insolar/core"
	ecdsahelper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

// verbose switches on verbose mode
var verbose = false

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

// GetSeed makes get_seed request and extracts it
func GetSeed(url string) ([]byte, error) {
	body, err := GetResponseBody(url, PostParams{
		"query_type": "get_seed",
	})
	if err != nil {
		return nil, errors.Wrap(err, "[ getSeed ]")
	}

	type seedResponse struct{ Seed []byte }
	seedResp := seedResponse{}

	err = json.Unmarshal(body, &seedResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ getSeed ]")
	}

	return seedResp.Seed, nil
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

	serRequest, err := core.MarshalArgs(
		core.NewRefFromBase58(userCfg.Caller),
		reqCfg.Method,
		params,
		seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with serializing request")
	}

	verboseInfo(ctx, "Signing request ...")
	signature, err := ecdsahelper.Sign(serRequest, userCfg.privateKeyObject)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with signing request")
	}
	verboseInfo(ctx, "Signing request completed")

	body, err := GetResponseBody(url, PostParams{
		"params":    params,
		"method":    reqCfg.Method,
		"reference": userCfg.Caller,
		"seed":      seed,
		"signature": signature,
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

// InfoResponse represents response from /info
type InfoResponse struct {
	Prototypes map[string]string `json:"prototypes"`
	RootDomain string            `json:"root_domain"`
	RootMember string            `json:"root_member"`
}

// Info sends request to /info and return result
func Info(url string) (*InfoResponse, error) {
	body, err := GetResponseBody(url+"/info", PostParams{})
	if err != nil {
		return nil, errors.Wrap(err, "[ Info ] problem with sending request:")
	}

	infoResp := InfoResponse{}

	err = json.Unmarshal(body, &infoResp)
	if err != nil {
		return nil, errors.Wrap(err, "[ Info ] problem with unmarshal response:")
	}

	return &infoResp, nil
}
