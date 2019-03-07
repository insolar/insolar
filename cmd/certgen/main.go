/*
 *    Copyright 2019 Insolar Technologies
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
	"io"
	"os"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/spf13/pflag"
)

const defaultURL = "http://localhost:19101/api"

var ks = platformpolicy.NewKeyProcessor()

var (
	role        string
	api         string
	keysFile    string
	certFileOut string
	rootConfig  string
	verbose     bool
	reuseKeys   bool
)

func parseInputParams() {
	pflag.StringVarP(&role, "role", "r", "virtual", "The role of the new node")
	pflag.StringVarP(&api, "url", "h", defaultURL, "Insolar API URL")
	pflag.BoolVarP(&verbose, "verbose", "v", false, "Be verbose (default false)")
	pflag.BoolVarP(&reuseKeys, "reuse-keys", "u", false, "Read keys from file instead og generating of new ones")
	pflag.StringVarP(&keysFile, "keys-file", "k", "keys.json", "The OUT/IN ( depends on 'reuse-keys' ) file for public/private keys of the node")
	pflag.StringVarP(&certFileOut, "cert-file", "c", "cert.json", "The OUT file the node certificate")
	pflag.StringVarP(&rootConfig, "root-conf", "t", "", "Config that contains public/private keys of root member")
	pflag.Parse()
}

func generateKeys() (crypto.PublicKey, crypto.PrivateKey) {
	privKey, err := ks.GeneratePrivateKey()
	checkError("Failed to generate private key:", err)
	pubKey := ks.ExtractPublicKey(privKey)
	fmt.Println("Generate reys")
	return pubKey, privKey
}

func loadKeys() (crypto.PublicKey, crypto.PrivateKey) {
	keyStore, err := keystore.NewKeyStore(keysFile)
	checkError("Failed to laod keys", err)

	privKey, err := keyStore.GetPrivateKey("")
	checkError("Failed to GetPrivateKey", err)

	fmt.Println("Load keys")
	return ks.ExtractPublicKey(privKey), privKey
}

func getKeys() (crypto.PublicKey, crypto.PrivateKey) {
	if reuseKeys {
		return loadKeys()
	}
	return generateKeys()
}

type RegisterResult struct {
	Result  string `json:"result"`
	TraceID string `json:"traceID"`
}

func extractReference(response []byte, requestTypeMsg string) core.RecordRef {
	r := RegisterResult{}
	err := json.Unmarshal(response, &r)
	checkError(fmt.Sprintf("Failed to parse response from '%s' node request", requestTypeMsg), err)
	if verbose {
		fmt.Println("Response:", string(response))
	}

	ref, err := core.NewRefFromBase58(r.Result)
	checkError(fmt.Sprintf("Failed to construct ref from '%s' node response", requestTypeMsg), err)

	return *ref
}

func registerNode(key crypto.PublicKey, staticRole core.StaticRole) core.RecordRef {
	userCfg := getUserConfig()

	keySerialized, err := ks.ExportPublicKeyPEM(key)
	checkError("Failed to export public key:", err)
	request := requester.RequestConfigJSON{
		Method: "RegisterNode",
		Params: []interface{}{keySerialized, staticRole.String()},
	}

	ctx := inslogger.ContextWithTrace(context.Background(), "insolarUtility")
	response, err := requester.Send(ctx, api, userCfg, &request)
	checkError("Failed to execute register node request", err)

	fmt.Println("Register node")
	return extractReference(response, "registerNode")
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

func writeKeys(pubKey crypto.PublicKey, privKey crypto.PrivateKey) {
	privKeyStr, err := ks.ExportPrivateKeyPEM(privKey)
	checkError("Failed to deserialize private key:", err)

	pubKeyStr, err := ks.ExportPublicKeyPEM(pubKey)
	checkError("Failed to deserialize public key:", err)

	result, err := json.MarshalIndent(map[string]interface{}{
		"private_key": string(privKeyStr),
		"public_key":  string(pubKeyStr),
	}, "", "    ")
	checkError("Failed to serialize file with private/public keys:", err)
	f, err := openFile(keysFile)
	checkError("Failed to open file with private/public keys:", err)
	_, err = f.Write([]byte(result))
	checkError("Failed to write file with private/public keys:", err)
}

func writeCertificate(cert []byte) {
	f, err := openFile(certFileOut)
	checkError("Failed to open file with certificate:", err)
	_, err = f.Write(cert)
	checkError("Failed to write file with certificate:", err)
}

func checkError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, ": ", err)
		os.Exit(1)
	}
}

func openFile(path string) (io.Writer, error) {
	return os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}

func getUserConfig() *requester.UserConfigJSON {
	requester.SetVerbose(verbose)
	userCfg, err := requester.ReadUserConfigFromFile(rootConfig)
	checkError("Failed to read root config:", err)
	info, err := requester.Info(api)
	checkError("Failed to execute info request to API:", err)
	userCfg.Caller = info.RootMember

	return userCfg
}

func getNodeRefByPk(key crypto.PublicKey) core.RecordRef {
	userCfg := getUserConfig()

	keySerialized, err := ks.ExportPublicKeyPEM(key)
	checkError("Failed to export public key:", err)
	request := requester.RequestConfigJSON{
		Method: "GetNodeRef",
		Params: []interface{}{keySerialized},
	}

	ctx := inslogger.ContextWithTrace(context.Background(), "insolarUtility")
	response, err := requester.Send(ctx, api, userCfg, &request)
	checkError("Failed to execute GetNodeRefByPK node request", err)

	fmt.Println("Extract node by PK")
	return extractReference(response, "getNodeRefByPk")
}

func getNodeRef(pubKey crypto.PublicKey, staticRole core.StaticRole) core.RecordRef {
	if reuseKeys {
		return getNodeRefByPk(pubKey)
	}

	return registerNode(pubKey, staticRole)
}

func main() {
	parseInputParams()
	staticRole := core.GetStaticRoleFromString(role)
	if staticRole == core.StaticRoleUnknown {
		fmt.Println("Invalid role:", role)
		os.Exit(1)
	}

	pub, priv := getKeys()
	ref := getNodeRef(pub, staticRole)
	cert := fetchCertificate(ref)

	if !reuseKeys {
		writeKeys(pub, priv)
		fmt.Println("Write keys")
	}
	writeCertificate(cert)
	fmt.Println("Write certificate")
}
