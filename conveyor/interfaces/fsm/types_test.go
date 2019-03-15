/*
 *    Copyright 2019 INS Ecosystem
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

package fsm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStateConvertingZeroState(t *testing.T) {
	// hex( int( 22*'1' + 10*'0', 2) )
	original := ElementState(0xfffffc00)

	sm, state := original.Parse()
	require.Equal(t, ID(0x3fffff), sm)
	require.Equal(t, StateID(0), state)

	require.Equal(t, original, NewElementState(sm, state))
}

func TestStateConvertingZeroSM(t *testing.T) {
	// hex( int( 22*'0' + 10*'1', 2) )

	original := ElementState(0x3ff)

	sm, state := original.Parse()
	require.Equal(t, ID(0), sm)
	require.Equal(t, StateID(0x3ff), state)

	require.Equal(t, original, NewElementState(sm, state))
}

func TestStateConvertingZeroAll(t *testing.T) {
	original := ElementState(0)

	sm, state := original.Parse()
	require.Equal(t, ID(0), sm)
	require.Equal(t, StateID(0), state)

	require.Equal(t, original, NewElementState(sm, state))
}

func TestStateConvertingAllOnes(t *testing.T) {
	//  hex(int(32 * '1', 2))
	original := ElementState(0xffffffff)

	sm, state := original.Parse()
	require.Equal(t, ID(0x3fffff), sm)
	require.Equal(t, StateID(0x3ff), state)

	require.Equal(t, original, NewElementState(sm, state))
}

func TestJoinStatesStateOverFlow(t *testing.T) {
	require.PanicsWithValue(t, "Invalid state: 333333333", func() { NewElementState(1, 333333333) })
}
