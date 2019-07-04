package main

import (
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/insolar/x-crypto/elliptic"

	"github.com/insolar/x-crypto/x509"

	"github.com/insolar/x-crypto/ecdsa"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/log"
	"github.com/spf13/pflag"
)

var (
	apiURL         string
	method         string
	params         string
	memberRef      string
	paramsFile     string
	memberKeysPath string
	privateKeyHex  string
)

const defaultURL = "http://localhost:19101/api"

func parseInputParams() {
	fmt.Println("Parse data")
	pflag.StringVarP(&memberKeysPath, "memberkeys", "k", "", "path to file with Member keys")
	pflag.StringVarP(&privateKeyHex, "privateKeyHex", "h", "", "private key in hex foramt")
	pflag.StringVarP(&apiURL, "url", "u", defaultURL, "api url")
	pflag.StringVarP(&paramsFile, "paramsFile", "f", "", "json file params")

	pflag.StringVarP(&method, "method", "m", "", "Insolar method name")
	pflag.StringVarP(&params, "params", "p", "", "params JSON")
	pflag.StringVarP(&memberRef, "address", "i", "", "Insolar member ref")

	pflag.Parse()
}

func main() {
	parseInputParams()
	err := log.SetLevel("error")
	check("can't set 'error' level on logger: ", err)
	request := &requester.Request{
		JSONRPC: "2.0",
		ID:      0,
		Method:  "api.call",
	}

	if paramsFile != "" {
		request, err = readRequestParams(paramsFile)
		check("[ simpleRequester ]", err)
	} else {
		if method != "" {
			request.Params.CallSite = method
		}
		if params != "" {
			request.Params.CallParams = params
		}
		request.Params.Reference = memberRef
	}

	if request.Params.Reference == "" {
		response, err := requester.Info(apiURL)
		check("[ simpleRequester ]", err)
		request.Params.Reference = response.RootMember
	}

	if request.Params.CallSite == "" {
		fmt.Println("Method cannot be null", err)
		os.Exit(1)
	}

	if len(memberKeysPath) > 0 {
		rawConf, err := ioutil.ReadFile(memberKeysPath)
		check("[ simpleRequester ]", err)

		stringParams, _ := json.Marshal(request.Params.CallParams)
		fmt.Println("callParams: " + string(stringParams))
		fmt.Println("Method: " + request.Method)
		fmt.Println("Reference: " + request.Params.Reference)

		keys := &memberKeys{}
		err = json.Unmarshal(rawConf, keys)
		check("[ simpleRequester ] failed to unmarshal", err)
		response, err := execute(apiURL, *keys, *request)
		check("[ simpleRequester ] failed to execute", err)
		fmt.Println("Execute result: \n", response)
	}
	if len(privateKeyHex) > 0 {
		stringParams, _ := json.Marshal(request.Params.CallParams)
		fmt.Println("callParams: " + string(stringParams))
		fmt.Println("Method: " + request.Method)
		fmt.Println("Reference: " + request.Params.Reference)

		i := new(big.Int)
		i.SetString(privateKeyHex, 16)

		privateKey := new(ecdsa.PrivateKey)
		privateKey.PublicKey.Curve = elliptic.P256K()
		if 8*len(i.Bytes()) != privateKey.Params().BitSize {
			fmt.Println("invalid length, need %d bits", privateKey.Params().BitSize)
			os.Exit(1)
		}

		privateKey.D = new(big.Int).SetBytes(i.Bytes())

		privateKey.PublicKey.X, privateKey.PublicKey.Y = privateKey.PublicKey.Curve.ScalarBaseMult(i.Bytes())

		publicKey := &privateKey.PublicKey

		pemPriv := exportPrivateToPem(privateKey)
		check("[ simpleRequester ] failed to convert private to pem", err)
		pemPub := exportPublicToPem(publicKey)
		check("[ simpleRequester ] failed to convert public to pem", err)

		keys := &memberKeys{Private: pemPriv, Public: pemPub}
		check("[ simpleRequester ] failed to unmarshal", err)
		response, err := execute(apiURL, *keys, *request)
		check("[ simpleRequester ] failed to execute", err)
		fmt.Println("Execute result: \n", response)
	} else {
		fmt.Println("Private key cannot be null", err)
		os.Exit(1)
	}
}

func exportPrivateToPem(ecdsaPrivateKey *ecdsa.PrivateKey) string {
	x509Encoded, err := x509.MarshalECPrivateKey(ecdsaPrivateKey)
	if err != nil {
		return ""
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})
	return string(pemEncoded)
}

func exportPublicToPem(publicKey *ecdsa.PublicKey) string {
	x509EncodedPub, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return ""
	}
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})
	return string(pemEncoded)
}
