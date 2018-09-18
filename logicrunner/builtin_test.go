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

package logicrunner

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/builtin/helloworld"
	"github.com/insolar/insolar/messagerouter/message"
	"github.com/insolar/insolar/messagerouter/response"

	"github.com/insolar/insolar/ledger/ledgertestutil"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
)

func byteRecorRef(b byte) core.RecordRef {
	var ref core.RecordRef
	ref[core.RecordRefSize-1] = b
	return ref
}

func TestBareHelloworld(t *testing.T) {
	l, cleaner := ledgertestutil.TmpLedger(t, "")
	defer cleaner()

	am := l.GetArtifactManager()
	lr, err := NewLogicRunner(configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err, "Initialize runner")

	mr := &testMessageRouter{lr}

	assert.NoError(t, lr.Start(core.Components{
		"core.Ledger":        l,
		"core.MessageRouter": mr,
	}), "starting logicrunner")

	hw := helloworld.NewHelloWorld()

	am.SetArchPref([]core.MachineType{core.MachineTypeBuiltin})

	domain := byteRecorRef(2)
	request := byteRecorRef(3)
	_, _, classRef, err := testutil.AMPublishCode(t, am, domain, request, core.MachineTypeBuiltin, []byte("helloworld"))

	contract, err := am.ActivateObj(request, domain, *classRef, *am.RootRef(), testutil.CBORMarshal(t, hw))
	assert.Equal(t, true, contract != nil, "contract created")

	// #1
	resp, err := lr.Execute(&message.CallMethodMessage{
		Request:   request,
		ObjectRef: *contract,
		Method:    "Greet",
		Arguments: testutil.CBORMarshal(t, []interface{}{"Vany"}),
	})
	assert.NoError(t, err, "contract call")

	d := testutil.CBORUnMarshal(t, resp.(*response.CommonResponse).Data)
	r := testutil.CBORUnMarshal(t, resp.(*response.CommonResponse).Result)
	assert.Equal(t, []interface{}([]interface{}{"Hello Vany's world"}), r)
	assert.Equal(t, map[interface{}]interface{}(map[interface{}]interface{}{"Greeted": uint64(1)}), d)

	// #2
	resp, err = lr.Execute(&message.CallMethodMessage{
		Request:   request,
		ObjectRef: *contract,
		Method:    "Greet",
		Arguments: testutil.CBORMarshal(t, []interface{}{"Ruz"}),
	})
	assert.NoError(t, err, "contract call")

	d = testutil.CBORUnMarshal(t, resp.(*response.CommonResponse).Data)
	r = testutil.CBORUnMarshal(t, resp.(*response.CommonResponse).Result)
	assert.Equal(t, []interface{}([]interface{}{"Hello Ruz's world"}), r)
	assert.Equal(t, map[interface{}]interface{}(map[interface{}]interface{}{"Greeted": uint64(2)}), d)
}
