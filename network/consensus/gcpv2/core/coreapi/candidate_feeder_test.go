//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package coreapi

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"

	"github.com/stretchr/testify/require"
)

// TODO
/*func TestPickNextJoinCandidate(t *testing.T) {
	require.Equal(t, nil, (&SequentialCandidateFeeder{}).PickNextJoinCandidate())

	s := &SequentialCandidateFeeder{buf: make([]profiles.CandidateProfile, 1)}
	c := profiles.NewCandidateProfileMock(t)
	s.buf[0] = c
	require.Equal(t, c, s.PickNextJoinCandidate())
}*/

func TestRemoveJoinCandidate(t *testing.T) {
	require.False(t, (&SequentialCandidateFeeder{}).RemoveJoinCandidate(false, insolar.ShortNodeID(0)))

	s := &SequentialCandidateFeeder{buf: make([]profiles.CandidateProfile, 1)}
	c := profiles.NewCandidateProfileMock(t)

	s.buf[0] = c
	c.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return insolar.ShortNodeID(1) })
	require.False(t, s.RemoveJoinCandidate(false, insolar.ShortNodeID(2)))

	c.GetStaticNodeIDMock.Set(func() insolar.ShortNodeID { return insolar.ShortNodeID(1) })
	require.True(t, s.RemoveJoinCandidate(false, insolar.ShortNodeID(1)))

	require.Equal(t, []profiles.CandidateProfile(nil), s.buf)

	s.buf = make([]profiles.CandidateProfile, 2)
	s.buf[0] = c
	c2 := profiles.NewCandidateProfileMock(t)
	s.buf[1] = c2
	require.True(t, s.RemoveJoinCandidate(false, insolar.ShortNodeID(1)))

	require.Equal(t, 1, len(s.buf))

	require.True(t, len(s.buf) > 0 && s.buf[0] == c2)
}

func TestAddJoinCandidatePanicForNil(t *testing.T) {
	s := NewSequentialCandidateFeeder(0)
	require.NotNil(t, s)
	require.Panics(t, func() { s.AddJoinCandidate(nil) })
}

func TestAddJoinCandidate(t *testing.T) {
	s := NewSequentialCandidateFeeder(0)

	f1 := transport.NewFullIntroductionReaderMock(t)
	f2 := transport.NewFullIntroductionReaderMock(t)

	err := s.AddJoinCandidate(f1)
	assert.NoError(t, err)
	require.True(t, len(s.buf) == 1 && s.buf[0] == f1)

	// add second
	err = s.AddJoinCandidate(f2)
	assert.NoError(t, err)
	require.True(t, len(s.buf) == 2 && s.buf[1] == f2)
}

func TestAddJoinCandidateFullQueue(t *testing.T) {
	s := NewSequentialCandidateFeeder(1)

	f1 := transport.NewFullIntroductionReaderMock(t)
	f2 := transport.NewFullIntroductionReaderMock(t)

	err := s.AddJoinCandidate(f1)
	assert.NoError(t, err)
	require.True(t, len(s.buf) == 1 && s.buf[0] == f1)

	// add second
	err = s.AddJoinCandidate(f2)
	assert.Error(t, err)
}
