/*
 *    Copyright 2018 Insolar
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

	"github.com/spf13/viper"
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

func TestConfiguration_LoadEnv(t *testing.T) {
	holder := NewHolder()
	defaultCfg := NewConfiguration()

	os.Setenv("INSOLAR_HOST_TRANSPORT_ADDRESS", "127.0.0.2:5555")
	err := holder.LoadEnv()
	assert.NoError(t, err)
	assert.NotEqual(t, defaultCfg, holder.Configuration)
	assert.Equal(t, "127.0.0.2:5555", holder.Configuration.Host.Transport.Address)
}

func TestConfiguration_Init(t *testing.T) {
	var (
		holder *Holder
		err    error
	)
	holder, err = NewHolder().Init(false)
	assert.NoError(t, err)
	assert.NotNil(t, holder)

	holder = NewHolder().MustInit(false)
	assert.NotNil(t, holder)

	holder, err = NewHolder().Init(true)
	assert.Error(t, err)
	assert.IsType(t, viper.ConfigFileNotFoundError{}, err)
	assert.Nil(t, holder)

	assert.Panics(t, func() { NewHolder().MustInit(true) })
}
