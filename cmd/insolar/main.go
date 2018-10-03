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
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	ecdsa_helper "github.com/insolar/insolar/cryptohelpers/ecdsa"
	"github.com/insolar/insolar/genesis/bootstrapcertificate"
	"github.com/insolar/insolar/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const defaultStdoutPath = "-"

func chooseOutput(path string) (io.Writer, error) {
	var res io.Writer
	if path == defaultStdoutPath {
		res = os.Stdout
	} else {
		var err error
		res, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't open file for writing")
		}
	}
	return res, nil
}

func exit(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(1)
}

func check(msg string, err error) {
	if err != nil {
		exit(msg, err)
	}
}

var (
	output             string
	cmd                string
	numberCertificates uint
)

func parseInputParams() {
	var rootCmd = &cobra.Command{Use: "insolar"}
	rootCmd.Flags().StringVarP(&cmd, "cmd", "c", "",
		"available commands: default_config | random_ref | version | gen_keys | gen_certificates")
	rootCmd.Flags().StringVarP(&output, "output", "o", defaultStdoutPath, "output file (use - for STDOUT)")
	rootCmd.Flags().UintVarP(&numberCertificates, "num_serts", "n", 3, "number of certificates")
	err := rootCmd.Execute()

	if len(cmd) == 0 {
		rootCmd.Usage()
		os.Exit(0)
	}

	if err != nil {
		exit("Wrong input params:", err)
	}
}

func writeToOutput(out io.Writer, data string) {
	_, err := out.Write([]byte(data))
	if err != nil {
		exit("Can't write data to output", err)
	}
}

func printDefaultConfig(out io.Writer) {
	cfgHolder := configuration.NewHolder()

	writeToOutput(out, configuration.ToString(cfgHolder.Configuration))
}

func randomRef(out io.Writer) {
	ref := core.RandomRef()

	writeToOutput(out, ref.String()+"\n")
}

func generateKeysPair(out io.Writer) {
	privKey, err := ecdsa_helper.GeneratePrivateKey()
	check("Problems with generating of private key:", err)

	privKeyStr, err := ecdsa_helper.ExportPrivateKey(privKey)
	check("Problems with serialization of private key:", err)

	pubKeyStr, err := ecdsa_helper.ExportPublicKey(&privKey.PublicKey)
	check("Problems with serialization of public key:", err)

	result := fmt.Sprintf("Public key:\n %s\n", pubKeyStr)
	result += fmt.Sprintf("Private key:\n %s", privKeyStr)

	writeToOutput(out, result)
}

func serializeToJSON(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "    ")
}

func makeKeysJSON(keys []*ecdsa.PrivateKey) ([]byte, error) {
	kk := []map[string]string{}
	for _, key := range keys {
		pubKey, err := ecdsa_helper.ExportPublicKey(&key.PublicKey)
		check("[ makeKeysJSON ]", err)

		privKey, err := ecdsa_helper.ExportPrivateKey(key)
		check("[ makeKeysJSON ]", err)

		kk = append(kk, map[string]string{"public_key": pubKey, "private_key": privKey})
	}

	return serializeToJSON(map[string]interface{}{"keys": kk})
}

type certRecords = bootstrapcertificate.CertRecords

func generateCertificates(out io.Writer) {

	records := make(map[core.RecordRef]*ecdsa.PrivateKey)
	var recordsBuf bytes.Buffer
	cRecords := certRecords{}
	keys := []*ecdsa.PrivateKey{}
	for i := uint(0); i < numberCertificates; i++ {
		ref := core.RandomRef()
		privKey, err := ecdsa_helper.GeneratePrivateKey()
		check("[ generateCertificates ]:", err)

		records[ref] = privKey
		pubKey, err := ecdsa_helper.ExportPublicKey(&privKey.PublicKey)
		check("[ generateCertificates ]:", err)

		cRecords = append(cRecords, bootstrapcertificate.Record{NodeRef: ref.String(), PublicKey: pubKey})
		keys = append(keys, privKey)

		recordsBuf.WriteString(ref.String() + " " + pubKey)
	}

	cert, err := bootstrapcertificate.NewCertificateFromFields(cRecords, keys)
	check("[ generateCertificates ]:", err)

	certStr, err := cert.Dump()
	check("[ generateCertificates ]:", err)

	writeToOutput(out, certStr+"\n")

	keysList, err := makeKeysJSON(keys)
	check("[ generateCertificates ]:", err)
	writeToOutput(out, string(keysList)+"\n")
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
	case "gen_certificates":
		generateCertificates(out)
	}
}
