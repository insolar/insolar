// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package gateway

import (
	"bytes"
	"context"
	"crypto/rand"

	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/serialization"
	"github.com/insolar/insolar/network/storage"
)

func GetBootstrapPulse(ctx context.Context, accessor storage.PulseAccessor) insolar.Pulse {
	pulse, err := accessor.GetLatestPulse(ctx)
	if err != nil {
		pulse = *insolar.EphemeralPulse
	}

	return pulse
}

func EnsureGetPulse(ctx context.Context, accessor storage.PulseAccessor, pulseNumber insolar.PulseNumber) insolar.Pulse {
	pulse, err := accessor.GetPulse(ctx, pulseNumber)
	if err != nil {
		inslogger.FromContext(ctx).Panicf("Failed to fetch pulse: %d", pulseNumber)
	}

	return pulse
}

func getAnnounceSignature(
	node insolar.NetworkNode,
	isDiscovery bool,
	kp insolar.KeyProcessor,
	keystore insolar.KeyStore,
	scheme insolar.PlatformCryptographyScheme,
) ([]byte, *insolar.Signature, error) {

	brief := serialization.NodeBriefIntro{}
	brief.ShortID = node.ShortID()
	brief.SetPrimaryRole(adapters.StaticRoleToPrimaryRole(node.Role()))
	if isDiscovery {
		brief.SpecialRoles = member.SpecialRoleDiscovery
	}
	brief.StartPower = 10

	addr, err := endpoints.NewIPAddress(node.Address())
	if err != nil {
		return nil, nil, err
	}
	copy(brief.Endpoint[:], addr[:])

	pk, err := kp.ExportPublicKeyBinary(node.PublicKey())
	if err != nil {
		return nil, nil, err
	}

	copy(brief.NodePK[:], pk)

	buf := &bytes.Buffer{}
	err = brief.SerializeTo(nil, buf)
	if err != nil {
		return nil, nil, err
	}

	data := buf.Bytes()
	data = data[:len(data)-64]

	key, err := keystore.GetPrivateKey("")
	if err != nil {
		return nil, nil, err
	}

	digest := scheme.IntegrityHasher().Hash(data)
	sign, err := scheme.DigestSigner(key).Sign(digest)
	if err != nil {
		return nil, nil, err
	}

	return digest, sign, nil
}

func getKeyStore(cryptographyService insolar.CryptographyService) insolar.KeyStore {
	// TODO: hacked
	return cryptographyService.(*cryptography.NodeCryptographyService).KeyStore
}

type consensusProxy struct {
	Gatewayer network.Gatewayer
}

func (p consensusProxy) State() []byte {
	nshBytes := make([]byte, 64)
	_, _ = rand.Read(nshBytes)
	return nshBytes
}

func (p *consensusProxy) ChangePulse(ctx context.Context, newPulse insolar.Pulse) {
	p.Gatewayer.Gateway().(adapters.PulseChanger).ChangePulse(ctx, newPulse)
}

func (p *consensusProxy) UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte) {
	p.Gatewayer.Gateway().(adapters.StateUpdater).UpdateState(ctx, pulseNumber, nodes, cloudStateHash)
}
