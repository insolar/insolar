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
	"fmt"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
)

type Initializer struct {
	cfgPath           string
	genesisConfigPath string
	genesisKeyOut     string
}

func NewInitializer(cfgPath string, genesisConfigPath, genesisKeyOut string) *Initializer {
	return &Initializer{
		cfgPath:           cfgPath,
		genesisConfigPath: genesisConfigPath,
		genesisKeyOut:     genesisKeyOut,
	}
}

func (s *Initializer) Init() {
	cfgHolder := configuration.NewHolder()
	var err error
	if len(s.cfgPath) != 0 {
		err = cfgHolder.LoadFromFile(s.cfgPath)
	} else {
		err = cfgHolder.Load()
	}
	if err != nil {
		log.Warn("failed to load configuration from file: ", err.Error())
	}

	cfg := &cfgHolder.Configuration

	ctx, inslog := initLogger(context.Background(), cfg.Log)
	log.SetGlobalLogger(inslog)
	fmt.Println("Starts with configuration:\n", configuration.ToString(cfgHolder.Configuration))

	bootstrapComponents := initBootstrapComponents(ctx, *cfg)
	certManager := initCertificateManager(
		ctx,
		*cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.KeyProcessor,
	)

	cm, err := initComponents(
		ctx,
		*cfg,
		bootstrapComponents.CryptographyService,
		bootstrapComponents.PlatformCryptographyScheme,
		bootstrapComponents.KeyStore,
		bootstrapComponents.KeyProcessor,
		certManager,
		s.genesisConfigPath,
		s.genesisKeyOut,
	)
	checkError(ctx, err, "failed to init components")

	err = cm.Init(ctx)
	checkError(ctx, err, "failed to init components")

	err = cm.Start(ctx)
	checkError(ctx, err, "failed to start components")

	err = cm.Stop(ctx)
	checkError(ctx, err, "failed to stop components")
}

func initLogger(ctx context.Context, cfg configuration.Log) (context.Context, insolar.Logger) {
	inslog, err := log.NewLog(cfg)
	if err != nil {
		panic(err)
	}

	if newInslog, err := inslog.WithLevel(cfg.Level); err != nil {
		inslog.Error(err.Error())
	} else {
		inslog = newInslog
	}

	ctx = inslogger.SetLogger(ctx, inslog)
	return ctx, inslog
}

func checkError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}
	inslogger.FromContext(ctx).Fatalf("%v: %v", message, err.Error())
}
