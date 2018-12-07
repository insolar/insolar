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

package main

import (
	"context"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/stretchr/testify/require"
)

func TestInitComponents(t *testing.T) {
	ctx := context.Background()
	cfg := configuration.NewConfiguration()
	cfg.KeysPath = "testdata/bootstrap_keys.json"
	cfg.CertificatePath = "testdata/certificate.json"

	bootstrapComponents := initBootstrapComponents(ctx, cfg)
	cert := initCertificateManager(
		ctx,
		cfg,
		false,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.KeyProcessor,
	)
	cm, err := initComponents(
		ctx,
		cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.PlatformCryptographyScheme,
		bootstrapComponents.KeyStore,
		bootstrapComponents.KeyProcessor,
		cert,
		false,
		"",
		"",
	)
	require.NoError(t, err)
	require.NotNil(t, cm)

	err = cm.Init(ctx)
	require.NoError(t, err)

	// TODO: return it, when go to one role one node
	err = cm.Start(ctx)
	require.NoError(t, err)

	err = cm.Stop(ctx)
	require.NoError(t, err)
}
