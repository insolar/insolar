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
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

const CallUrl = "http://localhost:19192/api/call"

type TimeoutSuite struct {
	suite.Suite
	ctx   context.Context
	mc    *minimock.Controller
	cr    *testutils.ContractRequesterMock
	api   *Runner
	user  *requester.UserConfigJSON
	userPublicKey string
	delay bool
}

type APIresp struct {
	Result string
	Error  string
}


func (s *TimeoutSuite) TestRunner_callHandler() {
	seed, err := s.api.SeedGenerator.Next()
	s.NoError(err)
	s.api.SeedManager.Add(*seed)

	s.cr.SendRequestFunc = func(p context.Context, p1 *insolar.Reference, method string, p3 []interface{}) (insolar.Reply, error) {
		switch method {
		case "GetPublicKey":
			var contractErr *foundation.Error
			data, _ := insolar.MarshalArgs(s.userPublicKey, contractErr)
			return &reply.CallMethod{
				Result: data,
			}, nil
		default:
			var contractErr *foundation.Error
			data, _ := insolar.MarshalArgs("OK", contractErr)
			return &reply.CallMethod{
				Result: data,
			}, nil
		}
	}

	resp, err := requester.SendWithSeed(
		s.ctx,
		CallUrl,
		s.user,
		&requester.RequestConfigJSON{},
		seed[:],
	)
	s.NoError(err)

	var result APIresp
	err = json.Unmarshal(resp, &result)
	s.NoError(err)
	s.Equal("", result.Error)
	s.Equal("OK", result.Result)
}

func (s *TimeoutSuite) TestCallHandlerTimeout() {
	seed, err := s.api.SeedGenerator.Next()
	s.NoError(err)
	s.api.SeedManager.Add(*seed)

	ch := make(chan struct{}, 1)
	s.cr.SendRequestFunc = func(p context.Context, p1 *insolar.Reference, method string, p3 []interface{}) (insolar.Reply, error) {
		switch method {
		case "GetPublicKey":
			var contractErr *foundation.Error
			data, _ := insolar.MarshalArgs(s.userPublicKey, contractErr)
			return &reply.CallMethod{
				Result: data,
			}, nil
		default:
			<-ch
			var contractErr *foundation.Error
			data, _ := insolar.MarshalArgs("OK", contractErr)
			return &reply.CallMethod{
				Result: data,
			}, nil
		}
	}

	resp, err := requester.SendWithSeed(
		s.ctx,
		CallUrl,
		s.user,
		&requester.RequestConfigJSON{},
		seed[:],
	)
	s.NoError(err)

	close(ch)

	var result APIresp
	err = json.Unmarshal(resp, &result)
	s.NoError(err)
	s.Equal("Messagebus timeout exceeded", result.Error)
	s.Equal("", result.Result)
}

func TestTimeoutSuite(t *testing.T) {
	s := new(TimeoutSuite)

	ks := platformpolicy.NewKeyProcessor()
	sKey, err := ks.GeneratePrivateKey()
	require.NoError(t, err)
	sKeyString, err := ks.ExportPrivateKeyPEM(sKey)
	require.NoError(t, err)
	pKey := ks.ExtractPublicKey(sKey)
	pKeyString, err := ks.ExportPublicKeyPEM(pKey)
	require.NoError(t, err)

	userRef := testutils.RandomRef().String()
	s.user, err = requester.CreateUserConfig(userRef, string(sKeyString))
	require.NoError(t, err)

	s.userPublicKey = string(pKeyString)

	suite.Run(t, s)
}

func (s *TimeoutSuite) BeforeTest(suiteName, testName string) {
	t := s.T()

	s.ctx = inslogger.TestContext(t)
	s.mc = minimock.NewController(t)

	cfg := configuration.NewAPIRunner()
	cfg.Address = "localhost:19192"
	cfg.Timeout = 1

	var err error
	s.api, err = NewRunner(&cfg)
	require.NoError(t, err)

	rootDomainRef := testutils.RandomRef()

	cert := testutils.NewCertificateMock(t)
	cert.GetRootDomainReferenceFunc = func() (r *insolar.Reference) {
		return &rootDomainRef
	}

	cm := testutils.NewCertificateManagerMock(t)
	cm.GetCertificateFunc = func() (r insolar.Certificate) {
		return cert
	}

	s.cr = testutils.NewContractRequesterMock(s.mc)
	s.api.ContractRequester = s.cr
	s.api.CertificateManager = cm

	http.DefaultServeMux = new(http.ServeMux)
	err = s.api.Start(s.ctx)
	require.NoError(t, err)

	requester.SetTimeout(3600)
}

func (s *TimeoutSuite) AfterTest(suiteName, testName string) {
	s.mc.Wait(1*time.Minute)
	s.mc.Finish()

	_ = s.api.Stop(s.ctx)
}
