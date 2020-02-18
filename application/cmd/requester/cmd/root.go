//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package cmd

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/insolar/insolar/application/api/requester"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/log/logadapter"
	crypto "github.com/insolar/x-crypto"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const ApplicationShortDescription string = "The requester is a simple CLI for sending requests to Insolar Platform"

var (
	memberKeysPath       string
	apiURL               string
	inputRequestParams   string
	shouldPasteSeed      bool
	shouldPastePublicKey bool
	verbose              bool
	memberPrivateKey     crypto.PrivateKey
	request              *requester.ContractRequest
)

func parseInputParams(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVarP(&memberKeysPath, "memberkeys", "k", "", "Path to member key")
	flags.StringVarP(&inputRequestParams, "request", "r", "", "The request body or path to request params file")
	flags.BoolVarP(&shouldPasteSeed, "autocompleteseed", "s", true, "Should replace seed to correct value")
	flags.BoolVarP(&shouldPastePublicKey, "autocompletekey", "p", true, "Should replace publicKey to correct value")
	flags.BoolVarP(&verbose, "verbose", "v", false, "Print request information")
}

func verifyParams() {
	// verify that the member keys paramsFile is exist
	if !isFileExists(memberKeysPath) {
		log.Fatal("Member keys does not exists")
	}

	// try to read keys
	keys, err := secrets.ReadXCryptoKeysFile(memberKeysPath, false)
	if err != nil {
		log.Fatal("Cannot parse member keys. ", err)
	}
	memberPrivateKey = keys.Private

	if len(inputRequestParams) == 0 {
		log.Fatal("Request parameters cannot be empty")
	}
	if isFileExists(inputRequestParams) {
		fileContent, err := ioutil.ReadFile(inputRequestParams)
		if err != nil {
			log.Fatal("Cannot read request. ", err)
		}
		// save to inputRequestParams if we could read params file for unmarshalling
		inputRequestParams = string(fileContent)
	}

	err = json.Unmarshal([]byte(inputRequestParams), &request)
	if err != nil {
		log.Fatal("Cannot unmarshal request. ["+inputRequestParams+"]", err)
	}
}

func isUrl(str string) (bool, error) {
	parsedUrl, err := url.Parse(str)
	return err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != "", err
}

func isFileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func getContextWithLogger() context.Context {
	cfg := configuration.NewLog()
	ctx := context.Background()
	cfg.Formatter = "text"
	if verbose {
		cfg.Level = insolar.DebugLevel.String()
	} else {
		cfg.Level = insolar.WarnLevel.String()
	}

	defaultCfg := logadapter.DefaultLoggerSettings()
	defaultCfg.Instruments.CallerMode = insolar.NoCallerField
	defaultCfg.Instruments.MetricsMode = insolar.NoLogMetrics
	logger, _ := log.NewLogExt(cfg, defaultCfg, 0)
	ctx = inslogger.SetLogger(ctx, logger)
	log.SetGlobalLogger(logger)

	return ctx
}

func createUserConfig(privateKey crypto.PrivateKey) (*requester.UserConfigJSON, error) {
	privateKeyBytes, err := secrets.ExportPrivateKeyPEM(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to export private key")
	}
	privateKeyStr := string(privateKeyBytes)

	publicKey, err := secrets.ExportPublicKeyPEM(secrets.ExtractPublicKey(privateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract public key")
	}
	publicKeyStr := string(publicKey)

	return requester.CreateUserConfig("", privateKeyStr, publicKeyStr)
}

// requireUrlArg returns an error if there is not url args.
func requireUrlArg() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("The program required url as an argument")
		}
		return nil
	}
}

func getRequesterCommand() *cobra.Command {
	retCmd := &cobra.Command{
		Use:     "requester <insolar url>",
		Short:   ApplicationShortDescription,
		Args:    requireUrlArg(),
		Example: "./requester http://localhost:19101/api/rpc  -k /tmp/userkey  -r params.json  -v",
		Run: func(_ *cobra.Command, args []string) {

			// no need to check args size because of requireUrlArg
			apiURL = args[0]
			if len(apiURL) > 0 {
				ok, err := isUrl(apiURL)
				if !ok {
					log.Fatal("URL parameter is incorrect. ", err)
				}
			}

			verifyParams()
			ctx := getContextWithLogger()
			requester.SetVerbose(verbose)

			userConfig, e := createUserConfig(memberPrivateKey)
			if e != nil {
				log.Fatal(e)
			}
			if shouldPastePublicKey {
				request.Params.PublicKey = userConfig.PublicKey
			}

			var response []byte
			var err error
			if shouldPasteSeed {
				response, err = requester.Send(ctx, apiURL, userConfig, &request.Params)
			} else {
				response, err = requester.SendWithSeed(ctx, apiURL, userConfig, &request.Params, request.Params.Seed)
			}

			if err != nil {
				log.Fatal(err)
			}

			_, _ = os.Stdout.Write(response)
		},
	}

	return retCmd
}

func Execute() error {
	rootCmd := getRequesterCommand()
	parseInputParams(rootCmd)
	return rootCmd.Execute()
}
