// Copyright 2020 Insolar Network Ltd.
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

package main

import (
	"context"
	"crypto"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

var (
	memberKeysPath       string
	apiURL               string
	inputRequestParams   string
	shouldPasteSeed      bool
	shouldPastePublicKey bool
	memberPrivateKey     crypto.PrivateKey
	request              *requester.ContractRequest
)

func parseInputParams() {
	pflag.StringVarP(&memberKeysPath, "memberkeys", "k", "", "Path to member key")
	pflag.StringVarP(&apiURL, "url", "u", "", "API URL. for example http://localhost:19101/api/rpc")
	pflag.StringVarP(&inputRequestParams, "request", "r", "", "The request body or path to request params file")
	pflag.BoolVarP(&shouldPasteSeed, "autocompleteseed", "s", false, "Should replace seed to correct value")
	pflag.BoolVarP(&shouldPastePublicKey, "autocompletekey", "p", false, "Should replace publicKey to correct value")
	pflag.Parse()
}

func verifyParams() {
	if len(apiURL) > 0 {
		ok, err := isUrl(apiURL)
		if !ok {
			log.Fatal("URL parameter is incorrect. ", err)
		}
	}

	// verify that the member keys paramsFile is exist
	if !isFileExists(memberKeysPath) {
		log.Fatal("Member keys does not exists")
	}

	// try to read keys
	keys, err := secrets.ReadXCryptoKeysFile(memberKeysPath, false)
	if err != nil {
		log.Fatal("Cannot parse member keys. ", err)
	}
	memberPrivateKey = keys.Private

	if len(inputRequestParams) == 0 {
		log.Fatal("Request parameters cannot be empty.")
	}
	if isFileExists(inputRequestParams) {
		fileContent, err := ioutil.ReadFile(inputRequestParams)
		if err != nil {
			log.Fatal("Cannot read request. ", err)
		}
		// save to inputRequestParams if we could read params file for unmarshalling
		inputRequestParams = string(fileContent)
	}

	err = json.Unmarshal([]byte(inputRequestParams), &request)
	if err != nil {
		log.Fatal("Cannot unmarshal request. ["+inputRequestParams+"]", err)
	}
}

func isUrl(str string) (bool, error) {
	parsedUrl, err := url.Parse(str)
	return err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != "", err
}

func isFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	parseInputParams()
	verifyParams()

	userConfig, e := createUserConfig(memberPrivateKey)
	if e != nil {
		log.Fatal(e)
	}
	if shouldPastePublicKey {
		request.Params.PublicKey = userConfig.PublicKey
	}

	var response []byte
	var err error
	if shouldPasteSeed {
		response, err = requester.Send(context.Background(), apiURL, userConfig, &request.Params)
	} else {
		response, err = requester.SendWithSeed(context.Background(), apiURL, userConfig, &request.Params, request.Params.Seed)
	}

	if err != nil {
		log.Fatal(err)
	}

	print(string(response))
}

func createUserConfig(privateKey crypto.PrivateKey) (*requester.UserConfigJSON, error) {
	privateKeyBytes, err := secrets.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to export private key")
	}
	privateKeyStr := string(privateKeyBytes)

	publicKey, err := secrets.ExportPublicKeyPEM(secrets.ExtractPublicKey(privateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract public key")
	}
	publicKeyStr := string(publicKey)

	return requester.CreateUserConfig("", privateKeyStr, publicKeyStr)
}
