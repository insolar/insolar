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

package main

import (
	"github.com/insolar/insolar/version"
	"github.com/stretchr/testify/assert"
	"testing"

	upd "github.com/insolar/insolar/updater"
)

// Just to make Goland happy
func TestStub(t *testing.T) {
	updater := upd.NewUpdater()
	assert.NotNil(t, updater)
	assert.Equal(t, updater.CurrentVer, version.Version)
	assert.Equal(t, updater.BinariesList, []string{"insgocc", "insgorund", "insolar", "insolard", "pulsard", "insupdater"})
	assert.NotEqual(t, updater.ServersList, []string{""})
	assert.Equal(t, updater.LastSuccessServer, "")
	verifyAndUpdate(updater)
}
