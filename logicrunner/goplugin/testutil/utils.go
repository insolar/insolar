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

package testutil

import (
	"go/build"
	"os"
)

// ChangeGoPath prepends `path` to GOPATH environment variable
// accounting for possibly for default value. Returns original
// value of the enviroment variable, don't forget to restore
// it with defer:
//    defer os.Setenv("GOPATH", origGoPath)
func ChangeGoPath(path string) (string, error) {
	gopathOrigEnv := os.Getenv("GOPATH")
	gopath := gopathOrigEnv
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	err := os.Setenv("GOPATH", path+":"+gopath)
	if err != nil {
		return "", err
	}
	return gopathOrigEnv, nil
}

// WriteFile dumps `text` into file named `name` into directory `dir`.
// Creates directory if needed as well as file
func WriteFile(dir string, name string, text string) error {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	fh, err := os.OpenFile(dir+"/"+name, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	_, err = fh.WriteString(text)
	if err != nil {
		return err
	}

	err = fh.Close()
	if err != nil {
		return err
	}

	return nil
}
