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

package configuration

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfiguration_Load_Default(t *testing.T) {
	holder := NewHolder()
	err := holder.LoadFromFile("testdata/default.yml")
	assert.NoError(t, err)

	cfg := NewConfiguration()
	assert.Equal(t, cfg, holder.Configuration)
}

func TestConfiguration_Load_Changed(t *testing.T) {
	holder := NewHolder()
	err := holder.LoadFromFile("testdata/changed.yml")
	assert.NoError(t, err)

	cfg := NewConfiguration()
	assert.NotEqual(t, cfg, holder.Configuration)

	cfg.Log.Level = "Debug"
	assert.Equal(t, cfg, holder.Configuration)
}

func TestConfiguration_Save_Default(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)

	holder := NewHolder()
	err = holder.SaveAs(dir + "insolar.yml")
	assert.NoError(t, err)

	holder2 := NewHolder()
	err = holder2.LoadFromFile(dir + "insolar.yml")
	assert.NoError(t, err)

	assert.Equal(t, holder.Configuration, holder2.Configuration)
}

func TestConfiguration_Load_Invalid(t *testing.T) {
	holder := NewHolder()
	err := holder.LoadFromFile("testdata/invalid.yml")
	assert.Error(t, err)
}
