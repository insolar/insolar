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

package logicrunner_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin"

	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
)

func TestGoPlugin(t *testing.T) {
	if err := testutil.Build(); err != nil {
		t.Fatal("Logic runner build failed, skip tests:", err.Error())
	}

	l, cleaner := ledgertestutil.TmpLedger(t, "")
	defer cleaner()

	lr, err := logicrunner.NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err, "Initialize runner")

	lr.ArtifactManager = l.GetManager()
	mr := testutil.NewTestMessageRouter(lr)
	assert.NoError(t, lr.Start(core.Components{
		"core.Ledger":        l,
		"core.MessageRouter": mr,
	}), "starting logicrunner")

	insiderStorage, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(insiderStorage) // nolint: errcheck

	gp, err := goplugin.NewGoPlugin(
		&configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     "./goplugin/ginsider-cli/ginsider-cli",
			RunnerCodePath: insiderStorage,
		},
		mr,
		l.GetManager(),
	)
	assert.NoError(t, err)
	defer gp.Stop()

	err = lr.RegisterExecutor(core.MachineTypeGoPlugin, gp)
	assert.NoError(t, err)

	t.Run("Hello", func(t *testing.T) {
		hello(t, l, gp)
	})
}

func hello(t *testing.T, l core.Ledger, gp *goplugin.GoPlugin) {

	var helloCode = `
package main

import (
	"fmt"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Hello struct {
	foundation.BaseContract
}

func New() *Hello {
	return &Hello{};
}

func (b *Hello) String() string {
	return fmt.Sprint("Hello, Go is there!")
}
	`

	cb := testutil.NewContractBuilder(l.GetManager(), testutil.ICC)
	err := cb.Build(map[string]string{"hello": helloCode})
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		*cb.Codes["hello"],
		testutil.CBORMarshal(t, &struct{}{}),
		"String",
		testutil.CBORMarshal(t, []interface{}{}),
	)
	if err != nil {
		panic(err)
	}
	resParsed := testutil.CBORUnMarshalToSlice(t, res)
	assert.Equal(t, "Hello, Go is there!", resParsed[0])
}
