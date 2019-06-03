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

	"github.com/ThreeDotsLabs/watermill"
	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/genesisdataprovider"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/delegationtoken"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/node"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/light/artifactmanager"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/light/pulsemanager"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/light/replication"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/messagebus"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/network/termination"
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

	logger := log.NewWatermillLogAdapter(inslogger.FromContext(ctx))
	pubSub := gochannel.NewGoChannel(gochannel.Config{}, logger)

	// Network.
	var (
		NetworkService *servicenetwork.ServiceNetwork
		NodeNetwork    insolar.NodeNetwork
		Termination    insolar.TerminationHandler
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

	// Role calculations.
	var (
		Coordinator jet.Coordinator
		Pulses      *pulse.StorageMem
		Jets        jet.Storage
		Nodes       *node.Storage
	)
	{
		Nodes = node.NewStorage()
		Pulses = pulse.NewStorageMem()
		Jets = jet.NewStore()

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit)
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.NodeNet = NodeNetwork
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes

		Coordinator = c
	}

	// Communication.
	var (
		Tokens  insolar.DelegationTokenFactory
		Parcels message.ParcelFactory
		Bus     insolar.MessageBus
		WmBus   *bus.Bus
	)
	{
		var err error
		Tokens = delegationtoken.NewDelegationTokenFactory()
		Parcels = messagebus.NewParcelFactory()
		Bus, err = messagebus.NewMessageBus(cfg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start MessageBus")
		}
		WmBus = bus.NewBus(pubSub, Pulses, Coordinator)
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
		Handler      *artifactmanager.MessageHandler
	)
	{
		conf := cfg.Ledger
		idLocker := object.NewIDLocker()
		drops := drop.NewStorageMemory()
		blobs := blob.NewStorageMemory()
		records := object.NewRecordMemory()
		indexes := object.NewInMemoryIndex()
		writeController := hot.NewWriteController()

		c := component.Manager{}
		c.Inject(CryptoScheme)

		hots := recentstorage.NewProvider()
		waiter := hot.NewChannelWaiter()

		handler := artifactmanager.NewMessageHandler(indexes, indexes, indexes, &conf)
		handler.RecentStorageProvider = hots
		handler.Bus = Bus
		handler.PCS = CryptoScheme
		handler.JetCoordinator = Coordinator
		handler.CryptographyService = CryptoService
		handler.DelegationTokenFactory = Tokens
		handler.JetStorage = Jets
		handler.DropModifier = drops
		handler.BlobModifier = blobs
		handler.BlobAccessor = blobs
		handler.Blobs = blobs
		handler.IDLocker = idLocker
		handler.RecordModifier = records
		handler.RecordAccessor = records
		handler.Nodes = Nodes
		handler.HotDataWaiter = waiter
		handler.JetReleaser = waiter
		handler.WriteAccessor = writeController
		handler.Sender = WmBus

		jetCalculator := jet.NewCalculator(Coordinator, Jets)
		var lightCleaner = replication.NewCleaner(
			Jets,
			Nodes,
			drops,
			blobs,
			records,
			indexes,
			Pulses,
			Pulses,
			conf.LightChainLimit,
		)
		dataGatherer := replication.NewDataGatherer(drops, blobs, records, indexes)
		lthSyncer := replication.NewReplicatorDefault(
			jetCalculator,
			dataGatherer,
			lightCleaner,
			Bus,
			Pulses,
		)

		pm := pulsemanager.NewPulseManager(
			conf,
			drops,
			blobs,
			blobs,
			Pulses,
			records,
			records,
			indexes,
			lthSyncer,
			writeController,
		)
		pm.MessageHandler = handler
		pm.Bus = Bus
		pm.NodeNet = NodeNetwork
		pm.JetCoordinator = Coordinator
		pm.CryptographyService = CryptoService
		pm.PlatformCryptographyScheme = CryptoScheme
		pm.RecentStorageProvider = hots
		pm.JetReleaser = waiter
		pm.JetAccessor = Jets
		pm.JetModifier = Jets
		pm.NodeSetter = Nodes
		pm.Nodes = Nodes
		pm.DropModifier = drops
		pm.DropAccessor = drops
		pm.DropCleaner = drops
		pm.PulseAccessor = Pulses
		pm.PulseCalculator = Pulses
		pm.PulseAppender = Pulses

		PulseManager = pm
		Handler = handler
	}

	c.cmp.Inject(
		WmBus,
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
		artifacts.NewClient(WmBus),
		Genesis,
		API,
		KeyProcessor,
		Termination,
		CryptoScheme,
		CryptoService,
		CertManager,
		NodeNetwork,
		NetworkService,
		pubSub,
	)

	err = c.cmp.Init(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init components")
	}

	startWatermill(ctx, logger, pubSub, WmBus, NetworkService.SendMessageHandler, Handler.FlowDispatcher.Process)

	return c, nil
}

func (c *components) Start(ctx context.Context) error {
	return c.cmp.Start(ctx)
}

func (c *components) Stop(ctx context.Context) error {
	return c.cmp.Stop(ctx)
}

func startWatermill(
	ctx context.Context,
	logger watermill.LoggerAdapter,
	pubSub watermillMsg.PubSub,
	b *bus.Bus,
	outHandler, inHandler watermillMsg.HandlerFunc,
) {
	inRouter, err := watermillMsg.NewRouter(watermillMsg.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}
	outRouter, err := watermillMsg.NewRouter(watermillMsg.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	outRouter.AddNoPublisherHandler(
		"OutgoingHandler",
		bus.TopicOutgoing,
		pubSub,
		outHandler,
	)

	inRouter.AddMiddleware(
		middleware.InstantAck,
		b.IncomingMessageRouter,
		middleware.CorrelationID,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingHandler",
		bus.TopicIncoming,
		pubSub,
		inHandler,
	)

	startRouter(ctx, inRouter)
	startRouter(ctx, outRouter)
}

func startRouter(ctx context.Context, router *watermillMsg.Router) {
	go func() {
		if err := router.Run(); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}
