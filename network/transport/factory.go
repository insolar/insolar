// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package transport

import (
	"errors"

	"github.com/insolar/insolar/configuration"
)

// Factory interface provides methods for creating stream or datagram transports
type Factory interface {
	CreateStreamTransport(StreamHandler) (StreamTransport, error)
	CreateDatagramTransport(DatagramHandler) (DatagramTransport, error)
}

// NewFactory constructor creates new transport factory
func NewFactory(cfg configuration.Transport) Factory {
	return &factory{cfg: cfg}
}

type factory struct {
	cfg configuration.Transport
}

// CreateStreamTransport creates new TCP transport
func (f *factory) CreateStreamTransport(handler StreamHandler) (StreamTransport, error) {
	switch f.cfg.Protocol {
	case "TCP":
		return newTCPTransport(f.cfg.Address, f.cfg.FixedPublicAddress, handler), nil
	default:
		return nil, errors.New("invalid transport configuration")
	}
}

// CreateDatagramTransport creates new UDP transport
func (f *factory) CreateDatagramTransport(handler DatagramHandler) (DatagramTransport, error) {
	return newUDPTransport(f.cfg.Address, f.cfg.FixedPublicAddress, handler), nil
}
