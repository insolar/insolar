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

package heavy

import (
	"context"
	"fmt"
	"net"
	"path/filepath"

	"github.com/ThreeDotsLabs/watermill"
	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"github.com/insolar/component-manager"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/application/api"
	"github.com/insolar/insolar/application/genesis"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/contractrequester"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/jetcoordinator"
	"github.com/insolar/insolar/insolar/node"
	insolarPulse "github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/store"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/artifact"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/insolar/insolar/ledger/heavy/exporter"
	"github.com/insolar/insolar/ledger/heavy/handler"
	"github.com/insolar/insolar/ledger/heavy/pulsemanager"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/log/logwatermill"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/metrics"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/server/internal"
)

type components struct {
	cmp         *component.Manager
	NodeRef     string
	NodeRole    string
	rollback    *executor.DBRollback
	stateKeeper *executor.InitialStateKeeper
	inRouter    *watermillMsg.Router
	outRouter   *watermillMsg.Router

	replicator executor.HeavyReplicator
}

func initTemporaryCertificateManager(ctx context.Context, cfg *configuration.Configuration) (*certificate.CertificateManager, error) {
	earlyComponents := component.NewManager(nil)

	keyStore, err := keystore.NewKeyStore(cfg.KeysPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load KeyStore")
	}

	platformCryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	keyProcessor := platformpolicy.NewKeyProcessor()

	cryptographyService := cryptography.NewCryptographyService()
	earlyComponents.Register(platformCryptographyScheme, keyStore)
	earlyComponents.Inject(cryptographyService, keyProcessor)

	publicKey, err := cryptographyService.GetPublicKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve node public key")
	}

	certManager, err := certificate.NewManagerReadCertificate(publicKey, keyProcessor, cfg.CertificatePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new CertificateManager")
	}

	return certManager, nil
}

func newComponents(ctx context.Context, cfg configuration.Configuration, genesisCfg application.GenesisHeavyConfig) (*components, error) {
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

		c := component.NewManager(nil)
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
	c.cmp = component.NewManager(nil)
	c.NodeRef = CertManager.GetCertificate().GetNodeRef().String()
	c.NodeRole = CertManager.GetCertificate().GetRole().String()

	logger := inslogger.FromContext(ctx)

	// Watermill stuff.
	var (
		wmLogger   *logwatermill.WatermillLogAdapter
		publisher  watermillMsg.Publisher
		subscriber watermillMsg.Subscriber
	)
	{
		wmLogger = logwatermill.NewWatermillLogAdapter(logger)
		pubsub := gochannel.NewGoChannel(gochannel.Config{}, wmLogger)
		subscriber = pubsub
		publisher = pubsub
		// Wrapped watermill publisher for introspection.
		publisher = internal.PublisherWrapper(ctx, c.cmp, cfg.Introspection, publisher)
	}

	// Network.
	var (
		NetworkService *servicenetwork.ServiceNetwork
	)
	{
		var err error
		// External communication.
		NetworkService, err = servicenetwork.NewServiceNetwork(cfg, c.cmp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Network")
		}
	}

	// Storage.
	var (
		Coordinator jet.Coordinator
		Pulses      *insolarPulse.DB
		Nodes       *node.StorageDB
		DB          *store.BadgerDB
		Jets        *jet.DBStore
	)
	{
		var err error
		fullDataDirectoryPath, err := filepath.Abs(cfg.Ledger.Storage.DataDirectory)
		if err != nil {
			panic(errors.Wrap(err, "failed to get absolute path for DataDirectory"))
		}
		DB, err = store.NewBadgerDB(
			badger.DefaultOptions(fullDataDirectoryPath),
			store.ValueLogDiscardRatio(cfg.Ledger.Storage.BadgerValueLogGCDiscardRatio),
		)
		if err != nil {
			panic(errors.Wrap(err, "failed to initialize DB"))
		}
		Nodes = node.NewStorageDB(DB)
		Pulses = insolarPulse.NewDB(DB)
		Jets = jet.NewDBStore(DB)

		c := jetcoordinator.NewJetCoordinator(cfg.Ledger.LightChainLimit, *CertManager.GetCertificate().GetNodeRef())
		c.PulseCalculator = Pulses
		c.PulseAccessor = Pulses
		c.JetAccessor = Jets
		c.PlatformCryptographyScheme = CryptoScheme
		c.Nodes = Nodes

		Coordinator = c
	}

	// Communication.
	var (
		WmBus *bus.Bus
	)
	{
		WmBus = bus.NewBus(cfg.Bus, publisher, Pulses, Coordinator, CryptoScheme)
	}

	// API.
	var (
		Requester           *contractrequester.ContractRequester
		ArtifactsClient     = artifacts.NewClient(WmBus)
		AvailabilityChecker = api.NewNetworkChecker(cfg.AvailabilityChecker)
		APIWrapper          *api.RunnerWrapper
	)
	{
		var err error
		Requester, err = contractrequester.New(
			WmBus,
			Pulses,
			Coordinator,
			CryptoScheme,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ContractRequester")
		}

		API, err := api.NewRunner(
			&cfg.APIRunner,
			CertManager,
			Requester,
			NetworkService,
			NetworkService,
			Pulses,
			ArtifactsClient,
			Coordinator,
			NetworkService,
			AvailabilityChecker,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start ApiRunner")
		}

		AdminAPIRunner, err := api.NewRunner(
			&cfg.AdminAPIRunner,
			CertManager,
			Requester,
			NetworkService,
			NetworkService,
			Pulses,
			ArtifactsClient,
			Coordinator,
			NetworkService,
			AvailabilityChecker,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start AdminAPIRunner")
		}

		APIWrapper = api.NewWrapper(API, AdminAPIRunner)
	}

	metricsComp := metrics.NewMetrics(
		cfg.Metrics,
		metrics.GetInsolarRegistry(c.NodeRole),
		c.NodeRole,
	)

	var (
		PulseManager *pulsemanager.PulseManager
		Handler      *handler.Handler
		Genesis      *genesis.Genesis
		Records      *object.RecordDB
		JetKeeper    *executor.DBJetKeeper
	)
	{
		Records = object.NewRecordDB(DB)
		indexes := object.NewIndexDB(DB, Records)
		drops := drop.NewDB(DB)
		JetKeeper = executor.NewJetKeeper(Jets, DB, Pulses)

		backupMaker, err := executor.NewBackupMaker(ctx, DB, cfg.Ledger, JetKeeper.TopSyncPulse(), DB)
		if err != nil {
			return nil, errors.Wrap(err, "failed create backuper")
		}

		c.rollback = executor.NewDBRollback(JetKeeper, drops, Records, indexes, Jets, Pulses, JetKeeper, Nodes, backupMaker)
		c.stateKeeper = executor.NewInitialStateKeeper(JetKeeper, Jets, Coordinator, indexes, drops)

		sp := insolarPulse.NewStartPulse()

		PulseManager = pulsemanager.NewPulseManager(Requester.FlowDispatcher)
		PulseManager.NodeNet = NetworkService
		PulseManager.NodeSetter = Nodes
		PulseManager.Nodes = Nodes
		PulseManager.PulseAppender = Pulses
		PulseManager.PulseAccessor = Pulses
		PulseManager.JetModifier = Jets
		PulseManager.StartPulse = sp
		PulseManager.FinalizationKeeper = executor.NewFinalizationKeeperDefault(JetKeeper, Pulses, cfg.Ledger.LightChainLimit)

		gcRunInfo := executor.NewBadgerGCRunInfo(DB, cfg.Ledger.Storage.GCRunFrequency)
		replicator := executor.NewHeavyReplicatorDefault(Records, indexes, CryptoScheme, Pulses, drops, JetKeeper, backupMaker, Jets, gcRunInfo)
		c.replicator = replicator

		h := handler.New(cfg.Ledger, gcRunInfo)
		h.RecordAccessor = Records
		h.RecordModifier = Records
		h.JetCoordinator = Coordinator
		h.IndexAccessor = indexes
		h.IndexModifier = indexes
		h.DropModifier = drops
		h.PCS = CryptoScheme
		h.PulseAccessor = Pulses
		h.PulseCalculator = Pulses
		h.StartPulse = sp
		h.JetModifier = Jets
		h.JetAccessor = Jets
		h.JetTree = Jets
		h.DropDB = drops
		h.JetKeeper = JetKeeper
		h.InitialStateReader = c.stateKeeper
		h.BackupMaker = backupMaker
		h.Sender = WmBus
		h.Replicator = replicator

		Handler = h

		artifactManager := &artifact.Scope{
			PulseNumber:    pulse.MinTimePulse,
			PCS:            CryptoScheme,
			RecordAccessor: Records,
			RecordModifier: Records,
			IndexModifier:  indexes,
			IndexAccessor:  indexes,
		}
		Genesis = &genesis.Genesis{
			ArtifactManager: artifactManager,
			BaseRecord: &genesis.BaseRecord{
				DB:             DB,
				DropModifier:   drops,
				PulseAppender:  Pulses,
				PulseAccessor:  Pulses,
				RecordModifier: Records,
				IndexModifier:  indexes,
			},

			DiscoveryNodes:  genesisCfg.DiscoveryNodes,
			ContractsConfig: genesisCfg.ContractsConfig,
		}
	}

	// Exporter
	var (
		recordExporter *exporter.RecordServer
		pulseExporter  *exporter.PulseServer
	)
	{
		recordExporter = exporter.NewRecordServer(Pulses, Records, Records, JetKeeper)
		pulseExporter = exporter.NewPulseServer(Pulses, JetKeeper, Nodes)

		grpcServer := grpc.NewServer()
		exporter.RegisterRecordExporterServer(grpcServer, recordExporter)
		exporter.RegisterPulseExporterServer(grpcServer, pulseExporter)

		lis, err := net.Listen("tcp", cfg.Exporter.Addr)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open port for Exporter")
		}

		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				panic(fmt.Errorf("exporter failed to serve: %s", err))
			}
		}()
	}

	c.cmp.Inject(
		DB,
		WmBus,
		Handler,
		PulseManager,
		Jets,
		Pulses,
		Coordinator,
		metricsComp,
		Requester,
		ArtifactsClient,
		APIWrapper,
		AvailabilityChecker,
		KeyProcessor,
		CryptoScheme,
		CryptoService,
		CertManager,
		NetworkService,
		publisher,
	)
	err := c.cmp.Init(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init components")
	}

	if !genesisCfg.Skip {
		if err := Genesis.Start(ctx); err != nil {
			logger.Fatalf("genesis failed on heavy with error: %v", err)
		}
	}

	c.startWatermill(ctx, wmLogger, subscriber, WmBus, NetworkService.SendMessageHandler, Handler.Process, Requester.FlowDispatcher.Process)

	return c, nil
}

func (c *components) Start(ctx context.Context) error {
	err := c.rollback.Start(ctx)
	if err != nil {
		return errors.Wrapf(err, "rollback.Start return error: %s", err.Error())
	}

	err = c.stateKeeper.Start(ctx)
	if err != nil {
		return errors.Wrapf(err, "stateKeeper.Start return error: %s", err.Error())
	}
	return c.cmp.Start(ctx)
}

func (c *components) Stop(ctx context.Context) error {
	err := c.inRouter.Close()
	if err != nil {
		inslogger.FromContext(ctx).Error("Error while closing router", err)
	}
	err = c.outRouter.Close()
	if err != nil {
		inslogger.FromContext(ctx).Error("Error while closing router", err)
	}
	c.replicator.Stop()
	return c.cmp.Stop(ctx)
}

func (c *components) startWatermill(
	ctx context.Context,
	logger watermill.LoggerAdapter,
	sub watermillMsg.Subscriber,
	b *bus.Bus,
	outHandler, inHandler, resultsHandler watermillMsg.NoPublishHandlerFunc,
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
		sub,
		outHandler,
	)

	inRouter.AddMiddleware(
		middleware.InstantAck,
		b.IncomingMessageRouter,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingHandler",
		bus.TopicIncoming,
		sub,
		inHandler,
	)

	inRouter.AddNoPublisherHandler(
		"IncomingRequestResultHandler",
		bus.TopicIncomingRequestResults,
		sub,
		resultsHandler)

	startRouter(ctx, inRouter)
	c.inRouter = inRouter
	startRouter(ctx, outRouter)
	c.outRouter = outRouter
}

func startRouter(ctx context.Context, router *watermillMsg.Router) {
	go func() {
		if err := router.Run(ctx); err != nil {
			inslogger.FromContext(ctx).Error("Error while running router", err)
		}
	}()
	<-router.Running()
}
