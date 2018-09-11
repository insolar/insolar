package logicrunner

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/logicrunner/builtin/helloworld"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/insolar/insolar/messagerouter/message"
	"github.com/stretchr/testify/assert"
)

func TNewAmL(t *testing.T, dir string) (core.ArtifactManager, *ledger.Ledger) {
	l, err := ledger.NewLedger(configuration.Ledger{
		DataDirectory: dir,
	})
	assert.NoError(t, err)
	am := l.GetManager()
	assert.Equal(t, true, am != nil)
	return am, l
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

	am, l := TNewAmL(t, tmpDir)
	lr, err := NewLogicRunner(configuration.LogicRunner{
		BuiltIn: &configuration.BuiltIn{},
	})
	assert.NoError(t, err, "Initialize runner")

	assert.NoError(t, lr.Start(core.Components{
		"core.Ledger":        l,
		"core.MessageRouter": &testMessageRouter{},
	}), "starting logicrunner")

	hw := helloworld.NewHelloWorld()

	am.SetArchPref([]core.MachineType{core.MachineTypeBuiltin})

	domain := byteRecorRef(2)
	request := byteRecorRef(3)
	_, _, classRef, err := testutil.AMPublishCode(t, am, domain, request, core.MachineTypeBuiltin, []byte("helloworld"))

	contract, err := am.ActivateObj(request, domain, *classRef, *am.RootRef(), testutil.CBORMarshal(t, hw))
	assert.Equal(t, true, contract != nil, "contract created")

	// #1
	resp := lr.Execute(&message.CallMethodMessage{
		Request:   request,
		ObjectRef: *contract,
		Method:    "Greet",
		Arguments: testutil.CBORMarshal(t, []interface{}{"Vany"}),
	})
	assert.NoError(t, resp.Error, "contract call")

	d := testutil.CBORUnMarshal(t, resp.Data)
	r := testutil.CBORUnMarshal(t, resp.Result)
	assert.Equal(t, []interface{}([]interface{}{"Hello Vany's world"}), r)
	assert.Equal(t, map[interface{}]interface{}(map[interface{}]interface{}{"Greeted": uint64(1)}), d)

	// #2
	resp = lr.Execute(&message.CallMethodMessage{
		Request:   request,
		ObjectRef: *contract,
		Method:    "Greet",
		Arguments: testutil.CBORMarshal(t, []interface{}{"Ruz"}),
	})
	assert.NoError(t, resp.Error, "contract call")

	d = testutil.CBORUnMarshal(t, resp.Data)
	r = testutil.CBORUnMarshal(t, resp.Result)
	assert.Equal(t, []interface{}([]interface{}{"Hello Ruz's world"}), r)
	assert.Equal(t, map[interface{}]interface{}(map[interface{}]interface{}{"Greeted": uint64(2)}), d)
}
