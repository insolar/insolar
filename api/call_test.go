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
	"encoding/json"
	"testing"
	"time"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type TimeoutSuite struct {
	suite.Suite
	ctx   context.Context
	api   *Runner
	user  *requester.UserConfigJSON
	delay bool
}

type APIresp struct {
	Result *string
	Error  *string
}

func (suite *TimeoutSuite) TestRunner_callHandler() {
	seed, err := suite.api.SeedGenerator.Next()
	suite.NoError(err)
	suite.api.SeedManager.Add(*seed)

	resp, err := requester.SendWithSeed(
		suite.ctx,
		TestUrl,
		suite.user,
		&requester.RequestConfigJSON{},
		seed[:],
	)
	suite.NoError(err)

	var result APIresp
	err = json.Unmarshal(resp, &result)
	suite.NoError(err)
	suite.Nil(result.Error)
	suite.Equal("OK", *result.Result)
}

func (suite *TimeoutSuite) TestRunner_callHandlerTimeout() {
	seed, err := suite.api.SeedGenerator.Next()
	suite.NoError(err)
	suite.api.SeedManager.Add(*seed)

	suite.delay = true
	resp, err := requester.SendWithSeed(
		suite.ctx,
		TestUrl,
		suite.user,
		&requester.RequestConfigJSON{},
		seed[:],
	)
	suite.NoError(err)

	var result APIresp
	err = json.Unmarshal(resp, &result)
	suite.NoError(err)
	suite.Equal("Messagebus timeout exceeded", *result.Error)
	suite.Nil(result.Result)
}

func TestTimeoutSuite(t *testing.T) {
	timeoutSuite := new(TimeoutSuite)
	timeoutSuite.ctx, _ = inslogger.WithTraceField(context.Background(), "APItests")

	ks := platformpolicy.NewKeyProcessor()
	sKey, err := ks.GeneratePrivateKey()
	require.NoError(t, err)
	sKeyString, err := ks.ExportPrivateKeyPEM(sKey)
	require.NoError(t, err)
	pKey := ks.ExtractPublicKey(sKey)
	pKeyString, err := ks.ExportPublicKeyPEM(pKey)
	require.NoError(t, err)

	userRef := testutils.RandomRef().String()
	timeoutSuite.user, err = requester.CreateUserConfig(userRef, string(sKeyString))

	cfg := configuration.NewAPIRunner()
	timeoutSuite.api, err = NewRunner(&cfg)
	require.NoError(t, err)

	cert := testutils.NewCertificateMock(t)
	cert.GetRootDomainReferenceFunc = func() (r *core.RecordRef) {
		ref := testutils.RandomRef()
		return &ref
	}

	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() (r core.Certificate) {
		return cert
	}

	cr := testutils.NewContractRequesterMock(t)
	cr.SendRequestFunc = func(p context.Context, p1 *core.RecordRef, method string, p3 []interface{}) (core.Reply, error) {
		switch method {
		case "GetPublicKey":
			var result = string(pKeyString)
			var contractErr *foundation.Error
			data, _ := core.MarshalArgs(result, contractErr)
			return &reply.CallMethod{
				Result: data,
			}, nil
		default:
			if timeoutSuite.delay {
				time.Sleep(time.Second * 10)
			}
			var result = "OK"
			var contractErr *foundation.Error
			data, _ := core.MarshalArgs(result, contractErr)
			return &reply.CallMethod{
				Result: data,
			}, nil
		}
	}

	timeoutSuite.api.ContractRequester = cr
	timeoutSuite.api.CertificateManager = cm
	timeoutSuite.api.Start(timeoutSuite.ctx)

	suite.Run(t, timeoutSuite)

	timeoutSuite.api.Stop(timeoutSuite.ctx)
}
