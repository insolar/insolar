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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"

	ecdsa_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
)

// verbose switches on verbose mode
var verbose = false

func verboseInfo(msg string) {
	if verbose {
		log.Infoln(msg)
	}
}

// SetVerbose switchs on verbose mode
func SetVerbose(verb bool) {
	verbose = verb
}

// PostParams represents params struct
type PostParams = map[string]string

// GetResponseBody makes request and extracts body
func GetResponseBody(url string, postP map[string]string) ([]byte, error) {
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
func GetSeed(url string) (string, error) {
	body, err := GetResponseBody(url, PostParams{
		"query_type": "get_seed",
	})
	if err != nil {
		return "", errors.Wrap(err, "[ getSeed ]")
	}

	type seedResponse struct{ Seed string }
	seedResp := seedResponse{}

	err = json.Unmarshal(body, &seedResp)
	if err != nil {
		return "", errors.Wrap(err, "[ getSeed ]")
	}

	return seedResp.Seed, nil
}

func constructParams(params []interface{}) ([]byte, error) {
	args, err := core.MarshalArgs(params)
	if err != nil {
		return nil, errors.Wrap(err, "[ constructParams ]")
	}
	return args, nil
}

// SendWithSeed sends request with known seed
func SendWithSeed(url string, userCfg *UserConfigJSON, reqCfg *RequestConfigJSON, seed string) ([]byte, error) {
	if userCfg == nil || reqCfg == nil {
		return nil, errors.New("[ Send ] Configs must be initialized")
	}

	params, err := constructParams(reqCfg.Params)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with serializing params")
	}

	verboseInfo("Signing params ...")
	signature, err := ecdsa_helper.Sign(params, userCfg.privateKeyObject)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with signing request")
	}
	verboseInfo("Signing params completed")

	body, err := GetResponseBody(url, PostParams{
		"params":    string(params),
		"method":    reqCfg.Method,
		"requester": reqCfg.Requester,
		"target":    reqCfg.Target,
		"delegate":  reqCfg.Delegate,
		"seed":      seed,
		"signature": ecdsa_helper.ExportSignature(signature),
	})

	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with sending target request")
	}

	return body, nil
}

// Send: first gets seed and after that makes target request
func Send(url string, userCfg *UserConfigJSON, reqCfg *RequestConfigJSON) ([]byte, error) {
	verboseInfo("Sending GETSEED request ...")
	seed, err := GetSeed(url)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ] Problem with getting seed")
	}
	verboseInfo("GETSEED request completed. seed: " + seed)

	response, err := SendWithSeed(url, userCfg, reqCfg, seed)
	if err != nil {
		return nil, errors.Wrap(err, "[ Send ]")
	}

	return response, nil
}
