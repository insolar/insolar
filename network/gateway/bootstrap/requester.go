package bootstrap

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/consensusv1/packets"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
	"github.com/insolar/insolar/version"
)

type Requester interface {
	Authorize(context.Context, *host.Host, insolar.AuthorizationCertificate) (*packet.AuthorizeResponse, error)
	Bootstrap(context.Context, *packet.Permit, *packets.NodeJoinClaim, *insolar.Pulse) (*packet.BootstrapResponse, error)
	UpdateSchedule(context.Context, *packet.Permit, insolar.PulseNumber) (*packet.UpdateScheduleResponse, error)
}

func NewRequester(options *common.Options) Requester {
	return &requester{options: options}
}

type requester struct {
	HostNetwork         network.HostNetwork         `inject:""`
	OriginProvider      insolar.OriginProvider      `inject:""`
	CryptographyService insolar.CryptographyService `inject:""`

	options *common.Options
}

func (ac *requester) Authorize(ctx context.Context, host *host.Host, cert insolar.AuthorizationCertificate) (*packet.AuthorizeResponse, error) {
	inslogger.FromContext(ctx).Infof("Authorizing on host: %s", host.String())

	ctx, span := instracer.StartSpan(ctx, "AuthorizationController.Authorize")
	span.AddAttributes(
		trace.StringAttribute("node", host.NodeID.String()),
	)
	defer span.End()
	serializedCert, err := certificate.Serialize(cert)
	if err != nil {
		return nil, errors.Wrap(err, "Error serializing certificate")
	}

	authData := &packet.AuthorizeData{Certificate: serializedCert, Version: version.Version}
	response, err := ac.authorizeWithTimestamp(ctx, host, authData, time.Now().UTC().UnixNano())
	if err != nil {
		return nil, err
	}

	switch response.Code {
	case packet.Success:
		return response, nil
	case packet.WrongMandate:
		return response, errors.New("Failed to authorize, wrong mandate.")
	case packet.WrongVersion:
		return response, errors.New("Failed to authorize, wrong version.")
	}

	// retry with received timestamp
	// TODO: change one retry to many
	response, err = ac.authorizeWithTimestamp(ctx, host, authData, response.Timestamp)
	return response, err
}

func (ac *requester) authorizeWithTimestamp(ctx context.Context, host *host.Host, authData *packet.AuthorizeData, timestamp int64) (*packet.AuthorizeResponse, error) {

	authData.Timestamp = timestamp

	data, err := authData.Marshal()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal permit")
	}

	signature, err := ac.CryptographyService.Sign(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign permit")
	}

	req := &packet.AuthorizeRequest{AuthorizeData: authData, Signature: signature.Bytes()}

	f, err := ac.HostNetwork.SendRequestToHost(ctx, types.Authorize, req, host)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending authorize request")
	}
	response, err := f.WaitResponse(ac.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for authorize request")
	}

	if response.GetResponse().GetError() != nil {
		return nil, errors.New(response.GetResponse().GetError().Error)
	}

	if response.GetResponse() == nil || response.GetResponse().GetAuthorize() == nil {
		return nil, errors.Errorf("Authorize failed: got incorrect response: %s", response)
	}

	return response.GetResponse().GetAuthorize(), nil
}

func (ac *requester) Bootstrap(ctx context.Context, permit *packet.Permit, joinClaim *packets.NodeJoinClaim, p *insolar.Pulse) (*packet.BootstrapResponse, error) {

	req := &packet.BootstrapRequest{
		JoinClaim: joinClaim,
		Pulse:     pulse.ToProto(p),
		Permit:    permit,
	}

	f, err := ac.HostNetwork.SendRequestToHost(ctx, types.Bootstrap, req, permit.Payload.ReconnectTo)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending bootstrap request")
	}

	resp, err := f.WaitResponse(ac.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for bootstrap request")
	}

	respData := resp.GetResponse().GetBootstrap()
	switch respData.Code {
	case packet.UpdateShortID:
	// make claim with new  shortID
	// bootstrap again
	case packet.UpdateSchedule:
		ac.UpdateSchedule(ctx, permit, p.PulseNumber)
		panic("call bootstrap again")
	case packet.Reject:
		//swith to no network
		return respData, errors.New("Bootstrap request rejected")
	}

	// case Accepted
	return respData, nil

}

func (ac *requester) UpdateSchedule(ctx context.Context, permit *packet.Permit, pulse insolar.PulseNumber) (*packet.UpdateScheduleResponse, error) {

	req := &packet.UpdateScheduleRequest{
		LastNodePulse: pulse,
		Permit:        permit,
	}

	f, err := ac.HostNetwork.SendRequestToHost(ctx, types.UpdateSchedule, req, permit.Payload.ReconnectTo)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending bootstrap request")
	}

	resp, err := f.WaitResponse(ac.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for bootstrap request")
	}

	return resp.GetResponse().GetUpdateSchedule(), nil
}
