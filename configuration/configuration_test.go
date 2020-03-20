// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfiguration_Load_Default(t *testing.T) {
	holder := NewLightHolder("testdata/insolard-light.yaml")
	err := holder.Load()
	require.NoError(t, err)

	cfg := NewLightConfig()
	t.Log(ToString(cfg))
	require.Equal(t, cfg, *holder.Configuration)
}

func TestConfiguration_DoubleLoad(t *testing.T) {
	holder := NewLightHolder("testdata/insolard-light.yaml")
	err := holder.Load()
	require.NoError(t, err)

	holder2 := NewLightHolder("testdata/insolard-light.yaml")
	err = holder2.Load()
	require.NoError(t, err)
	require.Equal(t, holder2.Configuration, holder.Configuration)
}

func TestConfiguration_Load_Invalid(t *testing.T) {
	holder := NewLightHolder("testdata/invalid.yaml")
	err := holder.Load()
	require.Error(t, err)
}

func TestConfiguration_LoadEnv(t *testing.T) {
	holder := NewLightHolder("testdata/insolard-light.yaml")

	require.NoError(t, os.Setenv("INSOLAR_HOST_TRANSPORT_ADDRESS", "127.0.0.2:5555"))
	err := holder.Load()
	require.NoError(t, err)
	require.NoError(t, os.Unsetenv("INSOLAR_HOST_TRANSPORT_ADDRESS"))

	require.NoError(t, err)
	require.Equal(t, "127.0.0.2:5555", holder.Configuration.Host.Transport.Address)

	defaultCfg := NewLightConfig()
	require.Equal(t, "127.0.0.1:0", defaultCfg.Host.Transport.Address)
}

func TestConfiguration_Load_EmptyPath(t *testing.T) {
	holder := NewLightHolder("")
	err := holder.Load()
	require.Error(t, err)

	require.Panics(t, func() { NewLightHolder("").MustLoad() })
}

func TestConfiguration_Load_ENVOverridesEmpty(t *testing.T) {
	_ = os.Setenv("INSOLAR_HOST_TRANSPORT_ADDRESS", "127.0.0.2:5555")
	defer os.Unsetenv("INSOLAR_HOST_TRANSPORT_ADDRESS")
	holder := NewLightHolder("testdata/insolard-light-empty.yaml")
	err := holder.Load()
	require.NoError(t, err)

	require.Equal(t, "127.0.0.2:5555", holder.Configuration.Host.Transport.Address)
}

func TestNewVirtualHolder_LoadDefault(t *testing.T) {
	holder := NewVirtualHolder("testdata/insolard-virtual.yaml")
	err := holder.Load()
	require.NoError(t, err)

	cfg := NewVirtualConfig()
	require.Equal(t, cfg, *holder.Configuration)
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
	exitCode := m.Run()

	// restore back variables
	for varName, varValue := range variablesBackup {
		if err := os.Setenv(varName, varValue); err != nil {
			fmt.Printf("Failed to unset env variable '%s' with '%s': %s\n",
				varName, varValue, err.Error())
		}

	}
	os.Exit(exitCode)
}
