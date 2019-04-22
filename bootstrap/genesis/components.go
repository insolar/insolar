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

package genesis

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/internal/ledger/store"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/platformpolicy"
)

type bootstrapComponents struct {
	CryptographyService        insolar.CryptographyService
	PlatformCryptographyScheme insolar.PlatformCryptographyScheme
	KeyProcessor               insolar.KeyProcessor
}

func initBootstrapComponents(ctx context.Context, cfg configuration.Configuration) bootstrapComponents {
	keyStore, err := keystore.NewKeyStore(cfg.KeysPath)
	checkError(ctx, err, "failed to load KeyStore: ")

	platformCryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	keyProcessor := platformpolicy.NewKeyProcessor()

	cryptographyService := cryptography.NewCryptographyService()

	earlyComponents := component.Manager{}
	earlyComponents.Register(platformCryptographyScheme, keyStore)
	earlyComponents.Inject(cryptographyService, keyProcessor)

	return bootstrapComponents{
		CryptographyService:        cryptographyService,
		PlatformCryptographyScheme: platformCryptographyScheme,
		KeyProcessor:               keyProcessor,
	}
}

func createCertificateManager(
	ctx context.Context,
	cryptographyService insolar.CryptographyService,
	keyProcessor insolar.KeyProcessor,
) *certificate.CertificateManager {

	publicKey, err := cryptographyService.GetPublicKey()
	checkError(ctx, err, "failed to retrieve node public key")

	certManager, err := certificate.NewManagerCertificateWithKeys(publicKey, keyProcessor)
	checkError(ctx, err, "failed to start Certificate (bootstrap mode)")

	return certManager
}

type storageComponents struct {
	storeBadgerDB *store.BadgerDB

	dropDB   *drop.DB
	blobDB   *blob.DB
	recordDB *object.RecordDB
	pulseDB  *pulse.DB
	indexDB  *object.IndexDB
}

func initStorageComponents(conf configuration.Ledger) storageComponents {
	db, err := store.NewBadgerDB(conf.Storage.DataDirectory)
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize DB"))
	}

	return storageComponents{
		storeBadgerDB: db,

		dropDB:   drop.NewDB(db),
		blobDB:   blob.NewDB(db),
		recordDB: object.NewRecordDB(db),
		pulseDB:  pulse.NewDB(db),
		indexDB:  object.NewIndexDB(db),
	}
}
