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

package merkle

import (
	"testing"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/testutils/nodekeeper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type calculatorSuite struct {
	suite.Suite

	calculator Calculator
}

func TestCalculator(t *testing.T) {
	nk := nodekeeper.GetTestNodekeeper()
	l, clean := ledgertestutils.TmpLedger(t, "", core.Components{})
	c := &certificate.Certificate{}
	c.GenerateKeys()

	calculator := &calculator{}

	cm := component.Manager{}
	cm.Register(nk, l, c, calculator)

	assert.NotNil(t, calculator.Ledger)
	assert.NotNil(t, calculator.NodeNetwork)
	assert.NotNil(t, calculator.Certificate)

	s := &calculatorSuite{
		Suite:      suite.Suite{},
		calculator: calculator,
	}
	suite.Run(t, s)

	clean()
}
