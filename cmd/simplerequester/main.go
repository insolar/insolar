package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
)

type response struct {
	Error   string
	Result  interface{}
	TraceID string
}

const defaultURL = "http://localhost:19101/api"

func parseInputParams() {
	pflag.StringVarP(&memberKeysPath, "memberkeys", "k", "", "path to file with Member keys")
	pflag.StringVarP(&apiURL, "url", "u", defaultURL, "api url")
	pflag.StringVarP(&paramsFile, "paramsFile", "f", "", "json file params")

	pflag.StringVarP(&method, "method", "m", "", "Insolar method name")
	pflag.StringVarP(&params, "params", "p", "", "params JSON")
	pflag.StringVarP(&memberRef, "address", "i", "", "Insolar member ref")

	pflag.Parse()
}

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

func main() {
	parseInputParams()

	err := log.SetLevel("error")
	check("can't set 'error' level on logger: ", err)
	datas := &DataToSign{}

	if paramsFile != "" {
		datas, err = readRequestParams(paramsFile)
		check("[ simpleRequester ]", err)
	} else {
		if memberRef == "" {
			response, err := requester.Info(apiURL)
			check("[ simpleRequester ]", err)
			datas.Reference = response.RootMember
		} else {
			datas.Reference = memberRef
		}
		if method != "" {
			datas.Method = method
		}
		if params != "" {
			datas.Params = params
		}
	}
	if datas.Method == "" {
		fmt.Println("Method cannot be null", err)
		os.Exit(1)
	}

	seed, err := requester.GetSeed(apiURL)
	check("[ simpleRequester ]", err)
	datas.Seed = seed

	rawConf, err := ioutil.ReadFile(memberKeysPath)
	check("[ simpleRequester ]", err)

	fmt.Println("Params: " + datas.Params)
	fmt.Println("Method: " + datas.Method)
	fmt.Println("Reference: " + datas.Reference)
	fmt.Println("Seed: " + datas.Seed)

	keys := &memberKeys{}
	err = json.Unmarshal(rawConf, keys)
	check("[ simpleRequester ]", err)

	jws, jwk, err := createSignedData(keys, datas)
	check("[ simpleRequester ]", err)
	params := requester.PostParams{
		"jws": jws,
		"jwk": jwk,
	}
	fmt.Println("JWS: " + jws)
	fmt.Println("JWK: " + jwk)
	body, err := requester.GetResponseBody(apiURL+"/call", params)
	check("[ simpleRequester ]", err)
	if len(body) == 0 {
		fmt.Println("[ simpleRequester ] Response body is Empty")
	}

	response, err := getResponse(body)
	check("[ simpleRequester ]", err)

	fmt.Println("Execute result: ", response)
}
