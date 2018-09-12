/*
 *    Copyright 2018 INS Ecosystem
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
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func TestLaunchApi(t *testing.T) {
	cfg := configuration.NewApiRunner()
	api, err := NewApiRunner(cfg)
	assert.NoError(t, err)

	cs := core.Components{}
	err = api.Start(cs)
	assert.NoError(t, err)

	const TestUrl = "http://localhost:8080/api/v1?query_type=LOL"

	{
		resp, err := http.Get(TestUrl)
		assert.NoError(t, err)
		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(body[:]), `"message": "Bad request"`)
	}

	{
		postParams := map[string]string{"query_type": "get_balance", "reference": "test"}
		jsonValue, _ := json.Marshal(postParams)
		postResp, err := http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
		assert.NoError(t, err)
		body, err := ioutil.ReadAll(postResp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(body[:]), `"message": "Handler error: [ ProcessGetBalance ]: [ SendRequest ]: [ RouteCall ] message`)
	}

	{
		postParams := map[string]string{"query_type": "TEST", "reference": "test"}
		jsonValue, _ := json.Marshal(postParams)
		postResp, err := http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
		assert.NoError(t, err)
		body, err := ioutil.ReadAll(postResp.Body)
		assert.NoError(t, err)
		assert.Contains(t, string(body[:]), `"message": "Wrong query parameter 'query_type' = 'TEST'"`)
	}

	api.Stop()
	assert.NoError(t, err)
}

func TestSerialization(t *testing.T) {
	var a uint = 1
	var b bool = true
	var c string = "test"

	serArgs, err := MarshalArgs(a, b, c)
	assert.NoError(t, err)
	assert.NotNil(t, serArgs)

	var aR uint
	var bR bool
	var cR string
	rowResp, err := UnMarshalResponse(serArgs, []interface{}{aR, bR, cR})
	assert.NoError(t, err)
	assert.Len(t, rowResp, 3)
	assert.Equal(t, reflect.TypeOf(a), reflect.TypeOf(rowResp[0]))
	assert.Equal(t, reflect.TypeOf(b), reflect.TypeOf(rowResp[1]))
	assert.Equal(t, reflect.TypeOf(c), reflect.TypeOf(rowResp[2]))
	assert.Equal(t, a, rowResp[0].(uint))
	assert.Equal(t, b, rowResp[1].(bool))
	assert.Equal(t, c, rowResp[2].(string))
}

func TestGetQid(t *testing.T) {
	exists := make(map[string]bool)
	const NumIters = 1500
	for i := 0; i < NumIters; i++ {
		exists[GenQID()] = true
	}
	assert.Len(t, exists, NumIters)
}

func TestNewApiRunnerNilConfig(t *testing.T) {
	_, err := NewApiRunner(nil)
	assert.EqualError(t, err, "[ NewApiRunner ] config is nil")
}

func TestNewApiRunnerNoRequiredParams(t *testing.T) {
	cfg := configuration.ApiRunner{}
	_, err := NewApiRunner(&cfg)
	assert.EqualError(t, err, "[ NewApiRunner ] Port must not be 0")

	cfg.Port = 100
	_, err = NewApiRunner(&cfg)
	assert.EqualError(t, err, "[ NewApiRunner ] Location must exist")

	cfg.Location = "test"
	_, err = NewApiRunner(&cfg)
	assert.NoError(t, err)
}
