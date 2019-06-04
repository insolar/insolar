package simplerequester

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
	memberKeysPath string
	apiURL         string
	paramsFile     string
	method         string
	params         string
	memberRef      string
)

type response struct {
	Error   string
	Result  interface{}
	TraceID string
}

type RequestParams struct {
	Reference string `json:"reference"`
	Method    string `json:"reference"`
	Params    string `json:"reference"`
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
	var datas *DataToSign

	if paramsFile != "" {
		datas, err = ReadRequestParams(paramsFile)
		check("[ simpleRequester ]", err)
	} else {
		datas.Reference = memberRef
		datas.Method = method
		datas.Params = params
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

	keys := memberKeys{}
	err = json.Unmarshal(rawConf, &keys)
	check("[ simpleRequester ]", err)

	privateKey, err := importPrivateKeyPEM([]byte(keys.Private))

	jws, jwk, err := createSignedData(privateKey, datas)
	check("[ simpleRequester ]", err)
	params := requester.PostParams{
		"jws": jws,
		"jwk": jwk,
	}
	body, err := requester.GetResponseBody(apiURL+"/call", params)
	check("[ simpleRequester ]", err)

	response, err := getResponse(body)
	check("[ simpleRequester ]", err)

	fmt.Println("Execute result: ", response)
}
