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

package merkle

import (
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulsar/pulsartestutils"
	"github.com/stretchr/testify/suite"
)

func (t *merkleHelperSuite) TestMerkleHelperPulseHash() {
	pulse := &core.Pulse{
		PulseNumber:     core.PulseNumber(1337),
		NextPulseNumber: core.PulseNumber(1347),
		Entropy:         pulsartestutils.MockEntropyGenerator{}.GenerateEntropy(),
	}

	expectedHash, _ := hex.DecodeString(
		"bd18c009950389026c5c6f85c838b899d188ec0d667f77948aa72a49747c3ed31835b1bdbb8bd1d1de62846b5f308ae3eac5127c7d36d7d5464985004122cc90",
	)

	actualHash := t.mh.pulseHash(pulse)

	t.Assert().Equal(expectedHash, actualHash)
}

func (t *merkleHelperSuite) TestMerkleHelperNodeInfoHash() {
	pulseHash, _ := hex.DecodeString(
		"bd18c009950389026c5c6f85c838b899d188ec0d667f77948aa72a49747c3ed31835b1bdbb8bd1d1de62846b5f308ae3eac5127c7d36d7d5464985004122cc90",
	)

	stateHash, _ := hex.DecodeString(
		"11b2b580757fc04fccbbd88880c8f7e0cba5b1f87e61bcf5c18a20eca9bba83a8607119124bd54f6794f24104b9f16844ec26faa58782dec08829874edef8e75",
	)

	expectedHash, _ := hex.DecodeString(
		"a03ff8f1a5c845b0fb563b76866d5aa73aa81dfc25398382af3b087d8be061d6c686b1f890b547350057daa746d374ff92b7e03c42e2d0d5ae0aaa4e6b6d5ce8",
	)

	actualHash := t.mh.nodeInfoHash(pulseHash, stateHash)

	t.Assert().Equal(expectedHash, actualHash)
}

func (t *merkleHelperSuite) TestMerkleHelperNodeHash() {
	nodeInfoHash, _ := hex.DecodeString(
		"a03ff8f1a5c845b0fb563b76866d5aa73aa81dfc25398382af3b087d8be061d6c686b1f890b547350057daa746d374ff92b7e03c42e2d0d5ae0aaa4e6b6d5ce8",
	)

	signature, _ := hex.DecodeString("00000000")

	expectedHash, _ := hex.DecodeString(
		"09bbd5069e005d7a700838c54c07a828a90c0142087f00d0b779ad03645e7d7239a98f4cfa3f83c38646524b5c71bf5d0fb8a3ec303e4609aee9006ef7e27a9d",
	)

	actualHash := t.mh.nodeHash(signature, nodeInfoHash)

	t.Assert().Equal(expectedHash, actualHash)
}

func (t *merkleHelperSuite) TestMerkleHelperBucketEntryHash() {
	nodeHash, _ := hex.DecodeString(
		"09bbd5069e005d7a700838c54c07a828a90c0142087f00d0b779ad03645e7d7239a98f4cfa3f83c38646524b5c71bf5d0fb8a3ec303e4609aee9006ef7e27a9d",
	)

	entryIndex := 0

	expectedHash, _ := hex.DecodeString(
		"de4f82d2fa34ca006f7e7c644dd8a7d71662094b4057d1ecf03d32bf9f5bb3e9a3ab03d6f91e14f0c01c3ca00e67911402e5cc308dc1ff549e17e1d520b9c291",
	)

	actualHash := t.mh.bucketEntryHash(uint32(entryIndex), nodeHash)

	t.Assert().Equal(expectedHash, actualHash)
}

func (t *merkleHelperSuite) TestMerkleHelperBucketInfoHash() {
	nodeCount := 1337
	role := core.StaticRoleVirtual

	expectedHash, _ := hex.DecodeString(
		"eeb9dd175bb0d139083eadae8020f5b8623cb694263e8aec199c97213c383daf6ba0a58e734429b914cad1e401db1619526b1dabb57c5a020cd2fffed1f0cdeb",
	)

	actualHash := t.mh.bucketInfoHash(role, uint32(nodeCount))

	t.Assert().Equal(expectedHash, actualHash)
}

func (t *merkleHelperSuite) TestMerkleHelperBucketHash() {
	bucketInfoHash, _ := hex.DecodeString(
		"de4f82d2fa34ca006f7e7c644dd8a7d71662094b4057d1ecf03d32bf9f5bb3e9a3ab03d6f91e14f0c01c3ca00e67911402e5cc308dc1ff549e17e1d520b9c291",
	)

	bucketEntryHash, _ := hex.DecodeString(
		"eeb9dd175bb0d139083eadae8020f5b8623cb694263e8aec199c97213c383daf6ba0a58e734429b914cad1e401db1619526b1dabb57c5a020cd2fffed1f0cdeb",
	)

	expectedHash, _ := hex.DecodeString(
		"6fe03f36a4dd1599bcf671b81b33bd9dae3c6ddcbf8616fec26d9dff4ab56d31cbb272c6de14b590b010327f6c76a745b3b99df9512607a97c7508f183ec5183",
	)

	actualHash := t.mh.bucketHash(bucketInfoHash, bucketEntryHash)

	t.Assert().Equal(expectedHash, actualHash)
}

func (t *merkleHelperSuite) TestMerkleHelperGlobuleInfoHash() {
	prevCloudHash, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	globuleID := 1
	nodeCount := 1337

	expectedHash, _ := hex.DecodeString(
		"c6e8c3ad6446f235c308425825f130a1840c625618b0bb04682a683fe42a72bd03d5de800a5f0ff0080503e925efd18dd27242a35882acdb7feab7631dd33ba6",
	)

	actualHash := t.mh.globuleInfoHash(prevCloudHash, uint32(globuleID), uint32(nodeCount))

	t.Assert().Equal(expectedHash, actualHash)
}

func (t *merkleHelperSuite) TestMerkleHelperGlobuleHash() {
	globuleInfoHash, _ := hex.DecodeString(
		"c6e8c3ad6446f235c308425825f130a1840c625618b0bb04682a683fe42a72bd03d5de800a5f0ff0080503e925efd18dd27242a35882acdb7feab7631dd33ba6",
	)

	globuleNodeRoot, _ := hex.DecodeString(
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
	)

	expectedHash, _ := hex.DecodeString(
		"e1ee3ff9d599cd946629ac14e052e1537abaeadc633ef8295313651ab25c432686a56c45c3f3fb2a59ac2a47ed947d5337281fa949f81ceeb87231d52f1f8eae",
	)

	actualHash := t.mh.globuleHash(globuleInfoHash, globuleNodeRoot)

	t.Assert().Equal(expectedHash, actualHash)
}

type merkleHelperSuite struct {
	suite.Suite

	mh *merkleHelper
}

func TestMerkleHelper(t *testing.T) {
	mh := newMerkleHelper(platformpolicy.NewPlatformCryptographyScheme())

	s := &merkleHelperSuite{
		Suite: suite.Suite{},
		mh:    mh,
	}
	suite.Run(t, s)
}
