package main

import (
	"encoding/json"
	"github.com/insolar/insolar/api/requester"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type memberKeys struct {
	Private string `json:"private_key"`
	Public  string `json:"public_key"`
}

func getResponse(body []byte) (*response, error) {
	res := &response{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.Wrap(err, "[ getResponse ] problems with unmarshal response")
	}
	return res, nil
}

func readRequestParams(path string) (*requester.Request, error) {

	fileParams := &requester.Request{}
	err := readFile(path, fileParams)
	if err != nil {
		return nil, errors.Wrap(err, "[ readRequestParams ] ")
	}

	return fileParams, nil
}

func readFile(path string, configType interface{}) error {
	var rawConf []byte
	var err error
	if path == "-" {
		rawConf, err = ioutil.ReadAll(os.Stdin)
	} else {
		rawConf, err = ioutil.ReadFile(filepath.Clean(path))
	}
	if err != nil {
		return errors.Wrap(err, "[ readFile ] Problem with reading config")
	}

	err = json.Unmarshal(rawConf, &configType)
	if err != nil {
		return errors.Wrap(err, "[ readFile ] Problem with unmarshaling config")
	}

	return nil
}
