// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package bootstrap

import (
	"context"
	"math/rand"
	"time"

	"github.com/opentracing/opentracing-go/log"
	"github.com/pkg/errors"

	"github.com/insolar/insolar/network/consensus/adapters"

	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/hostnetwork/host"
	"github.com/insolar/insolar/network/hostnetwork/packet"
	"github.com/insolar/insolar/network/hostnetwork/packet/types"
)

//go:generate minimock -i github.com/insolar/insolar/network/gateway/bootstrap.Requester -o ./ -s _mock.go -g

type Requester interface {
	Authorize(context.Context, insolar.Certificate) (*packet.Permit, error)
	Bootstrap(context.Context, *packet.Permit, adapters.Candidate, *insolar.Pulse) (*packet.BootstrapResponse, error)
	UpdateSchedule(context.Context, *packet.Permit, insolar.PulseNumber) (*packet.UpdateScheduleResponse, error)
	Reconnect(context.Context, *host.Host, *packet.Permit) (*packet.ReconnectResponse, error)
}

func NewRequester(options *network.Options) Requester {
	return &requester{options: options}
}

type requester struct {
	HostNetwork         network.HostNetwork         `inject:""`
	OriginProvider      network.OriginProvider      `inject:""`
	CryptographyService insolar.CryptographyService `inject:""`

	options *network.Options
}

func (ac *requester) Authorize(ctx context.Context, cert insolar.Certificate) (*packet.Permit, error) {
	logger := inslogger.FromContext(ctx)

	discoveryNodes := network.ExcludeOrigin(cert.GetDiscoveryNodes(), *cert.GetNodeRef())

	rand.Shuffle(
		len(discoveryNodes),
		func(i, j int) {
			discoveryNodes[i], discoveryNodes[j] = discoveryNodes[j], discoveryNodes[i]
		},
	)

	bestResult := &packet.AuthorizeResponse{}

	for _, n := range discoveryNodes {
		h, err := host.NewHostN(n.GetHost(), *n.GetNodeRef())
		if err != nil {
			logger.Warnf("Error authorizing to mallformed host %s[%s]: %s",
				n.GetHost(), *n.GetNodeRef(), err.Error())
			continue
		}

		logger.Infof("Trying to authorize to node: %s", h.String())
		res, err := ac.authorize(ctx, h, cert)
		if err != nil {
			logger.Warnf("Error authorizing to host %s: %s", h.String(), err.Error())
			continue
		}

		if int(res.DiscoveryCount) < cert.GetMajorityRule() {
			logger.Infof(
				"Check MajorityRule failed on authorize, expect %d, got %d",
				cert.GetMajorityRule(),
				res.DiscoveryCount,
			)

			if res.DiscoveryCount > bestResult.DiscoveryCount {
				bestResult = res
			}

			continue
		}

		return res.Permit, nil
	}

	if network.OriginIsDiscovery(cert) && bestResult.Permit != nil {
		return bestResult.Permit, nil
	}

	return nil, errors.New("failed to authorize to any discovery node")
}

func (ac *requester) authorize(ctx context.Context, host *host.Host, cert insolar.AuthorizationCertificate) (*packet.AuthorizeResponse, error) {
	inslogger.FromContext(ctx).Infof("Authorizing on host: %s", host.String())

	ctx, span := instracer.StartSpan(ctx, "AuthorizationController.Authorize")
	span.LogFields(
		log.String("node", host.NodeID.String()),
	)
	defer span.Finish()
	serializedCert, err := certificate.Serialize(cert)
	if err != nil {
		return nil, errors.Wrap(err, "Error serializing certificate")
	}

	authData := &packet.AuthorizeData{Certificate: serializedCert, Version: ac.OriginProvider.GetOrigin().Version()}
	response, err := ac.authorizeWithTimestamp(ctx, host, authData, time.Now().Unix())
	if err != nil {
		return nil, err
	}

	switch response.Code {
	case packet.Success:
		return response, nil
	case packet.WrongMandate:
		return response, errors.New("failed to authorize, wrong mandate")
	case packet.WrongVersion:
		return response, errors.New("failed to authorize, wrong version")
	}

	// retry with received timestamp
	// TODO: change one retry to many
	response, err = ac.authorizeWithTimestamp(ctx, host, authData, response.Timestamp)
	return response, err
}

func (ac *requester) authorizeWithTimestamp(ctx context.Context, h *host.Host, authData *packet.AuthorizeData, timestamp int64) (*packet.AuthorizeResponse, error) {

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

	f, err := ac.HostNetwork.SendRequestToHost(ctx, types.Authorize, req, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending Authorize request")
	}
	response, err := f.WaitResponse(ac.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for Authorize request")
	}

	if response.GetResponse().GetError() != nil {
		return nil, errors.New(response.GetResponse().GetError().Error)
	}

	if response.GetResponse() == nil || response.GetResponse().GetAuthorize() == nil {
		return nil, errors.Errorf("Authorize failed: got incorrect response: %s", response)
	}

	return response.GetResponse().GetAuthorize(), nil
}

func (ac *requester) Bootstrap(ctx context.Context, permit *packet.Permit, candidate adapters.Candidate, p *insolar.Pulse) (*packet.BootstrapResponse, error) {

	req := &packet.BootstrapRequest{
		CandidateProfile: candidate.Profile(),
		Pulse:            *pulse.ToProto(p),
		Permit:           permit,
	}

	f, err := ac.HostNetwork.SendRequestToHost(ctx, types.Bootstrap, req, permit.Payload.ReconnectTo)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending Bootstrap request")
	}

	resp, err := f.WaitResponse(ac.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for Bootstrap request")
	}

	respData := resp.GetResponse().GetBootstrap()
	if respData == nil {
		return nil, errors.New("bad response for bootstrap")
	}

	switch respData.Code {
	case packet.UpdateShortID:
		return respData, errors.New("Bootstrap got UpdateShortID")
	case packet.UpdateSchedule:
		// ac.UpdateSchedule(ctx, permit, p.PulseNumber)
		// panic("call bootstrap again")
		return respData, errors.New("Bootstrap got UpdateSchedule")
	case packet.Reject:
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
		return nil, errors.Wrapf(err, "Error sending UpdateSchedule request")
	}

	resp, err := f.WaitResponse(ac.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for UpdateSchedule request")
	}

	return resp.GetResponse().GetUpdateSchedule(), nil
}

func (ac *requester) Reconnect(ctx context.Context, h *host.Host, permit *packet.Permit) (*packet.ReconnectResponse, error) {
	req := &packet.ReconnectRequest{
		ReconnectTo: *permit.Payload.ReconnectTo,
		Permit:      permit,
	}

	f, err := ac.HostNetwork.SendRequestToHost(ctx, types.Reconnect, req, h)
	if err != nil {
		return nil, errors.Wrapf(err, "Error sending Reconnect request")
	}

	resp, err := f.WaitResponse(ac.options.PacketTimeout)
	if err != nil {
		return nil, errors.Wrapf(err, "Error getting response for Reconnect request")
	}

	return resp.GetResponse().GetReconnect(), nil
}
