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

	initgoplugin := func() (*goplugin.GoPlugin, func()) {
		gopluginconfig := &configuration.GoPlugin{
			MainListen:     "127.0.0.1:7778",
			RunnerListen:   "127.0.0.1:7777",
			RunnerPath:     "./goplugin/ginsider-cli/ginsider-cli",
			RunnerCodePath: insiderStorage,
		}

		gp, err := goplugin.NewGoPlugin(gopluginconfig, mr, l.GetManager())
		assert.NoError(t, err)
		// defer gp.Stop()
		err = lr.RegisterExecutor(core.MachineTypeGoPlugin, gp)
		assert.NoError(t, err)

		return gp, func() {
			_ = gp.Stop()
		}
	}

	t.Run("Hello", func(t *testing.T) {
		gp, stop := initgoplugin()
		defer stop()
		hello(t, l, gp)
	})
	t.Run("callingContract", func(t *testing.T) {
		gp, stop := initgoplugin()
		defer stop()
		callingContract(t, l, gp)
	})
	t.Run("injectingDelegate", func(t *testing.T) {
		gp, stop := initgoplugin()
		defer stop()
		injectingDelegate(t, l, gp)
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

func templateContract(t *testing.T, l core.Ledger, name string, codetemplate string) string {
	tpl := template.Must(template.New(name).Parse(codetemplate))
	var tplbuf bytes.Buffer
	err := tpl.Execute(&tplbuf, struct{ RootRefStr string }{
		RootRefStr: l.GetManager().RootRef().String(),
	})
	assert.NoError(t, err, "contract one template should compile")
	// log.Println("contract", name, ":", tplbuf.String())
	return tplbuf.String()
}

func callingContract(t *testing.T, l core.Ledger, gp *goplugin.GoPlugin) {
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
	friend := holder.AsChild(core.String2Ref("{{ .RootRefStr }}"))

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

	cb := testutil.NewContractBuilder(l.GetManager(), testutil.ICC)
	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		*cb.Codes["one"],
		testutil.CBORMarshal(t, &struct{}{}),
		"Hello",
		testutil.CBORMarshal(t, []interface{}{"ins"}),
	)
	if err != nil {
		panic("gp.CallMethod: " + err.Error())
	}

	resParsed := testutil.CBORUnMarshalToSlice(t, res)
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", resParsed[0])
}

func injectingDelegate(t *testing.T, l core.Ledger, gp *goplugin.GoPlugin) {
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
	friend := holder.AsDelegate(core.String2Ref("{{ .RootRefStr }}"))

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

	cb := testutil.NewContractBuilder(l.GetManager(), testutil.ICC)
	err := cb.Build(map[string]string{
		"one": contractOneCode,
		"two": contractTwoCode,
	})
	assert.NoError(t, err)

	_, res, err := gp.CallMethod(
		*cb.Codes["one"],
		testutil.CBORMarshal(t, &struct{}{}),
		"Hello",
		testutil.CBORMarshal(t, []interface{}{"ins"}),
	)
	if err != nil {
		panic(err)
	}

	resParsed := testutil.CBORUnMarshalToSlice(t, res)
	assert.Equal(t, "Hi, ins! Two said: Hello you too, ins. 644 times!", resParsed[0])
}
