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

package main

import (
	"context"
	"crypto"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/platformpolicy"
)

func (g *certGen) generateKeys() {
	privKey, err := g.keyProcessor.GeneratePrivateKey()
	checkError("Failed to generate private key:", err)
	pubKey := g.keyProcessor.ExtractPublicKey(privKey)
	fmt.Println("Generate keys")
	g.pubKey, g.privKey = pubKey, privKey
}

func (g *certGen) loadKeys() {
	keyStore, err := keystore.NewKeyStore(g.keysFileOut)
	checkError("Failed to laod keys", err)

	g.privKey, err = keyStore.GetPrivateKey("")
	checkError("Failed to GetPrivateKey", err)

	fmt.Println("Load keys")
	g.pubKey = g.keyProcessor.ExtractPublicKey(g.privKey)
}

type RegisterResult struct {
	Result  string `json:"result"`
	TraceID string `json:"traceID"`
}

func extractReference(response []byte, requestTypeMsg string) insolar.Reference {
	r := RegisterResult{}
	err := json.Unmarshal(response, &r)
	checkError(fmt.Sprintf("Failed to parse response from '%s' node request", requestTypeMsg), err)
	if verbose {
		fmt.Println("Response:", string(response))
	}

	ref, err := insolar.NewReferenceFromBase58(r.Result)
	checkError(fmt.Sprintf("Failed to construct ref from '%s' node response", requestTypeMsg), err)

	return *ref
}

func (g *certGen) registerNode() insolar.Reference {
	userCfg := g.getUserConfig()

	keySerialized, err := g.keyProcessor.ExportPublicKeyPEM(g.pubKey)
	checkError("Failed to export public key:", err)
	request := requester.RequestConfigJSON{
		Method: "RegisterNode",
		Params: []interface{}{keySerialized, g.staticRole.String()},
	}

	ctx := inslogger.ContextWithTrace(context.Background(), "insolarUtility")
	response, err := requester.Send(ctx, g.API, userCfg, &request)
	checkError("Failed to execute register node request", err)

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

func (g *certGen) fetchCertificate(ref insolar.Reference) []byte {
	params := requester.PostParams{
		"ref": ref.String(),
	}
	response, err := requester.GetResponseBody(g.API+"/rpc", requester.PostParams{
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

func (g *certGen) writeKeys() {
	privKeyStr, err := g.keyProcessor.ExportPrivateKeyPEM(g.privKey)
	checkError("Failed to deserialize private key:", err)

	pubKeyStr, err := g.keyProcessor.ExportPublicKeyPEM(g.pubKey)
	checkError("Failed to deserialize public key:", err)

	result, err := json.MarshalIndent(map[string]interface{}{
		"private_key": string(privKeyStr),
		"public_key":  string(pubKeyStr),
	}, "", "    ")
	checkError("Failed to serialize file with private/public keys:", err)

	f, err := openFile(g.keysFileOut)
	checkError("Failed to open file with private/public keys:", err)

	_, err = f.Write([]byte(result))
	checkError("Failed to write file with private/public keys:", err)

	fmt.Println("Write keys to", g.keysFileOut)
}

func (g *certGen) writeCertificate(cert []byte) {
	f, err := openFile(g.certFileOut)
	checkError("Failed to open file with certificate:", err)

	_, err = f.Write(cert)
	checkError("Failed to write file with certificate:", err)

	fmt.Println("Write certificate to", g.certFileOut)
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

func (g *certGen) getUserConfig() *requester.UserConfigJSON {
	requester.SetVerbose(verbose)
	userCfg, err := requester.ReadUserConfigFromFile(g.rootKeysFile)
	checkError("Failed to read root config:", err)
	info, err := requester.Info(g.API)
	checkError("Failed to execute info request to API:", err)
	userCfg.Caller = info.RootMember

	return userCfg
}

func (g *certGen) getNodeRefByPk() insolar.Reference {
	userCfg := g.getUserConfig()

	keySerialized, err := g.keyProcessor.ExportPublicKeyPEM(g.privKey)
	checkError("Failed to export public key:", err)
	request := requester.RequestConfigJSON{
		Method: "GetNodeRef",
		Params: []interface{}{keySerialized},
	}

	ctx := inslogger.ContextWithTrace(context.Background(), "insolarUtility")
	response, err := requester.Send(ctx, g.API, userCfg, &request)
	checkError("Failed to execute GetNodeRefByPK node request", err)

	fmt.Println("Extract node by PK")
	return extractReference(response, "getNodeRefByPk")
}

type certGen struct {
	keyProcessor insolar.KeyProcessor

	rootKeysFile string
	API          string
	staticRole   insolar.StaticRole

	keysFileOut string
	certFileOut string

	pubKey  crypto.PublicKey
	privKey crypto.PrivateKey
}

func genCertificate(
	rootKeysFile string,
	role string,
	url string,
	keysFile string,
	certFile string,
	reuseKeys bool,
) {
	staticRole := insolar.GetStaticRoleFromString(role)
	if staticRole == insolar.StaticRoleUnknown {
		fmt.Println("Invalid role:", role)
		os.Exit(1)
	}

	g := &certGen{
		keyProcessor: platformpolicy.NewKeyProcessor(),
		rootKeysFile: rootKeysFile,
		API:          url,
		staticRole:   staticRole,
		keysFileOut:  keysFile,
		certFileOut:  certFile,
	}

	var ref insolar.Reference
	if reuseKeys {
		g.loadKeys()
		ref = g.getNodeRefByPk()
	} else {
		g.generateKeys()
		ref = g.registerNode()
		fmt.Println("Register node", ref)
	}

	cert := g.fetchCertificate(ref)

	if !reuseKeys {
		g.writeKeys()
	}
	g.writeCertificate(cert)
}
