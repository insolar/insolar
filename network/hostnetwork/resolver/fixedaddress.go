package resolver

import (
	"fmt"
	"net"
	"net/url"

	"github.com/pkg/errors"
)

type fixedAddressResolver struct {
	publicAddress string
}

func NewFixedAddressResolver(publicAddress string) PublicAddressResolver {
	return newFixedAddressResolver(publicAddress)
}

func newFixedAddressResolver(publicAddress string) *fixedAddressResolver {
	return &fixedAddressResolver{
		publicAddress: publicAddress,
	}
}

func (r *fixedAddressResolver) Resolve(address string) (string, error) {
	url, err := url.Parse(address)

	var port string
	if err != nil {
		_, port, _ = net.SplitHostPort(address)
	} else {
		port = url.Port()
	}

	if port == "" {
		return "", errors.New("Failed to extract port from uri: " + address)
	}
	return fmt.Sprintf("%s:%s", r.publicAddress, port), nil
}
