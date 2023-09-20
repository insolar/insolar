package resolver

type exactResolver struct {
}

// NewExactResolver returns new no-op resolver.
func NewExactResolver() PublicAddressResolver {
	return newExactResolver()
}

func newExactResolver() *exactResolver {
	return &exactResolver{}
}

// Resolve returns host's current network address.
func (er *exactResolver) Resolve(address string) (string, error) {
	return address, nil
}
