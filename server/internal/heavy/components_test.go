// Copyright 2020 Insolar Network Ltd.
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

// +build slowtest

package heavy

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestComponents(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "heavy-")
	defer os.RemoveAll(tmpdir)
	require.NoError(t, err)

	ctx := inslogger.UpdateLogger(context.Background(), func(logger insolar.Logger) (insolar.Logger, error) {
		return logger.Copy().WithBuffer(100, false).Build()
	})
	cfg := configuration.NewConfiguration()
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.CertificatePath = "testdata/certificate.json"
	cfg.Metrics.ListenAddress = "0.0.0.0:0"
	cfg.APIRunner.Address = "0.0.0.0:0"
	cfg.AdminAPIRunner.Address = "0.0.0.0:0"
	cfg.APIRunner.SwaggerPath = "../../../application/api/spec/api-exported.yaml"
	cfg.AdminAPIRunner.SwaggerPath = "../../../application/api/spec/api-exported.yaml"
	cfg.Ledger.Storage.DataDirectory = tmpdir
	cfg.Exporter.Addr = ":0"

	_, err = newComponents(ctx, cfg, application.GenesisHeavyConfig{Skip: true}, false)
	require.NoError(t, err)
}
