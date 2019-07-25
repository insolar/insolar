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

package tests

import (
	"bytes"
	"context"
	"crypto"
	"crypto/ecdsa"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/keystore"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensus"
	"github.com/insolar/insolar/network/consensus/adapters"
	"github.com/insolar/insolar/network/consensus/common/endpoints"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/profiles"
	"github.com/insolar/insolar/network/consensus/serialization"
	"github.com/insolar/insolar/network/node"
	"github.com/insolar/insolar/network/nodenetwork"
	"github.com/insolar/insolar/network/transport"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
)

var (
	keyProcessor = platformpolicy.NewKeyProcessor()
	scheme       = platformpolicy.NewPlatformCryptographyScheme()
)

func initNodes(
	ctx context.Context,
	mode consensus.Mode,
	nodes []insolar.NetworkNode,
	discoveryNodes []insolar.NetworkNode,
	strategy NetStrategy,
	nodeInfos []*nodeInfo,
) ([]consensus.Controller, []network.PulseHandler, []transport.DatagramTransport, []context.Context, []profiles.StaticProfile, error) {

	controllers := make([]consensus.Controller, len(nodes))
	transports := make([]transport.DatagramTransport, len(nodes))
	contexts := make([]context.Context, len(nodes))
	pulseHandlers := make([]network.PulseHandler, 0, len(nodes))
	staticProfiles := make([]profiles.StaticProfile, len(nodes))

	for i, n := range nodes {
		nodeKeeper := nodenetwork.NewNodeKeeper(n)
		nodeKeeper.SetInitialSnapshot(nodes)
		certificateManager := initCrypto(n, discoveryNodes)
		datagramHandler := adapters.NewDatagramHandler()

		conf := configuration.NewHostNetwork().Transport
		conf.Address = n.Address()

		transportFactory := transport.NewFactory(conf)
		datagramTransport, err := transportFactory.CreateDatagramTransport(datagramHandler)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}

		pulseHandler := adapters.NewPulseHandler()
		pulseHandlers = append(pulseHandlers, pulseHandler)

		delayTransport := strategy.GetLink(datagramTransport)
		transports[i] = delayTransport

		controllers[i] = consensus.New(ctx, consensus.Dep{
			PrimingCloudStateHash: [64]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			KeyProcessor:          keyProcessor,
			Scheme:                scheme,
			CertificateManager:    certificateManager,
			KeyStore:              keystore.NewInplaceKeyStore(nodeInfos[i].privateKey),
			NodeKeeper:            nodeKeeper,
			StateGetter:           &nshGen{nshDelay: defaultNshGenerationDelay},
			PulseChanger: &pulseChanger{
				nodeKeeper: nodeKeeper,
			},
			StateUpdater: &stateUpdater{
				nodeKeeper: nodeKeeper,
			},
			DatagramTransport: delayTransport,
		}).ControllerFor(mode, datagramHandler, pulseHandler)

		ctx, _ = inslogger.WithFields(ctx, map[string]interface{}{
			"node_id":      n.ShortID(),
			"node_address": n.Address(),
		})
		contexts[i] = ctx
		err = delayTransport.Start(ctx)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}

		staticProfiles[i] = adapters.NewStaticProfile(n, certificateManager.GetCertificate(), keyProcessor)
	}

	return controllers, pulseHandlers, transports, contexts, staticProfiles, nil
}

func initLogger(level insolar.LogLevel) context.Context {
	ctx := context.Background()
	logger := inslogger.FromContext(ctx).WithCaller(false)
	logger, _ = logger.WithLevelNumber(level)
	logger, _ = logger.WithFormat(insolar.TextFormat)
	ctx = inslogger.SetLogger(ctx, logger)
	return ctx
}

func generateNodeIdentities(countNeutral, countHeavy, countLight, countVirtual int) []nodeIdentity {
	r := make([]nodeIdentity, 0, countNeutral+countHeavy+countLight+countVirtual)

	r = _generateNodeIdentity(r, countNeutral, insolar.StaticRoleUnknown)
	r = _generateNodeIdentity(r, countHeavy, insolar.StaticRoleHeavyMaterial)
	r = _generateNodeIdentity(r, countLight, insolar.StaticRoleLightMaterial)
	r = _generateNodeIdentity(r, countVirtual, insolar.StaticRoleVirtual)

	return r
}

var portOffset = 10000

func _generateNodeIdentity(r []nodeIdentity, count int, role insolar.StaticRole) []nodeIdentity {
	for i := 0; i < count; i++ {
		port := portOffset
		r = append(r, nodeIdentity{
			role: role,
			addr: fmt.Sprintf("127.0.0.1:%d", port),
		})
		portOffset += 1
	}
	return r
}

func generateNodeInfos(nodeIdentities []nodeIdentity) []*nodeInfo {
	nodeInfos := make([]*nodeInfo, 0, len(nodeIdentities))
	for _, ni := range nodeIdentities {
		privateKey, _ := keyProcessor.GeneratePrivateKey()
		publicKey := keyProcessor.ExtractPublicKey(privateKey)

		nodeInfos = append(nodeInfos, &nodeInfo{
			nodeIdentity: ni,
			publicKey:    publicKey,
			privateKey:   privateKey,
		})
	}
	return nodeInfos
}

type nodeIdentity struct {
	role insolar.StaticRole
	addr string
}

type nodeInfo struct {
	nodeIdentity
	privateKey crypto.PrivateKey
	publicKey  crypto.PublicKey
}

func getAnnounceSignature(
	node insolar.NetworkNode,
	isDiscovery bool,
	kp insolar.KeyProcessor,
	key *ecdsa.PrivateKey,
	scheme insolar.PlatformCryptographyScheme,
) ([]byte, insolar.Signature) {

	brief := serialization.NodeBriefIntro{}
	brief.ShortID = node.ShortID()
	brief.SetPrimaryRole(adapters.StaticRoleToPrimaryRole(node.Role()))
	if isDiscovery {
		brief.SpecialRoles = member.SpecialRoleDiscovery
	}
	brief.StartPower = 10

	addr, err := endpoints.NewIPAddress(node.Address())
	if err != nil {
		panic(err)
	}
	copy(brief.Endpoint[:], addr[:])

	pk, err := kp.ExportPublicKeyBinary(node.PublicKey())
	if err != nil {
		panic(err)
	}

	copy(brief.NodePK[:], pk)

	buf := &bytes.Buffer{}
	err = brief.SerializeTo(nil, buf)
	if err != nil {
		panic(err)
	}

	data := buf.Bytes()
	data = data[:len(data)-64]

	digest := scheme.IntegrityHasher().Hash(data)
	sign, err := scheme.DigestSigner(key).Sign(digest)
	if err != nil {
		panic(err)
	}

	return digest, *sign
}

func nodesFromInfo(nodeInfos []*nodeInfo) ([]insolar.NetworkNode, []insolar.NetworkNode) {
	nodes := make([]insolar.NetworkNode, len(nodeInfos))
	discoveryNodes := make([]insolar.NetworkNode, 0)

	for i, info := range nodeInfos {
		var isDiscovery bool
		if info.role == insolar.StaticRoleHeavyMaterial || info.role == insolar.StaticRoleUnknown {
			isDiscovery = true
		}

		nn := newNetworkNode(info.addr, info.role, info.publicKey)
		nodes[i] = nn
		if isDiscovery {
			discoveryNodes = append(discoveryNodes, nn)
		}

		d, s := getAnnounceSignature(
			nn,
			isDiscovery,
			keyProcessor,
			info.privateKey.(*ecdsa.PrivateKey),
			scheme,
		)
		nn.(node.MutableNode).SetSignature(d, s)
	}

	return nodes, discoveryNodes
}

var shortNodeIdOffset = 1000

func newNetworkNode(addr string, role insolar.StaticRole, pk crypto.PublicKey) node.MutableNode {
	n := node.NewNode(
		testutils.RandomRef(),
		role,
		pk,
		addr,
		"",
	)
	mn := n.(node.MutableNode)
	mn.SetShortID(insolar.ShortNodeID(shortNodeIdOffset))

	shortNodeIdOffset += 1
	return mn
}

func initCrypto(node insolar.NetworkNode, discoveryNodes []insolar.NetworkNode) *certificate.CertificateManager {
	pubKey := node.PublicKey()

	publicKey, _ := keyProcessor.ExportPublicKeyPEM(pubKey)

	bootstrapNodes := make([]certificate.BootstrapNode, 0, len(discoveryNodes))
	for _, dn := range discoveryNodes {
		pubKey := dn.PublicKey()
		pubKeyBuf, _ := keyProcessor.ExportPublicKeyPEM(pubKey)

		bootstrapNode := certificate.NewBootstrapNode(
			pubKey,
			string(pubKeyBuf[:]),
			dn.Address(),
			dn.ID().String(),
		)
		bootstrapNodes = append(bootstrapNodes, *bootstrapNode)
	}

	cert := &certificate.Certificate{
		AuthorizationCertificate: certificate.AuthorizationCertificate{
			PublicKey: string(publicKey[:]),
			Reference: node.ID().String(),
			Role:      node.Role().String(),
		},
		BootstrapNodes: bootstrapNodes,
	}

	// dump cert and read it again from json for correct private files initialization
	jsonCert, _ := cert.Dump()
	cert, _ = certificate.ReadCertificateFromReader(pubKey, keyProcessor, strings.NewReader(jsonCert))
	return certificate.NewCertificateManager(cert)
}

const defaultNshGenerationDelay = time.Millisecond * 0

type nshGen struct {
	nshDelay time.Duration
}

func (ng *nshGen) State() []byte {
	delay := ng.nshDelay
	if delay != 0 {
		time.Sleep(delay)
	}

	nshBytes := make([]byte, 64)
	rand.Read(nshBytes)

	return nshBytes
}

type pulseChanger struct {
	nodeKeeper network.NodeKeeper
}

func (pc *pulseChanger) ChangePulse(ctx context.Context, pulse insolar.Pulse) {
	inslogger.FromContext(ctx).Info(">>>>>> Change pulse called")
	err := pc.nodeKeeper.MoveSyncToActive(ctx, pulse.PulseNumber)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}
}

type stateUpdater struct {
	nodeKeeper network.NodeKeeper
}

func (su *stateUpdater) UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte) {
	inslogger.FromContext(ctx).Info(">>>>>> Update state called")

	err := su.nodeKeeper.Sync(ctx, nodes, nil)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
	}
	su.nodeKeeper.SetCloudHash(cloudStateHash)
}
