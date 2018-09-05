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

package goplugin

import (
	"os"
	"os/exec"
	"testing"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/ugorji/go/codec"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/insolar/insolar/network/hostnetwork"
)

func TestTypeCompatibility(t *testing.T) {
	var _ core.MachineLogicExecutor = (*GoPlugin)(nil)
}

func init() {
	log.SetLevel(log.DebugLevel)
}

type HelloWorlder struct {
	Greeted int
}

func (r *HelloWorlder) ProxyEcho(gp *GoPlugin, s string) string {
	ch := new(codec.CborHandle)
	var data []byte
	err := codec.NewEncoderBytes(&data, ch).Encode(*r)
	if err != nil {
		panic(err)
	}

	var args [1]interface{}
	args[0] = s

	var argsSerialized []byte
	err = codec.NewEncoderBytes(&argsSerialized, ch).Encode(args)
	if err != nil {
		panic(err)
	}

	data, res, err := gp.CallMethod(core.String2Ref("secondary"), data, "Echo", argsSerialized)
	if err != nil {
		panic(err)
	}

	err = codec.NewDecoderBytes(data, ch).Decode(r)
	if err != nil {
		panic(err)
	}

	var resParsed []interface{}
	err = codec.NewDecoderBytes(res, ch).Decode(&resParsed)
	if err != nil {
		panic(err)
	}

	return resParsed[0].(string)
}

func buildCLI(name string) error {
	out, err := exec.Command("go", "build", "-o", "./"+name+"/"+name, "./"+name+"/").CombinedOutput()
	if err != nil {
		return errors.Wrapf(err, "can't build %s: %s", name, string(out))
	}
	return nil
}

func buildInciderCLI() error {
	return buildCLI("ginsider-cli")
}

func buildPreprocessor() error {
	return buildCLI("preprocessor")
}

func compileBinaries() error {
	err := buildInciderCLI()
	if err != nil {
		return errors.Wrap(err, "can't build ginsider")
	}

	d, _ := os.Getwd()

	err = os.Chdir(d + "/testplugins")
	if err != nil {
		return errors.Wrap(err, "couldn't chdir")
	}

	defer os.Chdir(d) // nolint: errcheck

	out, err := exec.Command("make", "secondary.so").CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't build pluigins: "+string(out))
	}
	return nil
}

// TODO: uncomment me after using artifact manager instead of disk write
//func TestHelloWorld(t *testing.T) {
//	if err := compileBinaries(); err != nil {
//		t.Fatal("Can't compile binaries", err)
//	}
//	dir, err := ioutil.TempDir("", "test-")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer os.RemoveAll(dir) // nolint: errcheck
//
//	gp, err := NewGoPlugin(
//		configuration.Goplugin{
//			MainListen:     "127.0.0.1:7778",
//			MainCodePath:   "./testplugins/",
//			RunnerListen:   "127.0.0.1:7777",
//			RunnerCodePath: dir,
//		},
//		nil,
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer gp.Stop()
//
//	hw := &HelloWorlder{77}
//	res := hw.ProxyEcho(gp, "hi there here we are")
//	if hw.Greeted != 78 {
//		t.Fatalf("Got unexpected value: %d, 78 is expected", hw.Greeted)
//	}
//
//	if res != "hi there here we are" {
//		t.Fatalf("Got unexpected value: %s, 'hi there here we are' is expected", res)
//	}
//}

const contractOneCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"
import "contract-proxy/two"

type One struct {
	foundation.BaseContract
}

func (r *One) Hello(s string) string {
	friend := two.GetObject("some")
	res := friend.Hello(s)

	return "Hi, " + s + "! Two said: " + res
}
`

const contractTwoCode = `
package main

import "github.com/insolar/insolar/logicrunner/goplugin/foundation"

type Two struct {
	foundation.BaseContract
}

func (r *Two) Hello(s string) string {
	return "Hello you too, " + s
}
`

func generateContractProxy(root string, name string) error {
	dstDir := root + "/src/contract-proxy/" + name

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return err
	}

	contractPath := root + "/src/contract/" + name + "/main.go"

	out, err := exec.Command("./preprocessor/preprocessor", "proxy", "-o", dstDir+"/main.go", "--code-reference", "testReference", contractPath).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't generate proxy: "+string(out))
	}
	return nil
}

func buildContractPlugin(root string, name string) error {
	dstDir := root + "/plugins/"

	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return err
	}

	origGoPath, err := testutil.ChangeGoPath(root)
	if err != nil {
		return err
	}
	defer os.Setenv("GOPATH", origGoPath) // nolint: errcheck

	//contractPath := root + "/src/contract/" + name + "/main.go"

	out, err := exec.Command("go", "build", "-buildmode=plugin", "-o", dstDir+"/"+name+".so", "contract/"+name).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't build contract: "+string(out))
	}
	return nil
}

func generateContractWrapper(root string, name string) error {
	contractPath := root + "/src/contract/" + name + "/main.go"
	wrapperPath := root + "/src/contract/" + name + "/main_wrapper.go"

	out, err := exec.Command("./preprocessor/preprocessor", "wrapper", "-o", wrapperPath, contractPath).CombinedOutput()
	if err != nil {
		return errors.Wrap(err, "can't generate wrapper for contract '"+name+"': "+string(out))
	}
	return nil
}

func buildContracts(root string, names ...string) error {
	for _, name := range names {
		err := generateContractProxy(root, name)
		if err != nil {
			return err
		}
		err = generateContractWrapper(root, name)
		if err != nil {
			return err
		}
	}

	for _, name := range names {
		err := buildContractPlugin(root, name)
		if err != nil {
			return err
		}
	}
	return nil
}

type testMessageRouter struct {
	plugin *GoPlugin
}

func (r *testMessageRouter) Route(ctx hostnetwork.Context, msg core.Message) (resp core.Response, err error) {
	ch := new(codec.CborHandle)

	var data []byte
	err = codec.NewEncoderBytes(&data, ch).Encode(
		&struct{}{},
	)
	if err != nil {
		return core.Response{}, err
	}
	resdata, reslist, err := r.plugin.CallMethod(core.String2Ref("two"), data, msg.Method, msg.Arguments)
	return core.Response{Data: resdata, Result: reslist, Error: err}, nil
}

// TODO: uncomment after artifact manager integration
//func TestContractCallingContract(t *testing.T) {
//	err := buildInciderCLI()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = buildPreprocessor()
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	cwd, err := os.Getwd()
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer os.Chdir(cwd) // nolint: errcheck
//
//	tmpDir, err := ioutil.TempDir("", "test-")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer os.RemoveAll(tmpDir) // nolint: errcheck
//
//	err = testutil.WriteFile(tmpDir+"/src/contract/one/", "main.go", contractOneCode)
//	if err != nil {
//		t.Fatal(err)
//	}
//	err = testutil.WriteFile(tmpDir+"/src/contract/two/", "main.go", contractTwoCode)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	err = buildContracts(tmpDir, "one", "two")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	insiderStorage := tmpDir + "/insider-storage/"
//
//	err = os.MkdirAll(insiderStorage, 0777)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	mr := &testMessageRouter{}
//
//	gp, err := NewGoPlugin(
//		configuration.Goplugin{
//			MainListen:     "127.0.0.1:7778",
//			MainCodePath:   tmpDir + "/plugins/",
//			RunnerListen:   "127.0.0.1:7777",
//			RunnerCodePath: insiderStorage,
//		},
//		mr,
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer gp.Stop()
//
//	mr.plugin = gp
//
//	ch := new(codec.CborHandle)
//	var data []byte
//	err = codec.NewEncoderBytes(&data, ch).Encode(
//		&struct{}{},
//	)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	var argsSerialized []byte
//	err = codec.NewEncoderBytes(&argsSerialized, ch).Encode(
//		[]interface{}{"ins"},
//	)
//	if err != nil {
//		panic(err)
//	}
//
//	_, res, err := gp.CallMethod(core.RecordRef("one"), data, "Hello", argsSerialized)
//	if err != nil {
//		panic(err)
//	}
//
//	var resParsed []interface{}
//	err = codec.NewDecoderBytes(res, ch).Decode(&resParsed)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if resParsed[0].(string) != "Hi, ins! Two said: Hello you too, ins" {
//		t.Fatal("unexpected result")
//	}
//}
