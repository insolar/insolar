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

package preprocessor

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/testutils"
)

func contractPath(name string, contractsDir string) string {
	return filepath.Join(contractsDir, name, name+".go")
}

func MakeTestName(file string, contractType string) string {
	return fmt.Sprintf("Generate contract %s from '%s'", contractType, file)
}

type RealContractsSuite struct {
	suite.Suite

	icc           string
	contractNames []string
	contractsDir  string
}

func (s *RealContractsSuite) SetupSuite() {
	if err := log.SetLevel("debug"); err != nil {
		log.Error("Failed to set logLevel to debug: ", err.Error())
	}

	var err error
	if _, s.icc, err = goplugintestutils.Build(); err != nil {
		s.Fail("Logic runner build failed, skip tests: ", err.Error())
	}

	if s.contractNames, err = GetRealContractsNames(); err != nil {
		s.Fail("Failed to load contracts names: ", err.Error())
	}

	if s.contractsDir, err = GetRealApplicationDir("contract"); err != nil {
		s.Fail("Failed to find contracts dir: ", err.Error())
	}
}

func (s *RealContractsSuite) TestGenerateProxies() {
	for _, name := range s.contractNames {
		file := contractPath(name, s.contractsDir)
		testName := MakeTestName(file, "proxy")

		s.T().Run(testName, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			parsed, err := ParseFile(file, insolar.MachineTypeGoPlugin)
			a.NoError(err)

			var buf bytes.Buffer
			err = parsed.WriteProxy(testutils.RandomRef().String(), &buf)
			a.NoError(err)

			code, err := ioutil.ReadAll(&buf)
			a.NoError(err)
			a.NotEqual(0, len(code))
		})
	}
}

func (s *RealContractsSuite) TestGenerateWrappers() {
	for _, name := range s.contractNames {
		file := contractPath(name, s.contractsDir)
		testName := MakeTestName(file, "wrapper")

		s.T().Run(testName, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			parsed, err := ParseFile(file, insolar.MachineTypeGoPlugin)
			a.NoError(err)

			var buf bytes.Buffer
			err = parsed.WriteWrapper(&buf, parsed.ContractName())
			a.NoError(err)

			code, err := ioutil.ReadAll(&buf)
			a.NoError(err)
			a.NotEqual(0, len(code))
		})
	}
}

func (s *RealContractsSuite) TestCompiling() {
	contracts := make(map[string]string)
	for _, name := range s.contractNames {
		code, err := ioutil.ReadFile(contractPath(name, s.contractsDir))
		s.NoError(err)
		contracts[name] = string(code)
	}

	am := goplugintestutils.NewTestArtifactManager()
	cb := goplugintestutils.NewContractBuilder(am, s.icc)

	err := cb.Build(contracts)
	s.NoError(err)
	cb.Clean()
}

func TestRealSmartContract(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RealContractsSuite))
}
