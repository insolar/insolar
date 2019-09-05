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

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/insolar/insolar/configuration"
)

// getenv + default value
func GetEnvDefault(key, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}

// Return '<host>:<port>' on public docker interface
func getURI(port uint) string {
	host := GetEnvDefault("IP", "127.0.0.1")
	return fmt.Sprintf("%s:%d", host, port)
}

const (
	defaultConfigPath            = "/etc/insolar/insolard.yaml"
	defaultTranportListenPort    = 7900
	defaultLogLevel              = "info"
	defaultMetricsListenPort     = 8001
	defaultRPCListenPort         = 18182
	defaultInsgorundListenPort   = 18181
	defaultAPIListenPort         = 19191
	defaultAdminAPIListenPort    = 19091
	defaultJaegerEndpointPort    = 6831
	defaultKeysPath              = "/etc/insolar/keys.json"
	defaultCertPath              = "/etc/insolar/cert.json"
	defaultDataDir               = "/var/lib/insolar/"
	defaultTransportFixedAddress = ""
)

func main() {
	hld := configuration.NewHolder()
	err := hld.LoadFromFile(defaultConfigPath)
	if err != nil {
		fmt.Println("Failed to open configuration:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	cfg := hld.Configuration

	insgorundListen := GetEnvDefault("INSGORUND_ENDPOINT", "insgorund:"+string(defaultInsgorundListenPort))

	insolardMetricsListen := getURI(defaultMetricsListenPort)
	insolardRPCListen := getURI(defaultRPCListenPort)
	insolardAPIListen := getURI(defaultAPIListenPort)
	insolardAdminAPIListen := getURI(defaultAdminAPIListenPort)

	insolardTransportListen := GetEnvDefault("INSOLARD_TRANSPORT_LISTEN", getURI(defaultTranportListenPort))
	insolardLogLevel := GetEnvDefault("INSOLARD_LOG_LEVEL", defaultLogLevel)
	insolardTracerEndpoint := GetEnvDefault("INSOLARD_JAEGER_ENDPOINT", getURI(defaultJaegerEndpointPort))
	insolardTransportFixedAddress := GetEnvDefault("INSOLARD_TRANSPORT_FIXED_ADDRESS", defaultTransportFixedAddress)

	fmt.Println("[debug] cfg->host->transport->address ==", insolardTransportListen)
	fmt.Println("[debug] cfg->log->level ==", insolardLogLevel)
	fmt.Println("[debug] cfg->log->formatter == json")
	fmt.Println("[debug] cfg->metrics->listenaddress ==", insolardMetricsListen)
	fmt.Println("[debug] cfg->logicrunner->rpclisten ==", insolardRPCListen)
	fmt.Println("[debug] cfg->logicrunner->goplugin->runnerlisten ==", insgorundListen)
	fmt.Println("[debug] cfg->apirunner->address ==", insolardAPIListen)
	fmt.Println("[debug] cfg->adminapirunner->address ==", insolardAdminAPIListen)
	fmt.Println("[debug] cfg->tracer->jaeger->agentendpoint ==", insolardTracerEndpoint)

	// transport related
	cfg.Host.Transport.Address = insolardTransportListen
	cfg.Host.Transport.FixedPublicAddress = insolardTransportFixedAddress
	// logger related
	cfg.Log.Level = insolardLogLevel
	// metrics related
	cfg.Metrics.ListenAddress = insolardMetricsListen
	cfg.Log.Formatter = "json" // ??
	// logic runner related
	cfg.LogicRunner.RPCListen = insolardRPCListen
	cfg.LogicRunner.GoPlugin.RunnerListen = insgorundListen
	// api runner related
	cfg.APIRunner.Address = insolardAPIListen
	cfg.AdminAPIRunner.Address = insolardAdminAPIListen
	// with tracer
	cfg.Tracer.Jaeger.AgentEndpoint = insolardTracerEndpoint
	// unstructured
	cfg.KeysPath = defaultKeysPath
	cfg.CertificatePath = defaultCertPath
	cfg.Ledger.Storage.DataDirectory = defaultDataDir

	data, err := yaml.Marshal(cfg)
	if err != nil {
		fmt.Println("Failed to marshall configuration:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	err = ioutil.WriteFile(defaultConfigPath, data, 0666)
	if err != nil {
		fmt.Println("Failed to save configuration:")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
