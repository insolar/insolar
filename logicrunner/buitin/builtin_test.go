package buitin

import (
	"testing"

	"io/ioutil"
	"os"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/ledger"
	"github.com/insolar/insolar/logicrunner/buitin/helloworld"
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

	bi := NewBuiltIn(am, nil)

	assert.Equal(t, true, bi != nil)

	return bi
}

func byteRecorRef(b byte) core.RecordRef {
	ret := make([]byte, 64)
	ret[63] = b
	return ret
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
	coderef, err := am.DeployCode(domain, request, []core.RecordRef{hwtype}, map[core.MachineType][]byte{0: nil})
	assert.NoError(t, err, "create code on ledger")

	ch := new(codec.CborHandle)
	var data []byte
	err = codec.NewEncoderBytes(&data, ch).Encode(hw)
	assert.NoError(t, err, "serialise new helloworld")

	classref, err := am.ActivateClass(domain, request, coderef, data)
	assert.NoError(t, err, "create template for contract data")

	contract, err := am.ActivateObj(request, domain, classref, data)
	assert.NoError(t, err, "create actual contract")

	assert.Equal(t, true, contract != nil)

	bi.registry[utils.RefString(coderef)] = bi.registry[utils.RefString(helloworld.CodeRef())]
	var args []byte
	err = codec.NewEncoderBytes(&args, ch).Encode([]interface{}{"Vany"})
	assert.NoError(t, err, "serialise args")
	state, res, err := bi.Exec(coderef, data, "Greet", args)
	assert.NoError(t, err, "contract call")
	t.Logf("%+v  %+v", state, res)

}
