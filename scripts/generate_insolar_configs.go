/*
 *    Copyright 2019 Insolar
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
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/insolar/insolar/cmd/pulsewatcher/config"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/genesis"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

const (
	defaultOutputConfigNameTmpl = "insolar_%d.yaml"
	defaultHost                 = "127.0.0.1"
	defaultJaegerEndPoint       = defaultHost + ":6831"
	defaultLogLevel             = "Debug"
	defaultGenesisFile          = "genesis.yaml"
	dataDirectoryTemplate       = "scripts/insolard/nodes/%d/data"
	certificatePathTemplate     = "scripts/insolard/nodes/%d/cert.json"
	pulsewatcherFileName        = "pulsewatcher.yaml"
)

var (
	genesisFile     string
	outputDir       string
	debugLevel      string
	gorundPortsPath string
)

func parseInputParams() {
	var rootCmd = &cobra.Command{}

	rootCmd.Flags().StringVarP(&genesisFile, "genesis", "g", defaultGenesisFile, "input genesis file")
	rootCmd.Flags().StringVarP(&outputDir, "output", "o", "", "output directory ( required )")
	rootCmd.Flags().StringVarP(&debugLevel, "debuglevel", "d", defaultLogLevel, "debug level")
	rootCmd.Flags().StringVarP(&gorundPortsPath, "gorundports", "p", "", "path to insgorund ports ( required )")

	err := rootCmd.Execute()
	check("Wrong input params:", err)

	if outputDir == "" || gorundPortsPath == "" {
		err := rootCmd.Usage()
		check("[ parseInputParams ]", err)
	}
}

func writeGorundPorts(gorundPorts [][]string) {
	var portsData string
	for _, ports := range gorundPorts {
		portsData += ports[0] + " " + ports[1] + "\n"
	}
	err := genesis.WriteFile("./", gorundPortsPath, portsData)
	check("Can't WriteFile: "+gorundPortsPath, err)
}

func writeInsolarConfigs(insolarConfigs []configuration.Configuration) {
	for index, conf := range insolarConfigs {
		data, err := yaml.Marshal(conf)
		check("Can't Marshal insolard config", err)
		fileName := fmt.Sprintf(defaultOutputConfigNameTmpl, index+1)
		err = genesis.WriteFile(outputDir, fileName, string(data))
		check("Can't WriteFile: "+fileName, err)
	}
}

func main() {
	parseInputParams()

	genesisConf, err := genesis.ParseGenesisConfig(genesisFile)
	check("Can't read genesis config", err)

	pwConfig := pulsewatcher.Config{}
	insolarConfigs := make([]configuration.Configuration, 0, len(genesisConf.DiscoveryNodes))

	gorundPorts := [][]string{}

	for index, node := range genesisConf.DiscoveryNodes {
		nodeIndex := index + 1
		conf := configuration.NewConfiguration()

		conf.Host.Transport.Address = node.Host
		conf.Host.Transport.Protocol = "TCP"

		rpcListenPort := 33300 + (index+nodeIndex)*nodeIndex
		conf.LogicRunner = configuration.NewLogicRunner()
		conf.LogicRunner.GoPlugin.RunnerListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort-1)
		conf.LogicRunner.RPCListen = fmt.Sprintf(defaultHost+":%d", rpcListenPort)
		if node.Role == "virtual" {
			gorundPorts = append(gorundPorts, []string{strconv.Itoa(rpcListenPort - 1), strconv.Itoa(rpcListenPort)})
		}

		conf.APIRunner.Address = fmt.Sprintf(defaultHost+":191%02d", nodeIndex)
		conf.Metrics.ListenAddress = fmt.Sprintf(defaultHost+":80%02d", nodeIndex)

		conf.Tracer.Jaeger.AgentEndpoint = defaultJaegerEndPoint
		conf.Log.Level = debugLevel
		conf.Log.Adapter = "logrus"
		conf.KeysPath = node.KeysFile
		conf.Ledger.Storage.DataDirectory = fmt.Sprintf(dataDirectoryTemplate, nodeIndex)
		conf.CertificatePath = fmt.Sprintf(certificatePathTemplate, nodeIndex)

		insolarConfigs = append(insolarConfigs, conf)

		pwConfig.Nodes = append(pwConfig.Nodes, conf.APIRunner.Address)
	}

	writeInsolarConfigs(insolarConfigs)
	writeGorundPorts(gorundPorts)

	pwConfig.Interval = 100 * time.Millisecond
	pwConfig.Timeout = 1 * time.Second
	err = pulsewatcher.WriteConfig(outputDir+"/utils", pulsewatcherFileName, pwConfig)
	check("couldn't write pulsewatcher config file", err)
}
