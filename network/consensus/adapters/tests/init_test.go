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
	"github.com/insolar/insolar/insolar/gen"
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
)

var (
	keyProcessor = platformpolicy.NewKeyProcessor()
	scheme       = platformpolicy.NewPlatformCryptographyScheme()
)

var (
	shortNodeIdOffset = 1000
	portOffset        = 10000
)

type candidate struct {
	profiles.StaticProfile
	profiles.StaticProfileExtension
}

func testCase(stopAfter, startCaseAfter time.Duration, test func()) {
	startedAt := time.Now()

	ticker := time.NewTicker(time.Second)
	stopTest := time.After(stopAfter)
	startCase := time.After(startCaseAfter)
	for {
		select {
		case <-ticker.C:
			fmt.Println("===", time.Since(startedAt), "=================================================")
		case <-stopTest:
			return
		case <-startCase:
			test()
		}
	}
}

type InitializedNodes struct {
	controllers    []consensus.Controller
	nodeKeepers    []network.NodeKeeper
	transports     []transport.DatagramTransport
	contexts       []context.Context
	pulseHandlers  []network.PulseHandler
	staticProfiles []profiles.StaticProfile
}

type GeneratedNodes struct {
	nodes          []insolar.NetworkNode
	meta           []*nodeMeta
	discoveryNodes []insolar.NetworkNode
}

func generateNodes(countNeutral, countHeavy, countLight, countVirtual int, discoveryNodes []insolar.NetworkNode) (*GeneratedNodes, error) {
	nodeIdentities := generateNodeIdentities(countNeutral, countHeavy, countLight, countVirtual)
	nodeInfos := generateNodeInfos(nodeIdentities)
	nodes, dn, err := nodesFromInfo(nodeInfos)

	if len(discoveryNodes) > 0 {
		dn = discoveryNodes
	}

	if err != nil {
		return nil, err
	}

	return &GeneratedNodes{
		nodes:          nodes,
		meta:           nodeInfos,
		discoveryNodes: dn,
	}, nil
}

func newNodes(size int) InitializedNodes {
	return InitializedNodes{
		controllers:    make([]consensus.Controller, size),
		transports:     make([]transport.DatagramTransport, size),
		contexts:       make([]context.Context, size),
		pulseHandlers:  make([]network.PulseHandler, 0, size),
		staticProfiles: make([]profiles.StaticProfile, size),
		nodeKeepers:    make([]network.NodeKeeper, size),
	}
}

func initNodes(ctx context.Context, mode consensus.Mode, nodes GeneratedNodes, strategy NetStrategy) (*InitializedNodes, error) {
	ns := newNodes(len(nodes.nodes))

	for i, n := range nodes.nodes {
		nodeKeeper := nodenetwork.NewNodeKeeper(n)
		nodeKeeper.SetInitialSnapshot(nodes.nodes)
		ns.nodeKeepers[i] = nodeKeeper

		certificateManager := initCrypto(n, nodes.discoveryNodes)
		datagramHandler := adapters.NewDatagramHandler()

		conf := configuration.NewHostNetwork().Transport
		conf.Address = n.Address()

		transportFactory := transport.NewFactory(conf)
		datagramTransport, err := transportFactory.CreateDatagramTransport(datagramHandler)
		if err != nil {
			return nil, err
		}

		pulseHandler := adapters.NewPulseHandler()
		ns.pulseHandlers = append(ns.pulseHandlers, pulseHandler)

		delayTransport := strategy.GetLink(datagramTransport)
		ns.transports[i] = delayTransport

		controller := consensus.New(ctx, consensus.Dep{
			KeyProcessor:       keyProcessor,
			Scheme:             scheme,
			CertificateManager: certificateManager,
			KeyStore:           keystore.NewInplaceKeyStore(nodes.meta[i].privateKey),

			NodeKeeper:        nodeKeeper,
			DatagramTransport: delayTransport,

			StateGetter: &nshGen{nshDelay: defaultNshGenerationDelay},
			PulseChanger: &pulseChanger{
				nodeKeeper: nodeKeeper,
			},
			StateUpdater: &stateUpdater{
				nodeKeeper: nodeKeeper,
			},
			EphemeralController: &ephemeralController{
				allowed: true,
			},
		}).ControllerFor(mode, datagramHandler, pulseHandler)

		ns.controllers[i] = controller
		ctx, _ = inslogger.WithFields(ctx, map[string]interface{}{
			"node_id":      n.ShortID(),
			"node_address": n.Address(),
		})
		ns.contexts[i] = ctx
		err = delayTransport.Start(ctx)
		if err != nil {
			return nil, err
		}

		ns.staticProfiles[i] = adapters.NewStaticProfile(n, certificateManager.GetCertificate(), keyProcessor)
	}

	return &ns, nil
}

func initPulsar(ctx context.Context, delta uint16, ns InitializedNodes) {
	pulsar := NewPulsar(delta, ns.pulseHandlers)
	go func() {
		for {
			pulsar.Pulse(ctx, 4+len(ns.staticProfiles)/10)
		}
	}()
}

func initLogger(level insolar.LogLevel) context.Context {
	ctx := context.Background()

	cfg := configuration.NewLog()
	cfg.LLBufferSize = 0
	cfg.Level = level.String()
	cfg.Formatter = insolar.TextFormat.String()

	ctx, _ = inslogger.InitNodeLogger(ctx, cfg, "", "")
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

func generateNodeInfos(nodeIdentities []nodeIdentity) []*nodeMeta {
	nodeInfos := make([]*nodeMeta, 0, len(nodeIdentities))
	for _, ni := range nodeIdentities {
		privateKey, _ := keyProcessor.GeneratePrivateKey()
		publicKey := keyProcessor.ExtractPublicKey(privateKey)

		nodeInfos = append(nodeInfos, &nodeMeta{
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

type nodeMeta struct {
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

	digest := scheme.IntegrityHasher().Hash(data)
	sign, err := scheme.DigestSigner(key).Sign(digest)
	if err != nil {
		return nil, nil, err
	}

	return digest, sign, nil
}

func nodesFromInfo(nodeInfos []*nodeMeta) ([]insolar.NetworkNode, []insolar.NetworkNode, error) {
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

		d, s, err := getAnnounceSignature(
			nn,
			isDiscovery,
			keyProcessor,
			info.privateKey.(*ecdsa.PrivateKey),
			scheme,
		)
		if err != nil {
			return nil, nil, err
		}
		nn.(node.MutableNode).SetSignature(d, *s)
	}

	return nodes, discoveryNodes, nil
}

func newNetworkNode(addr string, role insolar.StaticRole, pk crypto.PublicKey) node.MutableNode {
	n := node.NewNode(
		gen.Reference(),
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
			dn.Role().String(),
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
	pc.nodeKeeper.MoveSyncToActive(ctx, pulse.PulseNumber)
}

type stateUpdater struct {
	nodeKeeper network.NodeKeeper
}

func (su *stateUpdater) UpdateState(ctx context.Context, pulseNumber insolar.PulseNumber, nodes []insolar.NetworkNode, cloudStateHash []byte) {
	inslogger.FromContext(ctx).Info(">>>>>> Update state called")

	su.nodeKeeper.Sync(ctx, pulseNumber, nodes)
}

type ephemeralController struct {
	allowed bool
}

func (e *ephemeralController) EphemeralMode(nodes []insolar.NetworkNode) bool {
	return e.allowed
}
