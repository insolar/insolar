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

package compile

import (
	"os"
	"os/exec"

	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/pkg/errors"
)

func Compile(output string, name string) error {
	dstDir := output + "/plugins/"
	err := os.MkdirAll(dstDir, 0777)
	if err != nil {
		return err
	}

	origGoPath, err := testutil.ChangeGoPath(output)
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
