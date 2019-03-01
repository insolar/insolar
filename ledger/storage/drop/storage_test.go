/*
*    Copyright 2019 Insolar Technologies
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

package drop

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewStorageMemory(t *testing.T) {
	ms := NewStorageMemory()

	require.NotNil(t, ms.jets)
}

func TestDropStorageMemory_ForPulse(t *testing.T) {
	// ms := NewStorageMemory()
	//
	// var drops []jet.Drop
	// f := fuzz.New().Funcs(func(jd *jet.Drop, c fuzz.Continue) {
	// 	jd.Pulse = 123
	// }).NumElements(5, 10)
	// f.Fuzz(&drops)
	//
	// for _, jd := range drops {
	// 	err := ms.Set(inslogger.TestContext(t), core.ZeroJetID, jd)
	// 	require.NoError(t, err)
	// }
}
