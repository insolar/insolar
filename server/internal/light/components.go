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

package light

import (
	"context"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/genesisdataprovider"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/artifactmanager"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/pulsemanager"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/light/replication"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/network/termination"
	"github.com/insolar/insolar/networkcoordinator"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

type components struct {
	cmp               component.Manager
	NodeRef, NodeRole string
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
		NetworkService     insolar.Network
		NetworkCoordinator insolar.NetworkCoordinator
		NodeNetwork        insolar.NodeNetwork
		Termination        insolar.TerminationHandler
	)
	{
		var err error
		// External communication.
		NetworkService, err = servicenetwork.NewServiceNetwork(cfg, &c.cmp, false)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Network")
		}

		Termination = termination.NewHandler(NetworkService)

		// Node info.
		NodeNetwork, err = nodenetwork.NewNodeNetwork(cfg.Host.Transport, CertManager.GetCertificate())
		if err != nil {
			return nil, errors.Wrap(err, "failed to start NodeNetwork")
		}

		NetworkCoordinator, err = networkcoordinator.New()
		if err != nil {
			return nil, errors.Wrap(err, "failed to start NetworkCoordinator")
		}
	}

	// API.
	var (
		Requester insolar.ContractRequester
		Genesis   insolar.GenesisDataProvider
		API       insolar.APIRunner
	)
	{
		var err error
		Requester, err = contractrequester.New()
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ContractRequester")
		}

		Genesis, err = genesisdataprovider.New()
		if err != nil {
			return nil, errors.Wrap(err, "failed to start GenesisDataProvider")
		}

		API, err = api.NewRunner(&cfg.APIRunner)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ApiRunner")
		}
	}

	// Communication.
	var (
		Tokens  insolar.DelegationTokenFactory
		Parcels message.ParcelFactory
		Bus     insolar.MessageBus
	)
	{
		var err error
		Tokens = delegationtoken.NewDelegationTokenFactory()
		Parcels = messagebus.NewParcelFactory()
		Bus, err = messagebus.NewMessageBus(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start MessageBus")
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

	// Light components.
	var (
		PulseManager insolar.PulseManager
		Coordinator  jet.Coordinator
		Pulses       pulse.Accessor
		Jets         jet.Accessor
		Handler      *artifactmanager.MessageHandler
	)
	{
		conf := cfg.Ledger
		idLocker := object.NewIDLocker()
		pulses := pulse.NewStorageMem()
		drops := drop.NewStorageMemory()
		blobs := blob.NewStorageMemory()
		records := object.NewRecordMemory()
		indexes := object.NewIndexMemory()
		jets := jet.NewStore()
		nodes := node.NewStorage()

		c := component.Manager{}
		c.Inject(CryptoScheme)

		hots := recentstorage.NewRecentStorageProvider()
		waiter := hot.NewChannelWaiter()
		cord := jetcoordinator.NewJetCoordinator(conf.LightChainLimit)
		cord.PulseCalculator = pulses
		cord.PulseAccessor = pulses
		cord.JetAccessor = jets
		cord.NodeNet = NodeNetwork
		cord.PlatformCryptographyScheme = CryptoScheme
		cord.Nodes = nodes
		Coordinator = cord

		handler := artifactmanager.NewMessageHandler(indexes, indexes, &conf)
		handler.RecentStorageProvider = hots
		handler.Bus = Bus
		handler.PlatformCryptographyScheme = CryptoScheme
		handler.JetCoordinator = Coordinator
		handler.CryptographyService = CryptoService
		handler.DelegationTokenFactory = Tokens
		handler.JetStorage = jets
		handler.DropModifier = drops
		handler.BlobModifier = blobs
		handler.BlobAccessor = blobs
		handler.Blobs = blobs
		handler.IDLocker = idLocker
		handler.RecordModifier = records
		handler.RecordAccessor = records
		handler.Nodes = nodes
		handler.HotDataWaiter = waiter
		handler.JetReleaser = waiter
		handler.IndexStorage = indexes
		handler.IndexStateModifier = indexes
		handler.IndexStorage = indexes

		jetCalculator := jet.NewCalculator(Coordinator, jets)
		var lightCleaner = replication.NewCleaner(
			jets,
			nodes,
			drops,
			blobs,
			records,
			indexes,
			pulses,
			pulses,
			conf.LightChainLimit,
		)
		dataGatherer := replication.NewDataGatherer(drops, blobs, records, indexes)
		lthSyncer := replication.NewReplicatorDefault(
			jetCalculator,
			dataGatherer,
			lightCleaner,
			Bus,
			pulses,
		)

		pm := pulsemanager.NewPulseManager(
			conf,
			drops,
			blobs,
			blobs,
			pulses,
			records,
			records,
			indexes,
			indexes,
			lthSyncer,
		)
		pm.MessageHandler = handler
		pm.Bus = Bus
		pm.NodeNet = NodeNetwork
		pm.JetCoordinator = Coordinator
		pm.CryptographyService = CryptoService
		pm.PlatformCryptographyScheme = CryptoScheme
		pm.RecentStorageProvider = hots
		pm.JetReleaser = waiter
		pm.JetAccessor = jets
		pm.JetModifier = jets
		pm.IndexAccessor = indexes
		pm.IndexModifier = indexes
		pm.CollectionIndexAccessor = indexes
		pm.IndexCleaner = indexes
		pm.NodeSetter = nodes
		pm.Nodes = nodes
		pm.DropModifier = drops
		pm.DropAccessor = drops
		pm.DropCleaner = drops
		pm.PulseAccessor = pulses
		pm.PulseCalculator = pulses
		pm.PulseAppender = pulses

		PulseManager = pm
		Pulses = pulses
		Jets = jets
		Handler = handler
	}

	c.cmp.Inject(
		Handler,
		Jets,
		Pulses,
		Coordinator,
		PulseManager,
		metricsHandler,
		Bus,
		Requester,
		Tokens,
		Parcels,
		artifacts.NewClient(),
		Genesis,
		API,
		NetworkCoordinator,
		KeyProcessor,
		Termination,
		CryptoScheme,
		CryptoService,
		CertManager,
		NodeNetwork,
		NetworkService,
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
