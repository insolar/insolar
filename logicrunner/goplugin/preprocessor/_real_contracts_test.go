/*
 *    Copyright 2018 INS Ecosystem
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
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

var contractNames []string = []string{"member", "wallet", "allowance"}
var pathWithContracts string = "../../../genesis/experiment/"

func GetContractsList() []string {
	var result []string
	for i := 0; i < len(contractNames); i++ {
		result = append(result, pathWithContracts+contractNames[i]+"/"+contractNames[i]+".go")
	}

	return result
}

func MakeTestName(file string, contractType string) string {
	return fmt.Sprintf("Generate contract %s from '%s'", contractType, file)
}

func TestGenerateProxiesForRealSmartContracts(t *testing.T) {
	for _, file := range GetContractsList() {
		t.Run(MakeTestName(file, "proxy"), func(t *testing.T) {
			var buf bytes.Buffer
			err := generateContractProxy(file, "testRef", &buf)
			assert.NoError(t, err)

			code, err := ioutil.ReadAll(&buf)
			assert.NoError(t, err)
			assert.NotEqual(t, len(code), 0)
		})
	}
}

func TestGenerateWrappersForRealSmartContracts(t *testing.T) {
	for _, file := range GetContractsList() {
		t.Run(MakeTestName(file, "wrapper"), func(t *testing.T) {
			var buf bytes.Buffer
			err := generateContractWrapper(file, &buf)
			assert.NoError(t, err)

			code, err := ioutil.ReadAll(&buf)
			assert.NoError(t, err)
			assert.NotEqual(t, len(code), 0)
		})
	}
}
