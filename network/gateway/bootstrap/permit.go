package bootstrap

import (
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
)

// CreatePermit creates permit for joiner node to
func CreatePermit(authorityNodeRef insolar.Reference, reconnectHost *host.Host, joinerPublicKey []byte, signer insolar.Signer) (*packet.Permit, error) {
	payload := packet.PermitPayload{
		AuthorityNodeRef: authorityNodeRef,
		ExpireTimestamp:  time.Now().Unix() + 60 * 5,
		ReconnectTo:      reconnectHost,
		JoinerPublicKey:  joinerPublicKey,
	}

	data, err := payload.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal bootstrap permit")
	}
	signature, err := signer.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign bootstrap permit")
	}

	return &packet.Permit{Payload: payload, Signature: signature.Bytes()}, nil
}

// ValidatePermit validate granted permit and verifies signature of Authority Node
func ValidatePermit(permit *packet.Permit, cert insolar.Certificate, verifier insolar.CryptographyService) error {
	discovery := network.FindDiscoveryByRef(cert, permit.Payload.AuthorityNodeRef)
	if discovery == nil {
		return errors.New("failed to find a discovery node from reference in permit")
	}

	payload, err := permit.Payload.Marshal()
	if err != nil {
		return errors.New("failed to marshal bootstrap permission payload part")
	}
	verified := verifier.Verify(discovery.GetPublicKey(), insolar.SignatureFromBytes(permit.Signature), payload)
	if !verified {
		return errors.New("bootstrap permission payload verification failed")
	}
	return nil
}

// func (bc *Bootstrap) getRandActiveDiscoveryAddress() string {
// 	if len(bc.NodeKeeper.GetAccessor().GetActiveNodes()) <= 1 {
// 		return bc.NodeKeeper.GetOrigin().Address()
// 	}
//
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	index := r.Intn(len(bc.Certificate.GetDiscoveryNodes()))
// 	n := bc.NodeKeeper.GetAccessor().GetActiveNode(*bc.Certificate.GetDiscoveryNodes()[index].GetNodeRef())
// 	if (n != nil) && (n.GetState() == insolar.NodeReady) {
// 		return bc.Certificate.GetDiscoveryNodes()[index].GetHost()
// 	}
//
// 	return bc.NodeKeeper.GetOrigin().Address()
// }
