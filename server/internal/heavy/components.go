// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package heavy

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/dgraph-io/badger"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	component "github.com/insolar/component-manager"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.opencensus.io/stats"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/gbrlsnchs/jwt/v3"

	"github.com/insolar/insolar/api"
	"github.com/insolar/insolar/applicationbase/genesis"
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
	"github.com/insolar/insolar/ledger/heavy/migration"
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

type badgerLogger struct {
	insolar.Logger
}

func (b badgerLogger) Warningf(fmt string, args ...interface{}) {
	b.Warnf(fmt, args...)
}

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

func initTemporaryCertificateManager(ctx context.Context, cfg *configuration.GenericConfiguration) (*certificate.CertificateManager, error) {
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

func initWithPostgres(
	ctx context.Context,
	cfg configuration.ConfigHeavyPg,
	genesisCfg genesis.HeavyConfig,
	genesisOptions genesis.Options,
	genesisOnly bool,
	apiOptions api.Options,
) (*components, error) {
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
		NetworkService, err = servicenetwork.NewServiceNetwork(cfg.GenericConfiguration.Host, c.cmp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Network")
		}
	}

	// Storage.
	var (
		Coordinator    jet.Coordinator
		PulsesPostgres *insolarPulse.PostgresDB
		NodesPostgres  *node.PostgresStorageDB
		Pool           *pgxpool.Pool
		JetsPostgres   *jet.PostgresDBStore
	)
	{
		c := jetcoordinator.NewJetCoordinator(cfg.LightChainLimit, *CertManager.GetCertificate().GetNodeRef())
		c.PlatformCryptographyScheme = CryptoScheme
		Coordinator = c

		var err error
		Pool, err = pgxpool.Connect(context.Background(), cfg.Ledger.PostgreSQL.URL)
		if err != nil {
			panic(errors.Wrap(err, "Unable to connect to PostgreSQL"))
		}

		cwd, err := os.Getwd()
		if err != nil {
			panic(errors.Wrap(err, "os.Getwd failed"))
		}
		path := cfg.Ledger.PostgreSQL.MigrationPath
		logger.Infof("About to run PostgreSQL migration, cwd = %s, migration path = %s", cwd, path)
		ver, err := migration.MigrateDatabase(ctx, Pool, path)
		if err != nil {
			panic(errors.Wrap(err, "Unable to migrate database"))
		}
		logger.Infof("PostgreSQL database migration done, current schema version: %d", ver)

		PulsesPostgres = insolarPulse.NewPostgresDB(Pool)
		NodesPostgres = node.NewPostgresStorageDB(Pool)
		JetsPostgres = jet.NewPostgresDBStore(Pool)

		c.PulseCalculator = PulsesPostgres
		c.PulseAccessor = PulsesPostgres
		c.JetAccessor = JetsPostgres
		c.Nodes = NodesPostgres
	}

	// Communication.
	var (
		WmBus *bus.Bus
	)
	{
		WmBus = bus.NewBus(cfg.Bus, publisher, PulsesPostgres, Coordinator, CryptoScheme)
		WmBus = bus.NewBus(
			cfg.Bus,
			publisher,
			PulsesPostgres,
			Coordinator, CryptoScheme)
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
			PulsesPostgres,
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
			PulsesPostgres,
			ArtifactsClient,
			Coordinator,
			NetworkService,
			AvailabilityChecker,
			apiOptions,
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
			PulsesPostgres,
			ArtifactsClient,
			Coordinator,
			NetworkService,
			AvailabilityChecker,
			apiOptions,
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

		RecordsPostgres   *object.PostgresRecordDB
		IndexesPostgres   *object.PostgresIndexDB
		DropPostgres      *drop.PostgresDB
		PostgresJetKeeper *executor.PostgresDBJetKeeper
	)
	{
		RecordsPostgres = object.NewPostgresRecordDB(Pool)
		IndexesPostgres = object.NewPostgresIndexDB(Pool, RecordsPostgres)
		DropPostgres = drop.NewPostgresDB(Pool)
		PostgresJetKeeper = executor.NewPostgresJetKeeper(JetsPostgres, Pool, PulsesPostgres)

		c.rollback = executor.NewDBRollback(
			PostgresJetKeeper,
			DropPostgres,
			RecordsPostgres,
			IndexesPostgres,
			JetsPostgres,
			PulsesPostgres,
			PostgresJetKeeper,
			NodesPostgres)

		c.stateKeeper = executor.NewInitialStateKeeper(
			PostgresJetKeeper,
			JetsPostgres,
			Coordinator,
			IndexesPostgres,
			DropPostgres)

		sp := insolarPulse.NewStartPulse()

		PulseManager = pulsemanager.NewPulseManager(Requester.FlowDispatcher)
		PulseManager.NodeNet = NetworkService

		PulseManager.NodeSetter = NodesPostgres
		PulseManager.Nodes = NodesPostgres
		PulseManager.PulseAppender = PulsesPostgres
		PulseManager.PulseAccessor = PulsesPostgres
		PulseManager.JetModifier = JetsPostgres
		PulseManager.FinalizationKeeper = executor.NewFinalizationKeeperDefault(PostgresJetKeeper, PulsesPostgres, cfg.LightChainLimit)

		replicator := executor.NewHeavyReplicatorDefault(
			RecordsPostgres,
			IndexesPostgres,
			CryptoScheme,
			PulsesPostgres,
			DropPostgres,
			PostgresJetKeeper,
			&executor.PostgresBackupMaker{},
			JetsPostgres,
			&executor.PostgresGCRunInfo{},
		)
		c.replicator = replicator

		PulseManager.StartPulse = sp

		h := handler.New(cfg.LightChainLimit, &executor.PostgresGCRunInfo{})
		h.RecordAccessor = RecordsPostgres
		h.RecordModifier = RecordsPostgres
		h.JetCoordinator = Coordinator
		h.IndexAccessor = IndexesPostgres
		h.IndexModifier = IndexesPostgres
		h.DropModifier = DropPostgres
		h.PCS = CryptoScheme
		h.PulseAccessor = PulsesPostgres
		h.PulseCalculator = PulsesPostgres
		h.StartPulse = sp
		h.JetModifier = JetsPostgres
		h.JetAccessor = JetsPostgres
		h.JetTree = JetsPostgres
		h.JetKeeper = PostgresJetKeeper
		h.InitialStateReader = c.stateKeeper
		h.Sender = WmBus
		h.Replicator = replicator

		Handler = h

		artifactManager := &artifact.Scope{
			PulseNumber:    pulse.MinTimePulse,
			PCS:            CryptoScheme,
			RecordAccessor: RecordsPostgres,
			RecordModifier: RecordsPostgres,
			IndexModifier:  IndexesPostgres,
			IndexAccessor:  IndexesPostgres,
		}
		Genesis = &genesis.Genesis{
			ArtifactManager: artifactManager,
			IndexModifier:   IndexesPostgres,
			BaseRecord: &genesis.PostgresBaseRecord{
				Pool:           Pool,
				DropModifier:   DropPostgres,
				PulseAppender:  PulsesPostgres,
				PulseAccessor:  PulsesPostgres,
				RecordModifier: RecordsPostgres,
				IndexModifier:  IndexesPostgres,
			},

			DiscoveryNodes: genesisCfg.DiscoveryNodes,
			GenesisOptions: genesisOptions,
		}
	}

	// Exporter
	var (
		recordExporter *exporter.RecordServer
		pulseExporter  *exporter.PulseServer
	)
	{
		recordExporter = exporter.NewRecordServer(PulsesPostgres, RecordsPostgres, RecordsPostgres, PostgresJetKeeper, cfg.Exporter.Auth)
		pulseExporter = exporter.NewPulseServer(PulsesPostgres, PostgresJetKeeper, NodesPostgres, cfg.Exporter.Auth)

		grpcServer, err := newGRPCServer(cfg.Exporter)
		if err != nil {
			return nil, errors.Wrap(err, "failed to initiate a GRPC server")
		}
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
		WmBus,
		Handler,
		PulseManager,
		JetsPostgres,
		PulsesPostgres,
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

	if genesisOnly {
		logger.Info("Terminating, because --genesis-only flag was specified.")
		os.Exit(0)
	}

	c.startWatermill(ctx, wmLogger, subscriber, WmBus, NetworkService.SendMessageHandler, Handler.Process, Requester.FlowDispatcher.Process)

	return c, nil
}

func initWithBadger(
	ctx context.Context,
	cfg configuration.HeavyBadgerConfig,
	genesisCfg genesis.HeavyConfig,
	genesisOptions genesis.Options,
	genesisOnly bool,
	apiOptions api.Options,
) (*components, error) {
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
		NetworkService, err = servicenetwork.NewServiceNetwork(cfg.GenericConfiguration.Host, c.cmp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to start Network")
		}
	}

	// Storage.
	var (
		Coordinator jet.Coordinator
		Pulses      *insolarPulse.BadgerDB
		Nodes       *node.BadgerStorageDB
		DB          *store.BadgerDB
		Jets        *jet.BadgerDBStore
	)
	{
		var err error
		startTime := time.Now()
		fullDataDirectoryPath, err := filepath.Abs(cfg.Ledger.Storage.DataDirectory)
		if err != nil {
			panic(errors.Wrap(err, "failed to get absolute path for DataDirectory"))
		}
		options := badger.DefaultOptions(fullDataDirectoryPath)
		options.Logger = badgerLogger{Logger: logger.WithField("component", "badger")}
		DB, err = store.NewBadgerDB(
			options,
			store.ValueLogDiscardRatio(cfg.Ledger.Storage.BadgerValueLogGCDiscardRatio),
			store.OpenAndCloseBadgerOnStart(true),
		)
		if err != nil {
			panic(errors.Wrap(err, "failed to initialize DB"))
		}
		Nodes = node.NewBadgerStorageDB(DB)
		Pulses = insolarPulse.NewBadgerDB(DB)
		Jets = jet.NewBadgerDBStore(DB)

		timeBadgerStarted := time.Since(startTime)
		stats.Record(ctx, statBadgerStartTime.M(float64(timeBadgerStarted.Nanoseconds())/1e6))
		logger.Info("badger starts in ", timeBadgerStarted)

		c := jetcoordinator.NewJetCoordinator(cfg.LightChainLimit, *CertManager.GetCertificate().GetNodeRef())
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
			apiOptions,
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
			apiOptions,
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
		Records      *object.BadgerRecordDB
		JetKeeper    *executor.BadgerDBJetKeeper
	)
	{
		Records = object.NewBadgerRecordDB(DB)
		indexes := object.NewBadgerIndexDB(DB, Records)
		drops := drop.NewBadgerDB(DB)
		JetKeeper = executor.NewBadgerJetKeeper(Jets, DB, Pulses)

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
		PulseManager.FinalizationKeeper = executor.NewFinalizationKeeperDefault(JetKeeper, Pulses, cfg.LightChainLimit)

		gcRunInfo := executor.NewBadgerGCRunInfo(DB, cfg.Ledger.Storage.GCRunFrequency)
		replicator := executor.NewHeavyReplicatorDefault(Records, indexes, CryptoScheme, Pulses, drops, JetKeeper, backupMaker, Jets, gcRunInfo)
		c.replicator = replicator

		h := handler.New(cfg.LightChainLimit, gcRunInfo)
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
			IndexModifier:   indexes,
			BaseRecord: &genesis.BadgerBaseRecord{
				DB:             DB,
				DropModifier:   drops,
				PulseAppender:  Pulses,
				PulseAccessor:  Pulses,
				RecordModifier: Records,
				IndexModifier:  indexes,
			},

			DiscoveryNodes: genesisCfg.DiscoveryNodes,
			GenesisOptions: genesisOptions,
		}
	}

	// Exporter
	var (
		recordExporter *exporter.RecordServer
		pulseExporter  *exporter.PulseServer
	)
	{
		recordExporter = exporter.NewRecordServer(Pulses, Records, Records, JetKeeper, cfg.Exporter.Auth)
		pulseExporter = exporter.NewPulseServer(Pulses, JetKeeper, Nodes, cfg.Exporter.Auth)

		grpcServer, err := newGRPCServer(cfg.Exporter)
		if err != nil {
			return nil, errors.Wrap(err, "failed to initiate a GRPC server")
		}
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

	if genesisOnly {
		logger.Info("Terminating, because --genesis-only flag was specified.")
		os.Exit(0)
	}

	c.startWatermill(ctx, wmLogger, subscriber, WmBus, NetworkService.SendMessageHandler, Handler.Process, Requester.FlowDispatcher.Process)

	return c, nil
}

var (
	jwtIss string
	jwtKey *jwt.HMACSHA
)

func newGRPCServer(cfg configuration.Exporter) (*grpc.Server, error) {
	if cfg.Auth.Required {
		jwtIss = cfg.Auth.Issuer
		key := []byte(cfg.Auth.Secret)
		if len(key) != 64 {
			return nil, errors.New("exporter.auth.secret must be 512-bit")
		}
		jwtKey = jwt.NewHS512(key)
		return grpc.NewServer(grpc.UnaryInterceptor(authUnaryIntcp), grpc.StreamInterceptor(authStreamIntcp)), nil
	}
	return grpc.NewServer(), nil
}

func newComponents(
	ctx context.Context,
	cfg configuration.ConfigHolder,
	genesisCfg genesis.HeavyConfig,
	genesisOptions genesis.Options,
	genesisOnly bool,
	apiOptions api.Options,
) (*components, error) {
	heavyCfg := cfg.GetNodeConfig()
	switch realCfg := heavyCfg.(type) {
	case *configuration.ConfigHeavyPg:
		return initWithPostgres(ctx, *realCfg, genesisCfg, genesisOptions, genesisOnly, apiOptions)
	case *configuration.HeavyBadgerConfig:
		return initWithBadger(ctx, *realCfg, genesisCfg, genesisOptions, genesisOnly, apiOptions)
	}
	return nil, errors.New("can't start heavy, db configuration error")
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

func authUnaryIntcp(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	newCtx, err := authorize(ctx)
	if err != nil {
		return nil, err
	}
	return handler(newCtx, req)
}

func authStreamIntcp(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	newCtx, err := authorize(stream.Context())
	if err != nil {
		return err
	}
	wrapped := grpc_middleware.WrapServerStream(stream)
	wrapped.WrappedContext = newCtx
	return handler(srv, wrapped)
}

func authorize(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.InvalidArgument, "failed to retrieve metadata")
	}

	elem, ok := md["authorization"]
	if !ok || len(elem) == 0 {
		return nil, status.Error(codes.InvalidArgument, "auth data not supplied")
	}
	sub, err := validateJWT(strings.TrimPrefix(elem[0], "Bearer "))
	if err != nil {
		return nil, err
	}

	newMD := md.Copy()
	newMD.Set(exporter.ObsID, sub)
	return metadata.NewIncomingContext(ctx, newMD), nil
}

func validateJWT(token string) (string, error) {
	if jwtIss == "" || jwtKey == nil {
		return "", status.Error(codes.Internal, "required auth parameters are not configured")
	}

	payload := jwt.Payload{}
	_, err := jwt.Verify([]byte(token), jwtKey, &payload)
	if err != nil {
		return "", status.Error(codes.Unauthenticated, err.Error())
	}

	if payload.Issuer != jwtIss {
		return "", status.Errorf(codes.Unauthenticated, "unknown JWT issuer: %s", payload.Issuer)
	}

	if payload.ExpirationTime.Before(time.Now()) {
		return "", status.Error(codes.Unauthenticated, "JWT is expired")
	}

	return payload.Subject, nil
}
