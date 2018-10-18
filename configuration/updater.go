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

// Updater holds configuration for updater publishing.
type Updater struct {
	BinariesList []string
	ServersList  []string
	Delay        int64
}

// NewUpdater creates new default configuration for updater publishing.
func NewUpdater() Updater {
	return Updater{
		BinariesList: []string{"insgocc", "insgorund", "insolar", "insolard", "pulsard", "updateserv"},
		ServersList:  []string{"http://localhost:2345"},
		Delay:        60,
	}
}
