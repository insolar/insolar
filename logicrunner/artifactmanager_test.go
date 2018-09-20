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
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/goplugin"

	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
)

type testGoPluginCtx struct {
	preprocessor string
	ledger       core.Ledger
	goplugin     *goplugin.GoPlugin
}

func TestGoPlugin(t *testing.T) {
	runnerbin, preprocessorbin, err := testutil.Build()
	if err != nil {
		t.Fatal("Logic runner build failed, skip tests:", err.Error())
	}

	l, cleaner := ledgertestutil.TmpLedger(t, "")
	defer cleaner()

	lr, err := logicrunner.NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err, "Initialize runner")

	lr.ArtifactManager = l.GetArtifactManager()
	eb := testutil.NewTestEventBus(lr)
	assert.NoError(t, lr.Start(core.Components{
		"core.Ledger":   l,
		"core.EventBus": eb,
	}), "starting logicrunner")

	insiderStorage, err := ioutil.TempDir("", "test-")
	assert.NoError(t, err)
	defer os.RemoveAll(insiderStorage) // nolint: errcheck

	// TODO: don't reuse ports here
	initgoplugin := func() (*goplugin.GoPlugin, func()) {
		gopluginconfig := &configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     runnerbin,
			RunnerCodePath: insiderStorage,
		}

		gp, err := goplugin.NewGoPlugin(gopluginconfig, eb, l.GetArtifactManager())
		assert.NoError(t, err)
		// defer gp.Stop()
		err = lr.RegisterExecutor(core.MachineTypeGoPlugin, gp)
		assert.NoError(t, err)

		return gp, func() {
			_ = gp.Stop()
		}
	}

	tctx := &testGoPluginCtx{
		preprocessor: preprocessorbin,
		ledger:       l,
	}

	t.Run("hello", func(t *testing.T) {
		gp, stop := initgoplugin()
		tctx.goplugin = gp
		defer stop()
		tctx.hello(t)
	})
	t.Run("callingContract", func(t *testing.T) {
		gp, stop := initgoplugin()
		tctx.goplugin = gp
		defer stop()
		tctx.callingContract(t)
	})
	t.Run("injectingDelegate", func(t *testing.T) {
		gp, stop := initgoplugin()
		tctx.goplugin = gp
		defer stop()
		tctx.injectingDelegate(t)
	})
}

func (tctx *testGoPluginCtx) hello(t *testing.T) {
	l := tctx.ledger
	gp := tctx.goplugin

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
	am := l.GetArtifactManager()
	am.SetArchPref([]core.MachineType{core.MachineTypeGoPlugin})
	cb := testutil.NewContractBuilder(am, tctx.preprocessor)
	defer cb.Clean()
	err := cb.Build(map[string]string{"hello": helloCode})
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		&core.LogicCallContext{},
		*cb.Codes["hello"],
		testutil.CBORMarshal(t, &struct{}{}),
		"String",
		testutil.CBORMarshal(t, []interface{}{}),
	)
	assert.NoError(t, err)
	resParsed := testutil.CBORUnMarshalToSlice(t, res)
	assert.Equal(t, "Hello, Go is there!", resParsed[0])
}

func templateContract(t *testing.T, l core.Ledger, name string, codetemplate string) string {
	tpl := template.Must(template.New(name).Parse(codetemplate))
	var tplbuf bytes.Buffer
	err := tpl.Execute(&tplbuf, struct{ RootRefStr string }{
		RootRefStr: l.GetArtifactManager().RootRef().String(),
	})
	assert.NoError(t, err, "contract one template should compile")
	// log.Println("contract", name, ":", tplbuf.String())
	return tplbuf.String()
}

func (tctx *testGoPluginCtx) callingContract(t *testing.T) {
	l := tctx.ledger
	gp := tctx.goplugin

	var contractOneCodeTpl = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "contract-proxy/two"
import "github.com/insolar/insolar/core"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello(s string) string {
	holder := two.New()
	friend := holder.AsChild(core.NewRefFromBase58("{{ .RootRefStr }}"))

	res := friend.Hello(s)

	return "Hi, " + s + "! Two said: " + res
}
`
	contractOneCode := templateContract(t, l, "one", contractOneCodeTpl)

	var contractTwoCode = `
package main

import (
	"fmt"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
	X int
}

func New() *Two {
	return &Two{X:322};
}

func (r *Two) Hello(s string) string {
	r.X *= 2
	return fmt.Sprintf("Hello you too, %s. %d times!", s, r.X)
}
`

	cb := testutil.NewContractBuilder(l.GetArtifactManager(), tctx.preprocessor)
	defer cb.Clean()
	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		&core.LogicCallContext{},
		*cb.Codes["one"],
		testutil.CBORMarshal(t, &struct{}{}),
		"Hello",
		testutil.CBORMarshal(t, []interface{}{"ins"}),
	)
	assert.NoError(t, err)

	resParsed := testutil.CBORUnMarshalToSlice(t, res)
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", resParsed[0])
}

func (tctx *testGoPluginCtx) injectingDelegate(t *testing.T) {
	l := tctx.ledger
	gp := tctx.goplugin

	var contractOneCodeTpl = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "contract-proxy/two"
import "github.com/insolar/insolar/core"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello(s string) string {
	holder := two.New()
	friend := holder.AsDelegate(core.NewRefFromBase58("{{ .RootRefStr }}"))

	res := friend.Hello(s)

	return "Hi, " + s + "! Two said: " + res
}
`
	contractOneCode := templateContract(t, l, "one", contractOneCodeTpl)

	var contractTwoCode = `
package main

import (
	"fmt"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Two struct {
	foundation.BaseContract
	X int
}

func New() *Two {
	return &Two{X:322};
}

func (r *Two) Hello(s string) string {
	r.X *= 2
	return fmt.Sprintf("Hello you too, %s. %d times!", s, r.X)
}
`

	cb := testutil.NewContractBuilder(l.GetArtifactManager(), tctx.preprocessor)
	defer cb.Clean()
	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		&core.LogicCallContext{},
		*cb.Codes["one"],
		testutil.CBORMarshal(t, &struct{}{}),
		"Hello",
		testutil.CBORMarshal(t, []interface{}{"ins"}),
	)
	assert.NoError(t, err)

	resParsed := testutil.CBORUnMarshalToSlice(t, res)
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", resParsed[0])
}
