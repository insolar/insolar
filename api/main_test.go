/*
 *    Copyright 2018 Insolar
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

package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

const HOST = "http://localhost:19191"
const TestUrl = HOST + "/api/v1?query_type=LOL"

func TestMain(m *testing.M) {
	cfg := configuration.NewAPIRunner()
	bootstrapCfg := configuration.NewConfiguration()
	api, _ := NewRunner(&cfg)

	cs := core.Components{}
	b, _ := bootstrap.NewBootstrapper(bootstrapCfg)
	cs["core.Bootstrapper"] = b
	api.Start(cs)

	code := m.Run()

	api.Stop()

	os.Exit(code)
}

func TestWrongQueryParam(t *testing.T) {
	postParams := map[string]string{"query_type": "TEST", "reference": "test"}
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body[:]), `"message": "Wrong query parameter 'query_type' = 'TEST'"`)
}

func TestHandlerError(t *testing.T) {
	postParams := map[string]string{"query_type": "get_balance", "reference": "test"}
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestUrl, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body[:]), `"message": "Handler error: [ ProcessGetBalance ]: [ SendRequest ]: [ RouteCall ] message`)
}

func TestBadRequest(t *testing.T) {
	resp, err := http.Get(TestUrl)
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body[:]), `"message": "Bad request"`)
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

func TestNewApiRunnerNilConfig(t *testing.T) {
	_, err := NewRunner(nil)
	assert.EqualError(t, err, "[ NewAPIRunner ] config is nil")
}

func TestNewApiRunnerNoRequiredParams(t *testing.T) {
	cfg := configuration.APIRunner{}
	_, err := NewRunner(&cfg)
	assert.EqualError(t, err, "[ NewAPIRunner ] Port must not be 0")

	cfg.Port = 100
	_, err = NewRunner(&cfg)
	assert.EqualError(t, err, "[ NewAPIRunner ] Location must exist")

	cfg.Location = "test"
	_, err = NewRunner(&cfg)
	assert.NoError(t, err)
}

type TestsMessageRouter struct {
}

func (ar *TestsMessageRouter) Start(c core.Components) error {
	return nil
}

func (ar *TestsMessageRouter) Stop() error {
	return nil
}

const TestBalance = 100500

func (ar *TestsMessageRouter) Route(msg core.Message) (core.Response, error) {
	data, _ := MarshalArgs(TestBalance)

	resp := core.Response{
		Result: data,
	}

	return resp, nil
}

func TestWithFakeMessageRouter(t *testing.T) {
	mr := TestsMessageRouter{}

	const LOCATION = "/test/test"

	fw := wrapAPIV1Handler(&mr, core.RecordRef{})
	http.HandleFunc(LOCATION, fw)

	const TestUrl2 = HOST + LOCATION + "?query_type=PPPPPPPP"

	postParams := map[string]string{"query_type": "get_balance", "reference": "test"}
	jsonValue, _ := json.Marshal(postParams)
	postResp, err := http.Post(TestUrl2, "application/json", bytes.NewBuffer(jsonValue))
	assert.NoError(t, err)
	body, err := ioutil.ReadAll(postResp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body[:]), `"amount": `+strconv.Itoa(TestBalance))
}
