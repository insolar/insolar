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

package logicrunner

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
)

func TestTypeCompatibility(t *testing.T) {
	var _ core.LogicRunner = (*LogicRunner)(nil)
}

type testExecutor struct {
	constructorResponses []*testResp
	methodResponses      []*testResp
}

func (r *testExecutor) Stop() error {
	return nil
}

type testResp struct {
	data []byte
	res  core.Arguments
	err  error
}

func newTestExecutor() *testExecutor {
	return &testExecutor{
		constructorResponses: make([]*testResp, 0),
		methodResponses:      make([]*testResp, 0),
	}
}

func (r *testExecutor) CallMethod(ref core.RecordRef, data []byte, method string, args core.Arguments) ([]byte, core.Arguments, error) {
	if len(r.methodResponses) < 1 {
		panic(errors.New("no expected 'CallMethod' calls"))
	}

	res := r.methodResponses[0]
	r.methodResponses = r.methodResponses[1:]
	return res.data, res.res, res.err
}

func (r *testExecutor) CallConstructor(ref core.RecordRef, name string, args core.Arguments) ([]byte, error) {
	if len(r.constructorResponses) < 1 {
		panic(errors.New("no expected 'CallConstructor' calls"))
	}

	res := r.constructorResponses[0]
	r.constructorResponses = r.constructorResponses[1:]
	return res.data, res.err
}

func TestBasics(t *testing.T) {
	lr, err := NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err)

	comps := core.Components{
		"core.Ledger":        &testLedger{am: testutil.NewTestArtifactManager()},
		"core.MessageRouter": &testMessageRouter{},
	}
	assert.NoError(t, lr.Start(comps))
	assert.IsType(t, &LogicRunner{}, lr)

	_, err = lr.GetExecutor(core.MachineTypeGoPlugin)
	assert.Error(t, err)

	te := newTestExecutor()

	err = lr.RegisterExecutor(core.MachineTypeGoPlugin, te)
	assert.NoError(t, err)

	te2, err := lr.GetExecutor(core.MachineTypeGoPlugin)
	assert.NoError(t, err)
	assert.Equal(t, te, te2)
}

type testLedger struct {
	am core.ArtifactManager
}

func (r *testLedger) Start(components core.Components) error { return nil }
func (r *testLedger) Stop() error                            { return nil }
func (r *testLedger) GetManager() core.ArtifactManager       { return &r.am }

type testMessageRouter struct{}

func (testMessageRouter) Start(components core.Components) error { return nil }
func (testMessageRouter) Stop() error                            { return nil }
func (testMessageRouter) Route(msg core.Message) (resp core.Response, err error) {
	panic("implement me")
}

func TestExecution(t *testing.T) {
	am := testutil.NewTestArtifactManager()
	ld := &testLedger{am: am}
	mr := &testMessageRouter{}
	lr, err := NewLogicRunner(configuration.LogicRunner{})
	assert.NoError(t, err)
	lr.Start(core.Components{
		"core.Ledger":        ld,
		"core.MessageRouter": mr,
	})

	codeRef := core.String2Ref("someCode")
	dataRef := core.String2Ref("someObject")
	am.Objects[dataRef] = &testutil.TestObjectDescriptor{
		Data: []byte("origData"),
		Code: &testutil.TestCodeDescriptor{ARef: &codeRef},
	}

	te := newTestExecutor()
	te.methodResponses = append(te.methodResponses, &testResp{data: []byte("data"), res: core.Arguments("res")})

	err = lr.RegisterExecutor(core.MachineTypeGoPlugin, te)
	assert.NoError(t, err)

	resp := lr.Execute(core.Message{Constructor: false, Reference: dataRef})
	assert.NoError(t, resp.Error)
	assert.Equal(t, []byte("data"), resp.Data)
	assert.Equal(t, []byte("res"), resp.Result)

	te.constructorResponses = append(te.constructorResponses, &testResp{data: []byte("data"), res: core.Arguments("res")})
	resp = lr.Execute(core.Message{Constructor: true})
	assert.NoError(t, resp.Error)
	assert.Equal(t, []byte("data"), resp.Data)
	assert.Equal(t, []byte(nil), resp.Result)
}
