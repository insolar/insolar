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

package observer

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/server/internal/observer/stubs"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/platformpolicy"
)

type components struct {
	cmp      component.Manager
	NodeRef  string
	NodeRole string
}

func newComponents(ctx context.Context, cfg configuration.Configuration) (*components, error) {
	// Cryptography.
	var (
		KeyProcessor  insolar.KeyProcessor
		CryptoScheme  insolar.PlatformCryptographyScheme
		CryptoService insolar.CryptographyService
		CertManager   insolar.CertificateManager
	)
	{
		var err error
		// Private key storage.
		ks, err := keystore.NewKeyStore(cfg.KeysPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to load KeyStore")
		}
		// Public key manipulations.
		KeyProcessor = platformpolicy.NewKeyProcessor()
		// Platform cryptography.
		CryptoScheme = platformpolicy.NewPlatformCryptographyScheme()
		// Sign, verify, etc.
		CryptoService = cryptography.NewCryptographyService()

		c := component.Manager{}
		c.Inject(CryptoService, CryptoScheme, KeyProcessor, ks)

		publicKey, err := CryptoService.GetPublicKey()
		if err != nil {
			return nil, errors.Wrap(err, "failed to retrieve node public key")
		}

		// Node certificate.
		CertManager, err = certificate.NewManagerReadCertificate(publicKey, KeyProcessor, cfg.CertificatePath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Certificate")
		}
	}

	c := &components{}
	c.cmp = component.Manager{}
	c.NodeRef = CertManager.GetCertificate().GetNodeRef().String()
	c.NodeRole = CertManager.GetCertificate().GetRole().String()

	// Network.
	var (
		HostNetwork      network.HostNetwork
		RoutingTableStub network.RoutingTable
		TransportFactory transport.Factory
	)
	{
		var err error
		// External communication.
		certPath := cfg.CertificatePath
		data, err := ioutil.ReadFile(certPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read certificate from: %s", certPath)
		}
		cert := certificate.AuthorizationCertificate{}
		err = json.Unmarshal(data, &cert)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal cert data")
		}
		HostNetwork, err = hostnetwork.NewHostNetwork(cert.GetNodeRef().String())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create new HostNetwork")
		}

		RoutingTableStub = &stubs.RoutingTable{}

		TransportFactory = transport.NewFactory(cfg.Host.Transport)
	}

	// API.
	var (
		// Requester insolar.ContractRequester
		API insolar.APIRunner
	)
	{
		var err error
		API, err = stubs.NewRunner(&cfg.APIRunner, &cfg.Host.Transport, CertManager.GetCertificate())
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ApiRunner")
		}
	}

	// Storage.
	var (
		DB *store.BadgerDB
	)
	{
		var err error
		DB, err = store.NewBadgerDB(cfg.Ledger.Storage.DataDirectory)
		if err != nil {
			panic(errors.Wrap(err, "failed to initialize DB"))
		}
	}

	metricsHandler, err := metrics.NewMetrics(
		ctx,
		cfg.Metrics,
		metrics.GetInsolarRegistry(c.NodeRole),
		c.NodeRole,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to start Metrics")
	}

	c.cmp.Inject(
		DB,
		metricsHandler,
		API,
		KeyProcessor,
		CryptoScheme,
		CryptoService,
		CertManager,
		HostNetwork,
		RoutingTableStub,
		TransportFactory,
	)
	err = c.cmp.Init(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init components")
	}

	return c, nil
}

func (c *components) Start(ctx context.Context) error {
	return c.cmp.Start(ctx)
}

func (c *components) Stop(ctx context.Context) error {
	return c.cmp.Stop(ctx)
}
