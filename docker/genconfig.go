package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/insolar/insolar/configuration"
)

/* Return env variable or default value, it env is not presented */
func GetEnvDefault(key, defaultVal string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	return val
}

/* Return '<host>:<port>' on public docker interface */
func getURI(port uint) string {
	host := GetEnvDefault("IP", "127.0.0.1")
	return fmt.Sprintf("%s:%d", host, port)
}

const (
	defaultConfigPath            = "/etc/insolar/insolard.yaml"
	defaultTranportListenPort    = 13831
	defaultLogLevel              = "info"
	defaultMetricsListenPort     = 8001
	defaultRPCListenPort         = 33001
	defaultInsgorundListenPort   = 33002
	defaultApiListenPort         = 19101
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

	defaultInsgorundListen := "127.0.0.1:" + string(defaultInsgorundListenPort)

	insolardTransportListen := GetEnvDefault("INSOLARD_TRANSPORT_LISTEN", getURI(defaultTranportListenPort))
	insolardLogLevel := GetEnvDefault("INSOLARD_LOG_LEVEL", defaultLogLevel)
	insolardMetricsListen := GetEnvDefault("INSOLARD_METRICS_LISTEN", getURI(defaultMetricsListenPort))
	insolardRPCListen := GetEnvDefault("INSOLARD_RPC_LISTEN", getURI(defaultRPCListenPort))
	insgorundListen := GetEnvDefault("INSGORUND_LISTEN", defaultInsgorundListen)
	insolardAPIListen := GetEnvDefault("INSOLARD_API_LISTEN", getURI(defaultApiListenPort))
	insolardTracerEndpoint := GetEnvDefault("INSOLARD_JAEGER_ENDPOINT", getURI(defaultJaegerEndpointPort))
	insolardTransportFixedAddress := GetEnvDefault("INSOLARD_TRANSPORT_FIXED_ADDRESS", "")

	fmt.Println("[debug] cfg->host->transport->address ==", insolardTransportListen)
	fmt.Println("[debug] cfg->log->level ==", insolardLogLevel)
	fmt.Println("[debug] cfg->log->formatter == json")
	fmt.Println("[debug] cfg->metrics->listenaddress ==", insolardMetricsListen)
	fmt.Println("[debug] cfg->logicrunner->rpclisten ==", insolardRPCListen)
	fmt.Println("[debug] cfg->logicrunner->goplugin->runnerlisten ==", insgorundListen)
	fmt.Println("[debug] cfg->apirunner->address ==", insolardAPIListen)
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
