/*
 *    Copyright 2018 Insolar
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
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/insolar/insolar/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const defaultStdoutPath = "-"
const defaultURL = "http://localhost:19191/api"

func genDefaultConfig(r interface{}) ([]byte, error) {
	t := reflect.TypeOf(r)
	res := map[string]interface{}{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		switch field.Type.Kind() {
		case reflect.String:
			res[tag] = ""
		case reflect.Slice:
			res[tag] = []int{}
		}
	}

	rawJSON, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		return nil, errors.Wrap(err, "[ genDefaultConfig ]")
	}

	return rawJSON, nil
}

func chooseOutput(path string) (io.Writer, error) {
	var res io.Writer
	if path == defaultStdoutPath {
		res = os.Stdout
	} else {
		var err error
		res, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't open file for writing")
		}
	}
	return res, nil
}

func check(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
		os.Exit(1)
	}
}

var (
	output             string
	cmd                string
	numberCertificates uint
	configPath         string
	paramsPath         string
	verbose            bool
	sendUrls           string
	rootAsCaller       bool
)

func parseInputParams() {
	var rootCmd = &cobra.Command{}
	rootCmd.Flags().StringVarP(&cmd, "cmd", "c", "",
		"available commands: default_config | random_ref | version | gen_keys | gen_certificate | send_request | gen_send_configs")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "be verbose (default false)")
	rootCmd.Flags().StringVarP(&output, "output", "o", defaultStdoutPath, "output file (use - for STDOUT)")
	rootCmd.Flags().StringVarP(&sendUrls, "url", "u", defaultURL, "api url")
	rootCmd.Flags().UintVarP(&numberCertificates, "num_certs", "n", 3, "number of certificates")
	rootCmd.Flags().StringVarP(&configPath, "config", "g", "config.json", "path to configuration file")
	rootCmd.Flags().StringVarP(&paramsPath, "params", "p", "", "path to params file (default params.json)")
	rootCmd.Flags().BoolVarP(&rootAsCaller, "root_as_caller", "r", false, "use root member as caller")
	err := rootCmd.Execute()
	check("Wrong input params:", err)

	if len(cmd) == 0 {
		err = rootCmd.Usage()
		check("[ parseInputParams ]", err)
		os.Exit(0)
	}

}

func verboseInfo(msg string) {
	if verbose {
		log.Infoln(msg)
	}
}

func writeToOutput(out io.Writer, data string) {
	_, err := out.Write([]byte(data))
	check("Can't write data to output", err)
}

func printDefaultConfig(out io.Writer) {
	cfgHolder := configuration.NewHolder()
	writeToOutput(out, configuration.ToString(cfgHolder.Configuration))
}

func randomRef(out io.Writer) {
	ref := testutils.RandomRef()

	writeToOutput(out, ref.String()+"\n")
}

func generateKeysPair(out io.Writer) {
	ks := platformpolicy.NewKeyProcessor()

	privKey, err := ks.GeneratePrivateKey()
	check("Problems with generating of private key:", err)

	privKeyStr, err := ks.ExportPrivateKeyPEM(privKey)
	check("Problems with serialization of private key:", err)

	pubKeyStr, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(privKey))
	check("Problems with serialization of public key:", err)

	result, err := json.MarshalIndent(map[string]interface{}{
		"private_key": string(privKeyStr),
		"public_key":  string(pubKeyStr),
	}, "", "    ")
	check("Problems with marshaling keys:", err)

	writeToOutput(out, string(result))
}

func generateCertificate(out io.Writer) {
	boundCryptographyService, err := cryptography.NewStorageBoundCryptographyService(configPath)
	check("[ generateCertificate ] failed to create cryptography service", err)

	keyProcessor := platformpolicy.NewKeyProcessor()

	publicKey, err := boundCryptographyService.GetPublicKey()
	check("[ generateCertificate ] failed to retrieve public key", err)

	cert, err := certificate.NewCertificatesWithKeys(publicKey, keyProcessor)
	check("[ generateCertificate ] Can't create certificate", err)

	data, err := cert.Dump()
	check("[ generateCertificate ] Can't dump certificate", err)
	writeToOutput(out, data)
}

func sendRequest(out io.Writer) {
	requester.SetVerbose(verbose)
	userCfg, err := requester.ReadUserConfigFromFile(configPath)
	check("[ sendRequest ]", err)
	if rootAsCaller {
		info, err := requester.Info(sendUrls)
		check("[ sendRequest ]", err)
		userCfg.Caller = info.RootMember
	}

	pPath := paramsPath
	if len(pPath) == 0 {
		pPath = configPath
	}
	reqCfg, err := requester.ReadRequestConfigFromFile(pPath)
	check("[ sendRequest ]", err)

	verboseInfo(fmt.Sprintln("User Config: ", userCfg))
	verboseInfo(fmt.Sprintln("Requester Config: ", reqCfg))

	ctx := inslogger.ContextWithTrace(context.Background(), "insolarUtility")
	response, err := requester.Send(ctx, sendUrls, userCfg, reqCfg)
	check("[ sendRequest ]", err)

	writeToOutput(out, string(response))
}

func genSendConfigs(out io.Writer) {
	reqConf, err := genDefaultConfig(requester.RequestConfigJSON{})
	check("[ genSendConfigs ]", err)

	userConf, err := genDefaultConfig(requester.UserConfigJSON{})
	check("[ genSendConfigs ]", err)

	writeToOutput(out, "Request config:\n")
	writeToOutput(out, string(reqConf))
	writeToOutput(out, "\n\n")

	writeToOutput(out, "User config:\n")
	writeToOutput(out, string(userConf)+"\n")
}

func main() {
	parseInputParams()
	out, err := chooseOutput(output)
	check("Problems with parsing input:", err)

	switch cmd {
	case "default_config":
		printDefaultConfig(out)
	case "random_ref":
		randomRef(out)
	case "version":
		fmt.Println(version.GetFullVersion())
	case "gen_keys":
		generateKeysPair(out)
	case "gen_certificate":
		generateCertificate(out)
	case "send_request":
		sendRequest(out)
	case "gen_send_configs":
		genSendConfigs(out)
	}
}
