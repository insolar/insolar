package main

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/servicenetwork"
	"github.com/insolar/insolar/network/storage"
	"github.com/insolar/insolar/network/termination"
	"github.com/insolar/insolar/platformpolicy"
)

type publisherMock struct {
	Error error
}

func (pm *publisherMock) Publish(topic string, messages ...*message.Message) error { return pm.Error }
func (pm *publisherMock) Close() error                                             { return nil }

func main() {
	cfg := configuration.NewConfiguration()
	cfg.KeysPath = "cmd/insolard/testdata/bootstrap_keys.json"
	cfg.CertificatePath = "cmd/insolard/testdata/certificate.json"

	cm := component.NewManager(nil)

	db, err := storage.NewBadgerDB(cfg.Service)
	if err != nil {
		log.Fatal(err.Error())
	}

	net, err := servicenetwork.NewServiceNetwork(cfg, cm, true)
	if err != nil {
		log.Fatal(err.Error())
	}

	ks, err := keystore.NewKeyStore(cfg.KeysPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	cryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	cryptographyService := cryptography.NewCryptographyService()
	keyProcessor := platformpolicy.NewKeyProcessor()

	cm.Inject(cryptographyService, cryptographyScheme, keyProcessor, ks)

	publicKey, err := cryptographyService.GetPublicKey()
	if err != nil {
		log.Fatal(err.Error())
	}

	CertManager, err := certificate.NewManagerReadCertificate(publicKey, keyProcessor, cfg.CertificatePath)

	//CertManager := certificate.NewCertificateManager(nil)
	NodeNetwork, err := nodenetwork.NewNodeNetwork(cfg.Host.Transport, CertManager.GetCertificate())
	if err != nil {
		log.Fatal(err.Error())
	}

	ctx := context.Background()
	cm.Inject(net,
		CertManager,
		db,
		storage.NewPulseStorage(),
		cryptographyService,
		NodeNetwork,
		termination.NewHandler(net),
		&publisherMock{},
		//testutils.NewPulseManagerMock(nil),
		//testutils.NewPulseAccessor // todo make adapter for Accessor from storage
	)

	err = cm.Init(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = cm.Start(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

}
