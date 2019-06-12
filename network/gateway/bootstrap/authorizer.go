package bootstrap

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/component"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/version"
)

type AuthorizationController interface {
	component.Initer

	Authorize(ctx context.Context, discoveryNode *DiscoveryNode, cert insolar.AuthorizationCertificate) (*packet.AuthorizeResponse, error)
}

type authorizer struct {
	HostNetwork network.HostNetwork
	Gatewayer   network.Gatewayer

	OriginProvider      insolar.OriginProvider
	CryptographyService insolar.CryptographyService
}

func (ac *authorizer) Init(ctx context.Context) error {
	ac.HostNetwork.RegisterRequestHandler(types.Authorize, ac.HandleNodeAuthorizeRequest)
	return nil
}

func (ac *authorizer) Authorize(ctx context.Context, discoveryNode *DiscoveryNode, cert insolar.AuthorizationCertificate) (*packet.AuthorizeResponse, error) {
	inslogger.FromContext(ctx).Infof("Authorizing on host: %s", discoveryNode.Host)
	inslogger.FromContext(ctx).Infof("cert: %s", cert)

	ctx, span := instracer.StartSpan(ctx, "AuthorizationController.Authorize")
	span.AddAttributes(
		trace.StringAttribute("node", discoveryNode.Node.GetNodeRef().String()),
	)
	defer span.End()
	serializedCert, err := certificate.Serialize(cert)
	if err != nil {
		return nil, errors.Wrap(err, "Error serializing certificate")
	}

	authData := &packet.AuthorizeData{Certificate: serializedCert, Version: version.Version, Timestamp: time.Now().UTC().UnixNano()}

	// TODO: signature
	auth := &packet.AuthorizeRequest{AuthorizeData: authData, Signature: nil}

	future, err := ac.HostNetwork.SendRequestToHost(ctx, types.Authorize, auth, discoveryNode.Host)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending authorize request")
	}
	response, err := future.WaitResponse(time.Second * 5 /*ac.options.PacketTimeout*/)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for authorize request")
	}
	if response.GetResponse() == nil || response.GetResponse().GetAuthorize() == nil {
		return nil, errors.Errorf("Authorize failed: got incorrect response: %s", response)
	}
	data := response.GetResponse().GetAuthorize()
	switch data.Code {
	case packet.WrongVersion:
	case packet.WrongTimestamp:
	case packet.WrongMandate:
		panic("todo: ")
	case packet.Success:
		//TODO check majority rule
		return data.Permit
	}

	return data, nil
}

func (g *authorizer) HandleNodeAuthorizeRequest(ctx context.Context, request network.Packet) (network.Packet, error) {
	if request.GetRequest() == nil || request.GetRequest().GetAuthorize() == nil {
		return nil, errors.Errorf("process authorize: got invalid protobuf request message: %s", request)
	}
	data := request.GetRequest().GetAuthorize().AuthorizeData

	// check timestamp
	if data.Timestamp-time.Now().UTC().UnixNano() > 5 { // TODO: check delta
		return g.HostNetwork.BuildResponse(ctx, request, &packet.AuthorizeResponse{Code: packet.WrongTimestamp, Error: nil}), nil
	}

	cert, err := certificate.Deserialize(data.Certificate, platformpolicy.NewKeyProcessor())
	if err != nil {
		return g.HostNetwork.BuildResponse(ctx, request, &packet.AuthorizeResponse{Code: packet.WrongMandate, Error: err.Error()}), nil
	}

	valid, err := g.Gatewayer.Gateway().Auther().ValidateCert(ctx, cert)
	if !valid {
		if err == nil {
			err = errors.New("Certificate validation failed")
		}
		return g.HostNetwork.BuildResponse(ctx, request, &packet.AuthorizeResponse{Code: packet.WrongMandate, Error: err.Error()}), nil
	}

	// pulse, _ := g.PulseAccessor.Latest(ctx)

	// TODO: public key to bytes
	permit, _ := CreatePermit(g.OriginProvider.GetOrigin().ID(), nil /* bc.getRandActiveDiscoveryAddress() */, cert.GetPublicKey(), g.CryptographyService)

	//discoveryCount := FindDiscoveryInActiveList(g.NodeKeeper.GetAccessor().GetActiveNodes())
	var discoveryCount uint32
	return g.HostNetwork.BuildResponse(ctx, request, &packet.AuthorizeResponse{
		Code: packet.Success, Error: nil,
		Permit: permit, DiscoveryCount: discoveryCount,
		PulseNumber:  0, // TODO
		NetworkState: uint32(g.Gatewayer.Gateway().GetState()),
	}), nil
}
