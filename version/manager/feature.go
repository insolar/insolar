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

package manager

import (
	"errors"

	"github.com/blang/semver"
)

type Feature struct {
	Key          string
	StartVersion *semver.Version
	Description  string
}

func NewFeature(key string, startVersion string, description string) (*Feature, error) {
	if key == "" {
		return nil, errors.New("Key cannot be null")
	}
	if startVersion == "" {
		return nil, errors.New("Start version cannot be null")
	}
	version, err := ParseVersion(startVersion)
	if err != nil {
		return nil, err
	}
	return &Feature{
		Key:          key,
		StartVersion: version,
		Description:  description,
	}, nil
}
