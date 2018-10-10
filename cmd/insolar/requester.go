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

package main

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/insolar/insolar/api/requesters"
	"github.com/insolar/insolar/core"

	"github.com/pkg/errors"

	ecdsa_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
)

type apiRequester struct {
}

type configJSON struct {
	PrivateKey       string `json:"private_key"`
	privateKeyObject *ecdsa.PrivateKey
	Params           []interface{} `json:"params"`
	Method           string        `json:"method"`
	Caller           string        `json:"caller"`
	Callee           string        `json:"callee"`
}

// type config struct {
// 	PrivateKey *ecdsa.PrivateKey
// }

const url = "http://localhost:19191/api/v1?"

func readConfigFromFile(path string) (*configJSON, error) {
	rawConf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "[ readConfigFromFile ] Problem with reading config")
	}

	cfgJSON := &configJSON{}
	err = json.Unmarshal(rawConf, &cfgJSON)
	if err != nil {
		return nil, errors.Wrap(err, "[ readConfigFromFile ] Problem with unmarshaling config")
	}

	cfgJSON.privateKeyObject, err = ecdsa_helper.ImportPrivateKey(cfgJSON.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ readConfigFromFile ] Problem with reading private key")
	}

	return cfgJSON, nil
}

func constructParams(params []interface{}) ([]byte, error) {
	args := []interface{}{}
	return
}

func (r *apiRequester) Send(out io.Writer, confPath string) error {
	cfg, err := readConfigFromFile(confPath)
	if err != nil {
		return errors.Wrap(err, "[ Send ] Problem with reading config")
	}

	seed, err := requesters.GetSeed(url)
	if err != nil {
		return errors.Wrap(err, "[ Send ] Problem with getting seed")
	}

	params, err := constructParams(cfg.Params)
	if err != nil {
		return errors.Wrap(err, "[ Send ] Problem with creating request")
	}

	signature, err := ecdsa_helper.Sign(params, cfg.privateKeyObject)
	if err != nil {
		return errors.Wrap(err, "[ Send ] Problem with signing request")

	}

	pubKey, err := ecdsa_helper.ExportPublicKey(&cfg.privateKeyObject.PublicKey)
	check("[ Send ] ", err)

	body, err := requesters.GetResponseBody(url, requesters.PostParams{
		"params":    base64.StdEncoding.EncodeToString(params),
		"method":    cfg.Method,
		"caller":    core.RandomRef().String(),
		"callee":    core.RandomRef().String(),
		"seed":      seed,
		"signature": ecdsa_helper.ExportSignature(signature),
		//
		"query_type": "create_member",
		"name":       "PUTIN",
		"public_key": pubKey,
	})

	if err != nil {
		return errors.Wrap(err, "[ Send ] Problem with sending target request")
	}

	writeToOutput(out, string(body))

	return nil
}
