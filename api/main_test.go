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
	"context"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/configuration"
)

const HOST = "http://localhost:19191"
const TestUrl = HOST + "/api/call"

type MainAPISuite struct {
	suite.Suite
}

func (suite *MainAPISuite) TestGetRequest() {
	resp, err := http.Get(TestUrl)
	suite.NoError(err)
	body, err := ioutil.ReadAll(resp.Body)
	suite.NoError(err)
	suite.Contains(string(body[:]), `"[ UnmarshalRequest ] Empty body"`)
}

func (suite *MainAPISuite) TestSerialization() {
	var a uint = 1
	var b bool = true
	var c string = "test"

	serArgs, err := core.MarshalArgs(a, b, c)
	suite.NoError(err)
	suite.NotNil(serArgs)

	var aR uint
	var bR bool
	var cR string
	rowResp, err := core.UnMarshalResponse(serArgs, []interface{}{aR, bR, cR})
	suite.NoError(err)
	suite.Len(rowResp, 3)
	suite.Equal(reflect.TypeOf(a), reflect.TypeOf(rowResp[0]))
	suite.Equal(reflect.TypeOf(b), reflect.TypeOf(rowResp[1]))
	suite.Equal(reflect.TypeOf(c), reflect.TypeOf(rowResp[2]))
	suite.Equal(a, rowResp[0].(uint))
	suite.Equal(b, rowResp[1].(bool))
	suite.Equal(c, rowResp[2].(string))
}

func (suite *MainAPISuite) TestNewApiRunnerNilConfig() {
	_, err := NewRunner(nil)
	suite.Contains(err.Error(), "config is nil")
}

func (suite *MainAPISuite) TestNewApiRunnerNoRequiredParams() {
	cfg := configuration.APIRunner{}
	_, err := NewRunner(&cfg)
	suite.Contains(err.Error(), "Address must not be empty")

	cfg.Address = "address:100"
	_, err = NewRunner(&cfg)
	suite.Contains(err.Error(), "Call must exist")

	cfg.Call = "test"
	_, err = NewRunner(&cfg)
	suite.Contains(err.Error(), "RPC must exist")

	cfg.RPC = "test"
	_, err = NewRunner(&cfg)
	suite.Contains(err.Error(), "Timeout must not be null")

	cfg.Timeout = 2
	_, err = NewRunner(&cfg)
	suite.NoError(err)
}

func TestMainTestSuite(t *testing.T) {
	ctx, _ := inslogger.WithTraceField(context.Background(), "APItests")
	cfg := configuration.NewAPIRunner()
	api, _ := NewRunner(&cfg)

	cm := certificate.NewCertificateManager(&certificate.Certificate{})
	api.CertificateManager = cm
	api.Start(ctx)

	suite.Run(t, new(MainAPISuite))

	api.Stop(ctx)
}
