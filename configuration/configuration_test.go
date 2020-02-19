// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package configuration

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestConfiguration_Load_Default(t *testing.T) {
	holder := NewHolder()
	err := holder.LoadFromFile("testdata/default.yml")
	require.NoError(t, err)

	cfg := NewConfiguration()
	require.Equal(t, cfg, holder.Configuration)
}

func TestConfiguration_Load_Changed(t *testing.T) {
	holder := NewHolder()
	err := holder.LoadFromFile("testdata/changed.yml")
	require.NoError(t, err)

	cfg := NewConfiguration()
	require.NotEqual(t, cfg, holder.Configuration)

	cfg.Log.Level = "Debug"
	require.Equal(t, cfg, holder.Configuration)
}

func TestConfiguration_Save_Default(t *testing.T) {
	dir, err := ioutil.TempDir("", "")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	newConfigurationPath := path.Join(dir, "insolar.yml")

	holder := NewHolder()
	holder.Configuration.Host.Transport.FixedPublicAddress = "192.168.1.1"
	err = holder.SaveAs(newConfigurationPath)
	require.NoError(t, err)

	holder2 := NewHolder()
	err = holder2.LoadFromFile(newConfigurationPath)
	require.NoError(t, err)

	require.Nil(t, holder2.viper.Get("insolar"))
	require.Equal(t, holder.Configuration, holder2.Configuration)
}

func TestConfiguration_Load_Invalid(t *testing.T) {
	holder := NewHolder()
	err := holder.LoadFromFile("testdata/invalid.yml")
	require.Error(t, err)
}

func TestConfiguration_LoadEnv(t *testing.T) {
	holder := NewHolder()

	require.NoError(t, os.Setenv("INSOLAR_HOST_TRANSPORT_ADDRESS", "127.0.0.2:5555"))
	err := holder.LoadFromFile("testdata/default.yml")
	require.NoError(t, err)
	require.NoError(t, os.Unsetenv("INSOLAR_HOST_TRANSPORT_ADDRESS"))

	require.NoError(t, err)
	require.Equal(t, "127.0.0.2:5555", holder.Configuration.Host.Transport.Address)

	defaultCfg := NewConfiguration()
	require.Equal(t, "127.0.0.1:0", defaultCfg.Host.Transport.Address)
}

func TestConfiguration_Init(t *testing.T) {
	var (
		holder *Holder
		err    error
	)
	holder, err = NewHolder().Init(false)
	require.NoError(t, err)
	require.NotNil(t, holder)

	holder = NewHolder().MustInit(false)
	require.NotNil(t, holder)

	holder, err = NewHolder().Init(true)
	require.Error(t, err)
	require.IsType(t, viper.ConfigFileNotFoundError{}, err)
	require.Nil(t, holder)

	require.Panics(t, func() { NewHolder().MustInit(true) })
}

func TestConfiguration_WithFilePaths(t *testing.T) {
	okCases := [][]string{
		{"testdata/changed.yml"},
		{"testdata/default.yml"},
		{"testdata/changed.yml", "testdata/default.yml"},
		{"testdata/changed.yml", "testdata/changed.yml"},
		{"testdata/default.yml", "testdata/default.yml"},
		{"notexists", "notexists.yaml", "testdata/default.yml", "testdata/default.yml"},
	}
	for _, tCase := range okCases {
		holder, err := NewHolderWithFilePaths(tCase...).Init(true)
		require.NoErrorf(t, err, "case: %#v", tCase)
		require.NotNil(t, holder)
	}

	failCases := [][]string{
		{"testdata/invalid.yml"},
		{"nonexists"},
		{"notexists", "notexists"},
		{"testdata/default.yml", "testdata/invalid.yml"},
	}
	for _, tCase := range failCases {
		holder, err := NewHolderWithFilePaths(tCase...).Init(true)
		require.Errorf(t, err, "case: %#v", tCase)
		require.Nil(t, holder)
	}
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
