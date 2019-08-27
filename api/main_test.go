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
	"net/http"
	"testing"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/suite"
)

type MainAPISuite struct {
	suite.Suite
}

func (suite *MainAPISuite) TestNewApiRunnerNilConfig() {
	_, err := NewRunner(nil, nil, nil, nil, nil, nil, nil, nil, nil)
	suite.Contains(err.Error(), "config is nil")
}

func (suite *MainAPISuite) TestNewApiRunnerNoRequiredParams() {
	cfg := configuration.APIRunner{}
	_, err := NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil)
	suite.Contains(err.Error(), "Address must not be empty")

	cfg.Address = "address:100"
	_, err = NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil)
	suite.Contains(err.Error(), "RPC must exist")

	cfg.RPC = "test"
	_, err = NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil)
	suite.NoError(err)
}

func TestMainTestSuite(t *testing.T) {
	ctx, _ := inslogger.WithTraceField(context.Background(), "APItests")
	http.DefaultServeMux = new(http.ServeMux)
	cfg := configuration.NewAPIRunner(false)
	api, _ := NewRunner(&cfg, nil, nil, nil, nil, nil, nil, nil, nil)

	cm := certificate.NewCertificateManager(&certificate.Certificate{})
	api.CertificateManager = cm
	api.Start(ctx)

	suite.Run(t, new(MainAPISuite))

	api.Stop(ctx)
}
