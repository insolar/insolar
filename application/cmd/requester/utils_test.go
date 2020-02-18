//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/insolar/blob/master/network/LICENSE.md.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/insolar/insolar/insolar/secrets"
)

const (
	binLocation                string = "./../../../bin/requester"
	getSeedResponse            string = `{"jsonrpc":"2.0","result":{"seed":"s2pMpKZwIeqOaHbVkWQQFnzMPqotBKwEnBElZCTmz3w=","traceID":"c648d2ee1e414ce7f12f266cef47eeac"},"id":0}`
	memberCreateResponse       string = `{"jsonrpc":"2.0","result":{"callResult":{"reference":"insolar:1AhiuscmUtUs4gqtSQFcKj-PV5Z1hp3GVGG4SPtdxkvs"},"requestReference":"insolar:1AhiusW4AvZu_0tdLO0Vq_cH9Hoc62vU4W0vJApydA1I.record","traceID":"f441b82fb75cc744595524353140e0b2"},"id":8674665223082153551}`
	createMemberRequestExample string = `{
  "jsonrpc": "2.0",
  "method": "contract.call",
  "id": 1,
  "params": {
    "seed": "fhDEwRRbSnYnbMnALKMh8gXdzaSvRv/nfsGC9S7kqik=",
    "callSite": "member.create",
    "publicKey": "-----BEGIN PUBLIC KEY-----\\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEMSbA4KvO/jlwY+8WFDEdwhCLlsEC\\nF3/GYvu9iTWHwCctx1wTbGGjNLY03EjXyYxaf8coNbSbZeu+jXcWeMHG0A==\\n-----END PUBLIC KEY-----"
  }
}`
)

func handlers() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		localVarBody, err := ioutil.ReadAll(r.Body)
		_ = r.Body.Close()
		w.WriteHeader(http.StatusOK)
		if err == nil {
			body := strings.ToLower(string(localVarBody))
			if strings.Contains(body, "node.getseed") {
				if _, err := fmt.Fprint(w, getSeedResponse); err != nil {
					log.Println(err.Error())
				}
			} else if strings.Contains(body, "member.create") {
				if _, err := fmt.Fprint(w, memberCreateResponse); err != nil {
					log.Println(err.Error())
				}
			}
		} else {
			if _, err := fmt.Fprint(w, "{\"requestReference\" :\"sddasacs\"'}"); err != nil {
				log.Println(err.Error())
			}
		}
	})
	return r
}

func runCmd(args ...string) (string, error) {
	cmd := exec.Command(binLocation, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func getArgs(key, url, requestParams string, shouldPrivateKey, shouldSeed bool) []string {
	cmd := make([]string, 5)
	cmd[0] = fmt.Sprintf("%v", url)
	cmd[1] = fmt.Sprintf("-k=%v", key)
	cmd[2] = fmt.Sprintf("-r=%v", requestParams)
	cmd[3] = fmt.Sprintf("-p=%v", shouldPrivateKey)
	cmd[4] = fmt.Sprintf("-s=%v", shouldSeed)
	return cmd
}

func writePrivateKeyToFile(pathToFile string) {
	privKey, err := secrets.GeneratePrivateKeyEthereum()
	if err != nil {
		log.Fatal("Problems with generating of private key:", err)
	}

	privKeyStr, err := secrets.ExportPrivateKeyPEM(privKey)
	if err != nil {
		log.Fatal("Problems with serialization of private key:", err)
	}

	pubKeyStr, err := secrets.ExportPublicKeyPEM(secrets.ExtractPublicKey(privKey))
	if err != nil {
		log.Fatal("Problems with serialization of public key:", err)
	}

	result, err := json.MarshalIndent(map[string]interface{}{
		"private_key": string(privKeyStr),
		"public_key":  string(pubKeyStr),
	}, "", "    ")
	if err != nil {
		log.Fatal("Problems with marshaling keys:", err)
	}

	err = ioutil.WriteFile(pathToFile, result, 0644)
	if err != nil {
		log.Fatal("Cannot write to temp paramsFile", err)
	}
}

func getRequestParamsFile() *os.File {
	tempFile, err := ioutil.TempFile("", "requester-test-params-")
	if err != nil {
		log.Fatal("failed open tmp paramsFile:", err)
	}
	filePath := tempFile.Name()

	err = ioutil.WriteFile(filePath, []byte(createMemberRequestExample), 0644)
	if err != nil {
		log.Fatal("Cannot write to temp paramsFile", err)
	}
	return tempFile
}
