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

package main

import (
	"io/ioutil"
	"testing"
)

func Test_generateForFile(t *testing.T) {
	mode = "wrapper"
	w, err := generateForFile("../testplugins/secondary/main.go")
	if err != nil {
		t.Fatal(err)
	}
	// io.Copy(os.Stdout, w)
	b, err := ioutil.ReadAll(w)
	if err != nil {
		t.Fatal("reading from generated code", err)
	}
	if len(b) == 0 {
		t.Fatal("generator returns zero length code")
	}
}
