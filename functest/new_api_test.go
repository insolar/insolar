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

// +build functest

package functest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
)

func contractError(body []byte) error {
	var t map[string]interface{}
	err := json.Unmarshal(body, &t)
	if err != nil {
		return err
	}
	if e, ok := t["error"]; ok {
		if ee, ok := e.(string); ok && ee != "" {
			return errors.New(ee)
		}
	}
	return nil
}

func TestBadSeed(t *testing.T) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(root.ref, root.privKey)
	require.NoError(t, err)
	res, err := requester.SendWithSeed(ctx, TestCallUrl, rootCfg, &requester.RequestConfigJSON{
		Method: "CreateMember",
		Params: nil,
	}, []byte("111"))
	require.NoError(t, err)
	require.EqualError(t, contractError(res), "[ checkSeed ] Bad seed param")
}

func TestIncorrectSeed(t *testing.T) {
	ctx := context.TODO()
	rootCfg, err := requester.CreateUserConfig(root.ref, root.privKey)
	require.NoError(t, err)
	res, err := requester.SendWithSeed(ctx, TestCallUrl, rootCfg, &requester.RequestConfigJSON{
		Method: "CreateMember",
		Params: nil,
	}, []byte("12345678901234567890123456789012"))
	require.NoError(t, err)
	require.EqualError(t, contractError(res), "[ checkSeed ] Incorrect seed")
}

func customSend(data string) (map[string]interface{}, error) {
	req, err := http.NewRequest("POST", TestCallUrl, strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var out map[string]interface{}
	err = json.Unmarshal(body, &out)
	return out, err
}

func TestEmptyBody(t *testing.T) {
	res, err := customSend("")
	require.NoError(t, err)
	require.Equal(t, "[ UnmarshalRequest ] Empty body", res["error"])
}

func TestCrazyJSON(t *testing.T) {
	res, err := customSend("[dh")
	require.NoError(t, err)
	require.Contains(t, res["error"], "[ UnmarshalRequest ] Can't unmarshal input params: invalid")
}

func TestIncorrectSign(t *testing.T) {
	args, err := insolar.MarshalArgs(nil)
	require.NoError(t, err)
	seed, err := requester.GetSeed(TestAPIURL)
	require.NoError(t, err)
	body, err := requester.GetResponseBody(TestCallUrl, requester.PostParams{
		"params":    args,
		"method":    "SomeMethod",
		"reference": root.ref,
		"seed":      seed,
		"signature": []byte("1234567890"),
	})
	require.NoError(t, err)
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	require.NoError(t, err)
	require.Contains(t, res["error"], "Incorrect signature")
}

func TestExporter_ValidateResponse(t *testing.T) {
	ctx := context.Background()
	swagger, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile("testdata/api-exported.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = swagger.Validate(ctx)
	if err != nil {
		t.Fatal(err)
	}

	method := "exporter.Export"
	response, err := requester.GetRPCResponse(TestRPCUrl, method, map[string]int{
		"from": 0,
		"size": 10,
	})
	require.NoError(t, err)
	var j map[string]interface{}
	umErr := json.Unmarshal(response, &j)
	if umErr != nil {
		t.Fatal("failed unmarshal response:", umErr, "\n", string(response))
	}
	responseFmt, jerr := json.MarshalIndent(j, "", "    ")
	if jerr != nil {
		t.Fatal("failed marshal unmarshaled response:", string(response), jerr)
	}

	// ignore servers addresses validation
	swagger.Servers = nil
	router := openapi3filter.NewRouter().WithSwagger(swagger)
	route, pathParams, err := router.FindRoute("POST", &url.URL{
		Path: "/api/rpc#method=" + method,
	})
	if err != nil {
		t.Fatal("failed find route", err)
	}

	httpReq, _ := http.NewRequest(http.MethodPost, "/api/rpc", nil)
	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: requestValidationInput,
		Status:                 http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{
				"application/json", "charset=utf-8",
			},
		},
	}
	responseValidationInput.SetBodyBytes(response)

	err = openapi3filter.ValidateResponse(ctx, responseValidationInput)
	if err != nil {
		if _, ok := err.(*openapi3filter.ResponseError); ok {
			fmt.Print("got response:\n", string(responseFmt))
		}
	}
	require.NoError(t, err, "validate response failed")
}
