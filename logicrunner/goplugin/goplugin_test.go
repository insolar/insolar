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
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/pkg/errors"
	"github.com/ugorji/go/codec"
)

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

	data, res, err := gp.Exec("secondary", data, "Echo", argsSerialized)
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

func compileBinaries() error {
	d, _ := os.Getwd()

	err := os.Chdir(d + "/ginsider")
	if err != nil {
		return errors.Wrap(err, "couldn't chdir")
	}

	defer os.Chdir(d) // nolint: errcheck

	err = exec.Command("go", "build", "ginsider.go").Run()
	if err != nil {
		return errors.Wrap(err, "can't build ginsider")
	}

	err = os.Chdir(d + "/testplugins")
	if err != nil {
		return errors.Wrap(err, "couldn't chdir")
	}

	err = exec.Command("make", "secondary.so").Run()
	if err != nil {
		return errors.Wrap(err, "can't build pluigins")
	}
	return nil
}

func TestHelloWorld(t *testing.T) {
	if err := compileBinaries(); err != nil {
		t.Fatal("Can't compile binaries", err)
	}
	dir, err := ioutil.TempDir("", "test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // nolint: errcheck

	gp, err := NewGoPlugin(
		Options{
			Listen:   "127.0.0.1:7778",
			CodePath: "./testplugins/",
		},
		RunnerOptions{
			Listen:          "127.0.0.1:7777",
			CodeStoragePath: dir,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	defer gp.Stop()

	hw := &HelloWorlder{77}
	res := hw.ProxyEcho(gp, "hi there here we are")
	if hw.Greeted != 78 {
		t.Fatalf("Got unexpected value: %d, 78 is expected", hw.Greeted)
	}

	if res != "hi there here we are" {
		t.Fatalf("Got unexpected value: %s, 'hi there here we are' is expected", res)
	}
}
