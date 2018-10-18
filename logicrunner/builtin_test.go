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

	"github.com/insolar/insolar/pulsar/entropygenerator"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/logicrunner/builtin/helloworld"

	"github.com/insolar/insolar/ledger/ledgertestutils"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/insolar/insolar/pulsar"
	"github.com/insolar/insolar/testutils/testmessagebus"
)

func byteRecorRef(b byte) core.RecordRef {
	var ref core.RecordRef
	ref[core.RecordRefSize-1] = b
	return ref
}

func TestBareHelloworld(t *testing.T) {
	lr, err := NewLogicRunner(&configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})

	l, cleaner := ledgertestutils.TmpLedger(t, lr, "")
	defer cleaner()
	am := l.GetArtifactManager()
	assert.NoError(t, err, "Initialize runner")

	mb := testmessagebus.NewTestMessageBus()
	assert.NoError(t, lr.Start(core.Components{
		Ledger:     l,
		MessageBus: mb,
	}), "starting logicrunner")
	MessageBusTrivialBehavior(mb, lr)
	lr.OnPulse(*pulsar.NewPulse(configuration.NewPulsar().NumberDelta, 0, &entropygenerator.StandardEntropyGenerator{}))

	hw := helloworld.NewHelloWorld()

	domain := byteRecorRef(2)
	request := byteRecorRef(3)
	_, _, classRef, err := goplugintestutils.AMPublishCode(t, am, domain, request, core.MachineTypeBuiltin, []byte("helloworld"))
	assert.NoError(t, err)

	contract, err := am.RegisterRequest(&message.CallConstructor{ClassRef: byteRecorRef(4)})
	assert.NoError(t, err)

	_, err = am.ActivateObject(domain, *contract, *classRef, *am.GenesisRef(), goplugintestutils.CBORMarshal(t, hw))
	assert.NoError(t, err)
	assert.Equal(t, true, contract != nil, "contract created")

	// #1
	resp, err := lr.Execute(&message.CallMethod{
		ObjectRef: *contract,
		Method:    "Greet",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{"Vany"}),
	})
	assert.NoError(t, err, "contract call")

	d := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Data)
	r := goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}([]interface{}{"Hello Vany's world"}), r)
	assert.Equal(t, map[interface{}]interface{}(map[interface{}]interface{}{"Greeted": uint64(1)}), d)

	// #2
	resp, err = lr.Execute(&message.CallMethod{
		ObjectRef: *contract,
		Method:    "Greet",
		Arguments: goplugintestutils.CBORMarshal(t, []interface{}{"Ruz"}),
	})
	assert.NoError(t, err, "contract call")

	d = goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Data)
	r = goplugintestutils.CBORUnMarshal(t, resp.(*reply.CallMethod).Result)
	assert.Equal(t, []interface{}([]interface{}{"Hello Ruz's world"}), r)
	assert.Equal(t, map[interface{}]interface{}(map[interface{}]interface{}{"Greeted": uint64(2)}), d)
}
