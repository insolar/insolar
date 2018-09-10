package logicrunner

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/logicrunner/builtin"
	"github.com/insolar/insolar/logicrunner/builtin/helloworld"
	"github.com/insolar/insolar/messagerouter/message"
	"github.com/stretchr/testify/assert"
	"github.com/ugorji/go/codec"
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
	}))

	hw := helloworld.NewHelloWorld()

	am.SetArchPref([]core.MachineType{0})

	domain := byteRecorRef(2)
	request := byteRecorRef(3)
	hwtype, err := am.DeclareType(domain, request, []byte{})
	assert.NoError(t, err, "creating type on ledger")
	coderef, err := am.DeployCode(
		domain, request, []core.RecordRef{*hwtype}, map[core.MachineType][]byte{core.MachineTypeBuiltin: nil},
	)
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

	bi := lr.Executors[core.MachineTypeBuiltin].(*builtin.BuiltIn)
	bi.Registry[coderef.String()] = bi.Registry[helloworld.CodeRef().String()]
	var args []byte
	err = codec.NewEncoderBytes(&args, ch).Encode([]interface{}{"Vany"})
	assert.NoError(t, err, "serialise args")

	resp := lr.Execute(&message.CallMethodMessage{
		Request:   request,
		ObjectRef: *contract,
		Method:    "Greet",
		Arguments: args,
	})
	assert.NoError(t, resp.Error, "contract call")
	t.Logf("%+v", resp)
}
