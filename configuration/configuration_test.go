// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configuration

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfiguration_Load_Default(t *testing.T) {
	holder := NewHolder("testdata/default.yml")
	err := holder.Load()
	require.NoError(t, err)

	cfg := NewConfiguration()
	fmt.Println(ToString(cfg))
	require.Equal(t, cfg, holder.Configuration)
}

func TestConfiguration_DoubleLoad(t *testing.T) {
	holder := NewHolder("testdata/default.yml")
	err := holder.Load()
	require.NoError(t, err)

	holder2 := NewHolder("testdata/default.yml")
	err = holder2.Load()
	require.NoError(t, err)
	require.Equal(t, holder2.Configuration, holder.Configuration)
}

func TestConfiguration_GlobalLoad(t *testing.T) {
	holder := NewHolder("testdata/default.yml")
	err := holder.Load()
	require.NoError(t, err)

	require.NoError(t, os.Setenv("INSOLAR_LOG_LEVEL", "error"))
	holder2 := NewHolderGlobal("").MustLoad()
	require.NoError(t, os.Unsetenv("INSOLAR_LOG_LEVEL"))
	require.Equal(t, "error", holder2.Configuration.Log.Level)

}

func TestConfiguration_Load_Changed(t *testing.T) {
	holder := NewHolder("testdata/changed.yml")
	err := holder.Load()
	require.NoError(t, err)

	cfg := NewConfiguration()
	require.NotEqual(t, cfg, holder.Configuration)

	cfg.Log.Level = "Debug"
	require.Equal(t, cfg, holder.Configuration)
}

func TestConfiguration_Load_Invalid(t *testing.T) {
	holder := NewHolder("testdata/invalid.yml")
	err := holder.Load()
	require.Error(t, err)
}

func TestConfiguration_LoadEnv(t *testing.T) {
	holder := NewHolder("testdata/default.yml")

	require.NoError(t, os.Setenv("INSOLAR_HOST_TRANSPORT_ADDRESS", "127.0.0.2:5555"))
	err := holder.Load()
	require.NoError(t, err)
	require.NoError(t, os.Unsetenv("INSOLAR_HOST_TRANSPORT_ADDRESS"))

	require.NoError(t, err)
	require.Equal(t, "127.0.0.2:5555", holder.Configuration.Host.Transport.Address)

	defaultCfg := NewConfiguration()
	require.Equal(t, "127.0.0.1:0", defaultCfg.Host.Transport.Address)
}

func TestConfiguration_Load_EmptyPath(t *testing.T) {
	var (
		holder *Holder
		err    error
	)
	holder = NewHolder("")
	err = holder.Load()
	require.Error(t, err)

	require.Panics(t, func() { NewHolder("").MustLoad() })
}

func TestConfiguration_Load_ENVOverridesEmpty(t *testing.T) {
	var (
		holder *Holder
		err    error
	)
	holder = NewHolder("")
	err = holder.Load()
	require.Error(t, err)

	require.Panics(t, func() { NewHolder("").MustLoad() })
}

func TestMain(m *testing.M) {
	// backup and delete INSOLAR_ env variables, that may interfere with tests
	variablesBackup := make(map[string]string)
	for _, varPair := range os.Environ() {
		varPairSlice := strings.SplitN(varPair, "=", 2)
		varName, varValue := varPairSlice[0], varPairSlice[1]

		if strings.HasPrefix(varName, "INSOLAR_") {
			variablesBackup[varName] = varValue
			if err := os.Unsetenv(varName); err != nil {
				fmt.Printf("Failed to unset env variable '%s': %s\n",
					varName, err.Error())
			}
		}
	}

	// run tests
	m.Run()

	// restore back variables
	for varName, varValue := range variablesBackup {
		if err := os.Setenv(varName, varValue); err != nil {
			fmt.Printf("Failed to unset env variable '%s' with '%s': %s\n",
				varName, varValue, err.Error())
		}

	}
}
