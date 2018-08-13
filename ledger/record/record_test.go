/*
 *    Copyright 2018 INS Ecosystem
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

package record

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestCodeRecord_GetCode(t *testing.T) {
	rec := CodeRecord{
		TargetedCode: map[ArchType][]byte{
			1: {1},
			2: {2},
		},
	}

	code, err := rec.GetCode([]ArchType{15})
	assert.Error(t, err)

	code, err = rec.GetCode([]ArchType{3, 2, 1})
	assert.NoError(t, err)
	assert.Equal(t, []byte{2}, code)

	code, err = rec.GetCode([]ArchType{1})
	assert.NoError(t, err)
	assert.Equal(t, []byte{1}, code)
}

func TestPulseNumID(t *testing.T) {
	pulse0 := PulseNum(0)
	pulse1 := PulseNum(1)

	rec := &LockUnlockRequest{}
	idPulse0 := pulse0.ID(rec)
	idPulse1 := pulse1.ID(rec)
	assert.NotEqual(t, idPulse0, idPulse1)
}

func TestReference2Key(t *testing.T) {
	pulse0 := PulseNum(0)
	pulse1 := PulseNum(1)

	rec0 := &LockUnlockRequest{}
	rec1 := &LockUnlockRequest{}

	idPulse0 := pulse0.ID(rec0)
	idPulse1 := pulse1.ID(rec1)

	refPulse0 := &Reference{
		Record: idPulse0,
	}
	refPulse1 := &Reference{
		Record: idPulse1,
	}

	k0 := refPulse0.Key()
	k1 := refPulse1.Key()
	assert.NotEqual(t, k0, k1)
}
