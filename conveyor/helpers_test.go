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

package conveyor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStateConvertingZeroState(t *testing.T) {
	// hex( int( 22*'1' + 10*'0', 2) )
	original := uint32(0xfffffc00)

	sm, state := extractStates(original)
	require.Equal(t, uint32(0x3fffff), sm)
	require.Equal(t, uint32(0), state)

	require.Equal(t, original, joinStates(sm, state))
}

func TestStateConvertingZeroSM(t *testing.T) {
	// hex( int( 22*'0' + 10*'1', 2) )

	original := uint32(0x3ff)

	sm, state := extractStates(original)
	require.Equal(t, uint32(0), sm)
	require.Equal(t, uint32(0x3ff), state)

	require.Equal(t, original, joinStates(sm, state))
}

func TestStateConvertingZeroAll(t *testing.T) {
	original := uint32(0)

	sm, state := extractStates(original)
	require.Equal(t, uint32(0), sm)
	require.Equal(t, uint32(0), state)

	require.Equal(t, original, joinStates(sm, state))
}

func TestStateConvertingAllOnes(t *testing.T) {
	//  hex(int(32 * '1', 2))
	original := uint32(0xffffffff)

	sm, state := extractStates(original)
	require.Equal(t, uint32(0x3fffff), sm)
	require.Equal(t, uint32(0x3ff), state)

	require.Equal(t, original, joinStates(sm, state))
}
