package tests

import (
	"context"
	"io"
	"math/rand"

	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/member"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/proofs"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/statevector"
	"github.com/insolar/insolar/network/consensus/gcpv2/api/transport"
)

var EmuDefaultPacketBuilder transport.PacketBuilder = &emuPacketBuilder{}
var EmuPrimingHash = cryptkit.NewDigest(NewEmuNodeStateHash(1234567890), "stubHash").AsDigestHolder()
var EmuDefaultCryptography transport.CryptographyAssistant = &emuTransportCryptography{EmuPrimingHash}

func NewEmuTransport(sender transport.PacketSender) transport.Factory {
	return &emuTransport{sender}
}

type emuTransport struct {
	sender transport.PacketSender
}

func (r *emuTransport) GetPacketSender() transport.PacketSender {
	return r.sender
}

func (r *emuTransport) GetPacketBuilder(signer cryptkit.DigestSigner) transport.PacketBuilder {
	return EmuDefaultPacketBuilder
}

func (r *emuTransport) GetCryptographyFactory() transport.CryptographyAssistant {
	return EmuDefaultCryptography
}

type emuPackerCloner interface {
	clonePacketFor(target transport.TargetProfile, sendOptions transport.PacketSendOptions) transport.PacketParser
}

type emuPacketSender struct {
	cloner emuPackerCloner
}

func (r *emuPacketSender) SendTo(ctx context.Context, t transport.TargetProfile, sendOptions transport.PacketSendOptions, s transport.PacketSender) {
	c := r.cloner.clonePacketFor(t, sendOptions)
	s.SendPacketToTransport(ctx, t, sendOptions, c)
}

func (r *emuPacketSender) SendToMany(ctx context.Context, targetCount int, s transport.PacketSender, filter transport.ProfileFilter) {

	for i := 0; i < targetCount; i++ {
		sendTo, sendOptions := filter(ctx, i)
		if sendTo != nil {
			c := r.cloner.clonePacketFor(sendTo, sendOptions)
			s.SendPacketToTransport(ctx, sendTo, sendOptions, c)
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

type emuPacketBuilder struct {
}

func (r *emuPacketBuilder) GetNeighbourhoodSize() transport.NeighbourhoodSizes {
	return transport.NeighbourhoodSizes{NeighbourhoodSize: 5, NeighbourhoodTrustThreshold: 2, JoinersPerNeighbourhood: 2, JoinersBoost: 1}
}

func (r *emuPacketBuilder) PreparePhase0Packet(sender *transport.NodeAnnouncementProfile, pulsarPacket proofs.OriginalPulsarPacket,
	options transport.PacketPrepareOptions) transport.PreparedPacketSender {
	v := EmuPhase0NetPacket{
		basePacket: basePacket{
			src:       sender.GetNodeID(),
			nodeCount: sender.GetNodeCount(),
			mp:        sender.GetMembershipProfile(),
		},
		pulsePacket: pulsarPacket.(*EmuPulsarNetPacket)}

	if options&transport.AlternativePhasePacket != 0 {
		panic("illegal value")
	}

	return &emuPacketSender{&v}
}

func (r *EmuPhase0NetPacket) clonePacketFor(t transport.TargetProfile, sendOptions transport.PacketSendOptions) transport.PacketParser {
	c := *r
	c.tgt = t.GetNodeID()
	return &c
}

func (r *emuPacketBuilder) PreparePhase1Packet(sender *transport.NodeAnnouncementProfile, pulsarPacket proofs.OriginalPulsarPacket,
	welcome *proofs.NodeWelcomePackage, options transport.PacketPrepareOptions) transport.PreparedPacketSender {

	pp := pulsarPacket.(*EmuPulsarNetPacket)
	if pp == nil || !pp.pulseData.IsValidPulseData() {
		panic("pulse data is missing or invalid")
	}

	v := EmuPhase1NetPacket{
		EmuPhase0NetPacket: EmuPhase0NetPacket{
			basePacket: basePacket{
				src:                       sender.GetNodeID(),
				nodeCount:                 sender.GetNodeCount(),
				mp:                        sender.GetMembershipProfile(),
				isLeaving:                 sender.IsLeaving(),
				leaveReason:               sender.GetLeaveReason(),
				joiner:                    sender.GetJoinerAnnouncement(),
				cloudIntro:                welcome,
				BriefCandidateProfile:     sender.GetStatic(),
				CandidateProfileExtension: sender.GetStatic().GetExtension(),
			},
			pulsePacket: pp},
	}
	v.basePacket.adjustBySender(sender)

	v.pn = pp.pulseData.PulseNumber
	v.isAlternative = options&transport.AlternativePhasePacket != 0

	if v.isAlternative || options&transport.PrepareWithoutPulseData != 0 {
		v.pulsePacket = nil
	}

	if v.joiner != nil && v.joiner.HasFullIntro() && options&transport.OnlyBriefIntroAboutJoiner != 0 {
		v.joiner = transport.NewBriefJoinerAnnouncementByFull(v.joiner)
	}

	return &emuPacketSender{&v}
}

func (r *EmuPhase1NetPacket) clonePacketFor(t transport.TargetProfile, sendOptions transport.PacketSendOptions) transport.PacketParser {
	c := *r
	c.tgt = t.GetNodeID()

	if c.joiner != nil && c.joiner.HasFullIntro() && c.GetJoinerID() == t.GetNodeID() {
		c.joiner = transport.NewBriefJoinerAnnouncementByFull(c.joiner)
	}

	if !t.IsJoiner() {
		c.cloudIntro = nil
	}

	if !c.mp.IsJoiner() && !t.IsJoiner() {
		c.BriefCandidateProfile = nil
		c.CandidateProfileExtension = nil
	}

	if sendOptions&transport.SendWithoutPulseData != 0 {
		c.pulsePacket = nil
	}

	return &c
}

func (r *emuPacketBuilder) PreparePhase2Packet(sender *transport.NodeAnnouncementProfile,
	welcome *proofs.NodeWelcomePackage, neighbourhood []transport.MembershipAnnouncementReader,
	options transport.PacketPrepareOptions) transport.PreparedPacketSender {

	v := EmuPhase2NetPacket{
		basePacket: basePacket{
			src:                       sender.GetNodeID(),
			nodeCount:                 sender.GetNodeCount(),
			mp:                        sender.GetMembershipProfile(),
			isLeaving:                 sender.IsLeaving(),
			leaveReason:               sender.GetLeaveReason(),
			joiner:                    sender.GetJoinerAnnouncement(),
			cloudIntro:                welcome,
			BriefCandidateProfile:     sender.GetStatic(),
			CandidateProfileExtension: sender.GetStatic().GetExtension(),
		},
		pulseNumber:   sender.GetPulseNumber(),
		neighbourhood: neighbourhood,
	}
	v.basePacket.adjustBySender(sender)

	v.isAlternative = options&transport.AlternativePhasePacket != 0

	if v.joiner != nil && v.joiner.HasFullIntro() && options&transport.OnlyBriefIntroAboutJoiner != 0 {
		v.joiner = transport.NewBriefJoinerAnnouncementByFull(v.joiner)
	}

	return &emuPacketSender{&v}
}

func (r *EmuPhase2NetPacket) clonePacketFor(t transport.TargetProfile, sendOptions transport.PacketSendOptions) transport.PacketParser {
	c := *r
	c.tgt = t.GetNodeID()

	if c.joiner != nil && c.joiner.HasFullIntro() && c.GetJoinerID() == t.GetNodeID() {
		c.joiner = transport.NewBriefJoinerAnnouncementByFull(c.joiner)
	}

	if !t.IsJoiner() {
		c.cloudIntro = nil
	}

	if !c.mp.IsJoiner() && !t.IsJoiner() {
		c.BriefCandidateProfile = nil
		c.CandidateProfileExtension = nil
	}

	return &c
}

func (r *emuPacketBuilder) PreparePhase3Packet(sender *transport.NodeAnnouncementProfile, vectors statevector.Vector,
	options transport.PacketPrepareOptions) transport.PreparedPacketSender {

	v := EmuPhase3NetPacket{
		basePacket: basePacket{
			src:       sender.GetNodeID(),
			nodeCount: sender.GetNodeCount(),
			mp:        sender.GetMembershipProfile(),
		},
		pulseNumber: sender.GetPulseNumber(),
		vectors:     vectors,
	}
	v.basePacket.adjustBySender(sender)

	v.isAlternative = options&transport.AlternativePhasePacket != 0

	return &emuPacketSender{&v}
}

func (r *EmuPhase3NetPacket) clonePacketFor(t transport.TargetProfile, sendOptions transport.PacketSendOptions) transport.PacketParser {
	c := *r
	c.tgt = t.GetNodeID()
	return &c
}

type emuTransportCryptography struct {
	defaultDigest cryptkit.DigestHolder
}

func (r *emuTransportCryptography) CreateSequenceDigester() cryptkit.SequenceDigester {
	return &seqDigester{}
}

func (r *emuTransportCryptography) CreateGlobulaStateDigester() transport.StateDigester {
	return &gshDigester{&seqDigester{}, r.defaultDigest}
}

func (r *emuTransportCryptography) CreatePublicKeyStore(skh cryptkit.SignatureKeyHolder) cryptkit.PublicKeyStore {
	return nil
}

func (r *emuTransportCryptography) CreatePacketDigester() cryptkit.DataDigester {
	panic("not implemented")
}

func (r *emuTransportCryptography) CreateAnnouncementDigester() cryptkit.SequenceDigester {
	return &seqDigester{}
}

func (r *emuTransportCryptography) IsDigestMethodSupported(m cryptkit.DigestMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsValidDataSignature(data io.Reader, signature cryptkit.SignatureHolder) bool {
	return true
}

func (r *emuTransportCryptography) IsSignOfSignatureMethodSupported(m cryptkit.SignatureMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsDigestOfSignatureMethodSupported(m cryptkit.SignatureMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsSignMethodSupported(m cryptkit.SignMethod) bool {
	return true
}

func (r *emuTransportCryptography) IsValidDigestSignature(digest cryptkit.DigestHolder, signature cryptkit.SignatureHolder) bool {
	return true
}

func (r *emuTransportCryptography) SignDigest(digest cryptkit.Digest) cryptkit.Signature {
	return cryptkit.NewSignature(digest, digest.GetDigestMethod().SignedBy(r.GetSignMethod()))
}

func (r *emuTransportCryptography) GetSignMethod() cryptkit.SignMethod {
	return "emuSing"
}

func (r *emuTransportCryptography) CreateSignatureVerifierWithPKS(pks cryptkit.PublicKeyStore) cryptkit.SignatureVerifier {
	return r
}

func (r *emuTransportCryptography) GetDigestFactory() transport.ConsensusDigestFactory {
	return r
}

func (r *emuTransportCryptography) CreateNodeSigner(sks cryptkit.SecretKeyStore) cryptkit.DigestSigner {
	return r
}

type seqDigester struct {
	// TODO do test or a proper digest calc
	rnd      *rand.Rand
	lastSeed int64
}

func (s *seqDigester) AddNext(digest longbits.FoldableReader) {
	// it is a dirty emulation of digest
	if s.rnd == nil {
		s.rnd = rand.New(rand.NewSource(0))
	}
	s.lastSeed = int64(s.rnd.Uint64() ^ digest.FoldToUint64())
	s.rnd.Seed(s.lastSeed)
}

func (s *seqDigester) GetDigestMethod() cryptkit.DigestMethod {
	return "emuDigest64"
}

func (s *seqDigester) ForkSequence() cryptkit.SequenceDigester {
	cp := seqDigester{}
	if s.rnd != nil {
		cp.rnd = rand.New(rand.NewSource(s.lastSeed))
	}
	return &cp
}

func (s *seqDigester) FinishSequence() cryptkit.Digest {
	if s.rnd == nil {
		panic("nothing")
	}
	bits := longbits.NewBits64(s.rnd.Uint64())
	s.rnd = nil
	return cryptkit.NewDigest(&bits, s.GetDigestMethod())
}

type gshDigester struct {
	sd            cryptkit.SequenceDigester
	defaultDigest longbits.FoldableReader
}

func (p *gshDigester) AddNext(digest longbits.FoldableReader, fullRank member.FullRank) {
	if digest == nil {
		p.sd.AddNext(p.defaultDigest)
	} else {
		p.sd.AddNext(digest)
	}
}

func (p *gshDigester) GetDigestMethod() cryptkit.DigestMethod {
	return p.sd.GetDigestMethod()
}

func (p *gshDigester) ForkSequence() transport.StateDigester {
	return &gshDigester{p.sd.ForkSequence(), p.defaultDigest}
}

func (p *gshDigester) FinishSequence() cryptkit.Digest {
	return p.sd.FinishSequence()
}
