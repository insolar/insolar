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

package builtin

import (
	"testing"

	"io/ioutil"
	"os"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/logicrunner/builtin/helloworld"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
)

func TNewBuiltIn(t *testing.T, dir string) *BuiltIn {
	l, err := ledger.NewLedger(configuration.Ledger{
		DataDirectory: dir,
	})
	assert.NoError(t, err)
	am := l.GetManager()
	assert.Equal(t, true, am != nil)

	bi := NewBuiltIn(nil, am)

	assert.Equal(t, true, bi != nil)

	return bi
}

func byteRecorRef(b byte) core.RecordRef {
	var ref core.RecordRef
	ref[core.RecordRefSize-1] = b
	return ref
}

func TestBareHelloworld(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir) // nolint: errcheck

	bi := TNewBuiltIn(t, tmpDir)
	am := bi.AM
	hw := helloworld.NewHelloWorld()

	am.SetArchPref([]core.MachineType{0})

	domain := byteRecorRef(2)
	request := byteRecorRef(3)
	hwtype, err := am.DeclareType(domain, request, []byte{})
	assert.NoError(t, err, "creating type on ledger")
	coderef, err := am.DeployCode(domain, request, []core.RecordRef{*hwtype}, map[core.MachineType][]byte{0: nil})
	assert.NoError(t, err, "create code on ledger")

	ch := new(codec.CborHandle)
	var data []byte
	err = codec.NewEncoderBytes(&data, ch).Encode(hw)
	assert.NoError(t, err, "serialise new helloworld")

	classref, err := am.ActivateClass(domain, request, *coderef, data)
	assert.NoError(t, err, "create template for contract data")

	contract, err := am.ActivateObj(request, domain, *classref, *am.RootRef(), data)
	assert.NoError(t, err, "create actual contract")

	assert.Equal(t, true, contract != nil)

	bi.registry[coderef.String()] = bi.registry[helloworld.CodeRef().String()]
	var args []byte
	err = codec.NewEncoderBytes(&args, ch).Encode([]interface{}{"Vany"})
	assert.NoError(t, err, "serialise args")
	state, res, err := bi.Exec(*coderef, data, "Greet", args)
	assert.NoError(t, err, "contract call")
	t.Logf("%+v  %+v", state, res)
}
