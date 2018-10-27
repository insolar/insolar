/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the License);
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an AS IS BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package manager

import (
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

//
// type VersionTable struct {
// 	Feature []struct {
// 		StartVersion string `yaml:"startVersion"`
// 		Description  string `yaml:",description"`
// 	}
// }


//go:generate go run internal/generate/generatefoobar.go

func main() *VersionManager {
	tt:= template.Must(template.New("versiontable").Parse(tpl))
	dest := strings.ToLower(os.Args)+"_vt.go"
	file, _:=  os.Create(dest)
	vals:= map[string]string{
		"MyType": os.Args[1],
	}
	tt.Execute(file, vals)

	// versionTable := make(map[string]*Feature)
	versionTable := NewFeatureList()
	// f, _ := os.Open(path.Join("..","..","..","testdata","versiontable.yml"))
	// bytes,_ := ioutil.ReadAll(f)

	bytes :=
		"versiontable:\n"+
	"	insolar:\n"+
		"		startversion: v1.1.1\n"+
		"		description: Version manager for Insolar platform test\n"+
		"	insolar2:\n"+
		"		startversion: v1.1.1\n"+
		"		description: Version manager for Insolar platform test\n"


	baseVersion, _ := ParseVersion(string(bytes))


	vm := &VersionManager{
		versionTable,
		baseVersion,
		viper.New(),
	}
	vm.viper.SetDefault("versiontable", vm.VersionTable)
	vm.viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vm.viper.SetEnvPrefix("insolar")
	vm.viper.SetConfigType("yml")

	return vm
}