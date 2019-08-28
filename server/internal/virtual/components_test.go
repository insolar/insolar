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
// +build slowtest

package virtual

import (
	"context"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/require"
)

func TestComponents(t *testing.T) {
	ctx := context.Background()
	cfg := configuration.NewConfiguration()
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.CertificatePath = "testdata/certificate.json"
	cfg.Metrics.ListenAddress = "0.0.0.0:0"
	cfg.APIRunner.Address = "0.0.0.0:0"

	bootstrapComponents := initBootstrapComponents(ctx, cfg)
	cert := initCertificateManager(
		ctx,
		cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.KeyProcessor,
	)
	cm, _, stopWatermill := initComponents(
		ctx,
		cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.PlatformCryptographyScheme,
		bootstrapComponents.KeyStore,
		bootstrapComponents.KeyProcessor,
		cert,
	)
	require.NotNil(t, cm)
	require.NotNil(t, stopWatermill)

	err := cm.Init(ctx)
	require.NoError(t, err)
}
