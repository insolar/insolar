// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/insolar/insolar/applicationbase/bootstrap"
	pulsewatcher "github.com/insolar/insolar/cmd/pulsewatcher/config"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/defaults"
	"github.com/insolar/insolar/log"
)

func baseDir() string {
	return defaults.LaunchnetDir()
}

var (
	defaultOutputConfigNameTmpl      = "%d/insolard.yaml"
	defaultHost                      = "127.0.0.1"
	defaultJaegerEndPoint            = ""
	discoveryDataDirectoryTemplate   = withBaseDir("discoverynodes/%d/data")
	discoveryCertificatePathTemplate = withBaseDir("discoverynodes/certs/discovery_cert_%d.json")
	nodeDataDirectoryTemplate        = "nodes/%d/data"
	nodeCertificatePathTemplate      = "nodes/%d/cert.json"
	pulsewatcherFileName             = withBaseDir("pulsewatcher.yaml")

	prometheusConfigTmpl = "scripts/prom/server.yml.tmpl"
	prometheusFileName   = "prometheus.yaml"

	bootstrapConfigTmpl = "scripts/insolard/bootstrap_template.yaml"
	bootstrapFileName   = withBaseDir("bootstrap.yaml")

	pulsardConfigTmpl = "scripts/insolard/pulsar_template.yaml"
	pulsardFileName   = withBaseDir("pulsar.yaml")

	keeperdConfigTmpl = "scripts/insolard/keeperd_template.yaml"
	keeperdFileName   = withBaseDir("keeperd.yaml")

	insolardDefaultsConfigWithBadger   = "scripts/insolard/defaults/insolard_badger.yaml"
	insolardDefaultsConfigWithPostgres = "scripts/insolard/defaults/insolard_postgres.yaml"
)

var (
	outputDir       string
	debugLevel      string
	gorundPortsPath string
)

func parseInputParams() {
	var rootCmd = &cobra.Command{}

	rootCmd.Flags().StringVarP(
		&outputDir, "output", "o", baseDir(), "output directory")
	rootCmd.Flags().StringVarP(
		&debugLevel, "debuglevel", "d", "Debug", "debug level")
	rootCmd.Flags().StringVarP(
		&gorundPortsPath, "gorundports", "p", "", "path to insgorund ports (required)")

	err := rootCmd.Execute()
	check("Wrong input params:", err)

	if gorundPortsPath == "" {
		err := rootCmd.Usage()
		check("[ parseInputParams ]", err)
	}
}

func writeGorundPorts(gorundPorts [][]string) {
	var portsData string
	for _, ports := range gorundPorts {
		portsData += ports[0] + " " + ports[1] + "\n"
	}
	err := makeFileWithDir("./", gorundPortsPath, portsData)
	check("failed to create gorund ports file: "+gorundPortsPath, err)
}

func writeInsolardConfigs(dir string, insolardConfigs []interface{}) {
	fmt.Println("generate_insolar_configs.go: writeInsolardConfigs...")

	for index, conf := range insolardConfigs {
		data := configuration.ToString(conf)

		fileName := fmt.Sprintf(defaultOutputConfigNameTmpl, index+1)
		fileName = filepath.Join(dir, fileName)
		err := createFileWithDir(fileName, data)
		check("failed to create insolard config: "+fileName, err)
	}
}

func main() {
	fmt.Println("[main] about to call parseInputParams()")
	parseInputParams()

	fmt.Println("[main] about to call mustMakeDir()")
	mustMakeDir(outputDir)
	fmt.Println("[main] about to call writeGenesisConfig()")
	writeGenesisConfig()

	bootstrapConf, err := bootstrap.ParseConfig(bootstrapFileName)
	check("Can't read bootstrap config", err)

	pwConfig := pulsewatcher.Config{}
	discoveryNodesConfigs := make([]interface{}, 0, len(bootstrapConf.DiscoveryNodes))

	var gorundPorts [][]string

	promVars := &promConfigVars{
		Jobs: map[string][]string{},
	}

	fmt.Println("[main] about to enter for loop which calls newDefaultInsolardConfig() first time")

	// process discovery nodes
	for index, node := range bootstrapConf.DiscoveryNodes {
		nodeIndex := index + 1

		genericConfiguration := createGenericConfigDiscovery(node, nodeIndex, bootstrapConf)

		switch node.Role {
		case "light_material":
			configLight := createLightConfig()
			configLight.GenericConfiguration = genericConfiguration

			discoveryNodesConfigs = append(discoveryNodesConfigs, configLight)
			promVars.addTarget("insolard", configLight.Metrics.ListenAddress)
			pwConfig.Nodes = append(pwConfig.Nodes, configLight.AdminAPIRunner.Address)

		case "virtual":
			configVirtual := createVirtualConfig(index, nodeIndex)
			configVirtual.GenericConfiguration = genericConfiguration

			rpcListenPort := 33300 + (index+nodeIndex)*nodeIndex
			gorundPorts = append(gorundPorts, []string{strconv.Itoa(rpcListenPort - 1), strconv.Itoa(rpcListenPort)})

			discoveryNodesConfigs = append(discoveryNodesConfigs, configVirtual)
			promVars.addTarget("insolard", configVirtual.Metrics.ListenAddress)
			pwConfig.Nodes = append(pwConfig.Nodes, configVirtual.AdminAPIRunner.Address)

		case "heavy_material":
			fmt.Println("[createHeavyConfig] os.Getenv == ", os.Getenv("POSTGRES_ENABLE"))

			if len(os.Getenv("POSTGRES_ENABLE")) > 0 {
				fmt.Println("[createHeavyConfig] Using PostgreSQL config")
				holder := *configuration.NewHolderHeavyPg(insolardDefaultsConfigWithPostgres).MustLoad()
				conf := holder.Configuration
				conf.GenericConfiguration = genericConfiguration

				discoveryNodesConfigs = append(discoveryNodesConfigs, conf)
				promVars.addTarget("insolard", conf.Metrics.ListenAddress)
				pwConfig.Nodes = append(pwConfig.Nodes, conf.AdminAPIRunner.Address)
			} else {
				fmt.Println("[newDefaultInsolardConfig] Using Badger config")
				holder := *configuration.NewHolderHeavyBadger(insolardDefaultsConfigWithBadger).MustLoad()
				conf := holder.Configuration
				conf.Ledger.Storage.DataDirectory = fmt.Sprintf(discoveryDataDirectoryTemplate, nodeIndex)
				conf.GenericConfiguration = genericConfiguration

				discoveryNodesConfigs = append(discoveryNodesConfigs, conf)
				promVars.addTarget("insolard", conf.Metrics.ListenAddress)
				pwConfig.Nodes = append(pwConfig.Nodes, conf.AdminAPIRunner.Address)
			}
		}
	}

	fmt.Println("[main] leaving the loop which calls newDefaultInsolardConfig() first time")

	// process extra nodes
	nodeDataDirectoryTemplate = filepath.Join(outputDir, nodeDataDirectoryTemplate)
	nodeCertificatePathTemplate = filepath.Join(outputDir, nodeCertificatePathTemplate)

	nodesConfigs := make([]interface{}, 0, len(bootstrapConf.Nodes))
	for index, node := range bootstrapConf.Nodes {
		nodeIndex := index + 1

		genericConfiguration := createGenericConfigJoiner(node, nodeIndex, bootstrapConf)

		if node.Role == "heavy_material" {
			if len(os.Getenv("POSTGRES_ENABLE")) > 0 {
				heavyConfig := configuration.NewConfigurationHeavyPg()
				heavyConfig.GenericConfiguration = genericConfiguration

				nodesConfigs = append(nodesConfigs, genericConfiguration)
				promVars.addTarget("insolard", genericConfiguration.Metrics.ListenAddress)
				pwConfig.Nodes = append(pwConfig.Nodes, genericConfiguration.AdminAPIRunner.Address)
			} else {
				heavyConfig := configuration.NewConfigurationHeavyBadger()
				heavyConfig.Ledger.Storage.DataDirectory = fmt.Sprintf(nodeDataDirectoryTemplate, nodeIndex)
				heavyConfig.GenericConfiguration = genericConfiguration

				nodesConfigs = append(nodesConfigs, genericConfiguration)
				promVars.addTarget("insolard", genericConfiguration.Metrics.ListenAddress)
				pwConfig.Nodes = append(pwConfig.Nodes, genericConfiguration.AdminAPIRunner.Address)
			}

		}

		if node.Role == "virtual" {
			virtualConfig := configuration.NewConfigurationVirtual()
			rpcListenPort := 34300 + (index+nodeIndex+len(bootstrapConf.DiscoveryNodes)+1)*nodeIndex
			virtualConfig.LogicRunner = configuration.NewLogicRunner()
			virtualConfig.LogicRunner.GoPlugin.RunnerListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort-1)
			virtualConfig.LogicRunner.RPCListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort)
			virtualConfig.GenericConfiguration = genericConfiguration

			gorundPorts = append(gorundPorts, []string{strconv.Itoa(rpcListenPort - 1), strconv.Itoa(rpcListenPort)})
			nodesConfigs = append(nodesConfigs, genericConfiguration)
			promVars.addTarget("insolard", genericConfiguration.Metrics.ListenAddress)
			pwConfig.Nodes = append(pwConfig.Nodes, genericConfiguration.AdminAPIRunner.Address)
		}
	}

	writePromConfig(promVars)
	writeInsolardConfigs(filepath.Join(outputDir, "/discoverynodes"), discoveryNodesConfigs)
	writeInsolardConfigs(filepath.Join(outputDir, "/nodes"), nodesConfigs)
	writeGorundPorts(gorundPorts)

	pulsarConf := &pulsarConfigVars{}
	pulsarConf.DataDir = withBaseDir("pulsar_data")
	pulsarConf.BaseDir = baseDir()
	for _, node := range bootstrapConf.DiscoveryNodes {
		pulsarConf.BootstrapHosts = append(pulsarConf.BootstrapHosts, node.Host)
	}
	pulsarConf.AgentEndpoint = defaultJaegerEndPoint
	writePulsarConfig(pulsarConf)

	writeKeeperdConfig()

	pwConfig.Interval = 500 * time.Millisecond
	pwConfig.Timeout = 1 * time.Second
	mustMakeDir(filepath.Dir(pulsewatcherFileName))
	err = pulsewatcher.WriteConfig(pulsewatcherFileName, pwConfig)
	check("couldn't write pulsewatcher config file", err)
	fmt.Println("generate_insolar_configs.go: write to file", pulsewatcherFileName)
}

func createGenericConfigDiscovery(node bootstrap.Node, nodeIndex int, bootstrapConf *bootstrap.Config) configuration.GenericConfiguration {
	conf := configuration.NewGenericConfiguration()
	conf.Host.Transport.Address = node.Host
	conf.Host.Transport.Protocol = "TCP"

	conf.APIRunner.Address = fmt.Sprintf(defaultHost+":191%02d", nodeIndex)
	conf.APIRunner.SwaggerPath = "application/api/spec/api-exported.yaml"

	conf.AvailabilityChecker.Enabled = true
	conf.AvailabilityChecker.KeeperURL = "http://127.0.0.1:12012/check"

	conf.AdminAPIRunner.Address = fmt.Sprintf(defaultHost+":190%02d", nodeIndex)
	conf.AdminAPIRunner.SwaggerPath = "application/api/spec/api-exported.yaml"

	conf.Metrics.ListenAddress = fmt.Sprintf(defaultHost+":80%02d", nodeIndex)
	conf.Introspection.Addr = fmt.Sprintf(defaultHost+":555%02d", nodeIndex)

	conf.Tracer.Jaeger.AgentEndpoint = defaultJaegerEndPoint
	conf.Log.Level = debugLevel
	conf.Log.Adapter = "zerolog"
	conf.Log.Formatter = "json"

	conf.KeysPath = bootstrapConf.DiscoveryKeysDir + fmt.Sprintf(bootstrapConf.KeysNameFormat, nodeIndex)
	conf.CertificatePath = fmt.Sprintf(discoveryCertificatePathTemplate, nodeIndex)
	conf.LightChainLimit = 15
	return conf
}

func createGenericConfigJoiner(node bootstrap.Node, nodeIndex int, bootstrapConf *bootstrap.Config) configuration.GenericConfiguration {
	genericConfiguration := configuration.NewGenericConfiguration()
	genericConfiguration.Host.Transport.Address = node.Host
	genericConfiguration.Host.Transport.Protocol = "TCP"
	genericConfiguration.APIRunner.Address = fmt.Sprintf(defaultHost+":191%02d", nodeIndex+len(bootstrapConf.DiscoveryNodes))
	genericConfiguration.AdminAPIRunner.Address = fmt.Sprintf(defaultHost+":190%02d", nodeIndex+len(bootstrapConf.DiscoveryNodes))
	genericConfiguration.Metrics.ListenAddress = fmt.Sprintf(defaultHost+":80%02d", nodeIndex+len(bootstrapConf.DiscoveryNodes))
	genericConfiguration.Introspection.Addr = fmt.Sprintf(defaultHost+":555%02d", nodeIndex+len(bootstrapConf.DiscoveryNodes))

	genericConfiguration.Tracer.Jaeger.AgentEndpoint = defaultJaegerEndPoint
	genericConfiguration.Log.Level = debugLevel
	genericConfiguration.Log.Adapter = "zerolog"
	genericConfiguration.Log.Formatter = "json"

	genericConfiguration.KeysPath = node.KeysFile
	genericConfiguration.CertificatePath = fmt.Sprintf(nodeCertificatePathTemplate, nodeIndex)
	genericConfiguration.LightChainLimit = 15
	return genericConfiguration
}

func createLightConfig() configuration.ConfigLight {
	conf := configuration.NewConfigurationLight()

	conf.Ledger.JetSplit.ThresholdRecordsCount = 1
	conf.Ledger.JetSplit.ThresholdOverflowCount = 0
	conf.Ledger.JetSplit.DepthLimit = 4
	return conf
}

func createVirtualConfig(index, nodeIndex int) configuration.ConfigVirtual {
	conf := configuration.NewConfigurationVirtual()

	rpcListenPort := 33300 + (index+nodeIndex)*nodeIndex
	conf.LogicRunner = configuration.NewLogicRunner()
	conf.LogicRunner.GoPlugin.RunnerListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort-1)
	conf.LogicRunner.RPCListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort)
	return conf
}

type commonConfigVars struct {
	BaseDir string
}

func writeGenesisConfig() {
	templates, err := template.ParseFiles(bootstrapConfigTmpl)
	check("Can't parse template: "+bootstrapConfigTmpl, err)

	var b bytes.Buffer
	err = templates.Execute(&b, &commonConfigVars{BaseDir: baseDir()})
	check("Can't process template: "+bootstrapConfigTmpl, err)

	err = makeFile(bootstrapFileName, b.String())
	check("Can't makeFileWithDir: "+bootstrapFileName, err)
}

var defaultInsloardConf *configuration.GenericConfiguration

func newDefaultInsolardConfig() configuration.GenericConfiguration {
	if defaultInsloardConf == nil {
		c := configuration.NewGenericConfiguration()
		defaultInsloardConf = &c
	}
	return *defaultInsloardConf
}

type pulsarConfigVars struct {
	commonConfigVars
	BootstrapHosts []string
	DataDir        string
	AgentEndpoint  string
}

func writePulsarConfig(pcv *pulsarConfigVars) {
	templates, err := template.ParseFiles(pulsardConfigTmpl)
	check("Can't parse template: "+pulsardConfigTmpl, err)

	var b bytes.Buffer
	err = templates.Execute(&b, pcv)
	check("Can't process template: "+pulsardConfigTmpl, err)
	err = makeFile(pulsardFileName, b.String())
	check("Can't makeFileWithDir: "+pulsardFileName, err)
}

func writeKeeperdConfig() {
	templates, err := template.ParseFiles(keeperdConfigTmpl)
	check("Can't parse template: "+keeperdConfigTmpl, err)

	var b bytes.Buffer
	err = templates.Execute(&b, nil)
	check("Can't process template: "+keeperdConfigTmpl, err)
	err = makeFile(keeperdFileName, b.String())
	check("Can't makeFileWithDir: "+keeperdFileName, err)
}

type promConfigVars struct {
	Jobs map[string][]string
}

func (pcv *promConfigVars) addTarget(name string, address string) {
	jobs := pcv.Jobs
	addrPair := strings.SplitN(address, ":", 2)
	addr := "host.docker.internal:" + addrPair[1]
	jobs[name] = append(jobs[name], addr)
}

func writePromConfig(pcv *promConfigVars) {
	templates, err := template.ParseFiles(prometheusConfigTmpl)
	check("Can't parse template: "+prometheusConfigTmpl, err)

	var b bytes.Buffer
	err = templates.Execute(&b, pcv)
	check("Can't process template: "+prometheusConfigTmpl, err)

	err = makeFileWithDir(outputDir, prometheusFileName, b.String())
	check("Can't makeFileWithDir: "+prometheusFileName, err)
}

func makeFile(name string, text string) error {
	fmt.Println("generate_insolar_configs.go: write to file", name)
	return ioutil.WriteFile(name, []byte(text), 0644)
}

func createFileWithDir(file string, text string) error {
	mustMakeDir(filepath.Dir(file))
	return makeFile(file, text)
}

// makeFileWithDir dumps `text` into file named `name` into directory `dir`.
// Creates directory if needed as well as file
func makeFileWithDir(dir string, name string, text string) error {
	err := os.MkdirAll(dir, 0775)
	if err != nil {
		return err
	}
	file := filepath.Join(dir, name)
	return makeFile(file, text)
}

func mustMakeDir(dir string) {
	err := os.MkdirAll(dir, 0775)
	check("couldn't create directory "+dir, err)
	fmt.Println("generate_insolar_configs.go: creates dir", dir)
}

func withBaseDir(subpath string) string {
	return filepath.Join(baseDir(), subpath)
}

func check(msg string, err error) {
	if err == nil {
		return
	}

	logCfg := configuration.NewLog()
	logCfg.Formatter = "text"
	inslog, _ := log.NewGlobalLogger(logCfg)
	inslog.WithField("error", err).Fatal(msg)
}
