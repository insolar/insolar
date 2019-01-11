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
	"context"
	"crypto"
	"encoding/json"
	"fmt"
	"os"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

const defaultURL = "http://localhost:19191/api"

var ks = platformpolicy.NewKeyProcessor()

var (
	role        string
	api         string
	keysFileOut string
	certFileOut string
	rootConfig  string
	verbose     bool
)

func parseInputParams() {
	pflag.StringVarP(&role, "role", "r", "virtual", "The role of the new node")
	pflag.StringVarP(&api, "url", "h", defaultURL, "Insolar API URL")
	pflag.BoolVarP(&verbose, "verbose", "v", false, "Be verbose (default false)")
	pflag.StringVarP(&keysFileOut, "keys-file", "k", "keys.json", "The OUT file for public/private keys of the node")
	pflag.StringVarP(&certFileOut, "cert-file", "c", "cert.json", "The OUT file the node certificate")
	pflag.StringVar(&rootConfig, "root-conf", "", "Config that contains public/private keys of root member")
	pflag.Parse()
}

func generateKeys() (crypto.PublicKey, crypto.PrivateKey) {
	privKey, err := ks.GeneratePrivateKey()
	checkError("Failed to generate private key:", err)
	pubKey := ks.ExtractPublicKey(privKey)
	return pubKey, privKey
}

type RegisterResult struct {
	Result  string `json:"result"`
	TraceID string `json:"traceID"`
}

func registerNode(key crypto.PublicKey, staticRole core.StaticRole) core.RecordRef {
	requester.SetVerbose(verbose)
	userCfg, err := requester.ReadUserConfigFromFile(rootConfig)
	checkError("Failed to read root config:", err)
	info, err := requester.Info(api)
	checkError("Failed to execute info request to API:", err)
	userCfg.Caller = info.RootMember

	keySerialized, err := ks.ExportPublicKeyPEM(key)
	checkError("Failed to export public key:", err)
	request := requester.RequestConfigJSON{
		Method: "RegisterNode",
		Params: []interface{}{keySerialized, staticRole.String()},
	}

	ctx := inslogger.ContextWithTrace(context.Background(), "insolarUtility")
	response, err := requester.Send(ctx, api, userCfg, &request)
	checkError("Failed to execute register node request", err)

	r := RegisterResult{}
	err = json.Unmarshal(response, &r)
	checkError("Failed to parse response from register node request", err)
	ref, err := core.NewRefFromBase58(r.Result)
	checkError("Failed to construct ref from register node response", err)
	return *ref
}

type GetCertificateResult struct {
	Cert json.RawMessage `json:"cert"`
}

type GetCertificateResponse struct {
	Version string               `json:"jsonrpc"`
	ID      string               `json:"id"`
	Result  GetCertificateResult `json:"result"`
}

func fetchCertificate(ref core.RecordRef) []byte {
	params := requester.PostParams{
		"ref": ref.String(),
	}
	response, err := requester.GetResponseBody(api+"/rpc", requester.PostParams{
		"jsonrpc": "2.0",
		"method":  "cert.Get",
		"id":      "",
		"params":  params,
	})
	checkError("Failed to get certificate for the registered node:", err)

	r := GetCertificateResponse{}
	err = json.Unmarshal(response, &r)
	checkError("Failed to parse response from get certificate request:", err)

	cert, err := r.Result.Cert.MarshalJSON()
	checkError("Failed to marshal certificate from API response:", err)
	return cert
}

func writeKeys(crypto.PublicKey, crypto.PrivateKey) error {
	return errors.New("not implemented")
}

func writeCertificate(cert []byte) error {
	return errors.New("not implemented")
}

func checkError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {
	parseInputParams()
	staticRole := core.GetStaticRoleFromString(role)
	if staticRole == core.StaticRoleUnknown {
		fmt.Println("Invalid role:", role)
		os.Exit(1)
	}

	pub, priv := generateKeys()
	ref := registerNode(pub, staticRole)
	cert := fetchCertificate(ref)

	err := writeKeys(pub, priv)
	checkError("Failed to write file with public/private keys:", err)

	err = writeCertificate(cert)
	checkError("Failed to write file with node certificate:", err)

	fmt.Println("Successfully generated files with keys and certificate")
}
