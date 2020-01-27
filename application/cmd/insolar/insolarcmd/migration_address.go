//
// Copyright 2020 Insolar Technologies GmbH
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
//

package insolarcmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/insolar/insolar/application/api/sdk"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
)

const (
	extraAddressesPattern = "addresses-%v.json"
	shardsPerFileDefault  = 10
)

// GenerateMigrationAddresses writes to io.Writer json array of random hashes compatible with ethereum address format.
// Array size is count.
func GenerateMigrationAddresses(w io.Writer, count int) error {
	ma := make([]string, count)
	for i := 0; i < count; i++ {
		ma[i] = RandomMigrationAddressHex()
	}

	b, err := json.MarshalIndent(ma, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal json: %v", err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}
	return nil
}

// WritesShardedMigrationsAddressesToDir writes provided sharded addresses to separate files in outDir directory.
func WritesShardedMigrationsAddressesToDir(outDir string, addrsByShard [][]string) error {
	// write every shardsPerFileDefault shards to one file.
	addrsPerFile := map[string][]string{}
	for i, addrs := range addrsByShard {
		filename := fmt.Sprintf("addresses-%v.json", i/shardsPerFileDefault)
		addrsPerFile[filename] = append(addrsPerFile[filename], addrs...)
	}
	for filename, addrs := range addrsPerFile {
		fullpath := filepath.Join(outDir, filename)
		b, err := json.MarshalIndent(addrs, "", "    ")
		if err != nil {
			return fmt.Errorf("failed to marshal json: %v", err)
		}

		err = ioutil.WriteFile(fullpath, b, 0644)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
	}
	return nil
}

// NRandomMigrationAddressSplitByShard generates n random addresses split by shards.
func NRandomMigrationAddressesSplitByShard(n int, shards int) [][]string {
	addrsShards := make([][]string, shards)
	for i := 0; i < n; i++ {
		addr := RandomMigrationAddressHex()
		shardIdx := foundation.GetShardIndex(addr, shards)
		addrsShards[shardIdx] = append(addrsShards[shardIdx], addr)
	}
	return addrsShards
}

// RandomMigrationAddressHex returns random migration address as lowercase hex string.
func RandomMigrationAddressHex() string {
	return strings.ToLower("0x" + hex.EncodeToString(RandomMigrationAddress()))
}

// RandomMigrationAddress returns random migration address.
func RandomMigrationAddress() []byte {
	b, err := randomNBytes(20)
	if err != nil {
		panic(err)
	}
	return b
}

func randomNBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// AddMigrationAddresses adds additional migration addresses to platform via API.
func AddMigrationAddresses(adminUrls []string, publicUrls []string, memberKeysDirPath string, dir string) error {
	insSDK, err := sdk.NewSDK(adminUrls, publicUrls, memberKeysDirPath, sdk.DefaultOptions)
	if err != nil {
		return fmt.Errorf("SDK is not initialized: %v", err)
	}

	// method AddMigrationAddresses in contract use only 10 shards in one call
	for i := 0; ; i++ {
		filename := filepath.Clean(filepath.Join(dir, fmt.Sprintf(extraAddressesPattern, strconv.Itoa(i))))
		if !fileExists(filename) {
			break
		}

		rawConf, err := ioutil.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("error while reading file: %v", err)
		}

		var addresses []string
		err = json.Unmarshal(rawConf, &addresses)
		if err != nil {
			return fmt.Errorf("error while unmarshal content of file %s to list of addresses: %v", filename, err)
		}

		_, err = insSDK.AddMigrationAddresses(addresses)
		if err != nil {
			return fmt.Errorf("error while adding addresses from file %v: %v", filename, err)
		}
	}
	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
