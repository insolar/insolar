package main

import (
	"encoding/json"
	"fmt"
	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/log"
	"github.com/spf13/pflag"
	"io/ioutil"
)

var (
	apiURL         string
	method         string
	params         string
	memberRef      string
	paramsFile     string
	memberKeysPath string
)

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

func main() {
	parseInputParams()
	err := log.SetLevel("error")
	check("can't set 'error' level on logger: ", err)
	request := &requester.ContractRequest{
		Request: requester.Request{
			Version: "2.0",
			ID:      0,
			Method:  "contract.call",
		},
	}
	params := requester.Params{}

	if paramsFile != "" {
		request, err = readRequestParams(paramsFile)
		check("[ simpleRequester ]", err)
	} else {
		if method != "" {
			params.CallSite = method
		}
		params.Reference = memberRef
	}

	if params.Reference == "" {
		response, err := requester.Info(apiURL)
		check("[ simpleRequester ]", err)
		params.Reference = response.RootMember
	}
	// fmt.Println()
	// if params.CallSite == "" {
	// 	fmt.Println("Method cannot be null", err)
	// 	os.Exit(1)
	// }

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
