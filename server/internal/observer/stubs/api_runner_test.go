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

package stubs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestNewRunner(t *testing.T) {
	api, err := NewRunner(nil, nil, nil)
	require.Nil(t, api)
	require.Error(t, err)

	cfg := configuration.NewAPIRunner()
	cfg.Address = "localhost:19192"
	cfgTransport := configuration.Transport{}
	cert := certificate.Certificate{}
	api, err = NewRunner(&cfg, &cfgTransport, &cert)
	require.NoError(t, err)
	require.NotNil(t, api)
}

func TestApiRunnerStub_IsAPIRunner(t *testing.T) {
	cfg := configuration.NewAPIRunner()
	cfg.Address = "localhost:19192"
	cfgTransport := configuration.Transport{}
	cert := certificate.Certificate{}
	api, err := NewRunner(&cfg, &cfgTransport, &cert)

	require.NoError(t, err)
	require.True(t, api.IsAPIRunner())
}

func TestApiRunnerStub_StartStop(t *testing.T) {
	ctx := inslogger.TestContext(t)
	cfg := configuration.NewAPIRunner()
	cfg.Address = "localhost:19192"
	cfgTransport := configuration.Transport{}
	cert := certificate.Certificate{}
	api, err := NewRunner(&cfg, &cfgTransport, &cert)
	require.NoError(t, err)

	starter, ok := api.(component.Starter)
	require.True(t, ok)

	stopper, ok := api.(component.Stopper)
	require.True(t, ok)

	err = starter.Start(ctx)
	require.NoError(t, err)

	err = stopper.Stop(ctx)
	require.NoError(t, err)
}
