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

package api

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

const CallUrl = "http://localhost:19192/api/rpc"

type TimeoutSuite struct {
	suite.Suite

	mc    *minimock.Controller
	ctx   context.Context
	api   *Runner
	user  *requester.UserConfigJSON
	delay chan struct{}
}

func (suite *TimeoutSuite) TestRunner_callHandler_NoTimeout() {
	seed, err := suite.api.SeedGenerator.Next()
	suite.NoError(err)
	suite.api.SeedManager.Add(*seed, 0)

	close(suite.delay)
	suite.api.timeout = 60 * time.Second
	seedString := base64.StdEncoding.EncodeToString(seed[:])

	resp, err := requester.SendWithSeed(
		suite.ctx,
		CallUrl,
		suite.user,
		&requester.Params{CallSite: "member.create", CallParams: map[string]interface{}{}, PublicKey: suite.user.PublicKey},
		seedString,
	)
	suite.NoError(err)

	var result requester.ContractResponse
	err = json.Unmarshal(resp, &result)
	suite.NoError(err)
	suite.Nil(result.Error)
	suite.Equal("OK", result.Result.CallResult)
}

func (suite *TimeoutSuite) TestRunner_callHandler_Timeout() {
	seed, err := suite.api.SeedGenerator.Next()
	suite.NoError(err)
	suite.api.SeedManager.Add(*seed, 0)

	suite.api.timeout = 1 * time.Millisecond

	seedString := base64.StdEncoding.EncodeToString(seed[:])

	_, err = requester.SendWithSeed(
		suite.ctx,
		CallUrl,
		suite.user,
		nil,
		seedString,
	)
	suite.Error(err, "Client.Timeout exceeded while awaiting headers")
}

func TestTimeoutSuite(t *testing.T) {
	timeoutSuite := new(TimeoutSuite)
	timeoutSuite.ctx, _ = inslogger.WithTraceField(context.Background(), "APItests")
	timeoutSuite.mc = minimock.NewController(t)

	ks := platformpolicy.NewKeyProcessor()
	sKey, err := ks.GeneratePrivateKey()
	require.NoError(t, err)
	sKeyString, err := ks.ExportPrivateKeyPEM(sKey)
	require.NoError(t, err)
	pKey := ks.ExtractPublicKey(sKey)
	pKeyString, err := ks.ExportPublicKeyPEM(pKey)
	require.NoError(t, err)

	userRef := gen.Reference().String()
	timeoutSuite.user, err = requester.CreateUserConfig(userRef, string(sKeyString), string(pKeyString))

	http.DefaultServeMux = new(http.ServeMux)
	cfg := configuration.NewAPIRunner(false)
	cfg.Address = "localhost:19192"
	timeoutSuite.api, err = NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	timeoutSuite.api.timeout = 1 * time.Second

	cr := testutils.NewContractRequesterMock(timeoutSuite.mc)
	cr.SendRequestWithPulseMock.Set(func(p context.Context, p1 *insolar.Reference, method string, p3 []interface{}, p4 insolar.PulseNumber) (insolar.Reply, *insolar.Reference, error) {
		requestReference, _ := insolar.NewReferenceFromBase58("4K3NiGuqYGqKPnYp6XeGd2kdN4P9veL6rYcWkLKWXZCu.4FFB8zfQoGznSmzDxwv4njX1aR9ioL8GHSH17QXH2AFa")
		switch method {
		case "GetPublicKey":
			var result = string(pKeyString)
			data, _ := foundation.MarshalMethodResult(result, nil)
			return &reply.CallMethod{
				Result: data,
			}, requestReference, nil
		default:
			<-timeoutSuite.delay
			var result = "OK"
			data, _ := foundation.MarshalMethodResult(result, nil)
			return &reply.CallMethod{
				Result: data,
			}, requestReference, nil
		}
	})

	timeoutSuite.api.ContractRequester = cr
	timeoutSuite.api.Start(timeoutSuite.ctx)
	timeoutSuite.api.SeedManager = seedmanager.NewSpecified(time.Minute, time.Minute)

	requester.SetTimeout(25)
	suite.Run(t, timeoutSuite)

	timeoutSuite.api.Stop(timeoutSuite.ctx)
}

func TestDigestParser(t *testing.T) {
	invalidDigest := ""
	_, err := parseDigest(invalidDigest)
	require.Error(t, err)

	validDigest := "SHA-256=foo"
	_, err = parseDigest(validDigest)
	require.NoError(t, err)
}

func TestSignatureParser(t *testing.T) {
	invalidSignature := ""
	_, err := parseSignature(invalidSignature)

	validSignature := `keyId="member-pub-key", algorithm="ecdsa", headers="digest", signature=bar`
	_, err = parseSignature(validSignature)
	require.NoError(t, err)
}

func TestValidateRequestHeaders(t *testing.T) {
	body := []byte("foobar")
	h := sha256.New()
	_, err := h.Write(body)
	require.NoError(t, err)

	digest := h.Sum(nil)
	calculatedDigest := `SHA-256=` + base64.URLEncoding.EncodeToString(digest)
	signature := `keyId="member-pub-key", algorithm="ecdsa", headers="digest", signature=bar`
	sig, err := validateRequestHeaders(calculatedDigest, signature, body)
	require.NoError(t, err)
	require.Equal(t, "bar", sig)
}

func (suite *TimeoutSuite) BeforeTest(suiteName, testName string) {
	suite.delay = make(chan struct{}, 0)
}

func (suite *TimeoutSuite) AfterTest(suiteName, testName string) {
	suite.mc.Wait(1 * time.Minute)
	suite.mc.Finish()
}
