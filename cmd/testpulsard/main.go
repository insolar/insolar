/*
 *    Copyright 2019 Insolar Technologies
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
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/insolar/component-manager"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/pulsenetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/insolar/insolar/version"
)

type inputParams struct {
	configPath string
	port       string
}

func parseInputParams() inputParams {
	var rootCmd = &cobra.Command{Use: "insolard"}
	var result inputParams
	rootCmd.Flags().StringVarP(&result.configPath, "config", "c", "", "path to config file")
	rootCmd.Flags().StringVarP(&result.port, "port", "port", "", "port for test pulsar")
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Wrong input params:", err.Error())
	}

	return result
}

func main() {
	params := parseInputParams()

	jww.SetStdoutThreshold(jww.LevelDebug)
	vp := viper.New()
	pCfg := configuration.NewPulsarConfiguration()
	if len(params.configPath) != 0 {
		vp.SetConfigFile(params.configPath)
	}
	err := vp.ReadInConfig()
	if err != nil {
		log.Warn("failed to load configuration from file: ", err.Error())
	}
	err = vp.Unmarshal(&pCfg)
	if err != nil {
		log.Warn("failed to load configuration from file: ", err.Error())
	}

	ctx := context.Background()
	ctx, _ = inslogger.InitNodeLogger(ctx, pCfg.Log, "", "test_pulsar")
	testPulsar := initPulsar(ctx, pCfg)

	http.HandleFunc("/pulse", func(writer http.ResponseWriter, request *http.Request) {
		err := testPulsar.SendPulse(ctx)
		if err != nil {
			_, err := fmt.Fprintf(writer, "Error - %v", err)
			if err != nil {
				panic(err)
			}
		}

		_, err = fmt.Fprint(writer, "OK")
		if err != nil {
			panic(err)
		}
	})

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(params.port, nil); err != nil {
		panic(err)
	}
}

func initPulsar(ctx context.Context, cfg configuration.PulsarConfiguration) *pulsar.TestPulsar {
	fmt.Println("Version: ", version.GetFullVersion())
	fmt.Println("Starts with configuration:\n", configuration.ToString(cfg))

	keyStore, err := keystore.NewKeyStore(cfg.KeysPath)
	if err != nil {
		panic(err)
	}
	cryptographyScheme := platformpolicy.NewPlatformCryptographyScheme()
	cryptographyService := cryptography.NewCryptographyService()
	keyProcessor := platformpolicy.NewKeyProcessor()

	pulseDistributor, err := pulsenetwork.NewDistributor(cfg.Pulsar.PulseDistributor)
	if err != nil {
		panic(err)
	}

	cm := component.NewManager(nil)
	cm.Register(cryptographyScheme, keyStore, keyProcessor, transport.NewFactory(cfg.Pulsar.DistributionTransport))
	cm.Inject(cryptographyService, pulseDistributor)

	if err = cm.Init(ctx); err != nil {
		panic(err)
	}

	if err = cm.Start(ctx); err != nil {
		panic(err)
	}

	return pulsar.NewTestPulsar(cfg.Pulsar, pulseDistributor, &entropygenerator.StandardEntropyGenerator{})
}
