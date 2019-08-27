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
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/insolar/insolar/bootstrap"
	pulsewatcher "github.com/insolar/insolar/cmd/pulsewatcher/config"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar/defaults"
	"github.com/insolar/insolar/log"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
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

	insolardDefaultsConfig = "scripts/insolard/defaults/insolard.yaml"
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

func writeInsolardConfigs(dir string, insolardConfigs []configuration.Configuration) {
	fmt.Println("generate_insolar_configs.go: writeInsolardConfigs...")
	for index, conf := range insolardConfigs {
		data, err := yaml.Marshal(conf)
		check("Can't Marshal insolard config", err)

		fileName := fmt.Sprintf(defaultOutputConfigNameTmpl, index+1)
		fileName = filepath.Join(dir, fileName)
		err = createFileWithDir(fileName, string(data))
		check("failed to create insolard config: "+fileName, err)
	}
}

func main() {
	parseInputParams()

	mustMakeDir(outputDir)
	writeGenesisConfig()

	bootstrapConf, err := bootstrap.ParseConfig(bootstrapFileName)
	check("Can't read bootstrap config", err)

	pwConfig := pulsewatcher.Config{}
	discoveryNodesConfigs := make([]configuration.Configuration, 0, len(bootstrapConf.DiscoveryNodes))

	var gorundPorts [][]string

	promVars := &promConfigVars{
		Jobs: map[string][]string{},
	}

	// process discovery nodes
	for index, node := range bootstrapConf.DiscoveryNodes {
		nodeIndex := index + 1

		conf := newDefaultInsolardConfig()

		conf.Host.Transport.Address = node.Host
		conf.Host.Transport.Protocol = "TCP"

		rpcListenPort := 33300 + (index+nodeIndex)*nodeIndex
		conf.LogicRunner = configuration.NewLogicRunner()
		conf.LogicRunner.GoPlugin.RunnerListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort-1)
		conf.LogicRunner.RPCListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort)
		if node.Role == "virtual" {
			gorundPorts = append(gorundPorts, []string{strconv.Itoa(rpcListenPort - 1), strconv.Itoa(rpcListenPort)})
		}

		if node.Role == "light_material" {
			conf.Ledger.JetSplit.ThresholdRecordsCount = 1
			conf.Ledger.JetSplit.ThresholdOverflowCount = 0
			conf.Ledger.JetSplit.DepthLimit = 4
		}

		conf.APIRunner.Address = fmt.Sprintf(defaultHost+":191%02d", nodeIndex)
		conf.AdminAPIRunner.Address = fmt.Sprintf(defaultHost+":190%02d", nodeIndex)
		conf.Metrics.ListenAddress = fmt.Sprintf(defaultHost+":80%02d", nodeIndex)
		conf.Introspection.Addr = fmt.Sprintf(defaultHost+":555%02d", nodeIndex)

		conf.Tracer.Jaeger.AgentEndpoint = defaultJaegerEndPoint
		conf.Log.Level = debugLevel
		conf.Log.Adapter = "zerolog"
		conf.Log.Formatter = "json"

		conf.KeysPath = bootstrapConf.DiscoveryKeysDir + fmt.Sprintf(bootstrapConf.KeysNameFormat, nodeIndex)
		conf.Ledger.Storage.DataDirectory = fmt.Sprintf(discoveryDataDirectoryTemplate, nodeIndex)
		conf.CertificatePath = fmt.Sprintf(discoveryCertificatePathTemplate, nodeIndex)

		discoveryNodesConfigs = append(discoveryNodesConfigs, conf)

		promVars.addTarget(node.Role, conf)

		pwConfig.Nodes = append(pwConfig.Nodes, conf.AdminAPIRunner.Address)
	}

	// process extra nodes
	nodeDataDirectoryTemplate = filepath.Join(outputDir, nodeDataDirectoryTemplate)
	nodeCertificatePathTemplate = filepath.Join(outputDir, nodeCertificatePathTemplate)

	nodesConfigs := make([]configuration.Configuration, 0, len(bootstrapConf.DiscoveryNodes))
	for index, node := range bootstrapConf.Nodes {
		nodeIndex := index + 1

		conf := newDefaultInsolardConfig()

		conf.Host.Transport.Address = node.Host
		conf.Host.Transport.Protocol = "TCP"

		rpcListenPort := 34300 + (index+nodeIndex+len(bootstrapConf.DiscoveryNodes)+1)*nodeIndex
		conf.LogicRunner = configuration.NewLogicRunner()
		conf.LogicRunner.GoPlugin.RunnerListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort-1)
		conf.LogicRunner.RPCListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort)
		if node.Role == "virtual" {
			gorundPorts = append(gorundPorts, []string{strconv.Itoa(rpcListenPort - 1), strconv.Itoa(rpcListenPort)})
		}

		conf.APIRunner.Address = fmt.Sprintf(defaultHost+":191%02d", nodeIndex+len(bootstrapConf.DiscoveryNodes))
		conf.AdminAPIRunner.Address = fmt.Sprintf(defaultHost+":190%02d", nodeIndex+len(bootstrapConf.DiscoveryNodes))
		conf.Metrics.ListenAddress = fmt.Sprintf(defaultHost+":80%02d", nodeIndex+len(bootstrapConf.DiscoveryNodes))
		conf.Introspection.Addr = fmt.Sprintf(defaultHost+":555%02d", nodeIndex+len(bootstrapConf.DiscoveryNodes))

		conf.Tracer.Jaeger.AgentEndpoint = defaultJaegerEndPoint
		conf.Log.Level = debugLevel
		conf.Log.Adapter = "zerolog"
		conf.Log.Formatter = "json"

		conf.KeysPath = node.KeysFile
		conf.Ledger.Storage.DataDirectory = fmt.Sprintf(nodeDataDirectoryTemplate, nodeIndex)
		conf.CertificatePath = fmt.Sprintf(nodeCertificatePathTemplate, nodeIndex)

		nodesConfigs = append(nodesConfigs, conf)

		promVars.addTarget(node.Role, conf)

		pwConfig.Nodes = append(pwConfig.Nodes, conf.AdminAPIRunner.Address)
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

	pwConfig.Interval = 500 * time.Millisecond
	pwConfig.Timeout = 1 * time.Second
	mustMakeDir(filepath.Dir(pulsewatcherFileName))
	err = pulsewatcher.WriteConfig(pulsewatcherFileName, pwConfig)
	check("couldn't write pulsewatcher config file", err)
	fmt.Println("generate_insolar_configs.go: write to file", pulsewatcherFileName)
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

var defaultInsloardConf *configuration.Configuration

func newDefaultInsolardConfig() configuration.Configuration {
	if defaultInsloardConf == nil {
		holder := configuration.NewHolderWithFilePaths(insolardDefaultsConfig).MustInit(true)
		defaultInsloardConf = &holder.Configuration
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

type promConfigVars struct {
	Jobs map[string][]string
}

func (pcv *promConfigVars) addTarget(name string, conf configuration.Configuration) {
	jobs := pcv.Jobs
	addrPair := strings.SplitN(conf.Metrics.ListenAddress, ":", 2)
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
	inslog, _ := log.NewLog(logCfg)
	inslog.WithField("error", err).Fatal(msg)
}
