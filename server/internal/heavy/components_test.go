// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

// +build slowtest

package heavy

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/applicationbase/genesis"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

func TestComponents(t *testing.T) {
	allowedVersionContract = 2
	tests := []struct {
		name         string
		exporterAuth bool
		authSecret   string
		errorCheck   func(*testing.T, error)
	}{{"no exporter auth",
		false,
		"",
		func(t *testing.T, err error) { require.NoError(t, err) },
	}, {"exporter auth - 512-bit secret",
		true,
		"/B?E(H+MbQeThWmZq4t7w!z$C&F)J@NcRfUjXn2r5u8x/A?D*G-KaPdSgVkYp3s6",
		func(t *testing.T, err error) { require.NoError(t, err) },
	}, {"exporter auth - short secret",
		true,
		"E(G+KbPeShVmYq3t6w9z$C&F)J@McQfT",
		func(t *testing.T, err error) { require.Error(t, err) },
	}, {"exporter auth - long secret",
		true,
		"B?E(H+MbQeThVmYq3t6w9z$C&F)J@NcRfUjXnZr4u7x!A%D*G-KaPdSgVkYp3s5v8y/B?E(H+MbQeThWmZq4t7w9z$C&F)J@NcRfUjXn2r5u8x/A%D*G-KaPdSgVkYp3",
		func(t *testing.T, err error) { require.Error(t, err) },
	}, {"exporter auth - empty secret",
		true,
		"",
		func(t *testing.T, err error) { require.Error(t, err) },
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpdir, err := ioutil.TempDir("", "heavy-")
			defer os.RemoveAll(tmpdir)
			require.NoError(t, err)

			ctx := inslogger.UpdateLogger(context.Background(), func(logger insolar.Logger) (insolar.Logger, error) {
				return logger.Copy().WithBuffer(100, false).Build()
			})
			cfg := configuration.NewHeavyBadgerConfig()
			cfg.KeysPath = "testdata/bootstrap_keys.json"
			cfg.CertificatePath = "testdata/certificate.json"
			cfg.Metrics.ListenAddress = "0.0.0.0:0"
			cfg.APIRunner.Address = "0.0.0.0:0"
			cfg.AdminAPIRunner.Address = "0.0.0.0:0"
			cfg.APIRunner.SwaggerPath = "../../../api/testdata/api-exported.yaml"
			cfg.AdminAPIRunner.SwaggerPath = "../../../api/testdata/api-exported.yaml"
			cfg.Ledger.Storage.DataDirectory = tmpdir
			cfg.Exporter.Addr = ":0"
			cfg.Exporter.Auth.Required = test.exporterAuth
			cfg.Exporter.Auth.Secret = test.authSecret

			holder := &configuration.HeavyBadgerHolder{
				Configuration: &cfg,
			}
			_, err = newComponents(ctx, holder, genesis.HeavyConfig{Skip: true}, genesis.Options{}, false, api.Options{}, allowedVersionContract)
			test.errorCheck(t, err)
		})
	}
}
