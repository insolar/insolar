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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	verbose bool
)

func main() {
	var sendURL string
	addURLFlag := func(fs *pflag.FlagSet) {
		fs.StringVarP(&sendURL, "url", "u", defaultURL(), "API URL")
	}

	var rootCmd = &cobra.Command{
		Use:   "insolar",
		Short: "insolar is the command line client for Insolar Platform",
	}
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "be verbose (default false)")

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.GetFullVersion())
		},
	}
	rootCmd.AddCommand(versionCmd)

	var infoCmd = &cobra.Command{
		Use:   "get-info",
		Short: "info about root member, root domain",
		Run: func(cmd *cobra.Command, args []string) {
			getInfo(sendURL)
		},
	}
	addURLFlag(infoCmd.Flags())
	rootCmd.AddCommand(infoCmd)

	var logLevel string
	var createMemberCmd = &cobra.Command{
		Use:   "create-member",
		Short: "creates member with random keys pair",
		Args:  cobra.ExactArgs(1), // username
		Run: func(cmd *cobra.Command, args []string) {
			createMember(sendURL, args[0], logLevel)
		},
	}
	createMemberCmd.Flags().StringVarP(
		&logLevel, "log-level-server", "L", "", "log level passed on server via request")
	addURLFlag(createMemberCmd.Flags())
	rootCmd.AddCommand(createMemberCmd)

	var genKeysPairCmd = &cobra.Command{
		Use:   "gen-key-pair",
		Short: "generates public/private keys pair",
		Run: func(cmd *cobra.Command, args []string) {
			generateKeysPair()
		},
	}
	rootCmd.AddCommand(genKeysPairCmd)

	var genMigrationAddressesCmd = &cobra.Command{
		Use:   "gen-migration-addresses",
		Short: "generates fake migration addresses",
		Run: func(cmd *cobra.Command, args []string) {
			generateMigrationAddresses()
		},
	}
	rootCmd.AddCommand(genMigrationAddressesCmd)

	var rootKeysFile string

	var (
		paramsPath   string
		rootAsCaller bool
	)
	var sendRequestCmd = &cobra.Command{
		Use:   "send-request",
		Short: "sends request",
		Run: func(cmd *cobra.Command, args []string) {
			sendRequest(sendURL, rootKeysFile, paramsPath, rootAsCaller)
		},
	}
	addURLFlag(sendRequestCmd.Flags())
	sendRequestCmd.Flags().StringVarP(
		&rootKeysFile, "root-keys", "k", "config.json", "path to json with root key pair")
	sendRequestCmd.Flags().StringVarP(
		&paramsPath, "params", "p", "", "path to params file (default params.json)")
	sendRequestCmd.Flags().BoolVarP(
		&rootAsCaller, "root-caller", "r", false, "use root member as caller")
	rootCmd.AddCommand(sendRequestCmd)

	var (
		role      string
		reuseKeys bool
		keysFile  string
		certFile  string
	)
	var certgenCmd = &cobra.Command{
		Use:   "certgen",
		Short: "generates keys and cerificate by root config",
		Run: func(cmd *cobra.Command, args []string) {
			genCertificate(rootKeysFile, role, sendURL, keysFile, certFile, reuseKeys)
		},
	}
	addURLFlag(certgenCmd.Flags())
	certgenCmd.Flags().StringVarP(
		&rootKeysFile, "root-keys", "k", "", "Config that contains public/private keys of root member")
	certgenCmd.Flags().StringVarP(
		&role, "role", "r", "virtual", "The role of the new node")
	certgenCmd.Flags().BoolVarP(
		&reuseKeys, "reuse-keys", "", false, "Read keys from file instead og generating of new ones")
	certgenCmd.Flags().StringVarP(
		&keysFile,
		"node-keys",
		"",
		"keys.json",
		"The OUT/IN ( depends on 'reuse-keys' ) file for public/private keys of the node",
	)
	certgenCmd.Flags().StringVarP(
		&certFile, "node-cert", "c", "cert.json", "The OUT file the node certificate")
	rootCmd.AddCommand(certgenCmd)

	rootCmd.AddCommand(bootstrapCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func defaultURL() string {
	if u := os.Getenv("INSOLAR_API_URL"); u != "" {
		return u
	}
	return "http://localhost:19001/admin-api/rpc"
}

type mixedConfig struct {
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	Caller     string `json:"caller"`
}

func createMember(sendURL string, userName string, serverLogLevel string) {
	ks := platformpolicy.NewKeyProcessor()

	logLevelInsolar, err := insolar.ParseLevel(serverLogLevel)
	check("Failed to parse logging level", err)

	privKey, err := ks.GeneratePrivateKey()
	check("Problems with generating of private key:", err)

	privKeyStr, err := ks.ExportPrivateKeyPEM(privKey)
	check("Problems with serialization of private key:", err)

	pubKeyStr, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(privKey))
	check("Problems with serialization of public key:", err)

	cfg := mixedConfig{
		PrivateKey: string(privKeyStr),
		PublicKey:  string(pubKeyStr),
	}

	info, err := requester.Info(sendURL)
	check("Problems with obtaining info", err)

	ucfg, err := requester.CreateUserConfig(info.RootMember, cfg.PrivateKey, cfg.PublicKey)
	check("Problems with creating user config:", err)

	ctx := inslogger.ContextWithTrace(context.Background(), "insolarUtility")
	params := requester.Params{
		CallSite:   "member.create",
		CallParams: []interface{}{userName, cfg.PublicKey},
		PublicKey:  ucfg.PublicKey,
		LogLevel:   logLevelInsolar.String(),
	}

	r, err := requester.Send(ctx, sendURL, ucfg, &params)
	check("Problems with sending request", err)

	var rStruct struct {
		Result string `json:"result"`
	}
	err = json.Unmarshal(r, &rStruct)
	check("Problems with understanding result", err)

	cfg.Caller = rStruct.Result
	result, err := json.MarshalIndent(cfg, "", "    ")
	check("Problems with marshaling config:", err)

	mustWrite(os.Stdout, string(result))
}

func verboseInfo(msg string) {
	if verbose {
		fmt.Fprintln(os.Stderr, msg)
	}
}

func mustWrite(out io.Writer, data string) {
	_, err := out.Write([]byte(data))
	check("Can't write data to output", err)
}

func generateMigrationAddresses() {
	maLen := 20000
	ma := make([]string, maLen)

	for i := 0; i < maLen; i++ {
		ma[i] = "fake_ma_" + strconv.Itoa(i)
	}

	result, err := json.MarshalIndent(ma, "", "    ")
	check("Problems with marshaling migration addresses:", err)

	mustWrite(os.Stdout, string(result))
}

func generateKeysPair() {
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

	mustWrite(os.Stdout, string(result))
}

func sendRequest(sendURL string, rootKeysFile string, paramsPath string, rootAsCaller bool) {
	requester.SetVerbose(verbose)

	userCfg, err := requester.ReadUserConfigFromFile(rootKeysFile)
	check("[ sendRequest ]", err)

	if rootAsCaller || userCfg.Caller == "" {
		info, err := requester.Info(sendURL)
		check("[ sendRequest ]", err)
		userCfg.Caller = info.RootMember
	}

	pPath := paramsPath
	if len(pPath) == 0 {
		pPath = rootKeysFile
	}
	reqCfg, err := requester.ReadRequestParamsFromFile(pPath)
	check("[ sendRequest ]", err)

	verboseInfo(fmt.Sprintln("User Config: ", userCfg))
	verboseInfo(fmt.Sprintln("Requester Config: ", reqCfg))

	ctx := inslogger.ContextWithTrace(context.Background(), "insolarUtility")
	response, err := requester.Send(ctx, sendURL, userCfg, reqCfg)
	check("[ sendRequest ]", err)

	mustWrite(os.Stdout, string(response))
}

func getInfo(url string) {
	info, err := requester.Info(url)
	check("[ sendRequest ]", err)
	fmt.Printf("TraceID    : %s\n", info.TraceID)
	fmt.Printf("RootMember : %s\n", info.RootMember)
	fmt.Printf("NodeDomain : %s\n", info.NodeDomain)
	fmt.Printf("RootDomain : %s\n", info.RootDomain)
}

func check(msg string, err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, msg, err)
		os.Exit(1)
	}
}
