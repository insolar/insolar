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

package privatekey

import (
	"crypto"
)

type keyLoader struct {
	parseFunc func(key []byte) (crypto.PrivateKey, error)
}

// TODO: deprecated, use PEM format
func readJSON(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrap(err, "[ read ] couldn't read keys from: "+path)
	}
	var keys map[string]string
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return nil, errors.Wrap(err, "[ read ] failed to parse json.")
	}

	key, ok := keys["private_key"]
	if !ok {
		return nil, errors.Errorf("[ read ] couldn't read keys from: %s", path)
	}

	return []byte(key), nil
}
