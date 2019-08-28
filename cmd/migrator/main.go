package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/log"
	"github.com/spf13/pflag"
)

var (
	indexStr       string
	apiURL         string
	paramsFile     string
	memberKeysPath string
)

const defaultURL = "http://localhost:19001/admin-api"

func parseInputParams() {
	pflag.StringVarP(&memberKeysPath, "memberkeys", "k", "../.artifacts/launchnet/configs/", "path to files with migration Member keys")
	pflag.StringVarP(&indexStr, "index", "i", "", "index of migration deamon")
	pflag.StringVarP(&apiURL, "adminurl", "u", defaultURL, "admin api url")
	pflag.StringVarP(&paramsFile, "paramsFile", "f", "migration.json", "json file params")
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
	if paramsFile != "" {
		request, err = readRequestParams(paramsFile)
		check("[ Migrator ]", err)
	}

	index, err := strconv.Atoi(indexStr)
	check("[ Migrator ]", err)

	if request.Params.Reference == "" {
		response, err := requester.Info(apiURL)
		check("[ Migrator ]", err)
		request.Params.Reference = response.MigrationDaemonMembers[index]
	}

	b, _ := json.Marshal(request)
	fmt.Println(string(b))

	if request.Params.CallSite == "" {
		fmt.Println("Method cannot be null", err)
		os.Exit(1)
	}

	rawConf, err := ioutil.ReadFile(memberKeysPath + "migration_daemon_" + strconv.Itoa(index) + "_member_keys.json")
	check("[ Migrator ]", err)

	stringParams, _ := json.Marshal(request.Params.CallParams)
	fmt.Println("callParams: " + string(stringParams))
	fmt.Println("Method: " + request.Method)
	fmt.Println("Reference: " + request.Params.Reference)

	keys := &memberKeys{}
	err = json.Unmarshal(rawConf, keys)
	check("[ Migrator ] failed to unmarshal", err)
	response, err := execute(apiURL, *keys, *request)
	check("[ Migrator ] failed to execute", err)
	fmt.Println("Execute result: \n", response)
}
